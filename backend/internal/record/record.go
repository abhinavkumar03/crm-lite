package record

import (
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"

	fieldrepository "github.com/abhinavkumar03/crm-lite/backend/internal/field/repository"
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
}

func NewModule(db *pgxpool.Pool, auth gin.HandlerFunc, org gin.HandlerFunc) *Module {
	recordRepo := repository.New(db)
	fieldRepo := fieldrepository.New(db)
	validator := vservice.New(vrepository.New(db), fieldRepo)

	svc := service.New(recordRepo, fieldRepo, validator)
	h := handler.New(svc)

	return &Module{
		Handler: h,
		auth:    auth,
		org:     org,
	}
}

// RegisterRoutes mounts the generic record CRUD + query API under a module. The
// module id param is ":id" for consistency with the other metadata engines.
func (m *Module) RegisterRoutes(api *gin.RouterGroup) {
	records := api.Group("/modules/:id/records")
	records.Use(m.auth, m.org)

	records.GET("", m.Handler.List)
	records.POST("", m.Handler.Create)
	records.GET("/:recordId", m.Handler.Get)
	records.PUT("/:recordId", m.Handler.Update)
	records.DELETE("/:recordId", m.Handler.Delete)
}
