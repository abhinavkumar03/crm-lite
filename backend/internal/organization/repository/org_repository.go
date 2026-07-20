package repository

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"golang.org/x/crypto/bcrypt"

	"github.com/abhinavkumar03/crm-lite/backend/internal/organization/dto"
)

type Repository struct {
	db *pgxpool.Pool
}

func New(db *pgxpool.Pool) *Repository {
	return &Repository{db: db}
}

func (r *Repository) ListOrgsForUser(ctx context.Context, userID string) ([]dto.OrgSummary, error) {
	rows, err := r.db.Query(ctx, `
		SELECT o.id, o.name, o.slug, o.logo_url, o.description, COALESCE(rl.slug, ''),
		       (u.active_organization_id IS NOT NULL AND u.active_organization_id = o.id)
		FROM organization_members om
		JOIN organizations o ON o.id = om.organization_id
		JOIN users u ON u.id = om.user_id
		LEFT JOIN roles rl ON rl.id = om.role_id
		WHERE om.user_id = $1 AND om.status = 'active'
		  AND o.deleted_at IS NULL
		ORDER BY o.name ASC
	`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	out := make([]dto.OrgSummary, 0)
	for rows.Next() {
		var s dto.OrgSummary
		if err := rows.Scan(&s.ID, &s.Name, &s.Slug, &s.LogoURL, &s.Description, &s.RoleSlug, &s.IsActive); err != nil {
			return nil, err
		}
		out = append(out, s)
	}
	return out, rows.Err()
}

// OrgRow is the raw organization profile used by Get/Update current workspace.
type OrgRow struct {
	ID          string
	Name        string
	Slug        string
	Plan        string
	LogoURL     *string
	Description *string
	Industry    *string
	CompanySize *string
	Country     *string
	Status      string
	CreatedBy   *string
	Settings    []byte
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

func (r *Repository) GetOrgByID(ctx context.Context, orgID string) (*OrgRow, error) {
	var o OrgRow
	err := r.db.QueryRow(ctx, `
		SELECT id, name, slug, plan, logo_url, description, industry, company_size,
		       country, status, created_by::text, settings, created_at, updated_at
		FROM organizations
		WHERE id = $1 AND deleted_at IS NULL
	`, orgID).Scan(
		&o.ID, &o.Name, &o.Slug, &o.Plan, &o.LogoURL, &o.Description,
		&o.Industry, &o.CompanySize, &o.Country, &o.Status, &o.CreatedBy,
		&o.Settings, &o.CreatedAt, &o.UpdatedAt,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	o.CreatedBy = nullEmpty(o.CreatedBy)
	return &o, nil
}

type OrgProfileUpdate struct {
	Name        string
	LogoURL     *string
	Description *string
	Industry    *string
	CompanySize *string
	Country     *string
	Settings    []byte
}

func (r *Repository) UpdateOrg(ctx context.Context, orgID string, p OrgProfileUpdate) (*OrgRow, error) {
	var o OrgRow
	err := r.db.QueryRow(ctx, `
		UPDATE organizations
		SET name = $2,
		    logo_url = $3,
		    description = $4,
		    industry = $5,
		    company_size = $6,
		    country = $7,
		    settings = $8,
		    updated_at = NOW()
		WHERE id = $1 AND deleted_at IS NULL
		RETURNING id, name, slug, plan, logo_url, description, industry, company_size,
		          country, status, created_by::text, settings, created_at, updated_at
	`, orgID, p.Name, p.LogoURL, p.Description, p.Industry, p.CompanySize, p.Country, p.Settings).Scan(
		&o.ID, &o.Name, &o.Slug, &o.Plan, &o.LogoURL, &o.Description,
		&o.Industry, &o.CompanySize, &o.Country, &o.Status, &o.CreatedBy,
		&o.Settings, &o.CreatedAt, &o.UpdatedAt,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	o.CreatedBy = nullEmpty(o.CreatedBy)
	return &o, nil
}

// SoftDeleteOrg marks the org deleted and clears active_organization_id for members
// who had it selected. Returns member user IDs that need tenant cache invalidation.
func (r *Repository) SoftDeleteOrg(ctx context.Context, orgID string) ([]string, error) {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx)

	tag, err := tx.Exec(ctx, `
		UPDATE organizations
		SET deleted_at = NOW(), status = 'inactive', updated_at = NOW()
		WHERE id = $1 AND deleted_at IS NULL
	`, orgID)
	if err != nil {
		return nil, err
	}
	if tag.RowsAffected() == 0 {
		return nil, ErrNotFound
	}

	rows, err := tx.Query(ctx, `
		SELECT user_id::text FROM organization_members
		WHERE organization_id = $1 AND status = 'active'
	`, orgID)
	if err != nil {
		return nil, err
	}
	memberIDs := make([]string, 0)
	for rows.Next() {
		var id string
		if err := rows.Scan(&id); err != nil {
			rows.Close()
			return nil, err
		}
		memberIDs = append(memberIDs, id)
	}
	rows.Close()
	if err := rows.Err(); err != nil {
		return nil, err
	}

	_, err = tx.Exec(ctx, `
		UPDATE users
		SET active_organization_id = NULL, updated_at = NOW()
		WHERE active_organization_id = $1
	`, orgID)
	if err != nil {
		return nil, err
	}

	// Point survivors at another active (non-deleted) membership when available.
	for _, uid := range memberIDs {
		_, err = tx.Exec(ctx, `
			UPDATE users u
			SET active_organization_id = sub.org_id, updated_at = NOW()
			FROM (
				SELECT om.organization_id AS org_id
				FROM organization_members om
				JOIN organizations o ON o.id = om.organization_id
				WHERE om.user_id = $1
				  AND om.status = 'active'
				  AND o.deleted_at IS NULL
				ORDER BY om.created_at ASC
				LIMIT 1
			) sub
			WHERE u.id = $1 AND u.active_organization_id IS NULL
		`, uid)
		if err != nil {
			return nil, err
		}
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, err
	}
	return memberIDs, nil
}

// CountActiveMemberships returns how many non-deleted orgs the user still belongs to.
func (r *Repository) CountActiveMemberships(ctx context.Context, userID string) (int, error) {
	var n int
	err := r.db.QueryRow(ctx, `
		SELECT count(*)
		FROM organization_members om
		JOIN organizations o ON o.id = om.organization_id
		WHERE om.user_id = $1 AND om.status = 'active' AND o.deleted_at IS NULL
	`, userID).Scan(&n)
	return n, err
}

func (r *Repository) ListMembers(ctx context.Context, orgID string) ([]dto.MemberResponse, error) {
	rows, err := r.db.Query(ctx, `
		SELECT u.id, u.name, u.email,
		       om.role_id::text, COALESCE(rl.slug, ''),
		       om.manager_user_id::text, om.department_id::text, om.team_id::text,
		       om.branch_id::text, om.designation, om.hierarchy_level, om.status
		FROM organization_members om
		JOIN users u ON u.id = om.user_id
		LEFT JOIN roles rl ON rl.id = om.role_id
		WHERE om.organization_id = $1
		ORDER BY u.name ASC
	`, orgID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	out := make([]dto.MemberResponse, 0)
	for rows.Next() {
		var m dto.MemberResponse
		var roleID, managerID, deptID, teamID, branchID *string
		var designation *string
		if err := rows.Scan(
			&m.UserID, &m.Name, &m.Email,
			&roleID, &m.RoleSlug,
			&managerID, &deptID, &teamID, &branchID, &designation,
			&m.HierarchyLevel, &m.Status,
		); err != nil {
			return nil, err
		}
		m.RoleID = nullEmpty(roleID)
		m.ManagerUserID = nullEmpty(managerID)
		m.DepartmentID = nullEmpty(deptID)
		m.TeamID = nullEmpty(teamID)
		m.BranchID = nullEmpty(branchID)
		m.Designation = designation
		out = append(out, m)
	}
	return out, rows.Err()
}

func nullEmpty(p *string) *string {
	if p == nil || *p == "" {
		return nil
	}
	return p
}

func (r *Repository) CreateInvitation(
	ctx context.Context,
	orgID, email, roleID, invitedBy string,
	managerID, deptID, teamID *string,
) (*dto.InviteResponse, error) {
	token, err := randomToken(32)
	if err != nil {
		return nil, err
	}
	expires := time.Now().Add(7 * 24 * time.Hour)
	var id string
	err = r.db.QueryRow(ctx, `
		INSERT INTO organization_invitations (
			organization_id, email, role_id, manager_user_id, department_id, team_id,
			token, status, invited_by, expires_at
		) VALUES ($1,$2,$3,$4,$5,$6,$7,'pending',$8,$9)
		RETURNING id
	`, orgID, email, roleID, managerID, deptID, teamID, token, invitedBy, expires).Scan(&id)
	if err != nil {
		return nil, err
	}
	body := fmt.Sprintf(
		"[CRM Lite simulation] Invite for %s. Accept with token %s (expires %s).",
		email, token, expires.Format(time.RFC3339),
	)
	return &dto.InviteResponse{
		ID:             id,
		Email:          email,
		Token:          token,
		Status:         "pending",
		ExpiresAt:      expires,
		SimulatedEmail: body,
	}, nil
}

type PendingInvite struct {
	ID             string
	OrganizationID string
	Email          string
	RoleID         *string
	ManagerUserID  *string
	DepartmentID   *string
	TeamID         *string
	ExpiresAt      time.Time
}

func (r *Repository) GetPendingInvite(ctx context.Context, token string) (*PendingInvite, error) {
	var inv PendingInvite
	err := r.db.QueryRow(ctx, `
		SELECT id, organization_id, email, role_id::text, manager_user_id::text,
		       department_id::text, team_id::text, expires_at
		FROM organization_invitations
		WHERE token = $1 AND status = 'pending'
	`, token).Scan(
		&inv.ID, &inv.OrganizationID, &inv.Email, &inv.RoleID,
		&inv.ManagerUserID, &inv.DepartmentID, &inv.TeamID, &inv.ExpiresAt,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	inv.RoleID = nullEmpty(inv.RoleID)
	inv.ManagerUserID = nullEmpty(inv.ManagerUserID)
	inv.DepartmentID = nullEmpty(inv.DepartmentID)
	inv.TeamID = nullEmpty(inv.TeamID)
	return &inv, nil
}

func (r *Repository) AcceptInvite(ctx context.Context, inv *PendingInvite, name, password string) (userID string, err error) {
	if time.Now().After(inv.ExpiresAt) {
		_, _ = r.db.Exec(ctx, `UPDATE organization_invitations SET status='expired', updated_at=NOW() WHERE id=$1`, inv.ID)
		return "", ErrInviteExpired
	}

	err = r.db.QueryRow(ctx, `SELECT id FROM users WHERE email = $1`, inv.Email).Scan(&userID)
	if errors.Is(err, pgx.ErrNoRows) {
		if name == "" || password == "" {
			return "", ErrPasswordRequired
		}
		hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		if err != nil {
			return "", err
		}
		err = r.db.QueryRow(ctx, `
			INSERT INTO users (name, email, password_hash, active_organization_id)
			VALUES ($1,$2,$3,$4) RETURNING id
		`, name, inv.Email, string(hash), inv.OrganizationID).Scan(&userID)
		if err != nil {
			return "", err
		}
	} else if err != nil {
		return "", err
	}

	_, err = r.db.Exec(ctx, `
		INSERT INTO organization_members (
			organization_id, user_id, role_id, status,
			manager_user_id, department_id, team_id, hierarchy_level
		) VALUES ($1,$2,$3,'active',$4,$5,$6,60)
		ON CONFLICT (organization_id, user_id) DO UPDATE
		SET role_id = EXCLUDED.role_id,
		    status = 'active',
		    manager_user_id = COALESCE(EXCLUDED.manager_user_id, organization_members.manager_user_id),
		    department_id = COALESCE(EXCLUDED.department_id, organization_members.department_id),
		    team_id = COALESCE(EXCLUDED.team_id, organization_members.team_id)
	`, inv.OrganizationID, userID, inv.RoleID, inv.ManagerUserID, inv.DepartmentID, inv.TeamID)
	if err != nil {
		return "", err
	}

	_, err = r.db.Exec(ctx, `
		UPDATE users SET active_organization_id = COALESCE(active_organization_id, $2), updated_at = NOW()
		WHERE id = $1
	`, userID, inv.OrganizationID)
	if err != nil {
		return "", err
	}

	_, err = r.db.Exec(ctx, `
		UPDATE organization_invitations
		SET status = 'accepted', accepted_at = NOW(), updated_at = NOW()
		WHERE id = $1
	`, inv.ID)
	return userID, err
}

var (
	ErrInviteExpired    = errors.New("invitation expired")
	ErrPasswordRequired = errors.New("name and password required for new users")
	ErrNotFound         = errors.New("not found")
)

func randomToken(n int) (string, error) {
	b := make([]byte, n)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}

// --- structure CRUD -------------------------------------------------------

func (r *Repository) ListDepartments(ctx context.Context, orgID string) ([]dto.StructureItem, error) {
	rows, err := r.db.Query(ctx, `
		SELECT id, name, description FROM departments WHERE organization_id = $1 ORDER BY name
	`, orgID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	out := make([]dto.StructureItem, 0)
	for rows.Next() {
		var item dto.StructureItem
		if err := rows.Scan(&item.ID, &item.Name, &item.Description); err != nil {
			return nil, err
		}
		out = append(out, item)
	}
	return out, rows.Err()
}

func (r *Repository) CreateDepartment(ctx context.Context, orgID string, req dto.CreateDepartmentRequest) (*dto.StructureItem, error) {
	var item dto.StructureItem
	err := r.db.QueryRow(ctx, `
		INSERT INTO departments (organization_id, name, description)
		VALUES ($1,$2,$3) RETURNING id, name, description
	`, orgID, req.Name, req.Description).Scan(&item.ID, &item.Name, &item.Description)
	if err != nil {
		return nil, err
	}
	return &item, nil
}

func (r *Repository) ListTeams(ctx context.Context, orgID string) ([]dto.StructureItem, error) {
	rows, err := r.db.Query(ctx, `
		SELECT id, name, description, department_id::text FROM teams
		WHERE organization_id = $1 ORDER BY name
	`, orgID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	out := make([]dto.StructureItem, 0)
	for rows.Next() {
		var item dto.StructureItem
		var dept *string
		if err := rows.Scan(&item.ID, &item.Name, &item.Description, &dept); err != nil {
			return nil, err
		}
		item.DepartmentID = nullEmpty(dept)
		out = append(out, item)
	}
	return out, rows.Err()
}

func (r *Repository) CreateTeam(ctx context.Context, orgID string, req dto.CreateTeamRequest) (*dto.StructureItem, error) {
	var item dto.StructureItem
	var dept *string
	err := r.db.QueryRow(ctx, `
		INSERT INTO teams (organization_id, name, description, department_id)
		VALUES ($1,$2,$3,$4) RETURNING id, name, description, department_id::text
	`, orgID, req.Name, req.Description, req.DepartmentID).Scan(&item.ID, &item.Name, &item.Description, &dept)
	if err != nil {
		return nil, err
	}
	item.DepartmentID = nullEmpty(dept)
	return &item, nil
}

func (r *Repository) ListBranches(ctx context.Context, orgID string) ([]dto.StructureItem, error) {
	rows, err := r.db.Query(ctx, `
		SELECT id, name, location FROM branches WHERE organization_id = $1 ORDER BY name
	`, orgID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	out := make([]dto.StructureItem, 0)
	for rows.Next() {
		var item dto.StructureItem
		if err := rows.Scan(&item.ID, &item.Name, &item.Location); err != nil {
			return nil, err
		}
		out = append(out, item)
	}
	return out, rows.Err()
}

func (r *Repository) CreateBranch(ctx context.Context, orgID string, req dto.CreateBranchRequest) (*dto.StructureItem, error) {
	var item dto.StructureItem
	err := r.db.QueryRow(ctx, `
		INSERT INTO branches (organization_id, name, location)
		VALUES ($1,$2,$3) RETURNING id, name, location
	`, orgID, req.Name, req.Location).Scan(&item.ID, &item.Name, &item.Location)
	if err != nil {
		return nil, err
	}
	return &item, nil
}

// SubordinateUserIDs returns user ids in the reporting tree under managerUserID (inclusive of self).
func (r *Repository) SubordinateUserIDs(ctx context.Context, orgID, managerUserID string) ([]string, error) {
	rows, err := r.db.Query(ctx, `
		WITH RECURSIVE tree AS (
			SELECT user_id FROM organization_members
			WHERE organization_id = $1 AND user_id = $2 AND status = 'active'
			UNION ALL
			SELECT om.user_id
			FROM organization_members om
			JOIN tree t ON om.manager_user_id = t.user_id
			WHERE om.organization_id = $1 AND om.status = 'active'
		)
		SELECT user_id FROM tree
	`, orgID, managerUserID)
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

func (r *Repository) MemberMeta(ctx context.Context, orgID, userID string) (deptID, teamID *string, roleSlug string, hierarchyLevel int, err error) {
	var d, t *string
	err = r.db.QueryRow(ctx, `
		SELECT om.department_id::text, om.team_id::text, COALESCE(rl.slug,''), om.hierarchy_level
		FROM organization_members om
		LEFT JOIN roles rl ON rl.id = om.role_id
		WHERE om.organization_id = $1 AND om.user_id = $2 AND om.status = 'active'
	`, orgID, userID).Scan(&d, &t, &roleSlug, &hierarchyLevel)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil, "", 100, nil
	}
	return nullEmpty(d), nullEmpty(t), roleSlug, hierarchyLevel, err
}
