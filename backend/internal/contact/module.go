package contact

import (
	"github.com/gin-gonic/gin"
)

type Module struct {
	// We'll inject dependencies in the next phase.
}

func NewModule( /* dependencies */ ) *Module {
	return &Module{}
}

func (m *Module) RegisterRoutes(router *gin.RouterGroup) {
	// Routes will be added in the CRUD phase.
}
