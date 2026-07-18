package roles

import (
	"testing"

	"github.com/gin-gonic/gin"

	"github.com/abhinavkumar03/crm-lite/backend/internal/rbac"
)

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

	NewModule(nil, nil, noop, noop, noop, rbac.New(nil, nil)).RegisterRoutes(api)
}
