package apicontract

import (
	"time"

	coreapicontract "github.com/HiIamJeff67/notezy-backend/contracts/core/v1/api"
)

type GetMyInfoRequestDto struct {
	coreapicontract.RequestDto[
		struct {
			UserAgent string `json:"userAgent" validate:"required,isuseragent"`
		},
		struct{},
		struct{},
		struct{},
	]
}

type GetMyInfoResponseDto struct {
	CoverBackgroundURL *string   `json:"coverBackgroundURL"`
	AvatarURL          *string   `json:"avatarURL"`
	Header             *string   `json:"header"`
	Introduction       *string   `json:"introduction"`
	Gender             string    `json:"gender"`
	Country            *string   `json:"country"`
	BirthDate          time.Time `json:"birthDate"`
	UpdatedAt          time.Time `json:"updatedAt"`
}
