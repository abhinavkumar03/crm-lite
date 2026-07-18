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

	t.Run("openapi yaml", func(t *testing.T) {
		w := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/api/v1/openapi.yaml", nil)
		r.ServeHTTP(w, req)
		if w.Code != http.StatusOK {
			t.Fatalf("status %d", w.Code)
		}
		ct := w.Header().Get("Content-Type")
		if !strings.Contains(ct, "yaml") {
			t.Fatalf("unexpected content-type %q", ct)
		}
		body := w.Body.String()
		if !strings.Contains(body, "openapi:") {
			t.Fatal("missing openapi version")
		}
		if !strings.Contains(body, "/leads") || !strings.Contains(body, "/roles") {
			t.Fatal("spec missing expected paths")
		}
		if !strings.Contains(body, "ErrorResponse") {
			t.Fatal("spec missing ErrorResponse schema")
		}
	})

	t.Run("swagger ui", func(t *testing.T) {
		w := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/api/v1/docs", nil)
		r.ServeHTTP(w, req)
		if w.Code != http.StatusOK {
			t.Fatalf("status %d", w.Code)
		}
		if !strings.Contains(w.Body.String(), "swagger-ui") {
			t.Fatal("expected swagger ui html")
		}
		if !strings.Contains(w.Body.String(), "/api/v1/openapi.yaml") {
			t.Fatal("ui must point at openapi.yaml")
		}
	})
}
