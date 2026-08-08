package strings

import (
	stdstrings "strings"
	"unicode"
)

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

func SplitValues(value string) []string {
	values := stdstrings.Split(value, ",")
	result := make([]string, 0, len(values))
	for _, item := range values {
		if trimmedItem := stdstrings.TrimSpace(item); trimmedItem != "" {
			result = append(result, trimmedItem)
		}
	}

	return result
}
