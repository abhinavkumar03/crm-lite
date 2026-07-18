package service

import (
	"context"
	"encoding/json"
	"errors"
	"math"

	"github.com/abhinavkumar03/crm-lite/backend/internal/access"
	fieldentity "github.com/abhinavkumar03/crm-lite/backend/internal/field/entity"
	"github.com/abhinavkumar03/crm-lite/backend/internal/record/dto"
	"github.com/abhinavkumar03/crm-lite/backend/internal/record/entity"
	"github.com/abhinavkumar03/crm-lite/backend/internal/record/repository"
	"github.com/abhinavkumar03/crm-lite/backend/internal/shared/cache"
	vdto "github.com/abhinavkumar03/crm-lite/backend/internal/validationengine/dto"
)

const (
	DefaultPageSize = 20
	MaxPageSize     = 100

	storageDynamic = "dynamic"
)

var (
	ErrModuleNotFound = errors.New("module not found")
	ErrNotDynamic     = errors.New("module does not use dynamic record storage")
	ErrNotFound       = errors.New("record not found")
)

// ValidationError carries field-keyed validation failures so the handler can
// surface them with the same shape as the dry-run validate endpoint.
type ValidationError struct {
	Errors []vdto.FieldError
}

func (e *ValidationError) Error() string { return "validation failed" }

// RecordRepository is the persistence contract for records.
type RecordRepository interface {
	Create(ctx context.Context, rec *entity.Record) error
	GetByID(ctx context.Context, orgID, moduleID, id string) (*entity.Record, error)
	Update(ctx context.Context, rec *entity.Record) error
	Delete(ctx context.Context, orgID, moduleID, id string) (bool, error)
	List(ctx context.Context, orgID, moduleID string, q dto.ListQuery, meta map[string]repository.FieldMeta, extra repository.ExtraWhere) ([]entity.Record, int, error)
	DisplayValues(ctx context.Context, orgID, moduleID string, ids []string, displayField string) (map[string]string, error)
	UserDisplays(ctx context.Context, ids []string) (map[string]string, error)
}

// FieldReader exposes module + field metadata (satisfied by the field engine's
// repository — dependency inversion, no duplication).
type FieldReader interface {
	ModuleStorage(ctx context.Context, orgID, moduleID string) (string, bool, error)
	List(ctx context.Context, orgID, moduleID string) ([]fieldentity.Field, error)
}

// Validator evaluates a payload against a module's schema (the Phase 7 engine).
type Validator interface {
	Validate(ctx context.Context, orgID, moduleID string, data map[string]any) (vdto.ValidateResult, error)
}

type Service struct {
	repo      RecordRepository
	fields    FieldReader
	validator Validator
	cache     *cache.Cache
	access    *access.Service
}

func New(repo RecordRepository, fields FieldReader, validator Validator, appCache *cache.Cache, accessSvc *access.Service) *Service {
	return &Service{repo: repo, fields: fields, validator: validator, cache: appCache, access: accessSvc}
}

func (s *Service) invalidateDashboard(ctx context.Context, orgID string) {
	if s.cache != nil {
		s.cache.InvalidateDashboard(ctx, orgID)
	}
}

// ensureDynamicModule verifies the module exists in the org and uses dynamic
// storage (records live in the JSONB table only for 'dynamic' modules).
func (s *Service) ensureDynamicModule(ctx context.Context, orgID, moduleID string) error {
	strategy, ok, err := s.fields.ModuleStorage(ctx, orgID, moduleID)
	if err != nil {
		return err
	}
	if !ok {
		return ErrModuleNotFound
	}
	if strategy != storageDynamic {
		return ErrNotDynamic
	}
	return nil
}

func (s *Service) Create(ctx context.Context, orgID, moduleID, userID string, req dto.CreateRecordRequest) (*dto.RecordResponse, error) {
	if err := s.ensureDynamicModule(ctx, orgID, moduleID); err != nil {
		return nil, err
	}

	result, err := s.validator.Validate(ctx, orgID, moduleID, req.Data)
	if err != nil {
		return nil, err
	}
	if !result.Valid {
		return nil, &ValidationError{Errors: result.Errors}
	}

	data, err := json.Marshal(req.Data)
	if err != nil {
		return nil, err
	}

	owner := req.OwnerID
	if owner == nil {
		owner = &userID
	}

	vis := "organization"
	if req.Visibility != nil && *req.Visibility != "" {
		vis = *req.Visibility
	}

	rec := &entity.Record{
		OrganizationID: orgID,
		ModuleID:       moduleID,
		Data:           data,
		OwnerID:        owner,
		AssignedTo:     req.AssignedTo,
		TeamID:         req.TeamID,
		DepartmentID:   req.DepartmentID,
		Visibility:     vis,
		CreatedBy:      &userID,
		UpdatedBy:      &userID,
	}
	if err := s.repo.Create(ctx, rec); err != nil {
		return nil, err
	}

	s.invalidateDashboard(ctx, orgID)
	resp := toResponse(rec)
	return &resp, nil
}

func (s *Service) Update(ctx context.Context, orgID, moduleID, id, userID string, req dto.UpdateRecordRequest) (*dto.RecordResponse, error) {
	if err := s.ensureDynamicModule(ctx, orgID, moduleID); err != nil {
		return nil, err
	}

	existing, err := s.repo.GetByID(ctx, orgID, moduleID, id)
	if err != nil {
		return nil, err
	}
	if existing == nil {
		return nil, ErrNotFound
	}

	if s.access != nil {
		actor, err := s.access.LoadActor(ctx, orgID, userID)
		if err != nil {
			return nil, err
		}
		vis := existing.Visibility
		if !access.CanViewRecord(actor, existing.OwnerID, existing.AssignedTo, existing.CreatedBy, existing.DepartmentID, existing.TeamID, &vis) {
			return nil, ErrNotFound
		}
	}

	result, err := s.validator.Validate(ctx, orgID, moduleID, req.Data)
	if err != nil {
		return nil, err
	}
	if !result.Valid {
		return nil, &ValidationError{Errors: result.Errors}
	}

	data, err := json.Marshal(req.Data)
	if err != nil {
		return nil, err
	}

	existing.Data = data
	existing.UpdatedBy = &userID
	if req.OwnerID != nil {
		existing.OwnerID = req.OwnerID
	}
	if req.AssignedTo != nil {
		existing.AssignedTo = req.AssignedTo
	}
	if req.TeamID != nil {
		existing.TeamID = req.TeamID
	}
	if req.DepartmentID != nil {
		existing.DepartmentID = req.DepartmentID
	}
	if req.Visibility != nil && *req.Visibility != "" {
		existing.Visibility = *req.Visibility
	}

	if err := s.repo.Update(ctx, existing); err != nil {
		return nil, err
	}

	s.invalidateDashboard(ctx, orgID)
	resp := toResponse(existing)
	return &resp, nil
}

func (s *Service) Get(ctx context.Context, orgID, moduleID, id, userID string, expand bool) (*dto.RecordResponse, error) {
	if err := s.ensureDynamicModule(ctx, orgID, moduleID); err != nil {
		return nil, err
	}

	rec, err := s.repo.GetByID(ctx, orgID, moduleID, id)
	if err != nil {
		return nil, err
	}
	if rec == nil {
		return nil, ErrNotFound
	}

	if s.access != nil {
		actor, err := s.access.LoadActor(ctx, orgID, userID)
		if err != nil {
			return nil, err
		}
		vis := rec.Visibility
		if !access.CanViewRecord(actor, rec.OwnerID, rec.AssignedTo, rec.CreatedBy, rec.DepartmentID, rec.TeamID, &vis) {
			return nil, ErrNotFound
		}
	}

	resp := toResponse(rec)
	if expand {
		fields, err := s.fields.List(ctx, orgID, moduleID)
		if err != nil {
			return nil, err
		}
		if err := s.expand(ctx, orgID, fields, []*dto.RecordResponse{&resp}); err != nil {
			return nil, err
		}
	}
	return &resp, nil
}

func (s *Service) List(ctx context.Context, orgID, moduleID, userID string, q dto.ListQuery) (*dto.ListResult, error) {
	if err := s.ensureDynamicModule(ctx, orgID, moduleID); err != nil {
		return nil, err
	}

	fields, err := s.fields.List(ctx, orgID, moduleID)
	if err != nil {
		return nil, err
	}
	meta := repository.BuildMeta(fields)

	normalizeQuery(&q)

	extra := repository.ExtraWhere{SQL: "TRUE"}
	if s.access != nil && userID != "" {
		actor, err := s.access.LoadActor(ctx, orgID, userID)
		if err != nil {
			return nil, err
		}
		// BuildWhere uses $1,$2 for org/module — visibility args start at $3 + filter args.
		// Compute after BuildWhere by asking VisibilitySQL with a high start and then
		// re-building: List merges ExtraWhere after BuildWhere, so startArg = 3 + filter args.
		// Simpler: build base where first to know arg count.
		base := repository.BuildWhere(orgID, moduleID, q, meta)
		sql, args, _ := access.VisibilitySQL(actor, len(base.Args)+1)
		extra = repository.ExtraWhere{SQL: sql, Args: args}
	}

	records, total, err := s.repo.List(ctx, orgID, moduleID, q, meta, extra)
	if err != nil {
		return nil, err
	}

	responses := make([]dto.RecordResponse, len(records))
	ptrs := make([]*dto.RecordResponse, len(records))
	for i := range records {
		responses[i] = toResponse(&records[i])
		ptrs[i] = &responses[i]
	}

	if q.Expand {
		if err := s.expand(ctx, orgID, fields, ptrs); err != nil {
			return nil, err
		}
	}

	totalPages := 1
	if q.SkipTotal {
		// Hint that another page may exist when the page is full.
		if len(responses) == q.PageSize {
			totalPages = q.Page + 1
		} else {
			totalPages = q.Page
		}
	} else {
		totalPages = int(math.Max(1, math.Ceil(float64(total)/float64(q.PageSize))))
	}

	return &dto.ListResult{
		Records:    responses,
		Page:       q.Page,
		PageSize:   q.PageSize,
		Total:      total,
		TotalPages: totalPages,
	}, nil
}

func (s *Service) Delete(ctx context.Context, orgID, moduleID, id, userID string) error {
	if err := s.ensureDynamicModule(ctx, orgID, moduleID); err != nil {
		return err
	}

	existing, err := s.repo.GetByID(ctx, orgID, moduleID, id)
	if err != nil {
		return err
	}
	if existing == nil {
		return ErrNotFound
	}
	if s.access != nil {
		actor, err := s.access.LoadActor(ctx, orgID, userID)
		if err != nil {
			return err
		}
		vis := existing.Visibility
		if !access.CanViewRecord(actor, existing.OwnerID, existing.AssignedTo, existing.CreatedBy, existing.DepartmentID, existing.TeamID, &vis) {
			return ErrNotFound
		}
	}

	deleted, err := s.repo.Delete(ctx, orgID, moduleID, id)
	if err != nil {
		return err
	}
	if !deleted {
		return ErrNotFound
	}
	s.invalidateDashboard(ctx, orgID)
	return nil
}

func normalizeQuery(q *dto.ListQuery) {
	if q.Page < 1 {
		q.Page = 1
	}
	if q.PageSize < 1 {
		q.PageSize = DefaultPageSize
	}
	if q.PageSize > MaxPageSize {
		q.PageSize = MaxPageSize
	}
}

func toResponse(rec *entity.Record) dto.RecordResponse {
	data := map[string]any{}
	if len(rec.Data) > 0 {
		_ = json.Unmarshal(rec.Data, &data)
	}
	vis := rec.Visibility
	if vis == "" {
		vis = "organization"
	}
	return dto.RecordResponse{
		ID:           rec.ID,
		ModuleID:     rec.ModuleID,
		Data:         data,
		OwnerID:      rec.OwnerID,
		AssignedTo:   rec.AssignedTo,
		TeamID:       rec.TeamID,
		DepartmentID: rec.DepartmentID,
		Visibility:   vis,
		CreatedBy:    rec.CreatedBy,
		UpdatedBy:    rec.UpdatedBy,
		CreatedAt:    rec.CreatedAt,
		UpdatedAt:    rec.UpdatedAt,
	}
}
