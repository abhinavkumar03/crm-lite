package processor

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/hibiken/asynq"
	"go.uber.org/zap"

	"github.com/abhinavkumar03/crm-lite/backend/internal/jobs"
	rdto "github.com/abhinavkumar03/crm-lite/backend/internal/record/dto"
	"github.com/abhinavkumar03/crm-lite/backend/internal/workflow/engine"
	"github.com/abhinavkumar03/crm-lite/backend/internal/workflow/entity"
	"github.com/abhinavkumar03/crm-lite/backend/internal/workflow/repository"
)

// Processor evaluates and runs workflows from asynq jobs.
type Processor struct {
	repo   *repository.Repository
	deps   engine.ActionDeps
	logger *zap.Logger
	pub    *jobs.Producer
}

func New(repo *repository.Repository, deps engine.ActionDeps, pub *jobs.Producer, logger *zap.Logger) *Processor {
	p := &Processor{repo: repo, deps: deps, logger: logger, pub: pub}
	p.deps.Invoker = p
	return p
}

func (p *Processor) EnqueueEvaluate(ctx context.Context, orgID, moduleID, recordID, userID, trigger, source string, depth int, workflowID string, after map[string]any) error {
	if p.pub == nil {
		return nil
	}
	payload := map[string]interface{}{
		"org_id": orgID, "module_id": moduleID, "record_id": recordID,
		"trigger": trigger, "source": source, "depth": depth, "workflow_id": workflowID, "after": after,
	}
	return p.pub.Publish(ctx, jobs.Job{Type: jobs.JobWorkflowEvaluate, UserID: userID, Payload: payload})
}

func (p *Processor) EnqueueResume(ctx context.Context, orgID, executionID, userID string, fromStep int, delayUntil time.Time) error {
	if p.pub == nil {
		return nil
	}
	payload := map[string]interface{}{
		"org_id": orgID, "execution_id": executionID, "from_step": fromStep,
	}
	return p.pub.Publish(ctx, jobs.Job{Type: jobs.JobWorkflowResume, UserID: userID, Payload: payload}, asynq.ProcessAt(delayUntil))
}

// ProcessEvaluate handles workflow.evaluate jobs.
func (p *Processor) ProcessEvaluate(ctx context.Context, job jobs.Job) error {
	orgID := str(job.Payload, "org_id")
	moduleID := str(job.Payload, "module_id")
	recordID := str(job.Payload, "record_id")
	trigger := str(job.Payload, "trigger")
	source := str(job.Payload, "source")
	if source == "" {
		source = entity.SourceUser
	}
	depth := intVal(job.Payload, "depth")
	if depth > engine.MaxDepth {
		p.logger.Warn("workflow: max depth exceeded", zap.Int("depth", depth))
		return nil
	}
	if orgID == "" || moduleID == "" || trigger == "" {
		return fmt.Errorf("workflow.evaluate missing fields: %w", asynq.SkipRetry)
	}

	before := asMap(job.Payload["before"])
	after := asMap(job.Payload["after"])
	changed := stringSlice(job.Payload["changed_fields"])
	if len(changed) == 0 && trigger == entity.TriggerRecordUpdated {
		changed = engine.ChangedFields(before, after)
	}

	// Load live record when after missing (manual / scheduled).
	if len(after) == 0 && recordID != "" && p.deps.Records != nil {
		if rec, err := p.deps.Records.Get(ctx, orgID, moduleID, recordID, job.UserID, false); err == nil && rec != nil {
			after = recordToMap(rec)
		}
	}

	excludeID := str(job.Payload, "exclude_workflow_id")
	specificWF := str(job.Payload, "workflow_id")

	var candidates []entity.MatchCandidate
	var err error
	if specificWF != "" {
		m, err := p.repo.GetMatchByWorkflow(ctx, orgID, specificWF)
		if err != nil {
			return err
		}
		if m != nil && m.Workflow.Status == entity.StatusActive {
			candidates = []entity.MatchCandidate{*m}
		}
	} else {
		candidates, err = p.repo.ListActiveMatches(ctx, orgID, moduleID, trigger)
		if err != nil {
			return err
		}
	}

	moduleAPI, _ := p.repo.ModuleAPIName(ctx, orgID, moduleID)
	orgName, _ := p.repo.OrgName(ctx, orgID)

	for _, c := range candidates {
		if excludeID != "" && c.Workflow.ID == excludeID {
			continue
		}
		if c.Workflow.ModuleID != nil && *c.Workflow.ModuleID != moduleID {
			continue
		}
		matched := false
		for _, t := range c.Triggers {
			if engine.TriggerMatches(t, trigger, changed) {
				matched = true
				break
			}
		}
		if !matched && specificWF == "" {
			continue
		}
		ok, err := engine.EvaluateConditions(c.Conditions, engine.EvalContext{Before: before, After: after, Changed: changed})
		if err != nil {
			p.logger.Warn("workflow: condition error", zap.String("workflow_id", c.Workflow.ID), zap.Error(err))
			continue
		}
		if !ok {
			continue
		}
		if err := p.runWorkflow(ctx, c, orgID, moduleID, recordID, job.UserID, trigger, source, depth, after, moduleAPI, orgName); err != nil {
			p.logger.Error("workflow: run failed", zap.String("workflow_id", c.Workflow.ID), zap.Error(err))
		}
	}
	return nil
}

// ProcessResume continues a delayed workflow from from_step.
func (p *Processor) ProcessResume(ctx context.Context, job jobs.Job) error {
	orgID := str(job.Payload, "org_id")
	executionID := str(job.Payload, "execution_id")
	fromStep := intVal(job.Payload, "from_step")
	if orgID == "" || executionID == "" {
		return fmt.Errorf("workflow.resume missing fields: %w", asynq.SkipRetry)
	}
	e, _, err := p.repo.GetExecution(ctx, orgID, executionID)
	if err != nil || e == nil {
		return err
	}
	if e.VersionID == nil || e.ModuleID == nil || e.RecordID == nil {
		return fmt.Errorf("execution incomplete: %w", asynq.SkipRetry)
	}
	match, err := p.repo.GetMatchByWorkflow(ctx, orgID, e.WorkflowID)
	if err != nil || match == nil {
		return err
	}
	var after map[string]any
	if p.deps.Records != nil {
		if rec, err := p.deps.Records.Get(ctx, orgID, *e.ModuleID, *e.RecordID, job.UserID, false); err == nil && rec != nil {
			after = recordToMap(rec)
		}
	}
	moduleAPI, _ := p.repo.ModuleAPIName(ctx, orgID, *e.ModuleID)
	orgName, _ := p.repo.OrgName(ctx, orgID)
	return p.continueActions(ctx, *match, e, job.UserID, after, moduleAPI, orgName, fromStep)
}

// ProcessScheduledSweep finds date_based / scheduled triggers due.
// Scheduled workflows fan out to module records (capped). Date-based matches
// records whose date field equals today+offset. Both skip records already run
// today for the same workflow+trigger to avoid minute-level re-fire.
func (p *Processor) ProcessScheduledSweep(ctx context.Context) error {
	rows, err := p.repo.DB().Query(ctx, `
		SELECT DISTINCT w.organization_id, w.id, w.module_id, t.type, t.config
		FROM workflows w
		JOIN workflow_versions v ON v.id = w.published_version_id
		JOIN workflow_triggers t ON t.version_id = v.id
		WHERE w.status = 'active'
		  AND t.type IN ('scheduled', 'date_based')`)
	if err != nil {
		return err
	}
	defer rows.Close()

	now := time.Now().UTC()
	for rows.Next() {
		var orgID, wfID, triggerType string
		var moduleID *string
		var config json.RawMessage
		if err := rows.Scan(&orgID, &wfID, &moduleID, &triggerType, &config); err != nil {
			continue
		}
		if moduleID == nil || *moduleID == "" {
			continue
		}
		cfg := map[string]any{}
		_ = json.Unmarshal(config, &cfg)

		switch triggerType {
		case entity.TriggerScheduled:
			if !scheduledDue(cfg, now) {
				continue
			}
			limit := engine.ConfigInt(cfg, "batch_size", 100)
			recs, err := p.repo.ListModuleRecordIDs(ctx, orgID, *moduleID, limit)
			if err != nil {
				p.logger.Warn("workflow: scheduled list records failed", zap.Error(err))
				continue
			}
			for _, rec := range recs {
				ran, _ := p.repo.HasRunToday(ctx, orgID, wfID, rec.ID, entity.TriggerScheduled)
				if ran {
					continue
				}
				after := map[string]any{}
				_ = json.Unmarshal(rec.Data, &after)
				_ = p.EnqueueEvaluate(ctx, orgID, *moduleID, rec.ID, "system", entity.TriggerScheduled, entity.SourceScheduled, 0, wfID, after)
			}

		case entity.TriggerDateBased:
			field := engine.ConfigString(cfg, "field_api_name")
			if field == "" {
				continue
			}
			offsetDays := engine.ConfigInt(cfg, "offset_days", 0)
			target := now.AddDate(0, 0, offsetDays).Format("2006-01-02")
			qrows, err := p.repo.DB().Query(ctx, `
				SELECT id, data FROM records
				WHERE organization_id = $1 AND module_id = $2
				  AND LEFT(COALESCE(data->>$3, ''), 10) = $4
				LIMIT 200`, orgID, *moduleID, field, target)
			if err != nil {
				continue
			}
			for qrows.Next() {
				var rid string
				var data json.RawMessage
				if err := qrows.Scan(&rid, &data); err != nil {
					continue
				}
				ran, _ := p.repo.HasRunToday(ctx, orgID, wfID, rid, entity.TriggerDateBased)
				if ran {
					continue
				}
				after := map[string]any{}
				_ = json.Unmarshal(data, &after)
				_ = p.EnqueueEvaluate(ctx, orgID, *moduleID, rid, "system", entity.TriggerDateBased, entity.SourceScheduled, 0, wfID, after)
			}
			qrows.Close()
		}
	}
	return rows.Err()
}

func scheduledDue(cfg map[string]any, now time.Time) bool {
	// Daily at hour:minute UTC. Optional days_of_week: [0-6] (Sun=0).
	hour := engine.ConfigInt(cfg, "hour", -1)
	minute := engine.ConfigInt(cfg, "minute", 0)
	if hour < 0 || hour > 23 || minute < 0 || minute > 59 {
		return false
	}
	if now.Hour() != hour || now.Minute() != minute {
		return false
	}
	if days, ok := cfg["days_of_week"].([]any); ok && len(days) > 0 {
		wd := int(now.Weekday())
		allowed := false
		for _, d := range days {
			switch v := d.(type) {
			case float64:
				if int(v) == wd {
					allowed = true
				}
			case int:
				if v == wd {
					allowed = true
				}
			}
		}
		if !allowed {
			return false
		}
	}
	return true
}

func (p *Processor) runWorkflow(ctx context.Context, c entity.MatchCandidate, orgID, moduleID, recordID, userID, trigger, source string, depth int, after map[string]any, moduleAPI, orgName string) error {
	started := time.Now().UTC()
	vid := c.Version.ID
	mid := moduleID
	rid := recordID
	exec := &entity.Execution{
		OrganizationID: orgID,
		WorkflowID:     c.Workflow.ID,
		VersionID:      &vid,
		ModuleID:       &mid,
		TriggerType:    trigger,
		Status:         entity.ExecRunning,
		Source:         source,
		Depth:          depth,
		StartedAt:      &started,
	}
	if rid != "" {
		exec.RecordID = &rid
	}
	if err := p.repo.CreateExecution(ctx, exec); err != nil {
		return err
	}

	ownerName, ownerEmail := "", ""
	if ownerID, ok := after["_system"].(map[string]any); ok {
		if oid, ok := ownerID["owner_id"].(string); ok && oid != "" {
			ownerName, ownerEmail, _ = p.repo.UserDisplay(ctx, oid)
		}
	} else if oid, ok := after["owner_id"].(string); ok {
		ownerName, ownerEmail, _ = p.repo.UserDisplay(ctx, oid)
	}

	merge := engine.BuildMergeMap(engine.VariableContext{
		Record: stripSystem(after), ModuleAPI: moduleAPI,
		OwnerName: ownerName, OwnerEmail: ownerEmail, OrgName: orgName,
	})

	run := engine.RunState{
		OrgID: orgID, ModuleID: moduleID, RecordID: recordID, UserID: userID,
		WorkflowID: c.Workflow.ID, Depth: depth, Record: stripSystem(after), Merge: merge,
	}

	status := entity.ExecSucceeded
	var errSummary *string
	failed := 0
	succeeded := 0

	for i, act := range c.Actions {
		cfg := engine.RenderMap(rawMap(act.Config), merge)
		stepStart := time.Now().UTC()
		step := &entity.ExecutionStep{
			ExecutionID: exec.ID, OrganizationID: orgID, ActionID: &act.ID,
			SortOrder: i, ActionType: act.Type, Status: entity.StepRunning,
			Input: mustJSON(cfg), StartedAt: &stepStart,
		}

		result := engine.RunAction(ctx, withGuardedRecords(p.deps, c.Workflow.ID, depth), act, cfg, run)
		stepEnd := time.Now().UTC()
		step.FinishedAt = &stepEnd
		if result.Error != nil {
			failed++
			msg := result.Error.Error()
			step.Status = entity.StepFailed
			step.Error = &msg
			step.Output = mustJSON(map[string]any{})
			_ = p.repo.CreateExecutionStep(ctx, step)

			continueOnErr := c.Workflow.OnActionError == entity.OnErrorContinue
			if act.ContinueOnError != nil {
				continueOnErr = *act.ContinueOnError
			}
			if !continueOnErr {
				status = entity.ExecFailed
				errSummary = &msg
				break
			}
			status = entity.ExecPartial
			errSummary = &msg
			continue
		}
		succeeded++
		step.Status = entity.StepSucceeded
		step.Output = mustJSON(result.Output)
		_ = p.repo.CreateExecutionStep(ctx, step)

		if result.Delay != nil {
			// Schedule resume for remaining actions.
			_ = p.EnqueueResume(ctx, orgID, exec.ID, userID, i+1, time.Now().UTC().Add(*result.Delay))
			status = entity.ExecRunning
			msg := "delayed"
			errSummary = &msg
			break
		}

		// Refresh record after mutations.
		if (act.Type == entity.ActionUpdateRecord || act.Type == entity.ActionAssignOwner) && p.deps.Records != nil && recordID != "" {
			if rec, err := p.deps.Records.Get(ctx, orgID, moduleID, recordID, userID, false); err == nil && rec != nil {
				run.Record = stripSystem(recordToMap(rec))
				run.Merge = engine.BuildMergeMap(engine.VariableContext{
					Record: run.Record, ModuleAPI: moduleAPI,
					OwnerName: ownerName, OwnerEmail: ownerEmail, OrgName: orgName,
				})
			}
		}
	}

	if failed > 0 && succeeded > 0 && status != entity.ExecFailed {
		status = entity.ExecPartial
	}
	finished := time.Now().UTC()
	dur := int(finished.Sub(started).Milliseconds())
	_ = p.repo.FinishExecution(ctx, orgID, exec.ID, status, errSummary, started, finished, dur)

	if p.deps.Activities != nil && recordID != "" {
		action := "WORKFLOW_EXECUTED"
		desc := fmt.Sprintf("Workflow %q finished with status %s", c.Workflow.Name, status)
		if status == entity.ExecFailed {
			action = "WORKFLOW_FAILED"
		}
		_ = p.deps.Activities.LogRecordActivity(ctx, orgID, moduleID, recordID, userID, action, desc, map[string]any{
			"workflow_id": c.Workflow.ID, "execution_id": exec.ID, "status": status,
		})
	}
	return nil
}

func (p *Processor) continueActions(ctx context.Context, c entity.MatchCandidate, exec *entity.Execution, userID string, after map[string]any, moduleAPI, orgName string, fromStep int) error {
	started := time.Now().UTC()
	if exec.StartedAt != nil {
		started = *exec.StartedAt
	}
	ownerName, ownerEmail := "", ""
	merge := engine.BuildMergeMap(engine.VariableContext{
		Record: stripSystem(after), ModuleAPI: moduleAPI, OwnerName: ownerName, OwnerEmail: ownerEmail, OrgName: orgName,
	})
	run := engine.RunState{
		OrgID: exec.OrganizationID, ModuleID: deref(exec.ModuleID), RecordID: deref(exec.RecordID),
		UserID: userID, WorkflowID: c.Workflow.ID, Depth: exec.Depth, Record: stripSystem(after), Merge: merge,
	}
	status := entity.ExecSucceeded
	var errSummary *string
	for i := fromStep; i < len(c.Actions); i++ {
		act := c.Actions[i]
		cfg := engine.RenderMap(rawMap(act.Config), merge)
		stepStart := time.Now().UTC()
		step := &entity.ExecutionStep{
			ExecutionID: exec.ID, OrganizationID: exec.OrganizationID, ActionID: &act.ID,
			SortOrder: i, ActionType: act.Type, Status: entity.StepRunning,
			Input: mustJSON(cfg), StartedAt: &stepStart,
		}
		result := engine.RunAction(ctx, withGuardedRecords(p.deps, c.Workflow.ID, exec.Depth), act, cfg, run)
		stepEnd := time.Now().UTC()
		step.FinishedAt = &stepEnd
		if result.Error != nil {
			msg := result.Error.Error()
			step.Status = entity.StepFailed
			step.Error = &msg
			_ = p.repo.CreateExecutionStep(ctx, step)
			continueOnErr := c.Workflow.OnActionError == entity.OnErrorContinue
			if act.ContinueOnError != nil {
				continueOnErr = *act.ContinueOnError
			}
			if !continueOnErr {
				status = entity.ExecFailed
				errSummary = &msg
				break
			}
			status = entity.ExecPartial
			errSummary = &msg
			continue
		}
		step.Status = entity.StepSucceeded
		step.Output = mustJSON(result.Output)
		_ = p.repo.CreateExecutionStep(ctx, step)
		if result.Delay != nil {
			_ = p.EnqueueResume(ctx, exec.OrganizationID, exec.ID, userID, i+1, time.Now().UTC().Add(*result.Delay))
			status = entity.ExecRunning
			break
		}
	}
	finished := time.Now().UTC()
	dur := int(finished.Sub(started).Milliseconds())
	return p.repo.FinishExecution(ctx, exec.OrganizationID, exec.ID, status, errSummary, started, finished, dur)
}

func recordToMap(rec *rdto.RecordResponse) map[string]any {
	out := clone(rec.Data)
	sys := map[string]any{"visibility": rec.Visibility}
	if rec.OwnerID != nil {
		sys["owner_id"] = *rec.OwnerID
	}
	if rec.AssignedTo != nil {
		sys["assigned_to"] = *rec.AssignedTo
	}
	out["_system"] = sys
	return out
}

func stripSystem(m map[string]any) map[string]any {
	out := clone(m)
	delete(out, "_system")
	return out
}

func clone(m map[string]any) map[string]any {
	if m == nil {
		return map[string]any{}
	}
	out := make(map[string]any, len(m))
	for k, v := range m {
		out[k] = v
	}
	return out
}

func rawMap(b json.RawMessage) map[string]any {
	m := map[string]any{}
	_ = json.Unmarshal(b, &m)
	return m
}

func mustJSON(v any) json.RawMessage {
	b, err := json.Marshal(v)
	if err != nil {
		return json.RawMessage(`{}`)
	}
	return b
}

func str(payload map[string]interface{}, key string) string {
	if payload == nil {
		return ""
	}
	if v, ok := payload[key].(string); ok {
		return v
	}
	return ""
}

func intVal(payload map[string]interface{}, key string) int {
	if payload == nil {
		return 0
	}
	switch v := payload[key].(type) {
	case float64:
		return int(v)
	case int:
		return v
	}
	return 0
}

func asMap(v any) map[string]any {
	if v == nil {
		return map[string]any{}
	}
	if m, ok := v.(map[string]any); ok {
		return m
	}
	// JSON numbers from asynq may decode nested maps as map[string]interface{}
	b, err := json.Marshal(v)
	if err != nil {
		return map[string]any{}
	}
	var m map[string]any
	_ = json.Unmarshal(b, &m)
	if m == nil {
		return map[string]any{}
	}
	return m
}

func stringSlice(v any) []string {
	if v == nil {
		return nil
	}
	switch t := v.(type) {
	case []string:
		return t
	case []any:
		out := make([]string, 0, len(t))
		for _, item := range t {
			if s, ok := item.(string); ok {
				out = append(out, s)
			}
		}
		return out
	}
	return nil
}

func deref(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}
