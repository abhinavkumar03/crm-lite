package search

import (
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/abhinavkumar03/crm-lite/backend/internal/search/handler"
	"github.com/abhinavkumar03/crm-lite/backend/internal/search/repository"
	"github.com/abhinavkumar03/crm-lite/backend/internal/search/service"
)

type Module struct {
	Handler *handler.SearchHandler
	auth    gin.HandlerFunc
	org     gin.HandlerFunc
}

func NewModule(
	db *pgxpool.Pool,
	auth gin.HandlerFunc,
	org gin.HandlerFunc,
) *Module {
	repo := repository.New(db)
	svc := service.New(repo)
	h := handler.New(svc)

	return &Module{
		Handler: h,
		auth:    auth,
		org:     org,
	}
}

func (m *Module) RegisterRoutes(api *gin.RouterGroup) {
	search := api.Group("/search")
	search.Use(m.auth, m.org)
	search.GET("", m.Handler.Search)
}
