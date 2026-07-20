package service

import (
	"errors"
	"testing"

	"github.com/abhinavkumar03/crm-lite/backend/internal/workspace/dto"
	"github.com/abhinavkumar03/crm-lite/backend/internal/workspace/repository"
)

func TestResolveEditable(t *testing.T) {
	tests := []struct {
		name       string
		readOnly   bool
		lockMode   string
		mode       string
		editableBy string
		wantEdit   bool
		wantLocked bool
	}{
		{"create never", false, "never", "create", "ALL", true, false},
		{"edit never", false, "never", "edit", "ALL", true, false},
		{"create after_create", false, "after_create", "create", "ALL", true, false},
		{"edit after_create", false, "after_create", "edit", "ALL", false, true},
		{"create always", false, "always", "create", "ALL", false, true},
		{"edit always", false, "always", "edit", "ALL", false, true},
		{"read_only wins", true, "never", "create", "ALL", false, false},
		{"non ALL editable_by", false, "never", "create", "OWNER", false, false},
		{"empty mode defaults create", false, "after_create", "", "ALL", true, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotEdit, gotLocked := ResolveEditable(tt.readOnly, tt.lockMode, tt.mode, tt.editableBy)
			if gotEdit != tt.wantEdit || gotLocked != tt.wantLocked {
				t.Fatalf("ResolveEditable() = (%v,%v), want (%v,%v)", gotEdit, gotLocked, tt.wantEdit, tt.wantLocked)
			}
		})
	}
}

func TestIsListColumnLocked(t *testing.T) {
	if !IsListColumnLocked(dto.ActionsColumnKey, false) {
		t.Fatal("actions should be locked")
	}
	if !IsListColumnLocked("name", true) {
		t.Fatal("system name should be locked")
	}
	if IsListColumnLocked("name", false) {
		t.Fatal("custom name should not be locked")
	}
	if IsListColumnLocked("phone", true) {
		t.Fatal("system phone should not be locked")
	}
}

func TestDefaultListVisible(t *testing.T) {
	if !defaultListVisible(true, "phone") {
		t.Fatal("system fields should default visible")
	}
	if defaultListVisible(false, "custom_field") {
		t.Fatal("custom fields should default hidden")
	}
	if !defaultListVisible(true, "name") {
		t.Fatal("locked system name should be visible")
	}
}

func TestNormalizeListColumns_ensuresActions(t *testing.T) {
	cols := normalizeListColumns([]dto.ListColumn{
		{FieldKey: "name", Visible: true, Order: 1},
	})
	if len(cols) != 2 {
		t.Fatalf("expected 2 columns, got %d", len(cols))
	}
	last := cols[len(cols)-1]
	if last.FieldKey != dto.ActionsColumnKey || !last.System || !last.Visible || !last.Locked {
		t.Fatalf("actions column malformed: %+v", last)
	}
}

func TestNormalizeListColumns_forcesActionsSystem(t *testing.T) {
	cols := normalizeListColumns([]dto.ListColumn{
		{FieldKey: "name", Visible: true, Order: 1},
		{FieldKey: dto.ActionsColumnKey, Visible: false, Order: 2, System: false},
	})
	var actions *dto.ListColumn
	for i := range cols {
		if cols[i].FieldKey == dto.ActionsColumnKey {
			actions = &cols[i]
		}
	}
	if actions == nil {
		t.Fatal("missing actions")
	}
	if !actions.System || !actions.Visible || !actions.Locked {
		t.Fatalf("actions should be system+visible+locked: %+v", actions)
	}
}

func TestNormalizeFormConfig_rejectsEmpty(t *testing.T) {
	_, err := normalizeFormMode("bogus")
	if err != ErrInvalidMode {
		t.Fatalf("expected ErrInvalidMode, got %v", err)
	}
	mode, err := normalizeFormMode("EDIT")
	if err != nil || mode != "edit" {
		t.Fatalf("expected edit, got %q %v", mode, err)
	}
}

func TestNormalizeListRequest_appendsActions(t *testing.T) {
	known := map[string]repository.HydrateField{
		"name":  {APIName: "name", IsVisible: true, IsSearchable: true, IsSystem: true},
		"phone": {APIName: "phone", IsVisible: true},
	}
	cols, err := normalizeListRequest([]dto.ListColumn{
		{FieldKey: "phone", Visible: true, Order: 2},
		{FieldKey: "name", Visible: false, Order: 1}, // locked → forced visible
		{FieldKey: dto.ActionsColumnKey, Visible: false, Order: 99},
	}, known)
	if err != nil {
		t.Fatal(err)
	}
	if len(cols) != 3 {
		t.Fatalf("expected 3 cols, got %d", len(cols))
	}
	if cols[0].FieldKey != "name" || !cols[0].Visible || !cols[0].Locked {
		t.Fatalf("expected name first locked+visible: %+v", cols[0])
	}
	if cols[1].FieldKey != "phone" {
		t.Fatalf("expected phone second: %+v", cols[1])
	}
	last := cols[2]
	if last.FieldKey != dto.ActionsColumnKey || !last.System || !last.Visible {
		t.Fatalf("actions: %+v", last)
	}
}

func TestNormalizeListRequest_rejectsUnknown(t *testing.T) {
	known := map[string]repository.HydrateField{
		"name": {APIName: "name", IsVisible: true, IsSystem: true},
	}
	_, err := normalizeListRequest([]dto.ListColumn{
		{FieldKey: "bogus", Visible: true, Order: 1},
	}, known)
	if !errors.Is(err, ErrInvalidListCol) {
		t.Fatalf("expected ErrInvalidListCol, got %v", err)
	}
}

func TestNormalizeListRequest_rejectsDuplicateOrder(t *testing.T) {
	known := map[string]repository.HydrateField{
		"name":  {APIName: "name", IsVisible: true, IsSystem: true},
		"phone": {APIName: "phone", IsVisible: true},
	}
	_, err := normalizeListRequest([]dto.ListColumn{
		{FieldKey: "name", Visible: true, Order: 1},
		{FieldKey: "phone", Visible: true, Order: 1},
	}, known)
	if !errors.Is(err, ErrInvalidListCol) {
		t.Fatalf("expected ErrInvalidListCol, got %v", err)
	}
}

func TestReconcileListColumns_addsMissingCustomHidden(t *testing.T) {
	fields := []repository.HydrateField{
		{APIName: "name", IsVisible: true, IsSystem: true},
		{APIName: "custom_x", IsVisible: true, IsSystem: false},
	}
	cols, changed := reconcileListColumns([]dto.ListColumn{
		{FieldKey: "name", Visible: true, Order: 1},
	}, fields)
	if !changed {
		t.Fatal("expected change")
	}
	var custom *dto.ListColumn
	for i := range cols {
		if cols[i].FieldKey == "custom_x" {
			custom = &cols[i]
		}
	}
	if custom == nil {
		t.Fatal("missing custom_x")
	}
	if custom.Visible {
		t.Fatalf("custom field should default hidden: %+v", custom)
	}
}

func TestReconcileListColumns_dropsDeleted(t *testing.T) {
	fields := []repository.HydrateField{
		{APIName: "name", IsVisible: true, IsSystem: true},
	}
	cols, changed := reconcileListColumns([]dto.ListColumn{
		{FieldKey: "name", Visible: true, Order: 1},
		{FieldKey: "gone", Visible: true, Order: 2},
	}, fields)
	if !changed {
		t.Fatal("expected change")
	}
	for _, c := range cols {
		if c.FieldKey == "gone" {
			t.Fatal("deleted field should be dropped")
		}
	}
}

func TestBuildDefaultListColumns_systemVisibleCustomHidden(t *testing.T) {
	cols, err := buildDefaultListColumns([]repository.HydrateField{
		{APIName: "name", IsVisible: true, IsSystem: true, IsSearchable: true},
		{APIName: "custom_a", IsVisible: true, IsSystem: false},
	})
	if err != nil {
		t.Fatal(err)
	}
	if len(cols) != 3 {
		t.Fatalf("expected 3, got %d", len(cols))
	}
	if !cols[0].Visible || cols[0].FieldKey != "name" {
		t.Fatalf("name should be visible: %+v", cols[0])
	}
	if cols[1].Visible || cols[1].FieldKey != "custom_a" {
		t.Fatalf("custom should be hidden: %+v", cols[1])
	}
}

func TestEnsureFormOrphans(t *testing.T) {
	fields := []repository.HydrateField{
		{APIName: "name", IsVisible: true},
		{APIName: "phone", IsVisible: true},
		{APIName: "hidden", IsVisible: false},
		{APIName: "sys", IsVisible: true, IsSystem: true},
	}
	cfg := formLayoutConfig{
		Sections: []dto.LayoutSection{
			{Key: "identity", Label: "Identity", Fields: []string{"name"}},
		},
	}
	out, changed := ensureFormOrphans(cfg, fields)
	if !changed {
		t.Fatal("expected change")
	}
	got := out.Sections[0].Fields
	if len(got) != 2 || got[0] != "name" || got[1] != "phone" {
		t.Fatalf("expected name+phone orphans merge, got %v", got)
	}
}

func TestHydrateListColumns(t *testing.T) {
	known := map[string]repository.HydrateField{
		"name": {ID: "f1", APIName: "name", Label: "Name", IsSystem: true, IsVisible: true},
	}
	cols := hydrateListColumns([]dto.ListColumn{
		{FieldKey: "name", Visible: true, Order: 1},
		{FieldKey: dto.ActionsColumnKey, Visible: true, Order: 2, System: true},
	}, known)
	if cols[0].FieldID != "f1" || cols[0].Label != "Name" || !cols[0].Locked {
		t.Fatalf("hydrate name: %+v", cols[0])
	}
	if cols[1].Label != "Actions" || !cols[1].Locked {
		t.Fatalf("hydrate actions: %+v", cols[1])
	}
}
