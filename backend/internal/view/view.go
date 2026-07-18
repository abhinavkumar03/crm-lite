package view

import (
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/abhinavkumar03/crm-lite/backend/internal/rbac"
	"github.com/abhinavkumar03/crm-lite/backend/internal/view/handler"
	"github.com/abhinavkumar03/crm-lite/backend/internal/view/repository"
	"github.com/abhinavkumar03/crm-lite/backend/internal/view/service"
)

// Module is the saved-views engine composition root.
type Module struct {
	Handler *handler.ViewHandler
	auth    gin.HandlerFunc
	org     gin.HandlerFunc
	load    gin.HandlerFunc
	guard   *rbac.Guard
}

func NewModule(db *pgxpool.Pool, auth, org, load gin.HandlerFunc, guard *rbac.Guard) *Module {
	repo := repository.New(db)
	svc := service.New(repo)
	h := handler.New(svc)

	return &Module{
		Handler: h,
		auth:    auth,
		org:     org,
		load:    load,
		guard:   guard,
	}
}

// RegisterRoutes mounts saved views under a module. Viewing/saving views
// requires record.view plus module-level view access.
func (m *Module) RegisterRoutes(api *gin.RouterGroup) {
	views := api.Group("/modules/:id/views")
	views.Use(m.auth, m.org, m.load, m.guard.RequireModule(rbac.ActionView))

	views.GET("", m.Handler.List)
	views.POST("", m.Handler.Create)
	views.GET("/:viewId", m.Handler.GetByID)
	views.PUT("/:viewId", m.Handler.Update)
	views.DELETE("/:viewId", m.Handler.Delete)
	views.POST("/:viewId/default", m.Handler.SetDefault)
}
