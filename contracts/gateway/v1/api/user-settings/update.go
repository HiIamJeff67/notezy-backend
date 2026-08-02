package usersettingsdto

import (
	"time"

	apiv1 "github.com/HiIamJeff67/notezy-backend/contracts/gateway/v1/api"
	enums "github.com/HiIamJeff67/notezy-backend/internal/services/core/data/database/schemas/enums"
)

type UpdateMySettingRequestDto struct {
	apiv1.RequestDto[
		struct {
			UserAgent string `json:"userAgent" validate:"required,isuseragent"`
		},
		struct {
			Values struct {
				Language           *enums.Language `json:"language" validate:"omitnil,islanguage"`
				GeneralSettingCode *int64          `json:"generalSettingCode" validate:"omitnil,min=0,max=999999999"`
				PrivacySettingCode *int64          `json:"privacySettingCode" validate:"omitnil,min=0,max=999999999"`
			} `json:"values"`
			SetNull *map[string]bool `json:"setNull,omitempty"`
		},
		struct{},
		struct{},
	]
}

type UpdateMySettingResponseDto struct {
	UpdatedAt time.Time `json:"updatedAt"`
}
