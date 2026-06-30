package note

import (
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/abhinavkumar03/crm-lite/backend/internal/note/handler"
	noteRepo "github.com/abhinavkumar03/crm-lite/backend/internal/note/repository"
	"github.com/abhinavkumar03/crm-lite/backend/internal/note/service"

	contactRepo "github.com/abhinavkumar03/crm-lite/backend/internal/contact/repository"
	leadRepo "github.com/abhinavkumar03/crm-lite/backend/internal/lead/repository"
	taskRepo "github.com/abhinavkumar03/crm-lite/backend/internal/task/repository"
)

type Module struct {
	handler *handler.NoteHandler
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

	noteService := service.New(
		noteRepository,
		leadRepository,
		contactRepository,
		taskRepository,
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

	api.POST("/notes/lead/:leadId", m.handler.CreateLeadNote)
	api.GET("/notes/lead/:leadId", m.handler.ListLeadNotes)

	api.POST("/notes/task/:taskId", m.handler.CreateTaskNote)
	api.GET("/notes/task/:taskId", m.handler.ListTaskNotes)

	api.POST("/notes/contact/:contactId", m.handler.CreateTaskNote)
	api.GET("/notes/contact/:contactId", m.handler.ListTaskNotes)

	api.PUT("/notes/:noteId", m.handler.UpdateNote)
	api.DELETE("/notes/:noteId", m.handler.DeleteNote)
}
