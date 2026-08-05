package durablejobeventscontract

import eventcontract "github.com/HiIamJeff67/notezy-backend/contracts/events"

const (
	CoreDurableJobRoutineTaskTopic eventcontract.Topic = "notezy.core.durablejob.routine-task.v1"

	DurableJobRoutineTaskConsumerGroup           = "notezy-durablejob-routine-task-v1"
	CoreDurableJobRoutineTaskClaimConsumerGroup  = "notezy-core-durablejob-routine-task-v1"
	CoreDurableJobRoutineTaskResultConsumerGroup = "notezy-core-durablejob-routine-task-result-v1"
)

const (
	CoreDurableJobYjsMaintenanceHintTopic    eventcontract.Topic = "notezy.core.durablejob.yjs-maintenance-hint.v1"
	DurableJobCoreYjsMaintenanceRequestTopic eventcontract.Topic = "notezy.durablejob.core.yjs-maintenance-request.v1"
	CoreYjsWorkerMaintenanceCommandTopic     eventcontract.Topic = "notezy.core.yjsworker.maintenance-command.v1"
	YjsWorkerCoreMaintenanceResultTopic      eventcontract.Topic = "notezy.yjsworker.core.maintenance-result.v1"
	DurableJobCoreYjsMaintenanceResultTopic  eventcontract.Topic = "notezy.core.durablejob.yjs-maintenance-result.v1"
)

const (
	CoreDurableJobYjsMaintenanceConsumerGroup   = "notezy-core-durablejob-yjs-maintenance-v1"
	DurableJobYjsMaintenanceConsumerGroup       = "notezy-durablejob-yjs-maintenance-v1"
	YjsWorkerMaintenanceConsumerGroup           = "notezy-yjs-worker-maintenance-v1"
	CoreYjsWorkerMaintenanceResultConsumerGroup = "notezy-core-yjs-maintenance-result-v1"
)

const (
	AggregateType_BlockPack        eventcontract.AggregateType = "BlockPack"
	AggregateType_DurableJobWorker eventcontract.AggregateType = "DurableJobWorker"
)

const (
	EventType_RoutineTaskClaimRequested eventcontract.EventType = "RoutineTaskClaimRequested"
	EventType_RoutineTasksAssigned      eventcontract.EventType = "RoutineTasksAssigned"
	EventType_RoutineTasksCompleted     eventcontract.EventType = "RoutineTasksCompleted"
	EventType_RoutineTasksFailed        eventcontract.EventType = "RoutineTasksFailed"
	EventType_YjsMaintenanceHint        eventcontract.EventType = "YjsMaintenanceHint"
	EventType_YjsMaintenanceRequested   eventcontract.EventType = "YjsMaintenanceRequested"
	EventType_YjsMaintenanceCommand     eventcontract.EventType = "YjsMaintenanceCommand"
	EventType_YjsMaintenanceCompleted   eventcontract.EventType = "YjsMaintenanceCompleted"
)
