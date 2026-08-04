package configs

import "testing"

func TestKafkaConsumerDefaults(t *testing.T) {
	t.Setenv("KAFKA_CONSUMER_MAXIMUM_ATTEMPTS", "")
	t.Setenv("KAFKA_CONSUMER_INITIAL_RETRY_BACKOFF_MILLISECONDS", "")
	t.Setenv("KAFKA_CONSUMER_MAXIMUM_RETRY_BACKOFF_MILLISECONDS", "")
	t.Setenv("KAFKA_CONSUMER_MAXIMUM_POLL_RECORDS", "")

	consumerConfig := KafkaConsumer()
	if consumerConfig.MaximumAttempts != 3 ||
		consumerConfig.InitialRetryBackoff.Milliseconds() != 250 ||
		consumerConfig.MaximumRetryBackoff.Milliseconds() != 5000 ||
		consumerConfig.MaximumPollRecords != 100 {
		t.Fatalf("unexpected Kafka consumer defaults: %#v", consumerConfig)
	}
}

func TestKafkaConsumerKeepsRetryBackoffBounded(t *testing.T) {
	t.Setenv("KAFKA_CONSUMER_INITIAL_RETRY_BACKOFF_MILLISECONDS", "1000")
	t.Setenv("KAFKA_CONSUMER_MAXIMUM_RETRY_BACKOFF_MILLISECONDS", "500")

	consumerConfig := KafkaConsumer()
	if consumerConfig.MaximumRetryBackoff < consumerConfig.InitialRetryBackoff {
		t.Fatalf("Kafka consumer retry backoff is not bounded: %#v", consumerConfig)
	}
}
