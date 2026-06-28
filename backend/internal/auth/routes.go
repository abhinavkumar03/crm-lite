package auth

import (
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/abhinavkumar03/crm-lite/backend/internal/auth/handler"
	"github.com/abhinavkumar03/crm-lite/backend/internal/auth/repository"
	"github.com/abhinavkumar03/crm-lite/backend/internal/auth/service"
)

func RegisterRoutes(
	router *gin.RouterGroup,
	db *pgxpool.Pool,
) {

	repo := repository.New(db)

	svc := service.New(repo)

	h := handler.New(svc)

	authGroup := router.Group("/auth")

	authGroup.POST("/register", h.Register)
}
