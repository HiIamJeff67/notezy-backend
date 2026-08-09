package topics

import (
	"time"

	emaileventscontract "github.com/HiIamJeff67/notezy-backend/contracts/email/v1/events"
)

func CoreEmailRequestTopicSpec() TopicSpec {
	return TopicSpec{
		Name:                emaileventscontract.CoreEmailRequestTopic.String(),
		Partitions:          3,
		ReplicationFactor:   1,
		Retention:           7 * 24 * time.Hour,
		CleanupPolicy:       "delete",
		MinInSyncReplicas:   1,
		CreateDeadLetter:    true,
		DeadLetterRetention: 30 * 24 * time.Hour,
	}
}
