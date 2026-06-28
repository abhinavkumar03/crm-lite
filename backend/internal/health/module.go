package health

import (
	"github.com/abhinavkumar03/crm-lite/backend/internal/health/handler"
	"github.com/gin-gonic/gin"
)

type Module struct {
	Handler *handler.HealthHandler
}

func NewModule() *Module {
	return &Module{
		Handler: handler.New(),
	}
}

func (m *Module) RegisterRoutes(api *gin.RouterGroup) {
	api.GET("/health", m.Handler.Health)
}
