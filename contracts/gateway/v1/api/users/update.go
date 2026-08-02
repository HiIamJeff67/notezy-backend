package usersdto

import (
	"time"

	apiv1 "github.com/HiIamJeff67/notezy-backend/contracts/gateway/v1/api"
)

type UpdateMeRequestDto struct {
	apiv1.RequestDto[
		struct {
			UserAgent string `json:"userAgent" validate:"required,isuseragent"`
		},
		struct {
			Values struct {
				DisplayName *string `json:"displayName" validate:"omitnil,min=6,max=32,alphaandnum"`
				Status      *string `json:"status" validate:"omitnil,isstatus"`
			} `json:"values"`
			SetNull *map[string]bool `json:"setNull,omitempty"`
		},
		struct{},
		struct{},
	]
}

type UpdateMeResponseDto struct {
	UpdatedAt time.Time `json:"updatedAt"`
}
