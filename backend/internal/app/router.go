package app

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"github.com/abhinavkumar03/crm-lite/backend/internal/routes"
	"github.com/abhinavkumar03/crm-lite/backend/internal/shared/middleware"
)

func NewRouter(
	logger *zap.Logger,
) *gin.Engine {

	router := gin.New()

	router.Use(
		middleware.RequestID(),
		middleware.Logger(logger),
		middleware.Recovery(logger),
		middleware.SecurityHeaders(),
		middleware.CORS(),
	)

	routes.Register(router)

	return router
}
