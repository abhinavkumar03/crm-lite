package service

import (
	"bytes"
	"context"
	"encoding/csv"
	"strings"
	"testing"

	"github.com/abhinavkumar03/crm-lite/backend/internal/exporter/dto"
	"github.com/abhinavkumar03/crm-lite/backend/internal/exporter/entity"
	fieldentity "github.com/abhinavkumar03/crm-lite/backend/internal/field/entity"
	recorddto "github.com/abhinavkumar03/crm-lite/backend/internal/record/dto"
)

type fakeFields struct {
	fields []fieldentity.Field
}

func (f *fakeFields) ModuleStorage(context.Context, string, string) (string, bool, error) {
	return "dynamic", true, nil
}
func (f *fakeFields) List(context.Context, string, string) ([]fieldentity.Field, error) {
	return f.fields, nil
}

// fakeRows returns preconfigured pages, capturing the last query it received so
// tests can assert expansion/paging behaviour.
type fakeRows struct {
	pages   []*recorddto.ListResult
	call    int
	lastExp bool
}

func (r *fakeRows) List(_ context.Context, _, _ string, q recorddto.ListQuery) (*recorddto.ListResult, error) {
	r.lastExp = q.Expand
	if r.call >= len(r.pages) {
		return &recorddto.ListResult{}, nil
	}
	page := r.pages[r.call]
	r.call++
	return page, nil
}

func schema() *fakeFields {
	return &fakeFields{
		fields: []fieldentity.Field{
			{APIName: "name", Label: "Name", FieldType: fieldentity.TypeText, IsVisible: true},
			{APIName: "amount", Label: "Amount", FieldType: fieldentity.TypeCurrency, IsVisible: true},
			{APIName: "owner", Label: "Owner", FieldType: fieldentity.TypeUser, IsVisible: true},
			{APIName: "secret", Label: "Secret", FieldType: fieldentity.TypeText, IsVisible: false},
		},
	}
}

func TestBuildDefaultColumnsCSV(t *testing.T) {
	rows := &fakeRows{
		pages: []*recorddto.ListResult{
			{
				Records: []recorddto.RecordResponse{
					{
						Data:      map[string]any{"name": "Ada", "amount": float64(100), "owner": "u1", "secret": "hidden"},
						Relations: map[string]recorddto.RelationRef{"owner": {ID: "u1", Label: "Ada Lovelace"}},
					},
				},
				TotalPages: 1,
			},
		},
	}
	svc := New(nil, nil, rows, schema(), nil)

	result, count, err := svc.Build(context.Background(), "org", "mod", dto.ExportSpec{Format: "csv"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if count != 1 {
		t.Fatalf("row count = %d, want 1", count)
	}
	// A user column forces relation expansion so labels resolve.
	if !rows.lastExp {
		t.Fatal("expected expand to be requested for relation columns")
	}

	records := parseCSV(t, result.Content)
	if len(records) != 2 {
		t.Fatalf("csv records = %d, want 2", len(records))
	}
	// Only visible fields, in order; "secret" excluded.
	if strings.Join(records[0], ",") != "Name,Amount,Owner" {
		t.Fatalf("header = %v", records[0])
	}
	// Owner shows the resolved relation label, not the raw id.
	if records[1][2] != "Ada Lovelace" {
		t.Fatalf("owner cell = %q, want resolved label", records[1][2])
	}
	if records[1][1] != "100" {
		t.Fatalf("amount cell = %q", records[1][1])
	}
}

func TestBuildExplicitColumnsWithMeta(t *testing.T) {
	rows := &fakeRows{
		pages: []*recorddto.ListResult{
			{
				Records:    []recorddto.RecordResponse{{ID: "rec-1", Data: map[string]any{"name": "Ada"}}},
				TotalPages: 1,
			},
		},
	}
	svc := New(nil, nil, rows, schema(), nil)

	result, _, err := svc.Build(context.Background(), "org", "mod", dto.ExportSpec{
		Format:  "csv",
		Columns: []string{"id", "name", "ghost"}, // ghost is unknown -> dropped
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	records := parseCSV(t, result.Content)
	if strings.Join(records[0], ",") != "ID,Name" {
		t.Fatalf("header = %v, want ID,Name", records[0])
	}
	if records[1][0] != "rec-1" {
		t.Fatalf("id cell = %q", records[1][0])
	}
}

func TestBuildNoColumns(t *testing.T) {
	// A schema with no visible fields yields no columns.
	empty := &fakeFields{fields: []fieldentity.Field{{APIName: "x", Label: "X", FieldType: "text", IsVisible: false}}}
	svc := New(nil, nil, &fakeRows{}, empty, nil)
	if _, _, err := svc.Build(context.Background(), "o", "m", dto.ExportSpec{Format: "csv"}); err != ErrNoColumns {
		t.Fatalf("err = %v, want ErrNoColumns", err)
	}
}

func TestNormalizeFormat(t *testing.T) {
	if normalizeFormat("xlsx") != entity.FormatXLSX {
		t.Fatal("xlsx should be preserved")
	}
	if normalizeFormat("") != entity.FormatCSV || normalizeFormat("weird") != entity.FormatCSV {
		t.Fatal("unknown/empty format should default to csv")
	}
}

func parseCSV(t *testing.T, content []byte) [][]string {
	t.Helper()
	body := bytes.TrimPrefix(content, []byte("\ufeff"))
	records, err := csv.NewReader(strings.NewReader(string(body))).ReadAll()
	if err != nil {
		t.Fatalf("csv parse: %v", err)
	}
	return records
}
