package settings

import (
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"

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

// RegisterRoutes mounts the settings API. It is organization-scoped like the
// other multi-tenant engines.
func (m *Module) RegisterRoutes(api *gin.RouterGroup) {
	settings := api.Group("/settings")
	settings.Use(m.auth, m.org)

	settings.GET("", m.Handler.Get)
	settings.PUT("", m.Handler.Update)
}
