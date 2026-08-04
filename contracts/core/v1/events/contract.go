package eventscontract

import (
	"time"

	"github.com/google/uuid"

	realtimegatewaycontract "github.com/HiIamJeff67/notezy-backend/contracts/realtimegateway/v1"
)

const Version = "v1"

type Topic string

func (t Topic) String() string {
	return string(t)
}

const (
	CoreLifecycleTopic        Topic = "notezy.core.lifecycle.v1"
	YjsWorkerCoreCommandTopic Topic = "notezy.yjsworker.core.command.v1"
	CoreYjsWorkerReplyTopic   Topic = "notezy.core.yjsworker.reply.v1"
)

type AggregateType string

const (
	AggregateType_RootShelf AggregateType = "RootShelf"
	AggregateType_SubShelf  AggregateType = "SubShelf"
	AggregateType_BlockPack AggregateType = "BlockPack"
	AggregateType_User      AggregateType = "User"
)

type EventType string

const (
	EventType_BlockPackAccessRevoked     EventType = "BlockPackAccessRevoked"
	EventType_BlockPackRoomPolicyChanged EventType = "BlockPackRoomPolicyChanged"
	EventType_RootShelfPermissionRevoked EventType = "RootShelfPermissionRevoked"
	EventType_UserSessionsRevoked        EventType = "UserSessionsRevoked"
	EventType_YjsWorkerCommand           EventType = "YjsWorkerCommand"
	EventType_YjsWorkerCommandCompleted  EventType = "YjsWorkerCommandCompleted"
)

type RoomAdmissionEnforcementStrategy = realtimegatewaycontract.RoomAdmissionEnforcementStrategy

const (
	RoomAdmissionEnforcementStrategy_RejectNewSubscriber = realtimegatewaycontract.RoomAdmissionEnforcementStrategy_RejectNewSubscriber
)

type EventEnvelope[D any] struct {
	SchemaVersion string        `json:"schemaVersion"`
	EventId       uuid.UUID     `json:"eventId"`
	EventType     EventType     `json:"eventType"`
	AggregateType AggregateType `json:"aggregateType"`
	AggregateId   uuid.UUID     `json:"aggregateId"`
	KafkaKey      string        `json:"kafkaKey"`
	OccurredAt    time.Time     `json:"occurredAt"`
	CorrelationId string        `json:"correlationId"`
	CausationId   *uuid.UUID    `json:"causationId,omitempty"`
	Trace         TraceMetadata `json:"trace"`
	Data          D             `json:"data"`
}

type TraceMetadata struct {
	TraceParent string `json:"traceParent,omitempty"`
	TraceState  string `json:"traceState,omitempty"`
}
