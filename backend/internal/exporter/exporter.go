package exporter

import (
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/abhinavkumar03/crm-lite/backend/internal/exporter/handler"
	"github.com/abhinavkumar03/crm-lite/backend/internal/exporter/repository"
	"github.com/abhinavkumar03/crm-lite/backend/internal/exporter/service"
	fieldrepository "github.com/abhinavkumar03/crm-lite/backend/internal/field/repository"
	"github.com/abhinavkumar03/crm-lite/backend/internal/jobs"
	recordrepository "github.com/abhinavkumar03/crm-lite/backend/internal/record/repository"
	recordservice "github.com/abhinavkumar03/crm-lite/backend/internal/record/service"
	vrepository "github.com/abhinavkumar03/crm-lite/backend/internal/validationengine/repository"
	vservice "github.com/abhinavkumar03/crm-lite/backend/internal/validationengine/service"
)

// Module is the export-engine composition root (API side). It exposes sync
// downloads, async export jobs (processed by the worker) and reusable templates.
// It reuses the Phase 10 record runtime for fetching rows and the field engine's
// repository for metadata.
type Module struct {
	Handler *handler.ExportHandler
	auth    gin.HandlerFunc
	org     gin.HandlerFunc
}

// NewService builds the export service, wiring the record runtime as the row
// source. It is exported so the worker can construct the same service for the
// asynchronous processor without duplicating the dependency graph.
func NewService(db *pgxpool.Pool, producer *jobs.Producer) *service.Service {
	fieldRepo := fieldrepository.New(db)
	validator := vservice.New(vrepository.New(db), fieldRepo)
	recordSvc := recordservice.New(recordrepository.New(db), fieldRepo, validator)

	return service.New(
		repository.NewExportRepository(db),
		repository.NewTemplateRepository(db),
		recordSvc,
		fieldRepo,
		producer,
	)
}

func NewModule(db *pgxpool.Pool, auth gin.HandlerFunc, org gin.HandlerFunc, producer *jobs.Producer) *Module {
	return &Module{
		Handler: handler.New(NewService(db, producer)),
		auth:    auth,
		org:     org,
	}
}

// RegisterRoutes mounts the export API under a module. "export" (sync),
// "exports" (async jobs) and "export-templates" are distinct static sub-paths so
// they never collide with each other or with the sibling engines' param routes.
func (m *Module) RegisterRoutes(api *gin.RouterGroup) {
	g := api.Group("/modules/:id")
	g.Use(m.auth, m.org)

	g.GET("/export", m.Handler.ExportNow)

	exports := g.Group("/exports")
	exports.POST("", m.Handler.Create)
	exports.GET("", m.Handler.List)
	exports.GET("/:exportId", m.Handler.Get)
	exports.GET("/:exportId/download", m.Handler.Download)

	templates := g.Group("/export-templates")
	templates.GET("", m.Handler.ListTemplates)
	templates.POST("", m.Handler.CreateTemplate)
	templates.PUT("/:templateId", m.Handler.UpdateTemplate)
	templates.DELETE("/:templateId", m.Handler.DeleteTemplate)
}
