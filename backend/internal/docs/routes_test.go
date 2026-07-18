package docs

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestDocsRoutes(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	api := r.Group("/api/v1")
	NewModule().RegisterRoutes(api)

	t.Run("openapi yaml uses request host", func(t *testing.T) {
		w := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/api/v1/openapi.yaml", nil)
		req.Host = "api.example.com:8443"
		req.Header.Set("X-Forwarded-Proto", "https")
		r.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Fatalf("status %d", w.Code)
		}
		body := w.Body.String()
		if !strings.Contains(body, "openapi:") {
			t.Fatal("missing openapi version")
		}
		if strings.Contains(body, serverURLPlaceholder) {
			t.Fatal("placeholder should be replaced")
		}
		if !strings.Contains(body, "https://api.example.com:8443/api/v1") {
			t.Fatalf("expected dynamic server url, got snippet missing host")
		}
		if strings.Contains(body, "http://localhost:8080/api/v1") {
			t.Fatal("must not keep static localhost server url")
		}
		if !strings.Contains(body, "/leads") || !strings.Contains(body, "ErrorResponse") {
			t.Fatal("spec missing expected content")
		}
	})

	t.Run("swagger ui points at dynamic openapi url", func(t *testing.T) {
		w := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/api/v1/docs", nil)
		req.Host = "crm.local:8080"
		r.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Fatalf("status %d", w.Code)
		}
		body := w.Body.String()
		if !strings.Contains(body, "swagger-ui") {
			t.Fatal("expected swagger ui html")
		}
		if strings.Contains(body, "__OPENAPI_URL__") {
			t.Fatal("openapi url placeholder should be replaced")
		}
		if !strings.Contains(body, "http://crm.local:8080/api/v1/openapi.yaml") {
			t.Fatal("ui must load openapi from current host")
		}
	})
}

func TestAPIBaseURL(t *testing.T) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodGet, "/api/v1/openapi.yaml", nil)
	c.Request.Host = "proxy.example.com"
	c.Request.Header.Set("X-Forwarded-Proto", "https")
	c.Request.Header.Set("X-Forwarded-Host", "public.example.com")

	got := apiBaseURL(c)
	want := "https://public.example.com/api/v1"
	if got != want {
		t.Fatalf("got %q want %q", got, want)
	}
}
