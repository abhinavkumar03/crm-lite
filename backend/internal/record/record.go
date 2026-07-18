package record

import (
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"

	fieldrepository "github.com/abhinavkumar03/crm-lite/backend/internal/field/repository"
	"github.com/abhinavkumar03/crm-lite/backend/internal/rbac"
	"github.com/abhinavkumar03/crm-lite/backend/internal/record/handler"
	"github.com/abhinavkumar03/crm-lite/backend/internal/record/repository"
	"github.com/abhinavkumar03/crm-lite/backend/internal/record/service"
	vrepository "github.com/abhinavkumar03/crm-lite/backend/internal/validationengine/repository"
	vservice "github.com/abhinavkumar03/crm-lite/backend/internal/validationengine/service"
)

// Module is the record-runtime composition root. It reuses the field engine's
// repository (metadata) and the validation engine's service (payload
// validation), keeping this engine free of duplicated logic.
type Module struct {
	Handler *handler.RecordHandler
	auth    gin.HandlerFunc
	org     gin.HandlerFunc
	load    gin.HandlerFunc
	guard   *rbac.Guard
}

func NewModule(db *pgxpool.Pool, auth, org, load gin.HandlerFunc, guard *rbac.Guard) *Module {
	recordRepo := repository.New(db)
	fieldRepo := fieldrepository.New(db)
	validator := vservice.New(vrepository.New(db), fieldRepo)

	svc := service.New(recordRepo, fieldRepo, validator)
	h := handler.New(svc, guard)

	return &Module{
		Handler: h,
		auth:    auth,
		org:     org,
		load:    load,
		guard:   guard,
	}
}

// RegisterRoutes mounts the generic record CRUD + query API under a module.
// Each verb checks the matching record.* permission and role_module_access.
func (m *Module) RegisterRoutes(api *gin.RouterGroup) {
	records := api.Group("/modules/:id/records")
	records.Use(m.auth, m.org, m.load)

	records.GET("", m.guard.RequireModule(rbac.ActionView), m.Handler.List)
	records.GET("/:recordId", m.guard.RequireModule(rbac.ActionView), m.Handler.Get)
	records.POST("", m.guard.RequireModule(rbac.ActionCreate), m.Handler.Create)
	records.PUT("/:recordId", m.guard.RequireModule(rbac.ActionUpdate), m.Handler.Update)
	records.DELETE("/:recordId", m.guard.RequireModule(rbac.ActionDelete), m.Handler.Delete)
}
