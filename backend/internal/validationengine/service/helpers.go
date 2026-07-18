package service

import (
	"encoding/json"
	"fmt"
	"net/url"
	"regexp"
	"strconv"
	"strings"

	fieldentity "github.com/abhinavkumar03/crm-lite/backend/internal/field/entity"
)

var emailRegex = regexp.MustCompile(`^[^@\s]+@[^@\s]+\.[^@\s]+$`)

// stringTypes are field types whose values are treated as text for length /
// regex / format checks.
var stringTypes = map[string]bool{
	fieldentity.TypeText:     true,
	fieldentity.TypeTextarea: true,
	fieldentity.TypeEmail:    true,
	fieldentity.TypePhone:    true,
	fieldentity.TypeURL:      true,
	fieldentity.TypeRichtext: true,
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

func toString(v any) string {
	switch t := v.(type) {
	case string:
		return t
	case float64:
		return strconv.FormatFloat(t, 'f', -1, 64)
	case bool:
		return strconv.FormatBool(t)
	case nil:
		return ""
	default:
		return fmt.Sprint(t)
	}
}

func toFloat(v any) (float64, bool) {
	switch t := v.(type) {
	case float64:
		return t, true
	case int:
		return float64(t), true
	case string:
		f, err := strconv.ParseFloat(strings.TrimSpace(t), 64)
		return f, err == nil
	default:
		return 0, false
	}
}

func isURL(s string) bool {
	u, err := url.ParseRequestURI(strings.TrimSpace(s))
	return err == nil && u.Scheme != "" && u.Host != ""
}

func runeLen(s string) int {
	return len([]rune(s))
}

// parseParams decodes a rule's raw JSONB params into a generic map.
func parseParams(raw []byte) map[string]any {
	m := map[string]any{}
	if len(raw) == 0 {
		return m
	}
	_ = json.Unmarshal(raw, &m)
	return m
}

// paramFloat reads a numeric param ("value").
func paramFloat(params map[string]any, key string) (float64, bool) {
	v, ok := params[key]
	if !ok {
		return 0, false
	}
	return toFloat(v)
}

func paramString(params map[string]any, key string) (string, bool) {
	v, ok := params[key]
	if !ok {
		return "", false
	}
	s, ok := v.(string)
	return s, ok
}

// paramStrings reads a string-slice param ("values").
func paramStrings(params map[string]any, key string) []string {
	v, ok := params[key]
	if !ok {
		return nil
	}
	arr, ok := v.([]any)
	if !ok {
		return nil
	}
	out := make([]string, 0, len(arr))
	for _, item := range arr {
		out = append(out, toString(item))
	}
	return out
}

// parseOptionValues extracts the allowed values from a field's stored options,
// accepting both ["A","B"] and [{"label":..,"value":..}] shapes.
func parseOptionValues(raw []byte) []string {
	if len(raw) == 0 {
		return nil
	}
	var strs []string
	if err := json.Unmarshal(raw, &strs); err == nil {
		return strs
	}
	var objs []struct {
		Value string `json:"value"`
	}
	if err := json.Unmarshal(raw, &objs); err == nil {
		out := make([]string, 0, len(objs))
		for _, o := range objs {
			out = append(out, o.Value)
		}
		return out
	}
	return nil
}

func contains(values []string, target string) bool {
	for _, v := range values {
		if v == target {
			return true
		}
	}
	return false
}

// pick returns the first non-empty message from the candidates, falling back to
// def. Used to honor custom error messages (rule > field > engine default).
func pick(candidates []*string, def string) string {
	for _, c := range candidates {
		if c != nil && strings.TrimSpace(*c) != "" {
			return *c
		}
	}
	return def
}
