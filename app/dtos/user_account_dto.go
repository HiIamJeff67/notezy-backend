package dtos

import (
	"time"

	"github.com/google/uuid"

	enums "notezy-backend/app/models/schemas/enums"
)

/* ============================== Request DTO ============================== */

type GetMyAccountReqDto struct {
	UserId uuid.UUID // extracted from the access token of AuthMiddleware()
}

type UpdateMyAccountReqDto struct {
	UserId uuid.UUID // extracted from the access token of AuthMiddleware()
	PartialUpdateDto[struct {
		CountryCode       *enums.CountryCode `json:"countryCode" validate:"omitempty"`
		PhoneNumber       *string            `json:"phoneNumber" validate:"omitempty"`
		GoogleCredential  *string            `json:"googleCrendential" validate:"omitempty"`
		DiscordCredential *string            `json:"discordCrendential" validate:"omitempty"`
	}]
}

/* ============================== Response DTO ============================== */

type GetMyAccountResDto struct {
	CountryCode       enums.CountryCode `json:"countryCode"`
	PhoneNumber       *string           `json:"phoneNumber"`
	GoogleCredential  *string           `json:"googleCrendential"`
	DiscordCredential *string           `json:"discordCrendential"`
}

type UpdateMyAccountResDto struct {
	UpdatedAt time.Time `json:"updatedAt"`
}
