package validation

import (
	"strings"

	"github.com/go-playground/validator/v10"
)

func registerCustomValidators(v *validator.Validate) {

	// Example:
	// "crmname" validates that the value is not just whitespace.

	_ = v.RegisterValidation("crmname", func(fl validator.FieldLevel) bool {

		value := strings.TrimSpace(fl.Field().String())

		return len(value) >= 3
	})
}
