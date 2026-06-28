package lead

import (
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/abhinavkumar03/crm-lite/backend/internal/lead/handler"
	"github.com/abhinavkumar03/crm-lite/backend/internal/lead/repository"
	"github.com/abhinavkumar03/crm-lite/backend/internal/lead/service"
)

type Module struct {
	Handler *handler.LeadHandler
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

	leads := api.Group("/leads")

	leads.Use(m.auth)

	leads.POST("", m.Handler.Create)
	leads.GET("", m.Handler.List)
	leads.GET("/:id", m.Handler.GetByID)
	leads.PUT("/:id", m.Handler.Update)
	leads.DELETE("/:id", m.Handler.Delete)
}
