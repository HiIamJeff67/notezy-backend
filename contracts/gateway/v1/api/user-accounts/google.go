package useraccountsdto

import (
	"time"

	apiv1 "github.com/HiIamJeff67/notezy-backend/contracts/gateway/v1/api"
)

type BindGoogleAccountRequestDto struct {
	apiv1.RequestDto[
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
	apiv1.RequestDto[
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
