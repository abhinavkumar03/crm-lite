package notify

import (
	"fmt"
	"regexp"
)

// placeholderRe matches {{ key }} tokens (optionally spaced). Keys are simple
// dotted identifiers referencing entries in the message Data map.
var placeholderRe = regexp.MustCompile(`\{\{\s*([a-zA-Z0-9_.]+)\s*\}\}`)

// Render substitutes {{key}} tokens in text with values from data. Unknown keys
// are replaced with an empty string so a half-populated payload never leaks the
// raw template token to a recipient. This keeps templating dependency-free while
// still supporting server-rendered notifications.
func Render(text string, data map[string]any) string {
	if text == "" || len(data) == 0 {
		return text
	}
	return placeholderRe.ReplaceAllStringFunc(text, func(token string) string {
		match := placeholderRe.FindStringSubmatch(token)
		if len(match) < 2 {
			return ""
		}
		if v, ok := data[match[1]]; ok {
			return fmt.Sprintf("%v", v)
		}
		return ""
	})
}
