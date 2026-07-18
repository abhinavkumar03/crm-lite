package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"time"

	"github.com/hibiken/asynq"

	"github.com/abhinavkumar03/crm-lite/backend/internal/exporter/dto"
	"github.com/abhinavkumar03/crm-lite/backend/internal/exporter/entity"
	"github.com/abhinavkumar03/crm-lite/backend/internal/exporter/repository"
	"github.com/abhinavkumar03/crm-lite/backend/internal/exporter/writer"
	fieldentity "github.com/abhinavkumar03/crm-lite/backend/internal/field/entity"
	"github.com/abhinavkumar03/crm-lite/backend/internal/jobs"
	recorddto "github.com/abhinavkumar03/crm-lite/backend/internal/record/dto"
)

const (
	DefaultPageSize = 20
	MaxPageSize     = 100

	// MaxRows bounds a single export so an inline-stored file stays reasonable.
	MaxRows = 20000
	// fetchPageSize is how many records are pulled per record-runtime page.
	fetchPageSize = 100

	storageDynamic = "dynamic"
)

var (
	ErrModuleNotFound = errors.New("module not found")
	ErrNotDynamic     = errors.New("module does not use dynamic record storage")
	ErrNotFound       = errors.New("export job not found")
	ErrNoColumns      = errors.New("no exportable columns were resolved")
	ErrNotReady       = errors.New("export is not ready for download")
)

// RowFetcher fetches records through the Phase 10 runtime (search + filters +
// sort + relation expansion). Satisfied by the record service.
type RowFetcher interface {
	List(ctx context.Context, orgID, moduleID, userID string, q recorddto.ListQuery) (*recorddto.ListResult, error)
}

// FieldReader exposes module + field metadata (satisfied by the field repo).
type FieldReader interface {
	ModuleStorage(ctx context.Context, orgID, moduleID string) (string, bool, error)
	List(ctx context.Context, orgID, moduleID string) ([]fieldentity.Field, error)
}

// Enqueuer publishes jobs onto the async queue (satisfied by *jobs.Producer).
type Enqueuer interface {
	Publish(ctx context.Context, job jobs.Job, opts ...asynq.Option) error
}

type Service struct {
	exports   *repository.ExportRepository
	templates *repository.TemplateRepository
	rows      RowFetcher
	fields    FieldReader
	enqueuer  Enqueuer
}

func New(
	exports *repository.ExportRepository,
	templates *repository.TemplateRepository,
	rows RowFetcher,
	fields FieldReader,
	enqueuer Enqueuer,
) *Service {
	return &Service{
		exports:   exports,
		templates: templates,
		rows:      rows,
		fields:    fields,
		enqueuer:  enqueuer,
	}
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

// Build is the generation core shared by synchronous downloads and asynchronous
// jobs: it resolves columns, pages through the record runtime and serializes the
// file. It needs no persistence, which keeps it easy to reuse and to unit test.
func (s *Service) Build(ctx context.Context, orgID, moduleID string, spec dto.ExportSpec) (*writer.Result, int, error) {
	fields, err := s.fields.List(ctx, orgID, moduleID)
	if err != nil {
		return nil, 0, err
	}

	columns := resolveColumns(fields, spec.Columns)
	if len(columns) == 0 {
		return nil, 0, ErrNoColumns
	}

	expand := spec.Expand || needsExpand(columns)

	rows := make([]map[string]any, 0)
	page := 1
	for {
		q := recorddto.ListQuery{
			Page:      page,
			PageSize:  fetchPageSize,
			Search:    spec.Search,
			Sort:      spec.Sort,
			Order:     spec.Order,
			Filters:   spec.Filters,
			Expand:    expand,
			SkipTotal: true, // export only needs "is there another page?"
		}
		res, err := s.rows.List(ctx, orgID, moduleID, "", q)
		if err != nil {
			return nil, 0, err
		}

		rows = append(rows, buildRows(res.Records, columns)...)

		if len(rows) >= MaxRows {
			rows = rows[:MaxRows]
			break
		}
		// Short page (or empty) means we have reached the end — no COUNT needed.
		if len(res.Records) < fetchPageSize {
			break
		}
		page++
	}

	result, err := writer.Write(normalizeFormat(spec.Format), columns, rows)
	if err != nil {
		return nil, 0, err
	}
	return result, len(rows), nil
}

// ExportNow builds and returns the file synchronously (for small, immediate
// downloads). Returns the suggested filename, content type and bytes.
func (s *Service) ExportNow(ctx context.Context, orgID, moduleID string, spec dto.ExportSpec) (string, string, []byte, error) {
	if err := s.ensureDynamicModule(ctx, orgID, moduleID); err != nil {
		return "", "", nil, err
	}
	result, _, err := s.Build(ctx, orgID, moduleID, spec)
	if err != nil {
		return "", "", nil, err
	}
	return buildFilename(result.Ext), result.ContentType, result.Content, nil
}

// CreateAsync persists a pending export job and enqueues it for the worker.
func (s *Service) CreateAsync(ctx context.Context, orgID, moduleID, userID string, spec dto.ExportSpec) (*dto.ExportResponse, error) {
	if err := s.ensureDynamicModule(ctx, orgID, moduleID); err != nil {
		return nil, err
	}

	format := normalizeFormat(spec.Format)
	columnsJSON, _ := json.Marshal(orEmptyStrings(spec.Columns))
	filtersJSON, _ := json.Marshal(orEmptyFilters(spec.Filters))
	optionsJSON, _ := json.Marshal(optionsFromSpec(spec))

	job := &entity.ExportJob{
		OrganizationID: orgID,
		ModuleID:       moduleID,
		Filename:       buildFilename(format),
		Format:         format,
		Status:         entity.StatusPending,
		Columns:        columnsJSON,
		Filters:        filtersJSON,
		Options:        optionsJSON,
		CreatedBy:      &userID,
	}
	if err := s.exports.Create(ctx, job); err != nil {
		return nil, err
	}

	msg := jobs.Job{
		Type:   jobs.JobExportProcess,
		UserID: userID,
		Payload: map[string]interface{}{
			"export_id": job.ID,
			"org_id":    orgID,
		},
	}
	if err := s.enqueuer.Publish(ctx, msg); err != nil {
		_ = s.exports.MarkFailed(ctx, job.ID, "failed to enqueue: "+err.Error())
		return nil, err
	}

	resp := toExportResponse(job)
	return &resp, nil
}

// RunJob executes a persisted export in the worker: reconstruct the spec, build
// the file and store it. It always returns nil once it has recorded a terminal
// state so asynq does not retry a finished export.
func (s *Service) RunJob(ctx context.Context, orgID, id string) error {
	job, err := s.exports.GetByID(ctx, orgID, id)
	if err != nil {
		return err
	}
	if job == nil {
		return nil
	}
	if job.Status == entity.StatusCompleted || job.Status == entity.StatusFailed {
		return nil
	}

	if err := s.exports.MarkProcessing(ctx, id); err != nil {
		return err
	}

	result, rowCount, err := s.Build(ctx, orgID, job.ModuleID, specFromJob(job))
	if err != nil {
		_ = s.exports.MarkFailed(ctx, id, err.Error())
		return nil
	}
	if err := s.exports.Complete(ctx, id, rowCount, result.Content); err != nil {
		return err
	}
	return nil
}

func (s *Service) Get(ctx context.Context, orgID, id string) (*dto.ExportResponse, error) {
	job, err := s.exports.GetByID(ctx, orgID, id)
	if err != nil {
		return nil, err
	}
	if job == nil {
		return nil, ErrNotFound
	}
	resp := toExportResponse(job)
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

	items, total, err := s.exports.List(ctx, orgID, moduleID, q)
	if err != nil {
		return nil, err
	}

	responses := make([]dto.ExportResponse, 0, len(items))
	for i := range items {
		responses = append(responses, toExportResponse(&items[i]))
	}

	return &dto.ListResult{
		Exports:    responses,
		Page:       q.Page,
		PageSize:   q.PageSize,
		Total:      total,
		TotalPages: int(math.Max(1, math.Ceil(float64(total)/float64(q.PageSize)))),
	}, nil
}

// Download returns the stored file for a completed export.
func (s *Service) Download(ctx context.Context, orgID, id string) (filename, contentType string, content []byte, err error) {
	name, format, status, body, found, err := s.exports.GetForDownload(ctx, orgID, id)
	if err != nil {
		return "", "", nil, err
	}
	if !found {
		return "", "", nil, ErrNotFound
	}
	if status != entity.StatusCompleted || body == nil {
		return "", "", nil, ErrNotReady
	}
	return name, contentTypeFor(format), body, nil
}

// --- helpers ---------------------------------------------------------------

func normalizeFormat(format string) string {
	switch format {
	case entity.FormatXLSX:
		return entity.FormatXLSX
	default:
		return entity.FormatCSV
	}
}

func contentTypeFor(format string) string {
	if format == entity.FormatXLSX {
		return "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet"
	}
	return "text/csv"
}

func buildFilename(ext string) string {
	return fmt.Sprintf("export-%s.%s", time.Now().UTC().Format("20060102-150405"), ext)
}

// specOptions is the non-column/filter part of a spec, persisted in the job's
// options JSONB and restored when the worker runs it.
type specOptions struct {
	Search string `json:"search"`
	Sort   string `json:"sort"`
	Order  string `json:"order"`
	Expand bool   `json:"expand"`
}

func optionsFromSpec(spec dto.ExportSpec) specOptions {
	return specOptions{Search: spec.Search, Sort: spec.Sort, Order: spec.Order, Expand: spec.Expand}
}

func specFromJob(job *entity.ExportJob) dto.ExportSpec {
	var columns []string
	_ = json.Unmarshal(job.Columns, &columns)

	var filters []recorddto.FilterClause
	_ = json.Unmarshal(job.Filters, &filters)

	var opts specOptions
	_ = json.Unmarshal(job.Options, &opts)

	return dto.ExportSpec{
		Format:  job.Format,
		Columns: columns,
		Filters: filters,
		Search:  opts.Search,
		Sort:    opts.Sort,
		Order:   opts.Order,
		Expand:  opts.Expand,
	}
}

func toExportResponse(j *entity.ExportJob) dto.ExportResponse {
	columns := []string{}
	if len(j.Columns) > 0 {
		_ = json.Unmarshal(j.Columns, &columns)
	}
	return dto.ExportResponse{
		ID:         j.ID,
		ModuleID:   j.ModuleID,
		Filename:   j.Filename,
		Format:     j.Format,
		Status:     j.Status,
		Columns:    columns,
		RowCount:   j.RowCount,
		ByteSize:   j.ByteSize,
		Error:      j.Error,
		CreatedBy:  j.CreatedBy,
		StartedAt:  j.StartedAt,
		FinishedAt: j.FinishedAt,
		CreatedAt:  j.CreatedAt,
		UpdatedAt:  j.UpdatedAt,
	}
}

func orEmptyStrings(s []string) []string {
	if s == nil {
		return []string{}
	}
	return s
}

func orEmptyFilters(f []recorddto.FilterClause) []recorddto.FilterClause {
	if f == nil {
		return []recorddto.FilterClause{}
	}
	return f
}
