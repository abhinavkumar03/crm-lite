package handler

import (
	"github.com/abhinavkumar03/crm-lite/backend/internal/dashboard/service"
	"github.com/abhinavkumar03/crm-lite/backend/internal/shared/response"
	"github.com/gin-gonic/gin"
)

type DashboardHandler struct {
	service *service.Service
}

func New(service *service.Service) *DashboardHandler {
	return &DashboardHandler{
		service: service,
	}
}

func (h *DashboardHandler) Dashboard(c *gin.Context) {

	userID := c.GetString("userID")
	refresh := c.Query("refresh") == "true"

	data, err := h.service.GetDashboard(
		c.Request.Context(),
		userID,
		refresh,
	)

	if err != nil {

		response.InternalServerError(
			c,
			"Unable to fetch dashboard",
		)

		return
	}

	response.OK(
		c,
		"Dashboard fetched successfully",
		data,
	)
}
