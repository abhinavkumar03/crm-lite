package rbac

import (
	"github.com/gin-gonic/gin"

	"github.com/abhinavkumar03/crm-lite/backend/internal/shared/response"
)

// Require returns middleware that aborts with 403 unless the caller has perm.
// It must run after Load().
func (g *Guard) Require(perm string) gin.HandlerFunc {
	return func(c *gin.Context) {
		if !Has(c, perm) {
			response.Forbidden(c, "Missing permission: "+perm)
			c.Abort()
			return
		}
		c.Next()
	}
}

// RequireAny allows the request if the caller has at least one of the perms.
func (g *Guard) RequireAny(perms ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		for _, p := range perms {
			if Has(c, p) {
				c.Next()
				return
			}
		}
		response.Forbidden(c, "Insufficient permissions")
		c.Abort()
	}
}

// RequireModule checks both the global record.* permission and the per-module
// ACL (role_module_access). The module id is read from the ":id" path param.
// Action is one of view|create|update|delete.
func (g *Guard) RequireModule(action ModuleAction) gin.HandlerFunc {
	global := map[ModuleAction]string{
		ActionView:   PermRecordView,
		ActionCreate: PermRecordCreate,
		ActionUpdate: PermRecordUpdate,
		ActionDelete: PermRecordDelete,
	}

	return func(c *gin.Context) {
		perm, ok := global[action]
		if !ok {
			response.InternalServerError(c, "Unknown module action")
			c.Abort()
			return
		}
		if !Has(c, perm) {
			response.Forbidden(c, "Missing permission: "+perm)
			c.Abort()
			return
		}

		moduleID := c.Param("id")
		if moduleID == "" {
			c.Next()
			return
		}

		allowed, err := g.ModuleAllowed(c.Request.Context(), c, moduleID, action)
		if err != nil {
			response.InternalServerError(c, "Unable to check module access")
			c.Abort()
			return
		}
		if !allowed {
			response.Forbidden(c, "Module access denied")
			c.Abort()
			return
		}
		c.Next()
	}
}
