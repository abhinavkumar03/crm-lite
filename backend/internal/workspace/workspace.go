package workspace

import (
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/abhinavkumar03/crm-lite/backend/internal/rbac"
	"github.com/abhinavkumar03/crm-lite/backend/internal/workspace/handler"
	"github.com/abhinavkumar03/crm-lite/backend/internal/workspace/repository"
	"github.com/abhinavkumar03/crm-lite/backend/internal/workspace/service"
)

type Module struct {
	Handler *handler.Handler
	Service *service.Service
	auth    gin.HandlerFunc
	org     gin.HandlerFunc
	load    gin.HandlerFunc
	guard   *rbac.Guard
}

func NewModule(db *pgxpool.Pool, auth, org, load gin.HandlerFunc, guard *rbac.Guard) *Module {
	repo := repository.New(db)
	svc := service.New(repo)
	h := handler.New(svc)
	return &Module{Handler: h, Service: svc, auth: auth, org: org, load: load, guard: guard}
}

func (m *Module) RegisterRoutes(api *gin.RouterGroup) {
	mod := api.Group("/modules/:id")
	mod.Use(m.auth, m.org, m.load)

	mod.GET("/layouts/detail", m.guard.RequireModule(rbac.ActionView), m.Handler.GetLayout)
	mod.GET("/related", m.guard.RequireModule(rbac.ActionView), m.Handler.ListRelated)

	rec := mod.Group("/records/:recordId")
	rec.GET("/notes", m.guard.RequireModule(rbac.ActionView), m.Handler.ListNotes)
	rec.POST("/notes", m.guard.RequireModule(rbac.ActionUpdate), m.Handler.CreateNote)
	rec.DELETE("/notes/:noteId", m.guard.RequireModule(rbac.ActionUpdate), m.Handler.DeleteNote)

	rec.GET("/attachments", m.guard.RequireModule(rbac.ActionView), m.Handler.ListAttachments)
	rec.POST("/attachments", m.guard.RequireModule(rbac.ActionUpdate), m.Handler.CreateAttachment)
	rec.DELETE("/attachments/:attachmentId", m.guard.RequireModule(rbac.ActionUpdate), m.Handler.DeleteAttachment)

	rec.GET("/activities", m.guard.RequireModule(rbac.ActionView), m.Handler.ListActivities)
}
