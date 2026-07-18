package roles

import (
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/abhinavkumar03/crm-lite/backend/internal/rbac"
	"github.com/abhinavkumar03/crm-lite/backend/internal/roles/handler"
	"github.com/abhinavkumar03/crm-lite/backend/internal/roles/repository"
	"github.com/abhinavkumar03/crm-lite/backend/internal/roles/service"
)

// Module is the roles & permissions composition root. It exposes the permission
// catalog, role CRUD, the permission matrix, and module-/field-level ACL.
type Module struct {
	Handler *handler.RoleHandler
	auth    gin.HandlerFunc
	org     gin.HandlerFunc
	load    gin.HandlerFunc
	guard   *rbac.Guard
}

func NewModule(db *pgxpool.Pool, auth, org, load gin.HandlerFunc, guard *rbac.Guard) *Module {
	repo := repository.New(db)
	svc := service.New(repo, guard)
	h := handler.New(svc)

	return &Module{
		Handler: h,
		auth:    auth,
		org:     org,
		load:    load,
		guard:   guard,
	}
}

// RegisterRoutes mounts the RBAC API.
func (m *Module) RegisterRoutes(api *gin.RouterGroup) {
	stack := func(extra ...gin.HandlerFunc) []gin.HandlerFunc {
		out := make([]gin.HandlerFunc, 0, 3+len(extra))
		out = append(out, m.auth, m.org, m.load)
		return append(out, extra...)
	}

	// Caller's own effective access — any authenticated org member.
	me := api.Group("/me")
	me.Use(stack()...)
	me.GET("/access", m.Handler.Me)

	// Permission catalog (needed to render the matrix).
	perms := api.Group("/permissions")
	perms.Use(stack(m.guard.Require(rbac.PermRoleManage))...)
	perms.GET("", m.Handler.ListPermissions)

	roles := api.Group("/roles")
	roles.Use(stack(m.guard.Require(rbac.PermRoleManage))...)

	roles.GET("", m.Handler.List)
	roles.POST("", m.Handler.Create)
	roles.GET("/:id", m.Handler.Get)
	roles.PUT("/:id", m.Handler.Update)
	roles.DELETE("/:id", m.Handler.Delete)
	roles.PUT("/:id/permissions", m.Handler.SetPermissions)
	roles.PUT("/:id/module-access", m.Handler.SetModuleAccess)
	roles.PUT("/:id/field-access", m.Handler.SetFieldAccess)
}
