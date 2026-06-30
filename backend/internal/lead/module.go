package lead

import (
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"

	activityRepository "github.com/abhinavkumar03/crm-lite/backend/internal/activity/repository"
	activityService "github.com/abhinavkumar03/crm-lite/backend/internal/activity/service"

	contactRepository "github.com/abhinavkumar03/crm-lite/backend/internal/contact/repository"
	"github.com/abhinavkumar03/crm-lite/backend/internal/jobs"
	"github.com/abhinavkumar03/crm-lite/backend/internal/lead/handler"
	"github.com/abhinavkumar03/crm-lite/backend/internal/lead/repository"
	"github.com/abhinavkumar03/crm-lite/backend/internal/lead/service"
	taskRepository "github.com/abhinavkumar03/crm-lite/backend/internal/task/repository"
)

type Module struct {
	Handler *handler.LeadHandler
	auth    gin.HandlerFunc
}

func NewModule(
	db *pgxpool.Pool,
	auth gin.HandlerFunc,
	producer *jobs.Producer,
) *Module {

	repo := repository.New(db)

	activityRepo := activityRepository.New(db)

	contactRepo := contactRepository.New(db)
	taskRepo := taskRepository.New(db)

	activitySvc := activityService.New(
		activityRepo,
		repo,
		contactRepo,
		taskRepo,
	)

	svc := service.New(
		repo,
		producer,
		activitySvc,
	)

	h := handler.New(svc)

	return &Module{
		Handler: h,
		auth:    auth,
	}
}

func (m *Module) RegisterRoutes(api *gin.RouterGroup) {

	leads := api.Group("/leads")

	leads.Use(m.auth)

	leads.POST("", m.Handler.Create)
	leads.GET("", m.Handler.List)
	leads.GET("/:id", m.Handler.GetByID)
	leads.PUT("/:id", m.Handler.Update)
	leads.DELETE("/:id", m.Handler.Delete)
}
