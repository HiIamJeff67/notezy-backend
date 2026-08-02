package authdto

import (
	"time"

	apiv1 "github.com/HiIamJeff67/notezy-backend/contracts/gateway/v1/api"
)

type DeleteMeRequestDto struct {
	apiv1.RequestDto[
		struct {
			UserAgent string `json:"userAgent" validate:"required,isuseragent"`
		},
		struct {
			AuthCode string `json:"authCode" validate:"omitempty,isnumberstring,len=6"`
		},
		struct{},
		struct{},
	]
}

type DeleteMeResponseDto struct {
	DeletedAt time.Time `json:"deletedAt"`
}
