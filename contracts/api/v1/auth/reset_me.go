package authdto

import (
	"time"

	apiv1 "github.com/HiIamJeff67/notezy-backend/contracts/api/v1"
)

type ResetMeRequestDto struct {
	apiv1.RequestDto[
		struct {
			UserAgent string `json:"userAgent" validate:"required,isuseragent"`
		},
		struct {
			AuthCode string `json:"authCode" validate:"required,isnumberstring,len=6"`
		},
		struct{},
		struct{},
	]
}

type ResetMeResponseDto struct {
	UpdatedAt time.Time `json:"updatedAt"`
}
