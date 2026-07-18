package service

import (
	"strings"

	fieldentity "github.com/abhinavkumar03/crm-lite/backend/internal/field/entity"
)

// suggestMapping auto-matches source columns to module fields by comparing a
// normalized form of each header against each field's api_name and label. The
// result is header -> field api_name for every confident match; unmatched
// columns are simply left out for the user to map manually.
func suggestMapping(headers []string, fields []fieldentity.Field) map[string]string {
	byKey := map[string]string{} // normalized field key -> api_name
	for i := range fields {
		f := fields[i]
		if !isWritable(f) {
			continue
		}
		byKey[normalizeKey(f.APIName)] = f.APIName
		byKey[normalizeKey(f.Label)] = f.APIName
	}

	out := map[string]string{}
	used := map[string]bool{}
	for _, h := range headers {
		if api, ok := byKey[normalizeKey(h)]; ok && !used[api] {
			out[h] = api
			used[api] = true
		}
	}
	return out
}

// sanitizeMapping filters a client-provided mapping down to entries that are
// safe to persist: the source column must exist in the file and the target must
// be a writable field of the module. Each field may be targeted only once.
func sanitizeMapping(mapping map[string]string, headers []string, fields []fieldentity.Field) map[string]string {
	headerSet := map[string]bool{}
	for _, h := range headers {
		headerSet[h] = true
	}
	writable := map[string]bool{}
	for i := range fields {
		if isWritable(fields[i]) {
			writable[fields[i].APIName] = true
		}
	}

	out := map[string]string{}
	usedFields := map[string]bool{}
	for header, api := range mapping {
		if api == "" || !headerSet[header] || !writable[api] || usedFields[api] {
			continue
		}
		out[header] = api
		usedFields[api] = true
	}
	return out
}

// isWritable excludes read-only, system and derived fields from import targets.
func isWritable(f fieldentity.Field) bool {
	if f.IsReadOnly || f.IsSystem {
		return false
	}
	switch f.FieldType {
	case fieldentity.TypeFormula:
		return false
	default:
		return true
	}
}

// normalizeKey lowercases and strips non-alphanumeric characters so "First Name",
// "first_name" and "firstName" all collapse to the same key.
func normalizeKey(s string) string {
	var b strings.Builder
	for _, r := range strings.ToLower(strings.TrimSpace(s)) {
		if (r >= 'a' && r <= 'z') || (r >= '0' && r <= '9') {
			b.WriteRune(r)
		}
	}
	return b.String()
}
