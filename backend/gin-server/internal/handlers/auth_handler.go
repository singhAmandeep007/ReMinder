package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"gin-server/internal/logger"
	"gin-server/internal/services"
	"gin-server/internal/utils"
)

type AuthHandler struct {
	authService services.AuthService
	logger      logger.Logger
}

func NewAuthHandler(authService services.AuthService, logger logger.Logger) *AuthHandler {
	return &AuthHandler{authService: authService, logger: logger}
}

type RegisterRequest struct {
	Email    string `json:"email" binding:"required"`
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func (h *AuthHandler) Register(c *gin.Context) {
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Warnf("Invalid register request: %v", err)
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	user, err := h.authService.Register(c.Request.Context(), req.Email, req.Username, req.Password)
	if err != nil {
		h.logger.Warnf("Registration failed for username: %s, error: %v", req.Username, err)
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	h.logger.Infof("User registered: %s, userID: %s", req.Username, user.ID)
	utils.SuccessResponse(c, http.StatusCreated, gin.H{"message": "User registered successfully", "userID": user.ID})
}
