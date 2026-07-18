package handler

import (
	"errors"

	"github.com/gin-gonic/gin"

	"github.com/abhinavkumar03/crm-lite/backend/internal/field/dto"
	"github.com/abhinavkumar03/crm-lite/backend/internal/field/service"
	"github.com/abhinavkumar03/crm-lite/backend/internal/shared/response"
	"github.com/abhinavkumar03/crm-lite/backend/internal/shared/validation"
	"github.com/abhinavkumar03/crm-lite/backend/internal/tenant"
)

// Route param names. The module id reuses ":id" so the /modules/:id subtree is
// consistent with the module engine's routes.
const (
	paramModuleID = "id"
	paramFieldID  = "fieldId"
)

type FieldHandler struct {
	service *service.Service
}

func New(service *service.Service) *FieldHandler {
	return &FieldHandler{service: service}
}

func (h *FieldHandler) List(c *gin.Context) {
	fields, err := h.service.List(c.Request.Context(), tenant.OrgID(c), c.Param(paramModuleID))
	if err != nil {
		h.writeServiceError(c, err, "Unable to fetch fields")
		return
	}
	response.OK(c, "Fields fetched successfully", fields)
}

func (h *FieldHandler) GetByID(c *gin.Context) {
	field, err := h.service.GetByID(c.Request.Context(), tenant.OrgID(c), c.Param(paramModuleID), c.Param(paramFieldID))
	if err != nil {
		h.writeServiceError(c, err, "Unable to fetch field")
		return
	}
	if field == nil {
		response.NotFound(c, "Field not found")
		return
	}
	response.OK(c, "Field fetched successfully", field)
}

func (h *FieldHandler) Create(c *gin.Context) {
	var req dto.CreateFieldRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request body", nil)
		return
	}
	if err := validation.ValidateStruct(&req); err != nil {
		response.BadRequest(c, "Validation failed", validation.FormatErrors(err))
		return
	}

	field, err := h.service.Create(c.Request.Context(), tenant.OrgID(c), c.Param(paramModuleID), req)
	if err != nil {
		h.writeServiceError(c, err, "Unable to create field")
		return
	}
	response.Created(c, "Field created successfully", field)
}

func (h *FieldHandler) Update(c *gin.Context) {
	var req dto.UpdateFieldRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request body", nil)
		return
	}
	if err := validation.ValidateStruct(&req); err != nil {
		response.BadRequest(c, "Validation failed", validation.FormatErrors(err))
		return
	}

	field, err := h.service.Update(c.Request.Context(), tenant.OrgID(c), c.Param(paramModuleID), c.Param(paramFieldID), req)
	if err != nil {
		h.writeServiceError(c, err, "Unable to update field")
		return
	}
	response.OK(c, "Field updated successfully", field)
}

func (h *FieldHandler) Reorder(c *gin.Context) {
	var req dto.ReorderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request body", nil)
		return
	}
	if err := validation.ValidateStruct(&req); err != nil {
		response.BadRequest(c, "Validation failed", validation.FormatErrors(err))
		return
	}

	if err := h.service.Reorder(c.Request.Context(), tenant.OrgID(c), c.Param(paramModuleID), req.Items); err != nil {
		h.writeServiceError(c, err, "Unable to reorder fields")
		return
	}
	response.OK(c, "Fields reordered successfully", nil)
}

func (h *FieldHandler) Delete(c *gin.Context) {
	if err := h.service.Delete(c.Request.Context(), tenant.OrgID(c), c.Param(paramModuleID), c.Param(paramFieldID)); err != nil {
		h.writeServiceError(c, err, "Unable to delete field")
		return
	}
	response.OK(c, "Field deleted successfully", nil)
}

// writeServiceError maps domain errors to HTTP responses (consistent handling).
func (h *FieldHandler) writeServiceError(c *gin.Context, err error, fallback string) {
	switch {
	case errors.Is(err, service.ErrModuleNotFound):
		response.NotFound(c, "Module not found")
	case errors.Is(err, service.ErrNotFound):
		response.NotFound(c, "Field not found")
	case errors.Is(err, service.ErrDuplicateAPIName):
		response.Conflict(c, err.Error(), nil)
	case errors.Is(err, service.ErrSystemField):
		response.Conflict(c, err.Error(), nil)
	case errors.Is(err, service.ErrInvalidAPIName),
		errors.Is(err, service.ErrInvalidType),
		errors.Is(err, service.ErrOptionsRequired),
		errors.Is(err, service.ErrLookupRequired),
		errors.Is(err, service.ErrInvalidLength):
		response.BadRequest(c, err.Error(), nil)
	default:
		response.InternalServerError(c, fallback)
	}
}
