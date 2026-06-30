package repository

import (
	"context"
	"errors"
	"strings"

	"github.com/abhinavkumar03/crm-lite/backend/internal/contact/dto"
	"github.com/abhinavkumar03/crm-lite/backend/internal/contact/entity"
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

func (r *Repository) GetByID(
	ctx context.Context,
	id string,
	ownerID string,
) (*entity.Contact, error) {

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
	WHERE id = $1
	AND owner_id = $2;
	`

	var contact entity.Contact

	err := r.db.QueryRow(
		ctx,
		query,
		id,
		ownerID,
	).Scan(
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
	)

	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	return &contact, nil
}

func (r *Repository) Update(
	ctx context.Context,
	contact *entity.Contact,
) error {

	query := `
	UPDATE contacts
	SET
		first_name = $1,
		last_name = $2,
		email = $3,
		phone = $4,
		company = $5,
		job_title = $6,
		notes = $7,
		updated_at = NOW()
	WHERE id = $8
	AND owner_id = $9
	RETURNING updated_at;
	`

	return r.db.QueryRow(
		ctx,
		query,
		contact.FirstName,
		contact.LastName,
		contact.Email,
		contact.Phone,
		contact.Company,
		contact.JobTitle,
		contact.Notes,
		contact.ID,
		contact.OwnerID,
	).Scan(
		&contact.UpdatedAt,
	)
}

func (r *Repository) Delete(
	ctx context.Context,
	id string,
	ownerID string,
) (bool, error) {

	query := `
	DELETE FROM contacts
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
) ([]dto.ContactResponse, error) {

	search := "%" + strings.ToLower(query) + "%"

	rows, err := r.db.Query(
		ctx, `SELECT
			id,
			first_name,
			last_name,
			email,
			phone,
			company,
			job_title,
			notes,
			owner_id,
			created_at,
			updated_at
		FROM contacts
		WHERE owner_id = $1
		AND (
			LOWER(first_name) LIKE $2
			OR LOWER(last_name) LIKE $2
			OR LOWER(email) LIKE $2
			OR LOWER(phone) LIKE $2
			OR LOWER(company) LIKE $2
		)
		ORDER BY created_at DESC
		LIMIT 10;`,
		ownerID,
		search,
	)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	contacts := make([]dto.ContactResponse, 0)

	for rows.Next() {

		var contact dto.ContactResponse

		err := rows.Scan(
			&contact.ID,
			&contact.FirstName,
			&contact.LastName,
			&contact.Email,
			&contact.Phone,
			&contact.Company,
			&contact.JobTitle,
			&contact.Notes,
			&contact.OwnerID,
			&contact.CreatedAt,
			&contact.UpdatedAt,
		)

		if err != nil {

			return nil, err

		}

		contacts = append(
			contacts,
			contact,
		)

	}

	return contacts, nil

}
