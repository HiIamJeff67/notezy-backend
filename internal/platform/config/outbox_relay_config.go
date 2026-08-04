package configs

import (
	"os"
	"strconv"
	"time"
)

type OutboxRelayConfig struct {
	BatchSize       int
	PollInterval    time.Duration
	ClaimTimeout    time.Duration
	InitialBackoff  time.Duration
	MaximumBackoff  time.Duration
	Retention       time.Duration
	CleanupInterval time.Duration
}

func OutboxRelay() OutboxRelayConfig {
	batchSize, err := strconv.Atoi(os.Getenv("OUTBOX_RELAY_BATCH_SIZE"))
	if err != nil || batchSize <= 0 {
		batchSize = 100
	}

	pollIntervalMilliseconds, err := strconv.Atoi(os.Getenv("OUTBOX_RELAY_POLL_INTERVAL_MILLISECONDS"))
	if err != nil || pollIntervalMilliseconds <= 0 {
		pollIntervalMilliseconds = 1000
	}
	claimTimeoutSeconds, err := strconv.Atoi(os.Getenv("OUTBOX_RELAY_CLAIM_TIMEOUT_SECONDS"))
	if err != nil || claimTimeoutSeconds <= 0 {
		claimTimeoutSeconds = 30
	}
	initialBackoffMilliseconds, err := strconv.Atoi(os.Getenv("OUTBOX_RELAY_INITIAL_BACKOFF_MILLISECONDS"))
	if err != nil || initialBackoffMilliseconds <= 0 {
		initialBackoffMilliseconds = 1000
	}
	maximumBackoffSeconds, err := strconv.Atoi(os.Getenv("OUTBOX_RELAY_MAXIMUM_BACKOFF_SECONDS"))
	if err != nil || maximumBackoffSeconds <= 0 {
		maximumBackoffSeconds = 60
	}
	retentionHours, err := strconv.Atoi(os.Getenv("OUTBOX_RELAY_RETENTION_HOURS"))
	if err != nil || retentionHours <= 0 {
		retentionHours = 168
	}
	cleanupIntervalMinutes, err := strconv.Atoi(os.Getenv("OUTBOX_RELAY_CLEANUP_INTERVAL_MINUTES"))
	if err != nil || cleanupIntervalMinutes <= 0 {
		cleanupIntervalMinutes = 60
	}

	return OutboxRelayConfig{
		BatchSize:       batchSize,
		PollInterval:    time.Duration(pollIntervalMilliseconds) * time.Millisecond,
		ClaimTimeout:    time.Duration(claimTimeoutSeconds) * time.Second,
		InitialBackoff:  time.Duration(initialBackoffMilliseconds) * time.Millisecond,
		MaximumBackoff:  time.Duration(maximumBackoffSeconds) * time.Second,
		Retention:       time.Duration(retentionHours) * time.Hour,
		CleanupInterval: time.Duration(cleanupIntervalMinutes) * time.Minute,
	}
}
