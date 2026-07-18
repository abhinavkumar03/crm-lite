package writer

import (
	"bytes"
	"encoding/csv"
	"strings"
	"testing"

	fieldentity "github.com/abhinavkumar03/crm-lite/backend/internal/field/entity"
)

func columns() []Column {
	return []Column{
		{Key: "name", Label: "Name", Type: fieldentity.TypeText},
		{Key: "amount", Label: "Amount", Type: fieldentity.TypeCurrency},
		{Key: "active", Label: "Active", Type: fieldentity.TypeBoolean},
		{Key: "tags", Label: "Tags", Type: fieldentity.TypeMultiselect},
	}
}

func TestFormatCell(t *testing.T) {
	cases := []struct {
		typ  string
		in   any
		want string
	}{
		{fieldentity.TypeText, "hi", "hi"},
		{fieldentity.TypeCurrency, float64(19.99), "19.99"},
		{fieldentity.TypeNumber, float64(42), "42"},
		{fieldentity.TypeBoolean, true, "true"},
		{fieldentity.TypeMultiselect, []any{"a", "b"}, "a, b"},
		{fieldentity.TypeText, nil, ""},
	}
	for _, c := range cases {
		if got := FormatCell(c.typ, c.in); got != c.want {
			t.Fatalf("FormatCell(%q, %v) = %q, want %q", c.typ, c.in, got, c.want)
		}
	}
}

func TestWriteCSV(t *testing.T) {
	rows := []map[string]any{
		{"name": "Ada", "amount": float64(100), "active": true, "tags": []any{"vip", "eu"}},
		{"name": "Grace", "amount": nil, "active": false, "tags": nil},
	}
	res, err := Write(FormatCSV, columns(), rows)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.Ext != "csv" {
		t.Fatalf("ext = %q, want csv", res.Ext)
	}

	// Strip the BOM before parsing.
	body := bytes.TrimPrefix(res.Content, []byte("\ufeff"))
	records, err := csv.NewReader(strings.NewReader(string(body))).ReadAll()
	if err != nil {
		t.Fatalf("csv parse: %v", err)
	}
	if len(records) != 3 {
		t.Fatalf("records = %d, want 3 (header + 2 rows)", len(records))
	}
	if records[0][0] != "Name" || records[0][3] != "Tags" {
		t.Fatalf("header = %v", records[0])
	}
	if records[1][1] != "100" || records[1][2] != "true" || records[1][3] != "vip, eu" {
		t.Fatalf("row0 = %v", records[1])
	}
	if records[2][1] != "" {
		t.Fatalf("row1 amount = %q, want empty", records[2][1])
	}
}

func TestWriteXLSX(t *testing.T) {
	rows := []map[string]any{{"name": "Ada", "amount": float64(1), "active": true, "tags": nil}}
	res, err := Write(FormatXLSX, columns(), rows)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.Ext != "xlsx" || len(res.Content) == 0 {
		t.Fatalf("unexpected xlsx result: ext=%q size=%d", res.Ext, len(res.Content))
	}
}

func TestWriteUnsupported(t *testing.T) {
	if _, err := Write("pdf", columns(), nil); err != ErrUnsupported {
		t.Fatalf("err = %v, want ErrUnsupported", err)
	}
}
