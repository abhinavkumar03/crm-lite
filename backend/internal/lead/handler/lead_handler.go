package handler

import (
	"log"
	"strconv"
	"strings"

	"github.com/abhinavkumar03/crm-lite/backend/internal/lead/dto"
	"github.com/abhinavkumar03/crm-lite/backend/internal/lead/service"
	"github.com/abhinavkumar03/crm-lite/backend/internal/shared/response"
	"github.com/abhinavkumar03/crm-lite/backend/internal/shared/validation"
	"github.com/gin-gonic/gin"
)

type LeadHandler struct {
	service *service.Service
}

func New(service *service.Service) *LeadHandler {
	return &LeadHandler{
		service: service,
	}
}

func (h *LeadHandler) Create(c *gin.Context) {

	userID := c.GetString("userID")

	var req dto.CreateLeadRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request body", nil)
		return
	}

	if err := validation.ValidateStruct(&req); err != nil {
		response.BadRequest(
			c,
			"Validation failed",
			validation.FormatErrors(err),
		)
		return
	}

	lead, err := h.service.Create(
		c.Request.Context(),
		userID,
		req,
	)

	if err != nil {
		log.Println(err)
		response.InternalServerError(
			c,
			"Unable to create lead",
		)
		return
	}

	response.Created(
		c,
		"Lead created successfully",
		lead,
	)
}

func (h *Handler) List(c *gin.Context) {

	userID := c.GetString("userID")

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))

	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))

	if page < 1 {
		page = 1
	}

	if limit <= 0 {
		limit = 20
	}

	if limit > 100 {
		limit = 100
	}

	req := dto.ListLeadRequest{
		Page:      page,
		Limit:     limit,
		Search:    strings.TrimSpace(c.Query("search")),
		Status:    strings.TrimSpace(c.Query("status")),
		SortBy:    c.DefaultQuery("sort_by", "created_at"),
		SortOrder: c.DefaultQuery("sort_order", "desc"),
	}

	leads, err := h.service.List(
		c.Request.Context(),
		userID,
		req,
	)

	if err != nil {
		response.InternalServerError(
			c,
			"Unable to fetch leads",
		)
		return
	}

	response.OK(
		c,
		"Leads fetched successfully",
		leads,
	)
}

func (h *LeadHandler) GetByID(c *gin.Context) {

	userID := c.GetString("userID")
	leadID := c.Param("id")

	lead, err := h.service.GetByID(
		c.Request.Context(),
		leadID,
		userID,
	)

	if err != nil {
		response.InternalServerError(
			c,
			"Unable to fetch lead",
		)
		return
	}

	if lead == nil {
		response.NotFound(
			c,
			"Lead not found",
		)
		return
	}

	response.OK(
		c,
		"Lead fetched successfully",
		lead,
	)
}

func (h *LeadHandler) Update(c *gin.Context) {

	userID := c.GetString("userID")
	leadID := c.Param("id")

	var req dto.UpdateLeadRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(
			c,
			"Invalid request body",
			nil,
		)
		return
	}

	lead, err := h.service.Update(
		c.Request.Context(),
		leadID,
		userID,
		req,
	)

	if err != nil {
		response.InternalServerError(
			c,
			"Unable to update lead",
		)
		return
	}

	if lead == nil {
		response.NotFound(
			c,
			"Lead not found",
		)
		return
	}

	response.OK(
		c,
		"Lead updated successfully",
		lead,
	)
}

func (h *LeadHandler) Delete(c *gin.Context) {

	userID := c.GetString("userID")
	leadID := c.Param("id")

	deleted, err := h.service.Delete(
		c.Request.Context(),
		leadID,
		userID,
	)

	if err != nil {
		response.InternalServerError(
			c,
			"Unable to delete lead",
		)
		return
	}

	if !deleted {
		response.NotFound(
			c,
			"Lead not found",
		)
		return
	}

	response.OK(
		c,
		"Lead deleted successfully",
		nil,
	)
}
