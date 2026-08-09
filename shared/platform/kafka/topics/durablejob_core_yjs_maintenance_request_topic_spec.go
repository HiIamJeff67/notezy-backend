package topics

import (
	"time"

	durablejobeventscontract "github.com/HiIamJeff67/notezy-backend/contracts/durable-job/v1/events"
)

func DurableJobCoreYjsMaintenanceRequestTopicSpec() TopicSpec {
	return TopicSpec{
		Name:                durablejobeventscontract.DurableJobCoreYjsMaintenanceRequestTopic.String(),
		Partitions:          3,
		ReplicationFactor:   1,
		Retention:           7 * 24 * time.Hour,
		CleanupPolicy:       "delete",
		MinInSyncReplicas:   1,
		CreateDeadLetter:    true,
		DeadLetterRetention: 30 * 24 * time.Hour,
	}
}
