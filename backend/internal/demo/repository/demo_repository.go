package repository

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository struct {
	db *pgxpool.Pool
}

func New(db *pgxpool.Pool) *Repository { return &Repository{db: db} }

type Workflow struct {
	Key, Name, Description string
	Version, DurationMin   int
}

type Step struct {
	Key, Title, Description, Why, How, Expected string
	SortOrder                                   int
	Route, Target, ActionLabel                  *string
	ValidatorKey                                string
	ValidatorParams                             json.RawMessage
	IsSkippable                                 bool
	RequiredAction                              string
	SuccessEvent                                *string
	FailureMessage                              string
	Hint                                        string
	MaxAttempts                                 int
	AllowSelectors                              json.RawMessage
	Placement                                   string
}

type Session struct {
	ID                       string
	UserID                   string
	WorkflowKey              string
	WorkflowVersion          int
	SandboxOrganizationID    *string
	PreviousOrganizationID   *string
	Status                   string
	CurrentStepKey           *string
	StartedAt                time.Time
	CompletedAt              *time.Time
	KeepData                 *bool
	Stats                    json.RawMessage
}

func (r *Repository) GetWorkflow(ctx context.Context, key string) (*Workflow, error) {
	var w Workflow
	err := r.db.QueryRow(ctx, `
		SELECT workflow_key, name, COALESCE(description,''), version, duration_min
		FROM demo_workflows WHERE workflow_key = $1 AND is_active = TRUE
	`, key).Scan(&w.Key, &w.Name, &w.Description, &w.Version, &w.DurationMin)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}
	return &w, err
}

func (r *Repository) ListSteps(ctx context.Context, workflowKey string) ([]Step, error) {
	rows, err := r.db.Query(ctx, `
		SELECT step_key, sort_order, title, description, why_it_matters, how_it_works,
		       expected_result, route, target_selector, action_label, validator_key,
		       validator_params, is_skippable,
		       COALESCE(required_action, 'acknowledge'), success_event,
		       COALESCE(failure_message, ''), COALESCE(hint, ''), COALESCE(max_attempts, 0),
		       COALESCE(allow_selectors, '[]'::jsonb), COALESCE(placement, 'center')
		FROM demo_workflow_steps
		WHERE workflow_key = $1
		ORDER BY sort_order ASC
	`, workflowKey)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	out := make([]Step, 0)
	for rows.Next() {
		var s Step
		if err := rows.Scan(
			&s.Key, &s.SortOrder, &s.Title, &s.Description, &s.Why, &s.How, &s.Expected,
			&s.Route, &s.Target, &s.ActionLabel, &s.ValidatorKey, &s.ValidatorParams, &s.IsSkippable,
			&s.RequiredAction, &s.SuccessEvent, &s.FailureMessage, &s.Hint, &s.MaxAttempts,
			&s.AllowSelectors, &s.Placement,
		); err != nil {
			return nil, err
		}
		out = append(out, s)
	}
	return out, rows.Err()
}

func (r *Repository) ActiveSession(ctx context.Context, userID, workflowKey string) (*Session, error) {
	var s Session
	err := r.db.QueryRow(ctx, `
		SELECT id::text, user_id::text, workflow_key, workflow_version,
		       sandbox_organization_id::text, previous_organization_id::text,
		       status, current_step_key, started_at, completed_at, keep_data, stats
		FROM demo_sessions
		WHERE user_id = $1 AND workflow_key = $2 AND status IN ('active', 'completed')
		ORDER BY CASE status WHEN 'active' THEN 0 ELSE 1 END, started_at DESC
		LIMIT 1
	`, userID, workflowKey).Scan(
		&s.ID, &s.UserID, &s.WorkflowKey, &s.WorkflowVersion,
		&s.SandboxOrganizationID, &s.PreviousOrganizationID,
		&s.Status, &s.CurrentStepKey, &s.StartedAt, &s.CompletedAt, &s.KeepData, &s.Stats,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}
	return &s, err
}

func (r *Repository) ModuleIDByAPIName(ctx context.Context, orgID, apiName string) (string, error) {
	var id string
	err := r.db.QueryRow(ctx, `
		SELECT id::text FROM modules
		WHERE organization_id = $1 AND api_name = $2
		LIMIT 1
	`, orgID, apiName).Scan(&id)
	if errors.Is(err, pgx.ErrNoRows) {
		return "", nil
	}
	return id, err
}

func (r *Repository) GetSession(ctx context.Context, sessionID, userID string) (*Session, error) {
	var s Session
	err := r.db.QueryRow(ctx, `
		SELECT id::text, user_id::text, workflow_key, workflow_version,
		       sandbox_organization_id::text, previous_organization_id::text,
		       status, current_step_key, started_at, completed_at, keep_data, stats
		FROM demo_sessions WHERE id = $1 AND user_id = $2
	`, sessionID, userID).Scan(
		&s.ID, &s.UserID, &s.WorkflowKey, &s.WorkflowVersion,
		&s.SandboxOrganizationID, &s.PreviousOrganizationID,
		&s.Status, &s.CurrentStepKey, &s.StartedAt, &s.CompletedAt, &s.KeepData, &s.Stats,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}
	return &s, err
}

func (r *Repository) CreateSession(ctx context.Context, s *Session) error {
	s.ID = uuid.NewString()
	s.StartedAt = time.Now().UTC()
	if s.Stats == nil {
		s.Stats = json.RawMessage(`{}`)
	}
	_, err := r.db.Exec(ctx, `
		INSERT INTO demo_sessions (
			id, user_id, workflow_key, workflow_version, sandbox_organization_id,
			previous_organization_id, status, current_step_key, started_at, updated_at, stats
		) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$9,$10)
	`, s.ID, s.UserID, s.WorkflowKey, s.WorkflowVersion, s.SandboxOrganizationID,
		s.PreviousOrganizationID, s.Status, s.CurrentStepKey, s.StartedAt, s.Stats)
	return err
}

func (r *Repository) UpdateSession(ctx context.Context, s *Session) error {
	_, err := r.db.Exec(ctx, `
		UPDATE demo_sessions SET
			status = $2, current_step_key = $3, completed_at = $4,
			keep_data = $5, stats = $6, updated_at = NOW()
		WHERE id = $1
	`, s.ID, s.Status, s.CurrentStepKey, s.CompletedAt, s.KeepData, s.Stats)
	return err
}

func (r *Repository) InitStepProgress(ctx context.Context, sessionID string, steps []Step) error {
	for i, st := range steps {
		status := "locked"
		if i == 0 {
			status = "active"
		}
		_, err := r.db.Exec(ctx, `
			INSERT INTO demo_step_progress (session_id, step_key, status)
			VALUES ($1,$2,$3)
			ON CONFLICT (session_id, step_key) DO NOTHING
		`, sessionID, st.Key, status)
		if err != nil {
			return err
		}
	}
	return nil
}

func (r *Repository) ListProgress(ctx context.Context, sessionID string) (map[string]string, error) {
	rows, err := r.db.Query(ctx, `
		SELECT step_key, status FROM demo_step_progress WHERE session_id = $1
	`, sessionID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	out := map[string]string{}
	for rows.Next() {
		var k, st string
		if err := rows.Scan(&k, &st); err != nil {
			return nil, err
		}
		out[k] = st
	}
	return out, rows.Err()
}

func (r *Repository) SetStepStatus(ctx context.Context, sessionID, stepKey, status, lastErr string) error {
	var completed any
	if status == "completed" || status == "skipped" {
		completed = time.Now().UTC()
	}
	_, err := r.db.Exec(ctx, `
		UPDATE demo_step_progress
		SET status = $3, attempts = attempts + 1, last_error = NULLIF($4,''),
		    completed_at = COALESCE($5::timestamptz, completed_at)
		WHERE session_id = $1 AND step_key = $2
	`, sessionID, stepKey, status, lastErr, completed)
	return err
}

func (r *Repository) TrackResource(ctx context.Context, sessionID, resourceType, resourceID string, meta map[string]any) error {
	raw, _ := json.Marshal(meta)
	if raw == nil {
		raw = []byte(`{}`)
	}
	_, err := r.db.Exec(ctx, `
		INSERT INTO demo_resources (session_id, resource_type, resource_id, meta)
		VALUES ($1,$2,$3,$4)
	`, sessionID, resourceType, resourceID, raw)
	return err
}

func (r *Repository) LogEvent(ctx context.Context, sessionID, eventType string, payload map[string]any) error {
	raw, _ := json.Marshal(payload)
	if raw == nil {
		raw = []byte(`{}`)
	}
	_, err := r.db.Exec(ctx, `
		INSERT INTO demo_events (session_id, event_type, payload) VALUES ($1,$2,$3)
	`, sessionID, eventType, raw)
	return err
}

func (r *Repository) GetActiveOrgID(ctx context.Context, userID string) (*string, error) {
	var id *string
	err := r.db.QueryRow(ctx, `SELECT active_organization_id::text FROM users WHERE id = $1`, userID).Scan(&id)
	return id, err
}

func (r *Repository) SetActiveOrg(ctx context.Context, userID, orgID string) error {
	_, err := r.db.Exec(ctx, `
		UPDATE users SET active_organization_id = $2, updated_at = NOW() WHERE id = $1
	`, userID, orgID)
	return err
}

func (r *Repository) DeleteOrganization(ctx context.Context, orgID string) error {
	_, err := r.db.Exec(ctx, `DELETE FROM organizations WHERE id = $1`, orgID)
	return err
}

func (r *Repository) ModuleExists(ctx context.Context, orgID, apiName string) (bool, error) {
	var ok bool
	err := r.db.QueryRow(ctx, `
		SELECT EXISTS(SELECT 1 FROM modules WHERE organization_id = $1 AND api_name = $2)
	`, orgID, apiName).Scan(&ok)
	return ok, err
}

func (r *Repository) FieldExists(ctx context.Context, orgID, moduleAPI, fieldAPI string) (bool, error) {
	var ok bool
	err := r.db.QueryRow(ctx, `
		SELECT EXISTS(
			SELECT 1 FROM fields f
			JOIN modules m ON m.id = f.module_id
			WHERE m.organization_id = $1 AND m.api_name = $2 AND f.api_name = $3
		)
	`, orgID, moduleAPI, fieldAPI).Scan(&ok)
	return ok, err
}

func (r *Repository) RecordExists(ctx context.Context, orgID, moduleAPI string) (bool, error) {
	var ok bool
	err := r.db.QueryRow(ctx, `
		SELECT EXISTS(
			SELECT 1 FROM records r
			JOIN modules m ON m.id = r.module_id
			WHERE r.organization_id = $1 AND m.api_name = $2
		)
	`, orgID, moduleAPI).Scan(&ok)
	return ok, err
}

func (r *Repository) NoteExists(ctx context.Context, orgID, moduleAPI string) (bool, error) {
	var ok bool
	err := r.db.QueryRow(ctx, `
		SELECT EXISTS(
			SELECT 1 FROM notes n
			JOIN modules m ON m.id = n.module_id
			WHERE n.organization_id = $1 AND m.api_name = $2 AND n.entity_type = 'RECORD'
		)
	`, orgID, moduleAPI).Scan(&ok)
	return ok, err
}

func (r *Repository) ActivityExists(ctx context.Context, orgID, moduleAPI string) (bool, error) {
	var ok bool
	err := r.db.QueryRow(ctx, `
		SELECT EXISTS(
			SELECT 1 FROM activities a
			JOIN modules m ON m.id = a.module_id
			WHERE a.organization_id = $1 AND m.api_name = $2 AND a.entity_type = 'RECORD'
		)
	`, orgID, moduleAPI).Scan(&ok)
	return ok, err
}

func (r *Repository) GetModuleID(ctx context.Context, orgID, apiName string) (string, error) {
	var id string
	err := r.db.QueryRow(ctx, `
		SELECT id::text FROM modules WHERE organization_id = $1 AND api_name = $2
	`, orgID, apiName).Scan(&id)
	if errors.Is(err, pgx.ErrNoRows) {
		return "", err
	}
	return id, err
}

// SeedProductDemoModule creates a showcase module with varied field types + sample records.
func (r *Repository) SeedProductDemoModule(ctx context.Context, orgID, userID string) error {
	var exists bool
	if err := r.db.QueryRow(ctx, `
		SELECT EXISTS(SELECT 1 FROM modules WHERE organization_id = $1 AND api_name = 'product_demo')
	`, orgID).Scan(&exists); err != nil {
		return err
	}
	if exists {
		return nil
	}

	var moduleID string
	err := r.db.QueryRow(ctx, `
		INSERT INTO modules (
			organization_id, api_name, singular_label, plural_label,
			icon, color, storage_strategy, is_system, sort_order,
			is_enabled, is_visible_sidebar
		) VALUES ($1,'product_demo','Product Demo','Product Demos','sparkles','#0ea5e9','dynamic',FALSE,50,TRUE,TRUE)
		RETURNING id
	`, orgID).Scan(&moduleID)
	if err != nil {
		return err
	}

	type fld struct {
		API, Label, Type string
		Required         bool
		Options          string
	}
	fields := []fld{
		{"name", "Demo Name", "text", true, "[]"},
		{"score", "Score", "number", false, "[]"},
		{"priority", "Priority", "dropdown", false, `[{"label":"Low","value":"low"},{"label":"Medium","value":"medium"},{"label":"High","value":"high"}]`},
		{"tags", "Tags", "multiselect", false, `[{"label":"Onboarding","value":"onboarding"},{"label":"Sales","value":"sales"},{"label":"Support","value":"support"}]`},
		{"due_date", "Due Date", "date", false, "[]"},
		{"notes", "Notes", "textarea", false, "[]"},
	}
	for i, f := range fields {
		_, err := r.db.Exec(ctx, `
			INSERT INTO fields (
				organization_id, module_id, api_name, label, field_type,
				is_required, is_searchable, is_filterable, options, sort_order, is_system
			) VALUES ($1,$2,$3,$4,$5,$6,TRUE,TRUE,$7::jsonb,$8,FALSE)
		`, orgID, moduleID, f.API, f.Label, f.Type, f.Required, f.Options, i+1)
		if err != nil {
			return err
		}
	}

	samples := []string{
		`{"name":"Welcome Kit","score":90,"priority":"high","tags":["onboarding"],"due_date":"2026-08-01","notes":"Seeded demo record"}`,
		`{"name":"Pipeline Review","score":70,"priority":"medium","tags":["sales"],"due_date":"2026-08-15","notes":"Another sandbox row"}`,
	}
	for _, data := range samples {
		_, err := r.db.Exec(ctx, `
			INSERT INTO records (organization_id, module_id, data, owner_id, created_by, updated_by, visibility)
			VALUES ($1,$2,$3::jsonb,$4,$4,$4,'organization')
		`, orgID, moduleID, data, userID)
		if err != nil {
			return err
		}
	}
	return nil
}
