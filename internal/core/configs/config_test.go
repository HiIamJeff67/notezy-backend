package config

import (
	"testing"
)

func TestLoadConfig(t *testing.T) {
	t.Setenv("CORE_LISTEN_ADDRESS", "127.0.0.1:7778")
	t.Setenv("OAUTH_GOOGLE_CLIENT_ID", "client-id")
	t.Setenv("OAUTH_GOOGLE_CLIENT_SECRET", "client-secret")
	t.Setenv("OAUTH_GOOGLE_REDIRECT_URL", "http://gateway/auth/google/callback")
	t.Setenv("OUTBOX_RELAY_BATCH_SIZE", "100")
	t.Setenv("OUTBOX_RELAY_POLL_INTERVAL", "1s")
	t.Setenv("OUTBOX_RELAY_CLAIM_TIMEOUT", "30s")
	t.Setenv("OUTBOX_RELAY_INITIAL_BACKOFF", "1s")
	t.Setenv("OUTBOX_RELAY_MAXIMUM_BACKOFF", "1m")
	t.Setenv("OUTBOX_RELAY_RETENTION", "168h")
	t.Setenv("OUTBOX_RELAY_CLEANUP_INTERVAL", "1h")
	t.Setenv("KAFKA_CONSUMER_MAXIMUM_ATTEMPTS", "3")
	t.Setenv("KAFKA_CONSUMER_INITIAL_RETRY_BACKOFF", "250ms")
	t.Setenv("KAFKA_CONSUMER_MAXIMUM_RETRY_BACKOFF", "5s")
	t.Setenv("KAFKA_CONSUMER_MAXIMUM_POLL_RECORDS", "100")
	t.Setenv("STORAGE_KEY_SALT", "salt")
	t.Setenv("CORE_USER_DATA_CACHE_SERVER_START", "0")
	t.Setenv("CORE_USER_DATA_CACHE_SERVER_SIZE", "4")
	t.Setenv("CORE_USER_DATA_CACHE_EXPIRES_IN", "1h")

	config, err := LoadConfig()
	if err != nil {
		t.Fatalf("LoadConfig() error = %v", err)
	}
	if config.OutboxRelay.BatchSize != 100 {
		t.Fatalf("LoadConfig() = %#v", config)
	}
}
