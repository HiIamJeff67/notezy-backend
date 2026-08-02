package configs

import (
	"os"
	"strconv"
	"time"
)

func CoreBaseURL() string {
	baseURL := os.Getenv("CORE_BASE_URL")
	if baseURL == "" {
		return "http://127.0.0.1:7778"
	}

	return baseURL
}

func CoreListenAddress() string {
	address := os.Getenv("CORE_LISTEN_ADDRESS")
	if address == "" {
		return "0.0.0.0:7778"
	}

	return address
}

func CoreClientTimeout() time.Duration {
	timeoutSeconds, err := strconv.Atoi(os.Getenv("CORE_CLIENT_TIMEOUT_SECONDS"))
	if err != nil || timeoutSeconds <= 0 {
		return 10 * time.Second
	}

	return time.Duration(timeoutSeconds) * time.Second
}
