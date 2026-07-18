package seeders

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
	"golang.org/x/crypto/bcrypt"
)

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

type OrganizationSeeder struct{}

func (OrganizationSeeder) Name() string { return "organization" }

func (OrganizationSeeder) Run(ctx context.Context, db *pgxpool.Pool) error {
	_, err := db.Exec(ctx, `
		INSERT INTO organizations (
			name, slug, plan, status, industry, company_size, country, settings
		) VALUES (
			$1, $2, 'pro', 'active', 'Technology', '51-200', 'IN',
			'{"general":{"timezone":"Asia/Kolkata","date_format":"YYYY-MM-DD","time_format":"24h","currency":"INR","locale":"en-IN","week_start":"monday"},"automation":{"notifications_enabled":true,"default_channel":"whatsapp","daily_digest":false}}'::jsonb
		)
		ON CONFLICT (slug) DO UPDATE SET
			industry = EXCLUDED.industry,
			company_size = EXCLUDED.company_size,
			country = EXCLUDED.country,
			status = EXCLUDED.status
	`, demoOrgName, demoOrgSlug)
	return err
}

type roleDef struct {
	Slug        string
	Name        string
	Description string
	Level       int
	Perms       []string
}

var roleDefs = []roleDef{
	{"owner", "Owner", "Organization owner", 0, []string{"*"}},
	{"super_admin", "Super Admin", "Full administrative access", 10, []string{"*"}},
	{"admin", "Administrator", "Organization administrator", 20, []string{"*"}},
	{"sales_manager", "Sales Manager", "Manages sales team and records", 40, []string{
		"module.view", "record.view", "record.create", "record.update", "record.delete",
		"import.run", "export.run", "automation.manage", "validation.manage",
		"analytics.view", "user.manage",
	}},
	{"sales_executive", "Sales Executive", "Day-to-day sales work", 60, []string{
		"module.view", "record.view", "record.create", "record.update", "export.run",
	}},
	{"viewer", "Viewer", "Read-only access", 100, []string{
		"module.view", "record.view", "analytics.view",
	}},
	// Legacy aliases kept for older demos / docs.
	{"manager", "Sales Manager (legacy)", "Alias of sales_manager", 40, []string{
		"module.view", "record.view", "record.create", "record.update", "record.delete",
		"import.run", "export.run", "automation.manage", "validation.manage",
	}},
	{"sales_rep", "Sales Representative (legacy)", "Alias of sales_executive", 60, []string{
		"module.view", "record.view", "record.create", "record.update", "export.run",
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
			INSERT INTO roles (organization_id, name, slug, description, is_system, hierarchy_level)
			VALUES ($1, $2, $3, $4, TRUE, $5)
			ON CONFLICT (organization_id, slug)
			DO UPDATE SET name = EXCLUDED.name, hierarchy_level = EXCLUDED.hierarchy_level
			RETURNING id
		`, orgID, rd.Name, rd.Slug, rd.Description, rd.Level).Scan(&roleID)
		if err != nil {
			return fmt.Errorf("upsert role %q: %w", rd.Slug, err)
		}

		if len(rd.Perms) == 1 && rd.Perms[0] == "*" {
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

type userDef struct {
	Name     string
	Email    string
	Password string
	Role     string
}

var userDefs = []userDef{
	{"Aarav Sharma", "admin@crmlite.com", "Admin@12345", "owner"},
	{"Priya Verma", "priya@crmlite.com", "Password@123", "sales_manager"},
	{"Vikram Reddy", "vikram@crmlite.com", "Password@123", "sales_executive"},
	{"Sneha Iyer", "sneha@crmlite.com", "Password@123", "sales_executive"},
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

type MembershipsSeeder struct{}

func (MembershipsSeeder) Name() string { return "memberships" }

func (MembershipsSeeder) Run(ctx context.Context, db *pgxpool.Pool) error {
	orgID, err := getOrgID(ctx, db)
	if err != nil {
		return err
	}

	for _, u := range userDefs {
		level := 100
		switch u.Role {
		case "owner":
			level = 0
		case "sales_manager", "manager":
			level = 40
		case "sales_executive", "sales_rep":
			level = 60
		case "viewer":
			level = 100
		}
		_, err := db.Exec(ctx, `
			INSERT INTO organization_members (organization_id, user_id, role_id, hierarchy_level)
			SELECT $1, usr.id, rol.id, $4
			FROM users usr
			JOIN roles rol ON rol.organization_id = $1 AND rol.slug = $3
			WHERE usr.email = $2
			ON CONFLICT (organization_id, user_id) DO UPDATE
			SET role_id = EXCLUDED.role_id, hierarchy_level = EXCLUDED.hierarchy_level
		`, orgID, u.Email, u.Role, level)
		if err != nil {
			return fmt.Errorf("membership for %q: %w", u.Email, err)
		}
	}

	// Reporting tree: Priya manages Vikram + Sneha.
	_, err = db.Exec(ctx, `
		UPDATE organization_members om
		SET manager_user_id = mgr.id
		FROM users mgr, users rep
		WHERE om.organization_id = $1
		  AND om.user_id = rep.id
		  AND mgr.email = 'priya@crmlite.com'
		  AND rep.email IN ('vikram@crmlite.com', 'sneha@crmlite.com')
	`, orgID)
	if err != nil {
		return err
	}

	// Active org for all demo users.
	_, err = db.Exec(ctx, `
		UPDATE users SET active_organization_id = $1
		WHERE email = ANY($2)
	`, orgID, []string{
		"admin@crmlite.com", "priya@crmlite.com", "vikram@crmlite.com",
		"sneha@crmlite.com", "arjun@crmlite.com",
	})
	return err
}

// OrgStructureSeeder seeds departments and teams, then attaches members.
type OrgStructureSeeder struct{}

func (OrgStructureSeeder) Name() string { return "org_structure" }

func (OrgStructureSeeder) Run(ctx context.Context, db *pgxpool.Pool) error {
	orgID, err := getOrgID(ctx, db)
	if err != nil {
		return err
	}

	var salesDept string
	err = db.QueryRow(ctx, `
		INSERT INTO departments (organization_id, name, description)
		VALUES ($1, 'Sales', 'Revenue team')
		ON CONFLICT (organization_id, name) DO UPDATE SET description = EXCLUDED.description
		RETURNING id
	`, orgID).Scan(&salesDept)
	if err != nil {
		return err
	}

	var eastTeam string
	err = db.QueryRow(ctx, `
		INSERT INTO teams (organization_id, department_id, name, description)
		VALUES ($1, $2, 'East Region', 'Field sales east')
		ON CONFLICT (organization_id, name) DO UPDATE SET department_id = EXCLUDED.department_id
		RETURNING id
	`, orgID, salesDept).Scan(&eastTeam)
	if err != nil {
		return err
	}

	_, err = db.Exec(ctx, `
		INSERT INTO branches (organization_id, name, location)
		VALUES ($1, 'HQ Bengaluru', 'Bengaluru, IN')
		ON CONFLICT (organization_id, name) DO NOTHING
	`, orgID)
	if err != nil {
		return err
	}

	_, err = db.Exec(ctx, `
		UPDATE organization_members om
		SET department_id = $2, team_id = $3, designation = CASE u.email
			WHEN 'priya@crmlite.com' THEN 'Sales Manager'
			WHEN 'vikram@crmlite.com' THEN 'Account Executive'
			WHEN 'sneha@crmlite.com' THEN 'Account Executive'
			ELSE om.designation END
		FROM users u
		WHERE om.user_id = u.id
		  AND om.organization_id = $1
		  AND u.email IN ('priya@crmlite.com','vikram@crmlite.com','sneha@crmlite.com')
	`, orgID, salesDept, eastTeam)
	return err
}
