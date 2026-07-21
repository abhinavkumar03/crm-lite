package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"strings"
	"time"

	fieldentity "github.com/abhinavkumar03/crm-lite/backend/internal/field/entity"
	"github.com/abhinavkumar03/crm-lite/backend/internal/jobs"
	"github.com/abhinavkumar03/crm-lite/backend/internal/workflow/catalog"
	"github.com/abhinavkumar03/crm-lite/backend/internal/workflow/dto"
	"github.com/abhinavkumar03/crm-lite/backend/internal/workflow/entity"
	"github.com/abhinavkumar03/crm-lite/backend/internal/workflow/repository"
)

var (
	ErrNotFound      = errors.New("workflow not found")
	ErrInvalidInput  = errors.New("invalid workflow definition")
	ErrNotPublishable = errors.New("workflow cannot be published")
)

type FieldReader interface {
	List(ctx context.Context, orgID, moduleID string) ([]fieldentity.Field, error)
}

type Producer interface {
	Publish(ctx context.Context, job jobs.Job, opts ...interface{ /* placeholder */ }) error
}

// JobPublisher abstracts asynq producer without importing option types into tests.
type JobPublisher interface {
	Publish(ctx context.Context, job jobs.Job, opts ...interface{}) error
}

// AsynqPublisher wraps jobs.Producer.
type AsynqPublisher struct {
	P *jobs.Producer
}

func (a *AsynqPublisher) Publish(ctx context.Context, job jobs.Job, _ ...interface{}) error {
	if a == nil || a.P == nil {
		return nil
	}
	return a.P.Publish(ctx, job)
}

type Service struct {
	repo     *repository.Repository
	fields   FieldReader
	publisher JobPublisher
	enabled  bool
}

func New(repo *repository.Repository, fields FieldReader, publisher JobPublisher, enabled bool) *Service {
	return &Service{repo: repo, fields: fields, publisher: publisher, enabled: enabled}
}

func (s *Service) Enabled() bool { return s.enabled }

func (s *Service) Create(ctx context.Context, orgID, userID string, req dto.CreateWorkflowRequest) (*dto.WorkflowResponse, error) {
	if err := validateDefinition(req.Triggers, req.Conditions, req.Actions); err != nil {
		return nil, err
	}
	onErr := req.OnActionError
	if onErr == "" {
		onErr = entity.OnErrorStop
	}
	priority := req.Priority
	if priority == 0 {
		priority = 100
	}
	w := &entity.Workflow{
		OrganizationID: orgID,
		ModuleID:       req.ModuleID,
		Name:           strings.TrimSpace(req.Name),
		Description:    req.Description,
		Status:         entity.StatusDraft,
		OnActionError:  onErr,
		Priority:       priority,
		CreatedBy:      &userID,
		UpdatedBy:      &userID,
	}
	if err := s.repo.CreateWorkflow(ctx, w); err != nil {
		return nil, err
	}
	v := &entity.Version{
		WorkflowID:         w.ID,
		OrganizationID:     orgID,
		Version:            1,
		State:              entity.VersionDraft,
		DefinitionSnapshot: json.RawMessage(`{}`),
	}
	if err := s.repo.CreateVersion(ctx, v); err != nil {
		return nil, err
	}
	if err := s.persistDefinition(ctx, orgID, v.ID, req.Triggers, req.Conditions, req.Actions); err != nil {
		return nil, err
	}
	_ = s.refreshSnapshot(ctx, orgID, v.ID)
	return s.Get(ctx, orgID, w.ID)
}

func (s *Service) Update(ctx context.Context, orgID, id, userID string, req dto.UpdateWorkflowRequest) (*dto.WorkflowResponse, error) {
	w, err := s.repo.GetWorkflow(ctx, orgID, id)
	if err != nil {
		return nil, err
	}
	if w == nil {
		return nil, ErrNotFound
	}
	if w.Status == entity.StatusArchived {
		return nil, fmt.Errorf("%w: archived", ErrInvalidInput)
	}
	if req.Name != nil {
		w.Name = strings.TrimSpace(*req.Name)
	}
	if req.Description != nil {
		w.Description = *req.Description
	}
	if req.ModuleID != nil {
		w.ModuleID = req.ModuleID
	}
	if req.OnActionError != nil && *req.OnActionError != "" {
		w.OnActionError = *req.OnActionError
	}
	if req.Priority != nil {
		w.Priority = *req.Priority
	}
	w.UpdatedBy = &userID
	if err := s.repo.UpdateWorkflowHeader(ctx, w); err != nil {
		return nil, err
	}

	hasDef := req.Triggers != nil || req.Conditions != nil || req.Actions != nil
	if hasDef {
		if err := validateDefinition(req.Triggers, req.Conditions, req.Actions); err != nil {
			return nil, err
		}
		draft, err := s.ensureDraftVersion(ctx, orgID, w)
		if err != nil {
			return nil, err
		}
		if err := s.repo.DeleteVersionChildren(ctx, draft.ID); err != nil {
			return nil, err
		}
		if err := s.persistDefinition(ctx, orgID, draft.ID, req.Triggers, req.Conditions, req.Actions); err != nil {
			return nil, err
		}
		_ = s.refreshSnapshot(ctx, orgID, draft.ID)
	}
	return s.Get(ctx, orgID, id)
}

func (s *Service) ensureDraftVersion(ctx context.Context, orgID string, w *entity.Workflow) (*entity.Version, error) {
	draft, err := s.repo.LatestDraftVersion(ctx, orgID, w.ID)
	if err != nil {
		return nil, err
	}
	if draft != nil {
		return draft, nil
	}
	max, err := s.repo.MaxVersion(ctx, w.ID)
	if err != nil {
		return nil, err
	}
	// Copy published definition into a new draft when editing an active workflow.
	v := &entity.Version{
		WorkflowID:         w.ID,
		OrganizationID:     orgID,
		Version:            max + 1,
		State:              entity.VersionDraft,
		DefinitionSnapshot: json.RawMessage(`{}`),
	}
	if err := s.repo.CreateVersion(ctx, v); err != nil {
		return nil, err
	}
	if w.PublishedVersionID != nil {
		if err := s.cloneDefinition(ctx, orgID, *w.PublishedVersionID, v.ID); err != nil {
			return nil, err
		}
		_ = s.refreshSnapshot(ctx, orgID, v.ID)
	}
	return v, nil
}

func (s *Service) cloneDefinition(ctx context.Context, orgID, fromVersionID, toVersionID string) error {
	tr, err := s.repo.ListTriggers(ctx, fromVersionID)
	if err != nil {
		return err
	}
	for _, t := range tr {
		nt := t
		nt.ID = ""
		nt.VersionID = toVersionID
		nt.OrganizationID = orgID
		if err := s.repo.InsertTrigger(ctx, &nt); err != nil {
			return err
		}
	}
	conds, err := s.repo.ListConditions(ctx, fromVersionID)
	if err != nil {
		return err
	}
	idMap := map[string]string{}
	// Insert roots first then children by repeating passes.
	remaining := append([]entity.Condition(nil), conds...)
	for len(remaining) > 0 {
		progress := false
		next := remaining[:0]
		for _, c := range remaining {
			if c.ParentID != nil {
				if _, ok := idMap[*c.ParentID]; !ok {
					next = append(next, c)
					continue
				}
			}
			nc := c
			nc.ID = ""
			nc.VersionID = toVersionID
			nc.OrganizationID = orgID
			if c.ParentID != nil {
				mapped := idMap[*c.ParentID]
				nc.ParentID = &mapped
			}
			oldID := c.ID
			if err := s.repo.InsertCondition(ctx, &nc); err != nil {
				return err
			}
			idMap[oldID] = nc.ID
			progress = true
		}
		if !progress {
			return fmt.Errorf("%w: condition parent cycle", ErrInvalidInput)
		}
		remaining = next
	}
	acts, err := s.repo.ListActions(ctx, fromVersionID)
	if err != nil {
		return err
	}
	for _, a := range acts {
		na := a
		na.ID = ""
		na.VersionID = toVersionID
		na.OrganizationID = orgID
		if err := s.repo.InsertAction(ctx, &na); err != nil {
			return err
		}
	}
	return nil
}

func (s *Service) Get(ctx context.Context, orgID, id string) (*dto.WorkflowResponse, error) {
	w, err := s.repo.GetWorkflow(ctx, orgID, id)
	if err != nil {
		return nil, err
	}
	if w == nil {
		return nil, ErrNotFound
	}
	versionID := ""
	versionNum := 0
	var draftID *string
	if d, _ := s.repo.LatestDraftVersion(ctx, orgID, id); d != nil {
		versionID = d.ID
		versionNum = d.Version
		draftID = &d.ID
	} else if w.PublishedVersionID != nil {
		if v, _ := s.repo.GetVersion(ctx, orgID, *w.PublishedVersionID); v != nil {
			versionID = v.ID
			versionNum = v.Version
		}
	}
	resp := &dto.WorkflowResponse{
		ID:                 w.ID,
		ModuleID:           w.ModuleID,
		Name:               w.Name,
		Description:        w.Description,
		Status:             w.Status,
		OnActionError:      w.OnActionError,
		Priority:           w.Priority,
		PublishedVersionID: w.PublishedVersionID,
		DraftVersionID:     draftID,
		Version:            versionNum,
		CreatedBy:          w.CreatedBy,
		UpdatedBy:          w.UpdatedBy,
		CreatedAt:          w.CreatedAt,
		UpdatedAt:          w.UpdatedAt,
		Triggers:           []dto.TriggerResponse{},
		Actions:            []dto.ActionResponse{},
	}
	if w.ModuleID != nil {
		if name, err := s.repo.ModuleAPIName(ctx, orgID, *w.ModuleID); err == nil {
			resp.ModuleAPIName = &name
		}
	}
	if versionID == "" {
		return resp, nil
	}
	tr, _ := s.repo.ListTriggers(ctx, versionID)
	for _, t := range tr {
		resp.Triggers = append(resp.Triggers, dto.TriggerResponse{
			ID: t.ID, Type: t.Type, Config: rawToMap(t.Config),
		})
	}
	conds, _ := s.repo.ListConditions(ctx, versionID)
	if tree := buildConditionTree(conds); tree != nil {
		resp.Conditions = tree
	}
	acts, _ := s.repo.ListActions(ctx, versionID)
	for _, a := range acts {
		resp.Actions = append(resp.Actions, dto.ActionResponse{
			ID: a.ID, SortOrder: a.SortOrder, Type: a.Type, Config: rawToMap(a.Config),
			MaxRetries: a.MaxRetries, ContinueOnError: a.ContinueOnError,
		})
	}
	return resp, nil
}

func (s *Service) List(ctx context.Context, orgID string, status, moduleID string, page, pageSize int) (*dto.ListWorkflowsResult, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}
	items, total, err := s.repo.ListWorkflows(ctx, orgID, status, moduleID, page, pageSize)
	if err != nil {
		return nil, err
	}
	out := make([]dto.WorkflowSummary, 0, len(items))
	for _, w := range items {
		sum := dto.WorkflowSummary{
			ID: w.ID, ModuleID: w.ModuleID, Name: w.Name, Description: w.Description,
			Status: w.Status, Priority: w.Priority, CreatedAt: w.CreatedAt, UpdatedAt: w.UpdatedAt,
			TriggerTypes: []string{},
		}
		if w.ModuleID != nil {
			if name, err := s.repo.ModuleAPIName(ctx, orgID, *w.ModuleID); err == nil {
				sum.ModuleAPIName = &name
			}
		}
		vid := ""
		if d, _ := s.repo.LatestDraftVersion(ctx, orgID, w.ID); d != nil {
			vid = d.ID
			sum.Version = d.Version
		} else if w.PublishedVersionID != nil {
			if v, _ := s.repo.GetVersion(ctx, orgID, *w.PublishedVersionID); v != nil {
				vid = v.ID
				sum.Version = v.Version
			}
		}
		if vid != "" {
			tr, _ := s.repo.ListTriggers(ctx, vid)
			for _, t := range tr {
				sum.TriggerTypes = append(sum.TriggerTypes, t.Type)
			}
			acts, _ := s.repo.ListActions(ctx, vid)
			sum.ActionCount = len(acts)
		}
		out = append(out, sum)
	}
	return &dto.ListWorkflowsResult{
		Items: out, Page: page, PageSize: pageSize, Total: total,
		TotalPages: int(math.Max(1, math.Ceil(float64(total)/float64(pageSize)))),
	}, nil
}

func (s *Service) Publish(ctx context.Context, orgID, id, userID string, req dto.PublishRequest) (*dto.WorkflowResponse, error) {
	w, err := s.repo.GetWorkflow(ctx, orgID, id)
	if err != nil {
		return nil, err
	}
	if w == nil {
		return nil, ErrNotFound
	}
	draft, err := s.repo.LatestDraftVersion(ctx, orgID, id)
	if err != nil {
		return nil, err
	}
	if draft == nil {
		if w.PublishedVersionID != nil {
			// Re-activate without new draft.
			w.Status = entity.StatusActive
			w.UpdatedBy = &userID
			if err := s.repo.UpdateWorkflowHeader(ctx, w); err != nil {
				return nil, err
			}
			return s.Get(ctx, orgID, id)
		}
		return nil, ErrNotPublishable
	}
	tr, _ := s.repo.ListTriggers(ctx, draft.ID)
	acts, _ := s.repo.ListActions(ctx, draft.ID)
	if len(tr) == 0 || len(acts) == 0 {
		return nil, fmt.Errorf("%w: need at least one trigger and one action", ErrNotPublishable)
	}
	now := time.Now().UTC()
	draft.State = entity.VersionPublished
	draft.PublishedAt = &now
	draft.PublishedBy = &userID
	draft.Changelog = req.Changelog
	_ = s.refreshSnapshot(ctx, orgID, draft.ID)
	if err := s.repo.UpdateVersion(ctx, draft); err != nil {
		return nil, err
	}
	w.Status = entity.StatusActive
	w.PublishedVersionID = &draft.ID
	w.UpdatedBy = &userID
	if err := s.repo.UpdateWorkflowHeader(ctx, w); err != nil {
		return nil, err
	}
	return s.Get(ctx, orgID, id)
}

func (s *Service) Disable(ctx context.Context, orgID, id, userID string) (*dto.WorkflowResponse, error) {
	w, err := s.repo.GetWorkflow(ctx, orgID, id)
	if err != nil {
		return nil, err
	}
	if w == nil {
		return nil, ErrNotFound
	}
	w.Status = entity.StatusDisabled
	w.UpdatedBy = &userID
	if err := s.repo.UpdateWorkflowHeader(ctx, w); err != nil {
		return nil, err
	}
	return s.Get(ctx, orgID, id)
}

func (s *Service) Archive(ctx context.Context, orgID, id, userID string) error {
	w, err := s.repo.GetWorkflow(ctx, orgID, id)
	if err != nil {
		return err
	}
	if w == nil {
		return ErrNotFound
	}
	w.Status = entity.StatusArchived
	w.UpdatedBy = &userID
	return s.repo.UpdateWorkflowHeader(ctx, w)
}

func (s *Service) ListVersions(ctx context.Context, orgID, workflowID string) ([]dto.VersionSummary, error) {
	w, err := s.repo.GetWorkflow(ctx, orgID, workflowID)
	if err != nil {
		return nil, err
	}
	if w == nil {
		return nil, ErrNotFound
	}
	vers, err := s.repo.ListVersions(ctx, orgID, workflowID)
	if err != nil {
		return nil, err
	}
	out := make([]dto.VersionSummary, 0, len(vers))
	for _, v := range vers {
		out = append(out, dto.VersionSummary{
			ID: v.ID, Version: v.Version, State: v.State, Changelog: v.Changelog,
			PublishedAt: v.PublishedAt, PublishedBy: v.PublishedBy, CreatedAt: v.CreatedAt,
		})
	}
	return out, nil
}

func (s *Service) Rollback(ctx context.Context, orgID, workflowID, versionID, userID string) (*dto.WorkflowResponse, error) {
	w, err := s.repo.GetWorkflow(ctx, orgID, workflowID)
	if err != nil {
		return nil, err
	}
	if w == nil {
		return nil, ErrNotFound
	}
	v, err := s.repo.GetVersion(ctx, orgID, versionID)
	if err != nil {
		return nil, err
	}
	if v == nil || v.WorkflowID != workflowID {
		return nil, ErrNotFound
	}
	max, err := s.repo.MaxVersion(ctx, workflowID)
	if err != nil {
		return nil, err
	}
	draft := &entity.Version{
		WorkflowID:         workflowID,
		OrganizationID:     orgID,
		Version:            max + 1,
		State:              entity.VersionDraft,
		DefinitionSnapshot: v.DefinitionSnapshot,
		Changelog:          fmt.Sprintf("Rollback from v%d", v.Version),
	}
	if err := s.repo.CreateVersion(ctx, draft); err != nil {
		return nil, err
	}
	if err := s.cloneDefinition(ctx, orgID, versionID, draft.ID); err != nil {
		return nil, err
	}
	now := time.Now().UTC()
	draft.State = entity.VersionPublished
	draft.PublishedAt = &now
	draft.PublishedBy = &userID
	_ = s.refreshSnapshot(ctx, orgID, draft.ID)
	if err := s.repo.UpdateVersion(ctx, draft); err != nil {
		return nil, err
	}
	// Mark source as rolled_back if it was published.
	if v.State == entity.VersionPublished {
		v.State = entity.VersionRolledBack
		_ = s.repo.UpdateVersion(ctx, v)
	}
	w.Status = entity.StatusActive
	w.PublishedVersionID = &draft.ID
	w.UpdatedBy = &userID
	if err := s.repo.UpdateWorkflowHeader(ctx, w); err != nil {
		return nil, err
	}
	return s.Get(ctx, orgID, workflowID)
}

func (s *Service) ListExecutions(ctx context.Context, orgID, workflowID, moduleID, recordID, status string, page, pageSize int) (*dto.ListExecutionsResult, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}
	items, names, total, err := s.repo.ListExecutions(ctx, orgID, workflowID, moduleID, recordID, status, page, pageSize)
	if err != nil {
		return nil, err
	}
	out := make([]dto.ExecutionSummary, 0, len(items))
	for _, e := range items {
		out = append(out, toExecSummary(e, names[e.ID]))
	}
	return &dto.ListExecutionsResult{
		Items: out, Page: page, PageSize: pageSize, Total: total,
		TotalPages: int(math.Max(1, math.Ceil(float64(total)/float64(pageSize)))),
	}, nil
}

func (s *Service) GetExecution(ctx context.Context, orgID, id string) (*dto.ExecutionDetail, error) {
	e, name, err := s.repo.GetExecution(ctx, orgID, id)
	if err != nil {
		return nil, err
	}
	if e == nil {
		return nil, ErrNotFound
	}
	steps, err := s.repo.ListExecutionSteps(ctx, id)
	if err != nil {
		return nil, err
	}
	detail := &dto.ExecutionDetail{ExecutionSummary: toExecSummary(*e, name)}
	for _, st := range steps {
		detail.Steps = append(detail.Steps, dto.ExecutionStepResponse{
			ID: st.ID, SortOrder: st.SortOrder, ActionType: st.ActionType, ActionID: st.ActionID,
			Status: st.Status, Input: rawToMap(st.Input), Output: rawToMap(st.Output),
			Error: st.Error, StartedAt: st.StartedAt, FinishedAt: st.FinishedAt,
		})
	}
	if detail.Steps == nil {
		detail.Steps = []dto.ExecutionStepResponse{}
	}
	return detail, nil
}

func (s *Service) RetryExecution(ctx context.Context, orgID, executionID, userID string) (*dto.ExecutionSummary, error) {
	e, name, err := s.repo.GetExecution(ctx, orgID, executionID)
	if err != nil {
		return nil, err
	}
	if e == nil {
		return nil, ErrNotFound
	}
	if e.Status != entity.ExecFailed && e.Status != entity.ExecPartial {
		return nil, fmt.Errorf("%w: only failed or partial runs can be retried", ErrInvalidInput)
	}
	if e.RecordID == nil || e.ModuleID == nil {
		return nil, fmt.Errorf("%w: execution missing record context", ErrInvalidInput)
	}
	w, err := s.repo.GetWorkflow(ctx, orgID, e.WorkflowID)
	if err != nil {
		return nil, err
	}
	if w == nil || w.Status != entity.StatusActive || w.PublishedVersionID == nil {
		return nil, fmt.Errorf("%w: workflow is not active", ErrInvalidInput)
	}
	payload := map[string]interface{}{
		"org_id":      orgID,
		"module_id":   *e.ModuleID,
		"record_id":   *e.RecordID,
		"trigger":     entity.TriggerManual,
		"source":      entity.SourceManual,
		"depth":       0,
		"workflow_id": e.WorkflowID,
		"retry_of":    executionID,
	}
	if s.publisher != nil {
		if err := s.publisher.Publish(ctx, jobs.Job{Type: jobs.JobWorkflowEvaluate, UserID: userID, Payload: payload}); err != nil {
			return nil, err
		}
	}
	sum := toExecSummary(*e, name)
	return &sum, nil
}

func (s *Service) ManualRun(ctx context.Context, orgID, workflowID, userID string, req dto.ManualRunRequest) error {
	w, err := s.repo.GetWorkflow(ctx, orgID, workflowID)
	if err != nil {
		return err
	}
	if w == nil {
		return ErrNotFound
	}
	if w.Status != entity.StatusActive || w.PublishedVersionID == nil {
		return fmt.Errorf("%w: workflow not active", ErrInvalidInput)
	}
	moduleID := ""
	if w.ModuleID != nil {
		moduleID = *w.ModuleID
	}
	// Prefer an explicit module from the client when the workflow is module-agnostic.
	if moduleID == "" && req.ModuleID != "" {
		moduleID = req.ModuleID
	}
	if w.ModuleID != nil && req.ModuleID != "" && *w.ModuleID != req.ModuleID {
		return fmt.Errorf("%w: workflow module mismatch", ErrInvalidInput)
	}
	if moduleID == "" {
		return fmt.Errorf("%w: module_id required for manual run", ErrInvalidInput)
	}
	payload := map[string]interface{}{
		"org_id":      orgID,
		"module_id":   moduleID,
		"record_id":   req.RecordID,
		"trigger":     entity.TriggerManual,
		"source":      entity.SourceManual,
		"depth":       0,
		"workflow_id": workflowID,
	}
	if s.publisher != nil {
		return s.publisher.Publish(ctx, jobs.Job{Type: jobs.JobWorkflowEvaluate, UserID: userID, Payload: payload})
	}
	return nil
}

func (s *Service) EnsureBuiltinTemplates(ctx context.Context) error {
	for _, bt := range catalog.BuiltinTemplates() {
		def, _ := json.Marshal(bt.Definition)
		// Embed category in definition metadata for API responses.
		var defMap map[string]any
		_ = json.Unmarshal(def, &defMap)
		if defMap == nil {
			defMap = map[string]any{}
		}
		defMap["_category"] = bt.Category
		def, _ = json.Marshal(defMap)
		mod := bt.ModuleAPIName
		t := &entity.Template{
			Key: bt.Key, Name: bt.Name, Description: bt.Description,
			ModuleAPIName: &mod, Definition: def, IsBuiltin: true,
		}
		if err := s.repo.UpsertBuiltinTemplate(ctx, t); err != nil {
			return err
		}
	}
	return nil
}

func (s *Service) ListTemplates(ctx context.Context, orgID string) ([]dto.TemplateResponse, error) {
	_ = s.EnsureBuiltinTemplates(ctx)
	items, err := s.repo.ListTemplates(ctx, orgID)
	if err != nil {
		return nil, err
	}
	out := make([]dto.TemplateResponse, 0, len(items))
	for _, t := range items {
		def := rawToMap(t.Definition)
		cat, _ := def["_category"].(string)
		delete(def, "_category")
		out = append(out, dto.TemplateResponse{
			ID: t.ID, Key: t.Key, Name: t.Name, Description: t.Description,
			ModuleAPIName: t.ModuleAPIName, Category: cat, Definition: def, IsBuiltin: t.IsBuiltin,
		})
	}
	return out, nil
}

func (s *Service) CloneTemplate(ctx context.Context, orgID, templateID, userID string) (*dto.WorkflowResponse, error) {
	t, err := s.repo.GetTemplate(ctx, orgID, templateID)
	if err != nil {
		return nil, err
	}
	if t == nil {
		return nil, ErrNotFound
	}
	def := rawToMap(t.Definition)
	delete(def, "_category")
	req := dto.CreateWorkflowRequest{
		Name:          t.Name,
		Description:   t.Description,
		OnActionError: entity.OnErrorStop,
		Priority:      100,
	}
	if t.ModuleAPIName != nil {
		if mid, err := s.repo.ModuleIDByAPIName(ctx, orgID, *t.ModuleAPIName); err == nil {
			req.ModuleID = &mid
		}
	}
	if triggers, ok := def["triggers"].([]any); ok {
		for _, raw := range triggers {
			m, _ := raw.(map[string]any)
			if m == nil {
				continue
			}
			ti := dto.TriggerInput{Type: str(m["type"]), Config: asMap(m["config"])}
			req.Triggers = append(req.Triggers, ti)
		}
	}
	if cond, ok := def["conditions"].(map[string]any); ok {
		c := mapToCondition(cond)
		req.Conditions = &c
	}
	if actions, ok := def["actions"].([]any); ok {
		for _, raw := range actions {
			m, _ := raw.(map[string]any)
			if m == nil {
				continue
			}
			req.Actions = append(req.Actions, dto.ActionInput{
				Type: str(m["type"]), Config: asMap(m["config"]),
			})
		}
	}
	return s.Create(ctx, orgID, userID, req)
}

func (s *Service) BuilderMetadata(ctx context.Context, orgID string) (*dto.BuilderMetadataResponse, error) {
	mods, err := s.repo.ListOrgModules(ctx, orgID)
	if err != nil {
		return nil, err
	}
	builderMods := make([]dto.BuilderModule, 0, len(mods))
	for _, m := range mods {
		bm := dto.BuilderModule{ID: m.ID, APIName: m.APIName, Label: m.Label, Fields: []dto.BuilderField{}}
		if s.fields != nil {
			fields, err := s.fields.List(ctx, orgID, m.ID)
			if err == nil {
				for _, f := range fields {
					bf := dto.BuilderField{
						APIName: f.APIName, Label: f.Label, Type: f.FieldType, Required: f.IsRequired,
					}
					if len(f.Options) > 0 {
						var opts []string
						_ = json.Unmarshal(f.Options, &opts)
						bf.Options = opts
					}
					bm.Fields = append(bm.Fields, bf)
				}
			}
		}
		builderMods = append(builderMods, bm)
	}
	users, _ := s.repo.ListOrgUsers(ctx, orgID)
	bUsers := make([]dto.BuilderUser, 0, len(users))
	for _, u := range users {
		bUsers = append(bUsers, dto.BuilderUser{ID: u.ID, Name: u.Name, Email: u.Email})
	}
	templates, _ := s.ListTemplates(ctx, orgID)
	return &dto.BuilderMetadataResponse{
		Modules:   builderMods,
		Operators: catalogOperators(),
		Actions:   catalogActions(),
		Triggers:  catalogTriggers(),
		Variables: catalogVariables(),
		Users:     bUsers,
		Templates: templates,
	}, nil
}

func (s *Service) Metrics(ctx context.Context, orgID string) (*dto.MetricsResponse, error) {
	a, d, dr, ex, fail, avg, err := s.repo.Metrics(ctx, orgID)
	if err != nil {
		return nil, err
	}
	return &dto.MetricsResponse{
		ActiveWorkflows: a, DisabledWorkflows: d, DraftWorkflows: dr,
		ExecutedToday: ex, FailedToday: fail, AvgDurationMs: avg,
	}, nil
}

// EnqueueRecordEvent is called from the record MutationHook.
func (s *Service) EnqueueRecordEvent(ctx context.Context, orgID, moduleID, recordID, userID, trigger string, before, after map[string]any, changed []string, source string, depth int, excludeWorkflowID string) {
	if !s.enabled || s.publisher == nil {
		return
	}
	if depth > 3 {
		return
	}
	payload := map[string]interface{}{
		"org_id":         orgID,
		"module_id":      moduleID,
		"record_id":      recordID,
		"trigger":        trigger,
		"before":         before,
		"after":          after,
		"changed_fields": changed,
		"source":         source,
		"depth":          depth,
	}
	if excludeWorkflowID != "" {
		payload["exclude_workflow_id"] = excludeWorkflowID
	}
	_ = s.publisher.Publish(ctx, jobs.Job{
		Type: jobs.JobWorkflowEvaluate, UserID: userID, Payload: payload,
	})
}

func (s *Service) persistDefinition(ctx context.Context, orgID, versionID string, triggers []dto.TriggerInput, conditions *dto.ConditionInput, actions []dto.ActionInput) error {
	for _, t := range triggers {
		cfg, _ := json.Marshal(t.Config)
		if t.Config == nil {
			cfg = []byte(`{}`)
		}
		tr := &entity.Trigger{VersionID: versionID, OrganizationID: orgID, Type: t.Type, Config: cfg}
		if err := s.repo.InsertTrigger(ctx, tr); err != nil {
			return err
		}
	}
	if conditions != nil {
		if err := s.insertConditionTree(ctx, orgID, versionID, nil, conditions, 0); err != nil {
			return err
		}
	}
	for i, a := range actions {
		cfg, _ := json.Marshal(a.Config)
		if a.Config == nil {
			cfg = []byte(`{}`)
		}
		act := &entity.Action{
			VersionID: versionID, OrganizationID: orgID, SortOrder: i, Type: a.Type,
			Config: cfg, MaxRetries: a.MaxRetries, ContinueOnError: a.ContinueOnError,
		}
		if err := s.repo.InsertAction(ctx, act); err != nil {
			return err
		}
	}
	return nil
}

func (s *Service) insertConditionTree(ctx context.Context, orgID, versionID string, parentID *string, node *dto.ConditionInput, sort int) error {
	c := &entity.Condition{
		VersionID: versionID, OrganizationID: orgID, ParentID: parentID,
		NodeType: node.NodeType, Logic: node.Logic, FieldAPIName: node.FieldAPIName,
		Operator: node.Operator, SortOrder: sort,
	}
	if node.Value != nil {
		b, _ := json.Marshal(node.Value)
		c.Value = b
	}
	if err := s.repo.InsertCondition(ctx, c); err != nil {
		return err
	}
	for i := range node.Children {
		if err := s.insertConditionTree(ctx, orgID, versionID, &c.ID, &node.Children[i], i); err != nil {
			return err
		}
	}
	return nil
}

func (s *Service) refreshSnapshot(ctx context.Context, orgID, versionID string) error {
	v, err := s.repo.GetVersion(ctx, orgID, versionID)
	if err != nil || v == nil {
		return err
	}
	tr, _ := s.repo.ListTriggers(ctx, versionID)
	conds, _ := s.repo.ListConditions(ctx, versionID)
	acts, _ := s.repo.ListActions(ctx, versionID)
	snap := map[string]any{
		"triggers":   tr,
		"conditions": buildConditionTree(conds),
		"actions":    acts,
	}
	b, _ := json.Marshal(snap)
	v.DefinitionSnapshot = b
	return s.repo.UpdateVersion(ctx, v)
}

func (s *Service) Repo() *repository.Repository { return s.repo }
