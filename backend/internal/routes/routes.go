package routes

import (
	"github.com/gin-gonic/gin"

	"github.com/abhinavkumar03/crm-lite/backend/internal/handler"
)

func Register(router *gin.Engine) {
	api := router.Group("/api/v1")

	{
		api.GET("/health", handler.HealthCheck)
	}
}
