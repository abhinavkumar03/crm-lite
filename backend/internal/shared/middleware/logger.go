package middleware

import (
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"github.com/abhinavkumar03/crm-lite/backend/internal/shared/constants"
	sharedlogger "github.com/abhinavkumar03/crm-lite/backend/internal/shared/logger"
)

func Logger() gin.HandlerFunc {

	logger := sharedlogger.New()

	return func(c *gin.Context) {

		start := time.Now()

		c.Next()

		requestID, _ := c.Get(constants.ContextRequestID)

		logger.Info(
			"request completed",
			zap.String("request_id", requestID.(string)),
			zap.String("method", c.Request.Method),
			zap.String("path", c.Request.URL.Path),
			zap.Int("status", c.Writer.Status()),
			zap.Duration("duration", time.Since(start)),
			zap.String("client_ip", c.ClientIP()),
		)
	}
}
