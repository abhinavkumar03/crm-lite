package service

import (
	"time"

	"github.com/abhinavkumar03/crm-lite/backend/internal/exporter/writer"
	fieldentity "github.com/abhinavkumar03/crm-lite/backend/internal/field/entity"
	recorddto "github.com/abhinavkumar03/crm-lite/backend/internal/record/dto"
)

// metaColumns are record-level attributes that are always exportable even though
// they are not user-defined fields.
var metaColumns = map[string]writer.Column{
	"id":         {Key: "id", Label: "ID", Type: fieldentity.TypeText},
	"created_at": {Key: "created_at", Label: "Created At", Type: fieldentity.TypeDatetime},
	"updated_at": {Key: "updated_at", Label: "Updated At", Type: fieldentity.TypeDatetime},
}

// resolveColumns turns a requested api_name list into ordered writer columns. An
// empty request defaults to every visible field (in sort order). Unknown names
// are dropped so a stale template can never break an export.
func resolveColumns(fields []fieldentity.Field, requested []string) []writer.Column {
	byAPI := make(map[string]fieldentity.Field, len(fields))
	for i := range fields {
		byAPI[fields[i].APIName] = fields[i]
	}

	if len(requested) == 0 {
		cols := make([]writer.Column, 0, len(fields))
		for i := range fields {
			if fields[i].IsVisible {
				cols = append(cols, fieldColumn(fields[i]))
			}
		}
		return cols
	}

	cols := make([]writer.Column, 0, len(requested))
	for _, name := range requested {
		if f, ok := byAPI[name]; ok {
			cols = append(cols, fieldColumn(f))
			continue
		}
		if mc, ok := metaColumns[name]; ok {
			cols = append(cols, mc)
		}
	}
	return cols
}

func fieldColumn(f fieldentity.Field) writer.Column {
	return writer.Column{Key: f.APIName, Label: f.Label, Type: f.FieldType}
}

// needsExpand reports whether any selected column is a relation, so the record
// runtime resolves lookup/user ids into human labels for the export.
func needsExpand(columns []writer.Column) bool {
	for _, c := range columns {
		if c.Type == fieldentity.TypeLookup || c.Type == fieldentity.TypeUser {
			return true
		}
	}
	return false
}

// buildRows projects records onto the chosen columns, preferring resolved
// relation labels over raw ids and stringifying record timestamps.
func buildRows(records []recorddto.RecordResponse, columns []writer.Column) []map[string]any {
	out := make([]map[string]any, 0, len(records))
	for i := range records {
		rec := records[i]
		row := make(map[string]any, len(columns))
		for _, c := range columns {
			switch c.Key {
			case "id":
				row[c.Key] = rec.ID
			case "created_at":
				row[c.Key] = rec.CreatedAt.Format(time.RFC3339)
			case "updated_at":
				row[c.Key] = rec.UpdatedAt.Format(time.RFC3339)
			default:
				if rel, ok := rec.Relations[c.Key]; ok && rel.Label != "" {
					row[c.Key] = rel.Label
				} else {
					row[c.Key] = rec.Data[c.Key]
				}
			}
		}
		out = append(out, row)
	}
	return out
}
