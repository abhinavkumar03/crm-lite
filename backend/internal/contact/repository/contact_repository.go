package repository

import (
	"context"
	"errors"
	"fmt"
	"math"
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
) (*dto.ListContactsResponse, error) {

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
LOWER(first_name) LIKE $%d
OR LOWER(last_name) LIKE $%d
OR LOWER(email) LIKE $%d
OR LOWER(phone) LIKE $%d
OR LOWER(company) LIKE $%d
)`,
				argIndex,
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

	whereClause := strings.Join(
		conditions,
		" AND ",
	)

	countQuery := fmt.Sprintf(`
SELECT COUNT(*)
FROM contacts
WHERE %s
`, whereClause)

	var total int64

	err := r.db.QueryRow(
		ctx,
		countQuery,
		args...,
	).Scan(&total)

	if err != nil {
		return nil, err
	}

	allowedSorts := map[string]string{
		"first_name": "first_name",
		"last_name":  "last_name",
		"company":    "company",
		"created_at": "created_at",
	}

	sortBy := "created_at"

	if v, ok := allowedSorts[req.SortBy]; ok {
		sortBy = v
	}

	order := "DESC"

	if strings.ToUpper(req.SortOrder) == "ASC" {
		order = "ASC"
	}

	offset := (req.Page - 1) * req.Limit

	query := fmt.Sprintf(`
SELECT
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

	if err := rows.Err(); err != nil {
		return nil, err
	}

	totalPages := 0

	if total > 0 {
		totalPages = int(math.Ceil(
			float64(total) / float64(req.Limit),
		))
	}

	return &dto.ListContactsResponse{
		Data:       contacts,
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

func (r *Repository) ExistsByIDAndOwner(
	ctx context.Context,
	id string,
	ownerID string,
) (bool, error) {

	query := `
	SELECT EXISTS(
		SELECT 1
		FROM contacts
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
