package handler

import (
	"errors"

	"github.com/gin-gonic/gin"

	"github.com/abhinavkumar03/crm-lite/backend/internal/settings/dto"
	"github.com/abhinavkumar03/crm-lite/backend/internal/settings/service"
	"github.com/abhinavkumar03/crm-lite/backend/internal/shared/response"
	"github.com/abhinavkumar03/crm-lite/backend/internal/shared/validation"
	"github.com/abhinavkumar03/crm-lite/backend/internal/tenant"
)

type SettingsHandler struct {
	service *service.Service
}

func New(service *service.Service) *SettingsHandler {
	return &SettingsHandler{service: service}
}

// Get returns the current organization's settings (general + automation), with
// defaults filled in for anything never saved.
func (h *SettingsHandler) Get(c *gin.Context) {
	settings, err := h.service.Get(c.Request.Context(), tenant.OrgID(c))
	if err != nil {
		if errors.Is(err, service.ErrNotFound) {
			response.NotFound(c, "Organization not found")
			return
		}
		response.InternalServerError(c, "Unable to fetch settings")
		return
	}
	response.OK(c, "Settings fetched successfully", settings)
}

// Update applies a partial settings change (any subset of name/general/automation).
func (h *SettingsHandler) Update(c *gin.Context) {
	var req dto.UpdateSettingsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request body", nil)
		return
	}
	if err := validation.ValidateStruct(&req); err != nil {
		response.BadRequest(c, "Validation failed", validation.FormatErrors(err))
		return
	}

	settings, err := h.service.Update(c.Request.Context(), tenant.OrgID(c), req)
	if err != nil {
		if errors.Is(err, service.ErrNotFound) {
			response.NotFound(c, "Organization not found")
			return
		}
		response.InternalServerError(c, "Unable to update settings")
		return
	}
	response.OK(c, "Settings updated successfully", settings)
}
