package tenant

import (
	"github.com/gin-gonic/gin"

	"github.com/abhinavkumar03/crm-lite/backend/internal/shared/response"
)

// Context keys for role information, set alongside the organization id so
// downstream RBAC checks (a later phase) can read them without another query.
const (
	ContextRoleID   = "roleID"
	ContextRoleSlug = "roleSlug"
)

// Middleware resolves the authenticated user's organization and injects it into
// the request context. It must run after the auth middleware (which sets
// "userID").
func Middleware(resolver *Resolver) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := c.GetString("userID")
		if userID == "" {
			response.Unauthorized(c, "Authentication required")
			c.Abort()
			return
		}

		membership, err := resolver.MembershipForUser(c.Request.Context(), userID)
		if err != nil {
			response.InternalServerError(c, "Unable to resolve organization")
			c.Abort()
			return
		}
		if membership == nil {
			response.Forbidden(c, "User does not belong to an organization")
			c.Abort()
			return
		}

		c.Set(ContextOrgID, membership.OrganizationID)
		c.Set(ContextRoleID, membership.RoleID)
		c.Set(ContextRoleSlug, membership.RoleSlug)
		c.Next()
	}
}

// OrgID is a small helper to read the resolved organization id from context.
func OrgID(c *gin.Context) string {
	return c.GetString(ContextOrgID)
}

// RoleID returns the caller's role id (empty if the membership has no role).
func RoleID(c *gin.Context) string {
	return c.GetString(ContextRoleID)
}

// RoleSlug returns the caller's role slug (e.g. "admin", "viewer").
func RoleSlug(c *gin.Context) string {
	return c.GetString(ContextRoleSlug)
}
