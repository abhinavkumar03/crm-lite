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
	{"company", "Company", "Companies", "building-2", "#8b5cf6", "dynamic", "", 1},
	{"contact", "Contact", "Contacts", "users", "#06b6d4", "dynamic", "", 2},
	{"deal", "Deal", "Deals", "handshake", "#ec4899", "dynamic", "", 3},
	{"lead", "Lead", "Leads", "user-plus", "#f59e0b", "dynamic", "", 4},
	{"task", "Task", "Tasks", "check-square", "#10b981", "dynamic", "", 5},
}

type ModulesSeeder struct{}

func (ModulesSeeder) Name() string { return "modules" }

func (ModulesSeeder) Run(ctx context.Context, db *pgxpool.Pool) error {
	orgIDs, err := listDemoOrgIDs(ctx, db)
	if err != nil {
		return err
	}

	for _, orgID := range orgIDs {
		for _, m := range moduleDefs {
			var nativeTable any
			if m.NativeTable != "" {
				nativeTable = m.NativeTable
			}
			_, err := db.Exec(ctx, `
				INSERT INTO modules (
					organization_id, api_name, singular_label, plural_label,
					icon, color, storage_strategy, native_table, is_system, sort_order,
					is_enabled, is_visible_sidebar
				)
				VALUES ($1, $2, $3, $4, $5, $6, $7, $8, TRUE, $9, TRUE, TRUE)
				ON CONFLICT (organization_id, api_name) DO UPDATE
				SET singular_label    = EXCLUDED.singular_label,
				    plural_label      = EXCLUDED.plural_label,
				    icon              = EXCLUDED.icon,
				    color             = EXCLUDED.color,
				    sort_order        = EXCLUDED.sort_order,
				    storage_strategy  = EXCLUDED.storage_strategy,
				    native_table      = EXCLUDED.native_table,
				    is_enabled        = TRUE,
				    is_visible_sidebar = TRUE
			`, orgID, m.APIName, m.Singular, m.Plural, m.Icon, m.Color, m.Storage, nativeTable, m.Sort)
			if err != nil {
				return fmt.Errorf("upsert module %q: %w", m.APIName, err)
			}
		}

		if _, err := db.Exec(ctx, `
			DELETE FROM modules
			WHERE organization_id = $1
			  AND api_name IN ('lead', 'task')
			  AND storage_strategy = 'native'
		`, orgID); err != nil {
			return fmt.Errorf("remove native modules: %w", err)
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
	// company
	{"company", "name", "Company Name", "text", true, true, false, nil, ""},
	{"company", "industry", "Industry", "dropdown", false, false, true, industries, ""},
	{"company", "status", "Status", "dropdown", false, false, true, companyStatuses, ""},
	{"company", "city", "City", "text", false, true, true, nil, ""},
	{"company", "country", "Country", "text", false, true, true, nil, ""},
	{"company", "website", "Website", "url", false, false, false, nil, ""},
	{"company", "phone", "Phone", "phone", false, true, false, nil, ""},
	{"company", "email", "Email", "email", false, true, false, nil, ""},
	{"company", "employees", "Employees", "number", false, false, true, nil, ""},
	{"company", "annual_revenue", "Annual Revenue", "currency", false, false, true, nil, ""},
	{"company", "linkedin", "LinkedIn", "url", false, false, false, nil, ""},
	{"company", "priority", "Priority", "dropdown", false, false, true, []string{"Low", "Medium", "High", "Urgent"}, ""},
	{"company", "tags", "Tags", "multiselect", false, false, true, tagPool, ""},
	{"company", "last_contacted", "Last Contacted", "date", false, false, true, nil, ""},
	{"company", "next_follow_up", "Next Follow-up", "date", false, false, true, nil, ""},
	{"company", "description", "Description", "textarea", false, false, false, nil, ""},

	// contact
	{"contact", "first_name", "First Name", "text", true, true, false, nil, ""},
	{"contact", "last_name", "Last Name", "text", true, true, false, nil, ""},
	{"contact", "email", "Email", "email", false, true, true, nil, ""},
	{"contact", "phone", "Phone", "phone", false, true, false, nil, ""},
	{"contact", "mobile", "Mobile", "phone", false, true, false, nil, ""},
	{"contact", "job_title", "Job Title", "text", false, true, true, nil, ""},
	{"contact", "department", "Department", "text", false, false, true, nil, ""},
	{"contact", "company", "Company", "lookup", false, false, true, nil, "company"},
	{"contact", "city", "City", "text", false, true, true, nil, ""},
	{"contact", "country", "Country", "text", false, true, true, nil, ""},
	{"contact", "linkedin", "LinkedIn", "url", false, false, false, nil, ""},
	{"contact", "priority", "Priority", "dropdown", false, false, true, []string{"Low", "Medium", "High", "Urgent"}, ""},
	{"contact", "rating", "Rating", "dropdown", false, false, true, []string{"Hot", "Warm", "Cold"}, ""},
	{"contact", "last_contacted", "Last Contacted", "date", false, false, true, nil, ""},
	{"contact", "next_follow_up", "Next Follow-up", "date", false, false, true, nil, ""},
	{"contact", "tags", "Tags", "multiselect", false, false, true, tagPool, ""},
	{"contact", "notes", "Notes", "textarea", false, false, false, nil, ""},

	// deal
	{"deal", "title", "Deal Title", "text", true, true, false, nil, ""},
	{"deal", "amount", "Amount", "currency", false, false, true, nil, ""},
	{"deal", "stage", "Stage", "dropdown", false, false, true, dealStages, ""},
	{"deal", "probability", "Probability %", "number", false, false, true, nil, ""},
	{"deal", "close_date", "Close Date", "date", false, false, true, nil, ""},
	{"deal", "next_step", "Next Step", "text", false, false, false, nil, ""},
	{"deal", "company", "Company", "lookup", false, false, true, nil, "company"},
	{"deal", "contact_name", "Primary Contact", "text", false, true, false, nil, ""},
	{"deal", "source", "Source", "dropdown", false, false, true, leadSources, ""},
	{"deal", "priority", "Priority", "dropdown", false, false, true, []string{"Low", "Medium", "High", "Urgent"}, ""},
	{"deal", "city", "City", "text", false, true, true, nil, ""},
	{"deal", "country", "Country", "text", false, true, true, nil, ""},
	{"deal", "expected_revenue", "Expected Revenue", "currency", false, false, true, nil, ""},
	{"deal", "tags", "Tags", "multiselect", false, false, true, tagPool, ""},
	{"deal", "description", "Description", "textarea", false, false, false, nil, ""},

	// lead
	{"lead", "first_name", "First Name", "text", true, true, false, nil, ""},
	{"lead", "last_name", "Last Name", "text", true, true, false, nil, ""},
	{"lead", "email", "Email", "email", false, true, true, nil, ""},
	{"lead", "phone", "Phone", "phone", false, true, false, nil, ""},
	{"lead", "company_name", "Company Name", "text", false, true, true, nil, ""},
	{"lead", "company", "Company", "lookup", false, false, true, nil, "company"},
	{"lead", "job_title", "Job Title", "text", false, true, true, nil, ""},
	{"lead", "industry", "Industry", "dropdown", false, false, true, industries, ""},
	{"lead", "status", "Status", "dropdown", true, false, true, leadStatusOptions, ""},
	{"lead", "source", "Source", "dropdown", false, false, true, leadSources, ""},
	{"lead", "website", "Website", "url", false, false, false, nil, ""},
	{"lead", "employees", "Employees", "number", false, false, true, nil, ""},
	{"lead", "annual_revenue", "Annual Revenue", "currency", false, false, true, nil, ""},
	{"lead", "city", "City", "text", false, true, true, nil, ""},
	{"lead", "country", "Country", "text", false, true, true, nil, ""},
	{"lead", "priority", "Priority", "dropdown", false, false, true, []string{"Low", "Medium", "High", "Urgent"}, ""},
	{"lead", "rating", "Rating", "dropdown", false, false, true, []string{"Hot", "Warm", "Cold"}, ""},
	{"lead", "last_contacted", "Last Contacted", "date", false, false, true, nil, ""},
	{"lead", "next_follow_up", "Next Follow-up", "date", false, false, true, nil, ""},
	{"lead", "linkedin", "LinkedIn", "url", false, false, false, nil, ""},
	{"lead", "tags", "Tags", "multiselect", false, false, true, tagPool, ""},
	{"lead", "notes", "Notes", "textarea", false, false, false, nil, ""},

	// task
	{"task", "title", "Task Title", "text", true, true, false, nil, ""},
	{"task", "status", "Status", "dropdown", true, false, true, taskStatusOptions, ""},
	{"task", "priority", "Priority", "dropdown", false, false, true, taskPriorities, ""},
	{"task", "due_date", "Due Date", "date", false, false, true, nil, ""},
	{"task", "company", "Company", "lookup", false, false, true, nil, "company"},
	{"task", "related_lead", "Related Lead", "text", false, true, false, nil, ""},
	{"task", "related_contact", "Related Contact", "text", false, true, false, nil, ""},
	{"task", "city", "City", "text", false, true, true, nil, ""},
	{"task", "country", "Country", "text", false, true, true, nil, ""},
	{"task", "estimated_hours", "Estimated Hours", "number", false, false, true, nil, ""},
	{"task", "completed_at", "Completed At", "date", false, false, true, nil, ""},
	{"task", "reminder_date", "Reminder Date", "date", false, false, true, nil, ""},
	{"task", "tags", "Tags", "multiselect", false, false, true, tagPool, ""},
	{"task", "outcome", "Outcome", "text", false, false, false, nil, ""},
	{"task", "description", "Description", "textarea", false, false, false, nil, ""},
}

type FieldsSeeder struct{}

func (FieldsSeeder) Name() string { return "fields" }

func (FieldsSeeder) Run(ctx context.Context, db *pgxpool.Pool) error {
	orgIDs, err := listDemoOrgIDs(ctx, db)
	if err != nil {
		return err
	}

	for _, orgID := range orgIDs {
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
					lookup_module_id, sort_order, is_system, is_visible
				)
				VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, TRUE, TRUE)
				ON CONFLICT (module_id, api_name) DO UPDATE
				SET label = EXCLUDED.label,
				    field_type = EXCLUDED.field_type,
				    is_required = EXCLUDED.is_required,
				    is_searchable = EXCLUDED.is_searchable,
				    is_filterable = EXCLUDED.is_filterable,
				    options = EXCLUDED.options,
				    lookup_module_id = EXCLUDED.lookup_module_id,
				    sort_order = EXCLUDED.sort_order,
				    is_visible = TRUE
			`, orgID, moduleID, f.APIName, f.Label, f.Type,
				f.Required, f.Searchable, f.Filterable, string(optionsJSON),
				lookupID, sortByModule[f.Module])
			if err != nil {
				return fmt.Errorf("insert field %s.%s: %w", f.Module, f.APIName, err)
			}
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
	{"dashboard", "Dashboard", "Your command centre — module counts and recent records.", "[data-tour=dashboard]", "/dashboard", "bottom"},
	{"forms", "Form Designer", "Preview how create forms are generated from module metadata. Create real records from each module’s Add button.", "[data-tutorial-action=\"open-forms\"]", "/settings/forms", "right"},
	{"tables", "Your modules", "Each module in the workspace sidebar has its own page — view, edit, and delete records there.", "[data-tour=\"sidebar-nav\"]", "/dashboard", "right"},
	{"modules", "Modules", "Add new modules and custom fields with no code.", "[data-tour=\"nav-settings\"]", "/settings/modules", "right"},
	{"import", "Import", "Bring in data from CSV or Excel with a guided wizard.", "[data-tutorial-action=\"open-imports\"]", "/settings/imports", "right"},
	{"export", "Export", "Export filtered data for reporting.", "[data-tutorial-action=\"open-exports\"]", "/settings/exports", "right"},
	{"finish", "You're all set", "Restart this tour anytime from Help or the Explore CRM button.", "", "/dashboard", "center"},
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
