package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Success sends a generic successful response.
func Success(
	c *gin.Context,
	statusCode int,
	message string,
	data interface{},
) {
	c.JSON(statusCode, APIResponse{
		Success: true,
		Message: message,
		Data:    data,
	})
}

// OK sends HTTP 200.
func OK(
	c *gin.Context,
	message string,
	data interface{},
) {
	Success(
		c,
		http.StatusOK,
		message,
		data,
	)
}

// Created sends HTTP 201.
func Created(
	c *gin.Context,
	message string,
	data interface{},
) {
	Success(
		c,
		http.StatusCreated,
		message,
		data,
	)
}

// NoContent sends HTTP 204.
func NoContent(
	c *gin.Context,
) {
	c.Status(http.StatusNoContent)
}
