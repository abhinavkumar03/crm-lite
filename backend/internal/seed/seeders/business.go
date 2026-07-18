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
// native leads/contacts/tasks/activities/notes (so the existing dashboard has
// data) plus dynamic company/deal records stored in the JSONB record engine
// (so the metadata engine is demonstrated end-to-end).
//
// It is idempotent: if the admin already has leads, it does nothing, so a
// -fresh re-run will not create duplicates.
type BusinessDataSeeder struct{}

func (BusinessDataSeeder) Name() string { return "business_data" }

func (BusinessDataSeeder) Run(ctx context.Context, db *pgxpool.Pool) error {
	orgID, err := getOrgID(ctx, db)
	if err != nil {
		return err
	}

	var adminID string
	if err := db.QueryRow(ctx,
		`SELECT id FROM users WHERE email = $1`, "admin@crmlite.com",
	).Scan(&adminID); err != nil {
		return fmt.Errorf("lookup admin user: %w", err)
	}

	// Idempotency guard.
	var existing int
	if err := db.QueryRow(ctx,
		`SELECT count(*) FROM leads WHERE owner_id = $1`, adminID,
	).Scan(&existing); err != nil {
		return err
	}
	if existing > 0 {
		return nil // demo data already present
	}

	r := rand.New(rand.NewSource(42))
	now := time.Now()

	leadIDs, err := seedLeads(ctx, db, r, now, adminID, 50)
	if err != nil {
		return err
	}
	contactIDs, err := seedContacts(ctx, db, r, now, adminID, 30)
	if err != nil {
		return err
	}
	if err := seedTasks(ctx, db, r, now, adminID, leadIDs, contactIDs, 60); err != nil {
		return err
	}
	if err := seedActivities(ctx, db, r, now, adminID, leadIDs); err != nil {
		return err
	}
	if err := seedNotes(ctx, db, r, now, adminID, leadIDs, 15); err != nil {
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
	companyRecordIDs, err := seedCompanies(ctx, db, r, now, orgID, companyModuleID, adminID, 20)
	if err != nil {
		return err
	}
	if err := seedDeals(ctx, db, r, now, orgID, dealModuleID, adminID, companyRecordIDs, 15); err != nil {
		return err
	}

	return nil
}

func spread(r *rand.Rand, now time.Time) time.Time {
	return now.AddDate(0, 0, -r.Intn(180)).Add(-time.Duration(r.Intn(24)) * time.Hour)
}

func seedLeads(ctx context.Context, db *pgxpool.Pool, r *rand.Rand, now time.Time, owner string, n int) ([]string, error) {
	ids := make([]string, 0, n)
	for i := 0; i < n; i++ {
		first, last := fullName(r)
		name := first + " " + last
		company := pick(r, companyNames)
		createdAt := spread(r, now)

		var id string
		err := db.QueryRow(ctx, `
			INSERT INTO leads (owner_id, name, email, phone, company, status, notes, created_at, updated_at)
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $8)
			RETURNING id
		`, owner, name, emailFrom(first, last, r), phone(r), company,
			pick(r, leadStatuses), "Imported from "+pick(r, cities)+" campaign", createdAt).Scan(&id)
		if err != nil {
			return nil, fmt.Errorf("insert lead: %w", err)
		}
		ids = append(ids, id)
	}
	return ids, nil
}

func seedContacts(ctx context.Context, db *pgxpool.Pool, r *rand.Rand, now time.Time, owner string, n int) ([]string, error) {
	ids := make([]string, 0, n)
	for i := 0; i < n; i++ {
		first, last := fullName(r)
		createdAt := spread(r, now)

		var id string
		err := db.QueryRow(ctx, `
			INSERT INTO contacts (owner_id, first_name, last_name, email, phone, company, job_title, notes, created_at, updated_at)
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $9)
			RETURNING id
		`, owner, first, last, emailFrom(first, last, r), phone(r),
			pick(r, companyNames), pick(r, jobTitles), "Key contact based in "+pick(r, cities), createdAt).Scan(&id)
		if err != nil {
			return nil, fmt.Errorf("insert contact: %w", err)
		}
		ids = append(ids, id)
	}
	return ids, nil
}

func seedTasks(ctx context.Context, db *pgxpool.Pool, r *rand.Rand, now time.Time, owner string, leadIDs, contactIDs []string, n int) error {
	titles := []string{
		"Follow up call", "Send proposal", "Schedule demo", "Share pricing",
		"Onboarding kickoff", "Quarterly review", "Contract renewal", "Collect feedback",
	}
	for i := 0; i < n; i++ {
		createdAt := spread(r, now)
		dueDate := createdAt.AddDate(0, 0, r.Intn(30)-10) // some overdue, some upcoming

		var leadID, contactID any
		if r.Intn(2) == 0 && len(leadIDs) > 0 {
			leadID = leadIDs[r.Intn(len(leadIDs))]
		}
		if r.Intn(2) == 0 && len(contactIDs) > 0 {
			contactID = contactIDs[r.Intn(len(contactIDs))]
		}

		_, err := db.Exec(ctx, `
			INSERT INTO tasks (owner_id, lead_id, contact_id, title, description, status, due_date, created_at, updated_at)
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $8)
		`, owner, leadID, contactID, pick(r, titles), "Auto-generated demo task",
			pick(r, taskStatuses), dueDate, createdAt)
		if err != nil {
			return fmt.Errorf("insert task: %w", err)
		}
	}
	return nil
}

func seedActivities(ctx context.Context, db *pgxpool.Pool, r *rand.Rand, now time.Time, owner string, leadIDs []string) error {
	actions := []struct {
		Action string
		Desc   string
	}{
		{"LEAD_CREATED", "Lead created"},
		{"LEAD_CONTACTED", "Reached out to lead"},
		{"LEAD_STATUS_CHANGED", "Lead status changed"},
		{"NOTE_ADDED", "Note added to lead"},
	}
	for _, leadID := range leadIDs {
		count := 1 + r.Intn(3)
		for j := 0; j < count; j++ {
			a := actions[r.Intn(len(actions))]
			_, err := db.Exec(ctx, `
				INSERT INTO activities (id, entity_type, entity_id, action, description, performed_by, created_at)
				VALUES (gen_random_uuid(), 'LEAD', $1, $2, $3, $4, $5)
			`, leadID, a.Action, a.Desc, owner, spread(r, now))
			if err != nil {
				return fmt.Errorf("insert activity: %w", err)
			}
		}
	}
	return nil
}

func seedNotes(ctx context.Context, db *pgxpool.Pool, r *rand.Rand, now time.Time, owner string, leadIDs []string, n int) error {
	texts := []string{
		"Interested in enterprise plan.", "Requested a callback next week.",
		"Budget approved by finance.", "Comparing with competitor.",
		"Wants an on-site demo.", "Decision expected end of month.",
	}
	if len(leadIDs) == 0 {
		return nil
	}
	for i := 0; i < n; i++ {
		createdAt := spread(r, now)
		_, err := db.Exec(ctx, `
			INSERT INTO notes (id, entity_type, entity_id, note, created_by, created_at, updated_at)
			VALUES (gen_random_uuid(), 'LEAD', $1, $2, $3, $4, $4)
		`, leadIDs[r.Intn(len(leadIDs))], pick(r, texts), owner, createdAt)
		if err != nil {
			return fmt.Errorf("insert note: %w", err)
		}
	}
	return nil
}

func seedCompanies(ctx context.Context, db *pgxpool.Pool, r *rand.Rand, now time.Time, orgID, moduleID, owner string, n int) ([]string, error) {
	ids := make([]string, 0, n)
	for i := 0; i < n; i++ {
		name := companyNames[i%len(companyNames)]
		data := map[string]any{
			"name":      name,
			"industry":  pick(r, industries),
			"city":      pick(r, cities),
			"website":   website(name),
			"employees": (r.Intn(50) + 1) * 20,
			"tags":      pickTags(r),
		}
		id, err := insertRecord(ctx, db, orgID, moduleID, owner, data, spread(r, now))
		if err != nil {
			return nil, err
		}
		ids = append(ids, id)
	}
	return ids, nil
}

func seedDeals(ctx context.Context, db *pgxpool.Pool, r *rand.Rand, now time.Time, orgID, moduleID, owner string, companyIDs []string, n int) error {
	for i := 0; i < n; i++ {
		var companyRef any
		if len(companyIDs) > 0 {
			companyRef = companyIDs[r.Intn(len(companyIDs))]
		}
		data := map[string]any{
			"title":      fmt.Sprintf("Deal #%d", 1000+i),
			"amount":     (r.Intn(90) + 10) * 10000,
			"stage":      pick(r, dealStages),
			"close_date": now.AddDate(0, 0, r.Intn(90)).Format("2006-01-02"),
			"company":    companyRef,
		}
		if _, err := insertRecord(ctx, db, orgID, moduleID, owner, data, spread(r, now)); err != nil {
			return err
		}
	}
	return nil
}

func insertRecord(ctx context.Context, db *pgxpool.Pool, orgID, moduleID, owner string, data map[string]any, createdAt time.Time) (string, error) {
	payload, err := json.Marshal(data)
	if err != nil {
		return "", err
	}
	var id string
	err = db.QueryRow(ctx, `
		INSERT INTO records (organization_id, module_id, data, owner_id, created_by, updated_by, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $4, $4, $5, $5)
		RETURNING id
	`, orgID, moduleID, string(payload), owner, createdAt).Scan(&id)
	if err != nil {
		return "", fmt.Errorf("insert record: %w", err)
	}
	return id, nil
}
