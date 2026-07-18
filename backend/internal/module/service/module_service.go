package service

import (
	"context"
	"errors"
	"regexp"

	"github.com/abhinavkumar03/crm-lite/backend/internal/module/dto"
	"github.com/abhinavkumar03/crm-lite/backend/internal/module/entity"
)

// Domain errors mapped to HTTP status codes by the handler.
var (
	ErrInvalidAPIName   = errors.New("api_name must start with a letter and contain only lowercase letters, digits and underscores")
	ErrDuplicateAPIName = errors.New("a module with this api_name already exists")
	ErrSystemModule     = errors.New("system modules cannot be deleted")
	ErrNotFound         = errors.New("module not found")
)

var apiNamePattern = regexp.MustCompile(`^[a-z][a-z0-9_]*$`)

// Repository is the persistence contract this service depends on. It is defined
// in the consumer package (idiomatic Go) so the service is easy to mock/test.
type Repository interface {
	Create(ctx context.Context, m *entity.Module) error
	List(ctx context.Context, orgID string) ([]entity.Module, error)
	Navigation(ctx context.Context, orgID string) ([]entity.Module, error)
	GetByID(ctx context.Context, orgID, id string) (*entity.Module, error)
	Update(ctx context.Context, m *entity.Module) error
	SetEnabled(ctx context.Context, orgID, id string, enabled bool) (bool, error)
	Delete(ctx context.Context, orgID, id string) (bool, error)
	ExistsByAPIName(ctx context.Context, orgID, apiName string) (bool, error)
	MaxSortOrder(ctx context.Context, orgID string) (int, error)
	Reorder(ctx context.Context, orgID string, positions []entity.SortPosition) error
}

type Service struct {
	repo Repository
}

func New(repo Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) Create(ctx context.Context, orgID string, req dto.CreateModuleRequest) (*dto.ModuleResponse, error) {
	if !apiNamePattern.MatchString(req.APIName) {
		return nil, ErrInvalidAPIName
	}

	exists, err := s.repo.ExistsByAPIName(ctx, orgID, req.APIName)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, ErrDuplicateAPIName
	}

	nextSort, err := s.repo.MaxSortOrder(ctx, orgID)
	if err != nil {
		return nil, err
	}

	m := &entity.Module{
		OrganizationID:   orgID,
		APIName:          req.APIName,
		SingularLabel:    req.SingularLabel,
		PluralLabel:      req.PluralLabel,
		Description:      req.Description,
		Icon:             req.Icon,
		Color:            req.Color,
		StorageStrategy:  entity.StorageDynamic, // user-created modules are always dynamic
		IsSystem:         false,
		IsEnabled:        true,
		IsVisibleSidebar: derefBool(req.IsVisibleSidebar, true),
		SortOrder:        nextSort + 1,
		DefaultSortField: derefString(req.DefaultSortField, "created_at"),
		DefaultSortOrder: derefString(req.DefaultSortOrder, "desc"),
	}

	if err := s.repo.Create(ctx, m); err != nil {
		return nil, err
	}

	resp := toResponse(m)
	return &resp, nil
}

func (s *Service) List(ctx context.Context, orgID string) ([]dto.ModuleResponse, error) {
	modules, err := s.repo.List(ctx, orgID)
	if err != nil {
		return nil, err
	}
	out := make([]dto.ModuleResponse, 0, len(modules))
	for i := range modules {
		out = append(out, toResponse(&modules[i]))
	}
	return out, nil
}

func (s *Service) Navigation(ctx context.Context, orgID string) ([]dto.NavigationItem, error) {
	modules, err := s.repo.Navigation(ctx, orgID)
	if err != nil {
		return nil, err
	}
	out := make([]dto.NavigationItem, 0, len(modules))
	for i := range modules {
		m := &modules[i]
		out = append(out, dto.NavigationItem{
			ID:            m.ID,
			APIName:       m.APIName,
			SingularLabel: m.SingularLabel,
			PluralLabel:   m.PluralLabel,
			Icon:          m.Icon,
			Color:         m.Color,
			SortOrder:     m.SortOrder,
		})
	}
	return out, nil
}

func (s *Service) GetByID(ctx context.Context, orgID, id string) (*dto.ModuleResponse, error) {
	m, err := s.repo.GetByID(ctx, orgID, id)
	if err != nil {
		return nil, err
	}
	if m == nil {
		return nil, nil
	}
	resp := toResponse(m)
	return &resp, nil
}

func (s *Service) Update(ctx context.Context, orgID, id string, req dto.UpdateModuleRequest) (*dto.ModuleResponse, error) {
	m, err := s.repo.GetByID(ctx, orgID, id)
	if err != nil {
		return nil, err
	}
	if m == nil {
		return nil, ErrNotFound
	}

	if req.SingularLabel != nil {
		m.SingularLabel = *req.SingularLabel
	}
	if req.PluralLabel != nil {
		m.PluralLabel = *req.PluralLabel
	}
	if req.Description != nil {
		m.Description = req.Description
	}
	if req.Icon != nil {
		m.Icon = req.Icon
	}
	if req.Color != nil {
		m.Color = req.Color
	}
	if req.IsVisibleSidebar != nil {
		m.IsVisibleSidebar = *req.IsVisibleSidebar
	}
	if req.DefaultSortField != nil {
		m.DefaultSortField = *req.DefaultSortField
	}
	if req.DefaultSortOrder != nil {
		m.DefaultSortOrder = *req.DefaultSortOrder
	}

	if err := s.repo.Update(ctx, m); err != nil {
		return nil, err
	}

	resp := toResponse(m)
	return &resp, nil
}

// SetEnabled toggles a module. Returns ErrNotFound if the module doesn't exist.
func (s *Service) SetEnabled(ctx context.Context, orgID, id string, enabled bool) error {
	found, err := s.repo.SetEnabled(ctx, orgID, id, enabled)
	if err != nil {
		return err
	}
	if !found {
		return ErrNotFound
	}
	return nil
}

// Delete removes a module. System modules are protected.
func (s *Service) Delete(ctx context.Context, orgID, id string) error {
	m, err := s.repo.GetByID(ctx, orgID, id)
	if err != nil {
		return err
	}
	if m == nil {
		return ErrNotFound
	}
	if m.IsSystem {
		return ErrSystemModule
	}

	deleted, err := s.repo.Delete(ctx, orgID, id)
	if err != nil {
		return err
	}
	if !deleted {
		return ErrNotFound
	}
	return nil
}

func (s *Service) Reorder(ctx context.Context, orgID string, items []dto.ReorderItem) error {
	positions := make([]entity.SortPosition, 0, len(items))
	for _, it := range items {
		positions = append(positions, entity.SortPosition{ID: it.ID, SortOrder: it.SortOrder})
	}
	return s.repo.Reorder(ctx, orgID, positions)
}

func toResponse(m *entity.Module) dto.ModuleResponse {
	return dto.ModuleResponse{
		ID:               m.ID,
		APIName:          m.APIName,
		SingularLabel:    m.SingularLabel,
		PluralLabel:      m.PluralLabel,
		Description:      m.Description,
		Icon:             m.Icon,
		Color:            m.Color,
		StorageStrategy:  m.StorageStrategy,
		NativeTable:      m.NativeTable,
		IsSystem:         m.IsSystem,
		IsEnabled:        m.IsEnabled,
		IsVisibleSidebar: m.IsVisibleSidebar,
		SortOrder:        m.SortOrder,
		DefaultSortField: m.DefaultSortField,
		DefaultSortOrder: m.DefaultSortOrder,
		CreatedAt:        m.CreatedAt,
		UpdatedAt:        m.UpdatedAt,
	}
}

func derefBool(p *bool, def bool) bool {
	if p == nil {
		return def
	}
	return *p
}

func derefString(p *string, def string) string {
	if p == nil || *p == "" {
		return def
	}
	return *p
}
