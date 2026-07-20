package bootstrap

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/abhinavkumar03/crm-lite/backend/internal/rbac"
	settingsentity "github.com/abhinavkumar03/crm-lite/backend/internal/settings/entity"
)

// RoleSpec describes a system role seeded for every new organization.
type RoleSpec struct {
	Slug           string
	Name           string
	Description    string
	HierarchyLevel int
	AllPermissions bool
	PermissionKeys []string
}

var DefaultRoles = []RoleSpec{
	{Slug: "owner", Name: "Owner", Description: "Organization owner", HierarchyLevel: 0, AllPermissions: true},
	{Slug: "super_admin", Name: "Super Admin", Description: "Full administrative access", HierarchyLevel: 10, AllPermissions: true},
	{Slug: "admin", Name: "Admin", Description: "Organization administrator", HierarchyLevel: 20, AllPermissions: true},
	{Slug: "sales_manager", Name: "Sales Manager", Description: "Manages sales team and records", HierarchyLevel: 40, PermissionKeys: []string{
		"module.view", "record.view", "record.create", "record.update", "record.delete",
		"import.run", "export.run", "automation.manage", "validation.manage", "analytics.view", "user.manage",
	}},
	{Slug: "sales_executive", Name: "Sales Executive", Description: "Day-to-day sales work", HierarchyLevel: 60, PermissionKeys: []string{
		"module.view", "record.view", "record.create", "record.update", "export.run",
	}},
	{Slug: "viewer", Name: "Viewer", Description: "Read-only access", HierarchyLevel: 100, PermissionKeys: []string{
		"module.view", "record.view", "analytics.view",
	}},
}

// CreateOptions carries optional profile + locale prefs for a new workspace.
type CreateOptions struct {
	Name        string
	Slug        string
	Description string
	Industry    string
	CompanySize string
	Country     string
	LogoURL     string
	Timezone    string
	Currency    string
	Locale      string
	DateFormat  string
}

type Service struct {
	db *pgxpool.Pool
}

func New(db *pgxpool.Pool) *Service {
	return &Service{db: db}
}

// HasMembership reports whether the user already belongs to any organization.
func (s *Service) HasMembership(ctx context.Context, userID string) (bool, error) {
	var exists bool
	err := s.db.QueryRow(ctx, `
		SELECT EXISTS(
			SELECT 1 FROM organization_members
			WHERE user_id = $1 AND status = 'active'
		)
	`, userID).Scan(&exists)
	return exists, err
}

// CreateOrganization inserts the org, seeds roles/settings/modules, and adds
// the creator as Owner with active_organization_id set.
func (s *Service) CreateOrganization(ctx context.Context, opts CreateOptions, creatorUserID string) (orgID string, err error) {
	name := strings.TrimSpace(opts.Name)
	if name == "" {
		return "", errors.New("name required")
	}

	slug := normalizeSlug(opts.Slug)
	if slug == "" {
		slug = normalizeSlug(name)
	}
	if slug == "" {
		slug = "workspace"
	}

	general := settingsentity.DefaultGeneral()
	if tz := strings.TrimSpace(opts.Timezone); tz != "" {
		general.Timezone = tz
	}
	if cur := strings.TrimSpace(opts.Currency); cur != "" {
		general.Currency = cur
	}
	if loc := strings.TrimSpace(opts.Locale); loc != "" {
		general.Locale = loc
	}
	if df := strings.TrimSpace(opts.DateFormat); df != "" {
		general.DateFormat = df
	}

	settings, err := json.Marshal(map[string]any{
		"general":    general,
		"automation": settingsentity.DefaultAutomation(),
	})
	if err != nil {
		return "", err
	}

	var industry, companySize, country, logoURL, description *string
	if v := strings.TrimSpace(opts.Industry); v != "" {
		industry = &v
	}
	if v := strings.TrimSpace(opts.CompanySize); v != "" {
		companySize = &v
	}
	if v := strings.TrimSpace(opts.Country); v != "" {
		country = &v
	}
	if v := strings.TrimSpace(opts.LogoURL); v != "" {
		logoURL = &v
	}
	if v := strings.TrimSpace(opts.Description); v != "" {
		description = &v
	}

	for attempt := 0; attempt < 5; attempt++ {
		candidate := slug
		if attempt > 0 {
			candidate = slug + "-" + shortSuffix()
		}
		err = s.db.QueryRow(ctx, `
			INSERT INTO organizations (
				name, slug, plan, status, created_by, settings,
				industry, company_size, country, logo_url, description
			)
			VALUES ($1, $2, 'free', 'active', $3, $4, $5, $6, $7, $8, $9)
			RETURNING id
		`, name, candidate, creatorUserID, settings, industry, companySize, country, logoURL, description).Scan(&orgID)
		if err == nil {
			break
		}
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			continue
		}
		return "", fmt.Errorf("create organization: %w", err)
	}
	if orgID == "" {
		return "", fmt.Errorf("create organization: slug unavailable")
	}

	// Catalog must exist before role grants; seed may not have been run yet.
	if err := s.EnsurePermissionCatalog(ctx); err != nil {
		return "", err
	}

	roleIDs, err := s.seedRoles(ctx, orgID)
	if err != nil {
		return "", err
	}

	if err := s.seedDefaultModules(ctx, orgID); err != nil {
		return "", err
	}

	ownerRole := roleIDs["owner"]
	_, err = s.db.Exec(ctx, `
		INSERT INTO organization_members (organization_id, user_id, role_id, status, hierarchy_level)
		VALUES ($1, $2, $3, 'active', 0)
	`, orgID, creatorUserID, ownerRole)
	if err != nil {
		return "", fmt.Errorf("add owner membership: %w", err)
	}

	_, err = s.db.Exec(ctx, `
		UPDATE users SET active_organization_id = $2, updated_at = NOW() WHERE id = $1
	`, creatorUserID, orgID)
	if err != nil {
		return "", err
	}

	return orgID, nil
}

// EnsureFullCatalog upserts the full 5-module field/layout catalog for an org.
// Safe to call on existing workspaces created with an older minimal bootstrap.
func (s *Service) EnsureFullCatalog(ctx context.Context, orgID string) error {
	return s.seedDefaultModules(ctx, orgID)
}

func (s *Service) seedRoles(ctx context.Context, orgID string) (map[string]string, error) {
	ids := make(map[string]string, len(DefaultRoles))
	for _, rd := range DefaultRoles {
		var roleID string
		err := s.db.QueryRow(ctx, `
			INSERT INTO roles (organization_id, name, slug, description, is_system, hierarchy_level)
			VALUES ($1, $2, $3, $4, TRUE, $5)
			RETURNING id
		`, orgID, rd.Name, rd.Slug, rd.Description, rd.HierarchyLevel).Scan(&roleID)
		if err != nil {
			return nil, fmt.Errorf("seed role %s: %w", rd.Slug, err)
		}
		ids[rd.Slug] = roleID

		if err := s.grantRolePermissions(ctx, roleID, rd); err != nil {
			return nil, fmt.Errorf("grant role %s: %w", rd.Slug, err)
		}
	}
	return ids, nil
}

// EnsurePermissionCatalog upserts every known permission key so Owner/Admin
// grants never land on an empty catalog (common when migrate ran but seed did not).
func (s *Service) EnsurePermissionCatalog(ctx context.Context) error {
	for _, p := range rbac.PermissionCatalog {
		_, err := s.db.Exec(ctx, `
			INSERT INTO permissions (key, category, description)
			VALUES ($1, $2, $3)
			ON CONFLICT (key) DO UPDATE
			SET category = EXCLUDED.category,
			    description = EXCLUDED.description
		`, p.Key, p.Category, p.Desc)
		if err != nil {
			return fmt.Errorf("ensure permission %s: %w", p.Key, err)
		}
	}
	return nil
}

// RepairFullAccessRoles grants the full catalog to Owner / Super Admin / Admin
// for every organization. Safe to run repeatedly (ON CONFLICT DO NOTHING).
// Returns the role IDs that were targeted so callers can invalidate RBAC cache.
func (s *Service) RepairFullAccessRoles(ctx context.Context) ([]string, error) {
	if err := s.EnsurePermissionCatalog(ctx); err != nil {
		return nil, err
	}
	_, err := s.db.Exec(ctx, `
		INSERT INTO role_permissions (role_id, permission_id)
		SELECT r.id, p.id
		FROM roles r
		CROSS JOIN permissions p
		WHERE r.is_system = TRUE
		  AND r.slug IN ('owner', 'super_admin', 'admin')
		ON CONFLICT DO NOTHING
	`)
	if err != nil {
		return nil, fmt.Errorf("repair full-access roles: %w", err)
	}

	rows, err := s.db.Query(ctx, `
		SELECT id::text FROM roles
		WHERE is_system = TRUE AND slug IN ('owner', 'super_admin', 'admin')
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	ids := make([]string, 0)
	for rows.Next() {
		var id string
		if err := rows.Scan(&id); err != nil {
			return nil, err
		}
		ids = append(ids, id)
	}
	return ids, rows.Err()
}

func (s *Service) grantRolePermissions(ctx context.Context, roleID string, rd RoleSpec) error {
	if rd.AllPermissions {
		_, err := s.db.Exec(ctx, `
			INSERT INTO role_permissions (role_id, permission_id)
			SELECT $1, p.id FROM permissions p
			ON CONFLICT DO NOTHING
		`, roleID)
		return err
	}
	_, err := s.db.Exec(ctx, `
		INSERT INTO role_permissions (role_id, permission_id)
		SELECT $1, p.id FROM permissions p WHERE p.key = ANY($2)
		ON CONFLICT DO NOTHING
	`, roleID, rd.PermissionKeys)
	return err
}

func normalizeSlug(s string) string {
	s = strings.ToLower(strings.TrimSpace(s))
	var b strings.Builder
	prevDash := false
	for _, r := range s {
		ok := (r >= 'a' && r <= 'z') || (r >= '0' && r <= '9')
		if ok {
			b.WriteRune(r)
			prevDash = false
			continue
		}
		if !prevDash && b.Len() > 0 {
			b.WriteByte('-')
			prevDash = true
		}
	}
	out := strings.Trim(b.String(), "-")
	if len(out) > 80 {
		out = out[:80]
	}
	return out
}

func shortSuffix() string {
	var b [3]byte
	_, _ = rand.Read(b[:])
	return hex.EncodeToString(b[:])
}
