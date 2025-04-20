package middleware

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/singhAmandeep007/ReMinder/backend/gin-server/pkg/logger"
)

type loggerMiddleware struct {
	log *logger.Logger
}

func NewLoggerMiddleware(log *logger.Logger) LoggerMiddleware {
	return &loggerMiddleware{
		log: log,
	}
}

func (m *loggerMiddleware) Logger() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		query := c.Request.URL.RawQuery

		c.Next()

		latency := time.Since(start)
		clientIP := c.ClientIP()
		method := c.Request.Method
		statusCode := c.Writer.Status()
		errorMessage := c.Errors.ByType(gin.ErrorTypePrivate).String()

		if query != "" {
			path = path + "?" + query
		}

		m.log.Infof("[GIN] %s | %3d | %13v | %15s | %-7s %s %s",
			start.Format("2006/01/02 - 15:04:05"),
			statusCode,
			latency,
			clientIP,
			method,
			path,
			errorMessage,
		)
	}
}
