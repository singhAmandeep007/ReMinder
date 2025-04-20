package auth

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func init() {
	// Set Gin to test mode
	gin.SetMode(gin.TestMode)
}

// Helper to create a JWT manager for testing
func createTestManager() *AuthManager {
	config := DefaultConfig()
	config.AccessSecret = "test-access-secret"
	config.RefreshSecret = "test-refresh-secret"
	config.AccessTokenDuration = 5 * time.Minute
	config.RefreshTokenDuration = 24 * time.Hour

	return NewAuthManager(config)
}

func TestGenerateAndParseToken(t *testing.T) {
	manager := createTestManager()

	// Test data
	entityID := "123"

	custom := map[string]interface{}{
		"username": "testuser",
		"email":    "test@example.com",
		"role":     "user",
	}

	// Test cases
	tests := []struct {
		name      string
		tokenType TokenType
		wantErr   bool
	}{
		{"Access Token", AccessToken, false},
		{"Refresh Token", RefreshToken, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Generate token
			token, err := manager.GenerateToken(entityID, tt.tokenType, custom)
			assert.NoError(t, err)
			assert.NotEmpty(t, token)

			// Parse token
			claims, err := manager.ParseToken(token, tt.tokenType)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, entityID, claims.EntityID)
				assert.Equal(t, tt.tokenType, claims.TokenType)

				assert.Equal(t, "user", claims.Custom["role"])
				assert.Equal(t, "testuser", claims.Custom["username"])
				assert.Equal(t, "test@example.com", claims.Custom["email"])
			}
		})
	}
}

func TestInvalidToken(t *testing.T) {
	manager := createTestManager()

	// Test cases
	tests := []struct {
		name      string
		token     string
		tokenType TokenType
		wantErr   error
	}{
		{"Empty Token", "", AccessToken, ErrMissingToken},
		{"Invalid Format", "not.a.jwt.token", AccessToken, ErrInvalidToken},
		{"Wrong Token Type", "", AccessToken, ErrInvalidToken}, // Will be set in setup
	}

	// Generate a refresh token to test with access token type
	refreshToken, _ := manager.GenerateToken("123", RefreshToken, nil)

	for i, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			token := tt.token
			if i == 2 {
				token = refreshToken
			}

			_, err := manager.ParseToken(token, tt.tokenType)
			if tt.wantErr == ErrMissingToken {
				assert.Error(t, err)
			} else {
				assert.Error(t, err)
			}
		})
	}
}

func TestExpiredToken(t *testing.T) {
	// Create manager with very short duration
	config := DefaultConfig()
	config.AccessSecret = "test-access-secret"
	config.AccessTokenDuration = 1 * time.Millisecond
	manager := NewAuthManager(config)

	// Generate token
	token, err := manager.GenerateToken("123", AccessToken, nil)
	assert.NoError(t, err)

	// Wait for token to expire
	time.Sleep(10 * time.Millisecond)

	// Try to parse expired token
	_, err = manager.ParseToken(token, AccessToken)
	assert.Equal(t, ErrExpiredToken, err)
}

func TestTokenPair(t *testing.T) {
	manager := createTestManager()

	// Generate token pair
	accessToken, refreshToken, err := manager.GenerateTokenPair("123", nil)

	// Validate results
	assert.NoError(t, err)
	assert.NotEmpty(t, accessToken)
	assert.NotEmpty(t, refreshToken)

	// Parse both tokens
	accessClaims, err1 := manager.ParseToken(accessToken, AccessToken)
	refreshClaims, err2 := manager.ParseToken(refreshToken, RefreshToken)

	assert.NoError(t, err1)
	assert.NoError(t, err2)
	assert.Equal(t, "123", accessClaims.EntityID)
	assert.Equal(t, "123", refreshClaims.EntityID)
	assert.Equal(t, AccessToken, accessClaims.TokenType)
	assert.Equal(t, RefreshToken, refreshClaims.TokenType)
}

func TestRefreshTokens(t *testing.T) {
	manager := createTestManager()

	// Generate initial tokens
	_, refreshToken, err := manager.GenerateTokenPair("123", nil)
	assert.NoError(t, err)

	// Refresh tokens
	newAccess, newRefresh, err := manager.RefreshTokens(refreshToken)

	// Validate results
	assert.NoError(t, err)
	assert.NotEmpty(t, newAccess)
	assert.NotEmpty(t, newRefresh)

	// Verify new tokens are valid
	accessClaims, err1 := manager.ParseToken(newAccess, AccessToken)
	refreshClaims, err2 := manager.ParseToken(newRefresh, RefreshToken)

	assert.NoError(t, err1)
	assert.NoError(t, err2)
	assert.Equal(t, "123", accessClaims.EntityID)
	assert.Equal(t, "123", refreshClaims.EntityID)
	assert.Equal(t, AccessToken, accessClaims.TokenType)
	assert.Equal(t, RefreshToken, refreshClaims.TokenType)
	assert.Equal(t, accessClaims.Custom, refreshClaims.Custom)

}

func TestBlacklistToken(t *testing.T) {
	manager := createTestManager()

	// Generate token
	token, err := manager.GenerateToken("123", AccessToken, nil)
	assert.NoError(t, err)

	// Parse token to get claims
	claims, err := manager.ParseToken(token, AccessToken)
	assert.NoError(t, err)

	// Blacklist the token
	err = manager.BlacklistToken(claims.TokenID, time.Minute)
	assert.NoError(t, err)

	// Try to use the blacklisted token
	_, err = manager.ParseToken(token, AccessToken)
	assert.Error(t, err)
	assert.Equal(t, ErrInvalidToken, err)
}

func TestExtractTokenFromRequest(t *testing.T) {
	manager := createTestManager()

	// Test cases
	tests := []struct {
		name        string
		setupFunc   func() *http.Request
		wantErr     bool
		tokenConfig func()
	}{
		{
			name: "Valid Header Token",
			setupFunc: func() *http.Request {
				req := httptest.NewRequest("GET", "/", nil)
				req.Header.Set("Authorization", "Bearer test-token")
				return req
			},
			wantErr: false,
		},
		{
			name: "Missing Header Token",
			setupFunc: func() *http.Request {
				req := httptest.NewRequest("GET", "/", nil)
				return req
			},
			wantErr: true,
		},
		{
			name: "Cookie Token",
			setupFunc: func() *http.Request {
				req := httptest.NewRequest("GET", "/", nil)
				req.AddCookie(&http.Cookie{
					Name:  "jwt_access_token",
					Value: "test-token",
				})
				return req
			},
			wantErr: false,
			tokenConfig: func() {
				manager.Config.TokenLookup = "cookie:jwt_access_token"
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Apply token config if needed
			if tt.tokenConfig != nil {
				originalConfig := manager.Config
				tt.tokenConfig()
				defer func() { manager.Config = originalConfig }()
			}

			req := tt.setupFunc()
			token, err := manager.ExtractTokenFromRequest(req)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, "test-token", token)
			}
		})
	}
}

func TestTokenCookies(t *testing.T) {
	manager := createTestManager()
	manager.Config.SendCookies = true

	// Test setting cookies
	w := httptest.NewRecorder()
	manager.SetTokenCookies(w, "access-token", "refresh-token")

	// Get cookies from response
	cookies := w.Result().Cookies()

	// Verify cookies were set
	var foundAccess, foundRefresh bool
	for _, cookie := range cookies {
		if cookie.Name == manager.Config.AccessCookieName {
			assert.Equal(t, "access-token", cookie.Value)
			assert.True(t, cookie.HttpOnly)
			foundAccess = true
		}
		if cookie.Name == manager.Config.RefreshCookieName {
			assert.Equal(t, "refresh-token", cookie.Value)
			assert.True(t, cookie.HttpOnly)
			foundRefresh = true
		}
	}

	assert.True(t, foundAccess, "Access token cookie not found")
	assert.True(t, foundRefresh, "Refresh token cookie not found")

	// Test clearing cookies
	w = httptest.NewRecorder()
	manager.ClearTokenCookies(w)

	// Get cookies from response
	cookies = w.Result().Cookies()

	// Verify cookies were cleared
	for _, cookie := range cookies {
		if cookie.Name == manager.Config.AccessCookieName || cookie.Name == manager.Config.RefreshCookieName {
			assert.Equal(t, "", cookie.Value)
			assert.True(t, cookie.MaxAge < 0)
		}
	}
}

func TestIsAuthorized(t *testing.T) {
	manager := createTestManager()

	// Create test claims with roles
	claims := &CustomClaims{
		EntityID: "123",
		Custom: map[string]interface{}{
			"roles": []string{"user", "editor"},
		},
	}

	// Test cases
	tests := []struct {
		name          string
		requiredRoles []string
		expected      bool
	}{
		{"No Roles Required", []string{}, true},
		{"Has Required Role", []string{"user"}, true},
		{"Has One of Required Roles", []string{"admin", "editor"}, true},
		{"Missing Required Role", []string{"admin"}, false},
		{"Missing All Required Roles", []string{"admin", "manager"}, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := manager.IsAuthorized(claims, "roles", tt.requiredRoles)
			assert.Equal(t, tt.expected, result)
		})
	}
}
