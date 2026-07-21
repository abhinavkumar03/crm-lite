package workflow

import (
	"context"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"

	fieldrepository "github.com/abhinavkumar03/crm-lite/backend/internal/field/repository"
	"github.com/abhinavkumar03/crm-lite/backend/internal/jobs"
	"github.com/abhinavkumar03/crm-lite/backend/internal/rbac"
	rdto "github.com/abhinavkumar03/crm-lite/backend/internal/record/dto"
	"github.com/abhinavkumar03/crm-lite/backend/internal/workflow/engine"
	"github.com/abhinavkumar03/crm-lite/backend/internal/workflow/entity"
	"github.com/abhinavkumar03/crm-lite/backend/internal/workflow/handler"
	"github.com/abhinavkumar03/crm-lite/backend/internal/workflow/processor"
	"github.com/abhinavkumar03/crm-lite/backend/internal/workflow/repository"
	"github.com/abhinavkumar03/crm-lite/backend/internal/workflow/service"
)

// Module is the workflow automation composition root.
type Module struct {
	Handler   *handler.Handler
	Service   *service.Service
	Processor *processor.Processor
	Publisher *EventPublisher
	auth      gin.HandlerFunc
	org       gin.HandlerFunc
	load      gin.HandlerFunc
	guard     *rbac.Guard
	enabled   bool
}

type ModuleDeps struct {
	DB       *pgxpool.Pool
	Auth     gin.HandlerFunc
	Org      gin.HandlerFunc
	Load     gin.HandlerFunc
	Guard    *rbac.Guard
	Producer *jobs.Producer
	Enabled  bool
	Logger   *zap.Logger
	Records  engine.RecordMutator
	Notify   engine.Notifier
	Notes    engine.NoteWriter
	Activities engine.ActivityWriter
}

func NewModule(d ModuleDeps) *Module {
	repo := repository.New(d.DB)
	fieldRepo := fieldrepository.New(d.DB)
	pub := &service.AsynqPublisher{P: d.Producer}
	svc := service.New(repo, fieldRepo, pub, d.Enabled)
	h := handler.New(svc)

	deps := engine.ActionDeps{
		Records:    d.Records,
		Notify:     d.Notify,
		Notes:      d.Notes,
		Activities: d.Activities,
	}
	proc := processor.New(repo, deps, d.Producer, d.Logger)
	publisher := &EventPublisher{svc: svc, enabled: d.Enabled}

	return &Module{
		Handler: h, Service: svc, Processor: proc, Publisher: publisher,
		auth: d.Auth, org: d.Org, load: d.Load, guard: d.Guard, enabled: d.Enabled,
	}
}

func (m *Module) RegisterRoutes(api *gin.RouterGroup) {
	if !m.enabled {
		return
	}
	wf := api.Group("/workflows")
	wf.Use(m.auth, m.org, m.load)

	anyView := m.guard.RequireAny(rbac.PermWorkflowView, rbac.PermAutomationManage)
	anyCreate := m.guard.RequireAny(rbac.PermWorkflowCreate, rbac.PermAutomationManage)
	anyEdit := m.guard.RequireAny(rbac.PermWorkflowEdit, rbac.PermAutomationManage)
	anyDelete := m.guard.RequireAny(rbac.PermWorkflowDelete, rbac.PermAutomationManage)
	anyPublish := m.guard.RequireAny(rbac.PermWorkflowPublish, rbac.PermAutomationManage)
	anyExecute := m.guard.RequireAny(rbac.PermWorkflowExecute, rbac.PermAutomationManage)
	anyLogs := m.guard.RequireAny(rbac.PermWorkflowLogsView, rbac.PermAutomationManage)

	wf.GET("/builder-metadata", anyView, m.Handler.BuilderMetadata)
	wf.GET("/metrics", m.guard.RequireAny(rbac.PermWorkflowView, rbac.PermAutomationManage, rbac.PermAnalyticsView), m.Handler.Metrics)
	wf.GET("/templates", anyView, m.Handler.ListTemplates)
	wf.POST("/templates/:id/clone", anyCreate, m.Handler.CloneTemplate)
	wf.GET("/executions", anyLogs, m.Handler.ListExecutions)
	wf.GET("/executions/:id", anyLogs, m.Handler.GetExecution)
	wf.POST("/executions/:id/retry", anyExecute, m.Handler.RetryExecution)

	wf.GET("", anyView, m.Handler.List)
	wf.POST("", anyCreate, m.Handler.Create)
	wf.GET("/:id", anyView, m.Handler.Get)
	wf.PATCH("/:id", anyEdit, m.Handler.Update)
	wf.DELETE("/:id", anyDelete, m.Handler.Delete)
	wf.POST("/:id/publish", anyPublish, m.Handler.Publish)
	wf.POST("/:id/disable", anyPublish, m.Handler.Disable)
	wf.GET("/:id/versions", anyView, m.Handler.Versions)
	wf.POST("/:id/versions/:versionId/rollback", anyPublish, m.Handler.Rollback)
	wf.POST("/:id/run", anyExecute, m.Handler.ManualRun)
}

// EventPublisher implements record MutationHook.
type EventPublisher struct {
	svc     *service.Service
	enabled bool
}

func (p *EventPublisher) AfterCreate(ctx context.Context, orgID, moduleID, userID string, rec *rdto.RecordResponse) error {
	if p == nil || !p.enabled || p.svc == nil || rec == nil {
		return nil
	}
	meta := mutationMeta(ctx)
	after := recordSnapshot(rec)
	p.svc.EnqueueRecordEvent(ctx, orgID, moduleID, rec.ID, userID, entity.TriggerRecordCreated,
		nil, after, nil, meta.Source, meta.Depth, meta.ExcludeWorkflowID)
	return nil
}

func (p *EventPublisher) AfterUpdate(ctx context.Context, orgID, moduleID, userID string, before, after *rdto.RecordResponse) error {
	if p == nil || !p.enabled || p.svc == nil || after == nil {
		return nil
	}
	meta := mutationMeta(ctx)
	b := map[string]any{}
	if before != nil {
		b = recordSnapshot(before)
	}
	a := recordSnapshot(after)
	changed := engine.ChangedFields(b, a)
	p.svc.EnqueueRecordEvent(ctx, orgID, moduleID, after.ID, userID, entity.TriggerRecordUpdated,
		b, a, changed, meta.Source, meta.Depth, meta.ExcludeWorkflowID)
	if len(changed) > 0 {
		p.svc.EnqueueRecordEvent(ctx, orgID, moduleID, after.ID, userID, entity.TriggerFieldUpdated,
			b, a, changed, meta.Source, meta.Depth, meta.ExcludeWorkflowID)
	}
	return nil
}

func (p *EventPublisher) AfterDelete(ctx context.Context, orgID, moduleID, userID, recordID string, before *rdto.RecordResponse) error {
	if p == nil || !p.enabled || p.svc == nil {
		return nil
	}
	meta := mutationMeta(ctx)
	b := map[string]any{}
	if before != nil {
		b = recordSnapshot(before)
	}
	p.svc.EnqueueRecordEvent(ctx, orgID, moduleID, recordID, userID, entity.TriggerRecordDeleted,
		b, nil, nil, meta.Source, meta.Depth, meta.ExcludeWorkflowID)
	return nil
}

func mutationMeta(ctx context.Context) engine.MutationMeta {
	if m, ok := engine.MutationMetaFrom(ctx); ok {
		return m
	}
	return engine.MutationMeta{Source: entity.SourceUser, Depth: 0}
}

func recordSnapshot(rec *rdto.RecordResponse) map[string]any {
	out := map[string]any{}
	for k, v := range rec.Data {
		out[k] = v
	}
	sys := map[string]any{"visibility": rec.Visibility}
	if rec.OwnerID != nil {
		sys["owner_id"] = *rec.OwnerID
	}
	if rec.AssignedTo != nil {
		sys["assigned_to"] = *rec.AssignedTo
	}
	if rec.TeamID != nil {
		sys["team_id"] = *rec.TeamID
	}
	if rec.DepartmentID != nil {
		sys["department_id"] = *rec.DepartmentID
	}
	out["_system"] = sys
	return out
}
