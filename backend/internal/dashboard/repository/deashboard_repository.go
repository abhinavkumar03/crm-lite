package repository

import (
	"context"
	"encoding/json"
	"time"

	"github.com/abhinavkumar03/crm-lite/backend/internal/dashboard/dto"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository struct {
	db *pgxpool.Pool
}

func New(db *pgxpool.Pool) *Repository {
	return &Repository{db: db}
}

func (r *Repository) GetDashboard(
	ctx context.Context,
	orgID string,
) (*dto.DashboardResponse, error) {
	out := &dto.DashboardResponse{
		ModuleCounts:  make([]dto.ModuleCount, 0),
		RecentRecords: make([]dto.RecentRecord, 0),
	}

	err := r.db.QueryRow(ctx, `
		SELECT
			(SELECT COUNT(*) FROM modules
			  WHERE organization_id = $1
			    AND is_enabled = TRUE
			    AND storage_strategy = 'dynamic'),
			(SELECT COUNT(*) FROM records r
			  JOIN modules m ON m.id = r.module_id
			  WHERE r.organization_id = $1
			    AND m.storage_strategy = 'dynamic'
			    AND m.is_enabled = TRUE)
	`, orgID).Scan(&out.TotalModules, &out.TotalRecords)
	if err != nil {
		return nil, err
	}

	rows, err := r.db.Query(ctx, `
		SELECT m.id, m.api_name, m.plural_label, COALESCE(m.icon, ''), COALESCE(m.color, ''),
		       COUNT(r.id) AS record_count
		FROM modules m
		LEFT JOIN records r
		  ON r.module_id = m.id AND r.organization_id = m.organization_id
		WHERE m.organization_id = $1
		  AND m.is_enabled = TRUE
		  AND m.storage_strategy = 'dynamic'
		GROUP BY m.id, m.api_name, m.plural_label, m.icon, m.color, m.sort_order
		ORDER BY m.sort_order ASC, m.plural_label ASC
	`, orgID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var mc dto.ModuleCount
		if err := rows.Scan(
			&mc.ModuleID,
			&mc.APIName,
			&mc.PluralLabel,
			&mc.Icon,
			&mc.Color,
			&mc.RecordCount,
		); err != nil {
			return nil, err
		}
		out.ModuleCounts = append(out.ModuleCounts, mc)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	recent, err := r.db.Query(ctx, `
		SELECT r.id, r.module_id, m.plural_label, m.api_name, r.data, r.created_at
		FROM records r
		JOIN modules m ON m.id = r.module_id
		WHERE r.organization_id = $1
		  AND m.storage_strategy = 'dynamic'
		  AND m.is_enabled = TRUE
		ORDER BY r.created_at DESC
		LIMIT 8
	`, orgID)
	if err != nil {
		return nil, err
	}
	defer recent.Close()

	for recent.Next() {
		var (
			rec       dto.RecentRecord
			raw       []byte
			createdAt time.Time
		)
		if err := recent.Scan(
			&rec.ID,
			&rec.ModuleID,
			&rec.ModuleLabel,
			&rec.APIName,
			&raw,
			&createdAt,
		); err != nil {
			return nil, err
		}
		rec.Title = recordTitle(raw)
		rec.CreatedAt = createdAt.Format(time.RFC3339)
		out.RecentRecords = append(out.RecentRecords, rec)
	}
	if err := recent.Err(); err != nil {
		return nil, err
	}

	return out, nil
}

func recordTitle(raw []byte) string {
	var data map[string]any
	if err := json.Unmarshal(raw, &data); err != nil {
		return "Untitled"
	}
	for _, key := range []string{"name", "title", "email", "company"} {
		if v, ok := data[key]; ok {
			if s, ok := v.(string); ok && s != "" {
				return s
			}
		}
	}
	for _, v := range data {
		if s, ok := v.(string); ok && s != "" {
			return s
		}
	}
	return "Untitled"
}
