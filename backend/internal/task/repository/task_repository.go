package repository

import (
	"context"
	"errors"
	"fmt"
	"math"
	"strings"
	"time"

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
) (*dto.ListTasksResponse, error) {

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
LOWER(title) LIKE $%d
OR LOWER(description) LIKE $%d
)`,
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

	countQuery := fmt.Sprintf(`
SELECT COUNT(*)
FROM tasks
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
		"title":      "title",
		"status":     "status",
		"due_date":   "due_date",
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
	title,
	description,
	status,
	lead_id,
	contact_id,
	due_date,
	owner_id,
	created_at,
	updated_at
FROM tasks
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

	tasks := make([]dto.TaskResponse, 0)

	for rows.Next() {

		var (
			task    dto.TaskResponse
			dueDate *time.Time
		)

		err := rows.Scan(
			&task.ID,
			&task.Title,
			&task.Description,
			&task.Status,
			&task.LeadID,
			&task.ContactID,
			&dueDate,
			&task.OwnerID,
			&task.CreatedAt,
			&task.UpdatedAt,
		)

		if err != nil {
			return nil, err
		}

		if dueDate != nil {
			v := dueDate.Format(time.RFC3339)
			task.DueDate = &v
		}

		tasks = append(tasks, task)
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

	return &dto.ListTasksResponse{
		Data:       tasks,
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

func (r *Repository) Delete(
	ctx context.Context,
	id string,
	ownerID string,
) (bool, error) {

	query := `
	DELETE FROM tasks
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
