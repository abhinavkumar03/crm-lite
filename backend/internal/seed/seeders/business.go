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
// dynamic company / contact / deal / lead / task records in the JSONB record engine.
//
// Idempotent per module: skips modules that already have records.
type BusinessDataSeeder struct{}

func (BusinessDataSeeder) Name() string { return "business_data" }

func (BusinessDataSeeder) Run(ctx context.Context, db *pgxpool.Pool) error {
	orgIDs, err := listDemoOrgIDs(ctx, db)
	if err != nil {
		return err
	}

	ownerIDs, err := demoOwnerIDs(ctx, db)
	if err != nil {
		return err
	}

	now := time.Now()
	for orgIdx, orgID := range orgIDs {
		r := rand.New(rand.NewSource(42 + int64(orgIdx)*1000))

		companyModuleID, err := getModuleID(ctx, db, orgID, "company")
		if err != nil {
			return err
		}
		contactModuleID, err := getModuleID(ctx, db, orgID, "contact")
		if err != nil {
			return err
		}
		dealModuleID, err := getModuleID(ctx, db, orgID, "deal")
		if err != nil {
			return err
		}
		leadModuleID, err := getModuleID(ctx, db, orgID, "lead")
		if err != nil {
			return err
		}
		taskModuleID, err := getModuleID(ctx, db, orgID, "task")
		if err != nil {
			return err
		}

		companyRecordIDs, err := ensureCompanies(ctx, db, r, now, orgID, companyModuleID, ownerIDs, 20)
		if err != nil {
			return err
		}
		if err := ensureContacts(ctx, db, r, now, orgID, contactModuleID, ownerIDs, companyRecordIDs, 22); err != nil {
			return err
		}
		if err := ensureDeals(ctx, db, r, now, orgID, dealModuleID, ownerIDs, companyRecordIDs, 20); err != nil {
			return err
		}
		if err := ensureLeads(ctx, db, r, now, orgID, leadModuleID, ownerIDs, companyRecordIDs, 28); err != nil {
			return err
		}
		if err := ensureTasks(ctx, db, r, now, orgID, taskModuleID, ownerIDs, companyRecordIDs, 22); err != nil {
			return err
		}
	}
	return nil
}

func moduleRecordCount(ctx context.Context, db *pgxpool.Pool, orgID, moduleID string) (int, error) {
	var n int
	err := db.QueryRow(ctx,
		`SELECT count(*) FROM records WHERE organization_id = $1 AND module_id = $2`,
		orgID, moduleID,
	).Scan(&n)
	return n, err
}

func demoOwnerIDs(ctx context.Context, db *pgxpool.Pool) ([]string, error) {
	emails := []string{
		demoUserEmail,
		"admin@crm.com",
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

func ensureCompanies(ctx context.Context, db *pgxpool.Pool, r *rand.Rand, now time.Time, orgID, moduleID string, owners []string, n int) ([]string, error) {
	existing, err := moduleRecordCount(ctx, db, orgID, moduleID)
	if err != nil {
		return nil, err
	}
	ids, err := listRecordIDs(ctx, db, orgID, moduleID, n)
	if err != nil {
		return nil, err
	}
	if existing >= n {
		return ids, nil
	}
	for i := existing; i < n; i++ {
		name := companyNames[i%len(companyNames)]
		if i >= len(companyNames) {
			name = fmt.Sprintf("%s %d", name, i)
		}
		owner := owners[i%len(owners)]
		data := map[string]any{
			"name":           name,
			"industry":       pick(r, industries),
			"status":         pick(r, companyStatuses),
			"city":           pick(r, cities),
			"country":        pick(r, countries),
			"website":        website(name),
			"phone":          phone(r),
			"email":          fmt.Sprintf("hello@%s", website(name)[len("https://"):]),
			"employees":      (r.Intn(50) + 1) * 20,
			"annual_revenue": (r.Intn(90) + 10) * 1_000_000,
			"priority":       pick(r, []string{"Low", "Medium", "High", "Urgent"}),
			"tags":           pickTags(r),
			"description":    fmt.Sprintf("%s is a growing account in the %s sector.", name, pick(r, industries)),
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

func ensureContacts(ctx context.Context, db *pgxpool.Pool, r *rand.Rand, now time.Time, orgID, moduleID string, owners []string, companyIDs []string, n int) error {
	existing, err := moduleRecordCount(ctx, db, orgID, moduleID)
	if err != nil {
		return err
	}
	if existing >= n {
		return nil
	}
	for i := existing; i < n; i++ {
		first, last := fullName(r)
		var companyRef any
		if len(companyIDs) > 0 {
			companyRef = companyIDs[r.Intn(len(companyIDs))]
		}
		owner := owners[i%len(owners)]
		data := map[string]any{
			"first_name": first,
			"last_name":  last,
			"email":      emailFrom(first, last, r),
			"phone":      phone(r),
			"mobile":     phone(r),
			"job_title":  pick(r, jobTitles),
			"department": pick(r, []string{"Sales", "Marketing", "Operations", "Finance", "IT"}),
			"company":    companyRef,
			"city":       pick(r, cities),
			"country":    pick(r, countries),
			"priority":   pick(r, []string{"Low", "Medium", "High", "Urgent"}),
			"rating":     pick(r, []string{"Hot", "Warm", "Cold"}),
			"tags":       pickTags(r),
			"notes":      fmt.Sprintf("Met at industry meetup — follow up on Q%d priorities.", 1+r.Intn(4)),
		}
		vis := "organization"
		if i%4 == 0 {
			vis = "owner"
		}
		if _, err := insertRecord(ctx, db, orgID, moduleID, owner, data, spread(r, now), vis); err != nil {
			return err
		}
	}
	return nil
}

func ensureDeals(ctx context.Context, db *pgxpool.Pool, r *rand.Rand, now time.Time, orgID, moduleID string, owners []string, companyIDs []string, n int) error {
	existing, err := moduleRecordCount(ctx, db, orgID, moduleID)
	if err != nil {
		return err
	}
	if existing >= n {
		return nil
	}
	titles := []string{
		"Enterprise license expansion", "Annual support renewal", "Pilot rollout",
		"Platform migration", "Multi-year MSA", "Upsell analytics add-on",
		"Regional franchise deal", "Cloud migration package",
		"Security audit package", "Training & onboarding bundle",
		"API integration project", "Data warehouse upgrade",
		"Customer success retain", "Partner channel expansion",
		"Mobile app redesign", "Compliance readiness",
		"AI assistant pilot", "Warehouse automation",
		"Marketing automation suite", "Support SLA upgrade",
	}
	for i := existing; i < n; i++ {
		var companyRef any
		if len(companyIDs) > 0 {
			companyRef = companyIDs[r.Intn(len(companyIDs))]
		}
		owner := owners[(i+1)%len(owners)]
		stage := pick(r, dealStages)
		prob := map[string]int{
			"Prospecting": 10, "Qualification": 25, "Proposal": 50,
			"Negotiation": 75, "Closed Won": 100, "Closed Lost": 0,
		}[stage]
		data := map[string]any{
			"title":       titles[i%len(titles)],
			"amount":      (r.Intn(90) + 10) * 10000,
			"stage":       stage,
			"probability": prob,
			"close_date":  now.AddDate(0, 0, r.Intn(90)).Format("2006-01-02"),
			"next_step":   pick(r, []string{"Send proposal", "Schedule demo", "Legal review", "Price negotiation", "Kickoff call"}),
			"company":     companyRef,
			"description": fmt.Sprintf("Pipeline opportunity #%d — %s.", 1000+i, stage),
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

func ensureLeads(ctx context.Context, db *pgxpool.Pool, r *rand.Rand, now time.Time, orgID, moduleID string, owners []string, companyIDs []string, n int) error {
	existing, err := moduleRecordCount(ctx, db, orgID, moduleID)
	if err != nil {
		return err
	}
	if existing >= n {
		return nil
	}
	for i := existing; i < n; i++ {
		first, last := fullName(r)
		var companyRef any
		if len(companyIDs) > 0 {
			companyRef = companyIDs[r.Intn(len(companyIDs))]
		}
		owner := owners[i%len(owners)]
		data := map[string]any{
			"first_name":     first,
			"last_name":      last,
			"email":          emailFrom(first, last, r),
			"phone":          phone(r),
			"company_name":   pick(r, companyNames),
			"company":        companyRef,
			"job_title":      pick(r, jobTitles),
			"industry":       pick(r, industries),
			"status":         pick(r, leadStatuses),
			"source":         pick(r, leadSources),
			"website":        website(pick(r, companyNames)),
			"employees":      (r.Intn(40) + 1) * 10,
			"annual_revenue": (r.Intn(50) + 5) * 100_000,
			"city":           pick(r, cities),
			"country":        pick(r, countries),
			"priority":       pick(r, []string{"Low", "Medium", "High", "Urgent"}),
			"rating":         pick(r, []string{"Hot", "Warm", "Cold"}),
			"tags":           pickTags(r),
			"notes":          fmt.Sprintf("Inbound lead #%d — nurture for discovery call.", 2000+i),
		}
		vis := "organization"
		if i%3 == 0 {
			vis = "owner"
		}
		if _, err := insertRecord(ctx, db, orgID, moduleID, owner, data, spread(r, now), vis); err != nil {
			return err
		}
	}
	return nil
}

func ensureTasks(ctx context.Context, db *pgxpool.Pool, r *rand.Rand, now time.Time, orgID, moduleID string, owners []string, companyIDs []string, n int) error {
	existing, err := moduleRecordCount(ctx, db, orgID, moduleID)
	if err != nil {
		return err
	}
	if existing >= n {
		return nil
	}
	titles := []string{
		"Follow up after demo", "Send pricing deck", "Schedule discovery call",
		"Prepare contract draft", "Confirm stakeholder list", "Update CRM notes",
		"Share case study pack", "Book technical deep-dive", "Review security questionnaire",
		"Collect purchase order", "Intro partner AE", "Qualify budget timeline",
		"Renewal check-in", "Onboarding kickoff prep", "Clarify success metrics",
		"Escalate blocked deal", "Customer health review", "Log support ticket themes",
		"Propose pilot scope", "Confirm go-live date",
	}
	for i := existing; i < n; i++ {
		var companyRef any
		if len(companyIDs) > 0 {
			companyRef = companyIDs[r.Intn(len(companyIDs))]
		}
		owner := owners[(i+2)%len(owners)]
		data := map[string]any{
			"title":       titles[i%len(titles)],
			"status":      pick(r, taskStatuses),
			"priority":    pick(r, taskPriorities),
			"due_date":    now.AddDate(0, 0, r.Intn(45)-7).Format("2006-01-02"),
			"company":     companyRef,
			"description": fmt.Sprintf("Demo task #%d — keep the account moving.", 3000+i),
		}
		vis := "organization"
		if i%4 == 0 {
			vis = "owner"
		}
		if _, err := insertRecord(ctx, db, orgID, moduleID, owner, data, spread(r, now), vis); err != nil {
			return err
		}
	}
	return nil
}

func listRecordIDs(ctx context.Context, db *pgxpool.Pool, orgID, moduleID string, limit int) ([]string, error) {
	rows, err := db.Query(ctx, `
		SELECT id::text FROM records
		WHERE organization_id = $1 AND module_id = $2
		ORDER BY created_at
		LIMIT $3
	`, orgID, moduleID, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	out := make([]string, 0)
	for rows.Next() {
		var id string
		if err := rows.Scan(&id); err != nil {
			return nil, err
		}
		out = append(out, id)
	}
	return out, rows.Err()
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
