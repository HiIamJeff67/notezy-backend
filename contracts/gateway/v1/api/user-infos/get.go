package userinfosdto

import (
	"time"

	apiv1 "github.com/HiIamJeff67/notezy-backend/contracts/gateway/v1/api"
)

type GetMyInfoRequestDto struct {
	apiv1.RequestDto[
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
