package main

import "strings"

func kebabCase(value string) string {
	if value == "" {
		return value
	}

	var output strings.Builder
	for index, current := range value {
		if current == '_' || current == '-' || current == ' ' {
			if output.Len() > 0 && !strings.HasSuffix(output.String(), "-") {
				output.WriteByte('-')
			}
			continue
		}
		if current >= 'A' && current <= 'Z' {
			if index > 0 && !strings.HasSuffix(output.String(), "-") {
				output.WriteByte('-')
			}
			current += 'a' - 'A'
		}
		output.WriteRune(current)
	}
	return strings.Trim(output.String(), "-")
}

func parameterVariableName(pathName string) string {
	parts := strings.Split(pathName, "-")
	if len(parts) == 1 {
		return pathName
	}
	for index := 1; index < len(parts); index++ {
		if parts[index] == "" {
			continue
		}
		parts[index] = strings.ToUpper(parts[index][:1]) + parts[index][1:]
	}
	return strings.Join(parts, "")
}

func matchesPathParameter(fieldName, pathName string) bool {
	return fieldName == pathName || kebabCase(fieldName) == pathName
}
