package models

import (
	"notezy-backend/app/util"
	"regexp"

	"github.com/go-playground/validator"
)

// initialize the validator to validate the inputs, dtos
var Validator = validator.New()

func RegisterUserModelFieldsValidators(validate *validator.Validate) {
	validate.RegisterValidation("isstrongpassword", func(fl validator.FieldLevel) bool {
		password := fl.Field().String()
		re := regexp.MustCompile(`^(?=.*[A-Za-z])(?=.*\d)(?=.*[^\w\s]).{8,32}$`)
		return re.MatchString(password)
	})
}

func RegisterEnumValidators(validate *validator.Validate) {
	validate.RegisterValidation("isrole", func(fl validator.FieldLevel) bool {
		val := fl.Field().String()
		return util.IsStringIn(val, AllUserRoleStrings)
	})

	validate.RegisterValidation("isplan", func(fl validator.FieldLevel) bool {
		val := fl.Field().String()
		return util.IsStringIn(val, AllUserPlanStrings)
	})

	validate.RegisterValidation("isstatus", func(fl validator.FieldLevel) bool {
		val := fl.Field().String()
		return util.IsStringIn(val, AllUserStatusStrings)
	})

}

func init() {
	RegisterUserModelFieldsValidators(Validator)
	RegisterEnumValidators(Validator)
}
