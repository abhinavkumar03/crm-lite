package validation

import (
	"sync"

	"github.com/go-playground/validator/v10"
)

var (
	validate *validator.Validate
	once     sync.Once
)

// Validator returns a singleton validator instance.
func Validator() *validator.Validate {
	once.Do(func() {
		validate = validator.New()
		registerCustomValidators(validate)
	})

	return validate
}
