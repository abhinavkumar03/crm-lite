// Package docs serves the OpenAPI 3 specification and an embedded Swagger UI.
//
//	GET /api/v1/openapi.yaml  — machine-readable spec (server URL is request-dynamic)
//	GET /api/v1/docs          — interactive Swagger UI
package docs

import (
	_ "embed"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

//go:embed openapi.yaml
var openAPISpec []byte

//go:embed swagger.html
var swaggerHTML []byte

// serverURLPlaceholder is substituted at request time with the caller's base API URL.
const serverURLPlaceholder = "__API_SERVER_URL__"

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

// Spec returns the OpenAPI 3 YAML with servers.url set to this request's host.
func (m *Module) Spec(c *gin.Context) {
	base := apiBaseURL(c)
	body := strings.ReplaceAll(string(openAPISpec), serverURLPlaceholder, base)

	c.Header("Cache-Control", "no-store")
	c.Data(http.StatusOK, "application/yaml; charset=utf-8", []byte(body))
}

// UI returns the Swagger UI HTML shell. The OpenAPI URL is absolute for the
// current host so the UI works behind proxies and non-localhost ports.
func (m *Module) UI(c *gin.Context) {
	specURL := apiBaseURL(c) + "/openapi.yaml"
	html := strings.ReplaceAll(string(swaggerHTML), "__OPENAPI_URL__", specURL)

	c.Header("Cache-Control", "no-store")
	c.Data(http.StatusOK, "text/html; charset=utf-8", []byte(html))
}

// apiBaseURL builds {scheme}://{host}/api/v1 from the incoming request,
// honouring common reverse-proxy headers.
func apiBaseURL(c *gin.Context) string {
	scheme := "http"
	if c.Request.TLS != nil {
		scheme = "https"
	}
	if proto := c.GetHeader("X-Forwarded-Proto"); proto != "" {
		scheme = strings.TrimSpace(strings.Split(proto, ",")[0])
	}

	host := c.Request.Host
	if fwd := c.GetHeader("X-Forwarded-Host"); fwd != "" {
		host = strings.TrimSpace(strings.Split(fwd, ",")[0])
	}

	return scheme + "://" + host + "/api/v1"
}
