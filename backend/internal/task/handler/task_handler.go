package handler

import (
	"github.com/abhinavkumar03/crm-lite/backend/internal/shared/response"
	"github.com/abhinavkumar03/crm-lite/backend/internal/shared/validation"
	"github.com/abhinavkumar03/crm-lite/backend/internal/task/dto"
	"github.com/abhinavkumar03/crm-lite/backend/internal/task/service"
	"github.com/gin-gonic/gin"
)

type TaskHandler struct {
	service *service.Service
}

func New(service *service.Service) *TaskHandler {
	return &TaskHandler{
		service: service,
	}
}

func (h *TaskHandler) Create(c *gin.Context) {

	userID := c.GetString("userID")

	var req dto.CreateTaskRequest

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

	task, err := h.service.Create(
		c.Request.Context(),
		userID,
		req,
	)

	if err != nil {
		response.BadRequest(c, err.Error(), nil)
		return
	}

	response.Created(
		c,
		"Task created successfully",
		task,
	)
}
