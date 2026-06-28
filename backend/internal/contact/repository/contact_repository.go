package repository

import (
	"context"

	"github.com/abhinavkumar03/crm-lite/backend/internal/contact/entity"
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
	contact *entity.Contact,
) error {

	query := `
	INSERT INTO contacts
	(
		owner_id,
		first_name,
		last_name,
		email,
		phone,
		company,
		job_title,
		notes
	)
	VALUES
	(
		$1,$2,$3,$4,$5,$6,$7,$8
	)
	RETURNING
		id,
		created_at,
		updated_at;
	`

	return r.db.QueryRow(
		ctx,
		query,
		contact.OwnerID,
		contact.FirstName,
		contact.LastName,
		contact.Email,
		contact.Phone,
		contact.Company,
		contact.JobTitle,
		contact.Notes,
	).Scan(
		&contact.ID,
		&contact.CreatedAt,
		&contact.UpdatedAt,
	)
}
