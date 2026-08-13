package durablejobeventscontract

import eventcontract "github.com/HiIamJeff67/notezy-backend/contracts/types/events"

const (
	CoreDurableJobRoutineTaskTopic                     eventcontract.Topic = "notezy.core.durablejob.routine-task.v1"
	DurableJobRealtimeGatewayRoutineTaskLifecycleTopic eventcontract.Topic = "notezy.durablejob.realtime-gateway.routine-task-lifecycle.v1"
)

const (
	DurableJobCoreYjsMaintenanceRequestTopic eventcontract.Topic = "notezy.durablejob.core.yjs-maintenance-request.v1"
	DurableJobCoreYjsMaintenanceResultTopic  eventcontract.Topic = "notezy.core.durablejob.yjs-maintenance-result.v1"
)

const (
	AggregateType_BlockPack        eventcontract.AggregateType = "BlockPack"
	AggregateType_DurableJobWorker eventcontract.AggregateType = "DurableJobWorker"
	AggregateType_RoutineTask      eventcontract.AggregateType = "RoutineTask"
)

const (
	EventType_RoutineTaskClaimRequested eventcontract.EventType = "RoutineTaskClaimRequested"
	EventType_RoutineTasksAssigned      eventcontract.EventType = "RoutineTasksAssigned"
	EventType_RoutineTasksCompleted     eventcontract.EventType = "RoutineTasksCompleted"
	EventType_RoutineTasksFailed        eventcontract.EventType = "RoutineTasksFailed"
	EventType_RoutineTaskRunning        eventcontract.EventType = "RoutineTaskRunning"
	EventType_YjsMaintenanceRequested   eventcontract.EventType = "YjsMaintenanceRequested"
	EventType_YjsMaintenanceCompleted   eventcontract.EventType = "YjsMaintenanceCompleted"
)
