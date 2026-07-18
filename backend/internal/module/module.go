package module

import (
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/abhinavkumar03/crm-lite/backend/internal/module/handler"
	"github.com/abhinavkumar03/crm-lite/backend/internal/module/repository"
	"github.com/abhinavkumar03/crm-lite/backend/internal/module/service"
)

// Module is the dynamic-module engine composition root.
type Module struct {
	Handler *handler.ModuleHandler
	auth    gin.HandlerFunc
	org     gin.HandlerFunc
}

// NewModule wires the module engine. It takes the auth middleware and the
// tenant (organization-scoping) middleware, both applied to every route.
func NewModule(db *pgxpool.Pool, auth gin.HandlerFunc, org gin.HandlerFunc) *Module {
	repo := repository.New(db)
	svc := service.New(repo)
	h := handler.New(svc)

	return &Module{
		Handler: h,
		auth:    auth,
		org:     org,
	}
}

func (m *Module) RegisterRoutes(api *gin.RouterGroup) {
	modules := api.Group("/modules")
	modules.Use(m.auth, m.org)

	modules.GET("", m.Handler.List)
	modules.POST("", m.Handler.Create)
	modules.POST("/reorder", m.Handler.Reorder)
	modules.GET("/:id", m.Handler.GetByID)
	modules.PUT("/:id", m.Handler.Update)
	modules.PATCH("/:id/status", m.Handler.SetStatus)
	modules.DELETE("/:id", m.Handler.Delete)

	// Navigation is exposed separately to avoid a static/param route conflict
	// under /modules (GET /modules/navigation vs GET /modules/:id).
	nav := api.Group("/navigation")
	nav.Use(m.auth, m.org)
	nav.GET("", m.Handler.Navigation)
}
