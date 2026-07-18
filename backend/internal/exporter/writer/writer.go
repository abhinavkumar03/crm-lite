// Package writer serializes a resolved column/row set into a downloadable file.
// It is format-agnostic to the rest of the export engine: callers depend only on
// Write, so adding a new output format is a single new function here.
package writer

import (
	"bytes"
	"encoding/csv"
	"errors"
	"fmt"
	"strconv"
	"strings"

	fieldentity "github.com/abhinavkumar03/crm-lite/backend/internal/field/entity"
	"github.com/xuri/excelize/v2"
)

// Supported output formats.
const (
	FormatCSV  = "csv"
	FormatXLSX = "xlsx"
)

// ErrUnsupported is returned for an unknown format.
var ErrUnsupported = errors.New("unsupported export format (use csv or xlsx)")

// Column describes one output column: which key to read from each row, the human
// header to print, and the field type used to format the cell.
type Column struct {
	Key   string
	Label string
	Type  string
}

// Result is the serialized file plus the metadata needed to serve it.
type Result struct {
	Content     []byte
	ContentType string
	Ext         string
}

// Write serializes rows using the given columns in the requested format. Each
// row is a map keyed by Column.Key; values are formatted per the column type.
func Write(format string, columns []Column, rows []map[string]any) (*Result, error) {
	switch format {
	case FormatCSV:
		content, err := writeCSV(columns, rows)
		if err != nil {
			return nil, err
		}
		return &Result{Content: content, ContentType: "text/csv", Ext: "csv"}, nil
	case FormatXLSX:
		content, err := writeXLSX(columns, rows)
		if err != nil {
			return nil, err
		}
		return &Result{
			Content:     content,
			ContentType: "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet",
			Ext:         "xlsx",
		}, nil
	default:
		return nil, ErrUnsupported
	}
}

func writeCSV(columns []Column, rows []map[string]any) ([]byte, error) {
	var buf bytes.Buffer
	// UTF-8 BOM so Excel opens accented characters correctly.
	buf.WriteString("\ufeff")

	w := csv.NewWriter(&buf)

	header := make([]string, len(columns))
	for i, c := range columns {
		header[i] = c.Label
	}
	if err := w.Write(header); err != nil {
		return nil, err
	}

	record := make([]string, len(columns))
	for _, row := range rows {
		for i, c := range columns {
			record[i] = FormatCell(c.Type, row[c.Key])
		}
		if err := w.Write(record); err != nil {
			return nil, err
		}
	}

	w.Flush()
	if err := w.Error(); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func writeXLSX(columns []Column, rows []map[string]any) ([]byte, error) {
	f := excelize.NewFile()
	defer f.Close()

	const sheet = "Sheet1"

	header := make([]any, len(columns))
	for i, c := range columns {
		header[i] = c.Label
	}
	if err := f.SetSheetRow(sheet, "A1", &header); err != nil {
		return nil, err
	}

	cells := make([]any, len(columns))
	for r, row := range rows {
		for i, c := range columns {
			cells[i] = FormatCell(c.Type, row[c.Key])
		}
		axis, err := excelize.CoordinatesToCellName(1, r+2)
		if err != nil {
			return nil, err
		}
		if err := f.SetSheetRow(sheet, axis, &cells); err != nil {
			return nil, err
		}
	}

	out, err := f.WriteToBuffer()
	if err != nil {
		return nil, err
	}
	return out.Bytes(), nil
}

// FormatCell renders a single value as a string appropriate for its field type.
func FormatCell(fieldType string, v any) string {
	if v == nil {
		return ""
	}

	switch fieldType {
	case fieldentity.TypeBoolean, fieldentity.TypeCheckbox:
		if b, ok := v.(bool); ok {
			return strconv.FormatBool(b)
		}
	case fieldentity.TypeNumber, fieldentity.TypeCurrency:
		if f, ok := v.(float64); ok {
			return strconv.FormatFloat(f, 'f', -1, 64)
		}
	case fieldentity.TypeMultiselect:
		if arr, ok := v.([]any); ok {
			parts := make([]string, 0, len(arr))
			for _, item := range arr {
				parts = append(parts, fmt.Sprint(item))
			}
			return strings.Join(parts, ", ")
		}
	}

	switch t := v.(type) {
	case string:
		return t
	case float64:
		return strconv.FormatFloat(t, 'f', -1, 64)
	case bool:
		return strconv.FormatBool(t)
	default:
		return fmt.Sprint(t)
	}
}
