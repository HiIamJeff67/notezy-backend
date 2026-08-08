package emaileventscontract

import eventcontract "github.com/HiIamJeff67/notezy-backend/contracts/types/events"

const (
	CoreEmailRequestTopic  eventcontract.Topic = "notezy.core.email.request.v1"
	CoreEmailConsumerGroup                     = "notezy-email-core-v1"
)

const (
	AggregateType_EmailRequest eventcontract.AggregateType = "EmailRequest"
	EventType_EmailRequested   eventcontract.EventType     = "EmailRequested"
)
