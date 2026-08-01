package realtimedto

import (
	"github.com/google/uuid"

	apiv1 "github.com/HiIamJeff67/notezy-backend/contracts/api/v1"
)

type GetMyBlockPackRealtimeParticipantsRequestDto struct {
	apiv1.RequestDto[
		struct {
			UserAgent string `json:"userAgent" validate:"required,isuseragent"`
		},
		struct{},
		struct {
			BlockPackId uuid.UUID `json:"blockPackId" validate:"required"`
		},
		struct{},
	]
}

type RealtimeBlockPackParticipantResponseDto struct {
	UserPublicId      uuid.UUID `json:"userPublicId"`
	Name              string    `json:"name"`
	DisplayName       string    `json:"displayName"`
	ChannelPermission string    `json:"channelPermission"`
	ConnectionCount   int       `json:"connectionCount"`
}

type GetMyBlockPackRealtimeParticipantsResponseDto []RealtimeBlockPackParticipantResponseDto
