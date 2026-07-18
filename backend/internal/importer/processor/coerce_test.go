package processor

import (
	"reflect"
	"testing"

	fieldentity "github.com/abhinavkumar03/crm-lite/backend/internal/field/entity"
)

func TestCoerce(t *testing.T) {
	cases := []struct {
		name      string
		fieldType string
		raw       string
		want      any
	}{
		{"number ok", fieldentity.TypeNumber, "42", float64(42)},
		{"currency decimal", fieldentity.TypeCurrency, "19.99", float64(19.99)},
		{"number invalid falls back to raw", fieldentity.TypeNumber, "N/A", "N/A"},
		{"boolean yes", fieldentity.TypeBoolean, "Yes", true},
		{"checkbox 0", fieldentity.TypeCheckbox, "0", false},
		{"boolean invalid falls back", fieldentity.TypeBoolean, "maybe", "maybe"},
		{"multiselect split", fieldentity.TypeMultiselect, "a, b ; c", []any{"a", "b", "c"}},
		{"text passthrough", fieldentity.TypeText, "hello", "hello"},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := coerce(tc.fieldType, tc.raw)
			if !reflect.DeepEqual(got, tc.want) {
				t.Fatalf("coerce(%q, %q) = %#v, want %#v", tc.fieldType, tc.raw, got, tc.want)
			}
		})
	}
}
