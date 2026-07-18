package settings

import (
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/abhinavkumar03/crm-lite/backend/internal/rbac"
	"github.com/abhinavkumar03/crm-lite/backend/internal/settings/handler"
	"github.com/abhinavkumar03/crm-lite/backend/internal/settings/repository"
	"github.com/abhinavkumar03/crm-lite/backend/internal/settings/service"
)

// Module is the organization-settings composition root. It backs the "General"
// and "Automation" tabs of the Settings Center; module/field/validation
// management on the other tabs reuse the existing metadata engines.
type Module struct {
	Handler *handler.SettingsHandler
	auth    gin.HandlerFunc
	org     gin.HandlerFunc
	load    gin.HandlerFunc
	guard   *rbac.Guard
}

func NewModule(db *pgxpool.Pool, auth, org, load gin.HandlerFunc, guard *rbac.Guard) *Module {
	repo := repository.New(db)
	svc := service.New(repo)
	h := handler.New(svc)

	return &Module{
		Handler: h,
		auth:    auth,
		org:     org,
		load:    load,
		guard:   guard,
	}
}

// RegisterRoutes mounts the settings API. Reads are available to any org member
// (so the UI can render preferences); writes require settings.manage.
func (m *Module) RegisterRoutes(api *gin.RouterGroup) {
	settings := api.Group("/settings")
	settings.Use(m.auth, m.org, m.load)

	settings.GET("", m.Handler.Get)
	settings.PUT("", m.guard.Require(rbac.PermSettingsManage), m.Handler.Update)
}
