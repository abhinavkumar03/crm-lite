package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/abhinavkumar03/crm-lite/backend/internal/record/dto"
	"github.com/abhinavkumar03/crm-lite/backend/internal/record/entity"
)

const recordColumns = `
	id, organization_id, module_id, data, owner_id, created_by, updated_by,
	created_at, updated_at
`

type Repository struct {
	db *pgxpool.Pool
}

func New(db *pgxpool.Pool) *Repository {
	return &Repository{db: db}
}

func scanRecord(row pgx.Row, r *entity.Record) error {
	return row.Scan(
		&r.ID, &r.OrganizationID, &r.ModuleID, &r.Data, &r.OwnerID,
		&r.CreatedBy, &r.UpdatedBy, &r.CreatedAt, &r.UpdatedAt,
	)
}

func (r *Repository) Create(ctx context.Context, rec *entity.Record) error {
	return r.db.QueryRow(ctx, `
		INSERT INTO records (organization_id, module_id, data, owner_id, created_by, updated_by)
		VALUES ($1,$2,$3,$4,$5,$6)
		RETURNING id, created_at, updated_at
	`,
		rec.OrganizationID, rec.ModuleID, rec.Data, rec.OwnerID, rec.CreatedBy, rec.UpdatedBy,
	).Scan(&rec.ID, &rec.CreatedAt, &rec.UpdatedAt)
}

// CreateBatch inserts many records in one COPY. IDs are assigned by the DB
// default and are not populated on the input structs (import does not need them).
func (r *Repository) CreateBatch(ctx context.Context, recs []*entity.Record) error {
	if len(recs) == 0 {
		return nil
	}
	rows := make([][]any, len(recs))
	for i, rec := range recs {
		rows[i] = []any{
			rec.OrganizationID, rec.ModuleID, rec.Data,
			rec.OwnerID, rec.CreatedBy, rec.UpdatedBy,
		}
	}
	_, err := r.db.CopyFrom(
		ctx,
		pgx.Identifier{"records"},
		[]string{"organization_id", "module_id", "data", "owner_id", "created_by", "updated_by"},
		pgx.CopyFromRows(rows),
	)
	return err
}

func (r *Repository) GetByID(ctx context.Context, orgID, moduleID, id string) (*entity.Record, error) {
	var rec entity.Record
	err := scanRecord(r.db.QueryRow(ctx, `
		SELECT `+recordColumns+`
		FROM records
		WHERE id = $1 AND module_id = $2 AND organization_id = $3
	`, id, moduleID, orgID), &rec)

	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &rec, nil
}

func (r *Repository) Update(ctx context.Context, rec *entity.Record) error {
	return r.db.QueryRow(ctx, `
		UPDATE records SET
			data = $1,
			owner_id = $2,
			updated_by = $3,
			updated_at = NOW()
		WHERE id = $4 AND module_id = $5 AND organization_id = $6
		RETURNING updated_at
	`,
		rec.Data, rec.OwnerID, rec.UpdatedBy, rec.ID, rec.ModuleID, rec.OrganizationID,
	).Scan(&rec.UpdatedAt)
}

func (r *Repository) Delete(ctx context.Context, orgID, moduleID, id string) (bool, error) {
	tag, err := r.db.Exec(ctx, `
		DELETE FROM records WHERE id = $1 AND module_id = $2 AND organization_id = $3
	`, id, moduleID, orgID)
	if err != nil {
		return false, err
	}
	return tag.RowsAffected() > 0, nil
}

// List runs the dynamic query (search + filters + sort + pagination) and returns
// the page of records plus the total count for the same filter set. When
// q.SkipTotal is set, COUNT(*) is skipped and total is 0.
func (r *Repository) List(
	ctx context.Context,
	orgID, moduleID string,
	q dto.ListQuery,
	meta map[string]FieldMeta,
) ([]entity.Record, int, error) {
	where := BuildWhere(orgID, moduleID, q, meta)
	order := BuildOrderBy(q.Sort, q.Order, meta)

	var total int
	if !q.SkipTotal {
		if err := r.db.QueryRow(ctx,
			"SELECT COUNT(*) FROM records WHERE "+where.SQL, where.Args...,
		).Scan(&total); err != nil {
			return nil, 0, err
		}
	}

	limitPH := fmt.Sprintf("$%d", len(where.Args)+1)
	offsetPH := fmt.Sprintf("$%d", len(where.Args)+2)
	args := append(append([]any{}, where.Args...), q.PageSize, (q.Page-1)*q.PageSize)

	rows, err := r.db.Query(ctx,
		"SELECT "+recordColumns+" FROM records WHERE "+where.SQL+" "+order+
			" LIMIT "+limitPH+" OFFSET "+offsetPH, args...,
	)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	records := make([]entity.Record, 0)
	for rows.Next() {
		var rec entity.Record
		if err := scanRecord(rows, &rec); err != nil {
			return nil, 0, err
		}
		records = append(records, rec)
	}
	return records, total, rows.Err()
}

// DisplayValues resolves a display label for referenced records of a lookup
// target module. If displayField is not a valid identifier, the id is used.
func (r *Repository) DisplayValues(
	ctx context.Context,
	orgID, moduleID string,
	ids []string,
	displayField string,
) (map[string]string, error) {
	out := make(map[string]string)
	if len(ids) == 0 {
		return out, nil
	}

	labelExpr := "id::text"
	if apiNameRe.MatchString(displayField) {
		labelExpr = fmt.Sprintf("COALESCE(NULLIF(data->>'%s', ''), id::text)", displayField)
	}

	rows, err := r.db.Query(ctx, fmt.Sprintf(`
		SELECT id::text, %s
		FROM records
		WHERE organization_id = $1 AND module_id = $2 AND id = ANY($3)
	`, labelExpr), orgID, moduleID, ids)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var id, label string
		if err := rows.Scan(&id, &label); err != nil {
			return nil, err
		}
		out[id] = label
	}
	return out, rows.Err()
}

// UserDisplays resolves user ids to their display name.
func (r *Repository) UserDisplays(ctx context.Context, ids []string) (map[string]string, error) {
	out := make(map[string]string)
	if len(ids) == 0 {
		return out, nil
	}

	rows, err := r.db.Query(ctx, `
		SELECT id::text, name FROM users WHERE id = ANY($1)
	`, ids)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var id, name string
		if err := rows.Scan(&id, &name); err != nil {
			return nil, err
		}
		out[id] = name
	}
	return out, rows.Err()
}
