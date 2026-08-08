package config

import (
	"testing"
)

func TestLoadConfig(t *testing.T) {
	t.Setenv("DURABLEJOB_LISTEN_ADDRESS", "127.0.0.1:8082")
	t.Setenv("KAFKA_CONSUMER_MAXIMUM_ATTEMPTS", "3")
	t.Setenv("KAFKA_CONSUMER_INITIAL_RETRY_BACKOFF", "250ms")
	t.Setenv("KAFKA_CONSUMER_MAXIMUM_RETRY_BACKOFF", "5s")
	t.Setenv("KAFKA_CONSUMER_MAXIMUM_POLL_RECORDS", "100")
	t.Setenv("DURABLEJOB_YJS_MAINTENANCE_MAXIMUM_PENDING_HINTS", "1000")
	t.Setenv("DURABLEJOB_YJS_MAINTENANCE_MAXIMUM_DISPATCH_BATCH", "32")
	t.Setenv("DURABLEJOB_YJS_MAINTENANCE_MAXIMUM_DISPATCH_WORKERS", "8")
	t.Setenv("DURABLEJOB_YJS_MAINTENANCE_MAXIMUM_REQUEST_ATTEMPTS", "3")

	config, err := LoadConfig()
	if err != nil {
		t.Fatalf("LoadConfig() error = %v", err)
	}
	if config.KafkaConsumer.MaximumAttempts != 3 {
		t.Fatalf("LoadConfig() = %#v", config)
	}
}
