package repository

import (
	"context"
	"errors"
	"fmt"
	"math"
	"strings"

	"github.com/abhinavkumar03/crm-lite/backend/internal/lead/dto"
	"github.com/abhinavkumar03/crm-lite/backend/internal/lead/entity"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository struct {
	db *pgxpool.Pool
}

func New(db *pgxpool.Pool) *Repository {
	return &Repository{
		db: db,
	}
}

func (r *Repository) Create(
	ctx context.Context,
	lead *entity.Lead,
) error {

	query := `
	INSERT INTO leads
	(
		owner_id,
		name,
		email,
		phone,
		company,
		status,
		notes
	)
	VALUES
	(
		$1,$2,$3,$4,$5,$6,$7
	)
	RETURNING
		id,
		created_at,
		updated_at;
	`

	return r.db.QueryRow(
		ctx,
		query,
		lead.OwnerID,
		lead.Name,
		lead.Email,
		lead.Phone,
		lead.Company,
		lead.Status,
		lead.Notes,
	).Scan(
		&lead.ID,
		&lead.CreatedAt,
		&lead.UpdatedAt,
	)
}

func (r *Repository) List(
	ctx context.Context,
	ownerID string,
	req dto.ListLeadRequest,
) (*dto.ListLeadResponse, error) {

	conditions := []string{
		"owner_id = $1",
	}

	args := []interface{}{
		ownerID,
	}

	argIndex := 2

	if req.Search != "" {

		conditions = append(
			conditions,
			fmt.Sprintf(`(
LOWER(name) LIKE $%d
OR LOWER(email) LIKE $%d
OR LOWER(phone) LIKE $%d
OR LOWER(company) LIKE $%d
)`,
				argIndex,
				argIndex,
				argIndex,
				argIndex,
			),
		)

		args = append(
			args,
			"%"+strings.ToLower(req.Search)+"%",
		)

		argIndex++
	}

	if req.Status != "" {

		conditions = append(
			conditions,
			fmt.Sprintf(
				"status = $%d",
				argIndex,
			),
		)

		args = append(
			args,
			req.Status,
		)

		argIndex++
	}

	whereClause := strings.Join(
		conditions,
		" AND ",
	)

	var total int64

	countQuery := fmt.Sprintf(`
SELECT COUNT(*)
FROM leads
WHERE %s
`, whereClause)

	err := r.db.QueryRow(
		ctx,
		countQuery,
		args...,
	).Scan(&total)

	if err != nil {
		return nil, err
	}

	allowedSorts := map[string]string{
		"name":       "name",
		"company":    "company",
		"status":     "status",
		"created_at": "created_at",
	}

	sortBy := "created_at"

	if v, ok := allowedSorts[req.SortBy]; ok {
		sortBy = v
	}

	order := "DESC"

	if strings.EqualFold(req.SortOrder, "ASC") {
		order = "ASC"
	}

	offset := (req.Page - 1) * req.Limit

	query := fmt.Sprintf(`
SELECT
	id,
	name,
	email,
	phone,
	company,
	status,
	notes,
	owner_id,
	created_at,
	updated_at
FROM leads
WHERE %s
ORDER BY %s %s
LIMIT %d
OFFSET %d
`,
		whereClause,
		sortBy,
		order,
		req.Limit,
		offset,
	)

	rows, err := r.db.Query(
		ctx,
		query,
		args...,
	)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	leads := make([]dto.LeadResponse, 0)

	for rows.Next() {

		var lead dto.LeadResponse

		err := rows.Scan(
			&lead.ID,
			&lead.Name,
			&lead.Email,
			&lead.Phone,
			&lead.Company,
			&lead.Status,
			&lead.Notes,
			&lead.OwnerID,
			&lead.CreatedAt,
			&lead.UpdatedAt,
		)

		if err != nil {
			return nil, err
		}

		leads = append(
			leads,
			lead,
		)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	totalPages := 0

	if total > 0 {
		totalPages = int(
			math.Ceil(
				float64(total) /
					float64(req.Limit),
			),
		)
	}

	return &dto.ListLeadResponse{
		Data:       leads,
		Page:       req.Page,
		Limit:      req.Limit,
		Total:      total,
		TotalPages: totalPages,
	}, nil
}

func (r *Repository) GetByID(
	ctx context.Context,
	id string,
	ownerID string,
) (*entity.Lead, error) {

	query := `
	SELECT
		id,
		name,
		email,
		phone,
		company,
		status,
		notes,
		owner_id,
		created_at,
		updated_at
	FROM leads
	WHERE id = $1
	AND owner_id = $2;
	`

	var lead entity.Lead

	err := r.db.QueryRow(
		ctx,
		query,
		id,
		ownerID,
	).Scan(
		&lead.ID,
		&lead.Name,
		&lead.Email,
		&lead.Phone,
		&lead.Company,
		&lead.Status,
		&lead.Notes,
		&lead.OwnerID,
		&lead.CreatedAt,
		&lead.UpdatedAt,
	)

	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	return &lead, nil
}

func (r *Repository) Update(
	ctx context.Context,
	lead *entity.Lead,
) error {

	query := `
	UPDATE leads
	SET
		name = $1,
		email = $2,
		phone = $3,
		company = $4,
		status = $5,
		notes = $6,
		updated_at = NOW()
	WHERE id = $7
	AND owner_id = $8
	RETURNING updated_at;
	`

	return r.db.QueryRow(
		ctx,
		query,
		lead.Name,
		lead.Email,
		lead.Phone,
		lead.Company,
		lead.Status,
		lead.Notes,
		lead.ID,
		lead.OwnerID,
	).Scan(
		&lead.UpdatedAt,
	)
}

func (r *Repository) Delete(
	ctx context.Context,
	id string,
	ownerID string,
) (bool, error) {

	query := `
	DELETE FROM leads
	WHERE id = $1
	AND owner_id = $2;
	`

	result, err := r.db.Exec(
		ctx,
		query,
		id,
		ownerID,
	)

	if err != nil {
		return false, err
	}

	return result.RowsAffected() > 0, nil
}

func (r *Repository) Search(
	ctx context.Context,
	ownerID string,
	query string,
) ([]dto.LeadResponse, error) {

	search := "%" + strings.ToLower(query) + "%"

	rows, err := r.db.Query(
		ctx,
		`
SELECT
	id,
	name,
	email,
	phone,
	company,
	status,
	notes,
	owner_id,
	created_at,
	updated_at
FROM leads
WHERE owner_id=$1
AND (

	LOWER(name) LIKE $2

	OR LOWER(email) LIKE $2

	OR LOWER(phone) LIKE $2

	OR LOWER(company) LIKE $2

)

ORDER BY created_at DESC

LIMIT 10
`,
		ownerID,
		search,
	)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	leads := make([]dto.LeadResponse, 0)

	for rows.Next() {

		var lead dto.LeadResponse

		err := rows.Scan(

			&lead.ID,

			&lead.Name,

			&lead.Email,

			&lead.Phone,

			&lead.Company,

			&lead.Status,

			&lead.Notes,

			&lead.OwnerID,

			&lead.CreatedAt,

			&lead.UpdatedAt,
		)

		if err != nil {

			return nil, err

		}

		leads = append(
			leads,
			lead,
		)

	}

	return leads, nil

}

func (r *Repository) ExistsByIDAndOwner(
	ctx context.Context,
	id string,
	ownerID string,
) (bool, error) {

	query := `
	SELECT EXISTS(
		SELECT 1
		FROM leads
		WHERE id = $1
		AND owner_id = $2
	);
	`

	var exists bool

	err := r.db.QueryRow(
		ctx,
		query,
		id,
		ownerID,
	).Scan(&exists)

	if err != nil {
		return false, err
	}

	return exists, nil
}
