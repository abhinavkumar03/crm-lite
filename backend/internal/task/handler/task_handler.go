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

func (h *TaskHandler) GetByID(c *gin.Context) {

	userID := c.GetString("userID")
	taskID := c.Param("id")

	task, err := h.service.GetByID(
		c.Request.Context(),
		taskID,
		userID,
	)

	if err != nil {
		response.InternalServerError(
			c,
			"Unable to fetch task",
		)
		return
	}

	if task == nil {
		response.NotFound(
			c,
			"Task not found",
		)
		return
	}

	response.OK(
		c,
		"Task fetched successfully",
		task,
	)
}

func (h *TaskHandler) Update(c *gin.Context) {

	userID := c.GetString("userID")
	taskID := c.Param("id")

	var req dto.UpdateTaskRequest

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

	task, err := h.service.Update(
		c.Request.Context(),
		taskID,
		userID,
		req,
	)

	if err != nil {

		response.BadRequest(
			c,
			err.Error(),
			nil,
		)

		return
	}

	if task == nil {

		response.NotFound(
			c,
			"Task not found",
		)

		return
	}

	response.OK(
		c,
		"Task updated successfully",
		task,
	)
}

func (h *TaskHandler) Delete(c *gin.Context) {

	userID := c.GetString("userID")
	taskID := c.Param("id")

	deleted, err := h.service.Delete(
		c.Request.Context(),
		taskID,
		userID,
	)

	if err != nil {

		response.InternalServerError(
			c,
			"Unable to delete task",
		)

		return
	}

	if !deleted {

		response.NotFound(
			c,
			"Task not found",
		)

		return
	}

	response.OK(
		c,
		"Task deleted successfully",
		nil,
	)
}
