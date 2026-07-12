package realtimetypes

import "github.com/google/uuid"

type Channel struct {
	Type                 ChannelType
	Id                   uuid.UUID
	AcknowledgedSequence int64
}
