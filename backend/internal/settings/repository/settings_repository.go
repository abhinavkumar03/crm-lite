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

const orgCols = `
	id, name, slug, plan,
	logo_url, industry, company_size, country, status, created_by,
	settings, updated_at
`

func scanOrg(row pgx.Row) (*entity.Organization, error) {
	var o entity.Organization
	err := row.Scan(
		&o.ID, &o.Name, &o.Slug, &o.Plan,
		&o.LogoURL, &o.Industry, &o.CompanySize, &o.Country, &o.Status, &o.CreatedBy,
		&o.Settings, &o.UpdatedAt,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &o, nil
}

func (r *Repository) GetByID(ctx context.Context, orgID string) (*entity.Organization, error) {
	return scanOrg(r.db.QueryRow(ctx, `
		SELECT `+orgCols+` FROM organizations WHERE id = $1
	`, orgID))
}

type ProfileUpdate struct {
	Name        string
	LogoURL     *string
	Industry    *string
	CompanySize *string
	Country     *string
	Status      string
	Settings    []byte
}

func (r *Repository) Update(ctx context.Context, orgID string, p ProfileUpdate) (*entity.Organization, error) {
	return scanOrg(r.db.QueryRow(ctx, `
		UPDATE organizations
		SET name = $2,
		    logo_url = $3,
		    industry = $4,
		    company_size = $5,
		    country = $6,
		    status = $7,
		    settings = $8,
		    updated_at = NOW()
		WHERE id = $1
		RETURNING `+orgCols+`
	`, orgID, p.Name, p.LogoURL, p.Industry, p.CompanySize, p.Country, p.Status, p.Settings))
}
