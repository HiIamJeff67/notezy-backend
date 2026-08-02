package realtimedto

import (
	"time"

	apiv1 "github.com/HiIamJeff67/notezy-backend/contracts/gateway/v1/api"
)

type CreateMyRealtimeConnectionTicketRequestDto struct {
	apiv1.RequestDto[
		struct {
			UserAgent string `json:"userAgent" validate:"required,isuseragent"`
		},
		struct{},
		struct{},
		struct{},
	]
}

type CreateMyRealtimeConnectionTicketResponseDto struct {
	RealtimeEndpoint        string    `json:"realtimeEndpoint"`
	RealtimeProtocolVersion int       `json:"realtimeProtocolVersion"`
	ConnectionTicket        string    `json:"connectionTicket"`
	ExpiresAt               time.Time `json:"expiresAt"`
}
