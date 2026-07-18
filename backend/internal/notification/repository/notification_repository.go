package repository

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/abhinavkumar03/crm-lite/backend/internal/notification/dto"
	"github.com/abhinavkumar03/crm-lite/backend/internal/notification/entity"
)

const notificationColumns = `
	id, organization_id, channel, recipient, subject, body, template, data,
	status, provider, error, entity_type, entity_id, created_by, sent_at,
	created_at, updated_at
`

type Repository struct {
	db *pgxpool.Pool
}

func New(db *pgxpool.Pool) *Repository {
	return &Repository{db: db}
}

func scan(row pgx.Row, n *entity.Notification) error {
	return row.Scan(
		&n.ID, &n.OrganizationID, &n.Channel, &n.Recipient, &n.Subject, &n.Body,
		&n.Template, &n.Data, &n.Status, &n.Provider, &n.Error, &n.EntityType,
		&n.EntityID, &n.CreatedBy, &n.SentAt, &n.CreatedAt, &n.UpdatedAt,
	)
}

func (r *Repository) Create(ctx context.Context, n *entity.Notification) error {
	return r.db.QueryRow(ctx, `
		INSERT INTO notifications (
			organization_id, channel, recipient, subject, body, template, data,
			status, entity_type, entity_id, created_by
		)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11)
		RETURNING id, status, created_at, updated_at
	`,
		n.OrganizationID, n.Channel, n.Recipient, n.Subject, n.Body, n.Template,
		n.Data, n.Status, n.EntityType, n.EntityID, n.CreatedBy,
	).Scan(&n.ID, &n.Status, &n.CreatedAt, &n.UpdatedAt)
}

func (r *Repository) GetByID(ctx context.Context, orgID, id string) (*entity.Notification, error) {
	var n entity.Notification
	err := scan(r.db.QueryRow(ctx, `
		SELECT `+notificationColumns+`
		FROM notifications
		WHERE id = $1 AND organization_id = $2
	`, id, orgID), &n)

	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &n, nil
}

func (r *Repository) List(ctx context.Context, orgID string, q dto.ListQuery) ([]entity.Notification, int, error) {
	conds := []string{"organization_id = $1"}
	args := []any{orgID}
	next := 2

	if q.Status != "" {
		conds = append(conds, fmt.Sprintf("status = $%d", next))
		args = append(args, q.Status)
		next++
	}
	if q.Channel != "" {
		conds = append(conds, fmt.Sprintf("channel = $%d", next))
		args = append(args, q.Channel)
		next++
	}

	where := "WHERE " + strings.Join(conds, " AND ")

	var total int
	if err := r.db.QueryRow(ctx, "SELECT COUNT(*) FROM notifications "+where, args...).Scan(&total); err != nil {
		return nil, 0, err
	}

	limitPH := fmt.Sprintf("$%d", next)
	offsetPH := fmt.Sprintf("$%d", next+1)
	args = append(args, q.PageSize, (q.Page-1)*q.PageSize)

	rows, err := r.db.Query(ctx,
		"SELECT "+notificationColumns+" FROM notifications "+where+
			" ORDER BY created_at DESC LIMIT "+limitPH+" OFFSET "+offsetPH, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	items := make([]entity.Notification, 0)
	for rows.Next() {
		var n entity.Notification
		if err := scan(rows, &n); err != nil {
			return nil, 0, err
		}
		items = append(items, n)
	}
	return items, total, rows.Err()
}

// MarkSent transitions a notification to the sent state and stamps the provider.
func (r *Repository) MarkSent(ctx context.Context, id, provider string) error {
	_, err := r.db.Exec(ctx, `
		UPDATE notifications
		SET status = 'sent', provider = $2, error = NULL, sent_at = NOW(), updated_at = NOW()
		WHERE id = $1
	`, id, provider)
	return err
}

// MarkFailed records a delivery failure and its reason.
func (r *Repository) MarkFailed(ctx context.Context, id, reason string) error {
	_, err := r.db.Exec(ctx, `
		UPDATE notifications
		SET status = 'failed', error = $2, updated_at = NOW()
		WHERE id = $1
	`, id, reason)
	return err
}
