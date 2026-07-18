package importer

import (
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"

	fieldrepository "github.com/abhinavkumar03/crm-lite/backend/internal/field/repository"
	"github.com/abhinavkumar03/crm-lite/backend/internal/importer/handler"
	"github.com/abhinavkumar03/crm-lite/backend/internal/importer/repository"
	"github.com/abhinavkumar03/crm-lite/backend/internal/importer/service"
	"github.com/abhinavkumar03/crm-lite/backend/internal/jobs"
)

// Module is the import-engine composition root (API side). It exposes analyze +
// start + read endpoints; the actual row processing runs in the worker's
// Processor. It reuses the field engine's repository for module/field metadata.
type Module struct {
	Handler *handler.ImportHandler
	auth    gin.HandlerFunc
	org     gin.HandlerFunc
}

func NewModule(db *pgxpool.Pool, auth gin.HandlerFunc, org gin.HandlerFunc, producer *jobs.Producer) *Module {
	repo := repository.New(db)
	fieldRepo := fieldrepository.New(db)
	svc := service.New(repo, fieldRepo, producer)
	h := handler.New(svc)

	return &Module{
		Handler: h,
		auth:    auth,
		org:     org,
	}
}

// RegisterRoutes mounts the import API under a module. analyze is a distinct
// static sub-path so it never collides with the :importId parameter route.
func (m *Module) RegisterRoutes(api *gin.RouterGroup) {
	imports := api.Group("/modules/:id/imports")
	imports.Use(m.auth, m.org)

	imports.POST("/analyze", m.Handler.Analyze)
	imports.POST("", m.Handler.Create)
	imports.GET("", m.Handler.List)
	imports.GET("/:importId", m.Handler.Get)
}
