package app

import (
	"time"

	"github.com/singhAmandeep007/ReMinder/backend/gin-server/pkg/auth"
	"github.com/singhAmandeep007/ReMinder/backend/gin-server/pkg/config"
	"github.com/singhAmandeep007/ReMinder/backend/gin-server/pkg/logger"

	"github.com/singhAmandeep007/ReMinder/backend/gin-server/server/internal/api/handler"
	"github.com/singhAmandeep007/ReMinder/backend/gin-server/server/internal/api/middleware"
	"github.com/singhAmandeep007/ReMinder/backend/gin-server/server/internal/db"
	"github.com/singhAmandeep007/ReMinder/backend/gin-server/server/internal/repository"
	"github.com/singhAmandeep007/ReMinder/backend/gin-server/server/internal/service"
)

// Container holds all dependencies
type Container struct {
	// Core components
	DBManager *db.DBManager
	Log       *logger.Logger
	Cfg       *config.Config

	// Middlewares
	Middleware middleware.Middleware

	// Repositories
	UserRepository     repository.UserRepository
	ReminderRepository repository.ReminderRepository

	// Services
	AuthService service.AuthService

	// Handlers
	AuthHandler handler.AuthHandler
}

// NewContainer creates a new dependency container
func NewContainer(cfg *config.Config, log *logger.Logger, dbManager *db.DBManager) *Container {
	c := &Container{
		DBManager: dbManager,
		Log:       log,
		Cfg:       cfg,
	}

	// Initialize JWT manager with configuration
	config := auth.DefaultConfig()
	config.AccessSecret = "your-access-secret-key" // Use strong, environment-based secrets in production
	config.RefreshSecret = "your-refresh-secret-key"
	config.AccessTokenDuration = 15 * time.Minute
	config.RefreshTokenDuration = 7 * 24 * time.Hour
	config.IdentityKey = "user" // Key to store user ID in claims

	authManager := auth.NewAuthManager(config)

	// Initialize repositories
	c.UserRepository = repository.NewUserRepository(dbManager)

	// Initialize services
	c.AuthService = service.NewAuthService(c.UserRepository, log, authManager)

	// Initialize handlers
	c.AuthHandler = handler.NewAuthHandler(c.AuthService, authManager, log)

	// Initialize middlewares
	c.Middleware = middleware.NewMiddleware(log, authManager)

	return c
}
