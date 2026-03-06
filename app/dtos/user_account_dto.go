package dtos

import (
	"time"

	"github.com/google/uuid"

	enums "notezy-backend/app/models/schemas/enums"
)

/* ============================== Request DTO ============================== */

type GetMyAccountReqDto struct {
	NotezyRequest[
		struct {
			UserAgent string `json:"userAgent" validate:"required,isuseragent"`
		},
		struct {
			UserId uuid.UUID // extracted from the access token of AuthMiddleware()
		},
		any,
		any,
	]
}

type UpdateMyAccountReqDto struct {
	NotezyRequest[
		struct {
			UserAgent string `json:"userAgent" validate:"required,isuseragent"`
		},
		struct {
			UserId uuid.UUID // extracted from the access token of AuthMiddleware()
		},
		struct {
			AuthCode string `json:"authCode" validate:"required,isnumberstring,len=6"`
			PartialUpdateDto[struct {
				CountryCode *enums.CountryCode `json:"countryCode" validate:"omitnil,iscountrycode"`
				BackupEmail *string            `json:"backupEmail" validate:"omitnil,email"`
				PhoneNumber *string            `json:"phoneNumber" validate:"omitnil,max=0,max=15,isnumberstring"`
			}]
		},
		any,
	]
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
