package repository

import (
	"context"

	"github.com/abhinavkumar03/crm-lite/backend/internal/task/dto"
	"github.com/abhinavkumar03/crm-lite/backend/internal/task/entity"
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
