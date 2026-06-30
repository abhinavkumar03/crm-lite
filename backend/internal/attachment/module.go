package attachment

import (
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/abhinavkumar03/crm-lite/backend/internal/attachment/handler"
	"github.com/abhinavkumar03/crm-lite/backend/internal/attachment/repository"
	"github.com/abhinavkumar03/crm-lite/backend/internal/attachment/service"
	contactRepository "github.com/abhinavkumar03/crm-lite/backend/internal/contact/repository"
	leadRepository "github.com/abhinavkumar03/crm-lite/backend/internal/lead/repository"
	taskRepository "github.com/abhinavkumar03/crm-lite/backend/internal/task/repository"
)

type Module struct {
	handler *handler.AttachmentHandler
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

	handler := handler.New(service)

	return &Module{
		handler: handler,
		auth:    auth,
	}
}

func (m *Module) RegisterRoutes(
	router *gin.RouterGroup,
) {

	api := router.Group("")

	api.Use(m.auth)

	api.POST("/attachments/lead/:leadId", m.handler.CreateLeadAttachment)
	api.GET("/attachments/lead/:leadId", m.handler.ListLeadAttachments)

	api.POST("/attachments/task/:taskId", m.handler.CreateTaskAttachment)
	api.GET("/attachments/task/:taskId", m.handler.ListTaskAttachments)

	api.POST("/attachments/contact/:contactId", m.handler.CreateContactAttachment)
	api.GET("/attachments/contact/:contactId", m.handler.ListContactAttachments)

	api.DELETE("/attachments/:noteId", m.handler.DeleteAttachment)
}
