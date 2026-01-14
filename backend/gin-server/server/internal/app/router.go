package app

import (
	"net/http"
	"os"
	"path/filepath"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func NewRouter(container *Container) *gin.Engine {
	r := gin.Default()

	r.Use(cors.Default())

	// Apply global middlewares
	r.Use(container.Middleware.Logger())
	r.Use(container.Middleware.Recovery())

	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	// OpenAPI spec endpoint
	r.GET("/api/openapi.yaml", func(c *gin.Context) {
		// Try to find the OpenAPI spec file
		paths := []string{
			"../../../openapi/api.yaml",    // From server dir
			"../../openapi/api.yaml",       // From gin-server dir
			"./openapi/api.yaml",           // From root
			os.Getenv("OPENAPI_SPEC_PATH"), // Environment variable
		}

		for _, p := range paths {
			if p == "" {
				continue
			}
			absPath, err := filepath.Abs(p)
			if err != nil {
				continue
			}
			if _, err := os.Stat(absPath); err == nil {
				c.File(absPath)
				return
			}
		}

		c.JSON(http.StatusNotFound, gin.H{"message": "OpenAPI spec not found"})
	})

	// Swagger UI static files (if swagger-ui-dist is available)
	r.StaticFS("/swagger", http.Dir("./swagger-ui"))

	api := r.Group("/api")
	{
		// Public routes
		auth := api.Group("/auth")
		{
			auth.POST("/register", container.AuthHandler.Register)
			auth.POST("/login", container.AuthHandler.Login)
			auth.POST("/refresh", container.AuthHandler.RefreshToken)
		}

		// Protected routes
		protected := api.Group("/")
		protected.Use(container.Middleware.Authenticate())
		{
			// User routes
			users := protected.Group("/users")
			{
				users.GET("/me", container.AuthHandler.GetMe)
			}

			// 	reminders := protected.Group("/reminders")
			// 	{
			// 		reminders.POST("", reminderHandler.CreateReminder)
			// 		reminders.GET("", reminderHandler.ListReminders)
			// 		reminders.GET("/:id", reminderHandler.GetReminder)
			// 		reminders.PUT("/:id", reminderHandler.UpdateReminder)
			// 		reminders.DELETE("/:id", reminderHandler.DeleteReminder)

			// }
		}
	}

	return r
}
