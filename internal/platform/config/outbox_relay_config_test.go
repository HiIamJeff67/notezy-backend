package configs

import "testing"

func TestOutboxRelayDefaults(t *testing.T) {
	t.Setenv("OUTBOX_RELAY_BATCH_SIZE", "")
	t.Setenv("OUTBOX_RELAY_POLL_INTERVAL_MILLISECONDS", "")
	t.Setenv("OUTBOX_RELAY_CLAIM_TIMEOUT_SECONDS", "")
	t.Setenv("OUTBOX_RELAY_INITIAL_BACKOFF_MILLISECONDS", "")
	t.Setenv("OUTBOX_RELAY_MAXIMUM_BACKOFF_SECONDS", "")
	t.Setenv("OUTBOX_RELAY_RETENTION_HOURS", "")
	t.Setenv("OUTBOX_RELAY_CLEANUP_INTERVAL_MINUTES", "")

	config := OutboxRelay()
	if config.BatchSize != 100 || config.PollInterval.Milliseconds() != 1000 ||
		config.ClaimTimeout.Seconds() != 30 || config.InitialBackoff.Milliseconds() != 1000 ||
		config.MaximumBackoff.Seconds() != 60 || config.Retention.Hours() != 168 ||
		config.CleanupInterval.Minutes() != 60 {
		t.Fatalf("unexpected outbox relay config: %#v", config)
	}
}
