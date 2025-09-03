package dtos

import (
	"time"

	"github.com/google/uuid"

	enums "notezy-backend/app/models/schemas/enums"
)

/* ============================== Request DTO ============================== */

type GetMyInfoReqDto struct {
	NotezyRequest[
		any,
		struct {
			UserId uuid.UUID // extracted from the access token of AuthMiddleware()
		},
		any,
		any,
	]
}

type UpdateMyInfoReqDto struct {
	NotezyRequest[
		any,
		struct {
			UserId uuid.UUID // extracted from the access token of AuthMiddleware()
		},
		struct {
			PartialUpdateDto[struct {
				CoverBackgroundURL *string           `json:"coverBackgroundURL" validate:"omitnil,isimageurl"`
				AvatarURL          *string           `json:"avatarURL" validate:"omitnil,isimageurl"`
				Header             *string           `json:"header" validate:"omitnil,min=0,max=64"`
				Introduction       *string           `json:"introduction" validate:"omitnil,min=0,max=256"`
				Gender             *enums.UserGender `json:"gender" validate:"omitnil,isgender"`
				Country            *enums.Country    `json:"country" validate:"omitnil,iscountry"`
				BirthDate          *time.Time        `json:"birthDate" validate:"omitnil,notfuture"`
			}]
		},
		any,
	]
}

/* ============================== Response DTO ============================== */

type GetMyInfoResDto struct {
	CoverBackgroundURL *string          `json:"coverBackgroundURL"`
	AvatarURL          *string          `json:"avatarURL"`
	Header             *string          `json:"header"`
	Introduction       *string          `json:"introduction"`
	Gender             enums.UserGender `json:"gender"`
	Country            *enums.Country   `json:"country"`
	BirthDate          time.Time        `json:"birthDate"`
	UpdatedAt          time.Time        `json:"updatedAt"`
}

type UpdateMyInfoResDto struct {
	UpdatedAt time.Time `json:"updatedAt"`
}
