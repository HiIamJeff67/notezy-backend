package validation

import (
	"net/url"
	"regexp"
	"strings"

	"github.com/go-playground/validator/v10" // make sure we use the version 10

	util "notezy-backend/app/util"
	constants "notezy-backend/shared/constants"
)

func RegisterStringsValidation(validate *validator.Validate) {
	validate.RegisterValidation("isaccount", func(fl validator.FieldLevel) bool {
		val := fl.Field().String()

		// try email validation
		if util.IsEmailString(val) {
			return true
		}

		// try alphaandnum validation
		hasLetter := regexp.MustCompile(`[a-zA-Z]`).MatchString(val)
		hasDigit := regexp.MustCompile(`\d`).MatchString(val)

		return hasLetter && hasDigit
	})
	validate.RegisterValidation("alphaandnum", func(fl validator.FieldLevel) bool {
		val := fl.Field().String()
		return util.IsAlphaAndNumberString(val)
	})
	validate.RegisterValidation("isstrongpassword", func(fl validator.FieldLevel) bool {
		password := strings.TrimSpace(fl.Field().String())
		if len(password) < constants.MinPasswordLength || len(password) > constants.MaxPasswordLength {
			return false
		}
		hasUpperCaseLetter := regexp.MustCompile(`[A-Z]`).MatchString(password)
		hasLowerCaseLetter := regexp.MustCompile(`[a-z]`).MatchString(password)
		hasDigit := regexp.MustCompile(`\d`).MatchString(password)
		hasSpecialCharacter := regexp.MustCompile(`[^\w\s]`).MatchString(password)
		return hasUpperCaseLetter && hasLowerCaseLetter && hasDigit && hasSpecialCharacter
	})
	validate.RegisterValidation("isuseragent", func(fl validator.FieldLevel) bool {
		userAgentStr := strings.TrimSpace(fl.Field().String())
		if len(userAgentStr) < 3 || len(userAgentStr) > constants.MaxUserAgentLength {
			return false
		}

		// check if the userAgent contain some malicious content
		if strings.Contains(userAgentStr, "<script>") ||
			strings.Contains(userAgentStr, "javascript:") ||
			strings.Contains(userAgentStr, "data:") {
			return false
		}

		return true
	})
	validate.RegisterValidation("isnumberstring", func(fl validator.FieldLevel) bool {
		val := fl.Field().String()
		return util.IsNumberString(val)
	})
	validate.RegisterValidation("isurl", func(fl validator.FieldLevel) bool {
		urlStr := strings.TrimSpace(fl.Field().String())
		if len(urlStr) > constants.MaxURLLength {
			return false
		}
		parsedURL, err := url.Parse(urlStr)
		if err != nil {
			return false
		}

		scheme := strings.ToLower(parsedURL.Scheme)
		for _, validScheme := range constants.URLWhiteList {
			if scheme == validScheme {
				return true
			}
		}
		for _, invalidScheme := range constants.URLBlackList {
			if scheme == invalidScheme {
				return false
			}
		}

		return parsedURL.Scheme != "" && parsedURL.Host != ""
	})
	validate.RegisterValidation("isimageurl", func(fl validator.FieldLevel) bool {
		urlStr := strings.TrimSpace(fl.Field().String())
		if len(urlStr) > constants.MaxURLLength {
			return false
		}
		parsedURL, err := url.Parse(urlStr)
		if err != nil {
			return false
		}
		scheme := strings.ToLower(parsedURL.Scheme)
		if scheme == "http" || scheme == "https" {
			return true
		}

		return scheme != "" && parsedURL.Host != ""
	})
	validate.RegisterValidation("isprogramminglanguage", func(fl validator.FieldLevel) bool {
		programmingLanguageStr := fl.Field().String()
		if len(programmingLanguageStr) > constants.MaxProgrammingLanguageLength {
			return false
		}
		for _, l := range constants.SupportedProgrammingLanguages {
			if programmingLanguageStr == l {
				return true
			}
		}
		return false
	})
	validate.RegisterValidation("ishexcode", func(fl validator.FieldLevel) bool {
		hexCode := fl.Field().String()
		if !strings.HasPrefix(hexCode, "#") {
			return false
		}
		cleanHexCode := hexCode[1:]
		length := len(cleanHexCode)
		if length != 3 && length != 4 && length != 6 && length != 8 {
			return false
		}

		return regexp.MustCompile(`^[0-9a-fA-F]+$`).MatchString(cleanHexCode)
	})
}
