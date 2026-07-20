package app

import "github.com/gin-gonic/gin"

type Module interface {
	RegisterRoutes(router *gin.RouterGroup)
}

// PublicModule optionally mounts unauthenticated routes (webhooks, tracking pixels).
type PublicModule interface {
	RegisterPublicRoutes(router *gin.Engine)
}
