package notification

import (
	"testing"

	"github.com/gin-gonic/gin"

	"github.com/abhinavkumar03/crm-lite/backend/internal/rbac"

	"github.com/abhinavkumar03/crm-lite/backend/internal/jobs"
)

// TestRegisterRoutesNoConflict verifies the notification routes register without
// a Gin panic.
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

	// asynq clients connect lazily, so building a producer here does not require
	// a live Redis.
	producer := jobs.NewProducer(jobs.RedisOpt("localhost", "6379", "", 0))
	defer producer.Close()

	NewModule(nil, noop, noop, noop, rbac.New(nil), producer).RegisterRoutes(api)
}
