package repository

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/abhinavkumar03/crm-lite/backend/internal/settings/entity"
)

type Repository struct {
	db *pgxpool.Pool
}

func New(db *pgxpool.Pool) *Repository {
	return &Repository{db: db}
}

// GetByID returns the organization row, or (nil, nil) if it does not exist.
func (r *Repository) GetByID(ctx context.Context, orgID string) (*entity.Organization, error) {
	var o entity.Organization
	err := r.db.QueryRow(ctx, `
		SELECT id, name, slug, plan, settings, updated_at
		FROM organizations
		WHERE id = $1
	`, orgID).Scan(&o.ID, &o.Name, &o.Slug, &o.Plan, &o.Settings, &o.UpdatedAt)

	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &o, nil
}

// Update persists the org name and the settings JSONB and returns the fresh row.
func (r *Repository) Update(ctx context.Context, orgID, name string, settings []byte) (*entity.Organization, error) {
	var o entity.Organization
	err := r.db.QueryRow(ctx, `
		UPDATE organizations
		SET name = $2, settings = $3, updated_at = NOW()
		WHERE id = $1
		RETURNING id, name, slug, plan, settings, updated_at
	`, orgID, name, settings).Scan(&o.ID, &o.Name, &o.Slug, &o.Plan, &o.Settings, &o.UpdatedAt)

	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &o, nil
}
