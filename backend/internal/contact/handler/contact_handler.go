package handler

import (
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
