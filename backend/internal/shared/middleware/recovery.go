package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/abhinavkumar03/crm-lite/backend/internal/shared/response"
)

func Recovery() gin.HandlerFunc {

	return gin.CustomRecovery(func(c *gin.Context, recovered interface{}) {

		response.InternalServerError(
			c,
			http.StatusText(http.StatusInternalServerError),
		)
	})
}
