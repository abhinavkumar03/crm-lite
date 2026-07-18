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
	Industry    string
	CompanySize string
	Country     string
	LogoURL     string
	Timezone    string
	Currency    string
	Locale      string
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

	settings, err := json.Marshal(map[string]any{
		"general":    general,
		"automation": settingsentity.DefaultAutomation(),
	})
	if err != nil {
		return "", err
	}

	var industry, companySize, country, logoURL *string
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

	for attempt := 0; attempt < 5; attempt++ {
		candidate := slug
		if attempt > 0 {
			candidate = slug + "-" + shortSuffix()
		}
		err = s.db.QueryRow(ctx, `
			INSERT INTO organizations (
				name, slug, plan, status, created_by, settings,
				industry, company_size, country, logo_url
			)
			VALUES ($1, $2, 'free', 'active', $3, $4, $5, $6, $7, $8)
			RETURNING id
		`, name, candidate, creatorUserID, settings, industry, companySize, country, logoURL).Scan(&orgID)
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

		if rd.AllPermissions {
			_, err = s.db.Exec(ctx, `
				INSERT INTO role_permissions (role_id, permission_id)
				SELECT $1, p.id FROM permissions p
				ON CONFLICT DO NOTHING
			`, roleID)
		} else {
			_, err = s.db.Exec(ctx, `
				INSERT INTO role_permissions (role_id, permission_id)
				SELECT $1, p.id FROM permissions p WHERE p.key = ANY($2)
				ON CONFLICT DO NOTHING
			`, roleID, rd.PermissionKeys)
		}
		if err != nil {
			return nil, fmt.Errorf("grant role %s: %w", rd.Slug, err)
		}
	}
	return ids, nil
}

func (s *Service) seedDefaultModules(ctx context.Context, orgID string) error {
	type mod struct {
		API, Singular, Plural, Icon, Color string
		Sort                               int
	}
	mods := []mod{
		{"company", "Company", "Companies", "building-2", "#8b5cf6", 1},
		{"deal", "Deal", "Deals", "handshake", "#ec4899", 2},
	}
	moduleIDs := map[string]string{}
	for _, m := range mods {
		var id string
		err := s.db.QueryRow(ctx, `
			INSERT INTO modules (
				organization_id, api_name, singular_label, plural_label,
				icon, color, storage_strategy, is_system, sort_order,
				is_enabled, is_visible_sidebar
			) VALUES ($1,$2,$3,$4,$5,$6,'dynamic',TRUE,$7,TRUE,TRUE)
			RETURNING id
		`, orgID, m.API, m.Singular, m.Plural, m.Icon, m.Color, m.Sort).Scan(&id)
		if err != nil {
			return fmt.Errorf("seed module %s: %w", m.API, err)
		}
		moduleIDs[m.API] = id
	}

	type field struct {
		Module, API, Label, Type string
		Required, Searchable     bool
	}
	fields := []field{
		{"company", "name", "Company Name", "text", true, true},
		{"company", "industry", "Industry", "text", false, false},
		{"company", "city", "City", "text", false, true},
		{"deal", "title", "Deal Title", "text", true, true},
		{"deal", "amount", "Amount", "currency", false, false},
		{"deal", "stage", "Stage", "text", false, true},
	}
	for i, f := range fields {
		_, err := s.db.Exec(ctx, `
			INSERT INTO fields (
				organization_id, module_id, api_name, label, field_type,
				is_required, is_searchable, is_filterable, options, sort_order, is_system
			) VALUES ($1,$2,$3,$4,$5,$6,$7,TRUE,'[]'::jsonb,$8,TRUE)
		`, orgID, moduleIDs[f.Module], f.API, f.Label, f.Type, f.Required, f.Searchable, i+1)
		if err != nil {
			return fmt.Errorf("seed field %s.%s: %w", f.Module, f.API, err)
		}
	}
	return nil
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
