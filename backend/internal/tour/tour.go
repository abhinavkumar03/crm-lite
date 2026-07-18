package tour

import (
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/abhinavkumar03/crm-lite/backend/internal/tour/handler"
	"github.com/abhinavkumar03/crm-lite/backend/internal/tour/repository"
	"github.com/abhinavkumar03/crm-lite/backend/internal/tour/service"
)

// Module is the guided-tour composition root. It exposes a small, per-user
// progress API used by the client to drive interactive onboarding (read current
// progress, persist advancement, restart).
type Module struct {
	Handler *handler.TourHandler
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

// RegisterRoutes mounts the tour API. Progress is per-user but organization
// scoped like the other multi-tenant engines.
func (m *Module) RegisterRoutes(api *gin.RouterGroup) {
	tour := api.Group("/tour")
	tour.Use(m.auth, m.org)

	tour.GET("", m.Handler.Get)
	tour.PUT("", m.Handler.Update)
	tour.POST("/restart", m.Handler.Restart)
}
