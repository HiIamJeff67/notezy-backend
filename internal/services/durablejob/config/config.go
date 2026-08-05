package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

type Config struct {
	YjsDocumentInitializationWorkerUrl string
	ListenAddress                      string
	KafkaConsumer                      KafkaConsumerConfig
}

func LoadConfig() (Config, error) {
	config := Config{
		YjsDocumentInitializationWorkerUrl: strings.TrimRight(strings.TrimSpace(os.Getenv("YJS_DOCUMENT_INITIALIZATION_WORKER_URL")), "/"),
		ListenAddress:                      strings.TrimSpace(os.Getenv("DURABLEJOB_LISTEN_ADDRESS")),
	}
	if config.ListenAddress == "" || config.YjsDocumentInitializationWorkerUrl == "" {
		return Config{}, fmt.Errorf("DURABLEJOB_LISTEN_ADDRESS and YJS_DOCUMENT_INITIALIZATION_WORKER_URL are required")
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

	return config, nil
}
