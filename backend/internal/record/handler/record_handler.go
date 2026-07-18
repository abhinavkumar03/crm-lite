package handler

import (
	"encoding/json"
	"errors"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"

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
}

func New(service *service.Service) *RecordHandler {
	return &RecordHandler{service: service}
}

func userID(c *gin.Context) string { return c.GetString("userID") }

func (h *RecordHandler) List(c *gin.Context) {
	q := parseListQuery(c)
	result, err := h.service.List(c.Request.Context(), tenant.OrgID(c), c.Param(paramModuleID), q)
	if err != nil {
		h.writeError(c, err, "Unable to fetch records")
		return
	}
	response.OK(c, "Records fetched successfully", result)
}

func (h *RecordHandler) Get(c *gin.Context) {
	expand := c.Query("expand") == "true"
	rec, err := h.service.Get(c.Request.Context(), tenant.OrgID(c), c.Param(paramModuleID), c.Param(paramRecordID), expand)
	if err != nil {
		h.writeError(c, err, "Unable to fetch record")
		return
	}
	response.OK(c, "Record fetched successfully", rec)
}

func (h *RecordHandler) Create(c *gin.Context) {
	var req dto.CreateRecordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request body", nil)
		return
	}
	if err := validation.ValidateStruct(&req); err != nil {
		response.BadRequest(c, "Validation failed", validation.FormatErrors(err))
		return
	}

	rec, err := h.service.Create(c.Request.Context(), tenant.OrgID(c), c.Param(paramModuleID), userID(c), req)
	if err != nil {
		h.writeError(c, err, "Unable to create record")
		return
	}
	response.Created(c, "Record created successfully", rec)
}

func (h *RecordHandler) Update(c *gin.Context) {
	var req dto.UpdateRecordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request body", nil)
		return
	}
	if err := validation.ValidateStruct(&req); err != nil {
		response.BadRequest(c, "Validation failed", validation.FormatErrors(err))
		return
	}

	rec, err := h.service.Update(c.Request.Context(), tenant.OrgID(c), c.Param(paramModuleID), c.Param(paramRecordID), userID(c), req)
	if err != nil {
		h.writeError(c, err, "Unable to update record")
		return
	}
	response.OK(c, "Record updated successfully", rec)
}

func (h *RecordHandler) Delete(c *gin.Context) {
	if err := h.service.Delete(c.Request.Context(), tenant.OrgID(c), c.Param(paramModuleID), c.Param(paramRecordID)); err != nil {
		h.writeError(c, err, "Unable to delete record")
		return
	}
	response.OK(c, "Record deleted successfully", nil)
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
