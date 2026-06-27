package auth

import (
	"github.com/gin-gonic/gin"

	"github.com/abhinavkumar03/crm-lite/backend/internal/auth/handler"
)

func RegisterRoutes(
	router *gin.RouterGroup,
) {

	h := handler.New()

	auth := router.Group("/auth")

	{
		auth.POST("/register", h.Register)

		auth.POST("/login", h.Login)

		auth.GET("/profile", h.Profile)
	}
}
