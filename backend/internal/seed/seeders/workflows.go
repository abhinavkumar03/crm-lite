package seeders

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"

	workflowrepo "github.com/abhinavkumar03/crm-lite/backend/internal/workflow/repository"
	workflowservice "github.com/abhinavkumar03/crm-lite/backend/internal/workflow/service"
)

// WorkflowDemoSeeder seeds built-in templates, sample workflows, and execution history.
type WorkflowDemoSeeder struct{}

func (WorkflowDemoSeeder) Name() string { return "workflow_demo" }

func (WorkflowDemoSeeder) Run(ctx context.Context, db *pgxpool.Pool) error {
	if err := seedBuiltinWorkflowTemplates(ctx, db); err != nil {
		return err
	}
	orgIDs, err := listDemoOrgIDs(ctx, db)
	if err != nil {
		return err
	}
	for _, orgID := range orgIDs {
		if err := seedOrgWorkflows(ctx, db, orgID); err != nil {
			return err
		}
	}
	return nil
}

func seedBuiltinWorkflowTemplates(ctx context.Context, db *pgxpool.Pool) error {
	repo := workflowrepo.New(db)
	svc := workflowservice.New(repo, nil, nil, true)
	return svc.EnsureBuiltinTemplates(ctx)
}

func seedOrgWorkflows(ctx context.Context, db *pgxpool.Pool, orgID string) error {
	var userID string
	_ = db.QueryRow(ctx, `
		SELECT u.id FROM users u
		JOIN organization_members om ON om.user_id = u.id
		WHERE om.organization_id = $1 AND om.status = 'active'
		ORDER BY u.created_at LIMIT 1`, orgID).Scan(&userID)

	modules := map[string]string{}
	rows, err := db.Query(ctx, `SELECT api_name, id FROM modules WHERE organization_id = $1`, orgID)
	if err != nil {
		return err
	}
	defer rows.Close()
	for rows.Next() {
		var api, id string
		if err := rows.Scan(&api, &id); err != nil {
			return err
		}
		modules[api] = id
	}

	type wfSpec struct {
		name, desc, module, status, trigger string
		field                               string
		condField, condOp, condVal          string
		actions                             []struct{ typ string; cfg map[string]any }
	}

	specs := []wfSpec{
		{
			name: "Manual Follow-up", desc: "Manual trigger demo for record detail Run", module: "lead", status: "active",
			trigger: "manual",
			actions: []struct {
				typ string
				cfg map[string]any
			}{
				{"create_activity", map[string]any{"description": "Manual workflow run"}},
				{"create_note", map[string]any{"body": "Manual follow-up note from workflow"}},
			},
		},
		{
			name: "New Website Lead", desc: "Active demo: website lead create", module: "lead", status: "active",
			trigger: "record_created", condField: "source", condOp: "eq", condVal: "Website",
			actions: []struct {
				typ string
				cfg map[string]any
			}{
				{"create_activity", map[string]any{"description": "Website lead automation"}},
				{"create_note", map[string]any{"body": "Auto-note: website lead"}},
			},
		},
		{
			name: "Qualified Lead", desc: "Active demo: status → Qualified", module: "lead", status: "active",
			trigger: "field_updated", field: "status", condField: "status", condOp: "eq", condVal: "Qualified",
			actions: []struct {
				typ string
				cfg map[string]any
			}{
				{"create_activity", map[string]any{"description": "Qualified lead workflow"}},
				{"send_email", map[string]any{"subject": "Welcome", "body": "Hello {{lead.name}}"}},
			},
		},
		{
			name: "Lost Lead", desc: "Disabled demo workflow", module: "lead", status: "disabled",
			trigger: "field_updated", field: "status", condField: "status", condOp: "eq", condVal: "Lost",
			actions: []struct {
				typ string
				cfg map[string]any
			}{
				{"create_activity", map[string]any{"description": "Lost lead logged"}},
			},
		},
		{
			name: "High Value Lead", desc: "Draft demo workflow", module: "lead", status: "draft",
			trigger: "record_created",
			actions: []struct {
				typ string
				cfg map[string]any
			}{
				{"create_activity", map[string]any{"description": "High value lead draft"}},
			},
		},
		{
			name: "New Contact Welcome", desc: "Active contact create", module: "contact", status: "active",
			trigger: "record_created",
			actions: []struct {
				typ string
				cfg map[string]any
			}{
				{"create_note", map[string]any{"body": "Welcome contact note"}},
			},
		},
		{
			name: "Contact Updated", desc: "Draft contact update", module: "contact", status: "draft",
			trigger: "record_updated",
			actions: []struct {
				typ string
				cfg map[string]any
			}{
				{"create_activity", map[string]any{"description": "Contact updated"}},
			},
		},
		{
			name: "Task Created", desc: "Active task create", module: "task", status: "active",
			trigger: "record_created",
			actions: []struct {
				typ string
				cfg map[string]any
			}{
				{"create_activity", map[string]any{"description": "Task created automation"}},
			},
		},
		{
			name: "Task Completed", desc: "Active when status Complete", module: "task", status: "active",
			trigger: "field_updated", field: "status", condField: "status", condOp: "eq", condVal: "Completed",
			actions: []struct {
				typ string
				cfg map[string]any
			}{
				{"create_activity", map[string]any{"description": "Task completed"}},
			},
		},
		{
			name: "Overdue Task Sweep", desc: "Date-based overdue reminder", module: "task", status: "active",
			trigger: "date_based", field: "due_date",
			actions: []struct {
				typ string
				cfg map[string]any
			}{
				{"create_activity", map[string]any{"description": "Overdue task reminder"}},
			},
		},
	}

	var firstWF, firstVersion, firstRecord, firstModule string

	for _, s := range specs {
		mid, ok := modules[s.module]
		if !ok {
			continue
		}
		var existing string
		_ = db.QueryRow(ctx, `
			SELECT id FROM workflows WHERE organization_id = $1 AND name = $2 LIMIT 1`, orgID, s.name).Scan(&existing)
		if existing != "" {
			if firstWF == "" && s.status == "active" {
				firstWF = existing
				firstModule = mid
				_ = db.QueryRow(ctx, `
					SELECT published_version_id FROM workflows WHERE id = $1`, existing).Scan(&firstVersion)
				_ = db.QueryRow(ctx, `
					SELECT id FROM records WHERE organization_id = $1 AND module_id = $2 LIMIT 1`, orgID, mid).Scan(&firstRecord)
			}
			continue
		}

		var wfID string
		err := db.QueryRow(ctx, `
			INSERT INTO workflows (
				organization_id, module_id, name, description, status, on_action_error, priority, created_by, updated_by
			) VALUES ($1,$2,$3,$4,$5,'stop',100,$6,$6)
			RETURNING id`, orgID, mid, s.name, s.desc, s.status, nullStr(userID)).Scan(&wfID)
		if err != nil {
			return fmt.Errorf("seed workflow %s: %w", s.name, err)
		}

		state := "draft"
		var publishedAt *time.Time
		if s.status == "active" || s.status == "disabled" {
			state = "published"
			publishedAt = ptrTime(time.Now().UTC())
		}
		var versionID string
		err = db.QueryRow(ctx, `
			INSERT INTO workflow_versions (workflow_id, organization_id, version, state, definition_snapshot, changelog, published_at, published_by)
			VALUES ($1,$2,1,$3,'{}'::jsonb,$4,$5,$6)
			RETURNING id`,
			wfID, orgID, state, "Demo seed", publishedAt, nullStr(userID),
		).Scan(&versionID)
		if err != nil {
			return err
		}

		cfg := map[string]any{}
		if s.field != "" {
			if s.trigger == "field_updated" {
				cfg["field_api_name"] = s.field
			}
			if s.trigger == "date_based" {
				cfg["field_api_name"] = s.field
				cfg["offset_days"] = 0
			}
		}
		cfgJSON, _ := json.Marshal(cfg)
		_, err = db.Exec(ctx, `
			INSERT INTO workflow_triggers (version_id, organization_id, type, config)
			VALUES ($1,$2,$3,$4)`, versionID, orgID, s.trigger, cfgJSON)
		if err != nil {
			return err
		}

		if s.condField != "" {
			var rootID string
			err = db.QueryRow(ctx, `
				INSERT INTO workflow_conditions (version_id, organization_id, parent_id, node_type, logic, sort_order)
				VALUES ($1,$2,NULL,'group','and',0) RETURNING id`, versionID, orgID).Scan(&rootID)
			if err != nil {
				return err
			}
			val, _ := json.Marshal(s.condVal)
			_, err = db.Exec(ctx, `
				INSERT INTO workflow_conditions (
					version_id, organization_id, parent_id, node_type, field_api_name, operator, value, sort_order
				) VALUES ($1,$2,$3,'predicate',$4,$5,$6,0)`,
				versionID, orgID, rootID, s.condField, s.condOp, val)
			if err != nil {
				return err
			}
		}

		for i, a := range s.actions {
			acfg, _ := json.Marshal(a.cfg)
			_, err = db.Exec(ctx, `
				INSERT INTO workflow_actions (version_id, organization_id, sort_order, type, config)
				VALUES ($1,$2,$3,$4,$5)`, versionID, orgID, i, a.typ, acfg)
			if err != nil {
				return err
			}
		}

		if s.status == "active" || s.status == "disabled" {
			_, _ = db.Exec(ctx, `UPDATE workflows SET published_version_id = $1 WHERE id = $2`, versionID, wfID)
		}

		if firstWF == "" && s.status == "active" {
			firstWF = wfID
			firstVersion = versionID
			firstModule = mid
			_ = db.QueryRow(ctx, `
				SELECT id FROM records WHERE organization_id = $1 AND module_id = $2 LIMIT 1`, orgID, mid).Scan(&firstRecord)
		}
	}

	// Sample executions (success + failed).
	if firstWF != "" && firstRecord != "" {
		now := time.Now().UTC()
		started := now.Add(-10 * time.Minute)
		finished := now.Add(-9 * time.Minute)
		dur := 60000
		var execOK, execFail string
		_ = db.QueryRow(ctx, `
			INSERT INTO workflow_executions (
				organization_id, workflow_id, version_id, module_id, record_id,
				trigger_type, status, source, depth, started_at, finished_at, duration_ms
			) VALUES ($1,$2,$3,$4,$5,'record_created','succeeded','user',0,$6,$7,$8)
			RETURNING id`, orgID, firstWF, firstVersion, firstModule, firstRecord, started, finished, dur).Scan(&execOK)
		failMsg := "simulated email failure"
		_ = db.QueryRow(ctx, `
			INSERT INTO workflow_executions (
				organization_id, workflow_id, version_id, module_id, record_id,
				trigger_type, status, source, depth, error_summary, started_at, finished_at, duration_ms
			) VALUES ($1,$2,$3,$4,$5,'field_updated','failed','user',0,$6,$7,$8,$9)
			RETURNING id`, orgID, firstWF, firstVersion, firstModule, firstRecord, failMsg,
			started.Add(-time.Hour), finished.Add(-time.Hour), 1200).Scan(&execFail)

		if execOK != "" {
			_, _ = db.Exec(ctx, `
				INSERT INTO workflow_execution_steps (
					execution_id, organization_id, sort_order, action_type, status, input, output, started_at, finished_at
				) VALUES ($1,$2,0,'create_activity','succeeded','{}','{"ok":true}',$3,$4)`,
				execOK, orgID, started, finished)
		}
		if execFail != "" {
			_, _ = db.Exec(ctx, `
				INSERT INTO workflow_execution_steps (
					execution_id, organization_id, sort_order, action_type, status, input, output, error, started_at, finished_at
				) VALUES ($1,$2,0,'send_email','failed','{}','{}',$3,$4,$5)`,
				execFail, orgID, failMsg, started.Add(-time.Hour), finished.Add(-time.Hour))
		}
	}

	return nil
}

func nullStr(s string) *string {
	if s == "" {
		return nil
	}
	return &s
}

func ptrTime(t time.Time) *time.Time { return &t }
