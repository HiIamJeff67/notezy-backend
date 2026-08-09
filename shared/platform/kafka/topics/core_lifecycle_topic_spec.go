package topics

import (
	"time"

	coreeventscontract "github.com/HiIamJeff67/notezy-backend/contracts/core/v1/events"
)

func CoreLifecycleTopicSpec() TopicSpec {
	return TopicSpec{
		Name:                coreeventscontract.CoreLifecycleTopic.String(),
		Partitions:          3,
		ReplicationFactor:   1,
		Retention:           7 * 24 * time.Hour,
		CleanupPolicy:       "delete",
		MinInSyncReplicas:   1,
		CreateDeadLetter:    true,
		DeadLetterRetention: 30 * 24 * time.Hour,
	}
}
