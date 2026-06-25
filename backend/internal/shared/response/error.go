package response

import (
	"net/http"

	"github.com/gin-gonic/gin"

	apperrors "github.com/abhinavkumar03/crm-lite/backend/internal/shared/errors"
)

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

func FromAppError(
	c *gin.Context,
	err *apperrors.AppError,
) {
	Error(
		c,
		err.StatusCode,
		err.Message,
		[]gin.H{
			{
				"code": err.Code,
			},
		},
	)
}

func BadRequest(
	c *gin.Context,
	message string,
	errors interface{},
) {
	Error(c, http.StatusBadRequest, message, errors)
}

func Unauthorized(
	c *gin.Context,
	message string,
) {
	Error(c, http.StatusUnauthorized, message, nil)
}

func Forbidden(
	c *gin.Context,
	message string,
) {
	Error(c, http.StatusForbidden, message, nil)
}

func NotFound(
	c *gin.Context,
	message string,
) {
	Error(c, http.StatusNotFound, message, nil)
}

func Conflict(
	c *gin.Context,
	message string,
	errors interface{},
) {
	Error(c, http.StatusConflict, message, errors)
}

func InternalServerError(
	c *gin.Context,
	message string,
) {
	Error(c, http.StatusInternalServerError, message, nil)
}
