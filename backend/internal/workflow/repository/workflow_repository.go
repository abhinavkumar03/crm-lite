package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/abhinavkumar03/crm-lite/backend/internal/workflow/entity"
)

type Repository struct {
	db *pgxpool.Pool
}

func New(db *pgxpool.Pool) *Repository {
	return &Repository{db: db}
}

func (r *Repository) DB() *pgxpool.Pool { return r.db }

// --- Workflows ---

func (r *Repository) CreateWorkflow(ctx context.Context, w *entity.Workflow) error {
	return r.db.QueryRow(ctx, `
		INSERT INTO workflows (
			organization_id, module_id, name, description, status,
			on_action_error, priority, created_by, updated_by
		) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9)
		RETURNING id, created_at, updated_at`,
		w.OrganizationID, w.ModuleID, w.Name, w.Description, w.Status,
		w.OnActionError, w.Priority, w.CreatedBy, w.UpdatedBy,
	).Scan(&w.ID, &w.CreatedAt, &w.UpdatedAt)
}

func (r *Repository) UpdateWorkflowHeader(ctx context.Context, w *entity.Workflow) error {
	_, err := r.db.Exec(ctx, `
		UPDATE workflows SET
			module_id = $3, name = $4, description = $5, status = $6,
			on_action_error = $7, priority = $8, published_version_id = $9,
			updated_by = $10, updated_at = NOW()
		WHERE id = $1 AND organization_id = $2`,
		w.ID, w.OrganizationID, w.ModuleID, w.Name, w.Description, w.Status,
		w.OnActionError, w.Priority, w.PublishedVersionID, w.UpdatedBy,
	)
	return err
}

func (r *Repository) GetWorkflow(ctx context.Context, orgID, id string) (*entity.Workflow, error) {
	row := r.db.QueryRow(ctx, `
		SELECT id, organization_id, module_id, name, description, status,
			on_action_error, priority, published_version_id, created_by, updated_by,
			created_at, updated_at
		FROM workflows WHERE id = $1 AND organization_id = $2`, id, orgID)
	return scanWorkflow(row)
}

func (r *Repository) ListWorkflows(ctx context.Context, orgID string, status string, moduleID string, page, pageSize int) ([]entity.Workflow, int, error) {
	args := []any{orgID}
	where := []string{"organization_id = $1"}
	if status != "" {
		args = append(args, status)
		where = append(where, fmt.Sprintf("status = $%d", len(args)))
	}
	if moduleID != "" {
		args = append(args, moduleID)
		where = append(where, fmt.Sprintf("module_id = $%d", len(args)))
	}
	wSQL := strings.Join(where, " AND ")

	var total int
	if err := r.db.QueryRow(ctx, `SELECT COUNT(*) FROM workflows WHERE `+wSQL, args...).Scan(&total); err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	args = append(args, pageSize, offset)
	rows, err := r.db.Query(ctx, `
		SELECT id, organization_id, module_id, name, description, status,
			on_action_error, priority, published_version_id, created_by, updated_by,
			created_at, updated_at
		FROM workflows WHERE `+wSQL+`
		ORDER BY priority ASC, updated_at DESC
		LIMIT $`+fmt.Sprint(len(args)-1)+` OFFSET $`+fmt.Sprint(len(args)), args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var out []entity.Workflow
	for rows.Next() {
		w, err := scanWorkflowRows(rows)
		if err != nil {
			return nil, 0, err
		}
		out = append(out, *w)
	}
	return out, total, rows.Err()
}

func (r *Repository) ModuleAPIName(ctx context.Context, orgID, moduleID string) (string, error) {
	var name string
	err := r.db.QueryRow(ctx, `
		SELECT api_name FROM modules WHERE id = $1 AND organization_id = $2`, moduleID, orgID).Scan(&name)
	return name, err
}

func (r *Repository) ModuleIDByAPIName(ctx context.Context, orgID, apiName string) (string, error) {
	var id string
	err := r.db.QueryRow(ctx, `
		SELECT id FROM modules WHERE organization_id = $1 AND api_name = $2`, orgID, apiName).Scan(&id)
	return id, err
}

// --- Versions ---

func (r *Repository) CreateVersion(ctx context.Context, v *entity.Version) error {
	if v.DefinitionSnapshot == nil {
		v.DefinitionSnapshot = json.RawMessage(`{}`)
	}
	return r.db.QueryRow(ctx, `
		INSERT INTO workflow_versions (
			workflow_id, organization_id, version, state, definition_snapshot, changelog, published_at, published_by
		) VALUES ($1,$2,$3,$4,$5,$6,$7,$8)
		RETURNING id, created_at`,
		v.WorkflowID, v.OrganizationID, v.Version, v.State, v.DefinitionSnapshot, v.Changelog, v.PublishedAt, v.PublishedBy,
	).Scan(&v.ID, &v.CreatedAt)
}

func (r *Repository) UpdateVersion(ctx context.Context, v *entity.Version) error {
	_, err := r.db.Exec(ctx, `
		UPDATE workflow_versions SET
			state = $3, definition_snapshot = $4, changelog = $5,
			published_at = $6, published_by = $7
		WHERE id = $1 AND organization_id = $2`,
		v.ID, v.OrganizationID, v.State, v.DefinitionSnapshot, v.Changelog, v.PublishedAt, v.PublishedBy,
	)
	return err
}

func (r *Repository) GetVersion(ctx context.Context, orgID, id string) (*entity.Version, error) {
	row := r.db.QueryRow(ctx, `
		SELECT id, workflow_id, organization_id, version, state, definition_snapshot,
			changelog, published_at, published_by, created_at
		FROM workflow_versions WHERE id = $1 AND organization_id = $2`, id, orgID)
	return scanVersion(row)
}

func (r *Repository) LatestDraftVersion(ctx context.Context, orgID, workflowID string) (*entity.Version, error) {
	row := r.db.QueryRow(ctx, `
		SELECT id, workflow_id, organization_id, version, state, definition_snapshot,
			changelog, published_at, published_by, created_at
		FROM workflow_versions
		WHERE organization_id = $1 AND workflow_id = $2 AND state = 'draft'
		ORDER BY version DESC LIMIT 1`, orgID, workflowID)
	v, err := scanVersion(row)
	if err == pgx.ErrNoRows {
		return nil, nil
	}
	return v, err
}

func (r *Repository) MaxVersion(ctx context.Context, workflowID string) (int, error) {
	var n int
	err := r.db.QueryRow(ctx, `SELECT COALESCE(MAX(version), 0) FROM workflow_versions WHERE workflow_id = $1`, workflowID).Scan(&n)
	return n, err
}

func (r *Repository) ListVersions(ctx context.Context, orgID, workflowID string) ([]entity.Version, error) {
	rows, err := r.db.Query(ctx, `
		SELECT id, workflow_id, organization_id, version, state, definition_snapshot,
			changelog, published_at, published_by, created_at
		FROM workflow_versions
		WHERE organization_id = $1 AND workflow_id = $2
		ORDER BY version DESC`, orgID, workflowID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []entity.Version
	for rows.Next() {
		v, err := scanVersionRows(rows)
		if err != nil {
			return nil, err
		}
		out = append(out, *v)
	}
	return out, rows.Err()
}

func (r *Repository) DeleteVersionChildren(ctx context.Context, versionID string) error {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)
	for _, q := range []string{
		`DELETE FROM workflow_actions WHERE version_id = $1`,
		`DELETE FROM workflow_conditions WHERE version_id = $1`,
		`DELETE FROM workflow_triggers WHERE version_id = $1`,
	} {
		if _, err := tx.Exec(ctx, q, versionID); err != nil {
			return err
		}
	}
	return tx.Commit(ctx)
}

// --- Triggers / Conditions / Actions ---

func (r *Repository) InsertTrigger(ctx context.Context, t *entity.Trigger) error {
	if t.Config == nil {
		t.Config = json.RawMessage(`{}`)
	}
	return r.db.QueryRow(ctx, `
		INSERT INTO workflow_triggers (version_id, organization_id, type, config)
		VALUES ($1,$2,$3,$4) RETURNING id, created_at`,
		t.VersionID, t.OrganizationID, t.Type, t.Config,
	).Scan(&t.ID, &t.CreatedAt)
}

func (r *Repository) InsertCondition(ctx context.Context, c *entity.Condition) error {
	return r.db.QueryRow(ctx, `
		INSERT INTO workflow_conditions (
			version_id, organization_id, parent_id, node_type, logic,
			field_api_name, operator, value, sort_order
		) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9)
		RETURNING id, created_at`,
		c.VersionID, c.OrganizationID, c.ParentID, c.NodeType, c.Logic,
		c.FieldAPIName, c.Operator, c.Value, c.SortOrder,
	).Scan(&c.ID, &c.CreatedAt)
}

func (r *Repository) InsertAction(ctx context.Context, a *entity.Action) error {
	if a.Config == nil {
		a.Config = json.RawMessage(`{}`)
	}
	return r.db.QueryRow(ctx, `
		INSERT INTO workflow_actions (
			version_id, organization_id, sort_order, type, config, max_retries, continue_on_error
		) VALUES ($1,$2,$3,$4,$5,$6,$7)
		RETURNING id, created_at`,
		a.VersionID, a.OrganizationID, a.SortOrder, a.Type, a.Config, a.MaxRetries, a.ContinueOnError,
	).Scan(&a.ID, &a.CreatedAt)
}

func (r *Repository) ListTriggers(ctx context.Context, versionID string) ([]entity.Trigger, error) {
	rows, err := r.db.Query(ctx, `
		SELECT id, version_id, organization_id, type, config, created_at
		FROM workflow_triggers WHERE version_id = $1 ORDER BY created_at`, versionID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []entity.Trigger
	for rows.Next() {
		var t entity.Trigger
		if err := rows.Scan(&t.ID, &t.VersionID, &t.OrganizationID, &t.Type, &t.Config, &t.CreatedAt); err != nil {
			return nil, err
		}
		out = append(out, t)
	}
	return out, rows.Err()
}

func (r *Repository) ListConditions(ctx context.Context, versionID string) ([]entity.Condition, error) {
	rows, err := r.db.Query(ctx, `
		SELECT id, version_id, organization_id, parent_id, node_type, logic,
			field_api_name, operator, value, sort_order, created_at
		FROM workflow_conditions WHERE version_id = $1 ORDER BY sort_order, created_at`, versionID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []entity.Condition
	for rows.Next() {
		var c entity.Condition
		if err := rows.Scan(&c.ID, &c.VersionID, &c.OrganizationID, &c.ParentID, &c.NodeType, &c.Logic,
			&c.FieldAPIName, &c.Operator, &c.Value, &c.SortOrder, &c.CreatedAt); err != nil {
			return nil, err
		}
		out = append(out, c)
	}
	return out, rows.Err()
}

func (r *Repository) ListActions(ctx context.Context, versionID string) ([]entity.Action, error) {
	rows, err := r.db.Query(ctx, `
		SELECT id, version_id, organization_id, sort_order, type, config, max_retries, continue_on_error, created_at
		FROM workflow_actions WHERE version_id = $1 ORDER BY sort_order`, versionID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []entity.Action
	for rows.Next() {
		var a entity.Action
		if err := rows.Scan(&a.ID, &a.VersionID, &a.OrganizationID, &a.SortOrder, &a.Type, &a.Config,
			&a.MaxRetries, &a.ContinueOnError, &a.CreatedAt); err != nil {
			return nil, err
		}
		out = append(out, a)
	}
	return out, rows.Err()
}

// ListActiveMatches returns active workflows for org/module whose published version has the trigger type.
func (r *Repository) ListActiveMatches(ctx context.Context, orgID, moduleID, triggerType string) ([]entity.MatchCandidate, error) {
	rows, err := r.db.Query(ctx, `
		SELECT w.id, w.organization_id, w.module_id, w.name, w.description, w.status,
			w.on_action_error, w.priority, w.published_version_id, w.created_by, w.updated_by,
			w.created_at, w.updated_at,
			v.id, v.workflow_id, v.organization_id, v.version, v.state, v.definition_snapshot,
			v.changelog, v.published_at, v.published_by, v.created_at
		FROM workflows w
		JOIN workflow_versions v ON v.id = w.published_version_id
		WHERE w.organization_id = $1
		  AND w.status = 'active'
		  AND w.published_version_id IS NOT NULL
		  AND (w.module_id IS NULL OR w.module_id = $2)
		  AND EXISTS (
			SELECT 1 FROM workflow_triggers t
			WHERE t.version_id = v.id AND t.type = $3
		  )
		ORDER BY w.priority ASC, w.created_at ASC`, orgID, moduleID, triggerType)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []entity.MatchCandidate
	for rows.Next() {
		var m entity.MatchCandidate
		if err := rows.Scan(
			&m.Workflow.ID, &m.Workflow.OrganizationID, &m.Workflow.ModuleID, &m.Workflow.Name, &m.Workflow.Description, &m.Workflow.Status,
			&m.Workflow.OnActionError, &m.Workflow.Priority, &m.Workflow.PublishedVersionID, &m.Workflow.CreatedBy, &m.Workflow.UpdatedBy,
			&m.Workflow.CreatedAt, &m.Workflow.UpdatedAt,
			&m.Version.ID, &m.Version.WorkflowID, &m.Version.OrganizationID, &m.Version.Version, &m.Version.State, &m.Version.DefinitionSnapshot,
			&m.Version.Changelog, &m.Version.PublishedAt, &m.Version.PublishedBy, &m.Version.CreatedAt,
		); err != nil {
			return nil, err
		}
		out = append(out, m)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	for i := range out {
		tr, err := r.ListTriggers(ctx, out[i].Version.ID)
		if err != nil {
			return nil, err
		}
		cond, err := r.ListConditions(ctx, out[i].Version.ID)
		if err != nil {
			return nil, err
		}
		acts, err := r.ListActions(ctx, out[i].Version.ID)
		if err != nil {
			return nil, err
		}
		out[i].Triggers = tr
		out[i].Conditions = cond
		out[i].Actions = acts
	}
	return out, nil
}

func (r *Repository) GetMatchByWorkflow(ctx context.Context, orgID, workflowID string) (*entity.MatchCandidate, error) {
	w, err := r.GetWorkflow(ctx, orgID, workflowID)
	if err != nil {
		return nil, err
	}
	if w == nil || w.PublishedVersionID == nil {
		return nil, nil
	}
	v, err := r.GetVersion(ctx, orgID, *w.PublishedVersionID)
	if err != nil {
		return nil, err
	}
	tr, err := r.ListTriggers(ctx, v.ID)
	if err != nil {
		return nil, err
	}
	cond, err := r.ListConditions(ctx, v.ID)
	if err != nil {
		return nil, err
	}
	acts, err := r.ListActions(ctx, v.ID)
	if err != nil {
		return nil, err
	}
	return &entity.MatchCandidate{Workflow: *w, Version: *v, Triggers: tr, Conditions: cond, Actions: acts}, nil
}

// --- Executions ---

func (r *Repository) CreateExecution(ctx context.Context, e *entity.Execution) error {
	return r.db.QueryRow(ctx, `
		INSERT INTO workflow_executions (
			organization_id, workflow_id, version_id, module_id, record_id,
			trigger_type, status, source, depth, started_at
		) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10)
		RETURNING id, created_at`,
		e.OrganizationID, e.WorkflowID, e.VersionID, e.ModuleID, e.RecordID,
		e.TriggerType, e.Status, e.Source, e.Depth, e.StartedAt,
	).Scan(&e.ID, &e.CreatedAt)
}

func (r *Repository) FinishExecution(ctx context.Context, orgID, id, status string, errSummary *string, started, finished time.Time, durationMs int) error {
	_, err := r.db.Exec(ctx, `
		UPDATE workflow_executions SET
			status = $3, error_summary = $4, started_at = $5, finished_at = $6, duration_ms = $7
		WHERE id = $1 AND organization_id = $2`,
		id, orgID, status, errSummary, started, finished, durationMs,
	)
	return err
}

func (r *Repository) CreateExecutionStep(ctx context.Context, s *entity.ExecutionStep) error {
	if s.Input == nil {
		s.Input = json.RawMessage(`{}`)
	}
	if s.Output == nil {
		s.Output = json.RawMessage(`{}`)
	}
	return r.db.QueryRow(ctx, `
		INSERT INTO workflow_execution_steps (
			execution_id, organization_id, action_id, sort_order, action_type,
			status, input, output, error, started_at, finished_at
		) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11)
		RETURNING id, created_at`,
		s.ExecutionID, s.OrganizationID, s.ActionID, s.SortOrder, s.ActionType,
		s.Status, s.Input, s.Output, s.Error, s.StartedAt, s.FinishedAt,
	).Scan(&s.ID, &s.CreatedAt)
}

func (r *Repository) ListExecutions(ctx context.Context, orgID string, workflowID, moduleID, recordID, status string, page, pageSize int) ([]entity.Execution, map[string]string, int, error) {
	args := []any{orgID}
	where := []string{"e.organization_id = $1"}
	if workflowID != "" {
		args = append(args, workflowID)
		where = append(where, fmt.Sprintf("e.workflow_id = $%d", len(args)))
	}
	if moduleID != "" {
		args = append(args, moduleID)
		where = append(where, fmt.Sprintf("e.module_id = $%d", len(args)))
	}
	if recordID != "" {
		args = append(args, recordID)
		where = append(where, fmt.Sprintf("e.record_id = $%d", len(args)))
	}
	if status != "" {
		args = append(args, status)
		where = append(where, fmt.Sprintf("e.status = $%d", len(args)))
	}
	wSQL := strings.Join(where, " AND ")

	var total int
	if err := r.db.QueryRow(ctx, `SELECT COUNT(*) FROM workflow_executions e WHERE `+wSQL, args...).Scan(&total); err != nil {
		return nil, nil, 0, err
	}

	offset := (page - 1) * pageSize
	args = append(args, pageSize, offset)
	rows, err := r.db.Query(ctx, `
		SELECT e.id, e.organization_id, e.workflow_id, e.version_id, e.module_id, e.record_id,
			e.trigger_type, e.status, e.source, e.depth, e.error_summary,
			e.started_at, e.finished_at, e.duration_ms, e.created_at, w.name
		FROM workflow_executions e
		JOIN workflows w ON w.id = e.workflow_id
		WHERE `+wSQL+`
		ORDER BY e.created_at DESC
		LIMIT $`+fmt.Sprint(len(args)-1)+` OFFSET $`+fmt.Sprint(len(args)), args...)
	if err != nil {
		return nil, nil, 0, err
	}
	defer rows.Close()

	var out []entity.Execution
	names := map[string]string{}
	for rows.Next() {
		var e entity.Execution
		var name string
		if err := rows.Scan(&e.ID, &e.OrganizationID, &e.WorkflowID, &e.VersionID, &e.ModuleID, &e.RecordID,
			&e.TriggerType, &e.Status, &e.Source, &e.Depth, &e.ErrorSummary,
			&e.StartedAt, &e.FinishedAt, &e.DurationMs, &e.CreatedAt, &name); err != nil {
			return nil, nil, 0, err
		}
		out = append(out, e)
		names[e.ID] = name
	}
	return out, names, total, rows.Err()
}

func (r *Repository) GetExecution(ctx context.Context, orgID, id string) (*entity.Execution, string, error) {
	var e entity.Execution
	var name string
	err := r.db.QueryRow(ctx, `
		SELECT e.id, e.organization_id, e.workflow_id, e.version_id, e.module_id, e.record_id,
			e.trigger_type, e.status, e.source, e.depth, e.error_summary,
			e.started_at, e.finished_at, e.duration_ms, e.created_at, w.name
		FROM workflow_executions e
		JOIN workflows w ON w.id = e.workflow_id
		WHERE e.id = $1 AND e.organization_id = $2`, id, orgID,
	).Scan(&e.ID, &e.OrganizationID, &e.WorkflowID, &e.VersionID, &e.ModuleID, &e.RecordID,
		&e.TriggerType, &e.Status, &e.Source, &e.Depth, &e.ErrorSummary,
		&e.StartedAt, &e.FinishedAt, &e.DurationMs, &e.CreatedAt, &name)
	if err == pgx.ErrNoRows {
		return nil, "", nil
	}
	if err != nil {
		return nil, "", err
	}
	return &e, name, nil
}

func (r *Repository) ListExecutionSteps(ctx context.Context, executionID string) ([]entity.ExecutionStep, error) {
	rows, err := r.db.Query(ctx, `
		SELECT id, execution_id, organization_id, action_id, sort_order, action_type,
			status, input, output, error, started_at, finished_at, created_at
		FROM workflow_execution_steps WHERE execution_id = $1 ORDER BY sort_order`, executionID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []entity.ExecutionStep
	for rows.Next() {
		var s entity.ExecutionStep
		if err := rows.Scan(&s.ID, &s.ExecutionID, &s.OrganizationID, &s.ActionID, &s.SortOrder, &s.ActionType,
			&s.Status, &s.Input, &s.Output, &s.Error, &s.StartedAt, &s.FinishedAt, &s.CreatedAt); err != nil {
			return nil, err
		}
		out = append(out, s)
	}
	return out, rows.Err()
}

// --- Templates ---

func (r *Repository) ListTemplates(ctx context.Context, orgID string) ([]entity.Template, error) {
	rows, err := r.db.Query(ctx, `
		SELECT id, organization_id, key, name, description, module_api_name, definition, is_builtin, created_at, updated_at
		FROM workflow_templates
		WHERE organization_id IS NULL OR organization_id = $1
		ORDER BY is_builtin DESC, name ASC`, orgID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []entity.Template
	for rows.Next() {
		var t entity.Template
		if err := rows.Scan(&t.ID, &t.OrganizationID, &t.Key, &t.Name, &t.Description, &t.ModuleAPIName,
			&t.Definition, &t.IsBuiltin, &t.CreatedAt, &t.UpdatedAt); err != nil {
			return nil, err
		}
		out = append(out, t)
	}
	return out, rows.Err()
}

func (r *Repository) GetTemplate(ctx context.Context, orgID, id string) (*entity.Template, error) {
	var t entity.Template
	err := r.db.QueryRow(ctx, `
		SELECT id, organization_id, key, name, description, module_api_name, definition, is_builtin, created_at, updated_at
		FROM workflow_templates
		WHERE id = $1 AND (organization_id IS NULL OR organization_id = $2)`, id, orgID,
	).Scan(&t.ID, &t.OrganizationID, &t.Key, &t.Name, &t.Description, &t.ModuleAPIName,
		&t.Definition, &t.IsBuiltin, &t.CreatedAt, &t.UpdatedAt)
	if err == pgx.ErrNoRows {
		return nil, nil
	}
	return &t, err
}

func (r *Repository) UpsertBuiltinTemplate(ctx context.Context, t *entity.Template) error {
	if t.Definition == nil {
		t.Definition = json.RawMessage(`{}`)
	}
	return r.db.QueryRow(ctx, `
		INSERT INTO workflow_templates (organization_id, key, name, description, module_api_name, definition, is_builtin)
		VALUES (NULL, $1, $2, $3, $4, $5, TRUE)
		ON CONFLICT (key) WHERE is_builtin = TRUE AND organization_id IS NULL
		DO UPDATE SET name = EXCLUDED.name, description = EXCLUDED.description,
			module_api_name = EXCLUDED.module_api_name, definition = EXCLUDED.definition, updated_at = NOW()
		RETURNING id, created_at, updated_at`,
		t.Key, t.Name, t.Description, t.ModuleAPIName, t.Definition,
	).Scan(&t.ID, &t.CreatedAt, &t.UpdatedAt)
}

func (r *Repository) CountBuiltinTemplates(ctx context.Context) (int, error) {
	var n int
	err := r.db.QueryRow(ctx, `
		SELECT COUNT(*) FROM workflow_templates WHERE is_builtin = TRUE AND organization_id IS NULL`).Scan(&n)
	return n, err
}

// --- Metrics / builder helpers ---

func (r *Repository) Metrics(ctx context.Context, orgID string) (active, disabled, draft, executedToday, failedToday int, avgMs *float64, err error) {
	err = r.db.QueryRow(ctx, `
		SELECT
			COUNT(*) FILTER (WHERE status = 'active'),
			COUNT(*) FILTER (WHERE status = 'disabled'),
			COUNT(*) FILTER (WHERE status = 'draft')
		FROM workflows WHERE organization_id = $1`, orgID,
	).Scan(&active, &disabled, &draft)
	if err != nil {
		return
	}
	err = r.db.QueryRow(ctx, `
		SELECT
			COUNT(*) FILTER (WHERE created_at::date = CURRENT_DATE),
			COUNT(*) FILTER (WHERE created_at::date = CURRENT_DATE AND status = 'failed'),
			AVG(duration_ms) FILTER (WHERE created_at::date = CURRENT_DATE AND duration_ms IS NOT NULL)
		FROM workflow_executions WHERE organization_id = $1`, orgID,
	).Scan(&executedToday, &failedToday, &avgMs)
	return
}

func (r *Repository) ListOrgModules(ctx context.Context, orgID string) ([]struct{ ID, APIName, Label string }, error) {
	rows, err := r.db.Query(ctx, `
		SELECT id, api_name, plural_label FROM modules
		WHERE organization_id = $1 AND is_enabled = TRUE
		ORDER BY sort_order, plural_label`, orgID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []struct{ ID, APIName, Label string }
	for rows.Next() {
		var m struct{ ID, APIName, Label string }
		if err := rows.Scan(&m.ID, &m.APIName, &m.Label); err != nil {
			return nil, err
		}
		out = append(out, m)
	}
	return out, rows.Err()
}

func (r *Repository) ListOrgUsers(ctx context.Context, orgID string) ([]struct{ ID, Name, Email string }, error) {
	rows, err := r.db.Query(ctx, `
		SELECT u.id, COALESCE(u.name, u.email), u.email
		FROM users u
		JOIN organization_members m ON m.user_id = u.id
		WHERE m.organization_id = $1 AND m.status = 'active'
		ORDER BY u.name NULLS LAST, u.email
		LIMIT 200`, orgID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []struct{ ID, Name, Email string }
	for rows.Next() {
		var u struct{ ID, Name, Email string }
		if err := rows.Scan(&u.ID, &u.Name, &u.Email); err != nil {
			return nil, err
		}
		out = append(out, u)
	}
	return out, rows.Err()
}

func (r *Repository) OrgName(ctx context.Context, orgID string) (string, error) {
	var name string
	err := r.db.QueryRow(ctx, `SELECT name FROM organizations WHERE id = $1`, orgID).Scan(&name)
	return name, err
}

func (r *Repository) UserDisplay(ctx context.Context, userID string) (name, email string, err error) {
	err = r.db.QueryRow(ctx, `SELECT COALESCE(name, email), email FROM users WHERE id = $1`, userID).Scan(&name, &email)
	return
}

// HasRunToday reports whether a workflow already executed for a record+trigger today (UTC).
// When recordID is empty, checks org-level runs for that workflow+trigger.
func (r *Repository) HasRunToday(ctx context.Context, orgID, workflowID, recordID, triggerType string) (bool, error) {
	var n int
	var err error
	if recordID == "" {
		err = r.db.QueryRow(ctx, `
			SELECT COUNT(*) FROM workflow_executions
			WHERE organization_id = $1 AND workflow_id = $2 AND trigger_type = $3
			  AND record_id IS NULL
			  AND created_at::date = (NOW() AT TIME ZONE 'utc')::date
			  AND status IN ('succeeded', 'partial', 'running', 'queued')`,
			orgID, workflowID, triggerType).Scan(&n)
	} else {
		err = r.db.QueryRow(ctx, `
			SELECT COUNT(*) FROM workflow_executions
			WHERE organization_id = $1 AND workflow_id = $2 AND trigger_type = $3
			  AND record_id = $4
			  AND created_at::date = (NOW() AT TIME ZONE 'utc')::date
			  AND status IN ('succeeded', 'partial', 'running', 'queued')`,
			orgID, workflowID, triggerType, recordID).Scan(&n)
	}
	return n > 0, err
}

// ListModuleRecordIDs returns up to limit record ids (and data) for a module.
func (r *Repository) ListModuleRecordIDs(ctx context.Context, orgID, moduleID string, limit int) ([]struct {
	ID   string
	Data []byte
}, error) {
	if limit <= 0 {
		limit = 100
	}
	rows, err := r.db.Query(ctx, `
		SELECT id, data FROM records
		WHERE organization_id = $1 AND module_id = $2
		ORDER BY updated_at DESC
		LIMIT $3`, orgID, moduleID, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []struct {
		ID   string
		Data []byte
	}
	for rows.Next() {
		var rec struct {
			ID   string
			Data []byte
		}
		if err := rows.Scan(&rec.ID, &rec.Data); err != nil {
			return nil, err
		}
		out = append(out, rec)
	}
	return out, rows.Err()
}

// scanners

type scannable interface {
	Scan(dest ...any) error
}

func scanWorkflow(row scannable) (*entity.Workflow, error) {
	var w entity.Workflow
	err := row.Scan(&w.ID, &w.OrganizationID, &w.ModuleID, &w.Name, &w.Description, &w.Status,
		&w.OnActionError, &w.Priority, &w.PublishedVersionID, &w.CreatedBy, &w.UpdatedBy,
		&w.CreatedAt, &w.UpdatedAt)
	if err == pgx.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &w, nil
}

func scanWorkflowRows(rows pgx.Rows) (*entity.Workflow, error) {
	return scanWorkflow(rows)
}

func scanVersion(row scannable) (*entity.Version, error) {
	var v entity.Version
	err := row.Scan(&v.ID, &v.WorkflowID, &v.OrganizationID, &v.Version, &v.State, &v.DefinitionSnapshot,
		&v.Changelog, &v.PublishedAt, &v.PublishedBy, &v.CreatedAt)
	if err == pgx.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &v, nil
}

func scanVersionRows(rows pgx.Rows) (*entity.Version, error) {
	return scanVersion(rows)
}
