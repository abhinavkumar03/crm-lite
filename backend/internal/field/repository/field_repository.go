package repository

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/abhinavkumar03/crm-lite/backend/internal/field/entity"
)

const fieldColumns = `
	id, organization_id, module_id, api_name, label, field_type,
	is_required, is_unique, is_read_only, default_value, placeholder,
	description, help_text, min_length, max_length, regex, validation_message,
	options, lookup_module_id, sort_order, is_visible, is_searchable,
	is_filterable, is_nullable, is_indexed, is_system,
	lock_mode, editable_by, viewable_by, created_at, updated_at
`

type Repository struct {
	db *pgxpool.Pool
}

func New(db *pgxpool.Pool) *Repository {
	return &Repository{db: db}
}

func scanField(row pgx.Row, f *entity.Field) error {
	return row.Scan(
		&f.ID, &f.OrganizationID, &f.ModuleID, &f.APIName, &f.Label, &f.FieldType,
		&f.IsRequired, &f.IsUnique, &f.IsReadOnly, &f.DefaultValue, &f.Placeholder,
		&f.Description, &f.HelpText, &f.MinLength, &f.MaxLength, &f.Regex, &f.ValidationMessage,
		&f.Options, &f.LookupModuleID, &f.SortOrder, &f.IsVisible, &f.IsSearchable,
		&f.IsFilterable, &f.IsNullable, &f.IsIndexed, &f.IsSystem,
		&f.LockMode, &f.EditableBy, &f.ViewableBy, &f.CreatedAt, &f.UpdatedAt,
	)
}

// ModuleStorage returns the module's storage strategy and whether the module
// exists within the given organization. Used to validate ownership and to
// derive each field's persistence strategy.
func (r *Repository) ModuleStorage(ctx context.Context, orgID, moduleID string) (string, bool, error) {
	var strategy string
	err := r.db.QueryRow(ctx, `
		SELECT storage_strategy FROM modules
		WHERE id = $1 AND organization_id = $2
	`, moduleID, orgID).Scan(&strategy)

	if errors.Is(err, pgx.ErrNoRows) {
		return "", false, nil
	}
	if err != nil {
		return "", false, err
	}
	return strategy, true, nil
}

// ModuleExistsInOrg reports whether a module id belongs to the organization
// (used to validate lookup targets).
func (r *Repository) ModuleExistsInOrg(ctx context.Context, orgID, moduleID string) (bool, error) {
	var exists bool
	err := r.db.QueryRow(ctx, `
		SELECT EXISTS(SELECT 1 FROM modules WHERE id = $1 AND organization_id = $2)
	`, moduleID, orgID).Scan(&exists)
	return exists, err
}

func (r *Repository) Create(ctx context.Context, f *entity.Field) error {
	query := `
		INSERT INTO fields (
			organization_id, module_id, api_name, label, field_type,
			is_required, is_unique, is_read_only, default_value, placeholder,
			description, help_text, min_length, max_length, regex, validation_message,
			options, lookup_module_id, sort_order, is_visible, is_searchable,
			is_filterable, is_system, lock_mode, editable_by, viewable_by
		)
		VALUES (
			$1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,$15,$16,
			$17,$18,$19,$20,$21,$22,$23,$24,$25,$26
		)
		RETURNING id, is_nullable, is_indexed, created_at, updated_at
	`
	return r.db.QueryRow(ctx, query,
		f.OrganizationID, f.ModuleID, f.APIName, f.Label, f.FieldType,
		f.IsRequired, f.IsUnique, f.IsReadOnly, f.DefaultValue, f.Placeholder,
		f.Description, f.HelpText, f.MinLength, f.MaxLength, f.Regex, f.ValidationMessage,
		f.Options, f.LookupModuleID, f.SortOrder, f.IsVisible, f.IsSearchable,
		f.IsFilterable, f.IsSystem, f.LockMode, f.EditableBy, f.ViewableBy,
	).Scan(&f.ID, &f.IsNullable, &f.IsIndexed, &f.CreatedAt, &f.UpdatedAt)
}

func (r *Repository) List(ctx context.Context, orgID, moduleID string) ([]entity.Field, error) {
	rows, err := r.db.Query(ctx, `
		SELECT `+fieldColumns+`
		FROM fields
		WHERE organization_id = $1 AND module_id = $2
		ORDER BY sort_order ASC, label ASC
	`, orgID, moduleID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	fields := make([]entity.Field, 0)
	for rows.Next() {
		var f entity.Field
		if err := scanField(rows, &f); err != nil {
			return nil, err
		}
		fields = append(fields, f)
	}
	return fields, rows.Err()
}

func (r *Repository) GetByID(ctx context.Context, orgID, moduleID, id string) (*entity.Field, error) {
	var f entity.Field
	err := scanField(r.db.QueryRow(ctx, `
		SELECT `+fieldColumns+`
		FROM fields
		WHERE id = $1 AND module_id = $2 AND organization_id = $3
	`, id, moduleID, orgID), &f)

	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &f, nil
}

func (r *Repository) Update(ctx context.Context, f *entity.Field) error {
	return r.db.QueryRow(ctx, `
		UPDATE fields SET
			label = $1,
			is_required = $2,
			is_unique = $3,
			is_read_only = $4,
			default_value = $5,
			placeholder = $6,
			description = $7,
			help_text = $8,
			min_length = $9,
			max_length = $10,
			regex = $11,
			validation_message = $12,
			options = $13,
			is_visible = $14,
			is_searchable = $15,
			is_filterable = $16,
			lock_mode = $17,
			editable_by = $18,
			viewable_by = $19,
			updated_at = NOW()
		WHERE id = $20 AND module_id = $21 AND organization_id = $22
		RETURNING updated_at
	`,
		f.Label, f.IsRequired, f.IsUnique, f.IsReadOnly, f.DefaultValue,
		f.Placeholder, f.Description, f.HelpText, f.MinLength, f.MaxLength,
		f.Regex, f.ValidationMessage, f.Options, f.IsVisible, f.IsSearchable,
		f.IsFilterable, f.LockMode, f.EditableBy, f.ViewableBy,
		f.ID, f.ModuleID, f.OrganizationID,
	).Scan(&f.UpdatedAt)
}

func (r *Repository) Delete(ctx context.Context, orgID, moduleID, id string) (bool, error) {
	tag, err := r.db.Exec(ctx, `
		DELETE FROM fields WHERE id = $1 AND module_id = $2 AND organization_id = $3
	`, id, moduleID, orgID)
	if err != nil {
		return false, err
	}
	return tag.RowsAffected() > 0, nil
}

func (r *Repository) ExistsByAPIName(ctx context.Context, moduleID, apiName string) (bool, error) {
	var exists bool
	err := r.db.QueryRow(ctx, `
		SELECT EXISTS(
			SELECT 1 FROM fields WHERE module_id = $1 AND api_name = $2
		)
	`, moduleID, apiName).Scan(&exists)
	return exists, err
}

func (r *Repository) MaxSortOrder(ctx context.Context, moduleID string) (int, error) {
	var max int
	err := r.db.QueryRow(ctx, `
		SELECT COALESCE(MAX(sort_order), 0) FROM fields WHERE module_id = $1
	`, moduleID).Scan(&max)
	return max, err
}

// Reorder updates sort_order for the given fields atomically.
func (r *Repository) Reorder(ctx context.Context, orgID, moduleID string, positions []entity.SortPosition) error {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return err
	}
	defer func() { _ = tx.Rollback(ctx) }()

	for _, p := range positions {
		if _, err := tx.Exec(ctx, `
			UPDATE fields SET sort_order = $1, updated_at = NOW()
			WHERE id = $2 AND module_id = $3 AND organization_id = $4
		`, p.SortOrder, p.ID, moduleID, orgID); err != nil {
			return err
		}
	}

	return tx.Commit(ctx)
}
