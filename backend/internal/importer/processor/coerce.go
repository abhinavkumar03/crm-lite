package processor

import (
	"strconv"
	"strings"

	fieldentity "github.com/abhinavkumar03/crm-lite/backend/internal/field/entity"
)

// coerce converts a raw spreadsheet string into the Go type the validation
// engine and JSONB storage expect for the given field type. On a failed
// conversion it returns the raw string so the validator surfaces a clear,
// field-specific error rather than silently dropping the value.
func coerce(fieldType, raw string) any {
	switch fieldType {
	case fieldentity.TypeNumber, fieldentity.TypeCurrency:
		if f, err := strconv.ParseFloat(strings.TrimSpace(raw), 64); err == nil {
			return f
		}
		return raw
	case fieldentity.TypeBoolean, fieldentity.TypeCheckbox:
		if b, ok := parseBool(raw); ok {
			return b
		}
		return raw
	case fieldentity.TypeMultiselect:
		return splitList(raw)
	default:
		return raw
	}
}

func parseBool(raw string) (bool, bool) {
	switch strings.ToLower(strings.TrimSpace(raw)) {
	case "true", "1", "yes", "y", "on":
		return true, true
	case "false", "0", "no", "n", "off":
		return false, true
	default:
		return false, false
	}
}

// splitList turns "a, b ; c" into an []any so the validator's option check
// (which type-asserts []any) can evaluate each selection.
func splitList(raw string) []any {
	fields := strings.FieldsFunc(raw, func(r rune) bool {
		return r == ',' || r == ';'
	})
	out := make([]any, 0, len(fields))
	for _, f := range fields {
		if v := strings.TrimSpace(f); v != "" {
			out = append(out, v)
		}
	}
	return out
}
