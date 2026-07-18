package repository

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/abhinavkumar03/crm-lite/backend/internal/exporter/dto"
	"github.com/abhinavkumar03/crm-lite/backend/internal/exporter/entity"
)

// exportColumns omits content: the file bytes are large and only needed by the
// download endpoint, which fetches them separately via GetForDownload.
const exportColumns = `
	id, organization_id, module_id, filename, format, status, columns, filters,
	options, row_count, byte_size, error, created_by, started_at, finished_at,
	created_at, updated_at
`

type ExportRepository struct {
	db *pgxpool.Pool
}

func NewExportRepository(db *pgxpool.Pool) *ExportRepository {
	return &ExportRepository{db: db}
}

func scanExport(row pgx.Row, j *entity.ExportJob) error {
	return row.Scan(
		&j.ID, &j.OrganizationID, &j.ModuleID, &j.Filename, &j.Format, &j.Status,
		&j.Columns, &j.Filters, &j.Options, &j.RowCount, &j.ByteSize, &j.Error,
		&j.CreatedBy, &j.StartedAt, &j.FinishedAt, &j.CreatedAt, &j.UpdatedAt,
	)
}

func (r *ExportRepository) Create(ctx context.Context, j *entity.ExportJob) error {
	return r.db.QueryRow(ctx, `
		INSERT INTO export_jobs (
			organization_id, module_id, filename, format, status, columns, filters,
			options, created_by
		)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9)
		RETURNING id, status, created_at, updated_at
	`,
		j.OrganizationID, j.ModuleID, j.Filename, j.Format, j.Status, j.Columns,
		j.Filters, j.Options, j.CreatedBy,
	).Scan(&j.ID, &j.Status, &j.CreatedAt, &j.UpdatedAt)
}

func (r *ExportRepository) GetByID(ctx context.Context, orgID, id string) (*entity.ExportJob, error) {
	var j entity.ExportJob
	err := scanExport(r.db.QueryRow(ctx, `
		SELECT `+exportColumns+`
		FROM export_jobs
		WHERE id = $1 AND organization_id = $2
	`, id, orgID), &j)

	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &j, nil
}

// GetForDownload returns just what the download endpoint needs. content is nil
// unless the job has completed.
func (r *ExportRepository) GetForDownload(ctx context.Context, orgID, id string) (filename, format, status string, content []byte, found bool, err error) {
	err = r.db.QueryRow(ctx, `
		SELECT filename, format, status, content
		FROM export_jobs
		WHERE id = $1 AND organization_id = $2
	`, id, orgID).Scan(&filename, &format, &status, &content)

	if errors.Is(err, pgx.ErrNoRows) {
		return "", "", "", nil, false, nil
	}
	if err != nil {
		return "", "", "", nil, false, err
	}
	return filename, format, status, content, true, nil
}

func (r *ExportRepository) List(ctx context.Context, orgID, moduleID string, q dto.ListQuery) ([]entity.ExportJob, int, error) {
	conds := []string{"organization_id = $1", "module_id = $2"}
	args := []any{orgID, moduleID}
	next := 3

	if q.Status != "" {
		conds = append(conds, fmt.Sprintf("status = $%d", next))
		args = append(args, q.Status)
		next++
	}

	where := "WHERE " + strings.Join(conds, " AND ")

	var total int
	if err := r.db.QueryRow(ctx, "SELECT COUNT(*) FROM export_jobs "+where, args...).Scan(&total); err != nil {
		return nil, 0, err
	}

	limitPH := fmt.Sprintf("$%d", next)
	offsetPH := fmt.Sprintf("$%d", next+1)
	args = append(args, q.PageSize, (q.Page-1)*q.PageSize)

	rows, err := r.db.Query(ctx,
		"SELECT "+exportColumns+" FROM export_jobs "+where+
			" ORDER BY created_at DESC LIMIT "+limitPH+" OFFSET "+offsetPH, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	items := make([]entity.ExportJob, 0)
	for rows.Next() {
		var j entity.ExportJob
		if err := scanExport(rows, &j); err != nil {
			return nil, 0, err
		}
		items = append(items, j)
	}
	return items, total, rows.Err()
}

func (r *ExportRepository) MarkProcessing(ctx context.Context, id string) error {
	_, err := r.db.Exec(ctx, `
		UPDATE export_jobs
		SET status = 'processing', started_at = NOW(), updated_at = NOW()
		WHERE id = $1
	`, id)
	return err
}

// Complete stores the generated file and its size/row count.
func (r *ExportRepository) Complete(ctx context.Context, id string, rowCount int, content []byte) error {
	_, err := r.db.Exec(ctx, `
		UPDATE export_jobs
		SET status = 'completed', row_count = $2, byte_size = $3, content = $4,
		    error = NULL, finished_at = NOW(), updated_at = NOW()
		WHERE id = $1
	`, id, rowCount, len(content), content)
	return err
}

func (r *ExportRepository) MarkFailed(ctx context.Context, id, reason string) error {
	_, err := r.db.Exec(ctx, `
		UPDATE export_jobs
		SET status = 'failed', error = $2, finished_at = NOW(), updated_at = NOW()
		WHERE id = $1
	`, id, reason)
	return err
}
