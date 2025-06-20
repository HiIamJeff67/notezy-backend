package dtos

import (
	"notezy-backend/app/models/enums"
	"time"

	"github.com/google/uuid"
)

/* ============================== Request DTO ============================== */

type GetMyInfoReqDto struct {
	UserId uuid.UUID // extracted from the access token of AuthMiddleware()
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
