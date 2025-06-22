package dtos

import (
	"notezy-backend/app/models/enums"
	"time"

	"github.com/google/uuid"
)

/* ============================== Request DTO ============================== */

type GetMySettingReqDto struct {
	UserId uuid.UUID // extracted from the access token of AuthMiddleware()
}

type UpdateMySettingReqDto struct {
	UserId uuid.UUID // extracted from the access token of AuthMiddleware()
	PartialUpdateDto[struct {
		Language           enums.Language `json:"language" validate:"omitempty"`
		GeneralSettingCode int            `json:"generalSettingCode" validate:"omitempty"`
		PrivacySettingCode int            `json:"privacySettingCode" validate:"omitempty"`
	}]
}

/* ============================== Response DTO ============================== */

type GetMySettingResDto struct {
	Language           enums.Language `json:"language"`
	GeneralSettingCode int            `json:"generalSettingCode"`
	PrivacySettingCode int            `json:"privacySettingCode"`
}

type UpdateMySettingResDto struct {
	UpdatedAt time.Time `json:"updatedAt"`
}
