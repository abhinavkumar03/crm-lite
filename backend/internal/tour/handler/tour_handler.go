package handler

import (
	"github.com/gin-gonic/gin"

	"github.com/abhinavkumar03/crm-lite/backend/internal/shared/response"
	"github.com/abhinavkumar03/crm-lite/backend/internal/shared/validation"
	"github.com/abhinavkumar03/crm-lite/backend/internal/tenant"
	"github.com/abhinavkumar03/crm-lite/backend/internal/tour/dto"
	"github.com/abhinavkumar03/crm-lite/backend/internal/tour/service"
)

type TourHandler struct {
	service *service.Service
}

func New(service *service.Service) *TourHandler {
	return &TourHandler{service: service}
}

func userID(c *gin.Context) string { return c.GetString("userID") }

// Get returns the current user's tour progress. A brand-new user gets a default
// "active" state so the client can start immediately.
func (h *TourHandler) Get(c *gin.Context) {
	progress, err := h.service.Get(c.Request.Context(), tenant.OrgID(c), userID(c), c.Query("key"))
	if err != nil {
		response.InternalServerError(c, "Unable to fetch tour progress")
		return
	}
	response.OK(c, "Tour progress fetched successfully", progress)
}

// Update advances, completes or skips the tour. It is the single write path.
func (h *TourHandler) Update(c *gin.Context) {
	var req dto.UpdateProgressRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request body", nil)
		return
	}
	if err := validation.ValidateStruct(&req); err != nil {
		response.BadRequest(c, "Validation failed", validation.FormatErrors(err))
		return
	}

	progress, err := h.service.Save(c.Request.Context(), tenant.OrgID(c), userID(c), req)
	if err != nil {
		response.InternalServerError(c, "Unable to save tour progress")
		return
	}
	response.OK(c, "Tour progress saved successfully", progress)
}

// Restart resets the tour back to the first step.
func (h *TourHandler) Restart(c *gin.Context) {
	var req dto.RestartRequest
	// A body is optional here; ignore bind errors and fall back to the default
	// tour key.
	_ = c.ShouldBindJSON(&req)
	if err := validation.ValidateStruct(&req); err != nil {
		response.BadRequest(c, "Validation failed", validation.FormatErrors(err))
		return
	}

	progress, err := h.service.Restart(c.Request.Context(), tenant.OrgID(c), userID(c), req)
	if err != nil {
		response.InternalServerError(c, "Unable to restart tour")
		return
	}
	response.OK(c, "Tour restarted successfully", progress)
}
