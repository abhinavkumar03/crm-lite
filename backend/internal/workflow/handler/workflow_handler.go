package handler

import (
	"errors"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/abhinavkumar03/crm-lite/backend/internal/shared/response"
	"github.com/abhinavkumar03/crm-lite/backend/internal/tenant"
	"github.com/abhinavkumar03/crm-lite/backend/internal/workflow/dto"
	"github.com/abhinavkumar03/crm-lite/backend/internal/workflow/service"
)

type Handler struct {
	svc *service.Service
}

func New(svc *service.Service) *Handler {
	return &Handler{svc: svc}
}

func userID(c *gin.Context) string { return c.GetString("userID") }

func (h *Handler) List(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	result, err := h.svc.List(c.Request.Context(), tenant.OrgID(c), c.Query("status"), c.Query("module_id"), page, pageSize)
	if err != nil {
		response.InternalServerError(c, err.Error())
		return
	}
	response.OK(c, "ok", result)
}

func (h *Handler) Create(c *gin.Context) {
	var req dto.CreateWorkflowRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error(), nil)
		return
	}
	result, err := h.svc.Create(c.Request.Context(), tenant.OrgID(c), userID(c), req)
	if err != nil {
		writeErr(c, err)
		return
	}
	response.Created(c, "created", result)
}

func (h *Handler) Get(c *gin.Context) {
	result, err := h.svc.Get(c.Request.Context(), tenant.OrgID(c), c.Param("id"))
	if err != nil {
		writeErr(c, err)
		return
	}
	response.OK(c, "ok", result)
}

func (h *Handler) Update(c *gin.Context) {
	var req dto.UpdateWorkflowRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error(), nil)
		return
	}
	result, err := h.svc.Update(c.Request.Context(), tenant.OrgID(c), c.Param("id"), userID(c), req)
	if err != nil {
		writeErr(c, err)
		return
	}
	response.OK(c, "updated", result)
}

func (h *Handler) Delete(c *gin.Context) {
	if err := h.svc.Archive(c.Request.Context(), tenant.OrgID(c), c.Param("id"), userID(c)); err != nil {
		writeErr(c, err)
		return
	}
	response.NoContent(c)
}

func (h *Handler) Publish(c *gin.Context) {
	var req dto.PublishRequest
	_ = c.ShouldBindJSON(&req)
	result, err := h.svc.Publish(c.Request.Context(), tenant.OrgID(c), c.Param("id"), userID(c), req)
	if err != nil {
		writeErr(c, err)
		return
	}
	response.OK(c, "published", result)
}

func (h *Handler) Disable(c *gin.Context) {
	result, err := h.svc.Disable(c.Request.Context(), tenant.OrgID(c), c.Param("id"), userID(c))
	if err != nil {
		writeErr(c, err)
		return
	}
	response.OK(c, "disabled", result)
}

func (h *Handler) Versions(c *gin.Context) {
	result, err := h.svc.ListVersions(c.Request.Context(), tenant.OrgID(c), c.Param("id"))
	if err != nil {
		writeErr(c, err)
		return
	}
	response.OK(c, "ok", result)
}

func (h *Handler) Rollback(c *gin.Context) {
	result, err := h.svc.Rollback(c.Request.Context(), tenant.OrgID(c), c.Param("id"), c.Param("versionId"), userID(c))
	if err != nil {
		writeErr(c, err)
		return
	}
	response.OK(c, "rolled_back", result)
}

func (h *Handler) ManualRun(c *gin.Context) {
	var req dto.ManualRunRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error(), nil)
		return
	}
	if err := h.svc.ManualRun(c.Request.Context(), tenant.OrgID(c), c.Param("id"), userID(c), req); err != nil {
		writeErr(c, err)
		return
	}
	response.OK(c, "queued", gin.H{"queued": true})
}

func (h *Handler) BuilderMetadata(c *gin.Context) {
	result, err := h.svc.BuilderMetadata(c.Request.Context(), tenant.OrgID(c))
	if err != nil {
		response.InternalServerError(c, err.Error())
		return
	}
	response.OK(c, "ok", result)
}

func (h *Handler) ListExecutions(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	result, err := h.svc.ListExecutions(c.Request.Context(), tenant.OrgID(c),
		c.Query("workflow_id"), c.Query("module_id"), c.Query("record_id"), c.Query("status"), page, pageSize)
	if err != nil {
		response.InternalServerError(c, err.Error())
		return
	}
	response.OK(c, "ok", result)
}

func (h *Handler) GetExecution(c *gin.Context) {
	result, err := h.svc.GetExecution(c.Request.Context(), tenant.OrgID(c), c.Param("id"))
	if err != nil {
		writeErr(c, err)
		return
	}
	response.OK(c, "ok", result)
}

func (h *Handler) RetryExecution(c *gin.Context) {
	result, err := h.svc.RetryExecution(c.Request.Context(), tenant.OrgID(c), c.Param("id"), userID(c))
	if err != nil {
		writeErr(c, err)
		return
	}
	response.OK(c, "queued", result)
}

func (h *Handler) ListTemplates(c *gin.Context) {
	result, err := h.svc.ListTemplates(c.Request.Context(), tenant.OrgID(c))
	if err != nil {
		response.InternalServerError(c, err.Error())
		return
	}
	response.OK(c, "ok", result)
}

func (h *Handler) CloneTemplate(c *gin.Context) {
	result, err := h.svc.CloneTemplate(c.Request.Context(), tenant.OrgID(c), c.Param("id"), userID(c))
	if err != nil {
		writeErr(c, err)
		return
	}
	response.Created(c, "cloned", result)
}

func (h *Handler) Metrics(c *gin.Context) {
	result, err := h.svc.Metrics(c.Request.Context(), tenant.OrgID(c))
	if err != nil {
		response.InternalServerError(c, err.Error())
		return
	}
	response.OK(c, "ok", result)
}

func writeErr(c *gin.Context, err error) {
	switch {
	case errors.Is(err, service.ErrNotFound):
		response.NotFound(c, err.Error())
	case errors.Is(err, service.ErrInvalidInput), errors.Is(err, service.ErrNotPublishable):
		response.BadRequest(c, err.Error(), nil)
	default:
		response.InternalServerError(c, err.Error())
	}
}
