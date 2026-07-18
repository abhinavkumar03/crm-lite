package rbac

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestHasAndRequire(t *testing.T) {
	gin.SetMode(gin.TestMode)

	g := &Guard{} // Require/Has only need context; no DB

	router := gin.New()
	router.GET("/ok", func(c *gin.Context) {
		c.Set(ContextPermissions, []string{PermModuleView, PermExportRun})
		c.Next()
	}, g.Require(PermExportRun), func(c *gin.Context) {
		c.Status(http.StatusOK)
	})
	router.GET("/deny", func(c *gin.Context) {
		c.Set(ContextPermissions, []string{PermModuleView})
		c.Next()
	}, g.Require(PermImportRun), func(c *gin.Context) {
		c.Status(http.StatusOK)
	})

	w := httptest.NewRecorder()
	router.ServeHTTP(w, httptest.NewRequest(http.MethodGet, "/ok", nil))
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}

	w = httptest.NewRecorder()
	router.ServeHTTP(w, httptest.NewRequest(http.MethodGet, "/deny", nil))
	if w.Code != http.StatusForbidden {
		t.Fatalf("expected 403, got %d", w.Code)
	}
}

func TestRequireAny(t *testing.T) {
	gin.SetMode(gin.TestMode)
	g := &Guard{}

	router := gin.New()
	router.GET("/x", func(c *gin.Context) {
		c.Set(ContextPermissions, []string{PermRecordView})
		c.Next()
	}, g.RequireAny(PermModuleManage, PermRecordView), func(c *gin.Context) {
		c.Status(http.StatusOK)
	})

	w := httptest.NewRecorder()
	router.ServeHTTP(w, httptest.NewRequest(http.MethodGet, "/x", nil))
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}
}
