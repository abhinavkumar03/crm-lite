package handler

import (
	"encoding/json"
	"errors"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"

	"github.com/abhinavkumar03/crm-lite/backend/internal/rbac"
	"github.com/abhinavkumar03/crm-lite/backend/internal/record/dto"
	"github.com/abhinavkumar03/crm-lite/backend/internal/record/service"
	"github.com/abhinavkumar03/crm-lite/backend/internal/shared/response"
	"github.com/abhinavkumar03/crm-lite/backend/internal/shared/validation"
	"github.com/abhinavkumar03/crm-lite/backend/internal/tenant"
)

const (
	paramModuleID = "id"
	paramRecordID = "recordId"
)

type RecordHandler struct {
	service *service.Service
	guard   *rbac.Guard
}

func New(service *service.Service, guard *rbac.Guard) *RecordHandler {
	return &RecordHandler{service: service, guard: guard}
}

func userID(c *gin.Context) string { return c.GetString("userID") }

func (h *RecordHandler) List(c *gin.Context) {
	moduleID := c.Param(paramModuleID)
	q := parseListQuery(c)
	result, err := h.service.List(c.Request.Context(), tenant.OrgID(c), moduleID, q)
	if err != nil {
		h.writeError(c, err, "Unable to fetch records")
		return
	}
	access := h.fieldAccess(c, moduleID)
	for i := range result.Records {
		result.Records[i] = h.stripHidden(result.Records[i], access)
	}
	response.OK(c, "Records fetched successfully", result)
}

func (h *RecordHandler) Get(c *gin.Context) {
	moduleID := c.Param(paramModuleID)
	expand := c.Query("expand") == "true"
	rec, err := h.service.Get(c.Request.Context(), tenant.OrgID(c), moduleID, c.Param(paramRecordID), expand)
	if err != nil {
		h.writeError(c, err, "Unable to fetch record")
		return
	}
	stripped := h.stripHidden(*rec, h.fieldAccess(c, moduleID))
	response.OK(c, "Record fetched successfully", stripped)
}

func (h *RecordHandler) Create(c *gin.Context) {
	moduleID := c.Param(paramModuleID)
	var req dto.CreateRecordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request body", nil)
		return
	}
	if err := validation.ValidateStruct(&req); err != nil {
		response.BadRequest(c, "Validation failed", validation.FormatErrors(err))
		return
	}

	req.Data = h.stripNonWritable(req.Data, h.fieldAccess(c, moduleID))

	rec, err := h.service.Create(c.Request.Context(), tenant.OrgID(c), moduleID, userID(c), req)
	if err != nil {
		h.writeError(c, err, "Unable to create record")
		return
	}
	stripped := h.stripHidden(*rec, h.fieldAccess(c, moduleID))
	response.Created(c, "Record created successfully", stripped)
}

func (h *RecordHandler) Update(c *gin.Context) {
	moduleID := c.Param(paramModuleID)
	var req dto.UpdateRecordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request body", nil)
		return
	}
	if err := validation.ValidateStruct(&req); err != nil {
		response.BadRequest(c, "Validation failed", validation.FormatErrors(err))
		return
	}

	req.Data = h.stripNonWritable(req.Data, h.fieldAccess(c, moduleID))

	rec, err := h.service.Update(c.Request.Context(), tenant.OrgID(c), moduleID, c.Param(paramRecordID), userID(c), req)
	if err != nil {
		h.writeError(c, err, "Unable to update record")
		return
	}
	stripped := h.stripHidden(*rec, h.fieldAccess(c, moduleID))
	response.OK(c, "Record updated successfully", stripped)
}

func (h *RecordHandler) Delete(c *gin.Context) {
	if err := h.service.Delete(c.Request.Context(), tenant.OrgID(c), c.Param(paramModuleID), c.Param(paramRecordID)); err != nil {
		h.writeError(c, err, "Unable to delete record")
		return
	}
	response.OK(c, "Record deleted successfully", nil)
}

// fieldAccess returns api_name → access for the caller's role. Empty map = full write.
func (h *RecordHandler) fieldAccess(c *gin.Context, moduleID string) map[string]string {
	if h.guard == nil {
		return nil
	}
	access, err := h.guard.FieldAccessByAPIName(c.Request.Context(), tenant.RoleID(c), moduleID)
	if err != nil {
		return nil
	}
	return access
}

func (h *RecordHandler) stripHidden(rec dto.RecordResponse, access map[string]string) dto.RecordResponse {
	if len(access) == 0 || rec.Data == nil {
		return rec
	}
	data := make(map[string]any, len(rec.Data))
	for k, v := range rec.Data {
		if access[k] == rbac.FieldHidden {
			continue
		}
		data[k] = v
	}
	rec.Data = data
	if rec.Relations != nil {
		rels := make(map[string]dto.RelationRef, len(rec.Relations))
		for k, v := range rec.Relations {
			if access[k] == rbac.FieldHidden {
				continue
			}
			rels[k] = v
		}
		rec.Relations = rels
	}
	return rec
}

// stripNonWritable drops hidden and read-only keys from an incoming payload so
// a caller cannot overwrite fields their role cannot write.
func (h *RecordHandler) stripNonWritable(data map[string]any, access map[string]string) map[string]any {
	if len(access) == 0 || data == nil {
		return data
	}
	out := make(map[string]any, len(data))
	for k, v := range data {
		level := access[k]
		if level == rbac.FieldHidden || level == rbac.FieldRead {
			continue
		}
		out[k] = v
	}
	return out
}

// parseListQuery reads pagination, search, sort and filters from the query
// string. Filters accept either a JSON array (?filters=[...]) or simple
// equality shorthands (?filter.<field>=<value>).
func parseListQuery(c *gin.Context) dto.ListQuery {
	page, _ := strconv.Atoi(c.Query("page"))
	pageSize, _ := strconv.Atoi(c.Query("page_size"))

	q := dto.ListQuery{
		Page:     page,
		PageSize: pageSize,
		Search:   strings.TrimSpace(c.Query("search")),
		Sort:     c.Query("sort"),
		Order:    c.Query("order"),
		Expand:   c.Query("expand") == "true",
	}

	if raw := c.Query("filters"); raw != "" {
		var filters []dto.FilterClause
		if err := json.Unmarshal([]byte(raw), &filters); err == nil {
			q.Filters = filters
		}
	}

	for key, values := range c.Request.URL.Query() {
		if strings.HasPrefix(key, "filter.") && len(values) > 0 {
			q.Filters = append(q.Filters, dto.FilterClause{
				Field:    strings.TrimPrefix(key, "filter."),
				Operator: dto.OpEquals,
				Value:    values[0],
			})
		}
	}

	return q
}

func (h *RecordHandler) writeError(c *gin.Context, err error, fallback string) {
	var verr *service.ValidationError
	switch {
	case errors.As(err, &verr):
		response.BadRequest(c, "Validation failed", verr.Errors)
	case errors.Is(err, service.ErrModuleNotFound):
		response.NotFound(c, "Module not found")
	case errors.Is(err, service.ErrNotDynamic):
		response.BadRequest(c, "This module does not support the record runtime", nil)
	case errors.Is(err, service.ErrNotFound):
		response.NotFound(c, "Record not found")
	default:
		response.InternalServerError(c, fallback)
	}
}
