package handler

import (
	"strconv"

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

func (h *TaskHandler) List(c *gin.Context) {

	userID := c.GetString("userID")

	req := dto.ListTasksRequest{
		Page:   1,
		Limit:  10,
		Search: c.DefaultQuery("search", ""),
		Status: c.DefaultQuery("status", ""),
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

	tasks, err := h.service.List(
		c.Request.Context(),
		userID,
		req,
	)

	if err != nil {

		response.InternalServerError(
			c,
			"Unable to fetch tasks",
		)

		return
	}

	response.OK(
		c,
		"Tasks fetched successfully",
		tasks,
	)
}
