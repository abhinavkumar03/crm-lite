package record

import (
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/abhinavkumar03/crm-lite/backend/internal/access"
	fieldrepository "github.com/abhinavkumar03/crm-lite/backend/internal/field/repository"
	orgrepo "github.com/abhinavkumar03/crm-lite/backend/internal/organization/repository"
	"github.com/abhinavkumar03/crm-lite/backend/internal/rbac"
	"github.com/abhinavkumar03/crm-lite/backend/internal/record/handler"
	"github.com/abhinavkumar03/crm-lite/backend/internal/record/repository"
	"github.com/abhinavkumar03/crm-lite/backend/internal/record/service"
	"github.com/abhinavkumar03/crm-lite/backend/internal/shared/cache"
	vrepository "github.com/abhinavkumar03/crm-lite/backend/internal/validationengine/repository"
	vservice "github.com/abhinavkumar03/crm-lite/backend/internal/validationengine/service"
)

type Module struct {
	Handler *handler.RecordHandler
	Service *service.Service
	auth    gin.HandlerFunc
	org     gin.HandlerFunc
	load    gin.HandlerFunc
	guard   *rbac.Guard
}

func NewModule(db *pgxpool.Pool, auth, org, load gin.HandlerFunc, guard *rbac.Guard, appCache *cache.Cache) *Module {
	recordRepo := repository.New(db)
	fieldRepo := fieldrepository.New(db)
	validator := vservice.New(vrepository.New(db), fieldRepo)
	accessSvc := access.New(orgrepo.New(db))

	svc := service.New(recordRepo, fieldRepo, validator, appCache, accessSvc)
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

func (m *Module) RegisterRoutes(api *gin.RouterGroup) {
	records := api.Group("/modules/:id/records")
	records.Use(m.auth, m.org, m.load)

	records.GET("", m.guard.RequireModule(rbac.ActionView), m.Handler.List)
	records.GET("/:recordId", m.guard.RequireModule(rbac.ActionView), m.Handler.Get)
	records.POST("", m.guard.RequireModule(rbac.ActionCreate), m.Handler.Create)
	records.PUT("/:recordId", m.guard.RequireModule(rbac.ActionUpdate), m.Handler.Update)
	records.DELETE("/:recordId", m.guard.RequireModule(rbac.ActionDelete), m.Handler.Delete)
}
