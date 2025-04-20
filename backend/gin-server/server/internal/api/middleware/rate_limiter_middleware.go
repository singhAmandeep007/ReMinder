package middleware

import (
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/singhAmandeep007/ReMinder/backend/gin-server/pkg/logger"
	"github.com/singhAmandeep007/ReMinder/backend/gin-server/server/internal/utils"
)

type rateLimiterMiddleware struct {
	log     *logger.Logger
	limit   int
	window  time.Duration
	clients map[string][]time.Time
	mu      sync.Mutex
}

func NewRateLimiterMiddleware(log *logger.Logger, limit int, window time.Duration) RateLimiterMiddleware {
	return &rateLimiterMiddleware{
		log:     log,
		limit:   limit,
		window:  window,
		clients: make(map[string][]time.Time),
	}
}

func (m *rateLimiterMiddleware) RateLimiter() gin.HandlerFunc {
	return func(c *gin.Context) {
		clientIP := c.ClientIP()

		m.mu.Lock()
		defer m.mu.Unlock()

		now := time.Now()

		// Create new entry for this client if it doesn't exist
		if _, exists := m.clients[clientIP]; !exists {
			m.clients[clientIP] = []time.Time{}
		}

		// Remove timestamps outside the current window
		var validRequests []time.Time
		for _, timestamp := range m.clients[clientIP] {
			if now.Sub(timestamp) <= m.window {
				validRequests = append(validRequests, timestamp)
			}
		}

		m.clients[clientIP] = validRequests

		// Check if the client has exceeded the limit
		if len(m.clients[clientIP]) >= m.limit {
			utils.ErrorResponseWithAbort(c, http.StatusTooManyRequests, "Rate limit exceeded")
			return
		}

		// Add current request timestamp
		m.clients[clientIP] = append(m.clients[clientIP], now)

		c.Next()
	}
}
