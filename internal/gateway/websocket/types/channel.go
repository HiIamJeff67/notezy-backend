package realtimetypes

import "github.com/google/uuid"

type Channel struct {
	Type                 ChannelType
	Id                   uuid.UUID
	Permission           ChannelPermission
	AcknowledgedSequence int64
}
