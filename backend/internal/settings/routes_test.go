package settings

import (
	"testing"

	"github.com/gin-gonic/gin"
)

// TestRegisterRoutesNoConflict verifies the settings routes register without a
// Gin panic.
func TestRegisterRoutesNoConflict(t *testing.T) {
	gin.SetMode(gin.TestMode)

	noop := func(c *gin.Context) { c.Next() }
	router := gin.New()
	api := router.Group("/api/v1")

	defer func() {
		if r := recover(); r != nil {
			t.Fatalf("route registration panicked: %v", r)
		}
	}()

	NewModule(nil, noop, noop).RegisterRoutes(api)
}
