package dtos

import (
	"time"

	"github.com/google/uuid"

	enums "notezy-backend/app/models/schemas/enums"
)

/* ============================== Request DTO ============================== */

type GetMySettingReqDto struct {
	UserId uuid.UUID // extracted from the access token of AuthMiddleware()
}

type UpdateMySettingReqDto struct {
	UserId uuid.UUID // extracted from the access token of AuthMiddleware()
	PartialUpdateDto[struct {
		Language           enums.Language `json:"language" validate:"omitempty"`
		GeneralSettingCode int64          `json:"generalSettingCode" validate:"omitempty"`
		PrivacySettingCode int64          `json:"privacySettingCode" validate:"omitempty"`
	}]
}

/* ============================== Response DTO ============================== */

type GetMySettingResDto struct {
	Language           enums.Language `json:"language"`
	GeneralSettingCode int64          `json:"generalSettingCode"`
	PrivacySettingCode int64          `json:"privacySettingCode"`
}

type UpdateMySettingResDto struct {
	UpdatedAt time.Time `json:"updatedAt"`
}
