package handler

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"

	"github.com/abhinavkumar03/crm-lite/backend/internal/exporter/dto"
	"github.com/abhinavkumar03/crm-lite/backend/internal/exporter/service"
	"github.com/abhinavkumar03/crm-lite/backend/internal/exporter/writer"
	recorddto "github.com/abhinavkumar03/crm-lite/backend/internal/record/dto"
	"github.com/abhinavkumar03/crm-lite/backend/internal/shared/response"
	"github.com/abhinavkumar03/crm-lite/backend/internal/shared/validation"
	"github.com/abhinavkumar03/crm-lite/backend/internal/tenant"
)

const (
	paramModuleID   = "id"
	paramExportID   = "exportId"
	paramTemplateID = "templateId"
)

type ExportHandler struct {
	service *service.Service
}

func New(service *service.Service) *ExportHandler {
	return &ExportHandler{service: service}
}

func userID(c *gin.Context) string { return c.GetString("userID") }

// Download (sync) builds and streams the file immediately.
func (h *ExportHandler) ExportNow(c *gin.Context) {
	spec := parseSpecFromQuery(c)
	filename, contentType, content, err := h.service.ExportNow(
		c.Request.Context(), tenant.OrgID(c), c.Param(paramModuleID), spec,
	)
	if err != nil {
		h.writeError(c, err, "Unable to export")
		return
	}
	stream(c, filename, contentType, content)
}

// Create (async) enqueues an export job for the worker.
func (h *ExportHandler) Create(c *gin.Context) {
	var spec dto.ExportSpec
	if err := c.ShouldBindJSON(&spec); err != nil {
		response.BadRequest(c, "Invalid request body", nil)
		return
	}

	job, err := h.service.CreateAsync(
		c.Request.Context(), tenant.OrgID(c), c.Param(paramModuleID), userID(c), spec,
	)
	if err != nil {
		h.writeError(c, err, "Unable to start export")
		return
	}
	response.Created(c, "Export started successfully", job)
}

func (h *ExportHandler) List(c *gin.Context) {
	page, _ := strconv.Atoi(c.Query("page"))
	pageSize, _ := strconv.Atoi(c.Query("page_size"))

	q := dto.ListQuery{Page: page, PageSize: pageSize, Status: c.Query("status")}

	result, err := h.service.List(c.Request.Context(), tenant.OrgID(c), c.Param(paramModuleID), q)
	if err != nil {
		h.writeError(c, err, "Unable to fetch exports")
		return
	}
	response.OK(c, "Exports fetched successfully", result)
}

func (h *ExportHandler) Get(c *gin.Context) {
	job, err := h.service.Get(c.Request.Context(), tenant.OrgID(c), c.Param(paramExportID))
	if err != nil {
		h.writeError(c, err, "Unable to fetch export")
		return
	}
	response.OK(c, "Export fetched successfully", job)
}

func (h *ExportHandler) Download(c *gin.Context) {
	filename, contentType, content, err := h.service.Download(
		c.Request.Context(), tenant.OrgID(c), c.Param(paramExportID),
	)
	if err != nil {
		h.writeError(c, err, "Unable to download export")
		return
	}
	stream(c, filename, contentType, content)
}

// --- Templates -------------------------------------------------------------

func (h *ExportHandler) ListTemplates(c *gin.Context) {
	items, err := h.service.ListTemplates(c.Request.Context(), tenant.OrgID(c), c.Param(paramModuleID))
	if err != nil {
		h.writeError(c, err, "Unable to fetch templates")
		return
	}
	response.OK(c, "Templates fetched successfully", items)
}

func (h *ExportHandler) CreateTemplate(c *gin.Context) {
	var req dto.CreateTemplateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request body", nil)
		return
	}
	if err := validation.ValidateStruct(&req); err != nil {
		response.BadRequest(c, "Validation failed", validation.FormatErrors(err))
		return
	}

	tpl, err := h.service.CreateTemplate(c.Request.Context(), tenant.OrgID(c), c.Param(paramModuleID), userID(c), req)
	if err != nil {
		h.writeError(c, err, "Unable to create template")
		return
	}
	response.Created(c, "Template created successfully", tpl)
}

func (h *ExportHandler) UpdateTemplate(c *gin.Context) {
	var req dto.UpdateTemplateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request body", nil)
		return
	}
	if err := validation.ValidateStruct(&req); err != nil {
		response.BadRequest(c, "Validation failed", validation.FormatErrors(err))
		return
	}

	tpl, err := h.service.UpdateTemplate(
		c.Request.Context(), tenant.OrgID(c), c.Param(paramModuleID), c.Param(paramTemplateID), req,
	)
	if err != nil {
		h.writeError(c, err, "Unable to update template")
		return
	}
	response.OK(c, "Template updated successfully", tpl)
}

func (h *ExportHandler) DeleteTemplate(c *gin.Context) {
	if err := h.service.DeleteTemplate(
		c.Request.Context(), tenant.OrgID(c), c.Param(paramModuleID), c.Param(paramTemplateID),
	); err != nil {
		h.writeError(c, err, "Unable to delete template")
		return
	}
	response.OK(c, "Template deleted successfully", nil)
}

// --- helpers ---------------------------------------------------------------

// stream writes the file as an attachment download.
func stream(c *gin.Context, filename, contentType string, content []byte) {
	c.Header("Content-Disposition", "attachment; filename=\""+filename+"\"")
	c.Data(http.StatusOK, contentType, content)
}

// parseSpecFromQuery reads a sync export spec from the query string. Columns are
// a comma-separated list; filters accept a JSON array.
func parseSpecFromQuery(c *gin.Context) dto.ExportSpec {
	spec := dto.ExportSpec{
		Format: c.Query("format"),
		Search: strings.TrimSpace(c.Query("search")),
		Sort:   c.Query("sort"),
		Order:  c.Query("order"),
		Expand: c.Query("expand") == "true",
	}

	if raw := strings.TrimSpace(c.Query("columns")); raw != "" {
		for _, col := range strings.Split(raw, ",") {
			if v := strings.TrimSpace(col); v != "" {
				spec.Columns = append(spec.Columns, v)
			}
		}
	}

	if raw := c.Query("filters"); raw != "" {
		var filters []recorddto.FilterClause
		if err := json.Unmarshal([]byte(raw), &filters); err == nil {
			spec.Filters = filters
		}
	}

	return spec
}

func (h *ExportHandler) writeError(c *gin.Context, err error, fallback string) {
	switch {
	case errors.Is(err, service.ErrModuleNotFound):
		response.NotFound(c, "Module not found")
	case errors.Is(err, service.ErrNotDynamic):
		response.BadRequest(c, "This module does not support exports", nil)
	case errors.Is(err, service.ErrNotFound):
		response.NotFound(c, "Export not found")
	case errors.Is(err, service.ErrNoColumns):
		response.BadRequest(c, "Select at least one column to export", nil)
	case errors.Is(err, service.ErrNotReady):
		response.Conflict(c, "Export is still processing", nil)
	case errors.Is(err, writer.ErrUnsupported):
		response.BadRequest(c, "Unsupported export format (use csv or xlsx)", nil)
	default:
		response.InternalServerError(c, fallback)
	}
}
