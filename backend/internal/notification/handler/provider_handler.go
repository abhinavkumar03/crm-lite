package handler

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/abhinavkumar03/crm-lite/backend/internal/notification/dto"
	"github.com/abhinavkumar03/crm-lite/backend/internal/notification/service"
	"github.com/abhinavkumar03/crm-lite/backend/internal/shared/response"
	"github.com/abhinavkumar03/crm-lite/backend/internal/tenant"
)

type ProviderHandler struct {
	service *service.ProviderService
}

func NewProviderHandler(svc *service.ProviderService) *ProviderHandler {
	return &ProviderHandler{service: svc}
}

func (h *ProviderHandler) List(c *gin.Context) {
	items, err := h.service.List(c.Request.Context(), tenant.OrgID(c), c.Query("channel"))
	if err != nil {
		response.InternalServerError(c, "Unable to list providers")
		return
	}
	response.OK(c, "Providers fetched successfully", items)
}

func (h *ProviderHandler) Get(c *gin.Context) {
	item, err := h.service.Get(c.Request.Context(), tenant.OrgID(c), c.Param("id"))
	if err != nil {
		if errors.Is(err, service.ErrProviderNotFound) {
			response.NotFound(c, "Provider not found")
			return
		}
		response.InternalServerError(c, "Unable to fetch provider")
		return
	}
	response.OK(c, "Provider fetched successfully", item)
}

func (h *ProviderHandler) Create(c *gin.Context) {
	var req dto.ProviderUpsertRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error(), nil)
		return
	}
	item, err := h.service.Create(c.Request.Context(), tenant.OrgID(c), userID(c), req)
	if err != nil {
		if errors.Is(err, service.ErrInvalidProvider) {
			response.BadRequest(c, err.Error(), nil)
			return
		}
		response.InternalServerError(c, "Unable to create provider")
		return
	}
	response.Created(c, "Provider created successfully", item)
}

func (h *ProviderHandler) Update(c *gin.Context) {
	var req dto.ProviderUpsertRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error(), nil)
		return
	}
	item, err := h.service.Update(c.Request.Context(), tenant.OrgID(c), c.Param("id"), req)
	if err != nil {
		if errors.Is(err, service.ErrProviderNotFound) {
			response.NotFound(c, "Provider not found")
			return
		}
		response.InternalServerError(c, "Unable to update provider")
		return
	}
	response.OK(c, "Provider updated successfully", item)
}

func (h *ProviderHandler) Delete(c *gin.Context) {
	if err := h.service.Delete(c.Request.Context(), tenant.OrgID(c), c.Param("id")); err != nil {
		response.InternalServerError(c, "Unable to delete provider")
		return
	}
	response.OK(c, "Provider deleted successfully", nil)
}

func (h *ProviderHandler) Test(c *gin.Context) {
	var req dto.ProviderTestRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error(), nil)
		return
	}
	if err := h.service.Test(c.Request.Context(), tenant.OrgID(c), c.Param("id"), req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "Provider test failed: " + err.Error()})
		return
	}
	response.OK(c, "Provider test succeeded", nil)
}

func (h *ProviderHandler) ListSenders(c *gin.Context) {
	items, err := h.service.ListSenders(c.Request.Context(), tenant.OrgID(c), c.Query("channel"))
	if err != nil {
		response.InternalServerError(c, "Unable to list sender identities")
		return
	}
	response.OK(c, "Sender identities fetched successfully", items)
}

func (h *ProviderHandler) CreateSender(c *gin.Context) {
	var req dto.SenderUpsertRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error(), nil)
		return
	}
	item, err := h.service.CreateSender(c.Request.Context(), tenant.OrgID(c), req)
	if err != nil {
		response.InternalServerError(c, "Unable to create sender identity")
		return
	}
	response.Created(c, "Sender identity created successfully", item)
}

func (h *ProviderHandler) DeleteSender(c *gin.Context) {
	if err := h.service.DeleteSender(c.Request.Context(), tenant.OrgID(c), c.Param("id")); err != nil {
		response.InternalServerError(c, "Unable to delete sender identity")
		return
	}
	response.OK(c, "Sender identity deleted successfully", nil)
}
