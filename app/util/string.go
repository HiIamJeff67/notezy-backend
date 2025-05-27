package util

import "strings"

func JoinValues(values []string) string {
    return strings.Join(values, "', '")
}