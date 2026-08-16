package topics

import (
	"time"

	coreeventscontract "github.com/HiIamJeff67/notegic-backend/contracts/core/v1/events"
)

func CoreDurableJobYjsMaintenanceHintTopicSpec() TopicSpec {
	return TopicSpec{
		Name:                coreeventscontract.CoreDurableJobYjsMaintenanceHintTopic.String(),
		Partitions:          3,
		ReplicationFactor:   1,
		Retention:           7 * 24 * time.Hour,
		CleanupPolicy:       "delete",
		MinInSyncReplicas:   1,
		CreateDeadLetter:    true,
		DeadLetterRetention: 30 * 24 * time.Hour,
	}
}
