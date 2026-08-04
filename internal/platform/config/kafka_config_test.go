package configs

import "testing"

func TestKafkaDefaults(t *testing.T) {
	t.Setenv("KAFKA_BROKERS", "")
	t.Setenv("KAFKA_CLIENT_ID", "")
	t.Setenv("KAFKA_CONSUMER_GROUP", "")
	t.Setenv("KAFKA_DIAL_TIMEOUT_SECONDS", "")
	t.Setenv("KAFKA_TLS_ENABLED", "")

	kafkaConfig := Kafka()
	if len(kafkaConfig.Brokers) != 1 || kafkaConfig.Brokers[0] != "127.0.0.1:9094" {
		t.Fatalf("unexpected Kafka brokers: %#v", kafkaConfig.Brokers)
	}
	if kafkaConfig.ClientId != "notezy-runtime" || kafkaConfig.ConsumerGroup != "notezy-runtime" {
		t.Fatalf("unexpected Kafka identifiers: %#v", kafkaConfig)
	}
	if kafkaConfig.TLS.Enabled {
		t.Fatal("expected Kafka TLS to be disabled by default")
	}
}

func TestKafkaTrimsConfiguredBrokers(t *testing.T) {
	t.Setenv("KAFKA_BROKERS", " kafka-1:9092, ,kafka-2:9092 ")

	kafkaConfig := Kafka()
	if len(kafkaConfig.Brokers) != 2 ||
		kafkaConfig.Brokers[0] != "kafka-1:9092" ||
		kafkaConfig.Brokers[1] != "kafka-2:9092" {
		t.Fatalf("unexpected Kafka brokers: %#v", kafkaConfig.Brokers)
	}
}
