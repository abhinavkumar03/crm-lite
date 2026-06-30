package contact

import (
	activityRepository "github.com/abhinavkumar03/crm-lite/backend/internal/activity/repository"
	activityService "github.com/abhinavkumar03/crm-lite/backend/internal/activity/service"

	"github.com/abhinavkumar03/crm-lite/backend/internal/contact/handler"
	"github.com/abhinavkumar03/crm-lite/backend/internal/contact/repository"
	"github.com/abhinavkumar03/crm-lite/backend/internal/contact/service"
	leadRepository "github.com/abhinavkumar03/crm-lite/backend/internal/lead/repository"
	taskRepository "github.com/abhinavkumar03/crm-lite/backend/internal/task/repository"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Module struct {
	Handler *handler.ContactHandler
	auth    gin.HandlerFunc
}

func NewModule(
	db *pgxpool.Pool,
	auth gin.HandlerFunc,
) *Module {

	repo := repository.New(db)

	leadRepo := leadRepository.New(db)
	taskRepo := taskRepository.New(db)
	activityRepo := activityRepository.New(db)

	activitySvc := activityService.New(
		activityRepo,
		leadRepo,
		repo,
		taskRepo,
	)

	svc := service.New(
		repo,
		activitySvc,
	)

	h := handler.New(svc)

	return &Module{
		Handler: h,
		auth:    auth,
	}
}

func (m *Module) RegisterRoutes(api *gin.RouterGroup) {

	contact := api.Group("/contacts")

	contact.Use(m.auth)

	contact.POST("", m.Handler.Create)
	contact.GET("", m.Handler.List)
	contact.GET("/:id", m.Handler.GetByID)
	contact.PUT("/:id", m.Handler.Update)
	contact.DELETE("/:id", m.Handler.Delete)
}
