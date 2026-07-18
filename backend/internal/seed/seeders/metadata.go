package seeders

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

// ---------------------------------------------------------------------------
// Modules
// ---------------------------------------------------------------------------

type moduleDef struct {
	APIName     string
	Singular    string
	Plural      string
	Icon        string
	Color       string
	Storage     string // native | dynamic
	NativeTable string // set when Storage == native
	Sort        int
}

var moduleDefs = []moduleDef{
	{"lead", "Lead", "Leads", "target", "#10b981", "native", "leads", 1},
	{"contact", "Contact", "Contacts", "users", "#3b82f6", "native", "contacts", 2},
	{"task", "Task", "Tasks", "check-square", "#f59e0b", "native", "tasks", 3},
	{"company", "Company", "Companies", "building-2", "#8b5cf6", "dynamic", "", 4},
	{"deal", "Deal", "Deals", "handshake", "#ec4899", "dynamic", "", 5},
}

type ModulesSeeder struct{}

func (ModulesSeeder) Name() string { return "modules" }

func (ModulesSeeder) Run(ctx context.Context, db *pgxpool.Pool) error {
	orgID, err := getOrgID(ctx, db)
	if err != nil {
		return err
	}

	for _, m := range moduleDefs {
		var nativeTable any
		if m.NativeTable != "" {
			nativeTable = m.NativeTable
		}
		_, err := db.Exec(ctx, `
			INSERT INTO modules (
				organization_id, api_name, singular_label, plural_label,
				icon, color, storage_strategy, native_table, is_system, sort_order
			)
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8, TRUE, $9)
			ON CONFLICT (organization_id, api_name) DO UPDATE
			SET singular_label = EXCLUDED.singular_label,
			    plural_label   = EXCLUDED.plural_label,
			    icon           = EXCLUDED.icon,
			    color          = EXCLUDED.color,
			    sort_order     = EXCLUDED.sort_order
		`, orgID, m.APIName, m.Singular, m.Plural, m.Icon, m.Color, m.Storage, nativeTable, m.Sort)
		if err != nil {
			return fmt.Errorf("upsert module %q: %w", m.APIName, err)
		}
	}
	return nil
}

func getModuleID(ctx context.Context, db *pgxpool.Pool, orgID, apiName string) (string, error) {
	var id string
	err := db.QueryRow(ctx,
		`SELECT id FROM modules WHERE organization_id = $1 AND api_name = $2`,
		orgID, apiName,
	).Scan(&id)
	if err != nil {
		return "", fmt.Errorf("lookup module %q: %w", apiName, err)
	}
	return id, nil
}

// ---------------------------------------------------------------------------
// Fields
// ---------------------------------------------------------------------------

type fieldDef struct {
	Module     string
	APIName    string
	Label      string
	Type       string
	Required   bool
	Searchable bool
	Filterable bool
	Options    []string
	LookupMod  string // api_name of the module this lookup points to
}

var fieldDefs = []fieldDef{
	// lead
	{"lead", "name", "Name", "text", true, true, false, nil, ""},
	{"lead", "email", "Email", "email", false, true, false, nil, ""},
	{"lead", "phone", "Phone", "phone", false, false, false, nil, ""},
	{"lead", "company", "Company", "text", false, true, true, nil, ""},
	{"lead", "status", "Status", "dropdown", true, false, true, []string{"NEW", "CONTACTED", "QUALIFIED", "WON", "LOST"}, ""},
	{"lead", "notes", "Notes", "textarea", false, false, false, nil, ""},

	// contact
	{"contact", "first_name", "First Name", "text", true, true, false, nil, ""},
	{"contact", "last_name", "Last Name", "text", false, true, false, nil, ""},
	{"contact", "email", "Email", "email", false, true, false, nil, ""},
	{"contact", "phone", "Phone", "phone", false, false, false, nil, ""},
	{"contact", "company", "Company", "text", false, false, true, nil, ""},
	{"contact", "job_title", "Job Title", "text", false, false, false, nil, ""},
	{"contact", "notes", "Notes", "textarea", false, false, false, nil, ""},

	// task
	{"task", "title", "Title", "text", true, true, false, nil, ""},
	{"task", "description", "Description", "textarea", false, false, false, nil, ""},
	{"task", "status", "Status", "dropdown", true, false, true, []string{"PENDING", "IN_PROGRESS", "COMPLETED"}, ""},
	{"task", "due_date", "Due Date", "datetime", false, false, true, nil, ""},

	// company (dynamic)
	{"company", "name", "Company Name", "text", true, true, false, nil, ""},
	{"company", "industry", "Industry", "dropdown", false, false, true, industries, ""},
	{"company", "city", "City", "text", false, true, true, nil, ""},
	{"company", "website", "Website", "url", false, false, false, nil, ""},
	{"company", "employees", "Employees", "number", false, false, true, nil, ""},
	{"company", "tags", "Tags", "multiselect", false, false, true, tagPool, ""},

	// deal (dynamic)
	{"deal", "title", "Deal Title", "text", true, true, false, nil, ""},
	{"deal", "amount", "Amount", "currency", false, false, true, nil, ""},
	{"deal", "stage", "Stage", "dropdown", false, false, true, dealStages, ""},
	{"deal", "close_date", "Close Date", "date", false, false, true, nil, ""},
	{"deal", "company", "Company", "lookup", false, false, true, nil, "company"},
}

type FieldsSeeder struct{}

func (FieldsSeeder) Name() string { return "fields" }

func (FieldsSeeder) Run(ctx context.Context, db *pgxpool.Pool) error {
	orgID, err := getOrgID(ctx, db)
	if err != nil {
		return err
	}

	// Cache module ids.
	moduleIDs := map[string]string{}
	for _, m := range moduleDefs {
		id, err := getModuleID(ctx, db, orgID, m.APIName)
		if err != nil {
			return err
		}
		moduleIDs[m.APIName] = id
	}

	sortByModule := map[string]int{}
	for _, f := range fieldDefs {
		moduleID := moduleIDs[f.Module]
		sortByModule[f.Module]++

		optionsJSON, err := json.Marshal(f.Options)
		if err != nil {
			return err
		}
		if f.Options == nil {
			optionsJSON = []byte("[]")
		}

		var lookupID any
		if f.LookupMod != "" {
			lookupID = moduleIDs[f.LookupMod]
		}

		_, err = db.Exec(ctx, `
			INSERT INTO fields (
				organization_id, module_id, api_name, label, field_type,
				is_required, is_searchable, is_filterable, options,
				lookup_module_id, sort_order, is_system
			)
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, TRUE)
			ON CONFLICT (module_id, api_name) DO NOTHING
		`, orgID, moduleID, f.APIName, f.Label, f.Type,
			f.Required, f.Searchable, f.Filterable, string(optionsJSON),
			lookupID, sortByModule[f.Module])
		if err != nil {
			return fmt.Errorf("insert field %s.%s: %w", f.Module, f.APIName, err)
		}
	}
	return nil
}

// ---------------------------------------------------------------------------
// Tour steps
// ---------------------------------------------------------------------------

type tourStep struct {
	Key       string
	Title     string
	Body      string
	Target    string
	Route     string
	Placement string
}

var tourSteps = []tourStep{
	{"welcome", "Welcome to CRM Lite", "This quick tour shows the main areas of the CRM.", "", "/dashboard", "center"},
	{"sidebar", "Navigation", "Switch between modules from the sidebar. Modules are fully configurable.", "[data-tour=sidebar]", "/dashboard", "right"},
	{"dashboard", "Dashboard", "Your command centre — key metrics and recent activity.", "[data-tour=dashboard]", "/dashboard", "bottom"},
	{"leads", "Leads", "Manage your sales pipeline and track lead status.", "[data-tour=leads]", "/leads", "bottom"},
	{"contacts", "Contacts", "Keep all your business contacts in one place.", "[data-tour=contacts]", "/contacts", "bottom"},
	{"tasks", "Tasks", "Create tasks and link them to leads and contacts.", "[data-tour=tasks]", "/tasks", "bottom"},
	{"modules", "Dynamic Modules", "Add new modules and custom fields with no code.", "[data-tour=settings-modules]", "/settings", "right"},
	{"import", "Import", "Bring in data from CSV or Excel with a guided wizard.", "[data-tour=import]", "/settings", "bottom"},
	{"export", "Export", "Export filtered data for reporting.", "[data-tour=export]", "/settings", "bottom"},
	{"automation", "Automation", "Trigger notifications and actions on record events.", "[data-tour=automation]", "/settings", "bottom"},
	{"finish", "You're all set", "Restart this tour anytime from the Explore CRM button.", "", "/dashboard", "center"},
}

type TourStepsSeeder struct{}

func (TourStepsSeeder) Name() string { return "tour_steps" }

func (TourStepsSeeder) Run(ctx context.Context, db *pgxpool.Pool) error {
	orgID, err := getOrgID(ctx, db)
	if err != nil {
		return err
	}

	// Tour steps are fully seed-owned; replace them wholesale so edits to the
	// list are reflected on re-run.
	if _, err := db.Exec(ctx, `DELETE FROM tour_steps WHERE organization_id = $1`, orgID); err != nil {
		return err
	}

	for i, s := range tourSteps {
		_, err := db.Exec(ctx, `
			INSERT INTO tour_steps (
				organization_id, step_key, title, body,
				target_selector, route, placement, sort_order
			)
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		`, orgID, s.Key, s.Title, s.Body, nullable(s.Target), s.Route, s.Placement, i+1)
		if err != nil {
			return fmt.Errorf("insert tour step %q: %w", s.Key, err)
		}
	}
	return nil
}

// nullable returns nil for empty strings so NULL is stored instead of "".
func nullable(s string) any {
	if s == "" {
		return nil
	}
	return s
}
