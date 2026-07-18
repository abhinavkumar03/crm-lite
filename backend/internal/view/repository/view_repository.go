package repository

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/abhinavkumar03/crm-lite/backend/internal/view/entity"
)

const viewColumns = `
	id, organization_id, module_id, name, columns, filters, sort,
	is_default, is_public, owner_id, created_at, updated_at
`

type Repository struct {
	db *pgxpool.Pool
}

func New(db *pgxpool.Pool) *Repository {
	return &Repository{db: db}
}

func scanView(row pgx.Row, v *entity.View) error {
	return row.Scan(
		&v.ID, &v.OrganizationID, &v.ModuleID, &v.Name, &v.Columns, &v.Filters, &v.Sort,
		&v.IsDefault, &v.IsPublic, &v.OwnerID, &v.CreatedAt, &v.UpdatedAt,
	)
}

func (r *Repository) ModuleExists(ctx context.Context, orgID, moduleID string) (bool, error) {
	var exists bool
	err := r.db.QueryRow(ctx, `
		SELECT EXISTS(SELECT 1 FROM modules WHERE id = $1 AND organization_id = $2)
	`, moduleID, orgID).Scan(&exists)
	return exists, err
}

func (r *Repository) Create(ctx context.Context, v *entity.View) error {
	return r.db.QueryRow(ctx, `
		INSERT INTO views (
			organization_id, module_id, name, columns, filters, sort,
			is_default, is_public, owner_id
		)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9)
		RETURNING id, created_at, updated_at
	`,
		v.OrganizationID, v.ModuleID, v.Name, v.Columns, v.Filters, v.Sort,
		v.IsDefault, v.IsPublic, v.OwnerID,
	).Scan(&v.ID, &v.CreatedAt, &v.UpdatedAt)
}

// ListVisible returns views on the module that are public or owned by the user,
// default views first.
func (r *Repository) ListVisible(ctx context.Context, orgID, moduleID, userID string) ([]entity.View, error) {
	rows, err := r.db.Query(ctx, `
		SELECT `+viewColumns+`
		FROM views
		WHERE organization_id = $1 AND module_id = $2
		  AND (is_public = TRUE OR owner_id = $3)
		ORDER BY is_default DESC, name ASC
	`, orgID, moduleID, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	views := make([]entity.View, 0)
	for rows.Next() {
		var v entity.View
		if err := scanView(rows, &v); err != nil {
			return nil, err
		}
		views = append(views, v)
	}
	return views, rows.Err()
}

func (r *Repository) GetByID(ctx context.Context, orgID, moduleID, id string) (*entity.View, error) {
	var v entity.View
	err := scanView(r.db.QueryRow(ctx, `
		SELECT `+viewColumns+`
		FROM views
		WHERE id = $1 AND module_id = $2 AND organization_id = $3
	`, id, moduleID, orgID), &v)

	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &v, nil
}

func (r *Repository) Update(ctx context.Context, v *entity.View) error {
	return r.db.QueryRow(ctx, `
		UPDATE views SET
			name = $1,
			columns = $2,
			filters = $3,
			sort = $4,
			is_public = $5,
			updated_at = NOW()
		WHERE id = $6 AND module_id = $7 AND organization_id = $8
		RETURNING updated_at
	`,
		v.Name, v.Columns, v.Filters, v.Sort, v.IsPublic,
		v.ID, v.ModuleID, v.OrganizationID,
	).Scan(&v.UpdatedAt)
}

func (r *Repository) Delete(ctx context.Context, orgID, moduleID, id string) (bool, error) {
	tag, err := r.db.Exec(ctx, `
		DELETE FROM views WHERE id = $1 AND module_id = $2 AND organization_id = $3
	`, id, moduleID, orgID)
	if err != nil {
		return false, err
	}
	return tag.RowsAffected() > 0, nil
}

// SetDefault makes one view the module default and clears the flag on the rest,
// atomically.
func (r *Repository) SetDefault(ctx context.Context, orgID, moduleID, id string) (bool, error) {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return false, err
	}
	defer func() { _ = tx.Rollback(ctx) }()

	if _, err := tx.Exec(ctx, `
		UPDATE views SET is_default = FALSE, updated_at = NOW()
		WHERE organization_id = $1 AND module_id = $2 AND is_default = TRUE
	`, orgID, moduleID); err != nil {
		return false, err
	}

	tag, err := tx.Exec(ctx, `
		UPDATE views SET is_default = TRUE, updated_at = NOW()
		WHERE id = $1 AND module_id = $2 AND organization_id = $3
	`, id, moduleID, orgID)
	if err != nil {
		return false, err
	}

	if tag.RowsAffected() == 0 {
		return false, nil
	}

	return true, tx.Commit(ctx)
}
