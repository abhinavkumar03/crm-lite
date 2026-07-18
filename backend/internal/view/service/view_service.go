package service

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/abhinavkumar03/crm-lite/backend/internal/view/dto"
	"github.com/abhinavkumar03/crm-lite/backend/internal/view/entity"
)

var (
	ErrModuleNotFound = errors.New("module not found")
	ErrNotFound       = errors.New("view not found")
	ErrForbidden      = errors.New("you cannot modify a view you do not own")
)

// Repository is the persistence contract for saved views.
type Repository interface {
	ModuleExists(ctx context.Context, orgID, moduleID string) (bool, error)
	Create(ctx context.Context, v *entity.View) error
	ListVisible(ctx context.Context, orgID, moduleID, userID string) ([]entity.View, error)
	GetByID(ctx context.Context, orgID, moduleID, id string) (*entity.View, error)
	Update(ctx context.Context, v *entity.View) error
	Delete(ctx context.Context, orgID, moduleID, id string) (bool, error)
	SetDefault(ctx context.Context, orgID, moduleID, id string) (bool, error)
}

type Service struct {
	repo Repository
}

func New(repo Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) List(ctx context.Context, orgID, moduleID, userID string) ([]dto.ViewResponse, error) {
	ok, err := s.repo.ModuleExists(ctx, orgID, moduleID)
	if err != nil {
		return nil, err
	}
	if !ok {
		return nil, ErrModuleNotFound
	}

	views, err := s.repo.ListVisible(ctx, orgID, moduleID, userID)
	if err != nil {
		return nil, err
	}
	out := make([]dto.ViewResponse, 0, len(views))
	for i := range views {
		out = append(out, toResponse(&views[i], userID))
	}
	return out, nil
}

func (s *Service) GetByID(ctx context.Context, orgID, moduleID, id, userID string) (*dto.ViewResponse, error) {
	v, err := s.repo.GetByID(ctx, orgID, moduleID, id)
	if err != nil {
		return nil, err
	}
	if v == nil {
		return nil, nil
	}
	resp := toResponse(v, userID)
	return &resp, nil
}

func (s *Service) Create(ctx context.Context, orgID, moduleID, userID string, req dto.CreateViewRequest) (*dto.ViewResponse, error) {
	ok, err := s.repo.ModuleExists(ctx, orgID, moduleID)
	if err != nil {
		return nil, err
	}
	if !ok {
		return nil, ErrModuleNotFound
	}

	columns, filters, sort, err := marshalConfig(req.Columns, req.Filters, req.Sort)
	if err != nil {
		return nil, err
	}

	owner := userID
	v := &entity.View{
		OrganizationID: orgID,
		ModuleID:       moduleID,
		Name:           req.Name,
		Columns:        columns,
		Filters:        filters,
		Sort:           sort,
		IsDefault:      false,
		IsPublic:       derefBool(req.IsPublic, true),
		OwnerID:        &owner,
	}

	if err := s.repo.Create(ctx, v); err != nil {
		return nil, err
	}

	resp := toResponse(v, userID)
	return &resp, nil
}

func (s *Service) Update(ctx context.Context, orgID, moduleID, id, userID string, req dto.UpdateViewRequest) (*dto.ViewResponse, error) {
	v, err := s.repo.GetByID(ctx, orgID, moduleID, id)
	if err != nil {
		return nil, err
	}
	if v == nil {
		return nil, ErrNotFound
	}
	if !canModify(v, userID) {
		return nil, ErrForbidden
	}

	if req.Name != nil {
		v.Name = *req.Name
	}
	if req.Columns != nil {
		columns, err := json.Marshal(req.Columns)
		if err != nil {
			return nil, err
		}
		v.Columns = columns
	}
	if req.Filters != nil {
		filters, err := json.Marshal(req.Filters)
		if err != nil {
			return nil, err
		}
		v.Filters = filters
	}
	if req.Sort != nil {
		sort, err := json.Marshal(req.Sort)
		if err != nil {
			return nil, err
		}
		v.Sort = sort
	}
	if req.IsPublic != nil {
		v.IsPublic = *req.IsPublic
	}

	if err := s.repo.Update(ctx, v); err != nil {
		return nil, err
	}

	resp := toResponse(v, userID)
	return &resp, nil
}

func (s *Service) Delete(ctx context.Context, orgID, moduleID, id, userID string) error {
	v, err := s.repo.GetByID(ctx, orgID, moduleID, id)
	if err != nil {
		return err
	}
	if v == nil {
		return ErrNotFound
	}
	if !canModify(v, userID) {
		return ErrForbidden
	}

	deleted, err := s.repo.Delete(ctx, orgID, moduleID, id)
	if err != nil {
		return err
	}
	if !deleted {
		return ErrNotFound
	}
	return nil
}

func (s *Service) SetDefault(ctx context.Context, orgID, moduleID, id string) error {
	ok, err := s.repo.SetDefault(ctx, orgID, moduleID, id)
	if err != nil {
		return err
	}
	if !ok {
		return ErrNotFound
	}
	return nil
}

func canModify(v *entity.View, userID string) bool {
	return v.OwnerID == nil || *v.OwnerID == userID
}

func marshalConfig(columns []string, filters []dto.ViewFilter, sort *dto.ViewSort) (col, fil, srt []byte, err error) {
	col, err = json.Marshal(columns)
	if err != nil {
		return nil, nil, nil, err
	}

	if filters == nil {
		filters = []dto.ViewFilter{}
	}
	fil, err = json.Marshal(filters)
	if err != nil {
		return nil, nil, nil, err
	}

	if sort == nil {
		sort = &dto.ViewSort{}
	}
	srt, err = json.Marshal(sort)
	if err != nil {
		return nil, nil, nil, err
	}

	return col, fil, srt, nil
}

func toResponse(v *entity.View, userID string) dto.ViewResponse {
	columns := []string{}
	_ = json.Unmarshal(v.Columns, &columns)

	filters := []dto.ViewFilter{}
	_ = json.Unmarshal(v.Filters, &filters)

	sort := dto.ViewSort{}
	_ = json.Unmarshal(v.Sort, &sort)

	return dto.ViewResponse{
		ID:        v.ID,
		ModuleID:  v.ModuleID,
		Name:      v.Name,
		Columns:   columns,
		Filters:   filters,
		Sort:      sort,
		IsDefault: v.IsDefault,
		IsPublic:  v.IsPublic,
		OwnerID:   v.OwnerID,
		IsOwner:   v.OwnerID != nil && *v.OwnerID == userID,
		CreatedAt: v.CreatedAt,
		UpdatedAt: v.UpdatedAt,
	}
}

func derefBool(p *bool, def bool) bool {
	if p == nil {
		return def
	}
	return *p
}
