package handler

import (
	"net/http"

	"github.com/abhinavkumar03/crm-lite/backend/internal/attachment/dto"
	"github.com/abhinavkumar03/crm-lite/backend/internal/attachment/service"
	"github.com/gin-gonic/gin"
)

type AttachmentHandler struct {
	service *service.Service
}

func New(
	service *service.Service,
) *AttachmentHandler {
	return &AttachmentHandler{
		service: service,
	}
}

func (h *AttachmentHandler) createEntityAttachment(
	c *gin.Context,
	entityType string,
	entityID string,
) {

	var req dto.CreateAttachmentRequest

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
			"message": "Attachment added successfully",
		},
	)
}

func (h *AttachmentHandler) listEntityAttachments(
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

func (h *AttachmentHandler) CreateLeadAttachment(
	c *gin.Context,
) {

	h.createEntityAttachment(
		c,
		"LEAD",
		c.Param("leadId"),
	)
}

func (h *AttachmentHandler) ListLeadAttachments(
	c *gin.Context,
) {

	h.listEntityAttachments(
		c,
		"LEAD",
		c.Param("leadId"),
	)
}

func (h *AttachmentHandler) CreateContactAttachment(
	c *gin.Context,
) {

	h.createEntityAttachment(
		c,
		"CONTACT",
		c.Param("contactId"),
	)
}

func (h *AttachmentHandler) ListContactAttachments(
	c *gin.Context,
) {

	h.listEntityAttachments(
		c,
		"CONTACT",
		c.Param("contactId"),
	)
}

func (h *AttachmentHandler) CreateTaskAttachment(
	c *gin.Context,
) {

	h.createEntityAttachment(
		c,
		"TASK",
		c.Param("taskId"),
	)
}

func (h *AttachmentHandler) ListTaskAttachments(
	c *gin.Context,
) {

	h.listEntityAttachments(
		c,
		"TASK",
		c.Param("taskId"),
	)
}

func (h *AttachmentHandler) DeleteAttachment(
	c *gin.Context,
) {

	userID := c.GetString("userID")

	err := h.service.Delete(
		c.Request.Context(),
		userID,
		c.Param("attachmentId"),
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
			"message": "Attachment deleted successfully",
		},
	)
}
