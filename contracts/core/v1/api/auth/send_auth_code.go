package apicontract

import (
	"time"

	coreapicontract "github.com/HiIamJeff67/notezy-backend/contracts/core/v1/api"
)

type SendAuthCodeRequestDto struct {
	coreapicontract.RequestDto[
		struct {
			UserAgent string `json:"userAgent" validate:"required,isuseragent"`
		},
		struct {
			Email string `json:"email" validate:"required,email"`
		},
		struct{},
		struct{},
	]
}

type SendAuthCodeResponseDto struct {
	AuthCodeExpiredAt  time.Time `json:"authCodeExpiredAt"`
	BlockAuthCodeUntil time.Time `json:"blockAuthCodeUntil"`
	UpdatedAt          time.Time `json:"updatedAt"`
}
