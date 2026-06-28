package dashboard

import (
	"github.com/abhinavkumar03/crm-lite/backend/internal/dashboard/handler"
	"github.com/abhinavkumar03/crm-lite/backend/internal/dashboard/repository"
	"github.com/abhinavkumar03/crm-lite/backend/internal/dashboard/service"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Module struct {
	Handler *handler.DashboardHandler
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

	dashboard := api.Group("/dashboard")

	dashboard.Use(m.auth)

	dashboard.GET("", m.Handler.Dashboard)
}
