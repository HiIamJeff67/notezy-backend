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
				PhoneNumber *string            `json:"phoneNumber" validate:"omitnil,min=1,max=15,isnumberstring"`
			}]
		},
		any,
	]
}

type BindGoogleAccountReqDto struct {
	NotezyRequest[
		struct {
			UserAgent string `json:"userAgent" validate:"required,isuseragent"`
		},
		struct {
			UserId uuid.UUID // extracted from the access token of AuthMiddleware()
		},
		struct {
			AuthorizationCode string `json:"authorizationCode" validate:"required"`
		},
		any,
	]
}

type UnbindGoogleAccountReqDto struct {
	NotezyRequest[
		struct {
			UserAgent string `json:"userAgent" validate:"required,isuseragent"`
		},
		struct {
			UserId uuid.UUID // extracted from the access token of AuthMiddleware()
		},
		struct {
			AuthCode string `json:"authCode" validate:"required"`
		},
		any,
	]
}

/* ============================== Response DTO ============================== */

type GetMyAccountResDto struct {
	CountryCode         *enums.CountryCode `json:"countryCode"`
	PhoneNumber         *string            `json:"phoneNumber"`
	GoogleCredential    *string            `json:"googleCrendential"`
	DiscordCredential   *string            `json:"discordCrendential"`
	RootShelfCount      int32              `json:"rootShelfCount"`
	BlockPackCount      int32              `json:"blockPackCount"`
	BlockCount          int32              `json:"blockCount"`
	MaterialCount       int32              `json:"materialCount"`
	WorkflowCount       int32              `json:"workflowCount"`
	AdditionalItemCount int32              `json:"additionalItemCount"`
	UpdatedAt           time.Time          `json:"updatedAt"`
}

type UpdateMyAccountResDto struct {
	UpdatedAt time.Time `json:"updatedAt"`
}

type BindGoogleAccountResDto struct {
	UpdatedAt time.Time `json:"updatedAt"`
}

type UnbindGoogleAccountResDto struct {
	UpdatedAt time.Time `json:"updatedAt"`
}
