package service

import (
	"context"
	"encoding/json"
	"errors"
	"math"

	"github.com/hibiken/asynq"

	fieldentity "github.com/abhinavkumar03/crm-lite/backend/internal/field/entity"
	"github.com/abhinavkumar03/crm-lite/backend/internal/importer/dto"
	"github.com/abhinavkumar03/crm-lite/backend/internal/importer/entity"
	"github.com/abhinavkumar03/crm-lite/backend/internal/importer/parser"
	"github.com/abhinavkumar03/crm-lite/backend/internal/jobs"
)

const (
	DefaultPageSize = 20
	MaxPageSize     = 100

	// MaxRows bounds a single import so a staged file cannot bloat the row or
	// starve the worker. Larger files should be split by the caller.
	MaxRows = 5000

	// sampleSize is how many rows the analyze step returns for preview.
	sampleSize = 10

	storageDynamic = "dynamic"
)

var (
	ErrModuleNotFound = errors.New("module not found")
	ErrNotDynamic     = errors.New("module does not use dynamic record storage")
	ErrNotFound       = errors.New("import job not found")
	ErrNoMapping      = errors.New("no valid column-to-field mappings were provided")
	ErrTooManyRows    = errors.New("file exceeds the maximum number of rows")
)

// FieldReader exposes module + field metadata (satisfied by the field engine's
// repository — dependency inversion, no duplication).
type FieldReader interface {
	ModuleStorage(ctx context.Context, orgID, moduleID string) (string, bool, error)
	List(ctx context.Context, orgID, moduleID string) ([]fieldentity.Field, error)
}

// Enqueuer publishes jobs onto the async queue (satisfied by *jobs.Producer).
type Enqueuer interface {
	Publish(ctx context.Context, job jobs.Job, opts ...asynq.Option) error
}

// Repository is the persistence contract for import jobs.
type Repository interface {
	Create(ctx context.Context, j *entity.ImportJob) error
	GetByID(ctx context.Context, orgID, id string) (*entity.ImportJob, error)
	List(ctx context.Context, orgID, moduleID string, q dto.ListQuery) ([]entity.ImportJob, int, error)
	Finish(ctx context.Context, id, status string, processed, success, errorRows int, errs []byte) error
}

type Service struct {
	repo     Repository
	fields   FieldReader
	enqueuer Enqueuer
}

func New(repo Repository, fields FieldReader, enqueuer Enqueuer) *Service {
	return &Service{repo: repo, fields: fields, enqueuer: enqueuer}
}

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

// Analyze parses the file (without persisting it) and returns everything the
// mapping UI needs: columns, a preview sample and an auto-suggested mapping.
func (s *Service) Analyze(ctx context.Context, orgID, moduleID, filename string, data []byte) (*dto.AnalyzeResult, error) {
	if err := s.ensureDynamicModule(ctx, orgID, moduleID); err != nil {
		return nil, err
	}

	parsed, err := parser.ParseFile(filename, data)
	if err != nil {
		return nil, err
	}

	fields, err := s.fields.List(ctx, orgID, moduleID)
	if err != nil {
		return nil, err
	}

	sample := parsed.Rows
	if len(sample) > sampleSize {
		sample = sample[:sampleSize]
	}

	return &dto.AnalyzeResult{
		Headers:          parsed.Headers,
		SampleRows:       sample,
		SuggestedMapping: suggestMapping(parsed.Headers, fields),
		RowCount:         len(parsed.Rows),
	}, nil
}

// Create parses the full file, stages the rows, persists a pending job and
// enqueues it for asynchronous processing.
func (s *Service) Create(
	ctx context.Context,
	orgID, moduleID, userID, filename string,
	data []byte,
	mapping map[string]string,
	options map[string]any,
) (*dto.ImportResponse, error) {
	if err := s.ensureDynamicModule(ctx, orgID, moduleID); err != nil {
		return nil, err
	}

	parsed, err := parser.ParseFile(filename, data)
	if err != nil {
		return nil, err
	}
	if len(parsed.Rows) > MaxRows {
		return nil, ErrTooManyRows
	}

	fields, err := s.fields.List(ctx, orgID, moduleID)
	if err != nil {
		return nil, err
	}

	// Keep only mappings whose source column exists and whose target is a real,
	// writable field of the module — never trust the client's field names.
	cleanMapping := sanitizeMapping(mapping, parsed.Headers, fields)
	if len(cleanMapping) == 0 {
		return nil, ErrNoMapping
	}

	mappingJSON, err := json.Marshal(cleanMapping)
	if err != nil {
		return nil, err
	}
	rowsJSON, err := json.Marshal(parsed.Rows)
	if err != nil {
		return nil, err
	}
	optionsJSON, err := json.Marshal(orEmpty(options))
	if err != nil {
		return nil, err
	}

	job := &entity.ImportJob{
		OrganizationID: orgID,
		ModuleID:       moduleID,
		Filename:       filename,
		Status:         entity.StatusPending,
		Mapping:        mappingJSON,
		Options:        optionsJSON,
		SourceRows:     rowsJSON,
		TotalRows:      len(parsed.Rows),
		CreatedBy:      &userID,
	}
	if err := s.repo.Create(ctx, job); err != nil {
		return nil, err
	}

	msg := jobs.Job{
		Type:   jobs.JobImportProcess,
		UserID: userID,
		Payload: map[string]interface{}{
			"import_id": job.ID,
			"org_id":    orgID,
		},
	}
	if err := s.enqueuer.Publish(ctx, msg); err != nil {
		failJSON, _ := json.Marshal([]entity.RowError{{Message: "failed to enqueue: " + err.Error()}})
		_ = s.repo.Finish(ctx, job.ID, entity.StatusFailed, 0, 0, 0, failJSON)
		return nil, err
	}

	resp := toResponse(job)
	return &resp, nil
}

func (s *Service) Get(ctx context.Context, orgID, id string) (*dto.ImportResponse, error) {
	job, err := s.repo.GetByID(ctx, orgID, id)
	if err != nil {
		return nil, err
	}
	if job == nil {
		return nil, ErrNotFound
	}
	resp := toResponse(job)
	return &resp, nil
}

func (s *Service) List(ctx context.Context, orgID, moduleID string, q dto.ListQuery) (*dto.ListResult, error) {
	if err := s.ensureDynamicModule(ctx, orgID, moduleID); err != nil {
		return nil, err
	}

	if q.Page < 1 {
		q.Page = 1
	}
	if q.PageSize < 1 {
		q.PageSize = DefaultPageSize
	}
	if q.PageSize > MaxPageSize {
		q.PageSize = MaxPageSize
	}

	items, total, err := s.repo.List(ctx, orgID, moduleID, q)
	if err != nil {
		return nil, err
	}

	responses := make([]dto.ImportResponse, 0, len(items))
	for i := range items {
		responses = append(responses, toResponse(&items[i]))
	}

	return &dto.ListResult{
		Imports:    responses,
		Page:       q.Page,
		PageSize:   q.PageSize,
		Total:      total,
		TotalPages: int(math.Max(1, math.Ceil(float64(total)/float64(q.PageSize)))),
	}, nil
}

func toResponse(j *entity.ImportJob) dto.ImportResponse {
	mapping := map[string]string{}
	if len(j.Mapping) > 0 {
		_ = json.Unmarshal(j.Mapping, &mapping)
	}
	errs := []dto.RowErrorDTO{}
	if len(j.Errors) > 0 {
		_ = json.Unmarshal(j.Errors, &errs)
	}
	return dto.ImportResponse{
		ID:            j.ID,
		ModuleID:      j.ModuleID,
		Filename:      j.Filename,
		Status:        j.Status,
		Mapping:       mapping,
		TotalRows:     j.TotalRows,
		ProcessedRows: j.ProcessedRows,
		SuccessRows:   j.SuccessRows,
		ErrorRows:     j.ErrorRows,
		Errors:        errs,
		CreatedBy:     j.CreatedBy,
		StartedAt:     j.StartedAt,
		FinishedAt:    j.FinishedAt,
		CreatedAt:     j.CreatedAt,
		UpdatedAt:     j.UpdatedAt,
	}
}

func orEmpty(m map[string]any) map[string]any {
	if m == nil {
		return map[string]any{}
	}
	return m
}
