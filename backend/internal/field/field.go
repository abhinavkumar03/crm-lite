package field

import (
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/abhinavkumar03/crm-lite/backend/internal/field/handler"
	"github.com/abhinavkumar03/crm-lite/backend/internal/field/repository"
	"github.com/abhinavkumar03/crm-lite/backend/internal/field/service"
)

// Module is the dynamic-field engine composition root.
type Module struct {
	Handler *handler.FieldHandler
	auth    gin.HandlerFunc
	org     gin.HandlerFunc
}

func NewModule(db *pgxpool.Pool, auth gin.HandlerFunc, org gin.HandlerFunc) *Module {
	repo := repository.New(db)
	svc := service.New(repo)
	h := handler.New(svc)

	return &Module{
		Handler: h,
		auth:    auth,
		org:     org,
	}
}

// RegisterRoutes mounts fields as a sub-resource of a module. The module id
// param is ":id" to stay consistent with the module engine's /modules/:id tree
// (Gin requires a single param name per path position).
func (m *Module) RegisterRoutes(api *gin.RouterGroup) {
	fields := api.Group("/modules/:id/fields")
	fields.Use(m.auth, m.org)

	fields.GET("", m.Handler.List)
	fields.POST("", m.Handler.Create)
	fields.POST("/reorder", m.Handler.Reorder)
	fields.GET("/:fieldId", m.Handler.GetByID)
	fields.PUT("/:fieldId", m.Handler.Update)
	fields.DELETE("/:fieldId", m.Handler.Delete)
}
