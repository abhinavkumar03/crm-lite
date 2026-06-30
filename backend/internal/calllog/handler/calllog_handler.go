package handler

import (
	"net/http"

	"github.com/abhinavkumar03/crm-lite/backend/internal/calllog/dto"
	"github.com/abhinavkumar03/crm-lite/backend/internal/calllog/service"
	"github.com/gin-gonic/gin"
)

type CallLogHandler struct {
	service *service.Service
}

func New(service *service.Service) *CallLogHandler {
	return &CallLogHandler{
		service: service,
	}
}

func (h *CallLogHandler) createEntityCallLog(
	c *gin.Context,
	entityType string,
	entityID string,
) {

	var req dto.CreateCallLogRequest

	if err := c.ShouldBindJSON(&req); err != nil {

		c.JSON(
			http.StatusBadRequest,
			gin.H{
				"success": false,
				"message": err.Error(),
			},
		)

		return
	}

	userID := c.GetString("userID")

	err := h.service.Create(
		c.Request.Context(),
		userID,
		entityType,
		entityID,
		req,
	)

	if err != nil {

		c.JSON(
			http.StatusBadRequest,
			gin.H{
				"success": false,
				"message": err.Error(),
			},
		)

		return
	}

	c.JSON(
		http.StatusCreated,
		gin.H{
			"success": true,
			"message": "Call log created successfully",
		},
	)
}

func (h *CallLogHandler) listEntityCallLogs(
	c *gin.Context,
	entityType string,
	entityID string,
) {

	userID := c.GetString("userID")

	result, err := h.service.List(
		c.Request.Context(),
		userID,
		entityType,
		entityID,
	)

	if err != nil {

		c.JSON(
			http.StatusBadRequest,
			gin.H{
				"success": false,
				"message": err.Error(),
			},
		)

		return
	}

	c.JSON(
		http.StatusOK,
		gin.H{
			"success": true,
			"data":    result,
		},
	)
}

func (h *CallLogHandler) CreateLeadCallLog(
	c *gin.Context,
) {
	h.createEntityCallLog(
		c,
		"LEAD",
		c.Param("leadId"),
	)
}

func (h *CallLogHandler) ListLeadCallLogs(
	c *gin.Context,
) {
	h.listEntityCallLogs(
		c,
		"LEAD",
		c.Param("leadId"),
	)
}

func (h *CallLogHandler) CreateContactCallLog(
	c *gin.Context,
) {
	h.createEntityCallLog(
		c,
		"CONTACT",
		c.Param("contactId"),
	)
}

func (h *CallLogHandler) ListContactCallLogs(
	c *gin.Context,
) {
	h.listEntityCallLogs(
		c,
		"CONTACT",
		c.Param("contactId"),
	)
}

func (h *CallLogHandler) CreateTaskCallLog(
	c *gin.Context,
) {
	h.createEntityCallLog(
		c,
		"TASK",
		c.Param("taskId"),
	)
}

func (h *CallLogHandler) ListTaskCallLogs(
	c *gin.Context,
) {
	h.listEntityCallLogs(
		c,
		"TASK",
		c.Param("taskId"),
	)
}

func (h *CallLogHandler) UpdateCallLog(
	c *gin.Context,
) {

	var req dto.UpdateCallLogRequest

	if err := c.ShouldBindJSON(&req); err != nil {

		c.JSON(
			http.StatusBadRequest,
			gin.H{
				"success": false,
				"message": err.Error(),
			},
		)

		return
	}

	userID := c.GetString("userID")

	err := h.service.Update(
		c.Request.Context(),
		userID,
		c.Param("callLogId"),
		req,
	)

	if err != nil {

		c.JSON(
			http.StatusBadRequest,
			gin.H{
				"success": false,
				"message": err.Error(),
			},
		)

		return
	}

	c.JSON(
		http.StatusOK,
		gin.H{
			"success": true,
			"message": "Call log updated successfully",
		},
	)
}

func (h *CallLogHandler) DeleteCallLog(
	c *gin.Context,
) {

	userID := c.GetString("userID")

	err := h.service.Delete(
		c.Request.Context(),
		userID,
		c.Param("callLogId"),
	)

	if err != nil {

		c.JSON(
			http.StatusBadRequest,
			gin.H{
				"success": false,
				"message": err.Error(),
			},
		)

		return
	}

	c.JSON(
		http.StatusOK,
		gin.H{
			"success": true,
			"message": "Call log deleted successfully",
		},
	)
}
