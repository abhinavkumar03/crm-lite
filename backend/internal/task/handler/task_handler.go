package handler

import (
	"strconv"
	"strings"

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

	req := dto.ListTasksRequest{
		Page:      page,
		Limit:     limit,
		Search:    strings.TrimSpace(c.Query("search")),
		Status:    strings.TrimSpace(c.Query("status")),
		SortBy:    c.DefaultQuery("sort_by", "created_at"),
		SortOrder: c.DefaultQuery("sort_order", "desc"),
	}

	result, err := h.service.List(
		c.Request.Context(),
		userID,
		req,
	)

	if err != nil {
		println(err.Error())

		response.InternalServerError(
			c,
			err.Error(),
		)

		return
	}

	response.OK(
		c,
		"Tasks fetched successfully",
		result,
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
