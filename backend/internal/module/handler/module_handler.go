package handler

import (
	"errors"

	"github.com/gin-gonic/gin"

	"github.com/abhinavkumar03/crm-lite/backend/internal/module/dto"
	"github.com/abhinavkumar03/crm-lite/backend/internal/module/service"
	"github.com/abhinavkumar03/crm-lite/backend/internal/shared/response"
	"github.com/abhinavkumar03/crm-lite/backend/internal/shared/validation"
	"github.com/abhinavkumar03/crm-lite/backend/internal/tenant"
)

type ModuleHandler struct {
	service *service.Service
}

func New(service *service.Service) *ModuleHandler {
	return &ModuleHandler{service: service}
}

func (h *ModuleHandler) List(c *gin.Context) {
	modules, err := h.service.List(c.Request.Context(), tenant.OrgID(c))
	if err != nil {
		response.InternalServerError(c, "Unable to fetch modules")
		return
	}
	response.OK(c, "Modules fetched successfully", modules)
}

func (h *ModuleHandler) Navigation(c *gin.Context) {
	nav, err := h.service.Navigation(c.Request.Context(), tenant.OrgID(c))
	if err != nil {
		response.InternalServerError(c, "Unable to fetch navigation")
		return
	}
	response.OK(c, "Navigation fetched successfully", nav)
}

func (h *ModuleHandler) GetByID(c *gin.Context) {
	module, err := h.service.GetByID(c.Request.Context(), tenant.OrgID(c), c.Param("id"))
	if err != nil {
		response.InternalServerError(c, "Unable to fetch module")
		return
	}
	if module == nil {
		response.NotFound(c, "Module not found")
		return
	}
	response.OK(c, "Module fetched successfully", module)
}

func (h *ModuleHandler) Create(c *gin.Context) {
	var req dto.CreateModuleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request body", nil)
		return
	}
	if err := validation.ValidateStruct(&req); err != nil {
		response.BadRequest(c, "Validation failed", validation.FormatErrors(err))
		return
	}

	module, err := h.service.Create(c.Request.Context(), tenant.OrgID(c), req)
	if err != nil {
		h.writeServiceError(c, err, "Unable to create module")
		return
	}
	response.Created(c, "Module created successfully", module)
}

func (h *ModuleHandler) Update(c *gin.Context) {
	var req dto.UpdateModuleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request body", nil)
		return
	}
	if err := validation.ValidateStruct(&req); err != nil {
		response.BadRequest(c, "Validation failed", validation.FormatErrors(err))
		return
	}

	module, err := h.service.Update(c.Request.Context(), tenant.OrgID(c), c.Param("id"), req)
	if err != nil {
		h.writeServiceError(c, err, "Unable to update module")
		return
	}
	response.OK(c, "Module updated successfully", module)
}

func (h *ModuleHandler) SetStatus(c *gin.Context) {
	var req dto.SetStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request body", nil)
		return
	}
	if err := validation.ValidateStruct(&req); err != nil {
		response.BadRequest(c, "Validation failed", validation.FormatErrors(err))
		return
	}

	if err := h.service.SetEnabled(c.Request.Context(), tenant.OrgID(c), c.Param("id"), *req.Enabled); err != nil {
		h.writeServiceError(c, err, "Unable to update module status")
		return
	}
	response.OK(c, "Module status updated successfully", nil)
}

func (h *ModuleHandler) Reorder(c *gin.Context) {
	var req dto.ReorderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request body", nil)
		return
	}
	if err := validation.ValidateStruct(&req); err != nil {
		response.BadRequest(c, "Validation failed", validation.FormatErrors(err))
		return
	}

	if err := h.service.Reorder(c.Request.Context(), tenant.OrgID(c), req.Items); err != nil {
		response.InternalServerError(c, "Unable to reorder modules")
		return
	}
	response.OK(c, "Modules reordered successfully", nil)
}

func (h *ModuleHandler) Delete(c *gin.Context) {
	if err := h.service.Delete(c.Request.Context(), tenant.OrgID(c), c.Param("id")); err != nil {
		h.writeServiceError(c, err, "Unable to delete module")
		return
	}
	response.OK(c, "Module deleted successfully", nil)
}

// writeServiceError maps domain errors to HTTP responses (consistent handling).
func (h *ModuleHandler) writeServiceError(c *gin.Context, err error, fallback string) {
	switch {
	case errors.Is(err, service.ErrNotFound):
		response.NotFound(c, "Module not found")
	case errors.Is(err, service.ErrDuplicateAPIName):
		response.Conflict(c, err.Error(), nil)
	case errors.Is(err, service.ErrSystemModule):
		response.Conflict(c, err.Error(), nil)
	case errors.Is(err, service.ErrInvalidAPIName):
		response.BadRequest(c, err.Error(), nil)
	default:
		response.InternalServerError(c, fallback)
	}
}
