package validation

import (
	"github.com/go-playground/validator/v10" // make sure we use the version 10

	blocknote "notezy-backend/shared/lib/blocknote"
)

// initialize the validator to validate the inputs, dtos
var Validator = validator.New()

func init() {
	RegisterStringsValidation(Validator)
	RegisterEnumsValidation(Validator)
	RegisterTimesValidation(Validator)
	blocknote.RegisterShelfBlockValidation(Validator)
}
