package repository

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/abhinavkumar03/crm-lite/backend/internal/tour/entity"
)

type Repository struct {
	db *pgxpool.Pool
}

func New(db *pgxpool.Pool) *Repository {
	return &Repository{db: db}
}

func scan(row pgx.Row, p *entity.TourProgress) error {
	var steps []byte
	if err := row.Scan(
		&p.ID, &p.OrganizationID, &p.UserID, &p.TourKey, &p.Status,
		&p.CurrentStep, &steps, &p.StartedAt, &p.CompletedAt,
		&p.CreatedAt, &p.UpdatedAt,
	); err != nil {
		return err
	}
	p.CompletedSteps = decodeSteps(steps)
	return nil
}

// GetByUser returns the user's progress for a tour, or (nil, nil) if the user
// has never started it.
func (r *Repository) GetByUser(ctx context.Context, orgID, userID, tourKey string) (*entity.TourProgress, error) {
	var p entity.TourProgress
	err := scan(r.db.QueryRow(ctx, `
		SELECT id, organization_id, user_id, tour_key, status, current_step,
		       completed_steps, started_at, completed_at, created_at, updated_at
		FROM tour_progress
		WHERE organization_id = $1 AND user_id = $2 AND tour_key = $3
	`, orgID, userID, tourKey), &p)

	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &p, nil
}

// Upsert creates or updates the user's progress for a tour. The unique
// (organization_id, user_id, tour_key) constraint makes this idempotent.
func (r *Repository) Upsert(ctx context.Context, p *entity.TourProgress) error {
	steps, err := json.Marshal(nonNil(p.CompletedSteps))
	if err != nil {
		return err
	}

	return scan(r.db.QueryRow(ctx, `
		INSERT INTO tour_progress (
			organization_id, user_id, tour_key, status, current_step,
			completed_steps, started_at, completed_at
		)
		VALUES ($1, $2, $3, $4, $5, $6, COALESCE($7, NOW()), $8)
		ON CONFLICT (organization_id, user_id, tour_key) DO UPDATE SET
			status          = EXCLUDED.status,
			current_step    = EXCLUDED.current_step,
			completed_steps = EXCLUDED.completed_steps,
			completed_at    = EXCLUDED.completed_at,
			updated_at      = NOW()
		RETURNING id, organization_id, user_id, tour_key, status, current_step,
		          completed_steps, started_at, completed_at, created_at, updated_at
	`,
		p.OrganizationID, p.UserID, p.TourKey, p.Status, p.CurrentStep,
		steps, nullableTime(p.StartedAt), p.CompletedAt,
	), p)
}

// Restart resets progress to the beginning of the tour and returns the fresh
// record.
func (r *Repository) Restart(ctx context.Context, orgID, userID, tourKey string) (*entity.TourProgress, error) {
	var p entity.TourProgress
	err := scan(r.db.QueryRow(ctx, `
		INSERT INTO tour_progress (
			organization_id, user_id, tour_key, status, current_step,
			completed_steps, started_at, completed_at
		)
		VALUES ($1, $2, $3, 'active', 0, '[]'::jsonb, NOW(), NULL)
		ON CONFLICT (organization_id, user_id, tour_key) DO UPDATE SET
			status          = 'active',
			current_step    = 0,
			completed_steps = '[]'::jsonb,
			started_at      = NOW(),
			completed_at    = NULL,
			updated_at      = NOW()
		RETURNING id, organization_id, user_id, tour_key, status, current_step,
		          completed_steps, started_at, completed_at, created_at, updated_at
	`, orgID, userID, tourKey), &p)
	if err != nil {
		return nil, err
	}
	return &p, nil
}

func decodeSteps(raw []byte) []string {
	steps := []string{}
	if len(raw) > 0 {
		_ = json.Unmarshal(raw, &steps)
	}
	return steps
}

func nonNil(s []string) []string {
	if s == nil {
		return []string{}
	}
	return s
}

// nullableTime lets the caller omit started_at (zero time) so the INSERT default
// (NOW()) is used and existing rows keep their original start time.
func nullableTime(t time.Time) interface{} {
	if t.IsZero() {
		return nil
	}
	return t
}
