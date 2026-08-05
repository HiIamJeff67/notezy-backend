package kafka

import (
	"testing"
)

func TestNewConsumerRejectsMissingGroupOrTopics(t *testing.T) {
	_, err := NewConsumer(ConsumerConfig{}, "notezy.core.lifecycle.v1")
	if err == nil {
		t.Fatal("expected Kafka consumer group to be required")
	}

	_, err = NewConsumer(ConsumerConfig{
		ConsumerGroup: "test-group",
	})
	if err == nil {
		t.Fatal("expected Kafka consumer topic to be required")
	}
}

func TestDeadLetterTopic(t *testing.T) {
	if DeadLetterTopic("notezy.core.lifecycle.v1") != "notezy.core.lifecycle.v1.dlq" {
		t.Fatal("unexpected Kafka dead-letter topic name")
	}
}
