package apicontract

import (
	"time"

	coreapicontract "github.com/HiIamJeff67/notegic-backend/contracts/core/v1/api"
	enumcontract "github.com/HiIamJeff67/notegic-backend/contracts/types/enums"
)

type UpdateMySettingRequestDto struct {
	coreapicontract.RequestDto[
		struct {
			UserAgent string `json:"userAgent" validate:"required,isuseragent"`
		},
		struct {
			Values struct {
				Language           *enumcontract.Language `json:"language" validate:"omitnil,islanguage"`
				GeneralSettingCode *int64                 `json:"generalSettingCode" validate:"omitnil,min=0,max=999999999"`
				PrivacySettingCode *int64                 `json:"privacySettingCode" validate:"omitnil,min=0,max=999999999"`
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
