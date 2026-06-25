package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Error sends a standardized error response.
func Error(
	c *gin.Context,
	statusCode int,
	message string,
	errors interface{},
) {
	c.JSON(statusCode, APIResponse{
		Success: false,
		Message: message,
		Errors:  errors,
	})
}

// BadRequest returns HTTP 400.
func BadRequest(
	c *gin.Context,
	message string,
	errors interface{},
) {
	Error(
		c,
		http.StatusBadRequest,
		message,
		errors,
	)
}

// Unauthorized returns HTTP 401.
func Unauthorized(
	c *gin.Context,
	message string,
) {
	Error(
		c,
		http.StatusUnauthorized,
		message,
		nil,
	)
}

// Forbidden returns HTTP 403.
func Forbidden(
	c *gin.Context,
	message string,
) {
	Error(
		c,
		http.StatusForbidden,
		message,
		nil,
	)
}

// NotFound returns HTTP 404.
func NotFound(
	c *gin.Context,
	message string,
) {
	Error(
		c,
		http.StatusNotFound,
		message,
		nil,
	)
}

// Conflict returns HTTP 409.
func Conflict(
	c *gin.Context,
	message string,
	errors interface{},
) {
	Error(
		c,
		http.StatusConflict,
		message,
		errors,
	)
}

// InternalServerError returns HTTP 500.
func InternalServerError(
	c *gin.Context,
	message string,
) {
	Error(
		c,
		http.StatusInternalServerError,
		message,
		nil,
	)
}
