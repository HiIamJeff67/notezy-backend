package yjsworkereventscontract

import eventcontract "github.com/HiIamJeff67/notezy-backend/contracts/events"

const (
	YjsWorkerCoreCommandTopic  eventcontract.Topic = "notezy.yjsworker.core.command.v1"
	CoreYjsWorkerReplyTopic    eventcontract.Topic = "notezy.core.yjsworker.reply.v1"
	CoreYjsWorkerConsumerGroup                     = "notezy-core-yjsworker-v1"
)

const (
	EventType_YjsWorkerCommand          eventcontract.EventType = "YjsWorkerCommand"
	EventType_YjsWorkerCommandCompleted eventcontract.EventType = "YjsWorkerCommandCompleted"
)

const AggregateType_BlockPack eventcontract.AggregateType = "BlockPack"
