package apicontract

import (
	"time"

	coreapicontract "github.com/HiIamJeff67/notegic-backend/contracts/core/v1/api"
)

type LogoutRequestDto struct {
	coreapicontract.RequestDto[
		struct {
			UserAgent string `json:"userAgent" validate:"required,isuseragent"`
		},
		struct{},
		struct{},
		struct{},
	]
}

type LogoutResponseDto struct {
	UpdatedAt time.Time `json:"updatedAt"`
}
