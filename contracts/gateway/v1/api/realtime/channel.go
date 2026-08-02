package realtimedto

import (
	"time"

	"github.com/google/uuid"

	apiv1 "github.com/HiIamJeff67/notezy-backend/contracts/gateway/v1/api"
)

type CreateMyBlockPackChannelTicketRequestDto struct {
	apiv1.RequestDto[
		struct {
			UserAgent string `json:"userAgent" validate:"required,isuseragent"`
		},
		struct {
			BlockPackId uuid.UUID `json:"blockPackId" validate:"required"`
			Permission  string    `json:"permission" validate:"required,oneof=read write"`
		},
		struct{},
		struct{},
	]
}

type CreateMyBlockPackChannelTicketResponseDto struct {
	ChannelTicket           string    `json:"channelTicket"`
	ExpiresAt               time.Time `json:"expiresAt"`
	ChannelType             string    `json:"channelType"`
	ChannelId               uuid.UUID `json:"channelId"`
	Permission              string    `json:"permission"`
	RoomName                string    `json:"roomName"`
	FragmentName            string    `json:"fragmentName"`
	SchemaId                string    `json:"schemaId"`
	SchemaVersion           int       `json:"schemaVersion"`
	RealtimeProtocolVersion int       `json:"realtimeProtocolVersion"`
	LastUpdateSequence      int64     `json:"lastUpdateSequence"`
	CompactedUntilSequence  int64     `json:"compactedUntilSequence"`
}
