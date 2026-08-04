package realtimedto

import (
	"github.com/google/uuid"

	coreapicontract "github.com/HiIamJeff67/notezy-backend/contracts/core/v1/api"
)

type GetMyBlockPackRealtimeParticipantsRequestDto struct {
	coreapicontract.RequestDto[
		struct {
			UserAgent string `json:"userAgent" validate:"required,isuseragent"`
		},
		struct {
			Participants []RealtimeBlockPackParticipantRequestDto `json:"participants" validate:"max=100,dive"`
		},
		struct {
			BlockPackId uuid.UUID `json:"blockPackId" validate:"required"`
		},
		struct{},
	]
}

type RealtimeBlockPackParticipantRequestDto struct {
	UserPublicId      uuid.UUID `json:"userPublicId" validate:"required"`
	ChannelPermission string    `json:"channelPermission" validate:"required,oneof=read write"`
	ConnectionCount   int       `json:"connectionCount" validate:"min=1,max=8"`
}

type RealtimeBlockPackParticipantResponseDto struct {
	UserPublicId      uuid.UUID `json:"userPublicId"`
	Name              string    `json:"name"`
	DisplayName       string    `json:"displayName"`
	ChannelPermission string    `json:"channelPermission"`
	ConnectionCount   int       `json:"connectionCount"`
}

type GetMyBlockPackRealtimeParticipantsResponseDto []RealtimeBlockPackParticipantResponseDto
