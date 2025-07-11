package util

import (
	"strings"
	"unicode"
)

func JoinValues(values []string) string {
	return strings.Join(values, "', '")
}

func ConvertCamelCaseToSentenceCase(camelCaseString string) string {
	var result []rune
	for index, r := range camelCaseString {
		if unicode.IsUpper(r) && index != 0 {
			result = append(result, ' ')
		}
		result = append(result, unicode.ToLower(r))
	}
	return string(result)
}

func IsStringIn(s string, strs []string) bool {
	for _, str := range strs {
		if s == str {
			return true
		}
	}
	return false
}

func IsEmailString(s string) bool {
	return strings.Contains(s, "@")
}

func IsAlphaNumberString(s string) bool {
	for _, r := range s {
		if !unicode.IsLetter(r) && !unicode.IsDigit(r) {
			return false
		}
	}
	return true
}
