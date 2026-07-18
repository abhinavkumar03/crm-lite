package repository

import (
	"fmt"
	"regexp"
	"strings"

	fieldentity "github.com/abhinavkumar03/crm-lite/backend/internal/field/entity"
	"github.com/abhinavkumar03/crm-lite/backend/internal/record/dto"
)

// apiNameRe guards against SQL injection when a field api_name is interpolated
// into a JSON path. Field api_names are lowercase identifiers by construction;
// anything else is rejected before it reaches the query.
var apiNameRe = regexp.MustCompile(`^[a-z][a-z0-9_]*$`)

// FieldMeta is the minimal per-field metadata the query builder needs.
type FieldMeta struct {
	Type       string
	Searchable bool
	Filterable bool
}

// BuildMeta reduces full field entities to the query builder's view of them.
func BuildMeta(fields []fieldentity.Field) map[string]FieldMeta {
	meta := make(map[string]FieldMeta, len(fields))
	for _, f := range fields {
		meta[f.APIName] = FieldMeta{
			Type:       f.FieldType,
			Searchable: f.IsSearchable,
			Filterable: f.IsFilterable,
		}
	}
	return meta
}

func numericType(t string) bool {
	return t == fieldentity.TypeNumber || t == fieldentity.TypeCurrency
}

func dateType(t string) bool {
	return t == fieldentity.TypeDate || t == fieldentity.TypeDatetime
}

// valueExpr returns the SQL expression that extracts a field's value, cast to a
// comparable type where appropriate.
func valueExpr(apiName, fieldType string) string {
	base := fmt.Sprintf("data->>'%s'", apiName)
	switch {
	case numericType(fieldType):
		return "(" + base + ")::numeric"
	case dateType(fieldType):
		return "(" + base + ")::timestamptz"
	default:
		return base
	}
}

func castParam(fieldType, placeholder string) string {
	switch {
	case numericType(fieldType):
		return placeholder + "::numeric"
	case dateType(fieldType):
		return placeholder + "::timestamptz"
	default:
		return placeholder
	}
}

func toText(v any) string {
	if v == nil {
		return ""
	}
	return fmt.Sprintf("%v", v)
}

// WhereClause is the composable output of the filter/search builder.
type WhereClause struct {
	SQL  string
	Args []any
}

// BuildWhere composes the full WHERE clause (org/module scope + search +
// filters) with positional args. Only fields present in meta are honoured, and
// filters are only applied to filterable fields, so untrusted input can never
// widen the query beyond the module's own schema.
func BuildWhere(orgID, moduleID string, q dto.ListQuery, meta map[string]FieldMeta) WhereClause {
	conds := []string{"organization_id = $1", "module_id = $2"}
	args := []any{orgID, moduleID}
	next := 3

	if search := strings.TrimSpace(q.Search); search != "" {
		var ors []string
		for apiName, m := range meta {
			if !m.Searchable || !apiNameRe.MatchString(apiName) {
				continue
			}
			ors = append(ors, fmt.Sprintf("data->>'%s' ILIKE $%d", apiName, next))
		}
		if len(ors) > 0 {
			// Deterministic ordering keeps generated SQL stable/testable.
			ors = sortedStrings(ors)
			args = append(args, "%"+search+"%")
			conds = append(conds, "("+strings.Join(ors, " OR ")+")")
			next++
		}
	}

	for _, f := range q.Filters {
		m, ok := meta[f.Field]
		if !ok || !m.Filterable || !apiNameRe.MatchString(f.Field) {
			continue
		}
		cond, used := buildFilter(f, m.Type, next)
		if cond == "" {
			continue
		}
		conds = append(conds, cond)
		args = append(args, filterArgs(f, m.Type)...)
		next += used
	}

	return WhereClause{SQL: strings.Join(conds, " AND "), Args: args}
}

func buildFilter(f dto.FilterClause, fieldType string, arg int) (string, int) {
	col := fmt.Sprintf("data->>'%s'", f.Field)
	ph := fmt.Sprintf("$%d", arg)

	switch f.Operator {
	case dto.OpEquals:
		return fmt.Sprintf("%s = %s", col, ph), 1
	case dto.OpNotEquals:
		return fmt.Sprintf("%s IS DISTINCT FROM %s", col, ph), 1
	case dto.OpContains:
		return fmt.Sprintf("%s ILIKE %s", col, ph), 1
	case dto.OpGreaterThan:
		return fmt.Sprintf("%s > %s", valueExpr(f.Field, fieldType), castParam(fieldType, ph)), 1
	case dto.OpLessThan:
		return fmt.Sprintf("%s < %s", valueExpr(f.Field, fieldType), castParam(fieldType, ph)), 1
	case dto.OpGreaterEq:
		return fmt.Sprintf("%s >= %s", valueExpr(f.Field, fieldType), castParam(fieldType, ph)), 1
	case dto.OpLessEq:
		return fmt.Sprintf("%s <= %s", valueExpr(f.Field, fieldType), castParam(fieldType, ph)), 1
	case dto.OpIn:
		return fmt.Sprintf("%s = ANY(%s)", col, ph), 1
	default:
		return "", 0
	}
}

func filterArgs(f dto.FilterClause, _ string) []any {
	switch f.Operator {
	case dto.OpContains:
		return []any{"%" + toText(f.Value) + "%"}
	case dto.OpIn:
		items, ok := f.Value.([]any)
		if !ok {
			return []any{[]string{}}
		}
		out := make([]string, 0, len(items))
		for _, it := range items {
			out = append(out, toText(it))
		}
		return []any{out}
	default:
		return []any{toText(f.Value)}
	}
}

// BuildOrderBy returns a safe ORDER BY clause. Sorting is allowed on the record
// timestamps and on any field in meta; everything else falls back to created_at.
func BuildOrderBy(sort, order string, meta map[string]FieldMeta) string {
	dir := "DESC"
	if strings.EqualFold(order, "asc") {
		dir = "ASC"
	}

	switch sort {
	case "", "created_at":
		return "ORDER BY created_at " + dir
	case "updated_at":
		return "ORDER BY updated_at " + dir
	}

	if m, ok := meta[sort]; ok && apiNameRe.MatchString(sort) {
		return fmt.Sprintf("ORDER BY %s %s NULLS LAST", valueExpr(sort, m.Type), dir)
	}
	return "ORDER BY created_at " + dir
}

func sortedStrings(in []string) []string {
	out := make([]string, len(in))
	copy(out, in)
	for i := 1; i < len(out); i++ {
		for j := i; j > 0 && out[j-1] > out[j]; j-- {
			out[j-1], out[j] = out[j], out[j-1]
		}
	}
	return out
}
