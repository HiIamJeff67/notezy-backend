package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

type Config struct {
	ListenAddress      string
	OAuthGoogle        OAuthGoogleConfig
	OutboxRelay        OutboxRelayConfig
	KafkaConsumer      KafkaConsumerConfig
	StorageKeySalt     string
	EmailBaseUrl       string
	EmailClientTimeout time.Duration
}

func LoadConfig() (Config, error) {
	listenAddress := strings.TrimSpace(os.Getenv("CORE_LISTEN_ADDRESS"))
	if listenAddress == "" {
		return Config{}, fmt.Errorf("CORE_LISTEN_ADDRESS is required")
	}
	oauthGoogle := OAuthGoogleConfig{
		ClientId:     strings.TrimSpace(os.Getenv("OAUTH_GOOGLE_CLIENT_ID")),
		ClientSecret: os.Getenv("OAUTH_GOOGLE_CLIENT_SECRET"),
		RedirectUrl:  strings.TrimSpace(os.Getenv("OAUTH_GOOGLE_REDIRECT_URL")),
	}
	if oauthGoogle.ClientId == "" || oauthGoogle.ClientSecret == "" || oauthGoogle.RedirectUrl == "" {
		return Config{}, fmt.Errorf("OAUTH_GOOGLE_CLIENT_ID, OAUTH_GOOGLE_CLIENT_SECRET, and OAUTH_GOOGLE_REDIRECT_URL are required")
	}

	outboxRelay := OutboxRelayConfig{}
	batchSize, err := strconv.Atoi(strings.TrimSpace(os.Getenv("OUTBOX_RELAY_BATCH_SIZE")))
	if err != nil || batchSize <= 0 {
		return Config{}, fmt.Errorf("OUTBOX_RELAY_BATCH_SIZE must be a positive integer")
	}
	outboxRelay.BatchSize = batchSize
	outboxRelay.PollInterval, err = time.ParseDuration(strings.TrimSpace(os.Getenv("OUTBOX_RELAY_POLL_INTERVAL")))
	if err != nil || outboxRelay.PollInterval <= 0 {
		return Config{}, fmt.Errorf("OUTBOX_RELAY_POLL_INTERVAL must be a positive Go duration")
	}
	outboxRelay.ClaimTimeout, err = time.ParseDuration(strings.TrimSpace(os.Getenv("OUTBOX_RELAY_CLAIM_TIMEOUT")))
	if err != nil || outboxRelay.ClaimTimeout <= 0 {
		return Config{}, fmt.Errorf("OUTBOX_RELAY_CLAIM_TIMEOUT must be a positive Go duration")
	}
	outboxRelay.InitialBackoff, err = time.ParseDuration(strings.TrimSpace(os.Getenv("OUTBOX_RELAY_INITIAL_BACKOFF")))
	if err != nil || outboxRelay.InitialBackoff <= 0 {
		return Config{}, fmt.Errorf("OUTBOX_RELAY_INITIAL_BACKOFF must be a positive Go duration")
	}
	outboxRelay.MaximumBackoff, err = time.ParseDuration(strings.TrimSpace(os.Getenv("OUTBOX_RELAY_MAXIMUM_BACKOFF")))
	if err != nil || outboxRelay.MaximumBackoff <= 0 {
		return Config{}, fmt.Errorf("OUTBOX_RELAY_MAXIMUM_BACKOFF must be a positive Go duration")
	}
	if outboxRelay.MaximumBackoff < outboxRelay.InitialBackoff {
		return Config{}, fmt.Errorf("OUTBOX_RELAY_MAXIMUM_BACKOFF must not be smaller than OUTBOX_RELAY_INITIAL_BACKOFF")
	}
	outboxRelay.Retention, err = time.ParseDuration(strings.TrimSpace(os.Getenv("OUTBOX_RELAY_RETENTION")))
	if err != nil || outboxRelay.Retention <= 0 {
		return Config{}, fmt.Errorf("OUTBOX_RELAY_RETENTION must be a positive Go duration")
	}
	outboxRelay.CleanupInterval, err = time.ParseDuration(strings.TrimSpace(os.Getenv("OUTBOX_RELAY_CLEANUP_INTERVAL")))
	if err != nil || outboxRelay.CleanupInterval <= 0 {
		return Config{}, fmt.Errorf("OUTBOX_RELAY_CLEANUP_INTERVAL must be a positive Go duration")
	}

	kafkaConsumer := KafkaConsumerConfig{}
	kafkaConsumer.MaximumAttempts, err = strconv.Atoi(strings.TrimSpace(os.Getenv("KAFKA_CONSUMER_MAXIMUM_ATTEMPTS")))
	if err != nil || kafkaConsumer.MaximumAttempts <= 0 {
		return Config{}, fmt.Errorf("KAFKA_CONSUMER_MAXIMUM_ATTEMPTS must be a positive integer")
	}
	kafkaConsumer.InitialRetryBackoff, err = time.ParseDuration(strings.TrimSpace(os.Getenv("KAFKA_CONSUMER_INITIAL_RETRY_BACKOFF")))
	if err != nil || kafkaConsumer.InitialRetryBackoff <= 0 {
		return Config{}, fmt.Errorf("KAFKA_CONSUMER_INITIAL_RETRY_BACKOFF must be a positive Go duration")
	}
	kafkaConsumer.MaximumRetryBackoff, err = time.ParseDuration(strings.TrimSpace(os.Getenv("KAFKA_CONSUMER_MAXIMUM_RETRY_BACKOFF")))
	if err != nil || kafkaConsumer.MaximumRetryBackoff <= 0 {
		return Config{}, fmt.Errorf("KAFKA_CONSUMER_MAXIMUM_RETRY_BACKOFF must be a positive Go duration")
	}
	if kafkaConsumer.MaximumRetryBackoff < kafkaConsumer.InitialRetryBackoff {
		return Config{}, fmt.Errorf("KAFKA_CONSUMER_MAXIMUM_RETRY_BACKOFF must not be smaller than KAFKA_CONSUMER_INITIAL_RETRY_BACKOFF")
	}
	kafkaConsumer.MaximumPollRecords, err = strconv.Atoi(strings.TrimSpace(os.Getenv("KAFKA_CONSUMER_MAXIMUM_POLL_RECORDS")))
	if err != nil || kafkaConsumer.MaximumPollRecords <= 0 {
		return Config{}, fmt.Errorf("KAFKA_CONSUMER_MAXIMUM_POLL_RECORDS must be a positive integer")
	}
	storageKeySalt := os.Getenv("STORAGE_KEY_SALT")
	if storageKeySalt == "" {
		return Config{}, fmt.Errorf("STORAGE_KEY_SALT is required")
	}
	emailBaseUrl := strings.TrimRight(strings.TrimSpace(os.Getenv("EMAIL_SERVICE_BASE_URL")), "/")
	if emailBaseUrl == "" {
		return Config{}, fmt.Errorf("EMAIL_SERVICE_BASE_URL is required")
	}
	emailClientTimeout, err := time.ParseDuration(strings.TrimSpace(os.Getenv("EMAIL_SERVICE_CLIENT_TIMEOUT")))
	if err != nil || emailClientTimeout <= 0 {
		return Config{}, fmt.Errorf("EMAIL_SERVICE_CLIENT_TIMEOUT must be a positive Go duration")
	}

	return Config{
		ListenAddress:      listenAddress,
		OAuthGoogle:        oauthGoogle,
		OutboxRelay:        outboxRelay,
		KafkaConsumer:      kafkaConsumer,
		StorageKeySalt:     storageKeySalt,
		EmailBaseUrl:       emailBaseUrl,
		EmailClientTimeout: emailClientTimeout,
	}, nil
}
