package configs

import (
	"os"
	"strconv"
	"time"
)

type KafkaConsumerConfig struct {
	MaximumAttempts     int
	InitialRetryBackoff time.Duration
	MaximumRetryBackoff time.Duration
	MaximumPollRecords  int
}

func KafkaConsumer() KafkaConsumerConfig {
	maximumAttempts, err := strconv.Atoi(os.Getenv("KAFKA_CONSUMER_MAXIMUM_ATTEMPTS"))
	if err != nil || maximumAttempts <= 0 {
		maximumAttempts = 3
	}
	initialRetryBackoffMilliseconds, err := strconv.Atoi(os.Getenv("KAFKA_CONSUMER_INITIAL_RETRY_BACKOFF_MILLISECONDS"))
	if err != nil || initialRetryBackoffMilliseconds <= 0 {
		initialRetryBackoffMilliseconds = 250
	}
	maximumRetryBackoffMilliseconds, err := strconv.Atoi(os.Getenv("KAFKA_CONSUMER_MAXIMUM_RETRY_BACKOFF_MILLISECONDS"))
	if err != nil || maximumRetryBackoffMilliseconds < initialRetryBackoffMilliseconds {
		maximumRetryBackoffMilliseconds = 5000
	}
	maximumPollRecords, err := strconv.Atoi(os.Getenv("KAFKA_CONSUMER_MAXIMUM_POLL_RECORDS"))
	if err != nil || maximumPollRecords <= 0 {
		maximumPollRecords = 100
	}

	return KafkaConsumerConfig{
		MaximumAttempts:     maximumAttempts,
		InitialRetryBackoff: time.Duration(initialRetryBackoffMilliseconds) * time.Millisecond,
		MaximumRetryBackoff: time.Duration(maximumRetryBackoffMilliseconds) * time.Millisecond,
		MaximumPollRecords:  maximumPollRecords,
	}
}
