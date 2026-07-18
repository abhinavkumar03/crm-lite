package seeders

import (
	"context"
	"encoding/json"
	"fmt"
	"math/rand"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

// BusinessDataSeeder generates a realistic demo dataset owned by the admin user:
// dynamic company/deal records in the JSONB record engine.
//
// It is idempotent: if the org already has company records, it does nothing.
type BusinessDataSeeder struct{}

func (BusinessDataSeeder) Name() string { return "business_data" }

func (BusinessDataSeeder) Run(ctx context.Context, db *pgxpool.Pool) error {
	orgID, err := getOrgID(ctx, db)
	if err != nil {
		return err
	}

	ownerIDs, err := demoOwnerIDs(ctx, db)
	if err != nil {
		return err
	}

	companyModuleID, err := getModuleID(ctx, db, orgID, "company")
	if err != nil {
		return err
	}
	dealModuleID, err := getModuleID(ctx, db, orgID, "deal")
	if err != nil {
		return err
	}

	// Idempotency guard — dynamic records only.
	var existing int
	if err := db.QueryRow(ctx,
		`SELECT count(*) FROM records WHERE organization_id = $1 AND module_id = $2`,
		orgID, companyModuleID,
	).Scan(&existing); err != nil {
		return err
	}
	if existing > 0 {
		return nil
	}

	r := rand.New(rand.NewSource(42))
	now := time.Now()

	companyRecordIDs, err := seedCompanies(ctx, db, r, now, orgID, companyModuleID, ownerIDs, 20)
	if err != nil {
		return err
	}
	if err := seedDeals(ctx, db, r, now, orgID, dealModuleID, ownerIDs, companyRecordIDs, 15); err != nil {
		return err
	}

	return nil
}

func demoOwnerIDs(ctx context.Context, db *pgxpool.Pool) ([]string, error) {
	emails := []string{
		"admin@crmlite.com",
		"priya@crmlite.com",
		"vikram@crmlite.com",
		"sneha@crmlite.com",
	}
	ids := make([]string, 0, len(emails))
	for _, email := range emails {
		var id string
		if err := db.QueryRow(ctx, `SELECT id FROM users WHERE email = $1`, email).Scan(&id); err != nil {
			return nil, fmt.Errorf("lookup user %q: %w", email, err)
		}
		ids = append(ids, id)
	}
	return ids, nil
}

func spread(r *rand.Rand, now time.Time) time.Time {
	return now.AddDate(0, 0, -r.Intn(180)).Add(-time.Duration(r.Intn(24)) * time.Hour)
}

func seedCompanies(ctx context.Context, db *pgxpool.Pool, r *rand.Rand, now time.Time, orgID, moduleID string, owners []string, n int) ([]string, error) {
	ids := make([]string, 0, n)
	for i := 0; i < n; i++ {
		name := companyNames[i%len(companyNames)]
		owner := owners[i%len(owners)]
		data := map[string]any{
			"name":      name,
			"industry":  pick(r, industries),
			"city":      pick(r, cities),
			"website":   website(name),
			"employees": (r.Intn(50) + 1) * 20,
			"tags":      pickTags(r),
		}
		vis := "organization"
		switch i % 5 {
		case 0:
			vis = "owner"
		case 1:
			vis = "hierarchy"
		}
		id, err := insertRecord(ctx, db, orgID, moduleID, owner, data, spread(r, now), vis)
		if err != nil {
			return nil, err
		}
		ids = append(ids, id)
	}
	return ids, nil
}

func seedDeals(ctx context.Context, db *pgxpool.Pool, r *rand.Rand, now time.Time, orgID, moduleID string, owners []string, companyIDs []string, n int) error {
	for i := 0; i < n; i++ {
		var companyRef any
		if len(companyIDs) > 0 {
			companyRef = companyIDs[r.Intn(len(companyIDs))]
		}
		owner := owners[(i+1)%len(owners)]
		data := map[string]any{
			"title":      fmt.Sprintf("Deal #%d", 1000+i),
			"amount":     (r.Intn(90) + 10) * 10000,
			"stage":      pick(r, dealStages),
			"close_date": now.AddDate(0, 0, r.Intn(90)).Format("2006-01-02"),
			"company":    companyRef,
		}
		vis := "organization"
		if i%4 == 0 {
			vis = "hierarchy"
		} else if i%5 == 0 {
			vis = "owner"
		}
		if _, err := insertRecord(ctx, db, orgID, moduleID, owner, data, spread(r, now), vis); err != nil {
			return err
		}
	}
	return nil
}

func insertRecord(ctx context.Context, db *pgxpool.Pool, orgID, moduleID, owner string, data map[string]any, createdAt time.Time, visibility string) (string, error) {
	payload, err := json.Marshal(data)
	if err != nil {
		return "", err
	}
	if visibility == "" {
		visibility = "organization"
	}
	var id string
	err = db.QueryRow(ctx, `
		INSERT INTO records (
			organization_id, module_id, data, owner_id, created_by, updated_by,
			visibility, created_at, updated_at
		)
		VALUES ($1, $2, $3, $4, $4, $4, $5, $6, $6)
		RETURNING id
	`, orgID, moduleID, string(payload), owner, visibility, createdAt).Scan(&id)
	if err != nil {
		return "", fmt.Errorf("insert record: %w", err)
	}
	return id, nil
}
