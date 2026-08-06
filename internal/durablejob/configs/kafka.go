package config

import "time"

type KafkaConsumerConfig struct {
	MaximumAttempts     int
	InitialRetryBackoff time.Duration
	MaximumRetryBackoff time.Duration
	MaximumPollRecords  int
}
