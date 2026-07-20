package seeders

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/abhinavkumar03/crm-lite/backend/internal/rbac"
)

// PermissionsSeeder populates the global permission catalog.
type PermissionsSeeder struct{}

func (PermissionsSeeder) Name() string { return "permissions_catalog" }

func (PermissionsSeeder) Run(ctx context.Context, db *pgxpool.Pool) error {
	for _, p := range rbac.PermissionCatalog {
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
