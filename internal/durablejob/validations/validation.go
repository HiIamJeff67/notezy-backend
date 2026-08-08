package validation

import (
	"github.com/go-playground/validator/v10"

	sharedvalidation "github.com/HiIamJeff67/notezy-backend/shared/validations"

	blocknotevalidations "github.com/HiIamJeff67/notezy-backend/contracts/types/blocknote/validations"
)

func New() *validator.Validate {
	validator := validator.New()
	blocknotevalidations.RegisterShelfBlockValidation(validator)
	sharedvalidation.RegisterStringsValidation(validator)
	sharedvalidation.RegisterTimesValidation(validator)
	RegisterEnumsValidation(validator)
	return validator
}
