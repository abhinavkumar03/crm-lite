package handler

import (
	"errors"

	"github.com/gin-gonic/gin"

	"github.com/abhinavkumar03/crm-lite/backend/internal/field/dto"
	"github.com/abhinavkumar03/crm-lite/backend/internal/field/service"
	"github.com/abhinavkumar03/crm-lite/backend/internal/rbac"
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
	guard   *rbac.Guard
}

func New(service *service.Service, guard *rbac.Guard) *FieldHandler {
	return &FieldHandler{service: service, guard: guard}
}

func (h *FieldHandler) List(c *gin.Context) {
	moduleID := c.Param(paramModuleID)
	fields, err := h.service.List(c.Request.Context(), tenant.OrgID(c), moduleID)
	if err != nil {
		h.writeServiceError(c, err, "Unable to fetch fields")
		return
	}
	fields = h.filterHidden(c, moduleID, fields)
	response.OK(c, "Fields fetched successfully", fields)
}

func (h *FieldHandler) GetByID(c *gin.Context) {
	moduleID := c.Param(paramModuleID)
	field, err := h.service.GetByID(c.Request.Context(), tenant.OrgID(c), moduleID, c.Param(paramFieldID))
	if err != nil {
		h.writeServiceError(c, err, "Unable to fetch field")
		return
	}
	if field == nil {
		response.NotFound(c, "Field not found")
		return
	}
	// Treat a hidden field as not found for callers without field.manage.
	if h.isHidden(c, moduleID, field.ID) {
		response.NotFound(c, "Field not found")
		return
	}
	response.OK(c, "Field fetched successfully", field)
}

// filterHidden removes fields the caller's role cannot see. Callers with
// field.manage (settings admins) always see the full catalog.
func (h *FieldHandler) filterHidden(c *gin.Context, moduleID string, fields []dto.FieldResponse) []dto.FieldResponse {
	if h.guard == nil || rbac.Has(c, rbac.PermFieldManage) {
		return fields
	}
	access, err := h.guard.FieldAccessMap(c.Request.Context(), tenant.RoleID(c), moduleID)
	if err != nil || len(access) == 0 {
		return fields
	}
	out := make([]dto.FieldResponse, 0, len(fields))
	for _, f := range fields {
		if access[f.ID] == rbac.FieldHidden {
			continue
		}
		out = append(out, f)
	}
	return out
}

func (h *FieldHandler) isHidden(c *gin.Context, moduleID, fieldID string) bool {
	if h.guard == nil || rbac.Has(c, rbac.PermFieldManage) {
		return false
	}
	access, err := h.guard.FieldAccessMap(c.Request.Context(), tenant.RoleID(c), moduleID)
	if err != nil {
		return false
	}
	return access[fieldID] == rbac.FieldHidden
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
