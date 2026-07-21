package engine

import (
	"fmt"
	"strings"
	"time"

	"github.com/abhinavkumar03/crm-lite/backend/internal/notify"
)

// VariableContext is the merge bag for templates and action configs.
type VariableContext struct {
	Record     map[string]any
	ModuleAPI  string
	OwnerName  string
	OwnerEmail string
	OrgName    string
}

// BuildMergeMap produces keys for notify.Render and nested module alias.
func BuildMergeMap(vc VariableContext) map[string]any {
	data := map[string]any{
		"today":        time.Now().UTC().Format("2006-01-02"),
		"current_date": time.Now().UTC().Format("2006-01-02"),
		"workspace": map[string]any{
			"name": vc.OrgName,
		},
		"owner": map[string]any{
			"name":  vc.OwnerName,
			"email": vc.OwnerEmail,
		},
		"record": vc.Record,
	}
	if vc.ModuleAPI != "" {
		data[vc.ModuleAPI] = vc.Record
	}
	// Flat aliases for convenience.
	for k, v := range vc.Record {
		data[k] = v
		data["record."+k] = v
		if vc.ModuleAPI != "" {
			data[vc.ModuleAPI+"."+k] = v
		}
	}
	data["owner.name"] = vc.OwnerName
	data["owner.email"] = vc.OwnerEmail
	data["workspace.name"] = vc.OrgName
	return data
}

// RenderString replaces {{vars}} in text.
func RenderString(text string, data map[string]any) string {
	return notify.Render(text, data)
}

// RenderMap deep-renders string values in a config map.
func RenderMap(cfg map[string]any, data map[string]any) map[string]any {
	if cfg == nil {
		return map[string]any{}
	}
	out := make(map[string]any, len(cfg))
	for k, v := range cfg {
		out[k] = renderValue(v, data)
	}
	return out
}

func renderValue(v any, data map[string]any) any {
	switch t := v.(type) {
	case string:
		return RenderString(t, data)
	case map[string]any:
		return RenderMap(t, data)
	case []any:
		arr := make([]any, len(t))
		for i, item := range t {
			arr[i] = renderValue(item, data)
		}
		return arr
	default:
		return v
	}
}

// ResolveRecipient picks to address from config or record fields.
func ResolveRecipient(cfg map[string]any, record map[string]any, channel string) string {
	if to, ok := cfg["to"].(string); ok && strings.TrimSpace(to) != "" {
		return strings.TrimSpace(to)
	}
	if channel == "email" {
		for _, k := range []string{"email", "Email", "work_email"} {
			if v, ok := record[k].(string); ok && v != "" {
				return v
			}
		}
	}
	if channel == "whatsapp" {
		for _, k := range []string{"phone", "Phone", "mobile", "whatsapp"} {
			if v, ok := record[k].(string); ok && v != "" {
				return v
			}
		}
	}
	return ""
}

// ConfigString reads a string from config.
func ConfigString(cfg map[string]any, key string) string {
	if cfg == nil {
		return ""
	}
	if v, ok := cfg[key].(string); ok {
		return v
	}
	return ""
}

// ConfigInt reads an int from config.
func ConfigInt(cfg map[string]any, key string, fallback int) int {
	if cfg == nil {
		return fallback
	}
	switch v := cfg[key].(type) {
	case float64:
		return int(v)
	case int:
		return v
	case string:
		var n int
		_, _ = fmt.Sscanf(v, "%d", &n)
		if n != 0 {
			return n
		}
	}
	return fallback
}
