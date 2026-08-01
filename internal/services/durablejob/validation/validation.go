package validation

import (
	"github.com/go-playground/validator/v10" // make sure we use the version 10

	blocknote "github.com/HiIamJeff67/notezy-backend/internal/shared/lib/blocknote"
	sharedvalidation "github.com/HiIamJeff67/notezy-backend/internal/shared/validation"
)

var Validator *validator.Validate

func init() {
	Validator = sharedvalidation.New()
	RegisterEnumsValidation(Validator)
	blocknote.RegisterShelfBlockValidation(Validator)
}
