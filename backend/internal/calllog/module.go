package calllog

import (
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"

	activityRepository "github.com/abhinavkumar03/crm-lite/backend/internal/activity/repository"
	activityService "github.com/abhinavkumar03/crm-lite/backend/internal/activity/service"

	"github.com/abhinavkumar03/crm-lite/backend/internal/calllog/handler"
	noteRepo "github.com/abhinavkumar03/crm-lite/backend/internal/calllog/repository"
	"github.com/abhinavkumar03/crm-lite/backend/internal/calllog/service"

	contactRepo "github.com/abhinavkumar03/crm-lite/backend/internal/contact/repository"
	leadRepo "github.com/abhinavkumar03/crm-lite/backend/internal/lead/repository"
	taskRepo "github.com/abhinavkumar03/crm-lite/backend/internal/task/repository"
)

type Module struct {
	handler *handler.CallLogHandler
	auth    gin.HandlerFunc
}

func NewModule(
	db *pgxpool.Pool,
	auth gin.HandlerFunc,
) *Module {

	noteRepository := noteRepo.New(db)

	leadRepository := leadRepo.New(db)

	contactRepository := contactRepo.New(db)

	taskRepository := taskRepo.New(db)

	activityRepo := activityRepository.New(db)

	activitySvc := activityService.New(
		activityRepo,
		leadRepository,
		contactRepository,
		taskRepository,
	)

	noteService := service.New(
		noteRepository,
		leadRepository,
		contactRepository,
		taskRepository,
		activitySvc,
	)

	noteHandler := handler.New(noteService)

	return &Module{
		handler: noteHandler,
		auth:    auth,
	}
}

func (m *Module) RegisterRoutes(
	router *gin.RouterGroup,
) {

	api := router.Group("")

	api.Use(m.auth)

	api.POST("/calllogs/lead/:leadId", m.handler.CreateLeadCallLog)
	api.GET("/calllogs/lead/:leadId", m.handler.ListLeadCallLogs)

	api.POST("/calllogs/task/:taskId", m.handler.CreateTaskCallLog)
	api.GET("/calllogs/task/:taskId", m.handler.ListTaskCallLogs)

	api.POST("/calllogs/contact/:contactId", m.handler.CreateContactCallLog)
	api.GET("/calllogs/contact/:contactId", m.handler.ListContactCallLogs)

	api.PUT("/calllogs/:noteId", m.handler.UpdateCallLog)
	api.DELETE("/calllogs/:noteId", m.handler.DeleteCallLog)
}
