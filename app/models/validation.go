package models

import (
	logs "notezy-backend/app/logs"
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
		logs.FInfo("字母檢查: %v", hasLetter)
		logs.FInfo("數字檢查: %v", hasDigit)
		logs.FInfo("特殊字元檢查: %v", hasSpecialCharacter)
		return hasLetter && hasDigit && hasSpecialCharacter
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
	validate.RegisterValidation("isgender", func(fl validator.FieldLevel) bool {
		val := fl.Field().String()
		return util.IsStringIn(val, AllUserGenderStrings)
	})
	validate.RegisterValidation("iscountry", func(fl validator.FieldLevel) bool {
		val := fl.Field().String()
		return util.IsStringIn(val, AllCountryStrings)
	})
	validate.RegisterValidation("iscountrycode", func(fl validator.FieldLevel) bool {
		val := fl.Field().String()
		return util.IsStringIn(val, AllCountryCodeStrings)
	})
	validate.RegisterValidation("istheme", func(fl validator.FieldLevel) bool {
		val := fl.Field().String()
		return util.IsStringIn(val, AllThemeStrings)
	})
	validate.RegisterValidation("islanguage", func(fl validator.FieldLevel) bool {
		val := fl.Field().String()
		return util.IsStringIn(val, AllLanguageStrings)
	})
	validate.RegisterValidation("isbadgetype", func(fl validator.FieldLevel) bool {
		val := fl.Field().String()
		return util.IsStringIn(val, AllBadgeTypeStrings)
	})
}

func init() {
	RegisterUserModelFieldsValidators(Validator)
	RegisterEnumValidators(Validator)
}
