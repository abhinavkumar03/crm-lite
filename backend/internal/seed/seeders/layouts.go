package seeders

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

// LayoutsSeeder upserts detail, form, and list layouts with meaningful defaults.
type LayoutsSeeder struct{}

func (LayoutsSeeder) Name() string { return "layouts" }

type layoutSectionDef struct {
	Key         string   `json:"key"`
	Label       string   `json:"label"`
	Description string   `json:"description,omitempty"`
	Order       int      `json:"order,omitempty"`
	Collapsed   bool     `json:"collapsed,omitempty"`
	Columns     int      `json:"columns,omitempty"`
	Fields      []string `json:"fields"`
}

type layoutConfigDef struct {
	Sections []layoutSectionDef `json:"sections"`
	Tabs     []string           `json:"tabs,omitempty"`
}

type listColumnDef struct {
	FieldKey   string `json:"field_key"`
	Visible    bool   `json:"visible"`
	Order      int    `json:"order"`
	Sortable   bool   `json:"sortable"`
	Searchable bool   `json:"searchable"`
	System     bool   `json:"system"`
}

type listConfigDef struct {
	Columns []listColumnDef `json:"columns"`
}

var defaultTabs = []string{"overview", "notes", "attachments", "timeline", "related"}

var systemSection = layoutSectionDef{
	Key:    "system",
	Label:  "System Fields",
	Fields: []string{"owner_id", "assigned_to", "visibility", "created_at", "updated_at"},
}

var moduleLayouts = map[string][]layoutSectionDef{
	"company": {
		{Key: "basic", Label: "Basic Information", Order: 1, Columns: 2, Fields: []string{"name", "website", "industry", "status", "priority"}},
		{Key: "contact", Label: "Contact Information", Order: 2, Columns: 2, Fields: []string{"phone", "email", "linkedin"}},
		{Key: "company", Label: "Company Details", Order: 3, Columns: 2, Fields: []string{"city", "country", "employees", "annual_revenue"}},
		{Key: "sales", Label: "Sales Information", Order: 4, Columns: 2, Fields: []string{"last_contacted", "next_follow_up", "tags"}},
		{Key: "additional", Label: "Additional Details", Order: 5, Columns: 1, Fields: []string{"description"}},
		systemSection,
	},
	"contact": {
		{Key: "basic", Label: "Basic Information", Order: 1, Columns: 2, Fields: []string{"first_name", "last_name", "job_title", "department"}},
		{Key: "contact", Label: "Contact Information", Order: 2, Columns: 2, Fields: []string{"email", "phone", "mobile", "linkedin"}},
		{Key: "company", Label: "Company Details", Order: 3, Columns: 2, Fields: []string{"company", "city", "country"}},
		{Key: "sales", Label: "Sales Information", Order: 4, Columns: 2, Fields: []string{"priority", "rating", "last_contacted", "next_follow_up", "tags"}},
		{Key: "additional", Label: "Additional Details", Order: 5, Columns: 1, Fields: []string{"notes"}},
		systemSection,
	},
	"deal": {
		{Key: "basic", Label: "Basic Information", Order: 1, Columns: 2, Fields: []string{"title", "stage", "amount", "probability"}},
		{Key: "sales", Label: "Sales Information", Order: 2, Columns: 2, Fields: []string{"close_date", "next_step", "source", "priority", "expected_revenue"}},
		{Key: "company", Label: "Company Details", Order: 3, Columns: 2, Fields: []string{"company", "contact_name", "city", "country"}},
		{Key: "additional", Label: "Additional Details", Order: 4, Columns: 2, Fields: []string{"tags", "description"}},
		systemSection,
	},
	"lead": {
		{Key: "basic", Label: "Basic Information", Order: 1, Columns: 2, Fields: []string{"first_name", "last_name", "job_title", "status"}},
		{Key: "contact", Label: "Contact Information", Order: 2, Columns: 2, Fields: []string{"email", "phone", "linkedin", "website"}},
		{Key: "company", Label: "Company Details", Order: 3, Columns: 2, Fields: []string{"company_name", "company", "industry", "employees", "annual_revenue", "city", "country"}},
		{Key: "sales", Label: "Sales Information", Order: 4, Columns: 2, Fields: []string{"source", "priority", "rating", "last_contacted", "next_follow_up", "tags"}},
		{Key: "additional", Label: "Additional Details", Order: 5, Columns: 1, Fields: []string{"notes"}},
		systemSection,
	},
	"task": {
		{Key: "basic", Label: "Basic Information", Order: 1, Columns: 2, Fields: []string{"title", "status", "priority", "due_date"}},
		{Key: "relationships", Label: "Relationships", Order: 2, Columns: 2, Fields: []string{"company", "related_lead", "related_contact"}},
		{Key: "scheduling", Label: "Scheduling", Order: 3, Columns: 2, Fields: []string{"reminder_date", "estimated_hours", "completed_at", "city", "country"}},
		{Key: "additional", Label: "Additional Details", Order: 4, Columns: 2, Fields: []string{"tags", "outcome", "description"}},
		systemSection,
	},
}

// moduleListVisibleKeys controls which fields are Visible in the default list layout.
// Fields present in form sections but missing here are seeded as hidden columns.
var moduleListVisibleKeys = map[string][]string{
	"company": {"name", "industry", "status", "city", "phone"},
	"contact": {"first_name", "last_name", "email", "phone", "company"},
	"deal":    {"title", "amount", "stage", "close_date", "company"},
	"lead":    {"first_name", "last_name", "email", "status", "company", "phone"},
	"task":    {"title", "status", "priority", "due_date", "company"},
}

func (LayoutsSeeder) Run(ctx context.Context, db *pgxpool.Pool) error {
	orgIDs, err := listDemoOrgIDs(ctx, db)
	if err != nil {
		return err
	}

	for _, orgID := range orgIDs {
		for apiName, sections := range moduleLayouts {
			moduleID, err := getModuleID(ctx, db, orgID, apiName)
			if err != nil {
				return err
			}
			if err := upsertLayout(ctx, db, orgID, moduleID, apiName, "detail", "Default Detail", layoutConfigDef{
				Sections: sections, Tabs: defaultTabs,
			}); err != nil {
				return err
			}

			formSections := make([]layoutSectionDef, 0, len(sections))
			for _, sec := range sections {
				if sec.Key == "system" {
					continue
				}
				formSections = append(formSections, sec)
			}
			if err := upsertLayout(ctx, db, orgID, moduleID, apiName, "form", "Default Form", layoutConfigDef{
				Sections: formSections,
			}); err != nil {
				return err
			}

			listCfg := buildListConfig(apiName, formSections)
			if err := upsertListLayout(ctx, db, orgID, moduleID, apiName, listCfg); err != nil {
				return err
			}
		}
	}
	return nil
}

func buildListConfig(apiName string, sections []layoutSectionDef) listConfigDef {
	visibleSet := map[string]bool{}
	if keys, ok := moduleListVisibleKeys[apiName]; ok {
		for _, k := range keys {
			visibleSet[k] = true
		}
	}

	// Prefer curated listing order when defined.
	cols := make([]listColumnDef, 0)
	order := 0
	seen := map[string]bool{}
	if keys, ok := moduleListVisibleKeys[apiName]; ok {
		for _, key := range keys {
			if seen[key] {
				continue
			}
			seen[key] = true
			order++
			cols = append(cols, listColumnDef{
				FieldKey: key, Visible: true, Order: order,
				Sortable: true, Searchable: false, System: false,
			})
		}
	}

	for _, sec := range sections {
		for _, key := range sec.Fields {
			if seen[key] {
				continue
			}
			seen[key] = true
			order++
			visible := true
			if len(visibleSet) > 0 {
				visible = visibleSet[key]
			}
			cols = append(cols, listColumnDef{
				FieldKey: key, Visible: visible, Order: order,
				Sortable: true, Searchable: false, System: false,
			})
		}
	}
	cols = append(cols, listColumnDef{
		FieldKey: "_actions", Visible: true, Order: order + 1,
		Sortable: false, Searchable: false, System: true,
	})
	return listConfigDef{Columns: cols}
}

func upsertLayout(ctx context.Context, db *pgxpool.Pool, orgID, moduleID, apiName, layoutType, name string, cfg layoutConfigDef) error {
	raw, err := json.Marshal(cfg)
	if err != nil {
		return err
	}
	var existingID string
	err = db.QueryRow(ctx, `
		SELECT id::text FROM layouts
		WHERE organization_id = $1 AND module_id = $2
		  AND layout_type = $3 AND is_default = TRUE
		LIMIT 1
	`, orgID, moduleID, layoutType).Scan(&existingID)
	if err == nil && existingID != "" {
		_, err = db.Exec(ctx, `
			UPDATE layouts SET config = $2, updated_at = NOW() WHERE id = $1
		`, existingID, raw)
		if err != nil {
			return fmt.Errorf("update %s layout %s: %w", layoutType, apiName, err)
		}
		return nil
	}
	_, err = db.Exec(ctx, `
		INSERT INTO layouts (organization_id, module_id, name, layout_type, is_default, config)
		VALUES ($1, $2, $3, $4, TRUE, $5)
	`, orgID, moduleID, name, layoutType, raw)
	if err != nil {
		return fmt.Errorf("insert %s layout %s: %w", layoutType, apiName, err)
	}
	return nil
}

func upsertListLayout(ctx context.Context, db *pgxpool.Pool, orgID, moduleID, apiName string, cfg listConfigDef) error {
	raw, err := json.Marshal(cfg)
	if err != nil {
		return err
	}
	var existingID string
	err = db.QueryRow(ctx, `
		SELECT id::text FROM layouts
		WHERE organization_id = $1 AND module_id = $2
		  AND layout_type = 'list' AND is_default = TRUE
		LIMIT 1
	`, orgID, moduleID).Scan(&existingID)
	if err == nil && existingID != "" {
		_, err = db.Exec(ctx, `
			UPDATE layouts SET config = $2, updated_at = NOW() WHERE id = $1
		`, existingID, raw)
		if err != nil {
			return fmt.Errorf("update list layout %s: %w", apiName, err)
		}
		return nil
	}
	_, err = db.Exec(ctx, `
		INSERT INTO layouts (organization_id, module_id, name, layout_type, is_default, config)
		VALUES ($1, $2, 'Default List', 'list', TRUE, $3)
	`, orgID, moduleID, raw)
	if err != nil {
		return fmt.Errorf("insert list layout %s: %w", apiName, err)
	}
	return nil
}
