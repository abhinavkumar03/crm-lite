package repository

import (
	"strings"
	"testing"

	fieldentity "github.com/abhinavkumar03/crm-lite/backend/internal/field/entity"
	"github.com/abhinavkumar03/crm-lite/backend/internal/record/dto"
)

func testMeta() map[string]FieldMeta {
	return map[string]FieldMeta{
		"name":   {Type: fieldentity.TypeText, Searchable: true, Filterable: true},
		"email":  {Type: fieldentity.TypeEmail, Searchable: true, Filterable: true},
		"amount": {Type: fieldentity.TypeNumber, Searchable: false, Filterable: true},
		"stage":  {Type: fieldentity.TypeDropdown, Searchable: false, Filterable: true},
		"secret": {Type: fieldentity.TypeText, Searchable: false, Filterable: false},
	}
}

func TestBuildWhere_ScopeOnly(t *testing.T) {
	w := BuildWhere("org1", "mod1", dto.ListQuery{}, testMeta())

	if w.SQL != "organization_id = $1 AND module_id = $2" {
		t.Fatalf("unexpected SQL: %q", w.SQL)
	}
	if len(w.Args) != 2 || w.Args[0] != "org1" || w.Args[1] != "mod1" {
		t.Fatalf("unexpected args: %#v", w.Args)
	}
}

func TestBuildWhere_SearchAcrossSearchableFields(t *testing.T) {
	w := BuildWhere("o", "m", dto.ListQuery{Search: "acme"}, testMeta())

	// Only searchable fields (name, email) participate, ordered deterministically.
	if !strings.Contains(w.SQL, "data->>'email' ILIKE $3") ||
		!strings.Contains(w.SQL, "data->>'name' ILIKE $3") {
		t.Fatalf("search SQL missing searchable fields: %q", w.SQL)
	}
	if strings.Contains(w.SQL, "'amount'") || strings.Contains(w.SQL, "'secret'") {
		t.Fatalf("non-searchable field leaked into search: %q", w.SQL)
	}
	if len(w.Args) != 3 || w.Args[2] != "%acme%" {
		t.Fatalf("unexpected search args: %#v", w.Args)
	}
}

func TestBuildWhere_FilterEquality(t *testing.T) {
	w := BuildWhere("o", "m", dto.ListQuery{
		Filters: []dto.FilterClause{{Field: "stage", Operator: dto.OpEquals, Value: "open"}},
	}, testMeta())

	if !strings.Contains(w.SQL, "data->>'stage' = $3") {
		t.Fatalf("unexpected filter SQL: %q", w.SQL)
	}
	if len(w.Args) != 3 || w.Args[2] != "open" {
		t.Fatalf("unexpected filter args: %#v", w.Args)
	}
}

func TestBuildWhere_NumericGreaterThanCasts(t *testing.T) {
	w := BuildWhere("o", "m", dto.ListQuery{
		Filters: []dto.FilterClause{{Field: "amount", Operator: dto.OpGreaterThan, Value: 100}},
	}, testMeta())

	if !strings.Contains(w.SQL, "(data->>'amount')::numeric > $3::numeric") {
		t.Fatalf("numeric cast missing: %q", w.SQL)
	}
	if w.Args[2] != "100" {
		t.Fatalf("expected stringified value, got %#v", w.Args[2])
	}
}

func TestBuildWhere_InOperator(t *testing.T) {
	w := BuildWhere("o", "m", dto.ListQuery{
		Filters: []dto.FilterClause{{Field: "stage", Operator: dto.OpIn, Value: []any{"open", "won"}}},
	}, testMeta())

	if !strings.Contains(w.SQL, "data->>'stage' = ANY($3)") {
		t.Fatalf("unexpected IN SQL: %q", w.SQL)
	}
	arr, ok := w.Args[2].([]string)
	if !ok || len(arr) != 2 || arr[0] != "open" || arr[1] != "won" {
		t.Fatalf("unexpected IN args: %#v", w.Args[2])
	}
}

func TestBuildWhere_IgnoresUnknownAndNonFilterable(t *testing.T) {
	w := BuildWhere("o", "m", dto.ListQuery{
		Filters: []dto.FilterClause{
			{Field: "secret", Operator: dto.OpEquals, Value: "x"}, // not filterable
			{Field: "ghost", Operator: dto.OpEquals, Value: "y"},  // unknown
			{Field: "name", Operator: "bogus", Value: "z"},        // bad operator
		},
	}, testMeta())

	if w.SQL != "organization_id = $1 AND module_id = $2" {
		t.Fatalf("expected filters to be ignored, got: %q", w.SQL)
	}
	if len(w.Args) != 2 {
		t.Fatalf("expected no extra args, got: %#v", w.Args)
	}
}

func TestBuildWhere_RejectsInjectionInFilterField(t *testing.T) {
	w := BuildWhere("o", "m", dto.ListQuery{
		Filters: []dto.FilterClause{{Field: "name'; DROP TABLE records;--", Operator: dto.OpEquals, Value: "x"}},
	}, testMeta())

	if strings.Contains(w.SQL, "DROP TABLE") {
		t.Fatalf("injection leaked into SQL: %q", w.SQL)
	}
}

func TestBuildOrderBy(t *testing.T) {
	meta := testMeta()

	cases := map[string]struct {
		sort, order, want string
	}{
		"default":     {"", "", "ORDER BY created_at DESC"},
		"updated asc": {"updated_at", "asc", "ORDER BY updated_at ASC"},
		"field text":  {"name", "asc", "ORDER BY data->>'name' ASC NULLS LAST"},
		"field num":   {"amount", "desc", "ORDER BY (data->>'amount')::numeric DESC NULLS LAST"},
		"unknown":     {"ghost", "asc", "ORDER BY created_at ASC"},
	}

	for name, tc := range cases {
		if got := BuildOrderBy(tc.sort, tc.order, meta); got != tc.want {
			t.Errorf("%s: got %q want %q", name, got, tc.want)
		}
	}
}
