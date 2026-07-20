package handler

import (
	"errors"
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

func (h *NotificationHandler) Compose(c *gin.Context) {
	var req dto.ComposeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request body", nil)
		return
	}
	if err := validation.ValidateStruct(&req); err != nil {
		response.BadRequest(c, "Validation failed", validation.FormatErrors(err))
		return
	}

	notification, err := h.service.Compose(c.Request.Context(), tenant.OrgID(c), userID(c), req)
	if err != nil {
		writeServiceError(c, err, "Unable to queue notification")
		return
	}
	msg := "Notification queued successfully"
	switch req.Mode {
	case "draft":
		msg = "Draft saved successfully"
	case "schedule":
		msg = "Notification scheduled successfully"
	}
	response.Created(c, msg, notification)
}

// Send keeps the legacy handler name as an alias of Compose.
func (h *NotificationHandler) Send(c *gin.Context) { h.Compose(c) }

func (h *NotificationHandler) List(c *gin.Context) {
	page, _ := strconv.Atoi(c.Query("page"))
	pageSize, _ := strconv.Atoi(c.Query("page_size"))

	q := dto.ListQuery{
		Page:       page,
		PageSize:   pageSize,
		Status:     c.Query("status"),
		Channel:    c.Query("channel"),
		Q:          c.Query("q"),
		ModuleID:   c.Query("module_id"),
		EntityID:   c.Query("entity_id"),
		EntityType: c.Query("entity_type"),
		DateFrom:   c.Query("date_from"),
		DateTo:     c.Query("date_to"),
		TemplateID: c.Query("template_id"),
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

func (h *NotificationHandler) Retry(c *gin.Context) {
	notification, err := h.service.Retry(c.Request.Context(), tenant.OrgID(c), userID(c), c.Param("id"))
	if err != nil {
		writeServiceError(c, err, "Unable to retry notification")
		return
	}
	response.OK(c, "Notification retry queued", notification)
}

func (h *NotificationHandler) Cancel(c *gin.Context) {
	notification, err := h.service.Cancel(c.Request.Context(), tenant.OrgID(c), c.Param("id"))
	if err != nil {
		writeServiceError(c, err, "Unable to cancel notification")
		return
	}
	response.OK(c, "Notification cancelled", notification)
}

func (h *NotificationHandler) Metrics(c *gin.Context) {
	metrics, err := h.service.Metrics(c.Request.Context(), tenant.OrgID(c))
	if err != nil {
		response.InternalServerError(c, "Unable to fetch notification metrics")
		return
	}
	response.OK(c, "Notification metrics fetched successfully", metrics)
}

func (h *NotificationHandler) CreateTemplate(c *gin.Context) {
	var req dto.CreateTemplateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request body", nil)
		return
	}
	if err := validation.ValidateStruct(&req); err != nil {
		response.BadRequest(c, "Validation failed", validation.FormatErrors(err))
		return
	}
	tpl, err := h.service.CreateTemplate(c.Request.Context(), tenant.OrgID(c), userID(c), req)
	if err != nil {
		response.InternalServerError(c, "Unable to create template")
		return
	}
	response.Created(c, "Template created successfully", tpl)
}

func (h *NotificationHandler) ListTemplates(c *gin.Context) {
	page, _ := strconv.Atoi(c.Query("page"))
	pageSize, _ := strconv.Atoi(c.Query("page_size"))
	result, err := h.service.ListTemplates(
		c.Request.Context(), tenant.OrgID(c),
		c.Query("channel"), c.Query("category"), page, pageSize,
	)
	if err != nil {
		response.InternalServerError(c, "Unable to fetch templates")
		return
	}
	response.OK(c, "Templates fetched successfully", result)
}

func (h *NotificationHandler) GetTemplate(c *gin.Context) {
	tpl, err := h.service.GetTemplate(c.Request.Context(), tenant.OrgID(c), c.Param("id"))
	if err != nil {
		response.InternalServerError(c, "Unable to fetch template")
		return
	}
	if tpl == nil {
		response.NotFound(c, "Template not found")
		return
	}
	response.OK(c, "Template fetched successfully", tpl)
}

func (h *NotificationHandler) UpdateTemplate(c *gin.Context) {
	var req dto.UpdateTemplateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request body", nil)
		return
	}
	if err := validation.ValidateStruct(&req); err != nil {
		response.BadRequest(c, "Validation failed", validation.FormatErrors(err))
		return
	}
	tpl, err := h.service.UpdateTemplate(c.Request.Context(), tenant.OrgID(c), c.Param("id"), req)
	if err != nil {
		response.InternalServerError(c, "Unable to update template")
		return
	}
	if tpl == nil {
		response.NotFound(c, "Template not found")
		return
	}
	response.OK(c, "Template updated successfully", tpl)
}

func (h *NotificationHandler) DeleteTemplate(c *gin.Context) {
	if err := h.service.DeleteTemplate(c.Request.Context(), tenant.OrgID(c), c.Param("id")); err != nil {
		response.InternalServerError(c, "Unable to delete template")
		return
	}
	response.OK(c, "Template deleted successfully", nil)
}

func (h *NotificationHandler) PublishTemplate(c *gin.Context) {
	tpl, err := h.service.PublishTemplate(c.Request.Context(), tenant.OrgID(c), c.Param("id"))
	if err != nil {
		response.InternalServerError(c, "Unable to publish template")
		return
	}
	if tpl == nil {
		response.NotFound(c, "Template not found")
		return
	}
	response.OK(c, "Template published successfully", tpl)
}

func (h *NotificationHandler) UpdateDraft(c *gin.Context) {
	var req dto.ComposeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request body", nil)
		return
	}
	notification, err := h.service.UpdateDraft(c.Request.Context(), tenant.OrgID(c), userID(c), c.Param("id"), req)
	if err != nil {
		writeServiceError(c, err, "Unable to update draft")
		return
	}
	response.OK(c, "Draft updated successfully", notification)
}

func (h *NotificationHandler) SendDraft(c *gin.Context) {
	notification, err := h.service.SendDraft(c.Request.Context(), tenant.OrgID(c), userID(c), c.Param("id"))
	if err != nil {
		writeServiceError(c, err, "Unable to send draft")
		return
	}
	response.OK(c, "Draft queued for send", notification)
}

func (h *NotificationHandler) PreviewTemplate(c *gin.Context) {
	var req dto.PreviewTemplateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request body", nil)
		return
	}
	preview, err := h.service.PreviewTemplate(c.Request.Context(), tenant.OrgID(c), userID(c), c.Param("id"), req)
	if err != nil {
		response.InternalServerError(c, "Unable to preview template")
		return
	}
	if preview == nil {
		response.NotFound(c, "Template not found")
		return
	}
	response.OK(c, "Template preview rendered", preview)
}

func writeServiceError(c *gin.Context, err error, fallback string) {
	switch {
	case errors.Is(err, service.ErrNotFound):
		response.NotFound(c, "Notification not found")
	case errors.Is(err, service.ErrInvalidState),
		errors.Is(err, service.ErrScheduleRequired),
		errors.Is(err, service.ErrEmailSubject):
		response.BadRequest(c, err.Error(), nil)
	default:
		response.InternalServerError(c, fallback)
	}
}
