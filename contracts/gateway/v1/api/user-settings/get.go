package usersettingsdto

import apiv1 "github.com/HiIamJeff67/notezy-backend/contracts/gateway/v1/api"

import enums "github.com/HiIamJeff67/notezy-backend/internal/services/core/data/database/schemas/enums"

type GetMySettingRequestDto struct {
	apiv1.RequestDto[
		struct {
			UserAgent string `json:"userAgent" validate:"required,isuseragent"`
		},
		struct{},
		struct{},
		struct{},
	]
}

type GetMySettingResponseDto struct {
	Language           enums.Language `json:"language"`
	GeneralSettingCode int64          `json:"generalSettingCode"`
	PrivacySettingCode int64          `json:"privacySettingCode"`
}
