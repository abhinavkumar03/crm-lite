package activity

import (
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/abhinavkumar03/crm-lite/backend/internal/activity/handler"
	"github.com/abhinavkumar03/crm-lite/backend/internal/activity/repository"
	"github.com/abhinavkumar03/crm-lite/backend/internal/activity/service"
	contactRepository "github.com/abhinavkumar03/crm-lite/backend/internal/contact/repository"
	leadRepository "github.com/abhinavkumar03/crm-lite/backend/internal/lead/repository"
	taskRepository "github.com/abhinavkumar03/crm-lite/backend/internal/task/repository"
)

type Module struct {
	handler *handler.ActivityHandler
	auth    gin.HandlerFunc
}

func NewModule(
	db *pgxpool.Pool,
	auth gin.HandlerFunc,
) *Module {

	repository := repository.New(db)

	leadRepo := leadRepository.New(db)

	contactRepo := contactRepository.New(db)

	taskRepo := taskRepository.New(db)

	service := service.New(

		repository,

		leadRepo,

		contactRepo,

		taskRepo,
	)

	handler := handler.New(
		service,
	)

	return &Module{

		handler: handler,

		auth: auth,
	}
}

func (m *Module) RegisterRoutes(
	router *gin.RouterGroup,
) {

	api := router.Group("")

	api.Use(m.auth)

	api.GET("/activities/lead/:leadId", m.handler.ListLeadActivities)
	api.GET("/activities/contact/:contactId", m.handler.ListContactActivities)
	api.GET("/activities/task/:taskId", m.handler.ListTaskActivities)
}
