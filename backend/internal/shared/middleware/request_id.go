package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/abhinavkumar03/crm-lite/backend/internal/shared/constants"
)

func RequestID() gin.HandlerFunc {

	return func(c *gin.Context) {

		requestID := uuid.NewString()

		c.Set(constants.ContextRequestID, requestID)

		c.Writer.Header().Set(
			"X-Request-ID",
			requestID,
		)

		c.Next()
	}
}
