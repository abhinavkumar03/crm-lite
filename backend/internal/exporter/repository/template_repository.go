package repository

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/abhinavkumar03/crm-lite/backend/internal/exporter/entity"
)

const templateColumns = `
	id, organization_id, module_id, name, format, columns, filters, sort,
	created_by, created_at, updated_at
`

type TemplateRepository struct {
	db *pgxpool.Pool
}

func NewTemplateRepository(db *pgxpool.Pool) *TemplateRepository {
	return &TemplateRepository{db: db}
}

func scanTemplate(row pgx.Row, t *entity.ExportTemplate) error {
	return row.Scan(
		&t.ID, &t.OrganizationID, &t.ModuleID, &t.Name, &t.Format, &t.Columns,
		&t.Filters, &t.Sort, &t.CreatedBy, &t.CreatedAt, &t.UpdatedAt,
	)
}

func (r *TemplateRepository) Create(ctx context.Context, t *entity.ExportTemplate) error {
	return r.db.QueryRow(ctx, `
		INSERT INTO export_templates (
			organization_id, module_id, name, format, columns, filters, sort, created_by
		)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8)
		RETURNING id, created_at, updated_at
	`,
		t.OrganizationID, t.ModuleID, t.Name, t.Format, t.Columns, t.Filters, t.Sort, t.CreatedBy,
	).Scan(&t.ID, &t.CreatedAt, &t.UpdatedAt)
}

func (r *TemplateRepository) List(ctx context.Context, orgID, moduleID string) ([]entity.ExportTemplate, error) {
	rows, err := r.db.Query(ctx, `
		SELECT `+templateColumns+`
		FROM export_templates
		WHERE organization_id = $1 AND module_id = $2
		ORDER BY name ASC
	`, orgID, moduleID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	items := make([]entity.ExportTemplate, 0)
	for rows.Next() {
		var t entity.ExportTemplate
		if err := scanTemplate(rows, &t); err != nil {
			return nil, err
		}
		items = append(items, t)
	}
	return items, rows.Err()
}

func (r *TemplateRepository) GetByID(ctx context.Context, orgID, moduleID, id string) (*entity.ExportTemplate, error) {
	var t entity.ExportTemplate
	err := scanTemplate(r.db.QueryRow(ctx, `
		SELECT `+templateColumns+`
		FROM export_templates
		WHERE id = $1 AND module_id = $2 AND organization_id = $3
	`, id, moduleID, orgID), &t)

	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &t, nil
}

func (r *TemplateRepository) Update(ctx context.Context, t *entity.ExportTemplate) error {
	return r.db.QueryRow(ctx, `
		UPDATE export_templates
		SET name = $1, format = $2, columns = $3, filters = $4, sort = $5, updated_at = NOW()
		WHERE id = $6 AND module_id = $7 AND organization_id = $8
		RETURNING updated_at
	`,
		t.Name, t.Format, t.Columns, t.Filters, t.Sort, t.ID, t.ModuleID, t.OrganizationID,
	).Scan(&t.UpdatedAt)
}

func (r *TemplateRepository) Delete(ctx context.Context, orgID, moduleID, id string) (bool, error) {
	tag, err := r.db.Exec(ctx, `
		DELETE FROM export_templates WHERE id = $1 AND module_id = $2 AND organization_id = $3
	`, id, moduleID, orgID)
	if err != nil {
		return false, err
	}
	return tag.RowsAffected() > 0, nil
}
