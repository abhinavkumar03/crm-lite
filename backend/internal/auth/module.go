package auth

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/abhinavkumar03/crm-lite/backend/internal/auth/handler"
	jwt "github.com/abhinavkumar03/crm-lite/backend/internal/auth/jwt"
	"github.com/abhinavkumar03/crm-lite/backend/internal/auth/middleware"
	"github.com/abhinavkumar03/crm-lite/backend/internal/auth/repository"
	"github.com/abhinavkumar03/crm-lite/backend/internal/auth/service"
)

type Module struct {
	Handler *handler.AuthHandler
	AuthMW  *middleware.AuthMiddleware
}

func NewModule(db *pgxpool.Pool, jwtSecret string, jwtExpiration time.Duration) *Module {
	repo := repository.New(db)

	jwtSvc := jwt.NewService(jwtSecret, jwtExpiration)

	svc := service.New(repo, jwtSvc)

	h := handler.New(svc)

	mw := middleware.New(jwtSvc)

	return &Module{
		Handler: h,
		AuthMW:  mw,
	}
}

func (m *Module) RegisterRoutes(api *gin.RouterGroup) {
	auth := api.Group("/auth")

	auth.POST("/register", m.Handler.Register)
	auth.POST("/login", m.Handler.Login)

	protected := auth.Group("/")
	protected.Use(m.AuthMW.Handle())

	protected.GET("/profile", m.Handler.Profile)
}

func (m *Module) Middleware() gin.HandlerFunc {
	return m.AuthMW.Handle()
}
