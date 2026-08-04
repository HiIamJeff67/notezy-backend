package config

import (
	"testing"
	"time"
)

func TestLoadConfig(t *testing.T) {
	t.Setenv("DURABLEJOB_LISTEN_ADDRESS", "127.0.0.1:8082")
	t.Setenv("CORE_BASE_URL", "http://core:7778")
	t.Setenv("CORE_CLIENT_TIMEOUT", "10s")
	t.Setenv("YJS_MAINTENANCE_WORKER_URL", "http://yjs/maintenance")
	t.Setenv("YJS_DOCUMENT_INITIALIZATION_WORKER_URL", "http://yjs/document-initialization")
	t.Setenv("YJS_PROJECTION_WORKER_URL", "http://yjs/projection")
	t.Setenv("KAFKA_CONSUMER_MAXIMUM_ATTEMPTS", "3")
	t.Setenv("KAFKA_CONSUMER_INITIAL_RETRY_BACKOFF", "250ms")
	t.Setenv("KAFKA_CONSUMER_MAXIMUM_RETRY_BACKOFF", "5s")
	t.Setenv("KAFKA_CONSUMER_MAXIMUM_POLL_RECORDS", "100")

	config, err := LoadConfig()
	if err != nil {
		t.Fatalf("LoadConfig() error = %v", err)
	}
	if config.CoreClientTimeout != 10*time.Second || config.KafkaConsumer.MaximumAttempts != 3 {
		t.Fatalf("LoadConfig() = %#v", config)
	}
}
