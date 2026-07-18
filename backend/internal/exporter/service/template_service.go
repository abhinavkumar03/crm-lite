package service

import (
	"context"
	"encoding/json"

	"github.com/abhinavkumar03/crm-lite/backend/internal/exporter/dto"
	"github.com/abhinavkumar03/crm-lite/backend/internal/exporter/entity"
	recorddto "github.com/abhinavkumar03/crm-lite/backend/internal/record/dto"
)

// ListTemplates returns the module's saved export templates.
func (s *Service) ListTemplates(ctx context.Context, orgID, moduleID string) ([]dto.TemplateResponse, error) {
	if err := s.ensureDynamicModule(ctx, orgID, moduleID); err != nil {
		return nil, err
	}

	items, err := s.templates.List(ctx, orgID, moduleID)
	if err != nil {
		return nil, err
	}
	out := make([]dto.TemplateResponse, 0, len(items))
	for i := range items {
		out = append(out, toTemplateResponse(&items[i]))
	}
	return out, nil
}

func (s *Service) CreateTemplate(ctx context.Context, orgID, moduleID, userID string, req dto.CreateTemplateRequest) (*dto.TemplateResponse, error) {
	if err := s.ensureDynamicModule(ctx, orgID, moduleID); err != nil {
		return nil, err
	}

	columns, _ := json.Marshal(orEmptyStrings(req.Columns))
	filters, _ := json.Marshal(orEmptyFilters(req.Filters))
	sort, _ := json.Marshal(orEmptySort(req.Sort))

	t := &entity.ExportTemplate{
		OrganizationID: orgID,
		ModuleID:       moduleID,
		Name:           req.Name,
		Format:         normalizeFormat(req.Format),
		Columns:        columns,
		Filters:        filters,
		Sort:           sort,
		CreatedBy:      &userID,
	}
	if err := s.templates.Create(ctx, t); err != nil {
		return nil, err
	}

	resp := toTemplateResponse(t)
	return &resp, nil
}

func (s *Service) UpdateTemplate(ctx context.Context, orgID, moduleID, id string, req dto.UpdateTemplateRequest) (*dto.TemplateResponse, error) {
	t, err := s.templates.GetByID(ctx, orgID, moduleID, id)
	if err != nil {
		return nil, err
	}
	if t == nil {
		return nil, ErrNotFound
	}

	if req.Name != nil {
		t.Name = *req.Name
	}
	if req.Format != nil {
		t.Format = normalizeFormat(*req.Format)
	}
	if req.Columns != nil {
		t.Columns, _ = json.Marshal(req.Columns)
	}
	if req.Filters != nil {
		t.Filters, _ = json.Marshal(req.Filters)
	}
	if req.Sort != nil {
		t.Sort, _ = json.Marshal(req.Sort)
	}

	if err := s.templates.Update(ctx, t); err != nil {
		return nil, err
	}

	resp := toTemplateResponse(t)
	return &resp, nil
}

func (s *Service) DeleteTemplate(ctx context.Context, orgID, moduleID, id string) error {
	deleted, err := s.templates.Delete(ctx, orgID, moduleID, id)
	if err != nil {
		return err
	}
	if !deleted {
		return ErrNotFound
	}
	return nil
}

func toTemplateResponse(t *entity.ExportTemplate) dto.TemplateResponse {
	columns := []string{}
	if len(t.Columns) > 0 {
		_ = json.Unmarshal(t.Columns, &columns)
	}
	filters := []recorddto.FilterClause{}
	if len(t.Filters) > 0 {
		_ = json.Unmarshal(t.Filters, &filters)
	}
	sort := dto.TemplateSort{}
	if len(t.Sort) > 0 {
		_ = json.Unmarshal(t.Sort, &sort)
	}
	return dto.TemplateResponse{
		ID:        t.ID,
		ModuleID:  t.ModuleID,
		Name:      t.Name,
		Format:    t.Format,
		Columns:   columns,
		Filters:   filters,
		Sort:      sort,
		CreatedBy: t.CreatedBy,
		CreatedAt: t.CreatedAt,
		UpdatedAt: t.UpdatedAt,
	}
}

func orEmptySort(s *dto.TemplateSort) dto.TemplateSort {
	if s == nil {
		return dto.TemplateSort{}
	}
	return *s
}
