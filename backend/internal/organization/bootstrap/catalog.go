package bootstrap

import (
	"context"
	"encoding/json"
	"fmt"
)

// Full CRM catalog seeded for every new workspace (parity with demo seeders).

var defaultModuleCatalog = []struct {
	API, Singular, Plural, Icon, Color string
	Sort                               int
}{
	{"company", "Company", "Companies", "building-2", "#8b5cf6", 1},
	{"contact", "Contact", "Contacts", "users", "#06b6d4", 2},
	{"deal", "Deal", "Deals", "handshake", "#ec4899", 3},
	{"lead", "Lead", "Leads", "user-plus", "#f59e0b", 4},
	{"task", "Task", "Tasks", "check-square", "#10b981", 5},
}

var (
	optIndustries = []string{
		"IT Services", "Manufacturing", "Retail", "Fintech", "Healthcare",
		"Education", "Logistics", "Real Estate", "FMCG", "Automotive",
	}
	optCompanyStatuses = []string{"Prospect", "Active", "Partner", "Churned"}
	optLeadStatuses    = []string{"NEW", "CONTACTED", "QUALIFIED", "WON", "LOST"}
	optLeadSources    = []string{"Website", "Referral", "Cold Call", "Trade Show", "LinkedIn", "Email Campaign"}
	optTaskStatuses    = []string{"PENDING", "IN_PROGRESS", "COMPLETED"}
	optTaskPriorities = []string{"Low", "Medium", "High"}
	optDealStages     = []string{"Prospecting", "Qualification", "Proposal", "Negotiation", "Closed Won", "Closed Lost"}
	optPriorities     = []string{"Low", "Medium", "High", "Urgent"}
	optRatings        = []string{"Hot", "Warm", "Cold"}
	optTagPool        = []string{"vip", "inbound", "referral", "enterprise", "smb", "hot", "cold"}
)

type catalogField struct {
	Module, API, Label, Type string
	Required, Searchable     bool
	Filterable               bool
	Options                  []string
	LookupMod                string
}

// 15+ fields per primary module (lead/contact/task/company); deal keeps a full set too.
var defaultFieldCatalog = []catalogField{
	// company (15+)
	{"company", "name", "Company Name", "text", true, true, false, nil, ""},
	{"company", "industry", "Industry", "dropdown", false, false, true, optIndustries, ""},
	{"company", "status", "Status", "dropdown", false, false, true, optCompanyStatuses, ""},
	{"company", "city", "City", "text", false, true, true, nil, ""},
	{"company", "country", "Country", "text", false, true, true, nil, ""},
	{"company", "website", "Website", "url", false, false, false, nil, ""},
	{"company", "phone", "Phone", "phone", false, true, false, nil, ""},
	{"company", "email", "Email", "email", false, true, false, nil, ""},
	{"company", "employees", "Employees", "number", false, false, true, nil, ""},
	{"company", "annual_revenue", "Annual Revenue", "currency", false, false, true, nil, ""},
	{"company", "linkedin", "LinkedIn", "url", false, false, false, nil, ""},
	{"company", "priority", "Priority", "dropdown", false, false, true, optPriorities, ""},
	{"company", "tags", "Tags", "multiselect", false, false, true, optTagPool, ""},
	{"company", "last_contacted", "Last Contacted", "date", false, false, true, nil, ""},
	{"company", "next_follow_up", "Next Follow-up", "date", false, false, true, nil, ""},
	{"company", "description", "Description", "textarea", false, false, false, nil, ""},

	// contact (15+)
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
	{"contact", "priority", "Priority", "dropdown", false, false, true, optPriorities, ""},
	{"contact", "rating", "Rating", "dropdown", false, false, true, optRatings, ""},
	{"contact", "last_contacted", "Last Contacted", "date", false, false, true, nil, ""},
	{"contact", "next_follow_up", "Next Follow-up", "date", false, false, true, nil, ""},
	{"contact", "tags", "Tags", "multiselect", false, false, true, optTagPool, ""},
	{"contact", "notes", "Notes", "textarea", false, false, false, nil, ""},

	// deal
	{"deal", "title", "Deal Title", "text", true, true, false, nil, ""},
	{"deal", "amount", "Amount", "currency", false, false, true, nil, ""},
	{"deal", "stage", "Stage", "dropdown", false, false, true, optDealStages, ""},
	{"deal", "probability", "Probability %", "number", false, false, true, nil, ""},
	{"deal", "close_date", "Close Date", "date", false, false, true, nil, ""},
	{"deal", "next_step", "Next Step", "text", false, false, false, nil, ""},
	{"deal", "company", "Company", "lookup", false, false, true, nil, "company"},
	{"deal", "contact_name", "Primary Contact", "text", false, true, false, nil, ""},
	{"deal", "source", "Source", "dropdown", false, false, true, optLeadSources, ""},
	{"deal", "priority", "Priority", "dropdown", false, false, true, optPriorities, ""},
	{"deal", "city", "City", "text", false, true, true, nil, ""},
	{"deal", "country", "Country", "text", false, true, true, nil, ""},
	{"deal", "expected_revenue", "Expected Revenue", "currency", false, false, true, nil, ""},
	{"deal", "tags", "Tags", "multiselect", false, false, true, optTagPool, ""},
	{"deal", "description", "Description", "textarea", false, false, false, nil, ""},

	// lead (15+)
	{"lead", "first_name", "First Name", "text", true, true, false, nil, ""},
	{"lead", "last_name", "Last Name", "text", true, true, false, nil, ""},
	{"lead", "email", "Email", "email", false, true, true, nil, ""},
	{"lead", "phone", "Phone", "phone", false, true, false, nil, ""},
	{"lead", "company_name", "Company Name", "text", false, true, true, nil, ""},
	{"lead", "company", "Company", "lookup", false, false, true, nil, "company"},
	{"lead", "job_title", "Job Title", "text", false, true, true, nil, ""},
	{"lead", "industry", "Industry", "dropdown", false, false, true, optIndustries, ""},
	{"lead", "status", "Status", "dropdown", true, false, true, optLeadStatuses, ""},
	{"lead", "source", "Source", "dropdown", false, false, true, optLeadSources, ""},
	{"lead", "website", "Website", "url", false, false, false, nil, ""},
	{"lead", "employees", "Employees", "number", false, false, true, nil, ""},
	{"lead", "annual_revenue", "Annual Revenue", "currency", false, false, true, nil, ""},
	{"lead", "city", "City", "text", false, true, true, nil, ""},
	{"lead", "country", "Country", "text", false, true, true, nil, ""},
	{"lead", "priority", "Priority", "dropdown", false, false, true, optPriorities, ""},
	{"lead", "rating", "Rating", "dropdown", false, false, true, optRatings, ""},
	{"lead", "last_contacted", "Last Contacted", "date", false, false, true, nil, ""},
	{"lead", "next_follow_up", "Next Follow-up", "date", false, false, true, nil, ""},
	{"lead", "linkedin", "LinkedIn", "url", false, false, false, nil, ""},
	{"lead", "tags", "Tags", "multiselect", false, false, true, optTagPool, ""},
	{"lead", "notes", "Notes", "textarea", false, false, false, nil, ""},

	// task (15+)
	{"task", "title", "Task Title", "text", true, true, false, nil, ""},
	{"task", "status", "Status", "dropdown", true, false, true, optTaskStatuses, ""},
	{"task", "priority", "Priority", "dropdown", false, false, true, optTaskPriorities, ""},
	{"task", "due_date", "Due Date", "date", false, false, true, nil, ""},
	{"task", "company", "Company", "lookup", false, false, true, nil, "company"},
	{"task", "related_lead", "Related Lead", "text", false, true, false, nil, ""},
	{"task", "related_contact", "Related Contact", "text", false, true, false, nil, ""},
	{"task", "city", "City", "text", false, true, true, nil, ""},
	{"task", "country", "Country", "text", false, true, true, nil, ""},
	{"task", "estimated_hours", "Estimated Hours", "number", false, false, true, nil, ""},
	{"task", "completed_at", "Completed At", "date", false, false, true, nil, ""},
	{"task", "reminder_date", "Reminder Date", "date", false, false, true, nil, ""},
	{"task", "tags", "Tags", "multiselect", false, false, true, optTagPool, ""},
	{"task", "outcome", "Outcome", "text", false, false, false, nil, ""},
	{"task", "description", "Description", "textarea", false, false, false, nil, ""},
}

type layoutSection struct {
	Key     string   `json:"key"`
	Label   string   `json:"label"`
	Order   int      `json:"order"`
	Columns int      `json:"columns"`
	Fields  []string `json:"fields"`
}

var defaultFormLayouts = map[string][]layoutSection{
	"company": {
		{Key: "basic", Label: "Basic Information", Order: 1, Columns: 2, Fields: []string{"name", "website", "industry", "status", "priority"}},
		{Key: "contact", Label: "Contact Information", Order: 2, Columns: 2, Fields: []string{"phone", "email", "linkedin"}},
		{Key: "company", Label: "Company Details", Order: 3, Columns: 2, Fields: []string{"city", "country", "employees", "annual_revenue"}},
		{Key: "sales", Label: "Sales Information", Order: 4, Columns: 2, Fields: []string{"last_contacted", "next_follow_up", "tags"}},
		{Key: "additional", Label: "Additional Details", Order: 5, Columns: 1, Fields: []string{"description"}},
	},
	"contact": {
		{Key: "basic", Label: "Basic Information", Order: 1, Columns: 2, Fields: []string{"first_name", "last_name", "job_title", "department"}},
		{Key: "contact", Label: "Contact Information", Order: 2, Columns: 2, Fields: []string{"email", "phone", "mobile", "linkedin"}},
		{Key: "company", Label: "Company Details", Order: 3, Columns: 2, Fields: []string{"company", "city", "country"}},
		{Key: "sales", Label: "Sales Information", Order: 4, Columns: 2, Fields: []string{"priority", "rating", "last_contacted", "next_follow_up", "tags"}},
		{Key: "additional", Label: "Additional Details", Order: 5, Columns: 1, Fields: []string{"notes"}},
	},
	"deal": {
		{Key: "basic", Label: "Basic Information", Order: 1, Columns: 2, Fields: []string{"title", "stage", "amount", "probability"}},
		{Key: "sales", Label: "Sales Information", Order: 2, Columns: 2, Fields: []string{"close_date", "next_step", "source", "priority", "expected_revenue"}},
		{Key: "company", Label: "Company Details", Order: 3, Columns: 2, Fields: []string{"company", "contact_name", "city", "country"}},
		{Key: "additional", Label: "Additional Details", Order: 4, Columns: 2, Fields: []string{"tags", "description"}},
	},
	"lead": {
		{Key: "basic", Label: "Basic Information", Order: 1, Columns: 2, Fields: []string{"first_name", "last_name", "job_title", "status"}},
		{Key: "contact", Label: "Contact Information", Order: 2, Columns: 2, Fields: []string{"email", "phone", "linkedin", "website"}},
		{Key: "company", Label: "Company Details", Order: 3, Columns: 2, Fields: []string{"company_name", "company", "industry", "employees", "annual_revenue", "city", "country"}},
		{Key: "sales", Label: "Sales Information", Order: 4, Columns: 2, Fields: []string{"source", "priority", "rating", "last_contacted", "next_follow_up", "tags"}},
		{Key: "additional", Label: "Additional Details", Order: 5, Columns: 1, Fields: []string{"notes"}},
	},
	"task": {
		{Key: "basic", Label: "Basic Information", Order: 1, Columns: 2, Fields: []string{"title", "status", "priority", "due_date"}},
		{Key: "relationships", Label: "Relationships", Order: 2, Columns: 2, Fields: []string{"company", "related_lead", "related_contact"}},
		{Key: "scheduling", Label: "Scheduling", Order: 3, Columns: 2, Fields: []string{"reminder_date", "estimated_hours", "completed_at", "city", "country"}},
		{Key: "additional", Label: "Additional Details", Order: 4, Columns: 2, Fields: []string{"tags", "outcome", "description"}},
	},
}

var defaultListVisible = map[string][]string{
	"company": {"name", "industry", "status", "city", "phone", "priority"},
	"contact": {"first_name", "last_name", "email", "phone", "company", "job_title"},
	"deal":    {"title", "amount", "stage", "close_date", "company", "priority"},
	"lead":    {"first_name", "last_name", "email", "status", "company_name", "phone", "priority"},
	"task":    {"title", "status", "priority", "due_date", "company"},
}

func (s *Service) seedDefaultModules(ctx context.Context, orgID string) error {
	moduleIDs := map[string]string{}
	for _, m := range defaultModuleCatalog {
		var id string
		err := s.db.QueryRow(ctx, `
			INSERT INTO modules (
				organization_id, api_name, singular_label, plural_label,
				icon, color, storage_strategy, is_system, sort_order,
				is_enabled, is_visible_sidebar
			) VALUES ($1,$2,$3,$4,$5,$6,'dynamic',TRUE,$7,TRUE,TRUE)
			ON CONFLICT (organization_id, api_name) DO UPDATE
			SET singular_label = EXCLUDED.singular_label,
			    plural_label = EXCLUDED.plural_label,
			    icon = EXCLUDED.icon,
			    color = EXCLUDED.color,
			    sort_order = EXCLUDED.sort_order,
			    is_enabled = TRUE,
			    is_visible_sidebar = TRUE
			RETURNING id
		`, orgID, m.API, m.Singular, m.Plural, m.Icon, m.Color, m.Sort).Scan(&id)
		if err != nil {
			return fmt.Errorf("seed module %s: %w", m.API, err)
		}
		moduleIDs[m.API] = id
	}

	sortByModule := map[string]int{}
	for _, f := range defaultFieldCatalog {
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
		_, err = s.db.Exec(ctx, `
			INSERT INTO fields (
				organization_id, module_id, api_name, label, field_type,
				is_required, is_searchable, is_filterable, options,
				lookup_module_id, sort_order, is_system, is_visible,
				lock_mode, editable_by, viewable_by
			) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9::jsonb,$10,$11,TRUE,TRUE,'never','ALL','ALL')
			ON CONFLICT (module_id, api_name) DO UPDATE
			SET label = EXCLUDED.label,
			    field_type = EXCLUDED.field_type,
			    is_required = EXCLUDED.is_required,
			    is_searchable = EXCLUDED.is_searchable,
			    is_filterable = EXCLUDED.is_filterable,
			    options = EXCLUDED.options,
			    lookup_module_id = EXCLUDED.lookup_module_id,
			    sort_order = EXCLUDED.sort_order
		`, orgID, moduleIDs[f.Module], f.API, f.Label, f.Type,
			f.Required, f.Searchable, f.Filterable, string(optionsJSON),
			lookupID, sortByModule[f.Module])
		if err != nil {
			return fmt.Errorf("seed field %s.%s: %w", f.Module, f.API, err)
		}
	}

	for api, sections := range defaultFormLayouts {
		formRaw, err := json.Marshal(map[string]any{"sections": sections})
		if err != nil {
			return err
		}
		_, err = s.db.Exec(ctx, `
			INSERT INTO layouts (organization_id, module_id, name, layout_type, is_default, config)
			VALUES ($1, $2, 'Default Form', 'form', TRUE, $3)
			ON CONFLICT DO NOTHING
		`, orgID, moduleIDs[api], formRaw)
		if err != nil {
			// layouts may lack unique constraint — fall back to skip-if-exists
			var exists bool
			_ = s.db.QueryRow(ctx, `
				SELECT EXISTS(
					SELECT 1 FROM layouts
					WHERE organization_id = $1 AND module_id = $2 AND layout_type = 'form' AND is_default
				)
			`, orgID, moduleIDs[api]).Scan(&exists)
			if !exists {
				_, err = s.db.Exec(ctx, `
					INSERT INTO layouts (organization_id, module_id, name, layout_type, is_default, config)
					VALUES ($1, $2, 'Default Form', 'form', TRUE, $3)
				`, orgID, moduleIDs[api], formRaw)
				if err != nil {
					return fmt.Errorf("seed form layout %s: %w", api, err)
				}
			}
		}

		visibleSet := map[string]bool{}
		for _, k := range defaultListVisible[api] {
			visibleSet[k] = true
		}
		cols := make([]map[string]any, 0)
		order := 0
		seen := map[string]bool{}
		for _, sec := range sections {
			for _, key := range sec.Fields {
				if seen[key] {
					continue
				}
				seen[key] = true
				order++
				cols = append(cols, map[string]any{
					"field_key": key, "visible": visibleSet[key], "order": order,
					"sortable": true, "searchable": false, "system": false,
				})
			}
		}
		cols = append(cols, map[string]any{
			"field_key": "_actions", "visible": true, "order": order + 1,
			"sortable": false, "searchable": false, "system": true,
		})
		listRaw, err := json.Marshal(map[string]any{"columns": cols})
		if err != nil {
			return err
		}
		var listExists bool
		_ = s.db.QueryRow(ctx, `
			SELECT EXISTS(
				SELECT 1 FROM layouts
				WHERE organization_id = $1 AND module_id = $2 AND layout_type = 'list' AND is_default
			)
		`, orgID, moduleIDs[api]).Scan(&listExists)
		if !listExists {
			_, err = s.db.Exec(ctx, `
				INSERT INTO layouts (organization_id, module_id, name, layout_type, is_default, config)
				VALUES ($1, $2, 'Default List', 'list', TRUE, $3)
			`, orgID, moduleIDs[api], listRaw)
			if err != nil {
				return fmt.Errorf("seed list layout %s: %w", api, err)
			}
		}

		var detailExists bool
		_ = s.db.QueryRow(ctx, `
			SELECT EXISTS(
				SELECT 1 FROM layouts
				WHERE organization_id = $1 AND module_id = $2 AND layout_type = 'detail' AND is_default
			)
		`, orgID, moduleIDs[api]).Scan(&detailExists)
		if !detailExists {
			detailRaw, err := json.Marshal(map[string]any{
				"sections": sections,
				"tabs":     []string{"overview", "notes", "attachments", "timeline", "related"},
			})
			if err != nil {
				return err
			}
			_, err = s.db.Exec(ctx, `
				INSERT INTO layouts (organization_id, module_id, name, layout_type, is_default, config)
				VALUES ($1, $2, 'Default Detail', 'detail', TRUE, $3)
			`, orgID, moduleIDs[api], detailRaw)
			if err != nil {
				return fmt.Errorf("seed detail layout %s: %w", api, err)
			}
		}
	}
	return nil
}
