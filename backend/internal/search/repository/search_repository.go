package repository

import (
	"context"
	"encoding/json"
	"strings"

	"github.com/abhinavkumar03/crm-lite/backend/internal/search/dto"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository struct {
	db *pgxpool.Pool
}

func New(db *pgxpool.Pool) *Repository {
	return &Repository{db: db}
}

// Search scans dynamic module records whose JSONB data contains the query
// (case-insensitive). Limited to a small global result set for the topbar.
func (r *Repository) Search(
	ctx context.Context,
	orgID, query string,
	limit int,
) ([]dto.SearchHit, error) {
	if limit <= 0 {
		limit = 15
	}
	pattern := "%" + strings.ToLower(strings.TrimSpace(query)) + "%"

	rows, err := r.db.Query(ctx, `
		SELECT r.id, r.module_id, m.plural_label, m.api_name, r.data
		FROM records r
		JOIN modules m ON m.id = r.module_id
		WHERE r.organization_id = $1
		  AND m.storage_strategy = 'dynamic'
		  AND m.is_enabled = TRUE
		  AND LOWER(r.data::text) LIKE $2
		ORDER BY r.updated_at DESC
		LIMIT $3
	`, orgID, pattern, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	hits := make([]dto.SearchHit, 0)
	for rows.Next() {
		var (
			hit dto.SearchHit
			raw []byte
		)
		if err := rows.Scan(
			&hit.ID,
			&hit.ModuleID,
			&hit.ModuleLabel,
			&hit.APIName,
			&raw,
		); err != nil {
			return nil, err
		}
		hit.Title, hit.Subtitle = titlesFromData(raw)
		hits = append(hits, hit)
	}
	return hits, rows.Err()
}

func titlesFromData(raw []byte) (title, subtitle string) {
	var data map[string]any
	if err := json.Unmarshal(raw, &data); err != nil {
		return "Untitled", ""
	}
	pick := func(keys ...string) string {
		for _, k := range keys {
			if v, ok := data[k]; ok {
				if s, ok := v.(string); ok && s != "" {
					return s
				}
			}
		}
		return ""
	}
	title = pick("name", "title", "email")
	if title == "" {
		for _, v := range data {
			if s, ok := v.(string); ok && s != "" {
				title = s
				break
			}
		}
	}
	if title == "" {
		title = "Untitled"
	}
	subtitle = pick("email", "company", "city", "industry", "stage")
	if subtitle == title {
		subtitle = pick("company", "city", "industry", "stage")
	}
	return title, subtitle
}
