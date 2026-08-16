package apicontract

import coreapicontract "github.com/HiIamJeff67/notegic-backend/contracts/core/v1/api"

import enumcontract "github.com/HiIamJeff67/notegic-backend/contracts/types/enums"

type GetMySettingRequestDto struct {
	coreapicontract.RequestDto[
		struct {
			UserAgent string `json:"userAgent" validate:"required,isuseragent"`
		},
		struct{},
		struct{},
		struct{},
	]
}

type GetMySettingResponseDto struct {
	Language           enumcontract.Language `json:"language"`
	GeneralSettingCode int64                 `json:"generalSettingCode"`
	PrivacySettingCode int64                 `json:"privacySettingCode"`
}
