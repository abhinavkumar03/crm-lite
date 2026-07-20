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

func (h *Handler) UpdateLayout(c *gin.Context) {
	var req dto.UpdateDetailLayoutRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid payload", nil)
		return
	}
	if err := h.validate.Struct(req); err != nil {
		response.BadRequest(c, "Validation failed", nil)
		return
	}
	layout, err := h.svc.UpdateDetailLayout(c.Request.Context(), tenant.OrgID(c), c.Param("id"), req)
	if err != nil {
		h.mapErr(c, err)
		return
	}
	response.OK(c, "Layout updated", layout)
}

func (h *Handler) GetFormLayout(c *gin.Context) {
	mode := c.DefaultQuery("mode", "create")
	layout, err := h.svc.GetFormLayout(c.Request.Context(), tenant.OrgID(c), c.Param("id"), mode)
	if err != nil {
		h.mapErr(c, err)
		return
	}
	response.OK(c, "Form layout fetched", layout)
}

func (h *Handler) UpdateFormLayout(c *gin.Context) {
	var req dto.UpdateFormLayoutRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid payload", nil)
		return
	}
	if err := h.validate.Struct(req); err != nil {
		response.BadRequest(c, "Validation failed", nil)
		return
	}
	layout, err := h.svc.UpdateFormLayout(c.Request.Context(), tenant.OrgID(c), c.Param("id"), req)
	if err != nil {
		h.mapErr(c, err)
		return
	}
	response.OK(c, "Form layout updated", layout)
}

func (h *Handler) ReorderFormFields(c *gin.Context) {
	var req dto.FormReorderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid payload", nil)
		return
	}
	if err := h.validate.Struct(req); err != nil {
		response.BadRequest(c, "Validation failed", nil)
		return
	}
	layout, err := h.svc.ReorderFormFields(c.Request.Context(), tenant.OrgID(c), c.Param("id"), req)
	if err != nil {
		h.mapErr(c, err)
		return
	}
	response.OK(c, "Form fields reordered", layout)
}

func (h *Handler) CreateFormSection(c *gin.Context) {
	var req dto.CreateSectionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid payload", nil)
		return
	}
	if err := h.validate.Struct(req); err != nil {
		response.BadRequest(c, "Validation failed", nil)
		return
	}
	layout, err := h.svc.CreateFormSection(c.Request.Context(), tenant.OrgID(c), c.Param("id"), req)
	if err != nil {
		h.mapErr(c, err)
		return
	}
	response.Created(c, "Section created", layout)
}

func (h *Handler) UpdateFormSection(c *gin.Context) {
	var req dto.UpdateSectionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid payload", nil)
		return
	}
	if err := h.validate.Struct(req); err != nil {
		response.BadRequest(c, "Validation failed", nil)
		return
	}
	layout, err := h.svc.UpdateFormSection(c.Request.Context(), tenant.OrgID(c), c.Param("id"), c.Param("sectionId"), req)
	if err != nil {
		h.mapErr(c, err)
		return
	}
	response.OK(c, "Section updated", layout)
}

func (h *Handler) DeleteFormSection(c *gin.Context) {
	layout, err := h.svc.DeleteFormSection(c.Request.Context(), tenant.OrgID(c), c.Param("id"), c.Param("sectionId"))
	if err != nil {
		h.mapErr(c, err)
		return
	}
	response.OK(c, "Section deleted", layout)
}

func (h *Handler) GetListLayout(c *gin.Context) {
	includeHidden := c.Query("include_hidden") == "true" || c.Query("admin") == "true"
	layout, err := h.svc.GetListLayout(c.Request.Context(), tenant.OrgID(c), c.Param("id"), includeHidden)
	if err != nil {
		h.mapErr(c, err)
		return
	}
	response.OK(c, "List layout fetched", layout)
}

func (h *Handler) UpdateListLayout(c *gin.Context) {
	var req dto.UpdateListLayoutRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid payload", nil)
		return
	}
	if err := h.validate.Struct(req); err != nil {
		response.BadRequest(c, "Validation failed", nil)
		return
	}
	layout, err := h.svc.UpdateListLayout(c.Request.Context(), tenant.OrgID(c), c.Param("id"), req)
	if err != nil {
		h.mapErr(c, err)
		return
	}
	response.OK(c, "List layout updated", layout)
}

func (h *Handler) ReorderListColumns(c *gin.Context) {
	var req dto.ListReorderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid payload", nil)
		return
	}
	if err := h.validate.Struct(req); err != nil {
		response.BadRequest(c, "Validation failed", nil)
		return
	}
	layout, err := h.svc.ReorderListColumns(c.Request.Context(), tenant.OrgID(c), c.Param("id"), req)
	if err != nil {
		h.mapErr(c, err)
		return
	}
	response.OK(c, "List columns reordered", layout)
}

func (h *Handler) ToggleListColumn(c *gin.Context) {
	var req dto.ListToggleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid payload", nil)
		return
	}
	if err := h.validate.Struct(req); err != nil {
		response.BadRequest(c, "Validation failed", nil)
		return
	}
	layout, err := h.svc.ToggleListColumn(c.Request.Context(), tenant.OrgID(c), c.Param("id"), req)
	if err != nil {
		h.mapErr(c, err)
		return
	}
	response.OK(c, "List column toggled", layout)
}

func (h *Handler) ResetListLayout(c *gin.Context) {
	layout, err := h.svc.ResetListLayout(c.Request.Context(), tenant.OrgID(c), c.Param("id"))
	if err != nil {
		h.mapErr(c, err)
		return
	}
	response.OK(c, "List layout reset", layout)
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
	switch {
	case errors.Is(err, service.ErrNotFound), errors.Is(err, service.ErrSectionNotFound):
		response.NotFound(c, "Not found")
	case errors.Is(err, service.ErrSectionNotEmpty),
		errors.Is(err, service.ErrSectionExists),
		errors.Is(err, service.ErrSystemColumn),
		errors.Is(err, service.ErrLockedColumn),
		errors.Is(err, service.ErrInvalidListCol),
		errors.Is(err, service.ErrInvalidMode):
		response.BadRequest(c, err.Error(), nil)
	default:
		response.BadRequest(c, err.Error(), nil)
	}
}
