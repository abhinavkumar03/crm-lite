package app

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"github.com/abhinavkumar03/crm-lite/backend/internal/shared/config"
)

type Application struct {
	Config *config.Config
	Logger *zap.Logger
	Router *gin.Engine
}
