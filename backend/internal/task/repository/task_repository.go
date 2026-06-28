package repository

import (
	"context"
	"errors"

	"github.com/abhinavkumar03/crm-lite/backend/internal/task/dto"
	"github.com/abhinavkumar03/crm-lite/backend/internal/task/entity"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository struct {
	db *pgxpool.Pool
}

func New(db *pgxpool.Pool) *Repository {
	return &Repository{db: db}
}

func (r *Repository) Create(
	ctx context.Context,
	task *entity.Task,
) error {

	query := `
	INSERT INTO tasks
	(
		owner_id,
		lead_id,
		contact_id,
		title,
		description,
		status,
		due_date
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
		task.OwnerID,
		task.LeadID,
		task.ContactID,
		task.Title,
		task.Description,
		task.Status,
		task.DueDate,
	).Scan(
		&task.ID,
		&task.CreatedAt,
		&task.UpdatedAt,
	)
}

func (r *Repository) List(
	ctx context.Context,
	ownerID string,
	req dto.ListTasksRequest,
) ([]entity.Task, error) {

	offset := (req.Page - 1) * req.Limit

	query := `
	SELECT
		id,
		owner_id,
		lead_id,
		contact_id,
		title,
		description,
		status,
		due_date,
		created_at,
		updated_at
	FROM tasks
	WHERE owner_id = $1
	AND (
		$2 = ''
		OR title ILIKE '%' || $2 || '%'
	)
	AND (
		$3 = ''
		OR status = $3
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

	tasks := make([]entity.Task, 0)

	for rows.Next() {

		var task entity.Task

		err := rows.Scan(
			&task.ID,
			&task.OwnerID,
			&task.LeadID,
			&task.ContactID,
			&task.Title,
			&task.Description,
			&task.Status,
			&task.DueDate,
			&task.CreatedAt,
			&task.UpdatedAt,
		)

		if err != nil {
			return nil, err
		}

		tasks = append(tasks, task)
	}

	return tasks, rows.Err()
}

func (r *Repository) GetByID(
	ctx context.Context,
	id string,
	ownerID string,
) (*entity.Task, error) {

	query := `
	SELECT
		id,
		owner_id,
		lead_id,
		contact_id,
		title,
		description,
		status,
		due_date,
		created_at,
		updated_at
	FROM tasks
	WHERE id = $1
	AND owner_id = $2;
	`

	var task entity.Task

	err := r.db.QueryRow(
		ctx,
		query,
		id,
		ownerID,
	).Scan(
		&task.ID,
		&task.OwnerID,
		&task.LeadID,
		&task.ContactID,
		&task.Title,
		&task.Description,
		&task.Status,
		&task.DueDate,
		&task.CreatedAt,
		&task.UpdatedAt,
	)

	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	return &task, nil
}

func (r *Repository) Update(
	ctx context.Context,
	task *entity.Task,
) error {

	query := `
	UPDATE tasks
	SET
		lead_id = $1,
		contact_id = $2,
		title = $3,
		description = $4,
		status = $5,
		due_date = $6,
		updated_at = NOW()
	WHERE id = $7
	AND owner_id = $8
	RETURNING updated_at;
	`

	return r.db.QueryRow(
		ctx,
		query,
		task.LeadID,
		task.ContactID,
		task.Title,
		task.Description,
		task.Status,
		task.DueDate,
		task.ID,
		task.OwnerID,
	).Scan(
		&task.UpdatedAt,
	)
}
