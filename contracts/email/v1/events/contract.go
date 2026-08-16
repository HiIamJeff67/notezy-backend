package emaileventscontract

import eventcontract "github.com/HiIamJeff67/notegic-backend/contracts/types/events"

const (
	CoreEmailRequestTopic  eventcontract.Topic = "notegic.core.email.request.v1"
	CoreEmailConsumerGroup                     = "notegic-email-core-v1"
)

const (
	AggregateType_EmailRequest eventcontract.AggregateType = "EmailRequest"
	EventType_EmailRequested   eventcontract.EventType     = "EmailRequested"
)
