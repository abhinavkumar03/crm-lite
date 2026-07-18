package app

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"github.com/abhinavkumar03/crm-lite/backend/internal/shared/config"
	"github.com/abhinavkumar03/crm-lite/backend/internal/shared/middleware"
)

// NewRouter builds the Gin engine and the global middleware chain. Config is
// injected (rather than re-loaded here) so configuration is read exactly once
// at process start.
func NewRouter(
	logger *zap.Logger,
	cfg *config.Config,
	modules ...Module,
) *gin.Engine {

	router := gin.New()

	router.Use(
		middleware.RequestID(),
		middleware.Logger(logger),
		middleware.Recovery(logger),
		middleware.SecurityHeaders(),
	)

	// Single, credential-aware CORS policy driven by configured origins.
	router.Use(cors.New(cors.Config{
		AllowOrigins: cfg.FrontendURLs,
		AllowMethods: []string{
			"GET",
			"POST",
			"PUT",
			"PATCH",
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
