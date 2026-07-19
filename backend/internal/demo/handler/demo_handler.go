package handler

import (
	"errors"

	"github.com/gin-gonic/gin"

	"github.com/abhinavkumar03/crm-lite/backend/internal/demo/dto"
	"github.com/abhinavkumar03/crm-lite/backend/internal/demo/service"
	"github.com/abhinavkumar03/crm-lite/backend/internal/shared/response"
)

type Handler struct {
	svc *service.Service
}

func New(svc *service.Service) *Handler {
	return &Handler{svc: svc}
}

func (h *Handler) Catalog(c *gin.Context) {
	info, err := h.svc.Catalog(c.Request.Context())
	if err != nil {
		response.NotFound(c, "Demo workflow not found — run migrations")
		return
	}
	response.OK(c, "Demo catalog", info)
}

func (h *Handler) Workflow(c *gin.Context) {
	def, err := h.svc.GetWorkflowDefinition(c.Request.Context(), c.Param("key"))
	if err != nil {
		response.NotFound(c, "Workflow not found")
		return
	}
	response.OK(c, "Workflow definition", def)
}

func (h *Handler) Active(c *gin.Context) {
	sess, err := h.svc.GetActive(c.Request.Context(), c.GetString("userID"))
	if err != nil {
		response.InternalServerError(c, "Unable to load demo session")
		return
	}
	response.OK(c, "Active demo session", sess)
}

func (h *Handler) Start(c *gin.Context) {
	sess, err := h.svc.Start(c.Request.Context(), c.GetString("userID"))
	if err != nil {
		response.BadRequest(c, err.Error(), nil)
		return
	}
	response.Created(c, "Demo started — switched to sandbox organization", sess)
}

func (h *Handler) Restart(c *gin.Context) {
	sess, err := h.svc.Restart(c.Request.Context(), c.GetString("userID"))
	if err != nil {
		response.BadRequest(c, err.Error(), nil)
		return
	}
	response.OK(c, "Demo restarted", sess)
}

func (h *Handler) Validate(c *gin.Context) {
	var body dto.ValidateRequest
	_ = c.ShouldBindJSON(&body)
	if body.StepKey == "" {
		body.StepKey = c.Param("stepKey")
	}
	res, err := h.svc.Validate(
		c.Request.Context(), c.GetString("userID"), c.Param("sessionId"), body.StepKey, body.Route, body.ClientEvent,
	)
	if errors.Is(err, service.ErrNotFound) {
		response.NotFound(c, "Session or step not found")
		return
	}
	if err != nil {
		response.BadRequest(c, err.Error(), nil)
		return
	}
	if !res.OK {
		response.BadRequest(c, res.Message, res)
		return
	}
	response.OK(c, res.Message, res)
}

func (h *Handler) Skip(c *gin.Context) {
	var body struct {
		StepKey string `json:"step_key"`
	}
	_ = c.ShouldBindJSON(&body)
	if body.StepKey == "" {
		body.StepKey = c.Param("stepKey")
	}
	sess, err := h.svc.Skip(c.Request.Context(), c.GetString("userID"), c.Param("sessionId"), body.StepKey)
	if err != nil {
		response.BadRequest(c, err.Error(), nil)
		return
	}
	response.OK(c, "Step skipped", sess)
}

func (h *Handler) Complete(c *gin.Context) {
	sess, err := h.svc.Complete(c.Request.Context(), c.GetString("userID"), c.Param("sessionId"))
	if err != nil {
		response.BadRequest(c, err.Error(), nil)
		return
	}
	response.OK(c, "Demo completed", sess)
}

func (h *Handler) Cleanup(c *gin.Context) {
	var req dto.CleanupRequest
	_ = c.ShouldBindJSON(&req)
	sess, err := h.svc.Cleanup(c.Request.Context(), c.GetString("userID"), c.Param("sessionId"), req.KeepData)
	if err != nil {
		response.BadRequest(c, err.Error(), nil)
		return
	}
	response.OK(c, "Demo cleanup finished", sess)
}

func (h *Handler) LogEvent(c *gin.Context) {
	var req dto.EventRequest
	if err := c.ShouldBindJSON(&req); err != nil || req.EventType == "" {
		response.BadRequest(c, "event_type is required", nil)
		return
	}
	if err := h.svc.LogClientEvent(
		c.Request.Context(), c.GetString("userID"), c.Param("sessionId"), req.EventType, req.Payload,
	); err != nil {
		response.BadRequest(c, err.Error(), nil)
		return
	}
	response.OK(c, "Event recorded", nil)
}
