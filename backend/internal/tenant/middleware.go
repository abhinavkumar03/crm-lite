package tenant

import (
	"github.com/gin-gonic/gin"

	"github.com/abhinavkumar03/crm-lite/backend/internal/shared/response"
)

const (
	ContextRoleID   = "roleID"
	ContextRoleSlug = "roleSlug"
)

// Middleware resolves the authenticated user's organization and injects it into
// the request context. It must run after the auth middleware (which sets
// "userID"). Optional X-Organization-Id selects a membership the user already has.
func Middleware(resolver *Resolver) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := c.GetString("userID")
		if userID == "" {
			response.Unauthorized(c, "Authentication required")
			c.Abort()
			return
		}

		preferred := c.GetHeader(HeaderOrganizationID)
		membership, err := resolver.MembershipForUser(c.Request.Context(), userID, preferred)
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

func OrgID(c *gin.Context) string {
	return c.GetString(ContextOrgID)
}

func RoleID(c *gin.Context) string {
	return c.GetString(ContextRoleID)
}

func RoleSlug(c *gin.Context) string {
	return c.GetString(ContextRoleSlug)
}
