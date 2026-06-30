package handler

import (
	"strconv"
	"strings"

	"github.com/abhinavkumar03/crm-lite/backend/internal/contact/dto"
	"github.com/abhinavkumar03/crm-lite/backend/internal/contact/service"
	"github.com/abhinavkumar03/crm-lite/backend/internal/shared/response"
	"github.com/abhinavkumar03/crm-lite/backend/internal/shared/validation"
	"github.com/gin-gonic/gin"
)

type ContactHandler struct {
	service *service.Service
}

func New(service *service.Service) *ContactHandler {
	return &ContactHandler{
		service: service,
	}
}

func (h *ContactHandler) Create(c *gin.Context) {

	userID := c.GetString("userID")

	var req dto.CreateContactRequest

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

	contact, err := h.service.Create(
		c.Request.Context(),
		userID,
		req,
	)

	if err != nil {
		response.InternalServerError(
			c,
			"Unable to create contact",
		)
		return
	}

	response.Created(
		c,
		"Contact created successfully",
		contact,
	)
}

func (h *ContactHandler) List(c *gin.Context) {

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

	req := dto.ListContactsRequest{
		Page:      page,
		Limit:     limit,
		Search:    strings.TrimSpace(c.Query("search")),
		SortBy:    c.DefaultQuery("sort_by", "created_at"),
		SortOrder: c.DefaultQuery("sort_order", "desc"),
	}

	result, err := h.service.List(
		c.Request.Context(),
		userID,
		req,
	)

	if err != nil {

		response.InternalServerError(
			c,
			"Unable to fetch contacts",
		)

		return
	}

	response.OK(
		c,
		"Contacts fetched successfully",
		result,
	)
}

func (h *ContactHandler) GetByID(c *gin.Context) {

	userID := c.GetString("userID")
	contactID := c.Param("id")

	contact, err := h.service.GetByID(
		c.Request.Context(),
		contactID,
		userID,
	)

	if err != nil {
		response.InternalServerError(
			c,
			"Unable to fetch contact",
		)
		return
	}

	if contact == nil {
		response.NotFound(
			c,
			"Contact not found",
		)
		return
	}

	response.OK(
		c,
		"Contact fetched successfully",
		contact,
	)
}

func (h *ContactHandler) Update(c *gin.Context) {

	userID := c.GetString("userID")
	contactID := c.Param("id")

	var req dto.UpdateContactRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(
			c,
			"Invalid request body",
			nil,
		)
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

	contact, err := h.service.Update(
		c.Request.Context(),
		contactID,
		userID,
		req,
	)

	if err != nil {
		response.InternalServerError(
			c,
			"Unable to update contact",
		)
		return
	}

	if contact == nil {
		response.NotFound(
			c,
			"Contact not found",
		)
		return
	}

	response.OK(
		c,
		"Contact updated successfully",
		contact,
	)
}

func (h *ContactHandler) Delete(c *gin.Context) {

	userID := c.GetString("userID")
	contactID := c.Param("id")

	deleted, err := h.service.Delete(
		c.Request.Context(),
		contactID,
		userID,
	)

	if err != nil {
		response.InternalServerError(
			c,
			"Unable to delete contact",
		)
		return
	}

	if !deleted {
		response.NotFound(
			c,
			"Contact not found",
		)
		return
	}

	response.OK(
		c,
		"Contact deleted successfully",
		nil,
	)
}
