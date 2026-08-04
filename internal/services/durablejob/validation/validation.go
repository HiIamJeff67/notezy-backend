package validation

import (
	"github.com/go-playground/validator/v10"

	blocknotevalidations "github.com/HiIamJeff67/notezy-backend/shared/lib/blocknote/validations"
	sharedvalidation "github.com/HiIamJeff67/notezy-backend/shared/validations"
)

func New() *validator.Validate {
	validator := validator.New()
	blocknotevalidations.RegisterShelfBlockValidation(validator)
	sharedvalidation.RegisterStringsValidation(validator)
	sharedvalidation.RegisterTimesValidation(validator)
	RegisterEnumsValidation(validator)
	return validator
}
