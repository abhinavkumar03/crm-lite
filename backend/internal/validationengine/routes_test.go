package validationengine

import (
	"testing"

	"github.com/gin-gonic/gin"

	fieldengine "github.com/abhinavkumar03/crm-lite/backend/internal/field"
	moduleengine "github.com/abhinavkumar03/crm-lite/backend/internal/module"
)

// TestRegisterRoutesNoConflict mounts the module, field, and validation engines
// on one router. They all nest under /modules/:id (with the module engine also
// registering the static /modules/reorder), so this proves the combined route
// tree builds without a Gin panic.
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

	moduleengine.NewModule(nil, noop, noop).RegisterRoutes(api)
	fieldengine.NewModule(nil, noop, noop).RegisterRoutes(api)
	NewModule(nil, noop, noop).RegisterRoutes(api)
}
