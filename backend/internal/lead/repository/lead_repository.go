package repository

import (
	"context"

	"github.com/abhinavkumar03/crm-lite/backend/internal/lead/entity"
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
