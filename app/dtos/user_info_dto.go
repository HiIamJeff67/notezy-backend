package dtos

import (
	"time"

	"github.com/google/uuid"

	enums "notezy-backend/app/models/enums"
)

/* ============================== Request DTO ============================== */

type GetMyInfoReqDto struct {
	UserId uuid.UUID // extracted from the access token of AuthMiddleware()
}

type UpdateMyInfoReqDto struct {
	UserId uuid.UUID // extracted from the access token of AuthMiddleware()
	PartialUpdateDto[struct {
		CoverBackgroundURL *string           `json:"coverBackgroundURL" validate:"omitempty"`
		AvatarURL          *string           `json:"avatarURL" validate:"omitempty"`
		Header             *string           `json:"header" validate:"omitempty"`
		Introduction       *string           `json:"introduction" validate:"omitempty"`
		Gender             *enums.UserGender `json:"gender" validate:"omitempty"`
		Country            *enums.Country    `json:"country" validate:"omitempty"`
		BirthDate          *time.Time        `json:"birthDate" validate:"omitempty"`
	}]
}

/* ============================== Response DTO ============================== */

type GetMyInfoResDto struct {
	CoverBackgroundURL *string          `json:"coverBackgroundURL"`
	AvatarURL          *string          `json:"avatarURL"`
	Header             *string          `json:"header"`
	Introduction       *string          `json:"introduction"`
	Gender             enums.UserGender `json:"gender"`
	Country            enums.Country    `json:"country"`
	BirthDate          time.Time        `json:"birthDate"`
	UpdatedAt          time.Time        `json:"updatedAt"`
}

type UpdateMyInfoResDto struct {
	UpdatedAt time.Time `json:"updatedAt"`
}
