package seeders

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

// permission is a single entry in the global permission catalog.
type permission struct {
	Key      string
	Category string
	Desc     string
}

// permissionCatalog is the full set of permission keys the app understands.
// Roles are granted subsets of these (see RolesSeeder).
var permissionCatalog = []permission{
	{"module.view", "module", "View modules and navigation"},
	{"module.manage", "module", "Create, update, reorder and delete modules"},
	{"field.manage", "field", "Create, update and delete fields"},
	{"record.view", "record", "View records"},
	{"record.create", "record", "Create records"},
	{"record.update", "record", "Update records"},
	{"record.delete", "record", "Delete records"},
	{"import.run", "import", "Run data imports"},
	{"export.run", "export", "Run data exports"},
	{"automation.manage", "automation", "Create and manage automation rules"},
	{"validation.manage", "validation", "Create and manage validation rules"},
	{"settings.manage", "settings", "Manage organization settings"},
	{"user.manage", "user", "Invite and manage users"},
	{"role.manage", "role", "Create and manage roles and permissions"},
	{"organization.manage", "organization", "Create organizations and manage membership structure"},
	{"analytics.view", "analytics", "View dashboards and analytics"},
}

// PermissionsSeeder populates the global permission catalog.
type PermissionsSeeder struct{}

func (PermissionsSeeder) Name() string { return "permissions_catalog" }

func (PermissionsSeeder) Run(ctx context.Context, db *pgxpool.Pool) error {
	for _, p := range permissionCatalog {
		_, err := db.Exec(ctx, `
			INSERT INTO permissions (key, category, description)
			VALUES ($1, $2, $3)
			ON CONFLICT (key) DO NOTHING
		`, p.Key, p.Category, p.Desc)
		if err != nil {
			return err
		}
	}
	return nil
}
