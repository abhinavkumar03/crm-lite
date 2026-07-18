package field

import (
	"testing"

	"github.com/gin-gonic/gin"

	"github.com/abhinavkumar03/crm-lite/backend/internal/rbac"

	moduleengine "github.com/abhinavkumar03/crm-lite/backend/internal/module"
)

// TestRegisterRoutesNoConflict registers both the module engine and the field
// engine on the same router. Fields nest under /modules/:id/fields while the
// module engine registers /modules/:id and the static /modules/reorder, so this
// verifies Gin's tree tolerates the static/param siblings without panicking.
func TestRegisterRoutesNoConflict(t *testing.T) {
	gin.SetMode(gin.TestMode)

	noop := func(c *gin.Context) { c.Next() }
	fieldModule := NewModule(nil, noop, noop, noop, rbac.New(nil))
	moduleModule := moduleengine.NewModule(nil, noop, noop, noop, rbac.New(nil))

	router := gin.New()
	api := router.Group("/api/v1")

	defer func() {
		if r := recover(); r != nil {
			t.Fatalf("route registration panicked: %v", r)
		}
	}()

	moduleModule.RegisterRoutes(api)
	fieldModule.RegisterRoutes(api)
}
