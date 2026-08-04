package authdto

import (
	"time"

	coreapicontract "github.com/HiIamJeff67/notezy-backend/contracts/core/v1/api"
)

type ResetEmailRequestDto struct {
	coreapicontract.RequestDto[
		struct {
			UserAgent string `json:"userAgent" validate:"required,isuseragent"`
		},
		struct {
			NewEmail string `json:"newEmail" validate:"required,email"`
			AuthCode string `json:"authCode" validate:"required,isnumberstring,len=6"`
		},
		struct{},
		struct{},
	]
}

type ResetEmailResponseDto struct {
	UpdatedAt time.Time `json:"updatedAt"`
}
