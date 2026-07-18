// Package parser turns an uploaded spreadsheet (CSV or XLSX) into a normalized,
// header-keyed set of rows. It is intentionally format-agnostic to the rest of
// the import engine: callers depend only on Parsed, so a new format is a single
// new function here with no changes downstream.
package parser

import (
	"bytes"
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"strings"

	"github.com/xuri/excelize/v2"
)

var (
	// ErrEmptyFile is returned when the file has a header but no data rows.
	ErrEmptyFile = errors.New("file has no data rows")
	// ErrNoHeaders is returned when the file has no header row at all.
	ErrNoHeaders = errors.New("file has no header row")
	// ErrUnsupported is returned for extensions other than .csv/.xlsx.
	ErrUnsupported = errors.New("unsupported file type (use .csv or .xlsx)")
)

// Parsed is the normalized output: an ordered header list and rows keyed by
// header. Every row contains a value (possibly empty) for every header.
type Parsed struct {
	Headers []string
	Rows    []map[string]string
}

// ParseFile dispatches on the filename extension. data is the full file body.
func ParseFile(filename string, data []byte) (*Parsed, error) {
	switch {
	case hasSuffixFold(filename, ".csv"):
		return ParseCSV(bytes.NewReader(data))
	case hasSuffixFold(filename, ".xlsx"):
		return ParseExcel(data)
	default:
		return nil, ErrUnsupported
	}
}

// ParseCSV reads a comma-separated file. It is lenient about ragged rows so a
// trailing short row does not abort the whole import.
func ParseCSV(r io.Reader) (*Parsed, error) {
	reader := csv.NewReader(r)
	reader.FieldsPerRecord = -1
	reader.TrimLeadingSpace = true

	records, err := reader.ReadAll()
	if err != nil {
		return nil, fmt.Errorf("parser: read csv: %w", err)
	}
	return fromMatrix(records)
}

// ParseExcel reads the first worksheet of an .xlsx workbook.
func ParseExcel(data []byte) (*Parsed, error) {
	f, err := excelize.OpenReader(bytes.NewReader(data))
	if err != nil {
		return nil, fmt.Errorf("parser: open xlsx: %w", err)
	}
	defer f.Close()

	sheets := f.GetSheetList()
	if len(sheets) == 0 {
		return nil, ErrNoHeaders
	}
	rows, err := f.GetRows(sheets[0])
	if err != nil {
		return nil, fmt.Errorf("parser: read xlsx rows: %w", err)
	}
	return fromMatrix(rows)
}

// fromMatrix normalizes a raw string matrix (first row = headers) into Parsed.
func fromMatrix(matrix [][]string) (*Parsed, error) {
	if len(matrix) == 0 {
		return nil, ErrNoHeaders
	}

	headers := normalizeHeaders(matrix[0])
	if len(headers) == 0 {
		return nil, ErrNoHeaders
	}

	rows := make([]map[string]string, 0, len(matrix)-1)
	for _, raw := range matrix[1:] {
		if isBlankRow(raw) {
			continue
		}
		row := make(map[string]string, len(headers))
		for i, h := range headers {
			if i < len(raw) {
				row[h] = strings.TrimSpace(raw[i])
			} else {
				row[h] = ""
			}
		}
		rows = append(rows, row)
	}

	if len(rows) == 0 {
		return nil, ErrEmptyFile
	}
	return &Parsed{Headers: headers, Rows: rows}, nil
}

// normalizeHeaders trims each header, drops a leading UTF-8 BOM, and guarantees
// uniqueness (duplicates get a numeric suffix) so header-keyed rows never clash.
func normalizeHeaders(raw []string) []string {
	seen := map[string]int{}
	out := make([]string, 0, len(raw))
	for i, h := range raw {
		name := strings.TrimSpace(strings.TrimPrefix(h, "\ufeff"))
		if name == "" {
			name = fmt.Sprintf("column_%d", i+1)
		}
		if n, ok := seen[name]; ok {
			seen[name] = n + 1
			name = fmt.Sprintf("%s_%d", name, n+1)
		} else {
			seen[name] = 1
		}
		out = append(out, name)
	}
	return out
}

func isBlankRow(raw []string) bool {
	for _, v := range raw {
		if strings.TrimSpace(v) != "" {
			return false
		}
	}
	return true
}

func hasSuffixFold(s, suffix string) bool {
	return len(s) >= len(suffix) && strings.EqualFold(s[len(s)-len(suffix):], suffix)
}
