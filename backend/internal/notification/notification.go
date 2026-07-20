package notification

import (
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"

	activityrepo "github.com/abhinavkumar03/crm-lite/backend/internal/activity/repository"
	"github.com/abhinavkumar03/crm-lite/backend/internal/jobs"
	"github.com/abhinavkumar03/crm-lite/backend/internal/notification/handler"
	"github.com/abhinavkumar03/crm-lite/backend/internal/notification/repository"
	"github.com/abhinavkumar03/crm-lite/backend/internal/notification/service"
	"github.com/abhinavkumar03/crm-lite/backend/internal/rbac"
	"github.com/abhinavkumar03/crm-lite/backend/internal/shared/secrets"
)

// Module is the notification-engine composition root (API side). It exposes
// compose/list/detail/retry/cancel/metrics, template CRUD, provider management,
// and public webhook endpoints.
type Module struct {
	Handler        *handler.NotificationHandler
	ProviderHandler *handler.ProviderHandler
	WebhookHandler *handler.WebhookHandler
	ProviderService *service.ProviderService
	auth           gin.HandlerFunc
	org            gin.HandlerFunc
	load           gin.HandlerFunc
	guard          *rbac.Guard
}

type ModuleDeps struct {
	DB            *pgxpool.Pool
	Auth          gin.HandlerFunc
	Org           gin.HandlerFunc
	Load          gin.HandlerFunc
	Guard         *rbac.Guard
	Producer      *jobs.Producer
	SecretsBox    *secrets.Box
	Logger        *zap.Logger
	MetaSecret    string
	MetaVerify    string
	ResendSecret  string
}

func NewModule(db *pgxpool.Pool, auth, org, load gin.HandlerFunc, guard *rbac.Guard, producer *jobs.Producer) *Module {
	return NewModuleWithDeps(ModuleDeps{
		DB: db, Auth: auth, Org: org, Load: load, Guard: guard, Producer: producer,
	})
}

func NewModuleWithDeps(d ModuleDeps) *Module {
	repo := repository.New(d.DB)
	svc := service.New(repo, d.Producer)
	h := handler.New(svc)

	var providerSvc *service.ProviderService
	var providerH *handler.ProviderHandler
	if d.SecretsBox != nil {
		providerSvc = service.NewProviderService(repository.NewProviderRepository(d.DB), d.SecretsBox, d.Logger)
		providerH = handler.NewProviderHandler(providerSvc)
	}

	webhookH := handler.NewWebhookHandler(
		repo,
		activityrepo.New(d.DB),
		d.Logger,
		d.MetaSecret,
		d.MetaVerify,
		d.ResendSecret,
	)

	return &Module{
		Handler:         h,
		ProviderHandler: providerH,
		WebhookHandler:  webhookH,
		ProviderService: providerSvc,
		auth:            d.Auth,
		org:             d.Org,
		load:            d.Load,
		guard:           d.Guard,
	}
}

// RegisterRoutes mounts the notification center API.
func (m *Module) RegisterRoutes(api *gin.RouterGroup) {
	notifications := api.Group("/notifications")
	notifications.Use(m.auth, m.org, m.load)

	notifications.GET("/metrics", m.guard.RequireAny(rbac.PermNotificationView, rbac.PermAutomationManage, rbac.PermAnalyticsView), m.Handler.Metrics)
	notifications.GET("", m.guard.RequireAny(rbac.PermNotificationView, rbac.PermAutomationManage), m.Handler.List)
	notifications.POST("", m.guard.RequireAny(rbac.PermNotificationSend, rbac.PermAutomationManage), m.Handler.Compose)
	notifications.GET("/:id", m.guard.RequireAny(rbac.PermNotificationView, rbac.PermAutomationManage), m.Handler.Get)
	notifications.PATCH("/:id", m.guard.RequireAny(rbac.PermNotificationSend, rbac.PermAutomationManage), m.Handler.UpdateDraft)
	notifications.POST("/:id/send", m.guard.RequireAny(rbac.PermNotificationSend, rbac.PermAutomationManage), m.Handler.SendDraft)
	notifications.POST("/:id/retry", m.guard.RequireAny(rbac.PermNotificationSend, rbac.PermAutomationManage), m.Handler.Retry)
	notifications.POST("/:id/cancel", m.guard.RequireAny(rbac.PermNotificationSend, rbac.PermAutomationManage), m.Handler.Cancel)

	templates := api.Group("/notification-templates")
	templates.Use(m.auth, m.org, m.load)
	templates.GET("", m.guard.RequireAny(rbac.PermNotificationView, rbac.PermNotificationTemplatesManage, rbac.PermAutomationManage), m.Handler.ListTemplates)
	templates.POST("", m.guard.RequireAny(rbac.PermNotificationTemplatesManage, rbac.PermAutomationManage), m.Handler.CreateTemplate)
	templates.GET("/:id", m.guard.RequireAny(rbac.PermNotificationView, rbac.PermNotificationTemplatesManage, rbac.PermAutomationManage), m.Handler.GetTemplate)
	templates.PUT("/:id", m.guard.RequireAny(rbac.PermNotificationTemplatesManage, rbac.PermAutomationManage), m.Handler.UpdateTemplate)
	templates.POST("/:id/publish", m.guard.RequireAny(rbac.PermNotificationTemplatesManage, rbac.PermAutomationManage), m.Handler.PublishTemplate)
	templates.DELETE("/:id", m.guard.RequireAny(rbac.PermNotificationTemplatesManage, rbac.PermAutomationManage), m.Handler.DeleteTemplate)
	templates.POST("/:id/preview", m.guard.RequireAny(rbac.PermNotificationView, rbac.PermNotificationTemplatesManage, rbac.PermAutomationManage), m.Handler.PreviewTemplate)

	if m.ProviderHandler != nil {
		providers := api.Group("/communication-providers")
		providers.Use(m.auth, m.org, m.load)
		providers.GET("", m.guard.RequireAny(rbac.PermCommunicationProvidersManage, rbac.PermSettingsManage), m.ProviderHandler.List)
		providers.POST("", m.guard.RequireAny(rbac.PermCommunicationProvidersManage, rbac.PermSettingsManage), m.ProviderHandler.Create)
		providers.GET("/:id", m.guard.RequireAny(rbac.PermCommunicationProvidersManage, rbac.PermSettingsManage), m.ProviderHandler.Get)
		providers.PUT("/:id", m.guard.RequireAny(rbac.PermCommunicationProvidersManage, rbac.PermSettingsManage), m.ProviderHandler.Update)
		providers.DELETE("/:id", m.guard.RequireAny(rbac.PermCommunicationProvidersManage, rbac.PermSettingsManage), m.ProviderHandler.Delete)
		providers.POST("/:id/test", m.guard.RequireAny(rbac.PermCommunicationProvidersManage, rbac.PermSettingsManage), m.ProviderHandler.Test)

		senders := api.Group("/communication-senders")
		senders.Use(m.auth, m.org, m.load)
		senders.GET("", m.guard.RequireAny(rbac.PermCommunicationProvidersManage, rbac.PermSettingsManage), m.ProviderHandler.ListSenders)
		senders.POST("", m.guard.RequireAny(rbac.PermCommunicationProvidersManage, rbac.PermSettingsManage), m.ProviderHandler.CreateSender)
		senders.DELETE("/:id", m.guard.RequireAny(rbac.PermCommunicationProvidersManage, rbac.PermSettingsManage), m.ProviderHandler.DeleteSender)
	}
}

// RegisterPublicRoutes mounts unauthenticated webhook + tracking endpoints.
func (m *Module) RegisterPublicRoutes(r *gin.Engine) {
	if m.WebhookHandler == nil {
		return
	}
	r.GET("/webhooks/whatsapp/meta", m.WebhookHandler.MetaVerify)
	r.POST("/webhooks/whatsapp/meta", m.WebhookHandler.MetaStatus)
	r.POST("/webhooks/whatsapp/twilio", m.WebhookHandler.TwilioStatus)
	r.POST("/webhooks/email/resend", m.WebhookHandler.ResendEvents)
	r.GET("/t/o/:token", m.WebhookHandler.OpenPixel)
}
