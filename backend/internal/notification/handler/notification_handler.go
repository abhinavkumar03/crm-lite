package handler

import (
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/abhinavkumar03/crm-lite/backend/internal/notification/dto"
	"github.com/abhinavkumar03/crm-lite/backend/internal/notification/service"
	"github.com/abhinavkumar03/crm-lite/backend/internal/shared/response"
	"github.com/abhinavkumar03/crm-lite/backend/internal/shared/validation"
	"github.com/abhinavkumar03/crm-lite/backend/internal/tenant"
)

type NotificationHandler struct {
	service *service.Service
}

func New(service *service.Service) *NotificationHandler {
	return &NotificationHandler{service: service}
}

func userID(c *gin.Context) string { return c.GetString("userID") }

func (h *NotificationHandler) Send(c *gin.Context) {
	var req dto.SendNotificationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request body", nil)
		return
	}
	if err := validation.ValidateStruct(&req); err != nil {
		response.BadRequest(c, "Validation failed", validation.FormatErrors(err))
		return
	}

	notification, err := h.service.Send(c.Request.Context(), tenant.OrgID(c), userID(c), req)
	if err != nil {
		response.InternalServerError(c, "Unable to queue notification")
		return
	}
	response.Created(c, "Notification queued successfully", notification)
}

func (h *NotificationHandler) List(c *gin.Context) {
	page, _ := strconv.Atoi(c.Query("page"))
	pageSize, _ := strconv.Atoi(c.Query("page_size"))

	q := dto.ListQuery{
		Page:     page,
		PageSize: pageSize,
		Status:   c.Query("status"),
		Channel:  c.Query("channel"),
	}

	result, err := h.service.List(c.Request.Context(), tenant.OrgID(c), q)
	if err != nil {
		response.InternalServerError(c, "Unable to fetch notifications")
		return
	}
	response.OK(c, "Notifications fetched successfully", result)
}

func (h *NotificationHandler) Get(c *gin.Context) {
	notification, err := h.service.Get(c.Request.Context(), tenant.OrgID(c), c.Param("id"))
	if err != nil {
		response.InternalServerError(c, "Unable to fetch notification")
		return
	}
	if notification == nil {
		response.NotFound(c, "Notification not found")
		return
	}
	response.OK(c, "Notification fetched successfully", notification)
}
