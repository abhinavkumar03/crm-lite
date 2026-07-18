package handler

import (
	"errors"

	"github.com/gin-gonic/gin"

	"github.com/abhinavkumar03/crm-lite/backend/internal/shared/response"
	"github.com/abhinavkumar03/crm-lite/backend/internal/shared/validation"
	"github.com/abhinavkumar03/crm-lite/backend/internal/tenant"
	"github.com/abhinavkumar03/crm-lite/backend/internal/validationengine/dto"
	"github.com/abhinavkumar03/crm-lite/backend/internal/validationengine/service"
)

const (
	paramModuleID = "id"
	paramRuleID   = "ruleId"
)

type Handler struct {
	service *service.Service
}

func New(service *service.Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) ListRules(c *gin.Context) {
	rules, err := h.service.List(c.Request.Context(), tenant.OrgID(c), c.Param(paramModuleID))
	if err != nil {
		h.writeServiceError(c, err, "Unable to fetch validation rules")
		return
	}
	response.OK(c, "Validation rules fetched successfully", rules)
}

func (h *Handler) GetRule(c *gin.Context) {
	rule, err := h.service.GetByID(c.Request.Context(), tenant.OrgID(c), c.Param(paramModuleID), c.Param(paramRuleID))
	if err != nil {
		h.writeServiceError(c, err, "Unable to fetch validation rule")
		return
	}
	if rule == nil {
		response.NotFound(c, "Validation rule not found")
		return
	}
	response.OK(c, "Validation rule fetched successfully", rule)
}

func (h *Handler) CreateRule(c *gin.Context) {
	var req dto.CreateRuleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request body", nil)
		return
	}
	if err := validation.ValidateStruct(&req); err != nil {
		response.BadRequest(c, "Validation failed", validation.FormatErrors(err))
		return
	}

	rule, err := h.service.Create(c.Request.Context(), tenant.OrgID(c), c.Param(paramModuleID), req)
	if err != nil {
		h.writeServiceError(c, err, "Unable to create validation rule")
		return
	}
	response.Created(c, "Validation rule created successfully", rule)
}

func (h *Handler) UpdateRule(c *gin.Context) {
	var req dto.UpdateRuleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request body", nil)
		return
	}
	if err := validation.ValidateStruct(&req); err != nil {
		response.BadRequest(c, "Validation failed", validation.FormatErrors(err))
		return
	}

	rule, err := h.service.Update(c.Request.Context(), tenant.OrgID(c), c.Param(paramModuleID), c.Param(paramRuleID), req)
	if err != nil {
		h.writeServiceError(c, err, "Unable to update validation rule")
		return
	}
	response.OK(c, "Validation rule updated successfully", rule)
}

func (h *Handler) DeleteRule(c *gin.Context) {
	if err := h.service.Delete(c.Request.Context(), tenant.OrgID(c), c.Param(paramModuleID), c.Param(paramRuleID)); err != nil {
		h.writeServiceError(c, err, "Unable to delete validation rule")
		return
	}
	response.OK(c, "Validation rule deleted successfully", nil)
}

func (h *Handler) Schema(c *gin.Context) {
	schema, err := h.service.Schema(c.Request.Context(), tenant.OrgID(c), c.Param(paramModuleID))
	if err != nil {
		h.writeServiceError(c, err, "Unable to build validation schema")
		return
	}
	response.OK(c, "Validation schema fetched successfully", schema)
}

func (h *Handler) Validate(c *gin.Context) {
	var req dto.ValidateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request body", nil)
		return
	}

	result, err := h.service.Validate(c.Request.Context(), tenant.OrgID(c), c.Param(paramModuleID), req.Data)
	if err != nil {
		h.writeServiceError(c, err, "Unable to validate payload")
		return
	}
	response.OK(c, "Validation completed", result)
}

func (h *Handler) writeServiceError(c *gin.Context, err error, fallback string) {
	switch {
	case errors.Is(err, service.ErrModuleNotFound):
		response.NotFound(c, "Module not found")
	case errors.Is(err, service.ErrNotFound):
		response.NotFound(c, "Validation rule not found")
	case errors.Is(err, service.ErrInvalidRuleType),
		errors.Is(err, service.ErrFieldRequired),
		errors.Is(err, service.ErrModuleRule),
		errors.Is(err, service.ErrInvalidParams):
		response.BadRequest(c, err.Error(), nil)
	default:
		response.InternalServerError(c, fallback)
	}
}
