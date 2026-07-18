package importer

import (
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"

	fieldrepository "github.com/abhinavkumar03/crm-lite/backend/internal/field/repository"
	"github.com/abhinavkumar03/crm-lite/backend/internal/importer/handler"
	"github.com/abhinavkumar03/crm-lite/backend/internal/importer/repository"
	"github.com/abhinavkumar03/crm-lite/backend/internal/importer/service"
	"github.com/abhinavkumar03/crm-lite/backend/internal/jobs"
	"github.com/abhinavkumar03/crm-lite/backend/internal/rbac"
)

// Module is the import-engine composition root (API side). It exposes analyze +
// start + read endpoints; the actual row processing runs in the worker's
// Processor. It reuses the field engine's repository for module/field metadata.
type Module struct {
	Handler *handler.ImportHandler
	auth    gin.HandlerFunc
	org     gin.HandlerFunc
	load    gin.HandlerFunc
	guard   *rbac.Guard
}

func NewModule(db *pgxpool.Pool, auth, org, load gin.HandlerFunc, guard *rbac.Guard, producer *jobs.Producer) *Module {
	repo := repository.New(db)
	fieldRepo := fieldrepository.New(db)
	svc := service.New(repo, fieldRepo, producer)
	h := handler.New(svc)

	return &Module{
		Handler: h,
		auth:    auth,
		org:     org,
		load:    load,
		guard:   guard,
	}
}

// RegisterRoutes mounts the import API under a module. All endpoints require
// import.run; analyze/create also require module-level create access.
func (m *Module) RegisterRoutes(api *gin.RouterGroup) {
	imports := api.Group("/modules/:id/imports")
	imports.Use(m.auth, m.org, m.load, m.guard.Require(rbac.PermImportRun))

	imports.POST("/analyze", m.guard.RequireModule(rbac.ActionCreate), m.Handler.Analyze)
	imports.POST("", m.guard.RequireModule(rbac.ActionCreate), m.Handler.Create)
	imports.GET("", m.guard.RequireModule(rbac.ActionView), m.Handler.List)
	imports.GET("/:importId", m.guard.RequireModule(rbac.ActionView), m.Handler.Get)
}
