package validationengine

import (
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"

	fieldrepository "github.com/abhinavkumar03/crm-lite/backend/internal/field/repository"
	"github.com/abhinavkumar03/crm-lite/backend/internal/rbac"
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
	load    gin.HandlerFunc
	guard   *rbac.Guard
}

func NewModule(db *pgxpool.Pool, auth, org, load gin.HandlerFunc, guard *rbac.Guard) *Module {
	ruleRepo := repository.New(db)
	fieldRepo := fieldrepository.New(db)
	svc := service.New(ruleRepo, fieldRepo)
	h := handler.New(svc)

	return &Module{
		Handler: h,
		auth:    auth,
		org:     org,
		load:    load,
		guard:   guard,
	}
}

// RegisterRoutes mounts validation rules, the compiled schema, and a dry-run
// validate endpoint under a module.
func (m *Module) RegisterRoutes(api *gin.RouterGroup) {
	rules := api.Group("/modules/:id/validation-rules")
	rules.Use(m.auth, m.org, m.load)

	rules.GET("", m.guard.Require(rbac.PermValidationManage), m.Handler.ListRules)
	rules.POST("", m.guard.Require(rbac.PermValidationManage), m.Handler.CreateRule)
	rules.GET("/:ruleId", m.guard.Require(rbac.PermValidationManage), m.Handler.GetRule)
	rules.PUT("/:ruleId", m.guard.Require(rbac.PermValidationManage), m.Handler.UpdateRule)
	rules.DELETE("/:ruleId", m.guard.Require(rbac.PermValidationManage), m.Handler.DeleteRule)

	module := api.Group("/modules/:id")
	module.Use(m.auth, m.org, m.load)

	// Schema + dry-run are needed by forms; any caller who can view records may use them.
	module.GET("/validation-schema", m.guard.Require(rbac.PermRecordView), m.Handler.Schema)
	module.POST("/validate", m.guard.Require(rbac.PermRecordView), m.Handler.Validate)
}
