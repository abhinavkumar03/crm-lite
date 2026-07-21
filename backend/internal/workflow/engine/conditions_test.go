package engine_test

import (
	"encoding/json"
	"testing"

	"github.com/abhinavkumar03/crm-lite/backend/internal/workflow/engine"
	"github.com/abhinavkumar03/crm-lite/backend/internal/workflow/entity"
)

func TestEvaluateConditionsAND(t *testing.T) {
	rootID := "root"
	and := entity.LogicAnd
	field := "status"
	op := "eq"
	val, _ := json.Marshal("Qualified")
	conds := []entity.Condition{
		{ID: rootID, NodeType: entity.NodeGroup, Logic: &and},
		{ID: "p1", ParentID: &rootID, NodeType: entity.NodePredicate, FieldAPIName: &field, Operator: &op, Value: val},
	}
	ok, err := engine.EvaluateConditions(conds, engine.EvalContext{
		After: map[string]any{"status": "Qualified"},
	})
	if err != nil || !ok {
		t.Fatalf("expected pass, got ok=%v err=%v", ok, err)
	}
	ok, err = engine.EvaluateConditions(conds, engine.EvalContext{
		After: map[string]any{"status": "New"},
	})
	if err != nil || ok {
		t.Fatalf("expected fail, got ok=%v err=%v", ok, err)
	}
}

func TestChangedFields(t *testing.T) {
	before := map[string]any{"status": "New", "_system": map[string]any{"owner_id": "a"}}
	after := map[string]any{"status": "Qualified", "_system": map[string]any{"owner_id": "b"}}
	changed := engine.ChangedFields(before, after)
	foundStatus, foundOwner := false, false
	for _, c := range changed {
		if c == "status" {
			foundStatus = true
		}
		if c == "owner_id" {
			foundOwner = true
		}
	}
	if !foundStatus || !foundOwner {
		t.Fatalf("changed=%v", changed)
	}
}
