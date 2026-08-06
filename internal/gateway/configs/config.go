package config

import (
	"fmt"
	"os"
	"strings"
	"time"
)

type Config struct {
	ListenAddress                string
	TrustedProxies               []string
	AllowedDomains               []string
	CoreBaseUrl                  string
	CoreClientTimeout            time.Duration
	RealtimeGatewayBaseUrl       string
	RealtimeGatewayClientTimeout time.Duration
	Redis                        RedisConfig
}

func LoadConfig() (Config, error) {
	config := Config{
		ListenAddress:          strings.TrimSpace(os.Getenv("GATEWAY_LISTEN_ADDRESS")),
		TrustedProxies:         splitValues(os.Getenv("GIN_TRUSTED_PROXIES")),
		AllowedDomains:         splitValues(os.Getenv("ALLOWED_DOMAINS")),
		CoreBaseUrl:            strings.TrimRight(strings.TrimSpace(os.Getenv("CORE_BASE_URL")), "/"),
		RealtimeGatewayBaseUrl: strings.TrimRight(strings.TrimSpace(os.Getenv("REALTIME_GATEWAY_BASE_URL")), "/"),
	}
	if config.ListenAddress == "" || config.CoreBaseUrl == "" || config.RealtimeGatewayBaseUrl == "" {
		return Config{}, fmt.Errorf("GATEWAY_LISTEN_ADDRESS, CORE_BASE_URL, and REALTIME_GATEWAY_BASE_URL are required")
	}
	coreTimeout, err := time.ParseDuration(strings.TrimSpace(os.Getenv("CORE_CLIENT_TIMEOUT")))
	if err != nil || coreTimeout <= 0 {
		return Config{}, fmt.Errorf("CORE_CLIENT_TIMEOUT must be a positive Go duration")
	}
	realtimeGatewayTimeout, err := time.ParseDuration(strings.TrimSpace(os.Getenv("REALTIME_GATEWAY_CLIENT_TIMEOUT")))
	if err != nil || realtimeGatewayTimeout <= 0 {
		return Config{}, fmt.Errorf("REALTIME_GATEWAY_CLIENT_TIMEOUT must be a positive Go duration")
	}

	config.CoreClientTimeout = coreTimeout
	config.RealtimeGatewayClientTimeout = realtimeGatewayTimeout
	redisConfig, err := loadRedisConfig()
	if err != nil {
		return Config{}, err
	}
	config.Redis = redisConfig

	return config, nil
}

func splitValues(value string) []string {
	values := strings.Split(value, ",")
	result := make([]string, 0, len(values))
	for _, item := range values {
		if trimmedItem := strings.TrimSpace(item); trimmedItem != "" {
			result = append(result, trimmedItem)
		}
	}

	return result
}
