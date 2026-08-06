package config

import "time"

type OutboxRelayConfig struct {
	BatchSize       int
	PollInterval    time.Duration
	ClaimTimeout    time.Duration
	InitialBackoff  time.Duration
	MaximumBackoff  time.Duration
	Retention       time.Duration
	CleanupInterval time.Duration
}
