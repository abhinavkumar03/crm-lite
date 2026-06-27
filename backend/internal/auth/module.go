package auth

import (
	"github.com/gin-gonic/gin"

	"github.com/abhinavkumar03/crm-lite/backend/internal/auth/handler"
)

type Module struct {
	Handler *handler.AuthHandler
}

func NewModule(
	handler *handler.AuthHandler,
) *Module {

	return &Module{
		Handler: handler,
	}
}

func (m *Module) RegisterRoutes(api *gin.RouterGroup) {
	auth := api.Group("/auth")

	auth.POST("/register", m.Handler.Register)
	auth.POST("/login", m.Handler.Login)
	auth.GET("/profile", m.Handler.Profile)
}
