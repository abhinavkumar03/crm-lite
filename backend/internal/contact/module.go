package contact

import (
	"github.com/abhinavkumar03/crm-lite/backend/internal/contact/handler"
	"github.com/abhinavkumar03/crm-lite/backend/internal/contact/repository"
	"github.com/abhinavkumar03/crm-lite/backend/internal/contact/service"
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

	svc := service.New(repo)

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

}
