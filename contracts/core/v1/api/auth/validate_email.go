package apicontract

import (
	"time"

	coreapicontract "github.com/HiIamJeff67/notezy-backend/contracts/core/v1/api"
)

type ValidateEmailRequestDto struct {
	coreapicontract.RequestDto[
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

type ValidateEmailResponseDto struct {
	UpdatedAt time.Time `json:"updatedAt"`
}
