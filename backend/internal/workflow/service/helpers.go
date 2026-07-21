package service

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/abhinavkumar03/crm-lite/backend/internal/workflow/dto"
	"github.com/abhinavkumar03/crm-lite/backend/internal/workflow/entity"
)

func validateDefinition(triggers []dto.TriggerInput, conditions *dto.ConditionInput, actions []dto.ActionInput) error {
	validTriggers := map[string]bool{
		entity.TriggerRecordCreated: true, entity.TriggerRecordUpdated: true,
		entity.TriggerRecordDeleted: true, entity.TriggerFieldUpdated: true,
		entity.TriggerScheduled: true, entity.TriggerDateBased: true, entity.TriggerManual: true,
	}
	validActions := map[string]bool{
		entity.ActionUpdateRecord: true, entity.ActionCreateRecord: true, entity.ActionDeleteRecord: true,
		entity.ActionAssignOwner: true, entity.ActionSendEmail: true, entity.ActionSendWhatsApp: true,
		entity.ActionCreateNote: true, entity.ActionCreateActivity: true, entity.ActionWebhook: true,
		entity.ActionDelay: true, entity.ActionInvokeWorkflow: true, entity.ActionBranch: true,
	}
	for _, t := range triggers {
		if !validTriggers[t.Type] {
			return fmt.Errorf("%w: unknown trigger %s", ErrInvalidInput, t.Type)
		}
	}
	if conditions != nil {
		if err := validateCondition(conditions); err != nil {
			return err
		}
	}
	for _, a := range actions {
		if !validActions[a.Type] {
			return fmt.Errorf("%w: unknown action %s", ErrInvalidInput, a.Type)
		}
	}
	return nil
}

func validateCondition(c *dto.ConditionInput) error {
	if c.NodeType == entity.NodeGroup {
		if c.Logic == nil || (*c.Logic != entity.LogicAnd && *c.Logic != entity.LogicOr) {
			return fmt.Errorf("%w: group requires and/or logic", ErrInvalidInput)
		}
		for i := range c.Children {
			if err := validateCondition(&c.Children[i]); err != nil {
				return err
			}
		}
		return nil
	}
	if c.NodeType == entity.NodePredicate {
		if c.FieldAPIName == nil || strings.TrimSpace(*c.FieldAPIName) == "" {
			return fmt.Errorf("%w: predicate requires field_api_name", ErrInvalidInput)
		}
		if c.Operator == nil || strings.TrimSpace(*c.Operator) == "" {
			return fmt.Errorf("%w: predicate requires operator", ErrInvalidInput)
		}
		return nil
	}
	return fmt.Errorf("%w: invalid node_type", ErrInvalidInput)
}

func buildConditionTree(conds []entity.Condition) *dto.ConditionResponse {
	if len(conds) == 0 {
		return nil
	}
	byParent := map[string][]entity.Condition{}
	var roots []entity.Condition
	for _, c := range conds {
		if c.ParentID == nil {
			roots = append(roots, c)
		} else {
			byParent[*c.ParentID] = append(byParent[*c.ParentID], c)
		}
	}
	if len(roots) == 0 {
		return nil
	}
	var build func(c entity.Condition) dto.ConditionResponse
	build = func(c entity.Condition) dto.ConditionResponse {
		node := dto.ConditionResponse{
			ID: c.ID, NodeType: c.NodeType, Logic: c.Logic,
			FieldAPIName: c.FieldAPIName, Operator: c.Operator,
		}
		if len(c.Value) > 0 {
			var v any
			_ = json.Unmarshal(c.Value, &v)
			node.Value = v
		}
		for _, ch := range byParent[c.ID] {
			node.Children = append(node.Children, build(ch))
		}
		return node
	}
	root := build(roots[0])
	return &root
}

func toExecSummary(e entity.Execution, name string) dto.ExecutionSummary {
	return dto.ExecutionSummary{
		ID: e.ID, WorkflowID: e.WorkflowID, WorkflowName: name, VersionID: e.VersionID,
		ModuleID: e.ModuleID, RecordID: e.RecordID, TriggerType: e.TriggerType,
		Status: e.Status, Source: e.Source, Depth: e.Depth, ErrorSummary: e.ErrorSummary,
		StartedAt: e.StartedAt, FinishedAt: e.FinishedAt, DurationMs: e.DurationMs, CreatedAt: e.CreatedAt,
	}
}

func rawToMap(raw json.RawMessage) map[string]any {
	if len(raw) == 0 {
		return map[string]any{}
	}
	var m map[string]any
	if err := json.Unmarshal(raw, &m); err != nil {
		return map[string]any{}
	}
	if m == nil {
		return map[string]any{}
	}
	return m
}

func str(v any) string {
	if v == nil {
		return ""
	}
	s, _ := v.(string)
	return s
}

func asMap(v any) map[string]any {
	if v == nil {
		return map[string]any{}
	}
	if m, ok := v.(map[string]any); ok {
		return m
	}
	return map[string]any{}
}

func mapToCondition(m map[string]any) dto.ConditionInput {
	c := dto.ConditionInput{NodeType: str(m["node_type"])}
	if l := str(m["logic"]); l != "" {
		c.Logic = &l
	}
	if f := str(m["field_api_name"]); f != "" {
		c.FieldAPIName = &f
	}
	if o := str(m["operator"]); o != "" {
		c.Operator = &o
	}
	if v, ok := m["value"]; ok {
		c.Value = v
	}
	if children, ok := m["children"].([]any); ok {
		for _, raw := range children {
			if cm, ok := raw.(map[string]any); ok {
				c.Children = append(c.Children, mapToCondition(cm))
			}
		}
	}
	return c
}

func catalogOperators() []dto.BuilderOperator {
	return []dto.BuilderOperator{
		{Key: "eq", Label: "Equals", ValueArity: "one"},
		{Key: "neq", Label: "Not Equals", ValueArity: "one"},
		{Key: "gt", Label: "Greater Than", ValueArity: "one"},
		{Key: "lt", Label: "Less Than", ValueArity: "one"},
		{Key: "gte", Label: "Greater or Equal", ValueArity: "one"},
		{Key: "lte", Label: "Less or Equal", ValueArity: "one"},
		{Key: "between", Label: "Between", ValueArity: "two"},
		{Key: "starts_with", Label: "Starts With", ValueArity: "one"},
		{Key: "ends_with", Label: "Ends With", ValueArity: "one"},
		{Key: "contains", Label: "Contains", ValueArity: "one"},
		{Key: "is_empty", Label: "Is Empty", ValueArity: "none"},
		{Key: "is_not_empty", Label: "Is Not Empty", ValueArity: "none"},
	}
}

func catalogActions() []dto.BuilderAction {
	return []dto.BuilderAction{
		{Type: entity.ActionUpdateRecord, Label: "Update Record", Description: "Set fields on the triggering record", MVP: true},
		{Type: entity.ActionCreateRecord, Label: "Create Record", Description: "Create a record in a target module", MVP: true},
		{Type: entity.ActionAssignOwner, Label: "Assign Owner", Description: "Assign owner or assignee", MVP: true},
		{Type: entity.ActionSendEmail, Label: "Send Email", Description: "Send email via Notification Center", MVP: true},
		{Type: entity.ActionSendWhatsApp, Label: "Send WhatsApp", Description: "Send WhatsApp via Notification Center", MVP: true},
		{Type: entity.ActionCreateNote, Label: "Create Note", Description: "Attach a note to the record", MVP: true},
		{Type: entity.ActionCreateActivity, Label: "Create Activity", Description: "Write a timeline activity", MVP: true},
		{Type: entity.ActionInvokeWorkflow, Label: "Invoke Workflow", Description: "Run another published workflow", MVP: true},
		{Type: entity.ActionWebhook, Label: "Webhook", Description: "Call an external HTTP endpoint", MVP: true},
		{Type: entity.ActionDelay, Label: "Delay", Description: "Wait before continuing", MVP: true},
		{Type: entity.ActionDeleteRecord, Label: "Delete Record", Description: "Delete a record (future)", MVP: false},
		{Type: entity.ActionBranch, Label: "Branch", Description: "Conditional branch (future)", MVP: false},
	}
}

func catalogTriggers() []dto.BuilderTrigger {
	return []dto.BuilderTrigger{
		{Type: entity.TriggerRecordCreated, Label: "Record Created", Description: "When a record is created", MVP: true},
		{Type: entity.TriggerRecordUpdated, Label: "Record Updated", Description: "When a record is updated", MVP: true},
		{Type: entity.TriggerFieldUpdated, Label: "Field Updated", Description: "When a specific field changes", MVP: true},
		{Type: entity.TriggerRecordDeleted, Label: "Record Deleted", Description: "When a record is deleted", MVP: true},
		{Type: entity.TriggerManual, Label: "Manual", Description: "Run on demand", MVP: true},
		{Type: entity.TriggerScheduled, Label: "Scheduled", Description: "Cron / recurring schedule", MVP: true},
		{Type: entity.TriggerDateBased, Label: "Date Based", Description: "Based on a date field", MVP: true},
	}
}

func catalogVariables() []dto.BuilderVariable {
	return []dto.BuilderVariable{
		{Key: "record.*", Label: "Record fields", Description: "Any field on the triggering record"},
		{Key: "{{module}}.*", Label: "Module alias", Description: "Same as record.* using module api_name (e.g. lead.name)"},
		{Key: "owner.name", Label: "Owner name", Description: "Record owner display name"},
		{Key: "owner.email", Label: "Owner email", Description: "Record owner email"},
		{Key: "workspace.name", Label: "Workspace name", Description: "Organization / workspace name"},
		{Key: "today", Label: "Today", Description: "Current date (YYYY-MM-DD)"},
		{Key: "current_date", Label: "Current date", Description: "Alias of today"},
	}
}
