// Package tenant resolves the active organization for an authenticated user and
// exposes a Gin middleware that injects the organization id into the request
// context. Every metadata/runtime module is organization-scoped, so this is the
// single place tenancy is resolved.
package tenant

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/abhinavkumar03/crm-lite/backend/internal/shared/cache"
)

// ContextOrgID is the Gin context key under which the resolved organization id
// is stored.
const ContextOrgID = "orgID"

// Membership is the resolved tenant context for a user.
type Membership struct {
	OrganizationID string `json:"organization_id"`
	RoleID         string `json:"role_id"`
	RoleSlug       string `json:"role_slug"`
}

type Resolver struct {
	db    *pgxpool.Pool
	cache *cache.Cache
}

func NewResolver(db *pgxpool.Pool, c *cache.Cache) *Resolver {
	return &Resolver{db: db, cache: c}
}

// MembershipForUser returns the user's active membership (organization + role).
// If the user belongs to no organization, it returns (nil, nil). Results are
// cached briefly so the hot auth path does not hit Postgres on every request.
func (r *Resolver) MembershipForUser(ctx context.Context, userID string) (*Membership, error) {
	key := cache.MembershipKey(userID)
	var cached Membership
	if r.cache.GetJSON(ctx, key, &cached) {
		// Distinguish "cached nil membership" via empty OrganizationID — we
		// only cache positive hits; misses always re-query.
		if cached.OrganizationID != "" {
			return &cached, nil
		}
	}

	var m Membership
	err := r.db.QueryRow(ctx, `
		SELECT om.organization_id,
		       COALESCE(om.role_id::text, ''),
		       COALESCE(rl.slug, '')
		FROM organization_members om
		LEFT JOIN roles rl ON rl.id = om.role_id
		WHERE om.user_id = $1 AND om.status = 'active'
		ORDER BY om.created_at ASC
		LIMIT 1
	`, userID).Scan(&m.OrganizationID, &m.RoleID, &m.RoleSlug)

	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("tenant: resolve membership: %w", err)
	}

	r.cache.SetJSON(ctx, key, m, cache.TTLShort)
	return &m, nil
}
