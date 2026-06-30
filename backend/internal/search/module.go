package search

import (
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/abhinavkumar03/crm-lite/backend/internal/contact/repository"
	leadRepository "github.com/abhinavkumar03/crm-lite/backend/internal/lead/repository"
	"github.com/abhinavkumar03/crm-lite/backend/internal/search/handler"
	"github.com/abhinavkumar03/crm-lite/backend/internal/search/service"
	taskRepository "github.com/abhinavkumar03/crm-lite/backend/internal/task/repository"
)

type Module struct {
	Handler *handler.SearchHandler
	auth    gin.HandlerFunc
}

func NewModule(
	db *pgxpool.Pool,
	auth gin.HandlerFunc,
) *Module {

	leadRepo := leadRepository.New(db)

	contactRepo := repository.New(db)

	taskRepo := taskRepository.New(db)

	svc := service.New(
		leadRepo,
		contactRepo,
		taskRepo,
	)

	h := handler.New(svc)

	return &Module{
		Handler: h,
		auth:    auth,
	}
}

func (m *Module) RegisterRoutes(
	api *gin.RouterGroup,
) {

	search := api.Group("/search")

	search.Use(m.auth)

	search.GET(
		"",
		m.Handler.Search,
	)
}
