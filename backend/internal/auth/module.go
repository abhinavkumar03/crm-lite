package auth

import (
	"github.com/gin-gonic/gin"

	"github.com/abhinavkumar03/crm-lite/backend/internal/auth/handler"
)

type Module struct {
	Handler *handler.Handler
}

func NewModule(
	handler *handler.Handler,
) *Module {

	return &Module{
		Handler: handler,
	}
}

func (m *Module) RegisterRoutes(
	router *gin.RouterGroup,
) {

	auth := router.Group("/auth")

	m.Handler.RegisterRoutes(auth)
}
