package repository

import (
	"context"
	"time"

	"github.com/abhinavkumar03/crm-lite/backend/internal/dashboard/dto"
	leadDto "github.com/abhinavkumar03/crm-lite/backend/internal/lead/dto"
	taskDto "github.com/abhinavkumar03/crm-lite/backend/internal/task/dto"
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

func (r *Repository) GetMetrics(
	ctx context.Context,
	ownerID string,
) (*dto.DashboardResponse, error) {

	var dashboard dto.DashboardResponse

	err := r.db.QueryRow(
		ctx,
		`
SELECT

(SELECT COUNT(*) FROM leads WHERE owner_id=$1),

(SELECT COUNT(*) FROM leads
WHERE owner_id=$1
AND status='NEW'),

(SELECT COUNT(*) FROM leads
WHERE owner_id=$1
AND status='CONTACTED'),

(SELECT COUNT(*) FROM leads
WHERE owner_id=$1
AND status='QUALIFIED'),

(SELECT COUNT(*) FROM leads
WHERE owner_id=$1
AND status='WON'),

(SELECT COUNT(*) FROM leads
WHERE owner_id=$1
AND status='LOST'),

(SELECT COUNT(*) FROM contacts
WHERE owner_id=$1),

(SELECT COUNT(*) FROM tasks
WHERE owner_id=$1),

(SELECT COUNT(*) FROM tasks
WHERE owner_id=$1
AND status='PENDING'),

(SELECT COUNT(*) FROM tasks
WHERE owner_id=$1
AND status='IN_PROGRESS'),

(SELECT COUNT(*) FROM tasks
WHERE owner_id=$1
AND status='COMPLETED')
`,
		ownerID,
	).Scan(
		&dashboard.TotalLeads,
		&dashboard.NewLeads,
		&dashboard.ContactedLeads,
		&dashboard.QualifiedLeads,
		&dashboard.WonLeads,
		&dashboard.LostLeads,
		&dashboard.TotalContacts,
		&dashboard.TotalTasks,
		&dashboard.PendingTasks,
		&dashboard.InProgressTasks,
		&dashboard.CompletedTasks,
	)

	if err != nil {
		return nil, err
	}

	return &dashboard, nil
}

func (r *Repository) RecentLeads(
	ctx context.Context,
	ownerID string,
) ([]leadDto.LeadResponse, error) {

	query := `
	SELECT
		id,
		name,
		email,
		phone,
		company,
		status,
		notes
	FROM leads
	WHERE owner_id = $1
	ORDER BY created_at DESC
	LIMIT 5;
	`

	rows, err := r.db.Query(ctx, query, ownerID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var leads []leadDto.LeadResponse

	for rows.Next() {

		var lead leadDto.LeadResponse

		if err := rows.Scan(
			&lead.ID,
			&lead.Name,
			&lead.Email,
			&lead.Phone,
			&lead.Company,
			&lead.Status,
			&lead.Notes,
		); err != nil {
			return nil, err
		}

		leads = append(leads, lead)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return leads, nil
}

func (r *Repository) UpcomingTasks(
	ctx context.Context,
	ownerID string,
) ([]taskDto.TaskResponse, error) {

	query := `
	SELECT
		id,
		title,
		description,
		status,
		lead_id,
		contact_id,
		due_date
	FROM tasks
	WHERE owner_id = $1
	  AND status IN ('PENDING', 'IN_PROGRESS')
	ORDER BY due_date ASC NULLS LAST
	LIMIT 5;
	`

	rows, err := r.db.Query(ctx, query, ownerID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tasks []taskDto.TaskResponse

	for rows.Next() {

		var task taskDto.TaskResponse
		var dueDate *time.Time

		if err := rows.Scan(
			&task.ID,
			&task.Title,
			&task.Description,
			&task.Status,
			&task.LeadID,
			&task.ContactID,
			&dueDate,
		); err != nil {
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

	return tasks, nil
}
