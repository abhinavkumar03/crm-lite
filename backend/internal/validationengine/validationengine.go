package validationengine

import (
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"

	fieldrepository "github.com/abhinavkumar03/crm-lite/backend/internal/field/repository"
	"github.com/abhinavkumar03/crm-lite/backend/internal/validationengine/handler"
	"github.com/abhinavkumar03/crm-lite/backend/internal/validationengine/repository"
	"github.com/abhinavkumar03/crm-lite/backend/internal/validationengine/service"
)

// Module is the validation-engine composition root. It reuses the field engine's
// repository as its field-metadata source (dependency inversion).
type Module struct {
	Handler *handler.Handler
	auth    gin.HandlerFunc
	org     gin.HandlerFunc
}

func NewModule(db *pgxpool.Pool, auth gin.HandlerFunc, org gin.HandlerFunc) *Module {
	ruleRepo := repository.New(db)
	fieldRepo := fieldrepository.New(db)
	svc := service.New(ruleRepo, fieldRepo)
	h := handler.New(svc)

	return &Module{
		Handler: h,
		auth:    auth,
		org:     org,
	}
}

// RegisterRoutes mounts validation rules, the compiled schema, and a dry-run
// validate endpoint under a module. The module id param is ":id" for consistency
// with the module/field engines.
func (m *Module) RegisterRoutes(api *gin.RouterGroup) {
	rules := api.Group("/modules/:id/validation-rules")
	rules.Use(m.auth, m.org)

	rules.GET("", m.Handler.ListRules)
	rules.POST("", m.Handler.CreateRule)
	rules.GET("/:ruleId", m.Handler.GetRule)
	rules.PUT("/:ruleId", m.Handler.UpdateRule)
	rules.DELETE("/:ruleId", m.Handler.DeleteRule)

	module := api.Group("/modules/:id")
	module.Use(m.auth, m.org)

	module.GET("/validation-schema", m.Handler.Schema)
	module.POST("/validate", m.Handler.Validate)
}
