package app

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"github.com/abhinavkumar03/crm-lite/backend/internal/shared/config"
	"github.com/abhinavkumar03/crm-lite/backend/internal/shared/middleware"
)

func NewRouter(
	logger *zap.Logger,
	modules ...Module,
) *gin.Engine {

	router := gin.New()

	router.Use(
		middleware.RequestID(),
		middleware.Logger(logger),
		middleware.Recovery(logger),
		middleware.SecurityHeaders(),
		middleware.CORS(),
	)
	cfg := config.Load()

	router.Use(cors.New(cors.Config{
		AllowOrigins: cfg.FrontendURLs,
		AllowMethods: []string{
			"GET",
			"POST",
			"PUT",
			"DELETE",
			"OPTIONS",
		},
		AllowHeaders: []string{
			"Authorization",
			"Content-Type",
		},
		AllowCredentials: true,
	}))

	RegisterModules(router, modules...)

	return router
}
