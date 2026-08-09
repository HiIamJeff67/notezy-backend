package kafka

import (
	"testing"
	"time"

	kafkatopics "github.com/HiIamJeff67/notezy-backend/shared/platform/kafka/topics"
)

func TestValidateTopicSpecRequiresExplicitCreationSettings(t *testing.T) {
	if err := validateTopicSpec(kafkatopics.TopicSpec{Name: "notezy.test.v1"}); err == nil {
		t.Fatal("validate topic spec returned nil for incomplete settings")
	}

	specification := kafkatopics.TopicSpec{
		Name:                "notezy.test.v1",
		Partitions:          3,
		ReplicationFactor:   1,
		Retention:           7 * 24 * time.Hour,
		CleanupPolicy:       "delete",
		MinInSyncReplicas:   1,
		CreateDeadLetter:    true,
		DeadLetterRetention: 30 * 24 * time.Hour,
	}
	if err := validateTopicSpec(specification); err != nil {
		t.Fatalf("validate complete topic spec: %v", err)
	}
}

func TestTopicRequestUsesKafkaMillisecondRetention(t *testing.T) {
	request := topicRequest(kafkatopics.TopicSpec{
		Name:              "notezy.test.v1",
		Partitions:        3,
		ReplicationFactor: 1,
		Retention:         5 * time.Minute,
		CleanupPolicy:     "delete",
		MinInSyncReplicas: 1,
	})

	if got := *request.Configs[1].Value; got != "300000" {
		t.Fatalf("retention.ms = %q, want 300000", got)
	}
}

func TestValidateTopicSpecRejectsInvalidValues(t *testing.T) {
	tests := []kafkatopics.TopicSpec{
		{Name: ""},
		{Name: " notezy.test.v1 ", Partitions: 3, ReplicationFactor: 1, Retention: time.Hour, CleanupPolicy: "delete", MinInSyncReplicas: 1},
		{Name: "notezy.test.v1", Partitions: -1},
		{Name: "notezy.test.v1", ReplicationFactor: -1},
		{Name: "notezy.test.v1", Retention: -time.Second},
		{Name: "notezy.test.v1", MinInSyncReplicas: -1},
		{Name: "notezy.test.v1", Partitions: 1, ReplicationFactor: 1, Retention: time.Hour, CleanupPolicy: "delete", MinInSyncReplicas: 1, CreateDeadLetter: true},
	}

	for _, specification := range tests {
		if err := validateTopicSpec(specification); err == nil {
			t.Fatalf("validate topic spec %+v returned nil error", specification)
		}
	}
}
