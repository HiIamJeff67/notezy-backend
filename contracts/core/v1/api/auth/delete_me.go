package authdto

import (
	"time"

	coreapicontract "github.com/HiIamJeff67/notezy-backend/contracts/core/v1/api"
)

type DeleteMeRequestDto struct {
	coreapicontract.RequestDto[
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
