package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"sort"
	"strings"

	"github.com/abhinavkumar03/crm-lite/backend/internal/workspace/dto"
	"github.com/abhinavkumar03/crm-lite/backend/internal/workspace/repository"
)

var (
	ErrSectionNotEmpty = errors.New("section must be empty before delete")
	ErrSectionNotFound = errors.New("section not found")
	ErrSectionExists   = errors.New("section already exists")
	ErrSystemColumn    = errors.New("system columns cannot be hidden or removed")
	ErrLockedColumn    = errors.New("locked columns cannot be hidden")
	ErrInvalidListCol  = errors.New("invalid list column configuration")
	ErrInvalidMode     = errors.New("mode must be create or edit")
)

// lockedListKeys are system field api_names that cannot be hidden from listings.
var lockedListKeys = map[string]bool{
	"name":   true,
	"title":  true,
	"status": true,
	"stage":  true,
}

// IsListColumnLocked reports whether a listing column must stay visible.
func IsListColumnLocked(fieldKey string, isSystem bool) bool {
	if fieldKey == dto.ActionsColumnKey {
		return true
	}
	return isSystem && lockedListKeys[fieldKey]
}

// defaultListVisible is the default visibility for a newly added / reset column.
func defaultListVisible(isSystem bool, fieldKey string) bool {
	if IsListColumnLocked(fieldKey, isSystem) {
		return true
	}
	return isSystem
}

type formLayoutConfig struct {
	Sections []dto.LayoutSection `json:"sections"`
}

type listLayoutConfig struct {
	Columns []dto.ListColumn `json:"columns"`
}

// ResolveEditable applies lock_mode + read_only for create/edit form rendering.
func ResolveEditable(isReadOnly bool, lockMode, mode, editableBy string) (editable bool, locked bool) {
	mode = strings.ToLower(strings.TrimSpace(mode))
	if mode == "" {
		mode = "create"
	}
	lockMode = strings.TrimSpace(lockMode)
	if lockMode == "" {
		lockMode = "never"
	}
	locked = lockMode == "always" || (lockMode == "after_create" && mode == "edit")
	editable = !isReadOnly && !locked && strings.EqualFold(strings.TrimSpace(editableBy), "ALL")
	return editable, locked
}

// --- Form layout ------------------------------------------------------------

func (s *Service) GetFormLayout(ctx context.Context, orgID, moduleID, mode string) (*dto.FormLayoutResponse, error) {
	mode, err := normalizeFormMode(mode)
	if err != nil {
		return nil, err
	}
	layout, err := s.ensureFormLayout(ctx, orgID, moduleID)
	if err != nil {
		return nil, err
	}
	var cfg formLayoutConfig
	if err := json.Unmarshal(layout.Config, &cfg); err != nil {
		return nil, err
	}
	fields, err := s.repo.ListFieldsForHydrate(ctx, orgID, moduleID)
	if err != nil {
		return nil, err
	}

	// Ensure every visible field appears in some section (new fields / drift).
	cfg, changed := ensureFormOrphans(cfg, fields)
	if changed {
		raw, err := json.Marshal(cfg)
		if err != nil {
			return nil, err
		}
		layout, err = s.repo.UpsertDefaultFormLayout(ctx, orgID, moduleID, raw)
		if err != nil {
			return nil, err
		}
	}

	byAPI := make(map[string]repository.HydrateField, len(fields))
	for _, f := range fields {
		byAPI[f.APIName] = f
	}
	sections := make([]dto.FormLayoutSection, 0, len(cfg.Sections))
	for i, sec := range cfg.Sections {
		order := sec.Order
		if order == 0 {
			order = i + 1
		}
		cols := sec.Columns
		if cols < 1 || cols > 3 {
			cols = 2
		}
		hydrated := make([]dto.FormLayoutField, 0, len(sec.Fields))
		display := 0
		for _, key := range sec.Fields {
			f, ok := byAPI[key]
			if !ok || !f.IsVisible {
				continue
			}
			display++
			editable, locked := ResolveEditable(f.IsReadOnly, f.LockMode, mode, f.EditableBy)
			hydrated = append(hydrated, dto.FormLayoutField{
				ID:           f.ID,
				Key:          f.APIName,
				Label:        f.Label,
				Type:         f.FieldType,
				Required:     f.IsRequired,
				Editable:     editable,
				Locked:       locked,
				DisplayOrder: display,
				Placeholder:  f.Placeholder,
				Description:  f.Description,
				DefaultValue: f.DefaultValue,
				ValidationRules: dto.FormFieldValidationRules{
					Min: f.MinLength, Max: f.MaxLength, Regex: f.Regex,
				},
				Options:        parseFormOptions(f.Options),
				LookupModuleID: f.LookupModuleID,
				LockMode:       f.LockMode,
			})
		}
		sections = append(sections, dto.FormLayoutSection{
			ID:          sec.Key,
			Title:       sec.Label,
			Description: sec.Description,
			Order:       order,
			Collapsed:   sec.Collapsed,
			Columns:     cols,
			Fields:      hydrated,
		})
	}
	sort.SliceStable(sections, func(i, j int) bool { return sections[i].Order < sections[j].Order })
	return &dto.FormLayoutResponse{
		ID: layout.ID, Name: layout.Name, LayoutType: layout.Type,
		IsDefault: layout.IsDefault, Mode: mode, Sections: sections,
	}, nil
}

// ensureFormOrphans appends visible fields missing from the form layout into the
// first section (or creates general). Returns whether the config changed.
func ensureFormOrphans(cfg formLayoutConfig, fields []repository.HydrateField) (formLayoutConfig, bool) {
	seen := map[string]bool{}
	for _, sec := range cfg.Sections {
		for _, name := range sec.Fields {
			seen[name] = true
		}
	}
	orphans := make([]string, 0)
	for _, f := range fields {
		if !f.IsVisible || f.IsSystem {
			continue
		}
		if !seen[f.APIName] {
			orphans = append(orphans, f.APIName)
		}
	}
	if len(orphans) == 0 {
		return cfg, false
	}
	if len(cfg.Sections) == 0 {
		cfg.Sections = []dto.LayoutSection{{
			Key: "general", Label: "General Information", Order: 1, Columns: 2, Fields: orphans,
		}}
		return cfg, true
	}
	target := 0
	for i, sec := range cfg.Sections {
		if sec.Key != "system" {
			target = i
			break
		}
	}
	cfg.Sections[target].Fields = append(cfg.Sections[target].Fields, orphans...)
	return cfg, true
}

func (s *Service) UpdateFormLayout(ctx context.Context, orgID, moduleID string, req dto.UpdateFormLayoutRequest) (*dto.FormLayoutResponse, error) {
	cfg, err := s.normalizeFormConfig(ctx, orgID, moduleID, req.Sections)
	if err != nil {
		return nil, err
	}
	raw, err := json.Marshal(cfg)
	if err != nil {
		return nil, err
	}
	if _, err := s.repo.UpsertDefaultFormLayout(ctx, orgID, moduleID, raw); err != nil {
		return nil, err
	}
	if err := s.syncFieldSortOrder(ctx, orgID, moduleID, cfg.Sections); err != nil {
		return nil, err
	}
	return s.GetFormLayout(ctx, orgID, moduleID, "create")
}

func (s *Service) ReorderFormFields(ctx context.Context, orgID, moduleID string, req dto.FormReorderRequest) (*dto.FormLayoutResponse, error) {
	layout, err := s.ensureFormLayout(ctx, orgID, moduleID)
	if err != nil {
		return nil, err
	}
	var cfg formLayoutConfig
	if err := json.Unmarshal(layout.Config, &cfg); err != nil {
		return nil, err
	}
	fields, err := s.repo.ListFieldsForHydrate(ctx, orgID, moduleID)
	if err != nil {
		return nil, err
	}
	idToAPI := make(map[string]string, len(fields))
	for _, f := range fields {
		idToAPI[f.ID] = f.APIName
		idToAPI[f.APIName] = f.APIName
	}

	sectionIdx := -1
	for i := range cfg.Sections {
		if cfg.Sections[i].Key == req.SectionID {
			sectionIdx = i
			break
		}
	}
	if sectionIdx < 0 {
		return nil, ErrSectionNotFound
	}

	type ordered struct {
		key   string
		order int
	}
	items := make([]ordered, 0, len(req.Fields))
	seen := map[string]bool{}
	for _, it := range req.Fields {
		api := idToAPI[strings.TrimSpace(it.FieldID)]
		if api == "" {
			continue
		}
		if seen[api] {
			continue
		}
		seen[api] = true
		items = append(items, ordered{key: api, order: it.Order})
	}
	sort.SliceStable(items, func(i, j int) bool { return items[i].order < items[j].order })

	// Keep only fields that already belong to this section, in new order; append leftovers.
	inSection := map[string]bool{}
	for _, name := range cfg.Sections[sectionIdx].Fields {
		inSection[name] = true
	}
	newFields := make([]string, 0, len(cfg.Sections[sectionIdx].Fields))
	placed := map[string]bool{}
	for _, it := range items {
		if inSection[it.key] {
			newFields = append(newFields, it.key)
			placed[it.key] = true
		}
	}
	for _, name := range cfg.Sections[sectionIdx].Fields {
		if !placed[name] {
			newFields = append(newFields, name)
		}
	}
	cfg.Sections[sectionIdx].Fields = newFields

	raw, err := json.Marshal(cfg)
	if err != nil {
		return nil, err
	}
	if _, err := s.repo.UpsertDefaultFormLayout(ctx, orgID, moduleID, raw); err != nil {
		return nil, err
	}
	if err := s.syncFieldSortOrder(ctx, orgID, moduleID, cfg.Sections); err != nil {
		return nil, err
	}
	return s.GetFormLayout(ctx, orgID, moduleID, "create")
}

func (s *Service) CreateFormSection(ctx context.Context, orgID, moduleID string, req dto.CreateSectionRequest) (*dto.FormLayoutResponse, error) {
	layout, err := s.ensureFormLayout(ctx, orgID, moduleID)
	if err != nil {
		return nil, err
	}
	var cfg formLayoutConfig
	if err := json.Unmarshal(layout.Config, &cfg); err != nil {
		return nil, err
	}
	key := strings.TrimSpace(req.Key)
	if key == "" || key == "system" {
		return nil, fmt.Errorf("invalid section key")
	}
	for _, sec := range cfg.Sections {
		if sec.Key == key {
			return nil, ErrSectionExists
		}
	}
	cols := req.Columns
	if cols < 1 || cols > 3 {
		cols = 2
	}
	cfg.Sections = append(cfg.Sections, dto.LayoutSection{
		Key: key, Label: strings.TrimSpace(req.Label), Description: req.Description,
		Order: len(cfg.Sections) + 1, Collapsed: req.Collapsed, Columns: cols, Fields: []string{},
	})
	raw, err := json.Marshal(cfg)
	if err != nil {
		return nil, err
	}
	if _, err := s.repo.UpsertDefaultFormLayout(ctx, orgID, moduleID, raw); err != nil {
		return nil, err
	}
	return s.GetFormLayout(ctx, orgID, moduleID, "create")
}

func (s *Service) UpdateFormSection(ctx context.Context, orgID, moduleID, sectionID string, req dto.UpdateSectionRequest) (*dto.FormLayoutResponse, error) {
	layout, err := s.ensureFormLayout(ctx, orgID, moduleID)
	if err != nil {
		return nil, err
	}
	var cfg formLayoutConfig
	if err := json.Unmarshal(layout.Config, &cfg); err != nil {
		return nil, err
	}
	found := false
	for i := range cfg.Sections {
		if cfg.Sections[i].Key != sectionID {
			continue
		}
		found = true
		if req.Label != nil {
			cfg.Sections[i].Label = strings.TrimSpace(*req.Label)
		}
		if req.Description != nil {
			cfg.Sections[i].Description = *req.Description
		}
		if req.Columns != nil {
			cols := *req.Columns
			if cols < 1 || cols > 3 {
				return nil, fmt.Errorf("columns must be 1, 2, or 3")
			}
			cfg.Sections[i].Columns = cols
		}
		if req.Collapsed != nil {
			cfg.Sections[i].Collapsed = *req.Collapsed
		}
		if req.Order != nil {
			cfg.Sections[i].Order = *req.Order
		}
		break
	}
	if !found {
		return nil, ErrSectionNotFound
	}
	raw, err := json.Marshal(cfg)
	if err != nil {
		return nil, err
	}
	if _, err := s.repo.UpsertDefaultFormLayout(ctx, orgID, moduleID, raw); err != nil {
		return nil, err
	}
	return s.GetFormLayout(ctx, orgID, moduleID, "create")
}

func (s *Service) DeleteFormSection(ctx context.Context, orgID, moduleID, sectionID string) (*dto.FormLayoutResponse, error) {
	layout, err := s.ensureFormLayout(ctx, orgID, moduleID)
	if err != nil {
		return nil, err
	}
	var cfg formLayoutConfig
	if err := json.Unmarshal(layout.Config, &cfg); err != nil {
		return nil, err
	}
	idx := -1
	for i := range cfg.Sections {
		if cfg.Sections[i].Key == sectionID {
			idx = i
			break
		}
	}
	if idx < 0 {
		return nil, ErrSectionNotFound
	}
	if len(cfg.Sections[idx].Fields) > 0 {
		return nil, ErrSectionNotEmpty
	}
	if len(cfg.Sections) <= 1 {
		return nil, fmt.Errorf("at least one section is required")
	}
	cfg.Sections = append(cfg.Sections[:idx], cfg.Sections[idx+1:]...)
	raw, err := json.Marshal(cfg)
	if err != nil {
		return nil, err
	}
	if _, err := s.repo.UpsertDefaultFormLayout(ctx, orgID, moduleID, raw); err != nil {
		return nil, err
	}
	return s.GetFormLayout(ctx, orgID, moduleID, "create")
}

func (s *Service) ensureFormLayout(ctx context.Context, orgID, moduleID string) (*repository.Layout, error) {
	layout, err := s.repo.GetDefaultLayout(ctx, orgID, moduleID, repository.LayoutTypeForm)
	if err != nil {
		return nil, err
	}
	if layout != nil {
		return layout, nil
	}
	cfg, err := s.buildDefaultFormConfig(ctx, orgID, moduleID)
	if err != nil {
		return nil, err
	}
	return s.repo.UpsertDefaultFormLayout(ctx, orgID, moduleID, cfg)
}

func (s *Service) buildDefaultFormConfig(ctx context.Context, orgID, moduleID string) (json.RawMessage, error) {
	fields, err := s.repo.ListVisibleFieldAPINames(ctx, orgID, moduleID)
	if err != nil {
		return nil, err
	}
	cfg := formLayoutConfig{
		Sections: []dto.LayoutSection{{
			Key: "general", Label: "General Information", Order: 1, Columns: 2, Fields: fields,
		}},
	}
	return json.Marshal(cfg)
}

func (s *Service) normalizeFormConfig(ctx context.Context, orgID, moduleID string, sections []dto.LayoutSection) (*formLayoutConfig, error) {
	known, err := s.repo.ListNonSystemFields(ctx, orgID, moduleID)
	if err != nil {
		return nil, err
	}
	knownSet := make(map[string]bool, len(known))
	for _, f := range known {
		knownSet[f.APIName] = true
	}

	seenKeys := make(map[string]bool)
	seenFields := make(map[string]bool)
	out := make([]dto.LayoutSection, 0, len(sections))

	for i, sec := range sections {
		key := strings.TrimSpace(sec.Key)
		label := strings.TrimSpace(sec.Label)
		if key == "" || label == "" {
			return nil, fmt.Errorf("section key and label are required")
		}
		if key == "system" {
			continue // form layouts never include the detail system section
		}
		if seenKeys[key] {
			return nil, fmt.Errorf("duplicate section key %q", key)
		}
		seenKeys[key] = true

		fields := make([]string, 0, len(sec.Fields))
		for _, name := range sec.Fields {
			name = strings.TrimSpace(name)
			if name == "" || systemFieldNames[name] || !knownSet[name] || seenFields[name] {
				continue
			}
			fields = append(fields, name)
			seenFields[name] = true
		}
		order := sec.Order
		if order == 0 {
			order = i + 1
		}
		cols := sec.Columns
		if cols < 1 || cols > 3 {
			cols = 2
		}
		out = append(out, dto.LayoutSection{
			Key: key, Label: label, Description: sec.Description,
			Order: order, Collapsed: sec.Collapsed, Columns: cols, Fields: fields,
		})
	}
	if len(out) == 0 {
		return nil, fmt.Errorf("at least one section is required")
	}
	for _, f := range known {
		if !seenFields[f.APIName] {
			out[0].Fields = append(out[0].Fields, f.APIName)
			seenFields[f.APIName] = true
		}
	}
	return &formLayoutConfig{Sections: out}, nil
}

func normalizeFormMode(mode string) (string, error) {
	m := strings.ToLower(strings.TrimSpace(mode))
	if m == "" {
		return "create", nil
	}
	if m != "create" && m != "edit" {
		return "", ErrInvalidMode
	}
	return m, nil
}

func parseFormOptions(raw []byte) []dto.FormFieldOption {
	out := make([]dto.FormFieldOption, 0)
	if len(raw) == 0 {
		return out
	}
	var strs []string
	if err := json.Unmarshal(raw, &strs); err == nil {
		for _, s := range strs {
			out = append(out, dto.FormFieldOption{Label: s, Value: s})
		}
		return out
	}
	var objs []dto.FormFieldOption
	if err := json.Unmarshal(raw, &objs); err == nil {
		return objs
	}
	return out
}

// --- List layout ------------------------------------------------------------

func (s *Service) GetListLayout(ctx context.Context, orgID, moduleID string, includeHidden bool) (*dto.ListLayoutResponse, error) {
	layout, fields, err := s.ensureListLayoutReconciled(ctx, orgID, moduleID)
	if err != nil {
		return nil, err
	}
	var cfg listLayoutConfig
	if err := json.Unmarshal(layout.Config, &cfg); err != nil {
		return nil, err
	}
	known := hydrateFieldMap(fields)
	cols := hydrateListColumns(normalizeListColumns(cfg.Columns), known)
	if !includeHidden {
		visible := make([]dto.ListColumn, 0, len(cols))
		for _, c := range cols {
			if c.Visible {
				visible = append(visible, c)
			}
		}
		cols = visible
	}
	sort.SliceStable(cols, func(i, j int) bool { return cols[i].Order < cols[j].Order })
	return &dto.ListLayoutResponse{
		ID: layout.ID, Name: layout.Name, LayoutType: layout.Type,
		IsDefault: layout.IsDefault, Columns: cols,
	}, nil
}

func (s *Service) UpdateListLayout(ctx context.Context, orgID, moduleID string, req dto.UpdateListLayoutRequest) (*dto.ListLayoutResponse, error) {
	fields, err := s.repo.ListFieldsForHydrate(ctx, orgID, moduleID)
	if err != nil {
		return nil, err
	}
	known := hydrateFieldMap(fields)
	cols, err := normalizeListRequest(req.Columns, known)
	if err != nil {
		return nil, err
	}
	raw, err := json.Marshal(listLayoutConfig{Columns: stripHydratedListColumns(cols)})
	if err != nil {
		return nil, err
	}
	if _, err := s.repo.UpsertDefaultListLayout(ctx, orgID, moduleID, raw); err != nil {
		return nil, err
	}
	return s.GetListLayout(ctx, orgID, moduleID, true)
}

func (s *Service) ReorderListColumns(ctx context.Context, orgID, moduleID string, req dto.ListReorderRequest) (*dto.ListLayoutResponse, error) {
	layout, _, err := s.ensureListLayoutReconciled(ctx, orgID, moduleID)
	if err != nil {
		return nil, err
	}
	var cfg listLayoutConfig
	if err := json.Unmarshal(layout.Config, &cfg); err != nil {
		return nil, err
	}
	cols := normalizeListColumns(cfg.Columns)
	byKey := make(map[string]*dto.ListColumn, len(cols))
	for i := range cols {
		byKey[cols[i].FieldKey] = &cols[i]
	}
	for _, it := range req.Columns {
		if c, ok := byKey[it.FieldKey]; ok && c.FieldKey != dto.ActionsColumnKey {
			c.Order = it.Order
		}
	}
	forceActionsLast(cols)
	sort.SliceStable(cols, func(i, j int) bool { return cols[i].Order < cols[j].Order })
	for i := range cols {
		cols[i].Order = i + 1
	}
	raw, err := json.Marshal(listLayoutConfig{Columns: stripHydratedListColumns(cols)})
	if err != nil {
		return nil, err
	}
	if _, err := s.repo.UpsertDefaultListLayout(ctx, orgID, moduleID, raw); err != nil {
		return nil, err
	}
	return s.GetListLayout(ctx, orgID, moduleID, true)
}

func (s *Service) ToggleListColumn(ctx context.Context, orgID, moduleID string, req dto.ListToggleRequest) (*dto.ListLayoutResponse, error) {
	if req.Visible == nil {
		return nil, fmt.Errorf("visible is required")
	}
	layout, fields, err := s.ensureListLayoutReconciled(ctx, orgID, moduleID)
	if err != nil {
		return nil, err
	}
	known := hydrateFieldMap(fields)
	var cfg listLayoutConfig
	if err := json.Unmarshal(layout.Config, &cfg); err != nil {
		return nil, err
	}
	cols := normalizeListColumns(cfg.Columns)
	found := false
	for i := range cols {
		if cols[i].FieldKey != req.FieldKey {
			continue
		}
		found = true
		isSystem := cols[i].System
		if f, ok := known[cols[i].FieldKey]; ok {
			isSystem = f.IsSystem
		}
		if cols[i].FieldKey == dto.ActionsColumnKey {
			return nil, ErrSystemColumn
		}
		if !*req.Visible && IsListColumnLocked(cols[i].FieldKey, isSystem) {
			return nil, ErrLockedColumn
		}
		cols[i].Visible = *req.Visible
		break
	}
	if !found {
		return nil, ErrNotFound
	}
	raw, err := json.Marshal(listLayoutConfig{Columns: stripHydratedListColumns(cols)})
	if err != nil {
		return nil, err
	}
	if _, err := s.repo.UpsertDefaultListLayout(ctx, orgID, moduleID, raw); err != nil {
		return nil, err
	}
	return s.GetListLayout(ctx, orgID, moduleID, true)
}

// ResetListLayout rebuilds the default list layout (system visible, custom hidden).
func (s *Service) ResetListLayout(ctx context.Context, orgID, moduleID string) (*dto.ListLayoutResponse, error) {
	cfg, err := s.buildDefaultListConfig(ctx, orgID, moduleID)
	if err != nil {
		return nil, err
	}
	if _, err := s.repo.UpsertDefaultListLayout(ctx, orgID, moduleID, cfg); err != nil {
		return nil, err
	}
	return s.GetListLayout(ctx, orgID, moduleID, true)
}

func (s *Service) ensureListLayout(ctx context.Context, orgID, moduleID string) (*repository.Layout, error) {
	layout, _, err := s.ensureListLayoutReconciled(ctx, orgID, moduleID)
	return layout, err
}

func (s *Service) ensureListLayoutReconciled(ctx context.Context, orgID, moduleID string) (*repository.Layout, []repository.HydrateField, error) {
	fields, err := s.repo.ListFieldsForHydrate(ctx, orgID, moduleID)
	if err != nil {
		return nil, nil, err
	}
	layout, err := s.repo.GetDefaultLayout(ctx, orgID, moduleID, repository.LayoutTypeList)
	if err != nil {
		return nil, nil, err
	}
	if layout == nil {
		cfg, err := buildDefaultListColumns(fields)
		if err != nil {
			return nil, nil, err
		}
		raw, err := json.Marshal(listLayoutConfig{Columns: stripHydratedListColumns(cfg)})
		if err != nil {
			return nil, nil, err
		}
		layout, err = s.repo.UpsertDefaultListLayout(ctx, orgID, moduleID, raw)
		if err != nil {
			return nil, nil, err
		}
		return layout, fields, nil
	}

	var cfg listLayoutConfig
	if err := json.Unmarshal(layout.Config, &cfg); err != nil {
		return nil, nil, err
	}
	cols, changed := reconcileListColumns(cfg.Columns, fields)
	if changed {
		raw, err := json.Marshal(listLayoutConfig{Columns: stripHydratedListColumns(cols)})
		if err != nil {
			return nil, nil, err
		}
		layout, err = s.repo.UpsertDefaultListLayout(ctx, orgID, moduleID, raw)
		if err != nil {
			return nil, nil, err
		}
	}
	return layout, fields, nil
}

func (s *Service) buildDefaultListConfig(ctx context.Context, orgID, moduleID string) (json.RawMessage, error) {
	fields, err := s.repo.ListFieldsForHydrate(ctx, orgID, moduleID)
	if err != nil {
		return nil, err
	}
	cols, err := buildDefaultListColumns(fields)
	if err != nil {
		return nil, err
	}
	return json.Marshal(listLayoutConfig{Columns: stripHydratedListColumns(cols)})
}

func buildDefaultListColumns(fields []repository.HydrateField) ([]dto.ListColumn, error) {
	cols := make([]dto.ListColumn, 0, len(fields)+1)
	order := 0
	for _, f := range fields {
		if !f.IsVisible {
			continue
		}
		order++
		cols = append(cols, dto.ListColumn{
			FieldKey:   f.APIName,
			Visible:    defaultListVisible(f.IsSystem, f.APIName),
			Order:      order,
			Sortable:   true,
			Searchable: f.IsSearchable,
			System:     false,
			Locked:     IsListColumnLocked(f.APIName, f.IsSystem),
		})
	}
	cols = append(cols, actionsListColumn(order+1))
	return cols, nil
}

// reconcileListColumns adds missing fields, drops deleted ones, renumbers.
func reconcileListColumns(existing []dto.ListColumn, fields []repository.HydrateField) ([]dto.ListColumn, bool) {
	known := hydrateFieldMap(fields)
	cols := normalizeListColumns(existing)
	seen := map[string]bool{}
	out := make([]dto.ListColumn, 0, len(cols)+len(fields))
	changed := false

	for _, c := range cols {
		if c.FieldKey == dto.ActionsColumnKey {
			continue
		}
		f, ok := known[c.FieldKey]
		if !ok || !f.IsVisible {
			changed = true
			continue
		}
		seen[c.FieldKey] = true
		out = append(out, c)
	}

	for _, f := range fields {
		if !f.IsVisible || seen[f.APIName] {
			continue
		}
		changed = true
		out = append(out, dto.ListColumn{
			FieldKey:   f.APIName,
			Visible:    defaultListVisible(f.IsSystem, f.APIName),
			Order:      len(out) + 1,
			Sortable:   true,
			Searchable: f.IsSearchable,
			System:     false,
		})
	}

	out = append(out, actionsListColumn(len(out)+1))
	sort.SliceStable(out, func(i, j int) bool {
		if out[i].FieldKey == dto.ActionsColumnKey {
			return false
		}
		if out[j].FieldKey == dto.ActionsColumnKey {
			return true
		}
		return out[i].Order < out[j].Order
	})
	for i := range out {
		want := i + 1
		if out[i].Order != want {
			changed = true
			out[i].Order = want
		}
		if out[i].FieldKey == dto.ActionsColumnKey {
			out[i].System = true
			out[i].Visible = true
			out[i].Locked = true
		}
	}
	return out, changed
}

func normalizeListRequest(req []dto.ListColumn, known map[string]repository.HydrateField) ([]dto.ListColumn, error) {
	if len(req) == 0 {
		return nil, fmt.Errorf("%w: columns required", ErrInvalidListCol)
	}
	out := make([]dto.ListColumn, 0, len(req)+1)
	seenKey := map[string]bool{}
	seenOrder := map[int]bool{}

	for i, c := range req {
		key := strings.TrimSpace(c.FieldKey)
		if key == "" {
			return nil, fmt.Errorf("%w: empty field_key", ErrInvalidListCol)
		}
		if key == dto.ActionsColumnKey {
			continue
		}
		f, ok := known[key]
		if !ok || !f.IsVisible {
			return nil, fmt.Errorf("%w: unknown or archived field %q", ErrInvalidListCol, key)
		}
		if seenKey[key] {
			return nil, fmt.Errorf("%w: duplicate field_key %q", ErrInvalidListCol, key)
		}
		seenKey[key] = true

		order := c.Order
		if order == 0 {
			order = i + 1
		}
		if seenOrder[order] {
			return nil, fmt.Errorf("%w: duplicate display_order %d", ErrInvalidListCol, order)
		}
		seenOrder[order] = true

		locked := IsListColumnLocked(key, f.IsSystem)
		visible := c.Visible
		if locked {
			visible = true
		}

		out = append(out, dto.ListColumn{
			FieldKey:   key,
			Visible:    visible,
			Order:      order,
			Sortable:   c.Sortable,
			Searchable: c.Searchable,
			System:     false,
			Locked:     locked,
			Width:      c.Width,
		})
	}

	sort.SliceStable(out, func(i, j int) bool { return out[i].Order < out[j].Order })
	for i := range out {
		out[i].Order = i + 1
	}
	out = append(out, actionsListColumn(len(out)+1))
	return out, nil
}

func normalizeListColumns(cols []dto.ListColumn) []dto.ListColumn {
	out := make([]dto.ListColumn, 0, len(cols))
	hasActions := false
	for _, c := range cols {
		if c.FieldKey == dto.ActionsColumnKey {
			hasActions = true
			c.System = true
			c.Visible = true
			c.Locked = true
		}
		out = append(out, c)
	}
	if !hasActions {
		out = append(out, actionsListColumn(len(out)+1))
	}
	return out
}

func hydrateFieldMap(fields []repository.HydrateField) map[string]repository.HydrateField {
	known := make(map[string]repository.HydrateField, len(fields))
	for _, f := range fields {
		known[f.APIName] = f
	}
	return known
}

func hydrateListColumns(cols []dto.ListColumn, known map[string]repository.HydrateField) []dto.ListColumn {
	out := make([]dto.ListColumn, len(cols))
	for i, c := range cols {
		out[i] = c
		if c.FieldKey == dto.ActionsColumnKey {
			out[i].Label = "Actions"
			out[i].System = true
			out[i].Visible = true
			out[i].Locked = true
			continue
		}
		if f, ok := known[c.FieldKey]; ok {
			out[i].FieldID = f.ID
			out[i].Label = f.Label
			out[i].Locked = IsListColumnLocked(c.FieldKey, f.IsSystem)
			if out[i].Locked {
				out[i].Visible = true
			}
		}
	}
	return out
}

// stripHydratedListColumns removes computed fields before persisting JSON.
func stripHydratedListColumns(cols []dto.ListColumn) []dto.ListColumn {
	out := make([]dto.ListColumn, len(cols))
	for i, c := range cols {
		out[i] = dto.ListColumn{
			FieldKey:   c.FieldKey,
			Visible:    c.Visible,
			Order:      c.Order,
			Sortable:   c.Sortable,
			Searchable: c.Searchable,
			System:     c.System || c.FieldKey == dto.ActionsColumnKey,
			Width:      c.Width,
		}
		if out[i].FieldKey == dto.ActionsColumnKey {
			out[i].Visible = true
			out[i].System = true
		}
	}
	return out
}

func actionsListColumn(order int) dto.ListColumn {
	return dto.ListColumn{
		FieldKey: dto.ActionsColumnKey, Visible: true, Order: order,
		Sortable: false, Searchable: false, System: true, Locked: true,
	}
}

func forceActionsLast(cols []dto.ListColumn) {
	maxOrder := 0
	for i := range cols {
		if cols[i].FieldKey != dto.ActionsColumnKey && cols[i].Order > maxOrder {
			maxOrder = cols[i].Order
		}
	}
	for i := range cols {
		if cols[i].FieldKey == dto.ActionsColumnKey {
			cols[i].Order = maxOrder + 1
			cols[i].System = true
			cols[i].Visible = true
			cols[i].Locked = true
		}
	}
}
