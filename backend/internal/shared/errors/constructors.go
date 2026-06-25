package errors

import "net/http"

func BadRequest(message string) *AppError {
	return &AppError{
		StatusCode: http.StatusBadRequest,
		Code:       CodeBadRequest,
		Message:    message,
	}
}

func Validation(message string) *AppError {
	return &AppError{
		StatusCode: http.StatusBadRequest,
		Code:       CodeValidation,
		Message:    message,
	}
}

func Unauthorized(message string) *AppError {
	return &AppError{
		StatusCode: http.StatusUnauthorized,
		Code:       CodeUnauthorized,
		Message:    message,
	}
}

func Forbidden(message string) *AppError {
	return &AppError{
		StatusCode: http.StatusForbidden,
		Code:       CodeForbidden,
		Message:    message,
	}
}

func NotFound(message string) *AppError {
	return &AppError{
		StatusCode: http.StatusNotFound,
		Code:       CodeNotFound,
		Message:    message,
	}
}

func Conflict(message string) *AppError {
	return &AppError{
		StatusCode: http.StatusConflict,
		Code:       CodeConflict,
		Message:    message,
	}
}

func Internal(err error) *AppError {
	return &AppError{
		StatusCode: http.StatusInternalServerError,
		Code:       CodeInternal,
		Message:    "Internal server error",
		Err:        err,
	}
}
