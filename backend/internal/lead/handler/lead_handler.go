package handler

import (
	"log"

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
