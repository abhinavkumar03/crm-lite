package service

import (
	"testing"

	fieldentity "github.com/abhinavkumar03/crm-lite/backend/internal/field/entity"
)

func fields() []fieldentity.Field {
	return []fieldentity.Field{
		{APIName: "first_name", Label: "First Name", FieldType: fieldentity.TypeText},
		{APIName: "amount", Label: "Deal Amount", FieldType: fieldentity.TypeCurrency},
		{APIName: "score", Label: "Score", FieldType: fieldentity.TypeFormula}, // derived: not a target
		{APIName: "created_at", Label: "Created", FieldType: fieldentity.TypeDate, IsSystem: true},
	}
}

func TestSuggestMapping(t *testing.T) {
	headers := []string{"First Name", "amount", "Unknown"}
	got := suggestMapping(headers, fields())

	if got["First Name"] != "first_name" {
		t.Fatalf(`suggested["First Name"] = %q, want "first_name"`, got["First Name"])
	}
	if got["amount"] != "amount" {
		t.Fatalf(`suggested["amount"] = %q, want "amount"`, got["amount"])
	}
	if _, ok := got["Unknown"]; ok {
		t.Fatalf("unmatched header should not be suggested, got %v", got)
	}
}

func TestSanitizeMappingDropsInvalidTargets(t *testing.T) {
	headers := []string{"First Name", "Deal", "Ghost"}
	mapping := map[string]string{
		"First Name": "first_name", // valid
		"Deal":       "amount",     // valid
		"Ghost":      "score",      // formula field -> not writable -> dropped
		"Missing":    "first_name", // header not in file -> dropped
	}

	got := sanitizeMapping(mapping, headers, fields())
	if len(got) != 2 {
		t.Fatalf("sanitized mapping = %v, want 2 entries", got)
	}
	if got["First Name"] != "first_name" || got["Deal"] != "amount" {
		t.Fatalf("sanitized mapping = %v", got)
	}
	if _, ok := got["Ghost"]; ok {
		t.Fatalf("formula target should be dropped, got %v", got)
	}
}

func TestNormalizeKey(t *testing.T) {
	cases := map[string]string{
		"First Name": "firstname",
		"first_name": "firstname",
		"firstName":  "firstname",
		"  E-mail  ": "email",
	}
	for in, want := range cases {
		if got := normalizeKey(in); got != want {
			t.Fatalf("normalizeKey(%q) = %q, want %q", in, got, want)
		}
	}
}
