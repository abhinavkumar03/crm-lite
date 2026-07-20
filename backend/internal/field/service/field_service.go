package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"regexp"
	"strings"

	"github.com/abhinavkumar03/crm-lite/backend/internal/field/dto"
	"github.com/abhinavkumar03/crm-lite/backend/internal/field/entity"
)

// Domain errors mapped to HTTP status codes by the handler.
var (
	ErrModuleNotFound   = errors.New("module not found")
	ErrNotFound         = errors.New("field not found")
	ErrInvalidAPIName   = errors.New("api_name must start with a letter and contain only lowercase letters, digits and underscores")
	ErrDuplicateAPIName = errors.New("a field with this api_name already exists on the module")
	ErrInvalidType      = errors.New("unsupported field_type")
	ErrInvalidLockMode  = errors.New("lock_mode must be never, after_create, or always")
	ErrInvalidACL       = errors.New("editable_by and viewable_by must be ALL for now")
	ErrOptionsRequired  = errors.New("this field_type requires a non-empty options list")
	ErrLookupRequired   = errors.New("lookup fields require a valid lookup_module_id")
	ErrInvalidLength    = errors.New("min_length cannot be greater than max_length")
	ErrSystemField      = errors.New("system fields cannot be deleted")
)

var apiNamePattern = regexp.MustCompile(`^[a-z][a-z0-9_]*$`)

// choiceTypes require an options list.
var choiceTypes = map[string]bool{
	entity.TypeDropdown:    true,
	entity.TypeMultiselect: true,
	entity.TypeRadio:       true,
}

// typeAliases map PRD / common names onto canonical stored types.
var typeAliases = map[string]string{
	"select":       entity.TypeDropdown,
	"multi_select": entity.TypeMultiselect,
	"multiselect":  entity.TypeMultiselect,
	"user_lookup":  entity.TypeUser,
	"percent":      entity.TypePercentage,
}

var validTypes = func() map[string]bool {
	m := make(map[string]bool, len(entity.AllTypes))
	for _, t := range entity.AllTypes {
		m[t] = true
	}
	return m
}()

var validLockModes = func() map[string]bool {
	m := make(map[string]bool, len(entity.AllLockModes))
	for _, t := range entity.AllLockModes {
		m[t] = true
	}
	return m
}()

// Repository is the persistence contract this service depends on.
type Repository interface {
	ModuleStorage(ctx context.Context, orgID, moduleID string) (string, bool, error)
	ModuleExistsInOrg(ctx context.Context, orgID, moduleID string) (bool, error)
	Create(ctx context.Context, f *entity.Field) error
	List(ctx context.Context, orgID, moduleID string) ([]entity.Field, error)
	GetByID(ctx context.Context, orgID, moduleID, id string) (*entity.Field, error)
	Update(ctx context.Context, f *entity.Field) error
	Delete(ctx context.Context, orgID, moduleID, id string) (bool, error)
	ExistsByAPIName(ctx context.Context, moduleID, apiName string) (bool, error)
	MaxSortOrder(ctx context.Context, moduleID string) (int, error)
	Reorder(ctx context.Context, orgID, moduleID string, positions []entity.SortPosition) error
}

type Service struct {
	repo   Repository
	layout LayoutSync
}

// LayoutSync keeps detail layout sections in sync with field create/delete.
type LayoutSync interface {
	AppendFieldToSection(ctx context.Context, orgID, moduleID, sectionKey, apiName string) error
	RemoveFieldFromLayout(ctx context.Context, orgID, moduleID, apiName string) error
}

func New(repo Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) SetLayoutSync(layout LayoutSync) {
	s.layout = layout
}

func (s *Service) List(ctx context.Context, orgID, moduleID string) ([]dto.FieldResponse, error) {
	strategy, ok, err := s.repo.ModuleStorage(ctx, orgID, moduleID)
	if err != nil {
		return nil, err
	}
	if !ok {
		return nil, ErrModuleNotFound
	}

	fields, err := s.repo.List(ctx, orgID, moduleID)
	if err != nil {
		return nil, err
	}
	out := make([]dto.FieldResponse, 0, len(fields))
	for i := range fields {
		out = append(out, toResponse(&fields[i], strategy))
	}
	return out, nil
}

func (s *Service) GetByID(ctx context.Context, orgID, moduleID, id string) (*dto.FieldResponse, error) {
	strategy, ok, err := s.repo.ModuleStorage(ctx, orgID, moduleID)
	if err != nil {
		return nil, err
	}
	if !ok {
		return nil, ErrModuleNotFound
	}

	f, err := s.repo.GetByID(ctx, orgID, moduleID, id)
	if err != nil {
		return nil, err
	}
	if f == nil {
		return nil, nil
	}
	resp := toResponse(f, strategy)
	return &resp, nil
}

func (s *Service) Create(ctx context.Context, orgID, moduleID string, req dto.CreateFieldRequest) (*dto.FieldResponse, error) {
	strategy, ok, err := s.repo.ModuleStorage(ctx, orgID, moduleID)
	if err != nil {
		return nil, err
	}
	if !ok {
		return nil, ErrModuleNotFound
	}

	if !apiNamePattern.MatchString(req.APIName) {
		return nil, ErrInvalidAPIName
	}
	fieldType := canonicalizeType(req.FieldType)
	if !validTypes[fieldType] {
		return nil, ErrInvalidType
	}
	if choiceTypes[fieldType] && len(req.Options) == 0 {
		return nil, ErrOptionsRequired
	}
	if err := validateLength(req.MinLength, req.MaxLength); err != nil {
		return nil, err
	}
	lockMode, err := normalizeLockMode(req.LockMode)
	if err != nil {
		return nil, err
	}
	editableBy, viewableBy, err := normalizeACL(req.EditableBy, req.ViewableBy)
	if err != nil {
		return nil, err
	}

	if err := s.validateLookup(ctx, orgID, fieldType, req.LookupModuleID); err != nil {
		return nil, err
	}

	exists, err := s.repo.ExistsByAPIName(ctx, moduleID, req.APIName)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, ErrDuplicateAPIName
	}

	nextSort, err := s.repo.MaxSortOrder(ctx, moduleID)
	if err != nil {
		return nil, err
	}

	options, err := marshalOptions(req.Options)
	if err != nil {
		return nil, err
	}

	f := &entity.Field{
		OrganizationID:    orgID,
		ModuleID:          moduleID,
		APIName:           req.APIName,
		Label:             req.Label,
		FieldType:         fieldType,
		IsRequired:        req.IsRequired,
		IsUnique:          req.IsUnique,
		IsReadOnly:        req.IsReadOnly,
		DefaultValue:      req.DefaultValue,
		Placeholder:       req.Placeholder,
		Description:       req.Description,
		HelpText:          req.HelpText,
		MinLength:         req.MinLength,
		MaxLength:         req.MaxLength,
		Regex:             req.Regex,
		ValidationMessage: req.ValidationMessage,
		Options:           options,
		LookupModuleID:    req.LookupModuleID,
		SortOrder:         nextSort + 1,
		IsVisible:         derefBool(req.IsVisible, true),
		IsSearchable:      req.IsSearchable,
		IsFilterable:      req.IsFilterable,
		IsSystem:          false,
		LockMode:          lockMode,
		EditableBy:        editableBy,
		ViewableBy:        viewableBy,
	}

	if err := s.repo.Create(ctx, f); err != nil {
		return nil, err
	}

	if s.layout != nil {
		section := strings.TrimSpace(req.SectionKey)
		if section == "" || section == "system" {
			section = "general"
		}
		if err := s.layout.AppendFieldToSection(ctx, orgID, moduleID, section, f.APIName); err != nil {
			return nil, fmt.Errorf("field created but layout sync failed: %w", err)
		}
	}

	resp := toResponse(f, strategy)
	return &resp, nil
}

func (s *Service) Update(ctx context.Context, orgID, moduleID, id string, req dto.UpdateFieldRequest) (*dto.FieldResponse, error) {
	strategy, ok, err := s.repo.ModuleStorage(ctx, orgID, moduleID)
	if err != nil {
		return nil, err
	}
	if !ok {
		return nil, ErrModuleNotFound
	}

	f, err := s.repo.GetByID(ctx, orgID, moduleID, id)
	if err != nil {
		return nil, err
	}
	if f == nil {
		return nil, ErrNotFound
	}

	if req.Label != nil {
		f.Label = *req.Label
	}
	if req.IsRequired != nil {
		f.IsRequired = *req.IsRequired
	}
	if req.IsUnique != nil {
		f.IsUnique = *req.IsUnique
	}
	if req.IsReadOnly != nil {
		f.IsReadOnly = *req.IsReadOnly
	}
	if req.DefaultValue != nil {
		f.DefaultValue = req.DefaultValue
	}
	if req.Placeholder != nil {
		f.Placeholder = req.Placeholder
	}
	if req.Description != nil {
		f.Description = req.Description
	}
	if req.HelpText != nil {
		f.HelpText = req.HelpText
	}
	if req.MinLength != nil {
		f.MinLength = req.MinLength
	}
	if req.MaxLength != nil {
		f.MaxLength = req.MaxLength
	}
	if req.Regex != nil {
		f.Regex = req.Regex
	}
	if req.ValidationMessage != nil {
		f.ValidationMessage = req.ValidationMessage
	}
	if req.IsVisible != nil {
		f.IsVisible = *req.IsVisible
	}
	if req.IsSearchable != nil {
		f.IsSearchable = *req.IsSearchable
	}
	if req.IsFilterable != nil {
		f.IsFilterable = *req.IsFilterable
	}
	if req.LockMode != nil {
		mode, err := normalizeLockMode(*req.LockMode)
		if err != nil {
			return nil, err
		}
		f.LockMode = mode
	}
	if req.EditableBy != nil || req.ViewableBy != nil {
		ed := f.EditableBy
		vw := f.ViewableBy
		if req.EditableBy != nil {
			ed = *req.EditableBy
		}
		if req.ViewableBy != nil {
			vw = *req.ViewableBy
		}
		ed, vw, err = normalizeACL(ed, vw)
		if err != nil {
			return nil, err
		}
		f.EditableBy = ed
		f.ViewableBy = vw
	}
	if req.Options != nil {
		options, err := marshalOptions(req.Options)
		if err != nil {
			return nil, err
		}
		f.Options = options
	}

	if err := validateLength(f.MinLength, f.MaxLength); err != nil {
		return nil, err
	}
	if choiceTypes[f.FieldType] && len(f.Options) == 0 {
		return nil, ErrOptionsRequired
	}

	if err := s.repo.Update(ctx, f); err != nil {
		return nil, err
	}

	resp := toResponse(f, strategy)
	return &resp, nil
}

func (s *Service) Delete(ctx context.Context, orgID, moduleID, id string) error {
	_, ok, err := s.repo.ModuleStorage(ctx, orgID, moduleID)
	if err != nil {
		return err
	}
	if !ok {
		return ErrModuleNotFound
	}

	f, err := s.repo.GetByID(ctx, orgID, moduleID, id)
	if err != nil {
		return err
	}
	if f == nil {
		return ErrNotFound
	}
	if f.IsSystem {
		return ErrSystemField
	}

	deleted, err := s.repo.Delete(ctx, orgID, moduleID, id)
	if err != nil {
		return err
	}
	if !deleted {
		return ErrNotFound
	}
	if s.layout != nil {
		if err := s.layout.RemoveFieldFromLayout(ctx, orgID, moduleID, f.APIName); err != nil {
			return fmt.Errorf("field deleted but layout sync failed: %w", err)
		}
	}
	return nil
}

func (s *Service) Reorder(ctx context.Context, orgID, moduleID string, items []dto.ReorderItem) error {
	_, ok, err := s.repo.ModuleStorage(ctx, orgID, moduleID)
	if err != nil {
		return err
	}
	if !ok {
		return ErrModuleNotFound
	}

	positions := make([]entity.SortPosition, 0, len(items))
	for _, it := range items {
		positions = append(positions, entity.SortPosition{ID: it.ID, SortOrder: it.SortOrder})
	}
	return s.repo.Reorder(ctx, orgID, moduleID, positions)
}

func (s *Service) validateLookup(ctx context.Context, orgID, fieldType string, lookupModuleID *string) error {
	if fieldType != entity.TypeLookup {
		return nil
	}
	if lookupModuleID == nil || *lookupModuleID == "" {
		return ErrLookupRequired
	}
	ok, err := s.repo.ModuleExistsInOrg(ctx, orgID, *lookupModuleID)
	if err != nil {
		return err
	}
	if !ok {
		return ErrLookupRequired
	}
	return nil
}

func validateLength(min, max *int) error {
	if min != nil && max != nil && *min > *max {
		return ErrInvalidLength
	}
	return nil
}

// storageFor derives the persistence descriptor from the module's strategy.
// Native modules map fields to real columns; dynamic modules store them as keys
// inside records.data.
func storageFor(strategy, apiName string) dto.StorageDescriptor {
	if strategy == "native" {
		return dto.StorageDescriptor{Kind: dto.StorageColumn, Path: apiName}
	}
	return dto.StorageDescriptor{Kind: dto.StorageJSONB, Path: "data." + apiName}
}

func canonicalizeType(raw string) string {
	t := strings.TrimSpace(strings.ToLower(raw))
	if alias, ok := typeAliases[t]; ok {
		return alias
	}
	return t
}

func normalizeLockMode(raw string) (string, error) {
	mode := strings.TrimSpace(raw)
	if mode == "" {
		return entity.LockNever, nil
	}
	if !validLockModes[mode] {
		return "", ErrInvalidLockMode
	}
	return mode, nil
}

func normalizeACL(editableBy, viewableBy string) (string, string, error) {
	ed := strings.TrimSpace(editableBy)
	vw := strings.TrimSpace(viewableBy)
	if ed == "" {
		ed = "ALL"
	}
	if vw == "" {
		vw = "ALL"
	}
	// Schema is future-ready; v1 only accepts ALL.
	if !strings.EqualFold(ed, "ALL") || !strings.EqualFold(vw, "ALL") {
		return "", "", ErrInvalidACL
	}
	return "ALL", "ALL", nil
}

func toResponse(f *entity.Field, strategy string) dto.FieldResponse {
	lockMode := f.LockMode
	if lockMode == "" {
		lockMode = entity.LockNever
	}
	editableBy := f.EditableBy
	if editableBy == "" {
		editableBy = "ALL"
	}
	viewableBy := f.ViewableBy
	if viewableBy == "" {
		viewableBy = "ALL"
	}
	return dto.FieldResponse{
		ID:                f.ID,
		ModuleID:          f.ModuleID,
		APIName:           f.APIName,
		Label:             f.Label,
		FieldType:         f.FieldType,
		IsRequired:        f.IsRequired,
		IsUnique:          f.IsUnique,
		IsReadOnly:        f.IsReadOnly,
		DefaultValue:      f.DefaultValue,
		Placeholder:       f.Placeholder,
		Description:       f.Description,
		HelpText:          f.HelpText,
		MinLength:         f.MinLength,
		MaxLength:         f.MaxLength,
		Regex:             f.Regex,
		ValidationMessage: f.ValidationMessage,
		Options:           parseOptions(f.Options),
		LookupModuleID:    f.LookupModuleID,
		SortOrder:         f.SortOrder,
		IsVisible:         f.IsVisible,
		IsSearchable:      f.IsSearchable,
		IsFilterable:      f.IsFilterable,
		IsNullable:        f.IsNullable,
		IsIndexed:         f.IsIndexed,
		IsSystem:          f.IsSystem,
		LockMode:          lockMode,
		EditableBy:        editableBy,
		ViewableBy:        viewableBy,
		Storage:           storageFor(strategy, f.APIName),
		CreatedAt:         f.CreatedAt,
		UpdatedAt:         f.UpdatedAt,
	}
}

// parseOptions normalizes stored options into {label,value} pairs. It accepts
// both the legacy plain-string array form (["NEW","WON"]) and the structured
// object form ([{"label":"New","value":"NEW"}]).
func parseOptions(raw []byte) []dto.FieldOption {
	out := make([]dto.FieldOption, 0)
	if len(raw) == 0 {
		return out
	}

	var strs []string
	if err := json.Unmarshal(raw, &strs); err == nil {
		for _, s := range strs {
			out = append(out, dto.FieldOption{Label: s, Value: s})
		}
		return out
	}

	var objs []dto.FieldOption
	if err := json.Unmarshal(raw, &objs); err == nil {
		return objs
	}

	return out
}

func marshalOptions(opts []dto.FieldOption) ([]byte, error) {
	if len(opts) == 0 {
		return []byte("[]"), nil
	}
	return json.Marshal(opts)
}

func derefBool(p *bool, def bool) bool {
	if p == nil {
		return def
	}
	return *p
}
