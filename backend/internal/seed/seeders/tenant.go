package seeders

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
	"golang.org/x/crypto/bcrypt"
)

// getOrgID looks up the demo organization id by slug. It is used by every
// seeder that runs after OrganizationSeeder so seeders stay decoupled.
func getOrgID(ctx context.Context, db *pgxpool.Pool) (string, error) {
	var id string
	err := db.QueryRow(ctx,
		`SELECT id FROM organizations WHERE slug = $1`, demoOrgSlug,
	).Scan(&id)
	if err != nil {
		return "", fmt.Errorf("lookup org %q: %w", demoOrgSlug, err)
	}
	return id, nil
}

// ---------------------------------------------------------------------------
// Organization
// ---------------------------------------------------------------------------

type OrganizationSeeder struct{}

func (OrganizationSeeder) Name() string { return "organization" }

func (OrganizationSeeder) Run(ctx context.Context, db *pgxpool.Pool) error {
	_, err := db.Exec(ctx, `
		INSERT INTO organizations (name, slug, plan)
		VALUES ($1, $2, 'pro')
		ON CONFLICT (slug) DO NOTHING
	`, demoOrgName, demoOrgSlug)
	return err
}

// ---------------------------------------------------------------------------
// Roles + role_permissions
// ---------------------------------------------------------------------------

type roleDef struct {
	Slug        string
	Name        string
	Description string
	Perms       []string // permission keys; nil sentinel {"*"} means all
}

var roleDefs = []roleDef{
	{"admin", "Administrator", "Full access to everything", []string{"*"}},
	{"manager", "Sales Manager", "Manages records, imports, exports and automation", []string{
		"module.view", "record.view", "record.create", "record.update", "record.delete",
		"import.run", "export.run", "automation.manage", "validation.manage",
	}},
	{"sales_rep", "Sales Representative", "Works day-to-day with records", []string{
		"module.view", "record.view", "record.create", "record.update", "export.run",
	}},
	{"viewer", "Viewer", "Read-only access", []string{
		"module.view", "record.view",
	}},
}

type RolesSeeder struct{}

func (RolesSeeder) Name() string { return "roles" }

func (RolesSeeder) Run(ctx context.Context, db *pgxpool.Pool) error {
	orgID, err := getOrgID(ctx, db)
	if err != nil {
		return err
	}

	for _, rd := range roleDefs {
		var roleID string
		err := db.QueryRow(ctx, `
			INSERT INTO roles (organization_id, name, slug, description, is_system)
			VALUES ($1, $2, $3, $4, TRUE)
			ON CONFLICT (organization_id, slug)
			DO UPDATE SET name = EXCLUDED.name
			RETURNING id
		`, orgID, rd.Name, rd.Slug, rd.Description).Scan(&roleID)
		if err != nil {
			return fmt.Errorf("upsert role %q: %w", rd.Slug, err)
		}

		if len(rd.Perms) == 1 && rd.Perms[0] == "*" {
			// Grant every permission in the catalog.
			_, err = db.Exec(ctx, `
				INSERT INTO role_permissions (role_id, permission_id)
				SELECT $1, p.id FROM permissions p
				ON CONFLICT DO NOTHING
			`, roleID)
		} else {
			_, err = db.Exec(ctx, `
				INSERT INTO role_permissions (role_id, permission_id)
				SELECT $1, p.id FROM permissions p WHERE p.key = ANY($2)
				ON CONFLICT DO NOTHING
			`, roleID, rd.Perms)
		}
		if err != nil {
			return fmt.Errorf("grant perms to role %q: %w", rd.Slug, err)
		}
	}
	return nil
}

// ---------------------------------------------------------------------------
// Users
// ---------------------------------------------------------------------------

type userDef struct {
	Name     string
	Email    string
	Password string
	Role     string // role slug for membership
}

var userDefs = []userDef{
	{"Aarav Sharma", "admin@crmlite.com", "Admin@12345", "admin"},
	{"Priya Verma", "priya@crmlite.com", "Password@123", "manager"},
	{"Vikram Reddy", "vikram@crmlite.com", "Password@123", "sales_rep"},
	{"Sneha Iyer", "sneha@crmlite.com", "Password@123", "sales_rep"},
	{"Arjun Nair", "arjun@crmlite.com", "Password@123", "viewer"},
}

type UsersSeeder struct{}

func (UsersSeeder) Name() string { return "users" }

func (UsersSeeder) Run(ctx context.Context, db *pgxpool.Pool) error {
	for _, u := range userDefs {
		hash, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
		if err != nil {
			return err
		}
		_, err = db.Exec(ctx, `
			INSERT INTO users (name, email, password_hash)
			VALUES ($1, $2, $3)
			ON CONFLICT (email) DO NOTHING
		`, u.Name, u.Email, string(hash))
		if err != nil {
			return fmt.Errorf("insert user %q: %w", u.Email, err)
		}
	}
	return nil
}

// ---------------------------------------------------------------------------
// Memberships (users <-> organization with a role)
// ---------------------------------------------------------------------------

type MembershipsSeeder struct{}

func (MembershipsSeeder) Name() string { return "memberships" }

func (MembershipsSeeder) Run(ctx context.Context, db *pgxpool.Pool) error {
	orgID, err := getOrgID(ctx, db)
	if err != nil {
		return err
	}

	for _, u := range userDefs {
		_, err := db.Exec(ctx, `
			INSERT INTO organization_members (organization_id, user_id, role_id)
			SELECT $1, usr.id, rol.id
			FROM users usr
			JOIN roles rol ON rol.organization_id = $1 AND rol.slug = $3
			WHERE usr.email = $2
			ON CONFLICT (organization_id, user_id) DO NOTHING
		`, orgID, u.Email, u.Role)
		if err != nil {
			return fmt.Errorf("membership for %q: %w", u.Email, err)
		}
	}
	return nil
}
