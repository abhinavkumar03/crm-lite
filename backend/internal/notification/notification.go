package notification

import (
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/abhinavkumar03/crm-lite/backend/internal/jobs"
	"github.com/abhinavkumar03/crm-lite/backend/internal/notification/handler"
	"github.com/abhinavkumar03/crm-lite/backend/internal/notification/repository"
	"github.com/abhinavkumar03/crm-lite/backend/internal/notification/service"
)

// Module is the notification-engine composition root (API side). It exposes the
// send + read API; asynchronous delivery is handled by the worker's Processor.
type Module struct {
	Handler *handler.NotificationHandler
	auth    gin.HandlerFunc
	org     gin.HandlerFunc
}

func NewModule(db *pgxpool.Pool, auth gin.HandlerFunc, org gin.HandlerFunc, producer *jobs.Producer) *Module {
	repo := repository.New(db)
	svc := service.New(repo, producer)
	h := handler.New(svc)

	return &Module{
		Handler: h,
		auth:    auth,
		org:     org,
	}
}

// RegisterRoutes mounts the notification API. It is organization-scoped like the
// other multi-tenant engines.
func (m *Module) RegisterRoutes(api *gin.RouterGroup) {
	notifications := api.Group("/notifications")
	notifications.Use(m.auth, m.org)

	notifications.POST("", m.Handler.Send)
	notifications.GET("", m.Handler.List)
	notifications.GET("/:id", m.Handler.Get)
}
