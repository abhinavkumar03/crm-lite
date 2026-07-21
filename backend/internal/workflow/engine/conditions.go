package engine

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/abhinavkumar03/crm-lite/backend/internal/workflow/entity"
)

// EvalContext holds record snapshots for condition evaluation.
type EvalContext struct {
	Before  map[string]any
	After   map[string]any
	Changed []string
}

// EvaluateConditions returns true when the tree passes (empty tree = true).
func EvaluateConditions(conds []entity.Condition, ctx EvalContext) (bool, error) {
	if len(conds) == 0 {
		return true, nil
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
		return true, nil
	}
	return evalNode(roots[0], byParent, ctx)
}

func evalNode(n entity.Condition, byParent map[string][]entity.Condition, ctx EvalContext) (bool, error) {
	if n.NodeType == entity.NodeGroup {
		children := byParent[n.ID]
		if len(children) == 0 {
			return true, nil
		}
		logic := entity.LogicAnd
		if n.Logic != nil {
			logic = *n.Logic
		}
		if logic == entity.LogicOr {
			for _, ch := range children {
				ok, err := evalNode(ch, byParent, ctx)
				if err != nil {
					return false, err
				}
				if ok {
					return true, nil
				}
			}
			return false, nil
		}
		for _, ch := range children {
			ok, err := evalNode(ch, byParent, ctx)
			if err != nil {
				return false, err
			}
			if !ok {
				return false, nil
			}
		}
		return true, nil
	}
	return evalPredicate(n, ctx)
}

func evalPredicate(n entity.Condition, ctx EvalContext) (bool, error) {
	if n.FieldAPIName == nil || n.Operator == nil {
		return false, fmt.Errorf("predicate missing field/operator")
	}
	field := *n.FieldAPIName
	op := *n.Operator
	actual := lookupField(ctx.After, field)
	var expected any
	if len(n.Value) > 0 {
		_ = json.Unmarshal(n.Value, &expected)
	}
	switch op {
	case "eq", "equals":
		return compareEqual(actual, expected), nil
	case "neq", "ne", "not_equals":
		return !compareEqual(actual, expected), nil
	case "gt":
		return compareOrdered(actual, expected) > 0, nil
	case "lt":
		return compareOrdered(actual, expected) < 0, nil
	case "gte":
		return compareOrdered(actual, expected) >= 0, nil
	case "lte":
		return compareOrdered(actual, expected) <= 0, nil
	case "between":
		arr, ok := expected.([]any)
		if !ok || len(arr) < 2 {
			return false, nil
		}
		return compareOrdered(actual, arr[0]) >= 0 && compareOrdered(actual, arr[1]) <= 0, nil
	case "starts_with":
		return strings.HasPrefix(stringify(actual), stringify(expected)), nil
	case "ends_with":
		return strings.HasSuffix(stringify(actual), stringify(expected)), nil
	case "contains":
		return strings.Contains(strings.ToLower(stringify(actual)), strings.ToLower(stringify(expected))), nil
	case "is_empty":
		return isEmpty(actual), nil
	case "is_not_empty":
		return !isEmpty(actual), nil
	case "changed":
		for _, c := range ctx.Changed {
			if c == field {
				return true, nil
			}
		}
		return false, nil
	default:
		return false, fmt.Errorf("unknown operator %s", op)
	}
}

func lookupField(data map[string]any, field string) any {
	if data == nil {
		return nil
	}
	// System columns may be top-level on the after map under _system.
	if sys, ok := data["_system"].(map[string]any); ok {
		if v, ok := sys[field]; ok {
			return v
		}
	}
	if v, ok := data[field]; ok {
		return v
	}
	return nil
}

func isEmpty(v any) bool {
	if v == nil {
		return true
	}
	switch t := v.(type) {
	case string:
		return strings.TrimSpace(t) == ""
	case []any:
		return len(t) == 0
	case map[string]any:
		return len(t) == 0
	default:
		return false
	}
}

func stringify(v any) string {
	if v == nil {
		return ""
	}
	switch t := v.(type) {
	case string:
		return t
	case float64:
		return strconv.FormatFloat(t, 'f', -1, 64)
	case bool:
		return strconv.FormatBool(t)
	default:
		b, _ := json.Marshal(t)
		return string(b)
	}
}

func compareEqual(a, b any) bool {
	return stringify(a) == stringify(b) || numbersEqual(a, b)
}

func numbersEqual(a, b any) bool {
	af, aok := toFloat(a)
	bf, bok := toFloat(b)
	return aok && bok && af == bf
}

func compareOrdered(a, b any) int {
	af, aok := toFloat(a)
	bf, bok := toFloat(b)
	if aok && bok {
		if af < bf {
			return -1
		}
		if af > bf {
			return 1
		}
		return 0
	}
	as, bs := stringify(a), stringify(b)
	// Try dates.
	if at, err1 := time.Parse(time.RFC3339, as); err1 == nil {
		if bt, err2 := time.Parse(time.RFC3339, bs); err2 == nil {
			return at.Compare(bt)
		}
	}
	if at, err1 := time.Parse("2006-01-02", as); err1 == nil {
		if bt, err2 := time.Parse("2006-01-02", bs); err2 == nil {
			return at.Compare(bt)
		}
	}
	return strings.Compare(as, bs)
}

func toFloat(v any) (float64, bool) {
	switch t := v.(type) {
	case float64:
		return t, true
	case float32:
		return float64(t), true
	case int:
		return float64(t), true
	case int64:
		return float64(t), true
	case json.Number:
		f, err := t.Float64()
		return f, err == nil
	case string:
		f, err := strconv.ParseFloat(t, 64)
		return f, err == nil
	default:
		return 0, false
	}
}

// TriggerMatches checks whether a candidate trigger applies to the event.
func TriggerMatches(t entity.Trigger, triggerType string, changed []string) bool {
	if t.Type != triggerType {
		return false
	}
	if t.Type != entity.TriggerFieldUpdated {
		return true
	}
	cfg := map[string]any{}
	_ = json.Unmarshal(t.Config, &cfg)
	field, _ := cfg["field_api_name"].(string)
	if field == "" {
		return len(changed) > 0
	}
	for _, c := range changed {
		if c == field {
			return true
		}
	}
	return false
}

// ChangedFields computes keys that differ between before and after data maps.
func ChangedFields(before, after map[string]any) []string {
	keys := map[string]struct{}{}
	for k := range before {
		keys[k] = struct{}{}
	}
	for k := range after {
		keys[k] = struct{}{}
	}
	var out []string
	for k := range keys {
		if k == "_system" {
			continue
		}
		if !compareEqual(before[k], after[k]) {
			out = append(out, k)
		}
	}
	// System fields.
	bs, _ := before["_system"].(map[string]any)
	as, _ := after["_system"].(map[string]any)
	sysKeys := []string{"owner_id", "assigned_to", "team_id", "department_id", "visibility"}
	for _, k := range sysKeys {
		var bv, av any
		if bs != nil {
			bv = bs[k]
		}
		if as != nil {
			av = as[k]
		}
		if !compareEqual(bv, av) {
			out = append(out, k)
		}
	}
	return out
}
