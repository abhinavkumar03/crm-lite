package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/abhinavkumar03/crm-lite/backend/internal/workspace/dto"
	"github.com/abhinavkumar03/crm-lite/backend/internal/workspace/repository"
)

var (
	ErrNotFound = errors.New("not found")
	ErrForbidden = errors.New("forbidden")
)

type Service struct {
	repo *repository.Repository
}

func New(repo *repository.Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) ensureRecord(ctx context.Context, orgID, moduleID, recordID string) error {
	ok, err := s.repo.RecordExists(ctx, orgID, moduleID, recordID)
	if err != nil {
		return err
	}
	if !ok {
		return ErrNotFound
	}
	return nil
}

func (s *Service) GetDetailLayout(ctx context.Context, orgID, moduleID string) (*dto.LayoutResponse, error) {
	l, err := s.repo.GetDefaultDetailLayout(ctx, orgID, moduleID)
	if err != nil {
		return nil, err
	}
	if l == nil {
		cfg, err := s.buildDefaultConfig(ctx, orgID, moduleID)
		if err != nil {
			return nil, err
		}
		l, err = s.repo.UpsertDefaultDetailLayout(ctx, orgID, moduleID, cfg)
		if err != nil {
			return nil, err
		}
	}
	return &dto.LayoutResponse{
		ID: l.ID, Name: l.Name, Type: l.Type, IsDefault: l.IsDefault, Config: l.Config,
	}, nil
}

func (s *Service) UpdateDetailLayout(ctx context.Context, orgID, moduleID string, req dto.UpdateDetailLayoutRequest) (*dto.LayoutResponse, error) {
	cfg, err := s.normalizeLayoutConfig(ctx, orgID, moduleID, req)
	if err != nil {
		return nil, err
	}
	raw, err := json.Marshal(cfg)
	if err != nil {
		return nil, err
	}
	l, err := s.repo.UpsertDefaultDetailLayout(ctx, orgID, moduleID, raw)
	if err != nil {
		return nil, err
	}
	if err := s.syncFieldSortOrder(ctx, orgID, moduleID, cfg.Sections); err != nil {
		return nil, err
	}
	return &dto.LayoutResponse{
		ID: l.ID, Name: l.Name, Type: l.Type, IsDefault: l.IsDefault, Config: l.Config,
	}, nil
}

// AppendFieldToSection adds api_name to the given section on detail + form layouts,
// and appends a visible column to the list layout.
func (s *Service) AppendFieldToSection(ctx context.Context, orgID, moduleID, sectionKey, apiName string) error {
	if err := s.appendFieldToDetail(ctx, orgID, moduleID, sectionKey, apiName); err != nil {
		return err
	}
	if err := s.appendFieldToForm(ctx, orgID, moduleID, sectionKey, apiName); err != nil {
		return err
	}
	return s.appendFieldToList(ctx, orgID, moduleID, apiName)
}

func (s *Service) appendFieldToDetail(ctx context.Context, orgID, moduleID, sectionKey, apiName string) error {
	layout, err := s.GetDetailLayout(ctx, orgID, moduleID)
	if err != nil {
		return err
	}
	var cfg layoutConfig
	if err := json.Unmarshal(layout.Config, &cfg); err != nil {
		return err
	}
	if cfg.Sections == nil {
		cfg.Sections = []dto.LayoutSection{}
	}
	key := strings.TrimSpace(sectionKey)
	if key == "" || key == "system" {
		key = "general"
	}
	for i := range cfg.Sections {
		cfg.Sections[i].Fields = filterOut(cfg.Sections[i].Fields, apiName)
	}
	found := false
	for i := range cfg.Sections {
		if cfg.Sections[i].Key == key {
			cfg.Sections[i].Fields = append(cfg.Sections[i].Fields, apiName)
			found = true
			break
		}
	}
	if !found {
		label := strings.ReplaceAll(key, "_", " ")
		if key == "general" {
			label = "General Information"
		}
		cfg.Sections = append([]dto.LayoutSection{{
			Key: key, Label: label, Fields: []string{apiName},
		}}, cfg.Sections...)
	}
	if len(cfg.Tabs) == 0 {
		cfg.Tabs = []string{"overview", "notes", "attachments", "timeline", "related"}
	}
	raw, err := json.Marshal(cfg)
	if err != nil {
		return err
	}
	_, err = s.repo.UpsertDefaultDetailLayout(ctx, orgID, moduleID, raw)
	if err != nil {
		return err
	}
	return s.syncFieldSortOrder(ctx, orgID, moduleID, cfg.Sections)
}

func (s *Service) appendFieldToForm(ctx context.Context, orgID, moduleID, sectionKey, apiName string) error {
	layout, err := s.ensureFormLayout(ctx, orgID, moduleID)
	if err != nil {
		return err
	}
	var cfg formLayoutConfig
	if err := json.Unmarshal(layout.Config, &cfg); err != nil {
		return err
	}
	key := strings.TrimSpace(sectionKey)
	if key == "" || key == "system" {
		key = "general"
	}
	for i := range cfg.Sections {
		cfg.Sections[i].Fields = filterOut(cfg.Sections[i].Fields, apiName)
	}
	found := false
	for i := range cfg.Sections {
		if cfg.Sections[i].Key == key {
			cfg.Sections[i].Fields = append(cfg.Sections[i].Fields, apiName)
			found = true
			break
		}
	}
	if !found {
		label := "General Information"
		if key != "general" {
			label = strings.ReplaceAll(key, "_", " ")
		}
		cfg.Sections = append(cfg.Sections, dto.LayoutSection{
			Key: key, Label: label, Order: len(cfg.Sections) + 1, Columns: 2, Fields: []string{apiName},
		})
	}
	raw, err := json.Marshal(cfg)
	if err != nil {
		return err
	}
	_, err = s.repo.UpsertDefaultFormLayout(ctx, orgID, moduleID, raw)
	return err
}

func (s *Service) appendFieldToList(ctx context.Context, orgID, moduleID, apiName string) error {
	layout, fields, err := s.ensureListLayoutReconciled(ctx, orgID, moduleID)
	if err != nil {
		return err
	}
	var cfg listLayoutConfig
	if err := json.Unmarshal(layout.Config, &cfg); err != nil {
		return err
	}
	cols := normalizeListColumns(cfg.Columns)
	for _, c := range cols {
		if c.FieldKey == apiName {
			return nil
		}
	}
	isSystem := false
	searchable := false
	for _, f := range fields {
		if f.APIName == apiName {
			isSystem = f.IsSystem
			searchable = f.IsSearchable
			break
		}
	}
	// Insert before _actions.
	insertAt := len(cols)
	for i, c := range cols {
		if c.FieldKey == dto.ActionsColumnKey {
			insertAt = i
			break
		}
	}
	newCol := dto.ListColumn{
		FieldKey: apiName, Visible: defaultListVisible(isSystem, apiName), Order: insertAt + 1,
		Sortable: true, Searchable: searchable, System: false,
		Locked: IsListColumnLocked(apiName, isSystem),
	}
	cols = append(cols[:insertAt], append([]dto.ListColumn{newCol}, cols[insertAt:]...)...)
	for i := range cols {
		cols[i].Order = i + 1
		if cols[i].FieldKey == dto.ActionsColumnKey {
			cols[i].System = true
			cols[i].Visible = true
			cols[i].Locked = true
		}
	}
	raw, err := json.Marshal(listLayoutConfig{Columns: stripHydratedListColumns(cols)})
	if err != nil {
		return err
	}
	_, err = s.repo.UpsertDefaultListLayout(ctx, orgID, moduleID, raw)
	return err
}

// RemoveFieldFromLayout drops api_name from detail, form, and list layouts.
func (s *Service) RemoveFieldFromLayout(ctx context.Context, orgID, moduleID, apiName string) error {
	if err := s.removeFieldFromDetail(ctx, orgID, moduleID, apiName); err != nil {
		return err
	}
	if err := s.removeFieldFromForm(ctx, orgID, moduleID, apiName); err != nil {
		return err
	}
	return s.removeFieldFromList(ctx, orgID, moduleID, apiName)
}

func (s *Service) removeFieldFromDetail(ctx context.Context, orgID, moduleID, apiName string) error {
	layout, err := s.repo.GetDefaultDetailLayout(ctx, orgID, moduleID)
	if err != nil || layout == nil {
		return err
	}
	var cfg layoutConfig
	if err := json.Unmarshal(layout.Config, &cfg); err != nil {
		return err
	}
	changed := false
	for i := range cfg.Sections {
		before := len(cfg.Sections[i].Fields)
		cfg.Sections[i].Fields = filterOut(cfg.Sections[i].Fields, apiName)
		if len(cfg.Sections[i].Fields) != before {
			changed = true
		}
	}
	if !changed {
		return nil
	}
	raw, err := json.Marshal(cfg)
	if err != nil {
		return err
	}
	_, err = s.repo.UpsertDefaultDetailLayout(ctx, orgID, moduleID, raw)
	return err
}

func (s *Service) removeFieldFromForm(ctx context.Context, orgID, moduleID, apiName string) error {
	layout, err := s.repo.GetDefaultLayout(ctx, orgID, moduleID, repository.LayoutTypeForm)
	if err != nil || layout == nil {
		return err
	}
	var cfg formLayoutConfig
	if err := json.Unmarshal(layout.Config, &cfg); err != nil {
		return err
	}
	changed := false
	for i := range cfg.Sections {
		before := len(cfg.Sections[i].Fields)
		cfg.Sections[i].Fields = filterOut(cfg.Sections[i].Fields, apiName)
		if len(cfg.Sections[i].Fields) != before {
			changed = true
		}
	}
	if !changed {
		return nil
	}
	raw, err := json.Marshal(cfg)
	if err != nil {
		return err
	}
	_, err = s.repo.UpsertDefaultFormLayout(ctx, orgID, moduleID, raw)
	return err
}

func (s *Service) removeFieldFromList(ctx context.Context, orgID, moduleID, apiName string) error {
	layout, err := s.repo.GetDefaultLayout(ctx, orgID, moduleID, repository.LayoutTypeList)
	if err != nil || layout == nil {
		return err
	}
	var cfg listLayoutConfig
	if err := json.Unmarshal(layout.Config, &cfg); err != nil {
		return err
	}
	out := make([]dto.ListColumn, 0, len(cfg.Columns))
	changed := false
	for _, c := range cfg.Columns {
		if c.FieldKey == apiName {
			changed = true
			continue
		}
		out = append(out, c)
	}
	if !changed {
		return nil
	}
	out = normalizeListColumns(out)
	for i := range out {
		out[i].Order = i + 1
	}
	raw, err := json.Marshal(listLayoutConfig{Columns: stripHydratedListColumns(out)})
	if err != nil {
		return err
	}
	_, err = s.repo.UpsertDefaultListLayout(ctx, orgID, moduleID, raw)
	return err
}

type layoutConfig struct {
	Sections []dto.LayoutSection `json:"sections"`
	Tabs     []string            `json:"tabs"`
}

var systemFieldNames = map[string]bool{
	"owner_id": true, "assigned_to": true, "visibility": true,
	"created_at": true, "updated_at": true,
}

func (s *Service) normalizeLayoutConfig(ctx context.Context, orgID, moduleID string, req dto.UpdateDetailLayoutRequest) (*layoutConfig, error) {
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
	sections := make([]dto.LayoutSection, 0, len(req.Sections))

	for _, sec := range req.Sections {
		key := strings.TrimSpace(sec.Key)
		label := strings.TrimSpace(sec.Label)
		if key == "" || label == "" {
			return nil, fmt.Errorf("section key and label are required")
		}
		if seenKeys[key] {
			return nil, fmt.Errorf("duplicate section key %q", key)
		}
		seenKeys[key] = true

		fields := make([]string, 0, len(sec.Fields))
		for _, name := range sec.Fields {
			name = strings.TrimSpace(name)
			if name == "" {
				continue
			}
			if systemFieldNames[name] {
				if !seenFields[name] {
					fields = append(fields, name)
					seenFields[name] = true
				}
				continue
			}
			// Custom fields never belong in the system section.
			if key == "system" {
				continue
			}
			if !knownSet[name] {
				continue // prune unknown
			}
			if seenFields[name] {
				continue // only one section
			}
			fields = append(fields, name)
			seenFields[name] = true
		}
		sections = append(sections, dto.LayoutSection{Key: key, Label: label, Fields: fields})
	}

	if len(sections) == 0 {
		return nil, fmt.Errorf("at least one section is required")
	}

	// Append orphans to the first non-system section (or first section).
	orphanTarget := 0
	for i, sec := range sections {
		if sec.Key != "system" {
			orphanTarget = i
			break
		}
	}
	for _, f := range known {
		if !seenFields[f.APIName] {
			sections[orphanTarget].Fields = append(sections[orphanTarget].Fields, f.APIName)
			seenFields[f.APIName] = true
		}
	}

	tabs := req.Tabs
	if len(tabs) == 0 {
		tabs = []string{"overview", "notes", "attachments", "timeline", "related"}
	}
	return &layoutConfig{Sections: sections, Tabs: tabs}, nil
}

func (s *Service) syncFieldSortOrder(ctx context.Context, orgID, moduleID string, sections []dto.LayoutSection) error {
	refs, err := s.repo.ListNonSystemFields(ctx, orgID, moduleID)
	if err != nil {
		return err
	}
	byAPI := make(map[string]string, len(refs))
	for _, f := range refs {
		byAPI[f.APIName] = f.ID
	}
	positions := make([]repository.FieldSortPosition, 0, len(refs))
	order := 0
	for _, sec := range sections {
		for _, name := range sec.Fields {
			if id, ok := byAPI[name]; ok {
				positions = append(positions, repository.FieldSortPosition{ID: id, SortOrder: order})
				order++
				delete(byAPI, name)
			}
		}
	}
	// Any remaining (should be none after normalize) keep at the end.
	for _, id := range byAPI {
		positions = append(positions, repository.FieldSortPosition{ID: id, SortOrder: order})
		order++
	}
	if len(positions) == 0 {
		return nil
	}
	return s.repo.ReorderFields(ctx, orgID, moduleID, positions)
}

func filterOut(list []string, apiName string) []string {
	out := make([]string, 0, len(list))
	for _, v := range list {
		if v != apiName {
			out = append(out, v)
		}
	}
	return out
}

func (s *Service) buildDefaultConfig(ctx context.Context, orgID, moduleID string) (json.RawMessage, error) {
	fields, err := s.repo.ListVisibleFieldAPINames(ctx, orgID, moduleID)
	if err != nil {
		return nil, err
	}
	cfg := map[string]any{
		"sections": []map[string]any{
			{"key": "general", "label": "General Information", "fields": fields},
			{"key": "system", "label": "System Fields", "fields": []string{"owner_id", "assigned_to", "visibility", "created_at", "updated_at"}},
		},
		"tabs": []string{"overview", "notes", "attachments", "timeline", "related"},
	}
	return json.Marshal(cfg)
}

func (s *Service) ListNotes(ctx context.Context, orgID, moduleID, recordID string) ([]dto.NoteResponse, error) {
	if err := s.ensureRecord(ctx, orgID, moduleID, recordID); err != nil {
		return nil, err
	}
	notes, err := s.repo.ListNotes(ctx, orgID, moduleID, recordID)
	if err != nil {
		return nil, err
	}
	out := make([]dto.NoteResponse, 0, len(notes))
	for _, n := range notes {
		out = append(out, dto.NoteResponse{
			ID: n.ID, Title: n.Title, Body: n.Body, CreatedBy: n.CreatedBy,
			AuthorName: n.AuthorName, CreatedAt: n.CreatedAt, UpdatedAt: n.UpdatedAt,
		})
	}
	return out, nil
}

func (s *Service) CreateNote(ctx context.Context, orgID, moduleID, recordID, userID string, req dto.CreateNoteRequest) (*dto.NoteResponse, error) {
	if err := s.ensureRecord(ctx, orgID, moduleID, recordID); err != nil {
		return nil, err
	}
	body := strings.TrimSpace(req.Body)
	if body == "" {
		return nil, errors.New("body required")
	}
	n, err := s.repo.CreateNote(ctx, orgID, moduleID, recordID, userID, body, req.Title)
	if err != nil {
		return nil, err
	}
	_ = s.repo.CreateActivity(ctx, orgID, moduleID, recordID, userID, "NOTE_ADDED", "Note added", nil)
	return &dto.NoteResponse{
		ID: n.ID, Title: n.Title, Body: n.Body, CreatedBy: n.CreatedBy,
		AuthorName: n.AuthorName, CreatedAt: n.CreatedAt, UpdatedAt: n.UpdatedAt,
	}, nil
}

func (s *Service) DeleteNote(ctx context.Context, orgID, moduleID, recordID, noteID, userID string) error {
	if err := s.ensureRecord(ctx, orgID, moduleID, recordID); err != nil {
		return err
	}
	ok, err := s.repo.DeleteNote(ctx, orgID, noteID)
	if err != nil {
		return err
	}
	if !ok {
		return ErrNotFound
	}
	_ = s.repo.CreateActivity(ctx, orgID, moduleID, recordID, userID, "NOTE_DELETED", "Note deleted", nil)
	return nil
}

func (s *Service) ListAttachments(ctx context.Context, orgID, moduleID, recordID string) ([]dto.AttachmentResponse, error) {
	if err := s.ensureRecord(ctx, orgID, moduleID, recordID); err != nil {
		return nil, err
	}
	items, err := s.repo.ListAttachments(ctx, orgID, moduleID, recordID)
	if err != nil {
		return nil, err
	}
	out := make([]dto.AttachmentResponse, 0, len(items))
	for _, a := range items {
		out = append(out, dto.AttachmentResponse{
			ID: a.ID, FileName: a.FileName, FileURL: a.FileURL, PublicID: a.PublicID,
			ResourceType: a.ResourceType, FileSize: a.FileSize, UploadedBy: a.UploadedBy,
			UploaderName: a.UploaderName, CreatedAt: a.CreatedAt,
		})
	}
	return out, nil
}

func (s *Service) CreateAttachment(ctx context.Context, orgID, moduleID, recordID, userID string, req dto.CreateAttachmentRequest) (*dto.AttachmentResponse, error) {
	if err := s.ensureRecord(ctx, orgID, moduleID, recordID); err != nil {
		return nil, err
	}
	a, err := s.repo.CreateAttachment(
		ctx, orgID, moduleID, recordID, userID,
		req.FileName, req.FileURL, req.PublicID, req.ResourceType, req.FileSize,
	)
	if err != nil {
		return nil, err
	}
	meta, _ := json.Marshal(map[string]any{"file_name": req.FileName})
	_ = s.repo.CreateActivity(ctx, orgID, moduleID, recordID, userID, "ATTACHMENT_ADDED", "Attachment uploaded: "+req.FileName, meta)
	return &dto.AttachmentResponse{
		ID: a.ID, FileName: a.FileName, FileURL: a.FileURL, PublicID: a.PublicID,
		ResourceType: a.ResourceType, FileSize: a.FileSize, UploadedBy: a.UploadedBy,
		UploaderName: a.UploaderName, CreatedAt: a.CreatedAt,
	}, nil
}

func (s *Service) DeleteAttachment(ctx context.Context, orgID, moduleID, recordID, id, userID string) error {
	if err := s.ensureRecord(ctx, orgID, moduleID, recordID); err != nil {
		return err
	}
	ok, err := s.repo.DeleteAttachment(ctx, orgID, id)
	if err != nil {
		return err
	}
	if !ok {
		return ErrNotFound
	}
	_ = s.repo.CreateActivity(ctx, orgID, moduleID, recordID, userID, "ATTACHMENT_DELETED", "Attachment deleted", nil)
	return nil
}

func (s *Service) ListActivities(ctx context.Context, orgID, moduleID, recordID string) ([]dto.ActivityResponse, error) {
	if err := s.ensureRecord(ctx, orgID, moduleID, recordID); err != nil {
		return nil, err
	}
	items, err := s.repo.ListActivities(ctx, orgID, moduleID, recordID, 50)
	if err != nil {
		return nil, err
	}
	out := make([]dto.ActivityResponse, 0, len(items))
	for _, a := range items {
		out = append(out, dto.ActivityResponse{
			ID: a.ID, Action: a.Action, Description: a.Description,
			PerformedBy: a.PerformedBy, ActorName: a.ActorName,
			Metadata: a.Metadata, CreatedAt: a.CreatedAt,
		})
	}
	return out, nil
}

// LogRecordActivity is used by the record service on CUD.
func (s *Service) LogRecordActivity(ctx context.Context, orgID, moduleID, recordID, userID, action, description string, metadata map[string]any) error {
	var raw json.RawMessage
	if metadata != nil {
		b, err := json.Marshal(metadata)
		if err != nil {
			return err
		}
		raw = b
	}
	return s.repo.CreateActivity(ctx, orgID, moduleID, recordID, userID, action, description, raw)
}

func (s *Service) ListRelated(ctx context.Context, orgID, moduleID string) ([]dto.RelatedDescriptorResponse, error) {
	items, err := s.repo.ListRelatedDescriptors(ctx, orgID, moduleID)
	if err != nil {
		return nil, err
	}
	out := make([]dto.RelatedDescriptorResponse, 0, len(items))
	for _, d := range items {
		out = append(out, dto.RelatedDescriptorResponse{
			ChildModuleID: d.ChildModuleID, ChildModuleName: d.ChildModuleName,
			ChildAPIName: d.ChildAPIName, LookupFieldAPI: d.LookupFieldAPI,
			LookupFieldLabel: d.LookupFieldLabel,
		})
	}
	return out, nil
}

func (s *Service) RelatedLookupField(ctx context.Context, orgID, childModuleID, parentModuleID string) (string, error) {
	api, err := s.repo.LookupFieldAPI(ctx, orgID, childModuleID, parentModuleID)
	if err != nil {
		return "", fmt.Errorf("%w: %v", ErrNotFound, err)
	}
	return api, nil
}
