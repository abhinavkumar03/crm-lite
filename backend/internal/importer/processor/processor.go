// Package processor executes staged import jobs in the worker. It reuses the
// Phase 7 validation engine (per-row validation) and the Phase 10 record
// repository (persistence), so an import obeys exactly the same rules as a
// single record created through the API.
package processor

import (
	"context"
	"encoding/json"
	"strings"

	"go.uber.org/zap"

	fieldentity "github.com/abhinavkumar03/crm-lite/backend/internal/field/entity"
	importentity "github.com/abhinavkumar03/crm-lite/backend/internal/importer/entity"
	importrepo "github.com/abhinavkumar03/crm-lite/backend/internal/importer/repository"
	recordentity "github.com/abhinavkumar03/crm-lite/backend/internal/record/entity"
	vdto "github.com/abhinavkumar03/crm-lite/backend/internal/validationengine/dto"
)

// maxErrors bounds the persisted error report so a pathological file cannot
// bloat the JSONB column. Counters still reflect the true totals.
const maxErrors = 500

// RecordWriter persists a record (satisfied by the record engine's repository).
type RecordWriter interface {
	Create(ctx context.Context, rec *recordentity.Record) error
}

// FieldReader exposes a module's field metadata (satisfied by the field repo).
type FieldReader interface {
	List(ctx context.Context, orgID, moduleID string) ([]fieldentity.Field, error)
}

// Validator evaluates a payload against a module's schema (the Phase 7 engine).
type Validator interface {
	Validate(ctx context.Context, orgID, moduleID string, data map[string]any) (vdto.ValidateResult, error)
}

type Processor struct {
	repo      *importrepo.Repository
	records   RecordWriter
	fields    FieldReader
	validator Validator
	logger    *zap.Logger
}

func New(
	repo *importrepo.Repository,
	records RecordWriter,
	fields FieldReader,
	validator Validator,
	logger *zap.Logger,
) *Processor {
	return &Processor{
		repo:      repo,
		records:   records,
		fields:    fields,
		validator: validator,
		logger:    logger,
	}
}

// Process maps, validates and inserts each staged row. It always returns nil on
// a run it managed to finish (even with row errors) so asynq does not retry and
// re-insert already-imported rows; a job that cannot be loaded is skipped.
func (p *Processor) Process(ctx context.Context, orgID, id string) error {
	job, err := p.repo.GetByID(ctx, orgID, id)
	if err != nil {
		return err
	}
	if job == nil {
		p.logger.Warn("import: job not found; skipping", zap.String("id", id))
		return nil
	}
	// Idempotency: never re-run a terminal job.
	if job.Status == importentity.StatusCompleted || job.Status == importentity.StatusFailed {
		return nil
	}

	if err := p.repo.MarkProcessing(ctx, id); err != nil {
		return err
	}

	rows, mapping, byAPIName, err := p.load(ctx, orgID, job)
	if err != nil {
		p.finishFailed(ctx, id, "failed to load import: "+err.Error())
		return nil
	}

	var (
		processed, success int
		rowErrors          []importentity.RowError
	)

	for i, row := range rows {
		lineNo := i + 1 // 1-based data row (header excluded)
		processed++

		data := buildData(row, mapping, byAPIName)

		result, verr := p.validator.Validate(ctx, orgID, job.ModuleID, data)
		if verr != nil {
			appendError(&rowErrors, importentity.RowError{Row: lineNo, Message: "validation error: " + verr.Error()})
			continue
		}
		if !result.Valid {
			for _, fe := range result.Errors {
				appendError(&rowErrors, importentity.RowError{Row: lineNo, Field: fe.Field, Message: fe.Message})
			}
			continue
		}

		if err := p.insert(ctx, orgID, job, data); err != nil {
			appendError(&rowErrors, importentity.RowError{Row: lineNo, Message: "failed to save: " + err.Error()})
			continue
		}
		success++
	}

	errorRows := processed - success
	status := importentity.StatusCompleted
	if success == 0 && processed > 0 {
		status = importentity.StatusFailed
	}

	errsJSON, _ := json.Marshal(rowErrors)
	if err := p.repo.Finish(ctx, id, status, processed, success, errorRows, errsJSON); err != nil {
		p.logger.Error("import: finish", zap.Error(err), zap.String("id", id))
	}
	p.logger.Info("import: completed",
		zap.String("id", id),
		zap.String("status", status),
		zap.Int("processed", processed),
		zap.Int("success", success),
		zap.Int("errors", errorRows),
	)
	return nil
}

// load fetches the staged rows, the persisted mapping and the module's fields
// indexed by api_name.
func (p *Processor) load(
	ctx context.Context,
	orgID string,
	job *importentity.ImportJob,
) ([]map[string]string, map[string]string, map[string]fieldentity.Field, error) {
	raw, err := p.repo.SourceRows(ctx, job.ID)
	if err != nil {
		return nil, nil, nil, err
	}
	var rows []map[string]string
	if len(raw) > 0 {
		if err := json.Unmarshal(raw, &rows); err != nil {
			return nil, nil, nil, err
		}
	}

	mapping := map[string]string{}
	if len(job.Mapping) > 0 {
		if err := json.Unmarshal(job.Mapping, &mapping); err != nil {
			return nil, nil, nil, err
		}
	}

	fields, err := p.fields.List(ctx, orgID, job.ModuleID)
	if err != nil {
		return nil, nil, nil, err
	}
	byAPIName := make(map[string]fieldentity.Field, len(fields))
	for i := range fields {
		byAPIName[fields[i].APIName] = fields[i]
	}

	return rows, mapping, byAPIName, nil
}

// buildData applies the mapping to one source row, coercing each value to the
// target field's type. Empty cells are omitted so required-field validation
// fires and optional fields stay absent.
func buildData(
	row map[string]string,
	mapping map[string]string,
	byAPIName map[string]fieldentity.Field,
) map[string]any {
	data := map[string]any{}
	for header, apiName := range mapping {
		f, ok := byAPIName[apiName]
		if !ok {
			continue
		}
		raw := strings.TrimSpace(row[header])
		if raw == "" {
			continue
		}
		data[apiName] = coerce(f.FieldType, raw)
	}
	return data
}

func (p *Processor) insert(ctx context.Context, orgID string, job *importentity.ImportJob, data map[string]any) error {
	b, err := json.Marshal(data)
	if err != nil {
		return err
	}
	rec := &recordentity.Record{
		OrganizationID: orgID,
		ModuleID:       job.ModuleID,
		Data:           b,
		OwnerID:        job.CreatedBy,
		CreatedBy:      job.CreatedBy,
		UpdatedBy:      job.CreatedBy,
	}
	return p.records.Create(ctx, rec)
}

func (p *Processor) finishFailed(ctx context.Context, id, message string) {
	errsJSON, _ := json.Marshal([]importentity.RowError{{Message: message}})
	if err := p.repo.Finish(ctx, id, importentity.StatusFailed, 0, 0, 0, errsJSON); err != nil {
		p.logger.Error("import: finish failed", zap.Error(err), zap.String("id", id))
	}
}

func appendError(errs *[]importentity.RowError, e importentity.RowError) {
	if len(*errs) >= maxErrors {
		return
	}
	*errs = append(*errs, e)
}
