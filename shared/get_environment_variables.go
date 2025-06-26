package shared

import (
	"os"
	"strconv"
)

/* ============================== Get Environment Variables Function ============================== */
func GetEnv(key string, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

func GetIntEnv(key string, fallback int) int {
	if valueStr, ok := os.LookupEnv(key); ok {
		value, err := strconv.Atoi(valueStr)
		if err != nil {
			return fallback
		}
		return value
	}
	return fallback
}
