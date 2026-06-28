package repository

import (
	"context"

	"github.com/abhinavkumar03/crm-lite/backend/internal/contact/dto"
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

func (r *Repository) List(
	ctx context.Context,
	ownerID string,
	req dto.ListContactsRequest,
) ([]entity.Contact, error) {

	offset := (req.Page - 1) * req.Limit

	query := `
	SELECT
		id,
		owner_id,
		first_name,
		last_name,
		email,
		phone,
		company,
		job_title,
		notes,
		created_at,
		updated_at
	FROM contacts
	WHERE owner_id = $1
	AND (
		$2 = ''
		OR first_name ILIKE '%' || $2 || '%'
		OR last_name ILIKE '%' || $2 || '%'
		OR email ILIKE '%' || $2 || '%'
		OR company ILIKE '%' || $2 || '%'
	)
	ORDER BY created_at DESC
	LIMIT $3
	OFFSET $4;
	`

	rows, err := r.db.Query(
		ctx,
		query,
		ownerID,
		req.Search,
		req.Limit,
		offset,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	contacts := make([]entity.Contact, 0)

	for rows.Next() {
		var contact entity.Contact

		if err := rows.Scan(
			&contact.ID,
			&contact.OwnerID,
			&contact.FirstName,
			&contact.LastName,
			&contact.Email,
			&contact.Phone,
			&contact.Company,
			&contact.JobTitle,
			&contact.Notes,
			&contact.CreatedAt,
			&contact.UpdatedAt,
		); err != nil {
			return nil, err
		}

		contacts = append(contacts, contact)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return contacts, nil
}
