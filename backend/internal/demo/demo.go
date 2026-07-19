package demo

import (
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/abhinavkumar03/crm-lite/backend/internal/demo/handler"
	"github.com/abhinavkumar03/crm-lite/backend/internal/demo/repository"
	"github.com/abhinavkumar03/crm-lite/backend/internal/demo/service"
	"github.com/abhinavkumar03/crm-lite/backend/internal/organization/bootstrap"
	"github.com/abhinavkumar03/crm-lite/backend/internal/tenant"
)

type Module struct {
	Handler *handler.Handler
	auth    gin.HandlerFunc
}

func NewModule(db *pgxpool.Pool, resolver *tenant.Resolver, auth gin.HandlerFunc) *Module {
	repo := repository.New(db)
	boot := bootstrap.New(db)
	svc := service.New(repo, boot, resolver)
	return &Module{Handler: handler.New(svc), auth: auth}
}

func (m *Module) RegisterRoutes(api *gin.RouterGroup) {
	d := api.Group("/demo")
	d.Use(m.auth)

	d.GET("/catalog", m.Handler.Catalog)
	d.GET("/workflows/:key", m.Handler.Workflow)
	d.GET("/session", m.Handler.Active)
	d.POST("/start", m.Handler.Start)
	d.POST("/restart", m.Handler.Restart)
	d.POST("/sessions/:sessionId/validate", m.Handler.Validate)
	d.POST("/sessions/:sessionId/skip", m.Handler.Skip)
	d.POST("/sessions/:sessionId/complete", m.Handler.Complete)
	d.POST("/sessions/:sessionId/cleanup", m.Handler.Cleanup)
	d.POST("/sessions/:sessionId/events", m.Handler.LogEvent)
}
