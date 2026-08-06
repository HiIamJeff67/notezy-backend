package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

type Config struct {
	ListenAddress     string
	TrustedProxies    []string
	AllowedDomains    []string
	RealtimeEnabled   bool
	BetaUserPublicIds []string
	YjsWorkerUrls     []string
	KafkaConsumer     KafkaConsumerConfig
	Redis             RedisConfig
}

type KafkaConsumerConfig struct {
	MaximumAttempts     int
	InitialRetryBackoff time.Duration
	MaximumRetryBackoff time.Duration
	MaximumPollRecords  int
}

func LoadConfig() (Config, error) {
	realtimeEnabled, err := strconv.ParseBool(strings.TrimSpace(os.Getenv("REALTIME_ENABLED")))
	if err != nil {
		return Config{}, fmt.Errorf("REALTIME_ENABLED must be a boolean")
	}
	config := Config{
		ListenAddress:     strings.TrimSpace(os.Getenv("REALTIME_GATEWAY_LISTEN_ADDRESS")),
		TrustedProxies:    splitValues(os.Getenv("GIN_TRUSTED_PROXIES")),
		AllowedDomains:    splitValues(os.Getenv("ALLOWED_DOMAINS")),
		RealtimeEnabled:   realtimeEnabled,
		BetaUserPublicIds: splitValues(os.Getenv("REALTIME_BETA_USER_PUBLIC_IDS")),
		YjsWorkerUrls:     splitValues(os.Getenv("YJS_WORKER_URLS")),
	}
	if config.ListenAddress == "" || len(config.YjsWorkerUrls) == 0 {
		return Config{}, fmt.Errorf("REALTIME_GATEWAY_LISTEN_ADDRESS and YJS_WORKER_URLS are required")
	}
	maximumAttempts, err := strconv.Atoi(strings.TrimSpace(os.Getenv("KAFKA_CONSUMER_MAXIMUM_ATTEMPTS")))
	if err != nil || maximumAttempts <= 0 {
		return Config{}, fmt.Errorf("KAFKA_CONSUMER_MAXIMUM_ATTEMPTS must be a positive integer")
	}
	initialRetryBackoff, err := time.ParseDuration(strings.TrimSpace(os.Getenv("KAFKA_CONSUMER_INITIAL_RETRY_BACKOFF")))
	if err != nil || initialRetryBackoff <= 0 {
		return Config{}, fmt.Errorf("KAFKA_CONSUMER_INITIAL_RETRY_BACKOFF must be a positive Go duration")
	}
	maximumRetryBackoff, err := time.ParseDuration(strings.TrimSpace(os.Getenv("KAFKA_CONSUMER_MAXIMUM_RETRY_BACKOFF")))
	if err != nil || maximumRetryBackoff <= 0 {
		return Config{}, fmt.Errorf("KAFKA_CONSUMER_MAXIMUM_RETRY_BACKOFF must be a positive Go duration")
	}
	if maximumRetryBackoff < initialRetryBackoff {
		return Config{}, fmt.Errorf("KAFKA_CONSUMER_MAXIMUM_RETRY_BACKOFF must not be smaller than KAFKA_CONSUMER_INITIAL_RETRY_BACKOFF")
	}
	maximumPollRecords, err := strconv.Atoi(strings.TrimSpace(os.Getenv("KAFKA_CONSUMER_MAXIMUM_POLL_RECORDS")))
	if err != nil || maximumPollRecords <= 0 {
		return Config{}, fmt.Errorf("KAFKA_CONSUMER_MAXIMUM_POLL_RECORDS must be a positive integer")
	}
	config.KafkaConsumer = KafkaConsumerConfig{
		MaximumAttempts:     maximumAttempts,
		InitialRetryBackoff: initialRetryBackoff,
		MaximumRetryBackoff: maximumRetryBackoff,
		MaximumPollRecords:  maximumPollRecords,
	}
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
