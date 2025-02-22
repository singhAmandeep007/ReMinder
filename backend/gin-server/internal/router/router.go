package router

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"

	"gin-server/internal/handlers"
)

func SetupRouter(authHandler *handlers.AuthHandler) *gin.Engine {
	r := gin.Default()

	r.Use(cors.Default())

	// Public routes
	authGroup := r.Group("/auth")
	{
		authGroup.POST("/register", authHandler.Register)
	}

	return r
}
