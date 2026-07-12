package dtos

import (
	"time"

	"github.com/google/uuid"

	realtimetypes "github.com/HiIamJeff67/notezy-backend/app/realtime/types"
)

/* ============================== Request DTO ============================== */

type CreateMyRealtimeConnectionTicketReqDto struct {
	NotezyRequest[
		struct {
			UserAgent string `json:"userAgent" validate:"required,isuseragent"`
		},
		struct {
			UserPublicId uuid.UUID // extracted from the AuthMiddleware()
		},
		any,
		any,
	]
}

type CreateMyBlockPackChannelTicketReqDto struct {
	NotezyRequest[
		struct {
			UserAgent string `json:"userAgent" validate:"required,isuseragent"`
		},
		struct {
			UserId       uuid.UUID // extracted from the access token of AuthMiddleware()
			UserPublicId uuid.UUID // extracted from the AuthMiddleware()
		},
		struct {
			BlockPackId uuid.UUID                       `json:"blockPackId" validate:"required"`
			Permission  realtimetypes.ChannelPermission `json:"permission" validate:"required,oneof=read write"`
		},
		any,
	]
}

/* ============================== Response DTO ============================== */

type CreateMyRealtimeConnectionTicketResDto struct {
	RealtimeEndpoint        string    `json:"realtimeEndpoint"`
	RealtimeProtocolVersion int       `json:"realtimeProtocolVersion"`
	ConnectionTicket        string    `json:"connectionTicket"`
	ExpiresAt               time.Time `json:"expiresAt"`
}

type CreateMyBlockPackChannelTicketResDto struct {
	ChannelTicket           string                          `json:"channelTicket"`
	ExpiresAt               time.Time                       `json:"expiresAt"`
	ChannelType             realtimetypes.ChannelType       `json:"channelType"`
	ChannelId               uuid.UUID                       `json:"channelId"`
	Permission              realtimetypes.ChannelPermission `json:"permission"`
	RoomName                string                          `json:"roomName"`
	FragmentName            string                          `json:"fragmentName"`
	SchemaId                string                          `json:"schemaId"`
	SchemaVersion           int                             `json:"schemaVersion"`
	RealtimeProtocolVersion int                             `json:"realtimeProtocolVersion"`
	LastUpdateSequence      int64                           `json:"lastUpdateSequence"`
	CompactedUntilSequence  int64                           `json:"compactedUntilSequence"`
}
