package stationsdto

import coreapicontract "github.com/HiIamJeff67/notezy-backend/contracts/core/v1/api"

type VisualizeMyTotalCountRequestDto struct {
	coreapicontract.RequestDto[
		struct {
			UserAgent string `json:"userAgent" validate:"required,isuseragent"`
		},
		struct{},
		struct{},
		struct {
			Permission string `json:"permission" validate:"required,isaccesscontrolpermission"`
		},
	]
}

type TotalCountDatumResponseDto struct {
	Id    string `json:"id"`
	X     string `json:"x"`
	Value int64  `json:"value"`
}

type VisualizeMyTotalCountResponseDto struct {
	Data []TotalCountDatumResponseDto `json:"data"`
}
