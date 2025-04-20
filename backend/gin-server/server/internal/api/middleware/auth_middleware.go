package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/singhAmandeep007/ReMinder/backend/gin-server/pkg/auth"
	"github.com/singhAmandeep007/ReMinder/backend/gin-server/pkg/logger"
	"github.com/singhAmandeep007/ReMinder/backend/gin-server/server/internal/utils"
)

const (
	UserRoleKey = "role"
)

type authMiddleware struct {
	log         *logger.Logger
	authManager *auth.AuthManager
}

func NewAuthMiddleware(log *logger.Logger, authManager *auth.AuthManager) AuthMiddleware {
	return &authMiddleware{
		log:         log,
		authManager: authManager,
	}
}

func (m *authMiddleware) Authenticate() gin.HandlerFunc {
	return func(c *gin.Context) {

		tokenString, err := m.authManager.ExtractTokenFromRequest(c.Request)
		if err != nil {
			utils.ErrorResponseWithAbort(c, http.StatusUnauthorized, "unauthorized")
			return
		}

		// Parse and validate token
		claims, err := m.authManager.ParseToken(tokenString, auth.AccessToken)
		if err != nil {
			if err == auth.ErrExpiredToken {
				utils.ErrorResponse(c, http.StatusUnauthorized, "token expired")
			} else {
				utils.ErrorResponse(c, http.StatusUnauthorized, "invalid token")
			}
			c.Abort()
			return
		}

		// Set claims in Gin context
		c.Set(m.authManager.Config.IdentityKey, claims)

		c.Next()
	}
}

func (m *authMiddleware) Authorize(roles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get claims from context
		claims, exists := utils.GetClaimsFromGinContext(c, m.authManager)
		if !exists {
			utils.ErrorResponseWithAbort(c, http.StatusUnauthorized, "unauthorized")
			return
		}

		// Check if user has required roles
		if !m.authManager.IsAuthorized(claims, UserRoleKey, roles) {
			utils.ErrorResponseWithAbort(c, http.StatusForbidden, "forbidden: insufficient permissions")
			return
		}

		c.Next()
	}
}
