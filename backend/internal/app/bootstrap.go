package app

import "github.com/gin-gonic/gin"

func RegisterModules(
	router *gin.Engine,
	modules ...Module,
) {

	api := router.Group("/api/v1")

	for _, module := range modules {
		module.RegisterRoutes(api)
		if pub, ok := module.(PublicModule); ok {
			pub.RegisterPublicRoutes(router)
		}
	}
}
