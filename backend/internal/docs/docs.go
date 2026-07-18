// Package docs serves the OpenAPI 3 specification and an embedded Swagger UI.
//
//	GET /api/v1/openapi.yaml  — machine-readable spec
//	GET /api/v1/docs          — interactive Swagger UI
package docs

import (
	_ "embed"
	"net/http"

	"github.com/gin-gonic/gin"
)

//go:embed openapi.yaml
var openAPISpec []byte

//go:embed swagger.html
var swaggerHTML []byte

// Module mounts documentation routes (no auth).
type Module struct{}

func NewModule() *Module {
	return &Module{}
}

func (m *Module) RegisterRoutes(api *gin.RouterGroup) {
	api.GET("/openapi.yaml", m.Spec)
	api.GET("/docs", m.UI)
	api.GET("/docs/", m.UI)
}

// Spec returns the OpenAPI 3 YAML document.
func (m *Module) Spec(c *gin.Context) {
	c.Header("Cache-Control", "public, max-age=60")
	c.Data(http.StatusOK, "application/yaml; charset=utf-8", openAPISpec)
}

// UI returns the Swagger UI HTML shell (loads the YAML from /api/v1/openapi.yaml).
func (m *Module) UI(c *gin.Context) {
	c.Header("Cache-Control", "public, max-age=60")
	c.Data(http.StatusOK, "text/html; charset=utf-8", swaggerHTML)
}
