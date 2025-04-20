package utils

import (
	"github.com/gin-gonic/gin"
	"github.com/singhAmandeep007/ReMinder/backend/gin-server/pkg/auth"
)

// GetClaimsFromGinContext extracts claims from the Gin context
func GetClaimsFromGinContext(c *gin.Context, authManager *auth.AuthManager) (customClaims *auth.CustomClaims, exists bool) {
	claims, exists := c.Get(authManager.Config.IdentityKey)
	if !exists {
		return nil, false
	}

	customClaims, ok := claims.(*auth.CustomClaims)
	return customClaims, ok
}
