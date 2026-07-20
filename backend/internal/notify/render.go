package notify

import (
	"fmt"
	"regexp"
	"strings"
)

// placeholderRe matches {{ key }} tokens (optionally spaced). Keys may be
// dotted identifiers (e.g. lead.name) referencing nested Data maps or flat keys.
var placeholderRe = regexp.MustCompile(`\{\{\s*([a-zA-Z0-9_.]+)\s*\}\}`)

// Render substitutes {{key}} tokens in text with values from data. Unknown keys
// are replaced with an empty string so a half-populated payload never leaks the
// raw template token to a recipient.
func Render(text string, data map[string]any) string {
	if text == "" || len(data) == 0 {
		return text
	}
	return placeholderRe.ReplaceAllStringFunc(text, func(token string) string {
		match := placeholderRe.FindStringSubmatch(token)
		if len(match) < 2 {
			return ""
		}
		if v, ok := lookup(data, match[1]); ok {
			return fmt.Sprintf("%v", v)
		}
		return ""
	})
}

// lookup resolves flat keys first, then dotted nested paths (lead.name).
func lookup(data map[string]any, key string) (any, bool) {
	if v, ok := data[key]; ok {
		return v, true
	}
	parts := strings.Split(key, ".")
	if len(parts) < 2 {
		return nil, false
	}
	var cur any = data
	for _, part := range parts {
		m, ok := cur.(map[string]any)
		if !ok {
			return nil, false
		}
		cur, ok = m[part]
		if !ok {
			return nil, false
		}
	}
	return cur, true
}
