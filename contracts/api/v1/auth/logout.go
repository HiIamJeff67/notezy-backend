package authdto

import (
	"time"

	apiv1 "github.com/HiIamJeff67/notezy-backend/contracts/api/v1"
)

type LogoutRequestDto struct {
	apiv1.RequestDto[
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
