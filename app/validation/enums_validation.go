package validation

import (
	"github.com/go-playground/validator/v10" // make sure we use the version 10

	enums "notezy-backend/app/models/schemas/enums"
	util "notezy-backend/app/util"
)

func RegisterEnumsValidation(validate *validator.Validate) {
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
	validate.RegisterValidation("ismaterialtype", func(fl validator.FieldLevel) bool {
		val := fl.Field().String()
		return util.IsStringIn(val, enums.AllMaterialTypeStrings)
	})
	validate.RegisterValidation("ismaterialcontenttype", func(fl validator.FieldLevel) bool {
		val := fl.Field().String()
		return util.IsStringIn(val, enums.AllMaterialContentTypeStrings)
	})
}
