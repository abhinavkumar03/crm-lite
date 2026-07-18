package rbac

import (
	"context"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/abhinavkumar03/crm-lite/backend/internal/shared/response"
	"github.com/abhinavkumar03/crm-lite/backend/internal/tenant"
)

// Guard loads a role's permission keys and answers ACL queries. It is the
// single composition root for RBAC middleware and access checks.
type Guard struct {
	db *pgxpool.Pool
}

func New(db *pgxpool.Pool) *Guard {
	return &Guard{db: db}
}

// Load is middleware that must run after tenant.Middleware. It fetches the
// role's permission keys once per request and stores them on the Gin context.
func (g *Guard) Load() gin.HandlerFunc {
	return func(c *gin.Context) {
		roleID := tenant.RoleID(c)
		if roleID == "" {
			c.Set(ContextPermissions, []string{})
			c.Next()
			return
		}

		perms, err := g.PermissionsForRole(c.Request.Context(), roleID)
		if err != nil {
			response.InternalServerError(c, "Unable to resolve permissions")
			c.Abort()
			return
		}
		c.Set(ContextPermissions, perms)
		c.Next()
	}
}

// PermissionsForRole returns every permission key granted to the role.
func (g *Guard) PermissionsForRole(ctx context.Context, roleID string) ([]string, error) {
	rows, err := g.db.Query(ctx, `
		SELECT p.key
		FROM role_permissions rp
		JOIN permissions p ON p.id = rp.permission_id
		WHERE rp.role_id = $1
		ORDER BY p.key
	`, roleID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	keys := make([]string, 0)
	for rows.Next() {
		var key string
		if err := rows.Scan(&key); err != nil {
			return nil, err
		}
		keys = append(keys, key)
	}
	return keys, rows.Err()
}

// Permissions reads the loaded keys from context (empty if Load was not run).
func Permissions(c *gin.Context) []string {
	v, ok := c.Get(ContextPermissions)
	if !ok {
		return nil
	}
	keys, _ := v.([]string)
	return keys
}

// Has reports whether the request carries the given permission key.
func Has(c *gin.Context, perm string) bool {
	for _, k := range Permissions(c) {
		if k == perm {
			return true
		}
	}
	return false
}
