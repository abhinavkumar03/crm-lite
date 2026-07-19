package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/abhinavkumar03/crm-lite/backend/internal/demo/dto"
	"github.com/abhinavkumar03/crm-lite/backend/internal/demo/repository"
	"github.com/abhinavkumar03/crm-lite/backend/internal/organization/bootstrap"
	"github.com/abhinavkumar03/crm-lite/backend/internal/tenant"
)

const DefaultWorkflow = "crm_interactive_walkthrough"

var (
	ErrNotFound = errors.New("not found")
	ErrConflict = errors.New("conflict")
)

type Service struct {
	repo      *repository.Repository
	bootstrap *bootstrap.Service
	tenant    *tenant.Resolver
}

func New(repo *repository.Repository, boot *bootstrap.Service, tenantResolver *tenant.Resolver) *Service {
	return &Service{repo: repo, bootstrap: boot, tenant: tenantResolver}
}

func (s *Service) Catalog(ctx context.Context) (*dto.WorkflowInfo, error) {
	w, err := s.repo.GetWorkflow(ctx, DefaultWorkflow)
	if err != nil || w == nil {
		return nil, ErrNotFound
	}
	return &dto.WorkflowInfo{
		Key: w.Key, Name: w.Name, Description: w.Description,
		Version: w.Version, DurationMin: w.DurationMin,
		Features: []string{
			"Dynamic Modules", "Dynamic Fields", "Records", "Record Workspace",
			"Notes", "Timeline", "Validation", "Import & Export",
			"Automation", "Roles & Permissions", "Sandbox Isolation",
		},
	}, nil
}

const OrientationWorkflow = "crm_orientation_tour"

// GetWorkflowDefinition returns workflow metadata + steps without a session
// (used by the orientation tour client and future tutorial builders).
func (s *Service) GetWorkflowDefinition(ctx context.Context, key string) (*dto.WorkflowDefinition, error) {
	w, err := s.repo.GetWorkflow(ctx, key)
	if err != nil || w == nil {
		return nil, ErrNotFound
	}
	steps, err := s.repo.ListSteps(ctx, key)
	if err != nil {
		return nil, err
	}
	out := make([]dto.StepDTO, 0, len(steps))
	for _, st := range steps {
		out = append(out, dto.StepDTO{
			Key: st.Key, SortOrder: st.SortOrder, Title: st.Title,
			Description: st.Description, WhyItMatters: st.Why, HowItWorks: st.How,
			ExpectedResult: st.Expected, Route: st.Route, TargetSelector: st.Target,
			ActionLabel: st.ActionLabel, ValidatorKey: st.ValidatorKey,
			ValidatorParams: st.ValidatorParams, IsSkippable: st.IsSkippable, Status: "locked",
			RequiredAction: st.RequiredAction, SuccessEvent: st.SuccessEvent,
			FailureMessage: st.FailureMessage, Hint: st.Hint, MaxAttempts: st.MaxAttempts,
			AllowSelectors: st.AllowSelectors, Placement: st.Placement,
		})
	}
	return &dto.WorkflowDefinition{
		Workflow: dto.WorkflowInfo{
			Key: w.Key, Name: w.Name, Description: w.Description,
			Version: w.Version, DurationMin: w.DurationMin,
		},
		Steps: out,
	}, nil
}

func (s *Service) GetActive(ctx context.Context, userID string) (*dto.SessionDTO, error) {
	sess, err := s.repo.ActiveSession(ctx, userID, DefaultWorkflow)
	if err != nil {
		return nil, err
	}
	if sess == nil {
		return nil, nil
	}
	return s.toDTO(ctx, sess)
}

func (s *Service) Start(ctx context.Context, userID string) (*dto.SessionDTO, error) {
	if existing, _ := s.repo.ActiveSession(ctx, userID, DefaultWorkflow); existing != nil {
		if existing.Status == "active" {
			return s.toDTO(ctx, existing)
		}
		// completed session awaiting cleanup — refuse start until cleaned/restarted
		return nil, fmt.Errorf("finish cleanup of the previous demo before starting again (or use Restart)")
	}

	w, err := s.repo.GetWorkflow(ctx, DefaultWorkflow)
	if err != nil || w == nil {
		return nil, ErrNotFound
	}
	steps, err := s.repo.ListSteps(ctx, DefaultWorkflow)
	if err != nil || len(steps) == 0 {
		return nil, ErrNotFound
	}

	prev, _ := s.repo.GetActiveOrgID(ctx, userID)

	short := userID
	if len(short) > 8 {
		short = short[:8]
	}
	orgID, err := s.bootstrap.CreateOrganization(ctx, bootstrap.CreateOptions{
		Name: "CRM Interactive Demo",
		Slug: fmt.Sprintf("demo-%s-%d", strings.ReplaceAll(short, "-", ""), time.Now().Unix()%100000),
	}, userID)
	if err != nil {
		return nil, fmt.Errorf("sandbox org: %w", err)
	}

	if err := s.seedProductDemo(ctx, orgID, userID); err != nil {
		return nil, err
	}

	first := steps[0].Key
	sess := &repository.Session{
		UserID:                  userID,
		WorkflowKey:             w.Key,
		WorkflowVersion:         w.Version,
		SandboxOrganizationID:   &orgID,
		PreviousOrganizationID:  prev,
		Status:                  "active",
		CurrentStepKey:          &first,
	}
	if err := s.repo.CreateSession(ctx, sess); err != nil {
		return nil, err
	}
	if err := s.repo.InitStepProgress(ctx, sess.ID, steps); err != nil {
		return nil, err
	}
	_ = s.repo.TrackResource(ctx, sess.ID, "organization", orgID, map[string]any{"role": "sandbox"})
	_ = s.repo.LogEvent(ctx, sess.ID, "started", map[string]any{"org_id": orgID})

	// Invalidate tenant cache so next request sees sandbox membership.
	_ = s.tenant.SetActiveOrganization(ctx, userID, orgID)

	return s.toDTO(ctx, sess)
}

func (s *Service) seedProductDemo(ctx context.Context, orgID, userID string) error {
	return s.repo.SeedProductDemoModule(ctx, orgID, userID)
}

func (s *Service) Restart(ctx context.Context, userID string) (*dto.SessionDTO, error) {
	if existing, _ := s.repo.ActiveSession(ctx, userID, DefaultWorkflow); existing != nil {
		_, _ = s.Cleanup(ctx, userID, existing.ID, false)
	}
	// Cleanup marks cleaned; Start only reuses status=active.
	return s.Start(ctx, userID)
}

func (s *Service) Validate(ctx context.Context, userID, sessionID, stepKey string, visitedRoute string, clientEvent *dto.ClientEvent) (*dto.ValidateResult, error) {
	sess, err := s.repo.GetSession(ctx, sessionID, userID)
	if err != nil || sess == nil {
		return nil, ErrNotFound
	}
	if sess.Status != "active" {
		return &dto.ValidateResult{OK: false, Message: "Session is not active"}, nil
	}
	steps, err := s.repo.ListSteps(ctx, sess.WorkflowKey)
	if err != nil {
		return nil, err
	}
	var step *repository.Step
	for i := range steps {
		if steps[i].Key == stepKey {
			step = &steps[i]
			break
		}
	}
	if step == nil {
		return nil, ErrNotFound
	}

	orgID := ""
	if sess.SandboxOrganizationID != nil {
		orgID = *sess.SandboxOrganizationID
	}

	_ = s.repo.SetStepStatus(ctx, sess.ID, stepKey, "validating", "")
	route := visitedRoute
	if clientEvent != nil && clientEvent.Path != "" && route == "" {
		route = clientEvent.Path
	}

	ok, msg := s.runValidator(ctx, orgID, step, route, clientEvent)
	if !ok {
		if step.FailureMessage != "" {
			msg = step.FailureMessage + " — " + msg
		}
		_ = s.repo.SetStepStatus(ctx, sess.ID, stepKey, "failed", msg)
		_ = s.repo.LogEvent(ctx, sess.ID, "validate_failed", map[string]any{"step": stepKey, "error": msg})
		return &dto.ValidateResult{OK: false, Message: msg}, nil
	}

	_ = s.repo.SetStepStatus(ctx, sess.ID, stepKey, "completed", "")
	_ = s.repo.LogEvent(ctx, sess.ID, "step_completed", map[string]any{"step": stepKey})
	s.trackValidatedResource(ctx, sess.ID, orgID, step)

	nextKey := nextStepKey(steps, stepKey)
	if nextKey != "" {
		_ = s.repo.SetStepStatus(ctx, sess.ID, nextKey, "active", "")
		sess.CurrentStepKey = &nextKey
	} else {
		now := time.Now().UTC()
		sess.Status = "completed"
		sess.CompletedAt = &now
	}
	_ = s.repo.UpdateSession(ctx, sess)

	dtoSess, err := s.toDTO(ctx, sess)
	if err != nil {
		return nil, err
	}
	return &dto.ValidateResult{OK: true, Message: "Step completed", Session: dtoSess}, nil
}

func (s *Service) trackValidatedResource(ctx context.Context, sessionID, orgID string, step *repository.Step) {
	params := map[string]any{}
	_ = json.Unmarshal(step.ValidatorParams, &params)
	switch step.ValidatorKey {
	case "module_exists":
		api, _ := params["api_name"].(string)
		if id, err := s.repo.ModuleIDByAPIName(ctx, orgID, api); err == nil && id != "" {
			_ = s.repo.TrackResource(ctx, sessionID, "module", id, map[string]any{"api_name": api})
		}
	case "field_exists":
		mod, _ := params["module_api_name"].(string)
		field, _ := params["api_name"].(string)
		_ = s.repo.TrackResource(ctx, sessionID, "field", uuidOrSession(sessionID), map[string]any{
			"module_api_name": mod, "api_name": field,
		})
	}
}

// uuidOrSession returns a stable placeholder UUID derived from the session when
// the resource has no dedicated id yet (fields are tracked by meta).
func uuidOrSession(sessionID string) string {
	if sessionID != "" {
		return sessionID
	}
	return "00000000-0000-0000-0000-000000000000"
}

func (s *Service) Skip(ctx context.Context, userID, sessionID, stepKey string) (*dto.SessionDTO, error) {
	sess, err := s.repo.GetSession(ctx, sessionID, userID)
	if err != nil || sess == nil {
		return nil, ErrNotFound
	}
	steps, _ := s.repo.ListSteps(ctx, sess.WorkflowKey)
	var step *repository.Step
	for i := range steps {
		if steps[i].Key == stepKey {
			step = &steps[i]
			break
		}
	}
	if step == nil {
		return nil, ErrNotFound
	}
	if !step.IsSkippable {
		return nil, fmt.Errorf("step cannot be skipped — complete the required action")
	}
	_ = s.repo.SetStepStatus(ctx, sess.ID, stepKey, "skipped", "")
	_ = s.repo.LogEvent(ctx, sess.ID, "step_skipped", map[string]any{"step": stepKey})
	nextKey := nextStepKey(steps, stepKey)
	if nextKey != "" {
		_ = s.repo.SetStepStatus(ctx, sess.ID, nextKey, "active", "")
		sess.CurrentStepKey = &nextKey
	} else {
		now := time.Now().UTC()
		sess.Status = "completed"
		sess.CompletedAt = &now
	}
	_ = s.repo.UpdateSession(ctx, sess)
	return s.toDTO(ctx, sess)
}

func (s *Service) Complete(ctx context.Context, userID, sessionID string) (*dto.SessionDTO, error) {
	sess, err := s.repo.GetSession(ctx, sessionID, userID)
	if err != nil || sess == nil {
		return nil, ErrNotFound
	}
	steps, err := s.repo.ListSteps(ctx, sess.WorkflowKey)
	if err != nil {
		return nil, err
	}
	progress, err := s.repo.ListProgress(ctx, sess.ID)
	if err != nil {
		return nil, err
	}
	for _, st := range steps {
		if st.IsSkippable {
			continue
		}
		status := progress[st.Key]
		if status != "completed" && status != "skipped" {
			return nil, fmt.Errorf("complete required steps first (pending: %s)", st.Key)
		}
	}
	now := time.Now().UTC()
	sess.Status = "completed"
	sess.CompletedAt = &now
	stats := map[string]any{
		"completed_at": now,
		"message":      "Walkthrough finished",
	}
	raw, _ := json.Marshal(stats)
	sess.Stats = raw
	_ = s.repo.UpdateSession(ctx, sess)
	_ = s.repo.LogEvent(ctx, sess.ID, "completed", stats)
	return s.toDTO(ctx, sess)
}

func (s *Service) Cleanup(ctx context.Context, userID, sessionID string, keepData bool) (*dto.SessionDTO, error) {
	sess, err := s.repo.GetSession(ctx, sessionID, userID)
	if err != nil || sess == nil {
		return nil, ErrNotFound
	}
	keep := keepData
	sess.KeepData = &keep

	if keepData {
		// Stay on the sandbox org so the user can explore freely.
		if sess.SandboxOrganizationID != nil && *sess.SandboxOrganizationID != "" {
			_ = s.tenant.SetActiveOrganization(ctx, userID, *sess.SandboxOrganizationID)
		}
	} else {
		if sess.PreviousOrganizationID != nil && *sess.PreviousOrganizationID != "" {
			_ = s.tenant.SetActiveOrganization(ctx, userID, *sess.PreviousOrganizationID)
		}
		if sess.SandboxOrganizationID != nil {
			_ = s.repo.DeleteOrganization(ctx, *sess.SandboxOrganizationID)
		}
	}

	sess.Status = "cleaned"
	_ = s.repo.UpdateSession(ctx, sess)
	_ = s.repo.LogEvent(ctx, sess.ID, "cleaned", map[string]any{"keep_data": keepData})
	return s.toDTO(ctx, sess)
}

func (s *Service) LogClientEvent(ctx context.Context, userID, sessionID, eventType string, payload json.RawMessage) error {
	sess, err := s.repo.GetSession(ctx, sessionID, userID)
	if err != nil || sess == nil {
		return ErrNotFound
	}
	if eventType == "abandoned" && sess.Status == "active" {
		sess.Status = "abandoned"
		_ = s.repo.UpdateSession(ctx, sess)
	}
	var body map[string]any
	if len(payload) > 0 {
		_ = json.Unmarshal(payload, &body)
	}
	if body == nil {
		body = map[string]any{}
	}
	return s.repo.LogEvent(ctx, sessionID, eventType, body)
}

func (s *Service) runValidator(ctx context.Context, orgID string, step *repository.Step, visitedRoute string, clientEvent *dto.ClientEvent) (bool, string) {
	params := map[string]any{}
	_ = json.Unmarshal(step.ValidatorParams, &params)

	switch step.ValidatorKey {
	case "none", "acknowledge":
		return true, "ok"
	case "route_visited":
		want, _ := params["route"].(string)
		if visitedRoute == "" {
			return false, "Open the suggested page, then continue"
		}
		if want != "" && !strings.HasPrefix(visitedRoute, want) {
			return false, fmt.Sprintf("Expected to visit %s (current: %s)", want, visitedRoute)
		}
		return true, "Route confirmed"
	case "ui_click":
		want, _ := params["selector"].(string)
		if clientEvent == nil || clientEvent.Selector == "" {
			return false, "Click the highlighted control to continue"
		}
		if want != "" && !strings.Contains(clientEvent.Selector, strings.Trim(want, "[]")) && clientEvent.Selector != want {
			return false, fmt.Sprintf("Expected click on %s", want)
		}
		return true, "Click recorded"
	case "module_exists":
		api, _ := params["api_name"].(string)
		ok, err := s.repo.ModuleExists(ctx, orgID, api)
		if err != nil || !ok {
			return false, fmt.Sprintf("Module %q not found yet — create it in Settings → Modules", api)
		}
		return true, "Module found"
	case "field_exists":
		mod, _ := params["module_api_name"].(string)
		field, _ := params["api_name"].(string)
		ok, err := s.repo.FieldExists(ctx, orgID, mod, field)
		if err != nil || !ok {
			return false, fmt.Sprintf("Field %q on %q not found — add it in Settings → Fields", field, mod)
		}
		return true, "Field found"
	case "record_exists":
		mod, _ := params["module_api_name"].(string)
		ok, err := s.repo.RecordExists(ctx, orgID, mod)
		if err != nil || !ok {
			return false, "No records yet — create one via Forms or Tables"
		}
		return true, "Record found"
	case "note_exists":
		mod, _ := params["module_api_name"].(string)
		ok, err := s.repo.NoteExists(ctx, orgID, mod)
		if err != nil || !ok {
			return false, "Add a note on a record in the Record Workspace Notes tab"
		}
		return true, "Note found"
	case "activity_exists":
		mod, _ := params["module_api_name"].(string)
		ok, err := s.repo.ActivityExists(ctx, orgID, mod)
		if err != nil || !ok {
			return false, "Timeline is empty — create/update a record or add a note first"
		}
		return true, "Activity found"
	case "workspace_opened":
		mod, _ := params["module_api_name"].(string)
		ok, err := s.repo.RecordExists(ctx, orgID, mod)
		if err != nil || !ok {
			return false, "Create a record first, then open it from Tables"
		}
		path := visitedRoute
		if clientEvent != nil && clientEvent.Path != "" {
			path = clientEvent.Path
		}
		if path != "" && strings.HasPrefix(path, "/tables/") {
			parts := strings.Split(strings.Trim(path, "/"), "/")
			if len(parts) >= 3 {
				return true, "Record workspace opened"
			}
		}
		// Accept if records exist and user confirms (Tables list is enough for soft step).
		return true, "Ready — open a record row when you can"
	case "import_completed", "export_completed", "automation_triggered":
		// Curriculum stubs — succeed when route visited for now if soft, else fail closed.
		if visitedRoute != "" {
			return true, "Observed"
		}
		return false, "Complete the action on the suggested page"
	default:
		return false, fmt.Sprintf("Unknown validator %q", step.ValidatorKey)
	}
}

func nextStepKey(steps []repository.Step, current string) string {
	for i, st := range steps {
		if st.Key == current && i+1 < len(steps) {
			return steps[i+1].Key
		}
	}
	return ""
}

func (s *Service) toDTO(ctx context.Context, sess *repository.Session) (*dto.SessionDTO, error) {
	steps, err := s.repo.ListSteps(ctx, sess.WorkflowKey)
	if err != nil {
		return nil, err
	}
	progress, err := s.repo.ListProgress(ctx, sess.ID)
	if err != nil {
		return nil, err
	}
	outSteps := make([]dto.StepDTO, 0, len(steps))
	completed := 0
	for _, st := range steps {
		status := progress[st.Key]
		if status == "" {
			status = "locked"
		}
		if status == "completed" || status == "skipped" {
			completed++
		}
		outSteps = append(outSteps, dto.StepDTO{
			Key: st.Key, SortOrder: st.SortOrder, Title: st.Title,
			Description: st.Description, WhyItMatters: st.Why, HowItWorks: st.How,
			ExpectedResult: st.Expected, Route: st.Route, TargetSelector: st.Target,
			ActionLabel: st.ActionLabel, ValidatorKey: st.ValidatorKey,
			ValidatorParams: st.ValidatorParams, IsSkippable: st.IsSkippable, Status: status,
			RequiredAction: st.RequiredAction, SuccessEvent: st.SuccessEvent,
			FailureMessage: st.FailureMessage, Hint: st.Hint, MaxAttempts: st.MaxAttempts,
			AllowSelectors: st.AllowSelectors, Placement: st.Placement,
		})
	}
	pct := 0
	if len(steps) > 0 {
		pct = (completed * 100) / len(steps)
	}
	return &dto.SessionDTO{
		ID: sess.ID, WorkflowKey: sess.WorkflowKey, WorkflowVersion: sess.WorkflowVersion,
		SandboxOrganizationID: sess.SandboxOrganizationID, Status: sess.Status,
		CurrentStepKey: sess.CurrentStepKey, StartedAt: sess.StartedAt,
		CompletedAt: sess.CompletedAt, Stats: sess.Stats, Steps: outSteps,
		ProgressPercent: pct,
	}, nil
}
