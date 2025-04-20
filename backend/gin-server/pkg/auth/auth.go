package auth

import (
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"

	"github.com/singhAmandeep007/ReMinder/backend/gin-server/pkg/memcache"
)

// Common errors
var (
	ErrMissingToken     = errors.New("missing auth token")
	ErrInvalidToken     = errors.New("invalid auth token")
	ErrExpiredToken     = errors.New("expired auth token")
	ErrInvalidSignature = errors.New("invalid token signature")
	ErrInvalidClaims    = errors.New("invalid token claims")
	ErrInvalidSubject   = errors.New("invalid token subject")
)

// TokenType defines different types of tokens
type TokenType string

const (
	// AccessToken is a short-lived token for API access
	AccessToken TokenType = "access"
	// RefreshToken is a longer-lived token for obtaining new access tokens
	RefreshToken TokenType = "refresh"
)

// Config holds the configuration for the JWT auth manager
type Config struct {
	AccessSecret          string         // Secret key for access tokens
	RefreshSecret         string         // Secret key for refresh tokens
	AccessTokenDuration   time.Duration  // Duration for access tokens
	RefreshTokenDuration  time.Duration  // Duration for refresh tokens
	TokenLookup           string         // Method to extract token: "header:Authorization" or "cookie:jwt"
	TokenHeadName         string         // Token prefix in header, e.g., "Bearer"
	AuthScheme            string         // Auth scheme in header, e.g., "Bearer"
	IdentityKey           string         // Key to store the identity in claims, e.g., "identity"
	DisableRefresh        bool           // Disable refresh token functionality if true
	SendCookies           bool           // Send tokens via cookies if true
	SecureCookies         bool           // Set Secure flag on cookies
	HTTPOnlyCookies       bool           // Set HTTPOnly flag on cookies
	CookieDomain          string         // Domain for cookies
	CookiePath            string         // Path for cookies
	CookieSameSite        http.SameSite  // SameSite attribute for cookies
	RefreshCookieName     string         // Name of refresh token cookie
	AccessCookieName      string         // Name of access token cookie
	BlacklistedTokenCache memcache.Cache // Cache for blacklisted tokens (optional)
}

// DefaultConfig returns a default configuration
func DefaultConfig() Config {
	return Config{
		AccessTokenDuration:   24 * time.Hour,
		RefreshTokenDuration:  7 * 24 * time.Hour,
		TokenLookup:           "header:Authorization",
		TokenHeadName:         "Bearer",
		AuthScheme:            "Bearer",
		IdentityKey:           "identity",
		DisableRefresh:        false,
		SendCookies:           false,
		SecureCookies:         true,
		HTTPOnlyCookies:       true,
		CookiePath:            "/",
		CookieSameSite:        http.SameSiteStrictMode,
		RefreshCookieName:     "jwt_refresh_token",
		AccessCookieName:      "jwt_access_token",
		BlacklistedTokenCache: memcache.NewInMemoryCache(24 * time.Hour),
	}
}

// CustomClaims extends jwt.RegisteredClaims with custom fields
type CustomClaims struct {
	jwt.RegisteredClaims
	TokenType TokenType              `json:"type"`          // Type of token (access/refresh)
	EntityID  string                 `json:"entityId "`     // EntityID identifier
	Custom    map[string]interface{} `json:"custom"`        // Custom user-defined claims
	TokenID   string                 `json:"jti,omitempty"` // Token ID for blacklisting
}

// AuthManager is the JWT authentication manager
type AuthManager struct {
	Config Config
}

// NewManager creates a new JWT authentication manager with the given configuration
func NewAuthManager(config Config) *AuthManager {
	return &AuthManager{
		Config: config,
	}
}

// GenerateTokenPair creates both access and refresh tokens
func (m *AuthManager) GenerateTokenPair(entityID string, custom map[string]interface{}) (string, string, error) {
	accessToken, err := m.GenerateToken(entityID, AccessToken, custom)
	if err != nil {
		return "", "", err
	}

	if m.Config.DisableRefresh {
		return accessToken, "", nil
	}

	refreshToken, err := m.GenerateToken(entityID, RefreshToken, custom)
	if err != nil {
		return "", "", err
	}

	return accessToken, refreshToken, nil
}

// GenerateToken creates a new token based on the provided information
func (m *AuthManager) GenerateToken(entityID string, tokenType TokenType, custom map[string]interface{}) (string, error) {
	now := time.Now()

	// Set appropriate duration and signing key based on token type
	var duration time.Duration
	var signingKey []byte

	if tokenType == AccessToken {
		duration = m.Config.AccessTokenDuration
		signingKey = []byte(m.Config.AccessSecret)
	} else {
		duration = m.Config.RefreshTokenDuration
		signingKey = []byte(m.Config.RefreshSecret)
	}

	// Generate a unique token ID
	tokenID := uuid.New().String()

	// Create the claims
	claims := CustomClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   entityID,
			ExpiresAt: jwt.NewNumericDate(now.Add(duration)),
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now),
			ID:        tokenID,
		},
		TokenType: tokenType,
		EntityID:  entityID,
		Custom:    custom,
		TokenID:   tokenID,
	}

	// Create the token with the claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Sign and get the complete encoded token as a string
	tokenString, err := token.SignedString(signingKey)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

// ParseToken parses and validates a token
func (m *AuthManager) ParseToken(tokenString string, tokenType TokenType) (*CustomClaims, error) {
	// Remove token prefix if it exists
	if strings.HasPrefix(tokenString, m.Config.TokenHeadName+" ") {
		tokenString = strings.TrimPrefix(tokenString, m.Config.TokenHeadName+" ")
	}

	// Select the appropriate secret key
	var secretKey []byte
	if tokenType == AccessToken {
		secretKey = []byte(m.Config.AccessSecret)
	} else {
		secretKey = []byte(m.Config.RefreshSecret)
	}

	// Parse the token
	token, err := jwt.ParseWithClaims(tokenString, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		// Validate the signing algorithm
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return secretKey, nil
	})

	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, ErrExpiredToken
		}
		return nil, ErrInvalidToken
	}

	// Check if token is valid
	if !token.Valid {
		return nil, ErrInvalidToken
	}

	// Cast claims
	claims, ok := token.Claims.(*CustomClaims)
	if !ok {
		return nil, ErrInvalidClaims
	}

	// Check if token type matches
	if claims.TokenType != tokenType {
		return nil, ErrInvalidToken
	}

	// Check if token is blacklisted
	if m.Config.BlacklistedTokenCache != nil {
		if _, blacklisted := m.Config.BlacklistedTokenCache.Get(claims.TokenID); blacklisted {
			return nil, ErrInvalidToken
		}
	}

	return claims, nil
}

func (m *AuthManager) RefreshTokens(refreshToken string) (string, string, error) {
	if m.Config.DisableRefresh {
		return "", "", errors.New("refresh functionality is disabled")
	}

	// Parse and validate the refresh token
	claims, err := m.ParseToken(refreshToken, RefreshToken)
	if err != nil {
		return "", "", err
	}

	// Generate new tokens
	accessToken, newRefreshToken, err := m.GenerateTokenPair(claims.EntityID, claims.Custom)
	if err != nil {
		return "", "", err
	}

	// Blacklist the old refresh token
	if m.Config.BlacklistedTokenCache != nil && claims.TokenID != "" {
		expiry := time.Until(claims.ExpiresAt.Time)
		m.Config.BlacklistedTokenCache.Set(claims.TokenID, true, expiry)
	}

	return accessToken, newRefreshToken, nil
}

// ExtractTokenFromRequest extracts the token from an HTTP request
func (m *AuthManager) ExtractTokenFromRequest(r *http.Request) (string, error) {
	parts := strings.Split(m.Config.TokenLookup, ":")
	if len(parts) != 2 {
		return "", errors.New("invalid token lookup config")
	}

	switch parts[0] {
	case "header":
		token := r.Header.Get(parts[1])
		if token == "" {
			return "", ErrMissingToken
		}

		// Handle authorization header schemes (Bearer, etc.)
		if m.Config.AuthScheme != "" {
			prefix := m.Config.AuthScheme + " "
			if !strings.HasPrefix(token, prefix) {
				return "", ErrInvalidToken
			}
			return strings.TrimPrefix(token, prefix), nil
		}

		return token, nil

	case "cookie":
		cookie, err := r.Cookie(parts[1])
		if err != nil {
			return "", ErrMissingToken
		}
		return cookie.Value, nil

	default:
		return "", errors.New("unsupported token lookup method")
	}
}

func (m *AuthManager) SetTokenCookies(w http.ResponseWriter, accessToken, refreshToken string) {
	if !m.Config.SendCookies {
		return
	}

	// Set access token cookie
	http.SetCookie(w, &http.Cookie{
		Name:     m.Config.AccessCookieName,
		Value:    accessToken,
		Path:     m.Config.CookiePath,
		Domain:   m.Config.CookieDomain,
		MaxAge:   int(m.Config.AccessTokenDuration.Seconds()),
		Secure:   m.Config.SecureCookies,
		HttpOnly: m.Config.HTTPOnlyCookies,
		SameSite: m.Config.CookieSameSite,
	})

	// Set refresh token cookie if refresh is enabled
	if !m.Config.DisableRefresh && refreshToken != "" {
		http.SetCookie(w, &http.Cookie{
			Name:     m.Config.RefreshCookieName,
			Value:    refreshToken,
			Path:     m.Config.CookiePath,
			Domain:   m.Config.CookieDomain,
			MaxAge:   int(m.Config.RefreshTokenDuration.Seconds()),
			Secure:   m.Config.SecureCookies,
			HttpOnly: m.Config.HTTPOnlyCookies,
			SameSite: m.Config.CookieSameSite,
		})
	}
}

// ClearTokenCookies removes token cookies
func (m *AuthManager) ClearTokenCookies(w http.ResponseWriter) {
	if !m.Config.SendCookies {
		return
	}

	// Clear access token cookie
	http.SetCookie(w, &http.Cookie{
		Name:     m.Config.AccessCookieName,
		Value:    "",
		Path:     m.Config.CookiePath,
		Domain:   m.Config.CookieDomain,
		MaxAge:   -1,
		Secure:   m.Config.SecureCookies,
		HttpOnly: m.Config.HTTPOnlyCookies,
		SameSite: m.Config.CookieSameSite,
	})

	// Clear refresh token cookie
	http.SetCookie(w, &http.Cookie{
		Name:     m.Config.RefreshCookieName,
		Value:    "",
		Path:     m.Config.CookiePath,
		Domain:   m.Config.CookieDomain,
		MaxAge:   -1,
		Secure:   m.Config.SecureCookies,
		HttpOnly: m.Config.HTTPOnlyCookies,
		SameSite: m.Config.CookieSameSite,
	})
}

// BlacklistToken adds a token to the blacklist
func (m *AuthManager) BlacklistToken(tokenID string, expiration time.Duration) error {
	if m.Config.BlacklistedTokenCache == nil {
		return errors.New("blacklist cache not configured")
	}

	return m.Config.BlacklistedTokenCache.Set(tokenID, true, expiration)
}

// IsAuthorized checks if the claims have the required roles
func (m *AuthManager) IsAuthorized(claims *CustomClaims, rolesKey string, requiredRoles []string) bool {
	if len(requiredRoles) == 0 {
		return true
	}

	// Convert user roles to a map for O(1) lookups
	userRoles := make(map[string]bool)
	for _, role := range claims.Custom[rolesKey].([]string) {
		userRoles[role] = true
	}

	// Check if user has any of the required roles
	for _, role := range requiredRoles {
		if userRoles[role] {
			return true
		}
	}

	return false
}
