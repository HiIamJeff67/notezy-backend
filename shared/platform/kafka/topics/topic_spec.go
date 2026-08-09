package topics

import "time"

// TopicSpec describes the complete broker-level settings required by one
// application topic. Topic owners must provide every creation setting
// explicitly; the shared Kafka layer does not infer broker defaults.
type TopicSpec struct {
	Name                string
	Partitions          int32
	ReplicationFactor   int16
	Retention           time.Duration
	CleanupPolicy       string
	MinInSyncReplicas   int
	CreateDeadLetter    bool
	DeadLetterRetention time.Duration
}
