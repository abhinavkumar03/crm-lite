package handler

import (
	"net/http"

	"github.com/abhinavkumar03/crm-lite/backend/internal/activity/service"
	"github.com/gin-gonic/gin"
)

type ActivityHandler struct {
	service *service.Service
}

func New(
	service *service.Service,
) *ActivityHandler {

	return &ActivityHandler{
		service: service,
	}
}

func (h *ActivityHandler) listEntityActivities(
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

			"data": result,
		},
	)
}

func (h *ActivityHandler) ListLeadActivities(
	c *gin.Context,
) {

	h.listEntityActivities(

		c,

		"LEAD",

		c.Param("leadId"),
	)
}
func (h *ActivityHandler) ListContactActivities(
	c *gin.Context,
) {

	h.listEntityActivities(

		c,

		"CONTACT",

		c.Param("contactId"),
	)
}

func (h *ActivityHandler) ListTaskActivities(
	c *gin.Context,
) {

	h.listEntityActivities(

		c,

		"TASK",

		c.Param("taskId"),
	)
}
