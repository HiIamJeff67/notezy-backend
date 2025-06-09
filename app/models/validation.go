package models

import (
	"notezy-backend/app/models/enums"
	"notezy-backend/app/util"
	"regexp"

	"github.com/go-playground/validator"
)

// initialize the validator to validate the inputs, dtos
var Validator = validator.New()

func RegisterUserModelFieldsValidators(validate *validator.Validate) {
	validate.RegisterValidation("isstrongpassword", func(fl validator.FieldLevel) bool {
		password := fl.Field().String()
		if len(password) < 8 || len(password) > 32 {
			return false
		}
		hasLetter := regexp.MustCompile(`[A-Za-z]`).MatchString(password)
		hasDigit := regexp.MustCompile(`\d`).MatchString(password)
		hasSpecialCharacter := regexp.MustCompile(`[^\w\s]`).MatchString(password)
		return hasLetter && hasDigit && hasSpecialCharacter
	})
}

func RegisterEnumValidators(validate *validator.Validate) {
	validate.RegisterValidation("isrole", func(fl validator.FieldLevel) bool {
		val := fl.Field().String()
		return util.IsStringIn(val, enums.AllUserRoleStrings)
	})
	validate.RegisterValidation("isplan", func(fl validator.FieldLevel) bool {
		val := fl.Field().String()
		return util.IsStringIn(val, enums.AllUserPlanStrings)
	})
	validate.RegisterValidation("isstatus", func(fl validator.FieldLevel) bool {
		val := fl.Field().String()
		return util.IsStringIn(val, enums.AllUserStatusStrings)
	})
	validate.RegisterValidation("isgender", func(fl validator.FieldLevel) bool {
		val := fl.Field().String()
		return util.IsStringIn(val, enums.AllUserGenderStrings)
	})
	validate.RegisterValidation("iscountry", func(fl validator.FieldLevel) bool {
		val := fl.Field().String()
		return util.IsStringIn(val, enums.AllCountryStrings)
	})
	validate.RegisterValidation("iscountrycode", func(fl validator.FieldLevel) bool {
		val := fl.Field().String()
		return util.IsStringIn(val, enums.AllCountryCodeStrings)
	})
	validate.RegisterValidation("islanguage", func(fl validator.FieldLevel) bool {
		val := fl.Field().String()
		return util.IsStringIn(val, enums.AllLanguageStrings)
	})
	validate.RegisterValidation("isbadgetype", func(fl validator.FieldLevel) bool {
		val := fl.Field().String()
		return util.IsStringIn(val, enums.AllBadgeTypeStrings)
	})
}

func init() {
	RegisterUserModelFieldsValidators(Validator)
	RegisterEnumValidators(Validator)
}
