package validators

import (
	"regexp"
	"unicode"
)

func IsEmailString(value string) bool {
	return regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`).MatchString(value)
}

func IsNumberString(value string) bool {
	for _, character := range value {
		if !unicode.IsDigit(character) {
			return false
		}
	}

	return true
}

func IsAlphaOrNumberString(value string) bool {
	for _, character := range value {
		if !unicode.IsLetter(character) && !unicode.IsDigit(character) {
			return false
		}
	}

	return true
}

func IsAlphaAndNumberString(value string) bool {
	var hasLetter bool
	var hasDigit bool

	for _, character := range value {
		if unicode.IsLetter(character) {
			hasLetter = true
		} else if unicode.IsDigit(character) {
			hasDigit = true
		} else {
			return false
		}
	}

	return hasLetter && hasDigit
}
