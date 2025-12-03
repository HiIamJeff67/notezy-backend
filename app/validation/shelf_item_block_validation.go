package validation

import (
	"regexp"

	"github.com/go-playground/validator/v10" // make sure we use the version 10

	constants "notezy-backend/shared/constants"
)

func RegisterShelfBlockValidation(validate *validator.Validate) {
	validate.RegisterValidation("isshelfname", func(fl validator.FieldLevel) bool {
		shelfNameStr := fl.Field().String()
		if len(shelfNameStr) > constants.MaxShelfNameLength {
			return false
		}
		return !regexp.MustCompile(`[\/\\:\*\?"<>\|]`).MatchString(shelfNameStr)
	})
	validate.RegisterValidation("isitemname", func(fl validator.FieldLevel) bool {
		itemNameStr := fl.Field().String()
		if len(itemNameStr) > constants.MaxItemNameLength {
			return false
		}
		return !regexp.MustCompile(`[\/\\:\*\?"<>\|]`).MatchString(itemNameStr)
	})
	validate.RegisterValidation("isfileblockname", func(fl validator.FieldLevel) bool {
		fileBlockNameStr := fl.Field().String()
		return len(fileBlockNameStr) <= constants.MaxFileBlockNameLength
	})
	validate.RegisterValidation("isfileblockcaption", func(fl validator.FieldLevel) bool {
		fileBlockCaptionStr := fl.Field().String()
		return len(fileBlockCaptionStr) <= constants.MaxFileBlockCaptionLength
	})
	validate.RegisterValidation("istextalignment", func(fl validator.FieldLevel) bool {
		textAlignmentStr := fl.Field().String()
		return textAlignmentStr == "left" || textAlignmentStr == "right" || textAlignmentStr == "center" || textAlignmentStr == "justify"
	})
	validate.RegisterValidation("isheadinglevel", func(fl validator.FieldLevel) bool {
		headingLevel := fl.Field().NumField()
		return headingLevel > 0 && headingLevel <= constants.MaxHeadingLevel
	})
}
