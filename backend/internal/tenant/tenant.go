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

const ContextOrgID = "orgID"

// HeaderOrganizationID is an optional override for multi-org users.
const HeaderOrganizationID = "X-Organization-Id"

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

// MembershipForUser returns the user's active membership.
// Preference: users.active_organization_id (if membership active), else oldest
// active membership. Optional preferredOrgID (from header) wins when the user
// is a member of that org.
func (r *Resolver) MembershipForUser(ctx context.Context, userID, preferredOrgID string) (*Membership, error) {
	if preferredOrgID != "" {
		m, err := r.membershipInOrg(ctx, userID, preferredOrgID)
		if err != nil {
			return nil, err
		}
		if m != nil {
			return m, nil
		}
	}

	cacheKey := cache.MembershipKey(userID)
	var cached Membership
	if preferredOrgID == "" && r.cache.GetJSON(ctx, cacheKey, &cached) {
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
		LEFT JOIN users u ON u.id = om.user_id
		WHERE om.user_id = $1 AND om.status = 'active'
		ORDER BY
		  CASE WHEN u.active_organization_id IS NOT NULL
		            AND om.organization_id = u.active_organization_id THEN 0
		       ELSE 1 END,
		  om.created_at ASC
		LIMIT 1
	`, userID).Scan(&m.OrganizationID, &m.RoleID, &m.RoleSlug)

	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("tenant: resolve membership: %w", err)
	}

	if preferredOrgID == "" {
		r.cache.SetJSON(ctx, cacheKey, m, cache.TTLShort)
	}
	return &m, nil
}

func (r *Resolver) membershipInOrg(ctx context.Context, userID, orgID string) (*Membership, error) {
	var m Membership
	err := r.db.QueryRow(ctx, `
		SELECT om.organization_id,
		       COALESCE(om.role_id::text, ''),
		       COALESCE(rl.slug, '')
		FROM organization_members om
		LEFT JOIN roles rl ON rl.id = om.role_id
		WHERE om.user_id = $1
		  AND om.organization_id = $2
		  AND om.status = 'active'
	`, userID, orgID).Scan(&m.OrganizationID, &m.RoleID, &m.RoleSlug)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &m, nil
}

// ListMemberships returns all active org memberships for a user.
func (r *Resolver) ListMemberships(ctx context.Context, userID string) ([]Membership, error) {
	rows, err := r.db.Query(ctx, `
		SELECT om.organization_id,
		       COALESCE(om.role_id::text, ''),
		       COALESCE(rl.slug, ''),
		       o.name,
		       o.slug
		FROM organization_members om
		JOIN organizations o ON o.id = om.organization_id
		LEFT JOIN roles rl ON rl.id = om.role_id
		WHERE om.user_id = $1 AND om.status = 'active'
		ORDER BY o.name ASC
	`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	type row struct {
		Membership
		Name string
		Slug string
	}
	out := make([]Membership, 0)
	for rows.Next() {
		var m Membership
		var name, slug string
		if err := rows.Scan(&m.OrganizationID, &m.RoleID, &m.RoleSlug, &name, &slug); err != nil {
			return nil, err
		}
		_ = name
		_ = slug
		out = append(out, m)
	}
	return out, rows.Err()
}

// SetActiveOrganization updates users.active_organization_id after validating membership.
func (r *Resolver) SetActiveOrganization(ctx context.Context, userID, orgID string) error {
	m, err := r.membershipInOrg(ctx, userID, orgID)
	if err != nil {
		return err
	}
	if m == nil {
		return ErrNotMember
	}
	_, err = r.db.Exec(ctx, `
		UPDATE users SET active_organization_id = $2, updated_at = NOW() WHERE id = $1
	`, userID, orgID)
	if err != nil {
		return err
	}
	r.cache.InvalidateMembership(ctx, userID)
	return nil
}

var ErrNotMember = errors.New("not a member of organization")
