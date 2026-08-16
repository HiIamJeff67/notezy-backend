package topics

// All returns the complete topic catalog owned by the Notegic runtimes.
// Every entry is constructed by its topic owner with explicit broker settings.
func All() []TopicSpec {
	return []TopicSpec{
		CoreLifecycleTopicSpec(),
		CoreDurableJobYjsMaintenanceHintTopicSpec(),
		CoreNotificationTopicSpec(),
		CoreDurableJobRoutineTaskTopicSpec(),
		DurableJobRealtimeGatewayRoutineTaskLifecycleTopicSpec(),
		DurableJobCoreYjsMaintenanceRequestTopicSpec(),
		DurableJobCoreYjsMaintenanceResultTopicSpec(),
		CoreEmailRequestTopicSpec(),
		NotificationTopicSpec(),
		YjsWorkerCoreCommandTopicSpec(),
		CoreYjsWorkerReplyTopicSpec(),
		YjsWorkerCoreMaintenanceCommandTopicSpec(),
		CoreYjsWorkerMaintenanceResultTopicSpec(),
	}
}
