package repository

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/abhinavkumar03/crm-lite/backend/internal/importer/dto"
	"github.com/abhinavkumar03/crm-lite/backend/internal/importer/entity"
)

// importColumns intentionally omits source_rows: it can be large and is only
// needed by the worker, which fetches it separately via SourceRows.
const importColumns = `
	id, organization_id, module_id, filename, status, mapping, options,
	total_rows, processed_rows, success_rows, error_rows, errors,
	created_by, started_at, finished_at, created_at, updated_at
`

type Repository struct {
	db *pgxpool.Pool
}

func New(db *pgxpool.Pool) *Repository {
	return &Repository{db: db}
}

func scan(row pgx.Row, j *entity.ImportJob) error {
	return row.Scan(
		&j.ID, &j.OrganizationID, &j.ModuleID, &j.Filename, &j.Status, &j.Mapping,
		&j.Options, &j.TotalRows, &j.ProcessedRows, &j.SuccessRows, &j.ErrorRows,
		&j.Errors, &j.CreatedBy, &j.StartedAt, &j.FinishedAt, &j.CreatedAt, &j.UpdatedAt,
	)
}

func (r *Repository) Create(ctx context.Context, j *entity.ImportJob) error {
	return r.db.QueryRow(ctx, `
		INSERT INTO import_jobs (
			organization_id, module_id, filename, status, mapping, options,
			source_rows, total_rows, created_by
		)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9)
		RETURNING id, status, created_at, updated_at
	`,
		j.OrganizationID, j.ModuleID, j.Filename, j.Status, j.Mapping, j.Options,
		j.SourceRows, j.TotalRows, j.CreatedBy,
	).Scan(&j.ID, &j.Status, &j.CreatedAt, &j.UpdatedAt)
}

func (r *Repository) GetByID(ctx context.Context, orgID, id string) (*entity.ImportJob, error) {
	var j entity.ImportJob
	err := scan(r.db.QueryRow(ctx, `
		SELECT `+importColumns+`
		FROM import_jobs
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

func (r *Repository) List(ctx context.Context, orgID, moduleID string, q dto.ListQuery) ([]entity.ImportJob, int, error) {
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
	if err := r.db.QueryRow(ctx, "SELECT COUNT(*) FROM import_jobs "+where, args...).Scan(&total); err != nil {
		return nil, 0, err
	}

	limitPH := fmt.Sprintf("$%d", next)
	offsetPH := fmt.Sprintf("$%d", next+1)
	args = append(args, q.PageSize, (q.Page-1)*q.PageSize)

	rows, err := r.db.Query(ctx,
		"SELECT "+importColumns+" FROM import_jobs "+where+
			" ORDER BY created_at DESC LIMIT "+limitPH+" OFFSET "+offsetPH, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	items := make([]entity.ImportJob, 0)
	for rows.Next() {
		var j entity.ImportJob
		if err := scan(rows, &j); err != nil {
			return nil, 0, err
		}
		items = append(items, j)
	}
	return items, total, rows.Err()
}

// SourceRows returns the staged, parsed rows for a job (worker only).
func (r *Repository) SourceRows(ctx context.Context, id string) ([]byte, error) {
	var raw []byte
	err := r.db.QueryRow(ctx, `SELECT source_rows FROM import_jobs WHERE id = $1`, id).Scan(&raw)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}
	return raw, err
}

// MarkProcessing transitions a queued job to processing and stamps started_at.
func (r *Repository) MarkProcessing(ctx context.Context, id string) error {
	_, err := r.db.Exec(ctx, `
		UPDATE import_jobs
		SET status = 'processing', started_at = NOW(), updated_at = NOW()
		WHERE id = $1
	`, id)
	return err
}

// Finish records the terminal outcome: final status, counters and the per-row
// error report.
func (r *Repository) Finish(
	ctx context.Context,
	id, status string,
	processed, success, errorRows int,
	errs []byte,
) error {
	_, err := r.db.Exec(ctx, `
		UPDATE import_jobs
		SET status = $2, processed_rows = $3, success_rows = $4, error_rows = $5,
		    errors = $6, finished_at = NOW(), updated_at = NOW()
		WHERE id = $1
	`, id, status, processed, success, errorRows, errs)
	return err
}
