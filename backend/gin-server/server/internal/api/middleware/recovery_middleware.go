package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/singhAmandeep007/ReMinder/backend/gin-server/pkg/logger"
	"github.com/singhAmandeep007/ReMinder/backend/gin-server/server/internal/utils"
)

type recoveryMiddleware struct {
	log *logger.Logger
}

func NewRecoveryMiddleware(log *logger.Logger) RecoveryMiddleware {
	return &recoveryMiddleware{
		log: log,
	}
}

func (m *recoveryMiddleware) Recovery() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				m.log.Errorf("Panic recovered: %v", err)
				utils.ErrorResponseWithAbort(c, http.StatusInternalServerError, "Internal server error")
			}
		}()

		c.Next()
	}
}
