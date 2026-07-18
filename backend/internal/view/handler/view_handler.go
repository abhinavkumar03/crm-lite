package handler

import (
	"errors"

	"github.com/gin-gonic/gin"

	"github.com/abhinavkumar03/crm-lite/backend/internal/shared/response"
	"github.com/abhinavkumar03/crm-lite/backend/internal/shared/validation"
	"github.com/abhinavkumar03/crm-lite/backend/internal/tenant"
	"github.com/abhinavkumar03/crm-lite/backend/internal/view/dto"
	"github.com/abhinavkumar03/crm-lite/backend/internal/view/service"
)

const (
	paramModuleID = "id"
	paramViewID   = "viewId"
)

type ViewHandler struct {
	service *service.Service
}

func New(service *service.Service) *ViewHandler {
	return &ViewHandler{service: service}
}

func userID(c *gin.Context) string {
	return c.GetString("userID")
}

func (h *ViewHandler) List(c *gin.Context) {
	views, err := h.service.List(c.Request.Context(), tenant.OrgID(c), c.Param(paramModuleID), userID(c))
	if err != nil {
		h.writeServiceError(c, err, "Unable to fetch views")
		return
	}
	response.OK(c, "Views fetched successfully", views)
}

func (h *ViewHandler) GetByID(c *gin.Context) {
	view, err := h.service.GetByID(c.Request.Context(), tenant.OrgID(c), c.Param(paramModuleID), c.Param(paramViewID), userID(c))
	if err != nil {
		h.writeServiceError(c, err, "Unable to fetch view")
		return
	}
	if view == nil {
		response.NotFound(c, "View not found")
		return
	}
	response.OK(c, "View fetched successfully", view)
}

func (h *ViewHandler) Create(c *gin.Context) {
	var req dto.CreateViewRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request body", nil)
		return
	}
	if err := validation.ValidateStruct(&req); err != nil {
		response.BadRequest(c, "Validation failed", validation.FormatErrors(err))
		return
	}

	view, err := h.service.Create(c.Request.Context(), tenant.OrgID(c), c.Param(paramModuleID), userID(c), req)
	if err != nil {
		h.writeServiceError(c, err, "Unable to create view")
		return
	}
	response.Created(c, "View created successfully", view)
}

func (h *ViewHandler) Update(c *gin.Context) {
	var req dto.UpdateViewRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request body", nil)
		return
	}
	if err := validation.ValidateStruct(&req); err != nil {
		response.BadRequest(c, "Validation failed", validation.FormatErrors(err))
		return
	}

	view, err := h.service.Update(c.Request.Context(), tenant.OrgID(c), c.Param(paramModuleID), c.Param(paramViewID), userID(c), req)
	if err != nil {
		h.writeServiceError(c, err, "Unable to update view")
		return
	}
	response.OK(c, "View updated successfully", view)
}

func (h *ViewHandler) SetDefault(c *gin.Context) {
	if err := h.service.SetDefault(c.Request.Context(), tenant.OrgID(c), c.Param(paramModuleID), c.Param(paramViewID)); err != nil {
		h.writeServiceError(c, err, "Unable to set default view")
		return
	}
	response.OK(c, "Default view updated successfully", nil)
}

func (h *ViewHandler) Delete(c *gin.Context) {
	if err := h.service.Delete(c.Request.Context(), tenant.OrgID(c), c.Param(paramModuleID), c.Param(paramViewID), userID(c)); err != nil {
		h.writeServiceError(c, err, "Unable to delete view")
		return
	}
	response.OK(c, "View deleted successfully", nil)
}

func (h *ViewHandler) writeServiceError(c *gin.Context, err error, fallback string) {
	switch {
	case errors.Is(err, service.ErrModuleNotFound):
		response.NotFound(c, "Module not found")
	case errors.Is(err, service.ErrNotFound):
		response.NotFound(c, "View not found")
	case errors.Is(err, service.ErrForbidden):
		response.Forbidden(c, err.Error())
	default:
		response.InternalServerError(c, fallback)
	}
}
