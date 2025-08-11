package validation

import (
	"net/url"
	"notezy-backend/app/util"
	"regexp"
	"strings"

	"github.com/go-playground/validator/v10" // make sure we use the version 10
)

func RegisterStringsValidation(validate *validator.Validate) {
	validate.RegisterValidation("account", func(fl validator.FieldLevel) bool {
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
		password := fl.Field().String()
		if len(password) < 8 {
			return false
		}
		hasUpperCaseLetter := regexp.MustCompile(`[A-Z]`).MatchString(password)
		hasLowerCaseLetter := regexp.MustCompile(`[a-z]`).MatchString(password)
		hasDigit := regexp.MustCompile(`\d`).MatchString(password)
		hasSpecialCharacter := regexp.MustCompile(`[^\w\s]`).MatchString(password)
		return hasUpperCaseLetter && hasLowerCaseLetter && hasDigit && hasSpecialCharacter
	})
	validate.RegisterValidation("isuseragent", func(fl validator.FieldLevel) bool {
		userAgent := strings.TrimSpace(fl.Field().String())

		if len(userAgent) < 3 || len(userAgent) > 2000 {
			return false
		}

		// check if the userAgent contain some malicious content
		if strings.Contains(userAgent, "<script>") ||
			strings.Contains(userAgent, "javascript:") ||
			strings.Contains(userAgent, "data:") {
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
		if urlStr == "" {
			return false
		}

		parsedURL, err := url.Parse(urlStr)
		if err != nil {
			return false
		}

		return parsedURL.Scheme != "" && parsedURL.Host != ""
	})
	validate.RegisterValidation("isimageurl", func(fl validator.FieldLevel) bool {
		urlStr := strings.TrimSpace(fl.Field().String())

		parsedURL, err := url.Parse(urlStr)
		if err != nil || parsedURL.Scheme == "" || parsedURL.Host == "" {
			return false
		}

		imageExtensions := []string{".jpg", ".jpeg", ".png", ".gif", ".webp", ".svg"}
		path := strings.ToLower(parsedURL.Path)

		for _, ext := range imageExtensions {
			if strings.HasSuffix(path, ext) {
				return true
			}
		}

		return false
	})
}
