package task

import (
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"

	activityRepository "github.com/abhinavkumar03/crm-lite/backend/internal/activity/repository"
	activityService "github.com/abhinavkumar03/crm-lite/backend/internal/activity/service"

	contactrepository "github.com/abhinavkumar03/crm-lite/backend/internal/contact/repository"
	leadrepository "github.com/abhinavkumar03/crm-lite/backend/internal/lead/repository"
	"github.com/abhinavkumar03/crm-lite/backend/internal/shared/cache"
	"github.com/abhinavkumar03/crm-lite/backend/internal/task/handler"
	"github.com/abhinavkumar03/crm-lite/backend/internal/task/repository"
	"github.com/abhinavkumar03/crm-lite/backend/internal/task/service"
)

type Module struct {
	Handler *handler.TaskHandler
	auth    gin.HandlerFunc
}

func NewModule(
	db *pgxpool.Pool,
	auth gin.HandlerFunc,
	c *cache.Cache,
) *Module {

	taskRepo := repository.New(db)
	leadRepo := leadrepository.New(db)
	contactRepo := contactrepository.New(db)

	activityRepo := activityRepository.New(db)

	activitySvc := activityService.New(
		activityRepo,
		leadRepo,
		contactRepo,
		taskRepo,
	)

	svc := service.New(
		taskRepo,
		leadRepo,
		contactRepo,
		activitySvc,
		c,
	)

	h := handler.New(svc)

	return &Module{
		Handler: h,
		auth:    auth,
	}
}

func (m *Module) RegisterRoutes(api *gin.RouterGroup) {

	tasks := api.Group("/tasks")
	tasks.Use(m.auth)

	tasks.POST("", m.Handler.Create)
	tasks.GET("", m.Handler.List)
	tasks.GET("/:id", m.Handler.GetByID)
	tasks.PUT("/:id", m.Handler.Update)
	tasks.DELETE("/:id", m.Handler.Delete)
}
