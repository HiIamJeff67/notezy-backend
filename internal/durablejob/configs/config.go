package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

type YjsMaintenanceStrategyConfig struct {
	MaximumPendingHints    int
	MaximumDispatchBatch   int
	MaximumDispatchWorkers int
	MaximumRequestAttempts int
}

type Config struct {
	ListenAddress          string
	KafkaConsumer          KafkaConsumerConfig
	YjsMaintenanceStrategy YjsMaintenanceStrategyConfig
}

func LoadConfig() (Config, error) {
	config := Config{
		ListenAddress: strings.TrimSpace(os.Getenv("DURABLEJOB_LISTEN_ADDRESS")),
	}
	if config.ListenAddress == "" {
		return Config{}, fmt.Errorf("DURABLEJOB_LISTEN_ADDRESS is required")
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

	maximumPendingHints, err := strconv.Atoi(strings.TrimSpace(os.Getenv("DURABLEJOB_YJS_MAINTENANCE_MAXIMUM_PENDING_HINTS")))
	if err != nil || maximumPendingHints <= 0 {
		return Config{}, fmt.Errorf("DURABLEJOB_YJS_MAINTENANCE_MAXIMUM_PENDING_HINTS must be a positive integer")
	}
	maximumDispatchBatch, err := strconv.Atoi(strings.TrimSpace(os.Getenv("DURABLEJOB_YJS_MAINTENANCE_MAXIMUM_DISPATCH_BATCH")))
	if err != nil || maximumDispatchBatch <= 0 {
		return Config{}, fmt.Errorf("DURABLEJOB_YJS_MAINTENANCE_MAXIMUM_DISPATCH_BATCH must be a positive integer")
	}
	maximumDispatchWorkers, err := strconv.Atoi(strings.TrimSpace(os.Getenv("DURABLEJOB_YJS_MAINTENANCE_MAXIMUM_DISPATCH_WORKERS")))
	if err != nil || maximumDispatchWorkers <= 0 {
		return Config{}, fmt.Errorf("DURABLEJOB_YJS_MAINTENANCE_MAXIMUM_DISPATCH_WORKERS must be a positive integer")
	}
	maximumRequestAttempts, err := strconv.Atoi(strings.TrimSpace(os.Getenv("DURABLEJOB_YJS_MAINTENANCE_MAXIMUM_REQUEST_ATTEMPTS")))
	if err != nil || maximumRequestAttempts <= 0 {
		return Config{}, fmt.Errorf("DURABLEJOB_YJS_MAINTENANCE_MAXIMUM_REQUEST_ATTEMPTS must be a positive integer")
	}
	config.YjsMaintenanceStrategy = YjsMaintenanceStrategyConfig{
		MaximumPendingHints:    maximumPendingHints,
		MaximumDispatchBatch:   maximumDispatchBatch,
		MaximumDispatchWorkers: maximumDispatchWorkers,
		MaximumRequestAttempts: maximumRequestAttempts,
	}

	return config, nil
}
