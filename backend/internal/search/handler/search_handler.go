package handler

import (
	"strings"

	"github.com/abhinavkumar03/crm-lite/backend/internal/search/service"
	"github.com/abhinavkumar03/crm-lite/backend/internal/shared/response"
	"github.com/abhinavkumar03/crm-lite/backend/internal/tenant"
	"github.com/gin-gonic/gin"
)

type SearchHandler struct {
	service *service.Service
}

func New(service *service.Service) *SearchHandler {
	return &SearchHandler{service: service}
}

func (h *SearchHandler) Search(c *gin.Context) {
	query := strings.TrimSpace(c.Query("q"))
	if query == "" {
		response.BadRequest(c, "Search query is required", nil)
		return
	}

	result, err := h.service.Search(
		c.Request.Context(),
		tenant.OrgID(c),
		query,
	)
	if err != nil {
		response.InternalServerError(c, "Unable to perform search")
		return
	}

	response.OK(c, "Search completed successfully", result)
}
