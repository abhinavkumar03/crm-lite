package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"github.com/abhinavkumar03/crm-lite/backend/internal/shared/response"
)

func Recovery(logger *zap.Logger) gin.HandlerFunc {
	return gin.CustomRecovery(func(c *gin.Context, recovered interface{}) {
		logger.Error("panic recovered", zap.Any("panic", recovered))
		response.InternalServerError(c, http.StatusText(http.StatusInternalServerError))
	})
}
