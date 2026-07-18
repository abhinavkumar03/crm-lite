// Package cache is a thin Redis JSON cache used across the API for hot, mostly
// read-only data (dashboard metrics, tenant membership, RBAC grants). Callers
// treat a cache miss or Redis error as transparent — they always fall through
// to the source of truth.
package cache

import (
	"context"
	"encoding/json"
	"strings"
	"time"

	"github.com/redis/go-redis/v9"
)

// Default TTLs for the common cache families.
const (
	TTLShort  = 2 * time.Minute
	TTLMedium = 5 * time.Minute
	TTLLong   = 15 * time.Minute
)

// Cache wraps a go-redis client with JSON helpers and key builders.
type Cache struct {
	rdb *redis.Client
}

func New(rdb *redis.Client) *Cache {
	if rdb == nil {
		return nil
	}
	return &Cache{rdb: rdb}
}

// Key joins parts with ':' into a namespaced Redis key.
func Key(parts ...string) string {
	return strings.Join(parts, ":")
}

// GetJSON unmarshals a cached value into dest. It returns true on a hit.
func (c *Cache) GetJSON(ctx context.Context, key string, dest any) bool {
	if c == nil || c.rdb == nil {
		return false
	}
	raw, err := c.rdb.Get(ctx, key).Bytes()
	if err != nil {
		return false
	}
	return json.Unmarshal(raw, dest) == nil
}

// SetJSON marshals value and stores it with the given TTL. Failures are ignored
// so a flaky Redis never breaks the request path.
func (c *Cache) SetJSON(ctx context.Context, key string, value any, ttl time.Duration) {
	if c == nil || c.rdb == nil {
		return
	}
	raw, err := json.Marshal(value)
	if err != nil {
		return
	}
	_ = c.rdb.Set(ctx, key, raw, ttl).Err()
}

// Delete removes one or more keys. Failures are ignored.
func (c *Cache) Delete(ctx context.Context, keys ...string) {
	if c == nil || c.rdb == nil || len(keys) == 0 {
		return
	}
	_ = c.rdb.Del(ctx, keys...).Err()
}

// ---------------------------------------------------------------------------
// Well-known key helpers
// ---------------------------------------------------------------------------

func DashboardKey(userID string) string { return Key("dashboard", userID) }

func MembershipKey(userID string) string { return Key("tenant", "membership", userID) }

func PermissionsKey(roleID string) string { return Key("rbac", "perms", roleID) }

func ModuleAccessKey(roleID string) string { return Key("rbac", "module", roleID) }

func FieldAccessKey(roleID, moduleID string) string {
	return Key("rbac", "field", roleID, moduleID)
}

// InvalidateRole drops every RBAC cache entry for a role (permissions + ACL).
func (c *Cache) InvalidateRole(ctx context.Context, roleID string) {
	if c == nil || roleID == "" {
		return
	}
	c.Delete(ctx, PermissionsKey(roleID), ModuleAccessKey(roleID))
	// Field ACL is per (role, module); scan is avoided — callers that change
	// field ACL also pass module ids, or we use a short TTL. Drop a wildcard
	// via SCAN only if needed; for portfolio scale, TTL (2m) is enough and
	// SetFieldAccess invalidates the known module keys below.
}

// InvalidateRoleFieldAccess drops field ACL for the given role+module pairs.
func (c *Cache) InvalidateRoleFieldAccess(ctx context.Context, roleID string, moduleIDs ...string) {
	if c == nil || roleID == "" {
		return
	}
	keys := make([]string, 0, len(moduleIDs))
	for _, mid := range moduleIDs {
		if mid != "" {
			keys = append(keys, FieldAccessKey(roleID, mid))
		}
	}
	c.Delete(ctx, keys...)
}

// InvalidateMembership drops the cached membership for a user.
func (c *Cache) InvalidateMembership(ctx context.Context, userID string) {
	if userID == "" {
		return
	}
	c.Delete(ctx, MembershipKey(userID))
}

// InvalidateDashboard drops the cached dashboard for a user.
func (c *Cache) InvalidateDashboard(ctx context.Context, userID string) {
	if userID == "" {
		return
	}
	c.Delete(ctx, DashboardKey(userID))
}
