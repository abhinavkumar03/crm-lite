package dto

import (
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/abhinavkumar03/crm-lite/backend/internal/shared/constants"
)

func PaginationFromRequest(
	c *gin.Context,
) PaginationRequest {

	page := constants.DefaultPage
	pageSize := constants.DefaultPageSize

	if value := c.Query("page"); value != "" {

		if parsed, err := strconv.Atoi(value); err == nil {
			page = parsed
		}
	}

	if value := c.Query("page_size"); value != "" {

		if parsed, err := strconv.Atoi(value); err == nil {
			pageSize = parsed
		}
	}

	return PaginationRequest{
		Page:     page,
		PageSize: pageSize,
	}
}
