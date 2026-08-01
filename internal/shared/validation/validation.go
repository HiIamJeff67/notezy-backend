package validation

import "github.com/go-playground/validator/v10"

func New() *validator.Validate {
	validate := validator.New()
	RegisterStringsValidation(validate)
	RegisterTimesValidation(validate)

	return validate
}
