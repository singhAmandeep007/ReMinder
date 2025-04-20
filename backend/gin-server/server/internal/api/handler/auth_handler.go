package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/singhAmandeep007/ReMinder/backend/gin-server/pkg/auth"
	"github.com/singhAmandeep007/ReMinder/backend/gin-server/pkg/logger"
	"github.com/singhAmandeep007/ReMinder/backend/gin-server/server/internal/api/dto/request"
	"github.com/singhAmandeep007/ReMinder/backend/gin-server/server/internal/service"
	"github.com/singhAmandeep007/ReMinder/backend/gin-server/server/internal/utils"
)

type AuthHandler interface {
	Register(c *gin.Context)
	Login(c *gin.Context)
	RefreshToken(c *gin.Context)
	GetMe(c *gin.Context)
}

// AuthHandler handles authentication-related requests
type authHandler struct {
	authService service.AuthService
	authManager *auth.AuthManager
	log         *logger.Logger
}

// NewAuthHandler creates a new AuthHandler instance
func NewAuthHandler(authService service.AuthService, authManager *auth.AuthManager, log *logger.Logger) AuthHandler {
	return &authHandler{
		log:         log,
		authService: authService,
		authManager: authManager,
	}
}

// Register handles user registration
func (h *authHandler) Register(c *gin.Context) {
	var req request.RegisterRequest

	// Bind the request body to the registration struct
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid request data")
		return
	}

	user, err := h.authService.Register(c.Request.Context(), req.Email, req.Password)

	if err != nil {
		h.log.Warnf("Registration failed for email: %s, error: %v", req.Email, err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	h.log.Infof("User registered userID: %s", user.ID)
	c.JSON(http.StatusCreated, gin.H{"message": "User registered successfully", "userID": user.ID})
}

// Login handles user login
func (h *authHandler) Login(c *gin.Context) {
	var req request.LoginRequest

	// Bind the request body to the login struct
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request data"})
		return
	}

	accessToken, refreshToken, err := h.authService.Login(c.Request.Context(), req.Email, req.Password)
	if err != nil {
		h.log.Warnf("Login failed for email: %s, error: %v", req.Email, err)
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	// Return the tokens
	utils.SuccessResponse(c, http.StatusOK, gin.H{
		"accessToken":  accessToken,
		"refreshToken": refreshToken,
	})
}

// RefreshToken handles token refresh
func (h *authHandler) RefreshToken(c *gin.Context) {
	var req request.RefreshTokenRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request data"})
		return
	}

	accessToken, refreshToken, err := h.authService.Refresh(c.Request.Context(), req.RefreshToken)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired refresh token"})
		return
	}

	utils.SuccessResponse(c, http.StatusOK, gin.H{
		"accessToken":  accessToken,
		"refreshToken": refreshToken,
	})
}

func (h *authHandler) GetMe(c *gin.Context) {
	claims, exists := utils.GetClaimsFromGinContext(c, h.authManager)
	if !exists {
		utils.ErrorResponse(c, http.StatusUnauthorized, "unauthorized")
		return
	}

	user, err := h.authService.GetMe(c.Request.Context(), claims.EntityID)
	if err != nil {
		h.log.Warnf("GetMe failed for userID: %s, error: %v", claims.EntityID, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve user information"})
		return
	}

	utils.SuccessResponse(c, http.StatusOK, user)
}
