package seeders

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
	"golang.org/x/crypto/bcrypt"

	"github.com/abhinavkumar03/crm-lite/backend/internal/organization/bootstrap"
)

func getOrgID(ctx context.Context, db *pgxpool.Pool) (string, error) {
	return getOrgIDBySlug(ctx, db, primaryOrgSlug)
}

func getOrgIDBySlug(ctx context.Context, db *pgxpool.Pool, slug string) (string, error) {
	var id string
	err := db.QueryRow(ctx,
		`SELECT id FROM organizations WHERE slug = $1 AND deleted_at IS NULL`, slug,
	).Scan(&id)
	if err != nil {
		return "", fmt.Errorf("lookup org %q: %w", slug, err)
	}
	return id, nil
}

func listDemoOrgIDs(ctx context.Context, db *pgxpool.Pool) ([]string, error) {
	ids := make([]string, 0, len(demoWorkspaces))
	for _, ws := range demoWorkspaces {
		id, err := getOrgIDBySlug(ctx, db, ws.Slug)
		if err != nil {
			return nil, err
		}
		ids = append(ids, id)
	}
	return ids, nil
}

type OrganizationSeeder struct{}

func (OrganizationSeeder) Name() string { return "organization" }

func (OrganizationSeeder) Run(ctx context.Context, db *pgxpool.Pool) error {
	var demoUserID string
	err := db.QueryRow(ctx, `SELECT id FROM users WHERE email = $1`, demoUserEmail).Scan(&demoUserID)
	if err != nil {
		return fmt.Errorf("demo user must exist before organizations: %w", err)
	}

	boot := bootstrap.New(db)
	for _, ws := range demoWorkspaces {
		var existing string
		err := db.QueryRow(ctx,
			`SELECT id FROM organizations WHERE slug = $1`, ws.Slug,
		).Scan(&existing)
		if err == nil && existing != "" {
			_, _ = db.Exec(ctx, `
				UPDATE organizations
				SET name = $2, description = $3, industry = $4, company_size = $5,
				    country = $6, logo_url = $7, plan = $8, status = 'active',
				    deleted_at = NULL, updated_at = NOW()
				WHERE id = $1
			`, existing, ws.Name, ws.Description, ws.Industry, ws.CompanySize,
				ws.Country, ws.LogoURL, ws.Plan)
			if err := boot.EnsureFullCatalog(ctx, existing); err != nil {
				return fmt.Errorf("ensure catalog %q: %w", ws.Slug, err)
			}
			continue
		}

		_, err = boot.CreateOrganization(ctx, bootstrap.CreateOptions{
			Name:        ws.Name,
			Slug:        ws.Slug,
			Description: ws.Description,
			Industry:    ws.Industry,
			CompanySize: ws.CompanySize,
			Country:     ws.Country,
			LogoURL:     ws.LogoURL,
			Timezone:    ws.Timezone,
			Currency:    ws.Currency,
			Locale:      ws.Locale,
			DateFormat:  ws.DateFormat,
		}, demoUserID)
		if err != nil {
			return fmt.Errorf("create workspace %q: %w", ws.Slug, err)
		}

		_, err = db.Exec(ctx, `
			UPDATE organizations SET plan = $2 WHERE slug = $1
		`, ws.Slug, ws.Plan)
		if err != nil {
			return err
		}
	}
	return nil
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
	orgIDs, err := listDemoOrgIDs(ctx, db)
	if err != nil {
		return err
	}
	for _, orgID := range orgIDs {
		if err := ensureRolesForOrg(ctx, db, orgID); err != nil {
			return err
		}
	}
	return nil
}

func ensureRolesForOrg(ctx context.Context, db *pgxpool.Pool, orgID string) error {
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
	{"Demo Owner", demoUserEmail, demoUserPass, "owner"},
	{"Organization Owner", "admin@crm.com", "Admin@123", "owner"},
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
			ON CONFLICT (email) DO UPDATE
			SET name = EXCLUDED.name,
			    password_hash = EXCLUDED.password_hash
		`, u.Name, u.Email, string(hash))
		if err != nil {
			return fmt.Errorf("upsert user %q: %w", u.Email, err)
		}
	}
	return nil
}

type MembershipsSeeder struct{}

func (MembershipsSeeder) Name() string { return "memberships" }

func (MembershipsSeeder) Run(ctx context.Context, db *pgxpool.Pool) error {
	orgIDs, err := listDemoOrgIDs(ctx, db)
	if err != nil {
		return err
	}
	primaryID := orgIDs[0]

	// Demo user is already owner on every workspace via bootstrap; ensure role.
	for _, orgID := range orgIDs {
		_, err := db.Exec(ctx, `
			INSERT INTO organization_members (organization_id, user_id, role_id, hierarchy_level, status)
			SELECT $1, usr.id, rol.id, 0, 'active'
			FROM users usr
			JOIN roles rol ON rol.organization_id = $1 AND rol.slug = 'owner'
			WHERE usr.email = $2
			ON CONFLICT (organization_id, user_id) DO UPDATE
			SET role_id = EXCLUDED.role_id, hierarchy_level = 0, status = 'active'
		`, orgID, demoUserEmail)
		if err != nil {
			return fmt.Errorf("demo owner membership: %w", err)
		}
	}

	// Team users belong to the primary workspace for RBAC demos.
	for _, u := range userDefs {
		if u.Email == demoUserEmail {
			continue
		}
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
			INSERT INTO organization_members (organization_id, user_id, role_id, hierarchy_level, status)
			SELECT $1, usr.id, rol.id, $4, 'active'
			FROM users usr
			JOIN roles rol ON rol.organization_id = $1 AND rol.slug = $3
			WHERE usr.email = $2
			ON CONFLICT (organization_id, user_id) DO UPDATE
			SET role_id = EXCLUDED.role_id, hierarchy_level = EXCLUDED.hierarchy_level, status = 'active'
		`, primaryID, u.Email, u.Role, level)
		if err != nil {
			return fmt.Errorf("membership for %q: %w", u.Email, err)
		}
	}

	_, err = db.Exec(ctx, `
		UPDATE organization_members om
		SET manager_user_id = mgr.id
		FROM users mgr, users rep
		WHERE om.organization_id = $1
		  AND om.user_id = rep.id
		  AND mgr.email = 'priya@crmlite.com'
		  AND rep.email IN ('vikram@crmlite.com', 'sneha@crmlite.com')
	`, primaryID)
	if err != nil {
		return err
	}

	emails := make([]string, 0, len(userDefs))
	for _, u := range userDefs {
		emails = append(emails, u.Email)
	}
	_, err = db.Exec(ctx, `
		UPDATE users SET active_organization_id = $1
		WHERE email = ANY($2)
	`, primaryID, emails)
	return err
}

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
