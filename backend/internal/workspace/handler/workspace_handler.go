package handler

import (
	"errors"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"

	"github.com/abhinavkumar03/crm-lite/backend/internal/shared/response"
	"github.com/abhinavkumar03/crm-lite/backend/internal/tenant"
	"github.com/abhinavkumar03/crm-lite/backend/internal/workspace/dto"
	"github.com/abhinavkumar03/crm-lite/backend/internal/workspace/service"
)

type Handler struct {
	svc      *service.Service
	validate *validator.Validate
}

func New(svc *service.Service) *Handler {
	return &Handler{svc: svc, validate: validator.New()}
}

func (h *Handler) GetLayout(c *gin.Context) {
	layout, err := h.svc.GetDetailLayout(c.Request.Context(), tenant.OrgID(c), c.Param("id"))
	if err != nil {
		response.InternalServerError(c, "Unable to load layout")
		return
	}
	response.OK(c, "Layout fetched", layout)
}

func (h *Handler) ListNotes(c *gin.Context) {
	items, err := h.svc.ListNotes(c.Request.Context(), tenant.OrgID(c), c.Param("id"), c.Param("recordId"))
	if err != nil {
		h.mapErr(c, err)
		return
	}
	response.OK(c, "Notes fetched", items)
}

func (h *Handler) CreateNote(c *gin.Context) {
	var req dto.CreateNoteRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid payload", nil)
		return
	}
	if err := h.validate.Struct(req); err != nil {
		response.BadRequest(c, "Validation failed", nil)
		return
	}
	item, err := h.svc.CreateNote(
		c.Request.Context(), tenant.OrgID(c), c.Param("id"), c.Param("recordId"),
		c.GetString("userID"), req,
	)
	if err != nil {
		h.mapErr(c, err)
		return
	}
	response.Created(c, "Note created", item)
}

func (h *Handler) DeleteNote(c *gin.Context) {
	err := h.svc.DeleteNote(
		c.Request.Context(), tenant.OrgID(c), c.Param("id"), c.Param("recordId"),
		c.Param("noteId"), c.GetString("userID"),
	)
	if err != nil {
		h.mapErr(c, err)
		return
	}
	response.OK(c, "Note deleted", nil)
}

func (h *Handler) ListAttachments(c *gin.Context) {
	items, err := h.svc.ListAttachments(c.Request.Context(), tenant.OrgID(c), c.Param("id"), c.Param("recordId"))
	if err != nil {
		h.mapErr(c, err)
		return
	}
	response.OK(c, "Attachments fetched", items)
}

func (h *Handler) CreateAttachment(c *gin.Context) {
	var req dto.CreateAttachmentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid payload", nil)
		return
	}
	if err := h.validate.Struct(req); err != nil {
		response.BadRequest(c, "Validation failed", nil)
		return
	}
	item, err := h.svc.CreateAttachment(
		c.Request.Context(), tenant.OrgID(c), c.Param("id"), c.Param("recordId"),
		c.GetString("userID"), req,
	)
	if err != nil {
		h.mapErr(c, err)
		return
	}
	response.Created(c, "Attachment created", item)
}

func (h *Handler) DeleteAttachment(c *gin.Context) {
	err := h.svc.DeleteAttachment(
		c.Request.Context(), tenant.OrgID(c), c.Param("id"), c.Param("recordId"),
		c.Param("attachmentId"), c.GetString("userID"),
	)
	if err != nil {
		h.mapErr(c, err)
		return
	}
	response.OK(c, "Attachment deleted", nil)
}

func (h *Handler) ListActivities(c *gin.Context) {
	items, err := h.svc.ListActivities(c.Request.Context(), tenant.OrgID(c), c.Param("id"), c.Param("recordId"))
	if err != nil {
		h.mapErr(c, err)
		return
	}
	response.OK(c, "Activities fetched", items)
}

func (h *Handler) ListRelated(c *gin.Context) {
	items, err := h.svc.ListRelated(c.Request.Context(), tenant.OrgID(c), c.Param("id"))
	if err != nil {
		response.InternalServerError(c, "Unable to list related modules")
		return
	}
	response.OK(c, "Related modules fetched", items)
}

func (h *Handler) mapErr(c *gin.Context, err error) {
	if errors.Is(err, service.ErrNotFound) {
		response.NotFound(c, "Not found")
		return
	}
	response.BadRequest(c, err.Error(), nil)
}
