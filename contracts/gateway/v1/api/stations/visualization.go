package stationsdto

import apiv1 "github.com/HiIamJeff67/notezy-backend/contracts/gateway/v1/api"

type VisualizeMyTotalCountRequestDto struct {
	apiv1.RequestDto[
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
