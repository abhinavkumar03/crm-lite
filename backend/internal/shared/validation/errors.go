package validation

import (
	"strings"

	"github.com/go-playground/validator/v10"

	"github.com/abhinavkumar03/crm-lite/backend/internal/shared/response"
)

func FormatErrors(err error) []response.ValidationError {

	var validationErrors []response.ValidationError

	if err == nil {
		return validationErrors
	}

	for _, fieldError := range err.(validator.ValidationErrors) {

		validationErrors = append(
			validationErrors,
			response.ValidationError{
				Field: strings.ToLower(fieldError.Field()),
				Message: validationMessage(
					fieldError,
				),
			},
		)
	}

	return validationErrors
}

func validationMessage(
	field validator.FieldError,
) string {

	switch field.Tag() {

	case "required":
		return "This field is required"

	case "email":
		return "Must be a valid email"

	case "min":
		return "Value is too short"

	case "max":
		return "Value is too long"

	case "crmname":
		return "Must contain at least 3 non-space characters"

	default:
		return "Invalid value"
	}
}
