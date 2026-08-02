package userinfosdto

import (
	"time"

	apiv1 "github.com/HiIamJeff67/notezy-backend/contracts/gateway/v1/api"
)

type UpdateMyInfoRequestDto struct {
	apiv1.RequestDto[
		struct {
			UserAgent string `json:"userAgent" validate:"required,isuseragent"`
		},
		struct {
			Values struct {
				CoverBackgroundURL *string    `json:"coverBackgroundURL,omitempty" validate:"omitnil,isimageurl"`
				AvatarURL          *string    `json:"avatarURL,omitempty" validate:"omitnil,isimageurl"`
				Header             *string    `json:"header,omitempty" validate:"omitnil,min=0,max=64"`
				Introduction       *string    `json:"introduction,omitempty" validate:"omitnil,min=0,max=256"`
				Gender             *string    `json:"gender,omitempty" validate:"omitnil,isgender"`
				Country            *string    `json:"country,omitempty" validate:"omitnil,iscountry"`
				BirthDate          *time.Time `json:"birthDate,omitempty" validate:"omitnil,notfuture"`
			} `json:"values"`
			SetNull *map[string]bool `json:"setNull,omitempty"`
		},
		struct{},
		struct{},
	]
}

type UpdateMyInfoResponseDto struct {
	UpdatedAt time.Time `json:"updatedAt"`
}
