package handler

import (
	"github.com/abhinavkumar03/crm-lite/backend/internal/shared/response"
	"github.com/gin-gonic/gin"
)

func HealthCheck(c *gin.Context) {

	response.OK(
		c,
		"Service is healthy",
		gin.H{
			"service": "crm-lite",
			"status":  "UP",
		},
	)
}
