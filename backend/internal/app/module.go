package app

import "github.com/gin-gonic/gin"

type Module interface {
	RegisterRoutes(router *gin.RouterGroup)
}
