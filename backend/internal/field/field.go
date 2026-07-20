package field

import (
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/abhinavkumar03/crm-lite/backend/internal/field/handler"
	"github.com/abhinavkumar03/crm-lite/backend/internal/field/repository"
	"github.com/abhinavkumar03/crm-lite/backend/internal/field/service"
	"github.com/abhinavkumar03/crm-lite/backend/internal/rbac"
)

// Module is the dynamic-field engine composition root.
type Module struct {
	Handler *handler.FieldHandler
	Service *service.Service
	auth    gin.HandlerFunc
	org     gin.HandlerFunc
	load    gin.HandlerFunc
	guard   *rbac.Guard
}

func NewModule(db *pgxpool.Pool, auth, org, load gin.HandlerFunc, guard *rbac.Guard) *Module {
	repo := repository.New(db)
	svc := service.New(repo)
	h := handler.New(svc, guard)

	return &Module{
		Handler: h,
		Service: svc,
		auth:    auth,
		org:     org,
		load:    load,
		guard:   guard,
	}
}

// RegisterRoutes mounts fields as a sub-resource of a module. Reads require
// module.view (forms need the schema); writes require field.manage. Hidden
// fields are filtered from list/get responses based on role_field_access.
func (m *Module) RegisterRoutes(api *gin.RouterGroup) {
	fields := api.Group("/modules/:id/fields")
	fields.Use(m.auth, m.org, m.load)

	fields.GET("", m.guard.Require(rbac.PermModuleView), m.Handler.List)
	fields.GET("/:fieldId", m.guard.Require(rbac.PermModuleView), m.Handler.GetByID)

	fields.POST("", m.guard.Require(rbac.PermFieldManage), m.Handler.Create)
	fields.POST("/reorder", m.guard.Require(rbac.PermFieldManage), m.Handler.Reorder)
	fields.PUT("/:fieldId", m.guard.Require(rbac.PermFieldManage), m.Handler.Update)
	fields.DELETE("/:fieldId", m.guard.Require(rbac.PermFieldManage), m.Handler.Delete)
}
