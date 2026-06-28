package repository

import (
	"context"
	"errors"

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
	req dto.ListLeadsRequest,
) ([]entity.Lead, error) {

	offset := (req.Page - 1) * req.Limit

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
	WHERE owner_id = $1
	AND (
		$2 = ''
		OR
		name ILIKE '%' || $2 || '%'
		OR
		email ILIKE '%' || $2 || '%'
	)
	AND (
		$3 = ''
		OR
		status = $3
	)
	ORDER BY created_at DESC
	LIMIT $4
	OFFSET $5;
	`

	rows, err := r.db.Query(
		ctx,
		query,
		ownerID,
		req.Search,
		req.Status,
		req.Limit,
		offset,
	)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var leads []entity.Lead

	for rows.Next() {

		var lead entity.Lead

		err = rows.Scan(
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

		leads = append(leads, lead)
	}

	return leads, rows.Err()
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
