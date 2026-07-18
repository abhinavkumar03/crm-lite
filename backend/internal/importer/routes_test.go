package importer

import (
	"testing"

	"github.com/gin-gonic/gin"

	"github.com/abhinavkumar03/crm-lite/backend/internal/rbac"

	fieldengine "github.com/abhinavkumar03/crm-lite/backend/internal/field"
	"github.com/abhinavkumar03/crm-lite/backend/internal/jobs"
	moduleengine "github.com/abhinavkumar03/crm-lite/backend/internal/module"
	"github.com/abhinavkumar03/crm-lite/backend/internal/record"
	"github.com/abhinavkumar03/crm-lite/backend/internal/validationengine"
	"github.com/abhinavkumar03/crm-lite/backend/internal/view"
)

// TestRegisterRoutesNoConflict mounts every engine nesting under /modules/:id
// (including the import routes) on one router, verifying the combined route tree
// builds without a Gin panic — in particular that the static "analyze" and
// ":importId" siblings coexist with the other engines' param routes.
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

	// asynq clients connect lazily, so building a producer here needs no Redis.
	producer := jobs.NewProducer(jobs.RedisOpt("localhost", "6379", "", 0))
	defer producer.Close()

	moduleengine.NewModule(nil, noop, noop, noop, rbac.New(nil)).RegisterRoutes(api)
	fieldengine.NewModule(nil, noop, noop, noop, rbac.New(nil)).RegisterRoutes(api)
	validationengine.NewModule(nil, noop, noop, noop, rbac.New(nil)).RegisterRoutes(api)
	view.NewModule(nil, noop, noop, noop, rbac.New(nil)).RegisterRoutes(api)
	record.NewModule(nil, noop, noop, noop, rbac.New(nil)).RegisterRoutes(api)
	NewModule(nil, noop, noop, noop, rbac.New(nil), producer).RegisterRoutes(api)
}
