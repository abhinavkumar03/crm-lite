package module

import (
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/abhinavkumar03/crm-lite/backend/internal/module/handler"
	"github.com/abhinavkumar03/crm-lite/backend/internal/module/repository"
	"github.com/abhinavkumar03/crm-lite/backend/internal/module/service"
	"github.com/abhinavkumar03/crm-lite/backend/internal/rbac"
)

// Module is the dynamic-module engine composition root.
type Module struct {
	Handler *handler.ModuleHandler
	auth    gin.HandlerFunc
	org     gin.HandlerFunc
	load    gin.HandlerFunc
	guard   *rbac.Guard
}

// NewModule wires the module engine. auth → org → rbac.Load run on every route;
// mutating endpoints additionally Require module.manage.
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

func (m *Module) RegisterRoutes(api *gin.RouterGroup) {
	modules := api.Group("/modules")
	modules.Use(m.auth, m.org, m.load)

	modules.GET("", m.guard.Require(rbac.PermModuleView), m.Handler.List)
	modules.GET("/:id", m.guard.Require(rbac.PermModuleView), m.Handler.GetByID)

	modules.POST("", m.guard.Require(rbac.PermModuleManage), m.Handler.Create)
	modules.POST("/reorder", m.guard.Require(rbac.PermModuleManage), m.Handler.Reorder)
	modules.PUT("/:id", m.guard.Require(rbac.PermModuleManage), m.Handler.Update)
	modules.PATCH("/:id/status", m.guard.Require(rbac.PermModuleManage), m.Handler.SetStatus)
	modules.DELETE("/:id", m.guard.Require(rbac.PermModuleManage), m.Handler.Delete)

	nav := api.Group("/navigation")
	nav.Use(m.auth, m.org, m.load)
	nav.GET("", m.guard.Require(rbac.PermModuleView), m.Handler.Navigation)
}
