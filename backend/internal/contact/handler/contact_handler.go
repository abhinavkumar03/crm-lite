package handler

import (
	"strconv"

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

	req := dto.ListContactsRequest{
		Page:   1,
		Limit:  10,
		Search: c.DefaultQuery("search", ""),
	}

	if page := c.Query("page"); page != "" {
		if p, err := strconv.Atoi(page); err == nil && p > 0 {
			req.Page = p
		}
	}

	if limit := c.Query("limit"); limit != "" {
		if l, err := strconv.Atoi(limit); err == nil && l > 0 && l <= 100 {
			req.Limit = l
		}
	}

	contacts, err := h.service.List(
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
		contacts,
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
