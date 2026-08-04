package useraccountsdto

import (
	"time"

	coreapicontract "github.com/HiIamJeff67/notezy-backend/contracts/core/v1/api"
)

type BindGoogleAccountRequestDto struct {
	coreapicontract.RequestDto[
		struct {
			UserAgent string `json:"userAgent" validate:"required,isuseragent"`
		},
		struct {
			AuthorizationCode string `json:"authorizationCode" validate:"required"`
		},
		struct{},
		struct{},
	]
}

type BindGoogleAccountResponseDto struct {
	UpdatedAt time.Time `json:"updatedAt"`
}

type UnbindGoogleAccountRequestDto struct {
	coreapicontract.RequestDto[
		struct {
			UserAgent string `json:"userAgent" validate:"required,isuseragent"`
		},
		struct {
			AuthCode string `json:"authCode" validate:"required"`
		},
		struct{},
		struct{},
	]
}

type UnbindGoogleAccountResponseDto struct {
	UpdatedAt time.Time `json:"updatedAt"`
}
