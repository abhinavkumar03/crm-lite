package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/abhinavkumar03/crm-lite/backend/internal/rbac"
	"github.com/abhinavkumar03/crm-lite/backend/internal/roles/entity"
)

type Repository struct {
	db *pgxpool.Pool
}

func New(db *pgxpool.Pool) *Repository {
	return &Repository{db: db}
}

func (r *Repository) ListPermissions(ctx context.Context) ([]entity.Permission, error) {
	rows, err := r.db.Query(ctx, `
		SELECT id::text, key, category, description
		FROM permissions
		ORDER BY category, key
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	items := make([]entity.Permission, 0)
	for rows.Next() {
		var p entity.Permission
		if err := rows.Scan(&p.ID, &p.Key, &p.Category, &p.Description); err != nil {
			return nil, err
		}
		items = append(items, p)
	}
	return items, rows.Err()
}

func (r *Repository) ListRoles(ctx context.Context, orgID string) ([]entity.Role, []int, error) {
	rows, err := r.db.Query(ctx, `
		SELECT r.id::text, r.organization_id::text, r.name, r.slug, r.description,
		       r.is_system, r.created_at, r.updated_at,
		       (SELECT COUNT(*) FROM organization_members om
		         WHERE om.role_id = r.id AND om.status = 'active') AS member_count
		FROM roles r
		WHERE r.organization_id = $1
		ORDER BY r.is_system DESC, r.name ASC
	`, orgID)
	if err != nil {
		return nil, nil, err
	}
	defer rows.Close()

	roles := make([]entity.Role, 0)
	counts := make([]int, 0)
	for rows.Next() {
		var role entity.Role
		var count int
		if err := rows.Scan(
			&role.ID, &role.OrganizationID, &role.Name, &role.Slug, &role.Description,
			&role.IsSystem, &role.CreatedAt, &role.UpdatedAt, &count,
		); err != nil {
			return nil, nil, err
		}
		roles = append(roles, role)
		counts = append(counts, count)
	}
	return roles, counts, rows.Err()
}

func (r *Repository) GetByID(ctx context.Context, orgID, roleID string) (*entity.Role, error) {
	var role entity.Role
	err := r.db.QueryRow(ctx, `
		SELECT id::text, organization_id::text, name, slug, description,
		       is_system, created_at, updated_at
		FROM roles
		WHERE id = $1 AND organization_id = $2
	`, roleID, orgID).Scan(
		&role.ID, &role.OrganizationID, &role.Name, &role.Slug, &role.Description,
		&role.IsSystem, &role.CreatedAt, &role.UpdatedAt,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &role, nil
}

func (r *Repository) MemberCount(ctx context.Context, roleID string) (int, error) {
	var n int
	err := r.db.QueryRow(ctx, `
		SELECT COUNT(*) FROM organization_members
		WHERE role_id = $1 AND status = 'active'
	`, roleID).Scan(&n)
	return n, err
}

func (r *Repository) Create(ctx context.Context, role *entity.Role) error {
	return r.db.QueryRow(ctx, `
		INSERT INTO roles (organization_id, name, slug, description, is_system)
		VALUES ($1, $2, $3, $4, FALSE)
		RETURNING id::text, created_at, updated_at
	`, role.OrganizationID, role.Name, role.Slug, role.Description).
		Scan(&role.ID, &role.CreatedAt, &role.UpdatedAt)
}

func (r *Repository) Update(ctx context.Context, role *entity.Role) error {
	return r.db.QueryRow(ctx, `
		UPDATE roles
		SET name = $3, description = $4, updated_at = NOW()
		WHERE id = $1 AND organization_id = $2
		RETURNING updated_at
	`, role.ID, role.OrganizationID, role.Name, role.Description).Scan(&role.UpdatedAt)
}

func (r *Repository) Delete(ctx context.Context, orgID, roleID string) error {
	tag, err := r.db.Exec(ctx, `
		DELETE FROM roles
		WHERE id = $1 AND organization_id = $2 AND is_system = FALSE
	`, roleID, orgID)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return pgx.ErrNoRows
	}
	return nil
}

func (r *Repository) PermissionKeys(ctx context.Context, roleID string) ([]string, error) {
	rows, err := r.db.Query(ctx, `
		SELECT p.key
		FROM role_permissions rp
		JOIN permissions p ON p.id = rp.permission_id
		WHERE rp.role_id = $1
		ORDER BY p.key
	`, roleID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	keys := make([]string, 0)
	for rows.Next() {
		var k string
		if err := rows.Scan(&k); err != nil {
			return nil, err
		}
		keys = append(keys, k)
	}
	return keys, rows.Err()
}

// SetPermissions replaces the role's grants. Unknown keys are ignored so a
// stale client matrix cannot break the write.
func (r *Repository) SetPermissions(ctx context.Context, roleID string, keys []string) error {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	if _, err := tx.Exec(ctx, `DELETE FROM role_permissions WHERE role_id = $1`, roleID); err != nil {
		return err
	}
	if len(keys) > 0 {
		_, err = tx.Exec(ctx, `
			INSERT INTO role_permissions (role_id, permission_id)
			SELECT $1, p.id FROM permissions p WHERE p.key = ANY($2)
			ON CONFLICT DO NOTHING
		`, roleID, keys)
		if err != nil {
			return err
		}
	}
	return tx.Commit(ctx)
}

func (r *Repository) SetModuleAccess(ctx context.Context, roleID string, access []rbac.ModuleAccess) error {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	if _, err := tx.Exec(ctx, `DELETE FROM role_module_access WHERE role_id = $1`, roleID); err != nil {
		return err
	}
	for _, a := range access {
		_, err = tx.Exec(ctx, `
			INSERT INTO role_module_access (role_id, module_id, can_view, can_create, can_update, can_delete)
			SELECT $1, m.id, $3, $4, $5, $6
			FROM modules m
			WHERE m.id = $2
			ON CONFLICT (role_id, module_id) DO UPDATE SET
				can_view = EXCLUDED.can_view,
				can_create = EXCLUDED.can_create,
				can_update = EXCLUDED.can_update,
				can_delete = EXCLUDED.can_delete,
				updated_at = NOW()
		`, roleID, a.ModuleID, a.CanView, a.CanCreate, a.CanUpdate, a.CanDelete)
		if err != nil {
			return fmt.Errorf("module access %s: %w", a.ModuleID, err)
		}
	}
	return tx.Commit(ctx)
}

func (r *Repository) SetFieldAccess(ctx context.Context, roleID string, access []rbac.FieldAccess) error {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	if _, err := tx.Exec(ctx, `DELETE FROM role_field_access WHERE role_id = $1`, roleID); err != nil {
		return err
	}
	for _, a := range access {
		if a.Access != rbac.FieldHidden && a.Access != rbac.FieldRead && a.Access != rbac.FieldWrite {
			return fmt.Errorf("invalid field access %q", a.Access)
		}
		_, err = tx.Exec(ctx, `
			INSERT INTO role_field_access (role_id, field_id, access)
			SELECT $1, f.id, $3
			FROM fields f
			WHERE f.id = $2
			ON CONFLICT (role_id, field_id) DO UPDATE SET
				access = EXCLUDED.access,
				updated_at = NOW()
		`, roleID, a.FieldID, a.Access)
		if err != nil {
			return fmt.Errorf("field access %s: %w", a.FieldID, err)
		}
	}
	return tx.Commit(ctx)
}

func (r *Repository) SlugExists(ctx context.Context, orgID, slug string) (bool, error) {
	var exists bool
	err := r.db.QueryRow(ctx, `
		SELECT EXISTS(
			SELECT 1 FROM roles WHERE organization_id = $1 AND slug = $2
		)
	`, orgID, slug).Scan(&exists)
	return exists, err
}
