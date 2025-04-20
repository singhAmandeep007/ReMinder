package app

import (
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
