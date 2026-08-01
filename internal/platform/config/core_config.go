package configs

import (
	"os"
	"strconv"
	"time"
)

func CoreBaseURL() string {
	return os.Getenv("CORE_BASE_URL")
}

func CoreListenAddress() string {
	return os.Getenv("CORE_LISTEN_ADDRESS")
}

func CoreClientTimeout() time.Duration {
	timeoutSeconds, _ := strconv.Atoi(os.Getenv("CORE_CLIENT_TIMEOUT_SECONDS"))
	return time.Duration(timeoutSeconds) * time.Second
}
