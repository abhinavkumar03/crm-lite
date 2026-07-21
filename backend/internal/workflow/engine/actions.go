package engine

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	ndto "github.com/abhinavkumar03/crm-lite/backend/internal/notification/dto"
	rdto "github.com/abhinavkumar03/crm-lite/backend/internal/record/dto"
	wdto "github.com/abhinavkumar03/crm-lite/backend/internal/workspace/dto"
	"github.com/abhinavkumar03/crm-lite/backend/internal/workflow/entity"
)

const MaxDepth = 3

// RecordMutator creates/updates records (system source).
type RecordMutator interface {
	Get(ctx context.Context, orgID, moduleID, id, userID string, expand bool) (*rdto.RecordResponse, error)
	Create(ctx context.Context, orgID, moduleID, userID string, req rdto.CreateRecordRequest) (*rdto.RecordResponse, error)
	Update(ctx context.Context, orgID, moduleID, id, userID string, req rdto.UpdateRecordRequest) (*rdto.RecordResponse, error)
}

// Notifier composes outbound messages.
type Notifier interface {
	Compose(ctx context.Context, orgID, userID string, req ndto.ComposeRequest) (*ndto.NotificationResponse, error)
}

// NoteWriter attaches notes.
type NoteWriter interface {
	CreateNote(ctx context.Context, orgID, moduleID, recordID, userID string, req wdto.CreateNoteRequest) (*wdto.NoteResponse, error)
}

// ActivityWriter writes timeline entries.
type ActivityWriter interface {
	LogRecordActivity(ctx context.Context, orgID, moduleID, recordID, userID, action, description string, metadata map[string]any) error
}

// WorkflowInvoker enqueues another workflow run.
type WorkflowInvoker interface {
	EnqueueEvaluate(ctx context.Context, orgID, moduleID, recordID, userID, trigger, source string, depth int, workflowID string, after map[string]any) error
	EnqueueResume(ctx context.Context, orgID, executionID, userID string, fromStep int, delayUntil time.Time) error
}

// ActionDeps bundles executors.
type ActionDeps struct {
	Records    RecordMutator
	Notify     Notifier
	Notes      NoteWriter
	Activities ActivityWriter
	Invoker    WorkflowInvoker
	HTTP       *http.Client
}

// ActionResult is the outcome of one action.
type ActionResult struct {
	Output map[string]any
	Error  error
	Delay  *time.Duration // if set, pause remaining steps
}

// RunAction executes a single action type.
func RunAction(ctx context.Context, deps ActionDeps, act entity.Action, cfg map[string]any, run RunState) ActionResult {
	switch act.Type {
	case entity.ActionUpdateRecord:
		return actionUpdateRecord(ctx, deps, cfg, run)
	case entity.ActionCreateRecord:
		return actionCreateRecord(ctx, deps, cfg, run)
	case entity.ActionAssignOwner:
		return actionAssignOwner(ctx, deps, cfg, run)
	case entity.ActionSendEmail:
		return actionSend(ctx, deps, cfg, run, "email")
	case entity.ActionSendWhatsApp:
		return actionSend(ctx, deps, cfg, run, "whatsapp")
	case entity.ActionCreateNote:
		return actionCreateNote(ctx, deps, cfg, run)
	case entity.ActionCreateActivity:
		return actionCreateActivity(ctx, deps, cfg, run)
	case entity.ActionInvokeWorkflow:
		return actionInvokeWorkflow(ctx, deps, cfg, run)
	case entity.ActionWebhook:
		return actionWebhook(ctx, deps, cfg, run)
	case entity.ActionDelay:
		return actionDelay(cfg)
	case entity.ActionDeleteRecord:
		return ActionResult{Error: fmt.Errorf("delete_record not implemented")}
	case entity.ActionBranch:
		return ActionResult{Error: fmt.Errorf("branch not implemented")}
	default:
		return ActionResult{Error: fmt.Errorf("unknown action type %s", act.Type)}
	}
}

// RunState is mutable context during a workflow run.
type RunState struct {
	OrgID      string
	ModuleID   string
	RecordID   string
	UserID     string
	WorkflowID string
	Depth      int
	Record     map[string]any
	Merge      map[string]any
}

func actionUpdateRecord(ctx context.Context, deps ActionDeps, cfg map[string]any, run RunState) ActionResult {
	if deps.Records == nil {
		return ActionResult{Error: fmt.Errorf("records unavailable")}
	}
	fields, _ := cfg["fields"].(map[string]any)
	if fields == nil {
		fields = map[string]any{}
	}
	data := cloneMap(run.Record)
	for k, v := range fields {
		data[k] = v
	}
	rec, err := deps.Records.Update(ctx, run.OrgID, run.ModuleID, run.RecordID, run.UserID, rdto.UpdateRecordRequest{Data: data})
	if err != nil {
		return ActionResult{Error: err}
	}
	return ActionResult{Output: map[string]any{"record_id": rec.ID}}
}

func actionCreateRecord(ctx context.Context, deps ActionDeps, cfg map[string]any, run RunState) ActionResult {
	if deps.Records == nil {
		return ActionResult{Error: fmt.Errorf("records unavailable")}
	}
	moduleID := ConfigString(cfg, "module_id")
	if moduleID == "" {
		return ActionResult{Error: fmt.Errorf("create_record requires module_id")}
	}
	fields, _ := cfg["fields"].(map[string]any)
	if fields == nil {
		fields = map[string]any{}
	}
	rec, err := deps.Records.Create(ctx, run.OrgID, moduleID, run.UserID, rdto.CreateRecordRequest{Data: fields})
	if err != nil {
		return ActionResult{Error: err}
	}
	return ActionResult{Output: map[string]any{"record_id": rec.ID, "module_id": moduleID}}
}

func actionAssignOwner(ctx context.Context, deps ActionDeps, cfg map[string]any, run RunState) ActionResult {
	if deps.Records == nil {
		return ActionResult{Error: fmt.Errorf("records unavailable")}
	}
	ownerID := ConfigString(cfg, "owner_id")
	assignedTo := ConfigString(cfg, "assigned_to")
	if ownerID == "" && assignedTo == "" {
		return ActionResult{Error: fmt.Errorf("assign_owner requires owner_id or assigned_to")}
	}
	req := rdto.UpdateRecordRequest{Data: cloneMap(run.Record)}
	if ownerID != "" {
		req.OwnerID = &ownerID
	}
	if assignedTo != "" {
		req.AssignedTo = &assignedTo
	}
	rec, err := deps.Records.Update(ctx, run.OrgID, run.ModuleID, run.RecordID, run.UserID, req)
	if err != nil {
		return ActionResult{Error: err}
	}
	return ActionResult{Output: map[string]any{"record_id": rec.ID, "owner_id": ownerID, "assigned_to": assignedTo}}
}

func actionSend(ctx context.Context, deps ActionDeps, cfg map[string]any, run RunState, channel string) ActionResult {
	if deps.Notify == nil {
		return ActionResult{Error: fmt.Errorf("notifier unavailable")}
	}
	to := ResolveRecipient(cfg, run.Record, channel)
	if to == "" {
		return ActionResult{Error: fmt.Errorf("no recipient for %s", channel)}
	}
	body := ConfigString(cfg, "body")
	if body == "" {
		body = " "
	}
	subject := ConfigString(cfg, "subject")
	req := ndto.ComposeRequest{
		Mode: "send", Channel: channel, To: to, Subject: subject, Body: body,
		Template: ConfigString(cfg, "template"), TemplateID: ConfigString(cfg, "template_id"),
		Data: run.Merge, ModuleID: run.ModuleID, EntityID: run.RecordID, EntityType: "RECORD",
	}
	resp, err := deps.Notify.Compose(ctx, run.OrgID, run.UserID, req)
	if err != nil {
		return ActionResult{Error: err}
	}
	return ActionResult{Output: map[string]any{"notification_id": resp.ID, "status": resp.Status}}
}

func actionCreateNote(ctx context.Context, deps ActionDeps, cfg map[string]any, run RunState) ActionResult {
	if deps.Notes == nil {
		return ActionResult{Error: fmt.Errorf("notes unavailable")}
	}
	body := ConfigString(cfg, "body")
	if body == "" {
		body = "Workflow note"
	}
	var title *string
	if t := ConfigString(cfg, "title"); t != "" {
		title = &t
	}
	n, err := deps.Notes.CreateNote(ctx, run.OrgID, run.ModuleID, run.RecordID, run.UserID, wdto.CreateNoteRequest{Title: title, Body: body})
	if err != nil {
		return ActionResult{Error: err}
	}
	return ActionResult{Output: map[string]any{"note_id": n.ID}}
}

func actionCreateActivity(ctx context.Context, deps ActionDeps, cfg map[string]any, run RunState) ActionResult {
	if deps.Activities == nil {
		return ActionResult{Error: fmt.Errorf("activities unavailable")}
	}
	action := ConfigString(cfg, "action")
	if action == "" {
		action = "WORKFLOW_ACTIVITY"
	}
	desc := ConfigString(cfg, "description")
	if desc == "" {
		desc = "Workflow activity"
	}
	meta, _ := cfg["metadata"].(map[string]any)
	if err := deps.Activities.LogRecordActivity(ctx, run.OrgID, run.ModuleID, run.RecordID, run.UserID, action, desc, meta); err != nil {
		return ActionResult{Error: err}
	}
	return ActionResult{Output: map[string]any{"action": action}}
}

func actionInvokeWorkflow(ctx context.Context, deps ActionDeps, cfg map[string]any, run RunState) ActionResult {
	if deps.Invoker == nil {
		return ActionResult{Error: fmt.Errorf("invoker unavailable")}
	}
	if run.Depth >= MaxDepth {
		return ActionResult{Error: fmt.Errorf("max workflow depth exceeded")}
	}
	wfID := ConfigString(cfg, "workflow_id")
	if wfID == "" {
		return ActionResult{Error: fmt.Errorf("invoke_workflow requires workflow_id")}
	}
	err := deps.Invoker.EnqueueEvaluate(ctx, run.OrgID, run.ModuleID, run.RecordID, run.UserID,
		entity.TriggerManual, entity.SourceWorkflow, run.Depth+1, wfID, run.Record)
	if err != nil {
		return ActionResult{Error: err}
	}
	return ActionResult{Output: map[string]any{"invoked_workflow_id": wfID}}
}

func actionWebhook(ctx context.Context, deps ActionDeps, cfg map[string]any, run RunState) ActionResult {
	url := ConfigString(cfg, "url")
	if url == "" {
		return ActionResult{Error: fmt.Errorf("webhook requires url")}
	}
	method := ConfigString(cfg, "method")
	if method == "" {
		method = http.MethodPost
	}
	payload := cfg["payload"]
	if payload == nil {
		payload = map[string]any{"record": run.Record, "module_id": run.ModuleID, "record_id": run.RecordID}
	}
	body, _ := json.Marshal(payload)
	client := deps.HTTP
	if client == nil {
		client = &http.Client{Timeout: 15 * time.Second}
	}
	req, err := http.NewRequestWithContext(ctx, method, url, bytes.NewReader(body))
	if err != nil {
		return ActionResult{Error: err}
	}
	req.Header.Set("Content-Type", "application/json")
	if headers, ok := cfg["headers"].(map[string]any); ok {
		for k, v := range headers {
			req.Header.Set(k, fmt.Sprint(v))
		}
	}
	resp, err := client.Do(req)
	if err != nil {
		return ActionResult{Error: err}
	}
	defer resp.Body.Close()
	respBody, _ := io.ReadAll(io.LimitReader(resp.Body, 4096))
	if resp.StatusCode >= 400 {
		return ActionResult{Error: fmt.Errorf("webhook status %d: %s", resp.StatusCode, string(respBody))}
	}
	return ActionResult{Output: map[string]any{"status": resp.StatusCode, "body": string(respBody)}}
}

func actionDelay(cfg map[string]any) ActionResult {
	seconds := ConfigInt(cfg, "seconds", 0)
	if seconds == 0 {
		days := ConfigInt(cfg, "days", 0)
		hours := ConfigInt(cfg, "hours", 0)
		minutes := ConfigInt(cfg, "minutes", 0)
		seconds = days*86400 + hours*3600 + minutes*60
	}
	if seconds <= 0 {
		return ActionResult{Error: fmt.Errorf("delay requires positive duration")}
	}
	d := time.Duration(seconds) * time.Second
	return ActionResult{Output: map[string]any{"delay_seconds": seconds}, Delay: &d}
}

func cloneMap(in map[string]any) map[string]any {
	out := make(map[string]any, len(in))
	for k, v := range in {
		if k == "_system" {
			continue
		}
		out[k] = v
	}
	return out
}
