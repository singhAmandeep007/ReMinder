package middleware

import (
	"time"

	"github.com/gin-gonic/gin"

	"github.com/singhAmandeep007/ReMinder/backend/gin-server/pkg/auth"
	"github.com/singhAmandeep007/ReMinder/backend/gin-server/pkg/logger"
)

type Middleware interface {
	Authenticate() gin.HandlerFunc
	Authorize(roles ...string) gin.HandlerFunc
	Logger() gin.HandlerFunc
	Recovery() gin.HandlerFunc
	RateLimiter() gin.HandlerFunc
}

type AuthMiddleware interface {
	Authenticate() gin.HandlerFunc
	Authorize(roles ...string) gin.HandlerFunc
}

type LoggerMiddleware interface {
	Logger() gin.HandlerFunc
}

type RateLimiterMiddleware interface {
	RateLimiter() gin.HandlerFunc
}

type RecoveryMiddleware interface {
	Recovery() gin.HandlerFunc
}

type middleware struct {
	authMiddleware
	loggerMiddleware
	recoveryMiddleware
	rateLimiterMiddleware
}

func NewMiddleware(log *logger.Logger, authManager *auth.AuthManager) Middleware {
	return &middleware{
		authMiddleware:        authMiddleware{log: log, authManager: authManager},
		loggerMiddleware:      loggerMiddleware{log: log},
		recoveryMiddleware:    recoveryMiddleware{log: log},
		rateLimiterMiddleware: rateLimiterMiddleware{log: log, limit: 100, window: 1 * time.Minute},
	}
}
