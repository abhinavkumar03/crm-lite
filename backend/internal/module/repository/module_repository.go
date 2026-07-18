package repository

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/abhinavkumar03/crm-lite/backend/internal/module/entity"
)

const moduleColumns = `
	id, organization_id, api_name, singular_label, plural_label,
	description, icon, color, storage_strategy, native_table,
	is_system, is_enabled, is_visible_sidebar, sort_order,
	default_sort_field, default_sort_order, created_at, updated_at
`

type Repository struct {
	db *pgxpool.Pool
}

func New(db *pgxpool.Pool) *Repository {
	return &Repository{db: db}
}

func scanModule(row pgx.Row, m *entity.Module) error {
	return row.Scan(
		&m.ID, &m.OrganizationID, &m.APIName, &m.SingularLabel, &m.PluralLabel,
		&m.Description, &m.Icon, &m.Color, &m.StorageStrategy, &m.NativeTable,
		&m.IsSystem, &m.IsEnabled, &m.IsVisibleSidebar, &m.SortOrder,
		&m.DefaultSortField, &m.DefaultSortOrder, &m.CreatedAt, &m.UpdatedAt,
	)
}

func (r *Repository) Create(ctx context.Context, m *entity.Module) error {
	query := `
		INSERT INTO modules (
			organization_id, api_name, singular_label, plural_label,
			description, icon, color, storage_strategy, native_table,
			is_system, is_enabled, is_visible_sidebar, sort_order,
			default_sort_field, default_sort_order
		)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,$15)
		RETURNING id, created_at, updated_at
	`
	return r.db.QueryRow(ctx, query,
		m.OrganizationID, m.APIName, m.SingularLabel, m.PluralLabel,
		m.Description, m.Icon, m.Color, m.StorageStrategy, m.NativeTable,
		m.IsSystem, m.IsEnabled, m.IsVisibleSidebar, m.SortOrder,
		m.DefaultSortField, m.DefaultSortOrder,
	).Scan(&m.ID, &m.CreatedAt, &m.UpdatedAt)
}

func (r *Repository) List(ctx context.Context, orgID string) ([]entity.Module, error) {
	rows, err := r.db.Query(ctx, `
		SELECT `+moduleColumns+`
		FROM modules
		WHERE organization_id = $1
		ORDER BY sort_order ASC, singular_label ASC
	`, orgID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return collect(rows)
}

func (r *Repository) Navigation(ctx context.Context, orgID string) ([]entity.Module, error) {
	rows, err := r.db.Query(ctx, `
		SELECT `+moduleColumns+`
		FROM modules
		WHERE organization_id = $1
		  AND is_enabled = TRUE
		  AND is_visible_sidebar = TRUE
		ORDER BY sort_order ASC, singular_label ASC
	`, orgID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return collect(rows)
}

func (r *Repository) GetByID(ctx context.Context, orgID, id string) (*entity.Module, error) {
	var m entity.Module
	err := scanModule(r.db.QueryRow(ctx, `
		SELECT `+moduleColumns+`
		FROM modules
		WHERE id = $1 AND organization_id = $2
	`, id, orgID), &m)

	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &m, nil
}

func (r *Repository) Update(ctx context.Context, m *entity.Module) error {
	return r.db.QueryRow(ctx, `
		UPDATE modules SET
			singular_label = $1,
			plural_label = $2,
			description = $3,
			icon = $4,
			color = $5,
			is_visible_sidebar = $6,
			default_sort_field = $7,
			default_sort_order = $8,
			updated_at = NOW()
		WHERE id = $9 AND organization_id = $10
		RETURNING updated_at
	`,
		m.SingularLabel, m.PluralLabel, m.Description, m.Icon, m.Color,
		m.IsVisibleSidebar, m.DefaultSortField, m.DefaultSortOrder,
		m.ID, m.OrganizationID,
	).Scan(&m.UpdatedAt)
}

func (r *Repository) SetEnabled(ctx context.Context, orgID, id string, enabled bool) (bool, error) {
	tag, err := r.db.Exec(ctx, `
		UPDATE modules SET is_enabled = $1, updated_at = NOW()
		WHERE id = $2 AND organization_id = $3
	`, enabled, id, orgID)
	if err != nil {
		return false, err
	}
	return tag.RowsAffected() > 0, nil
}

func (r *Repository) Delete(ctx context.Context, orgID, id string) (bool, error) {
	tag, err := r.db.Exec(ctx, `
		DELETE FROM modules WHERE id = $1 AND organization_id = $2
	`, id, orgID)
	if err != nil {
		return false, err
	}
	return tag.RowsAffected() > 0, nil
}

func (r *Repository) ExistsByAPIName(ctx context.Context, orgID, apiName string) (bool, error) {
	var exists bool
	err := r.db.QueryRow(ctx, `
		SELECT EXISTS(
			SELECT 1 FROM modules
			WHERE organization_id = $1 AND api_name = $2
		)
	`, orgID, apiName).Scan(&exists)
	return exists, err
}

func (r *Repository) MaxSortOrder(ctx context.Context, orgID string) (int, error) {
	var max int
	err := r.db.QueryRow(ctx, `
		SELECT COALESCE(MAX(sort_order), 0) FROM modules WHERE organization_id = $1
	`, orgID).Scan(&max)
	return max, err
}

// Reorder updates sort_order for the given modules atomically.
func (r *Repository) Reorder(ctx context.Context, orgID string, positions []entity.SortPosition) error {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return err
	}
	defer func() { _ = tx.Rollback(ctx) }()

	for _, p := range positions {
		if _, err := tx.Exec(ctx, `
			UPDATE modules SET sort_order = $1, updated_at = NOW()
			WHERE id = $2 AND organization_id = $3
		`, p.SortOrder, p.ID, orgID); err != nil {
			return err
		}
	}

	return tx.Commit(ctx)
}

func collect(rows pgx.Rows) ([]entity.Module, error) {
	modules := make([]entity.Module, 0)
	for rows.Next() {
		var m entity.Module
		if err := scanModule(rows, &m); err != nil {
			return nil, err
		}
		modules = append(modules, m)
	}
	return modules, rows.Err()
}
