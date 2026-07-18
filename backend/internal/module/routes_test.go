package module

import (
	"testing"

	"github.com/gin-gonic/gin"
)

// TestRegisterRoutesNoConflict ensures the module engine's route tree registers
// without a Gin panic (static vs param segment conflicts surface here).
func TestRegisterRoutesNoConflict(t *testing.T) {
	gin.SetMode(gin.TestMode)

	noop := func(c *gin.Context) { c.Next() }
	m := NewModule(nil, noop, noop)

	router := gin.New()
	api := router.Group("/api/v1")

	defer func() {
		if r := recover(); r != nil {
			t.Fatalf("route registration panicked: %v", r)
		}
	}()

	m.RegisterRoutes(api)
}
