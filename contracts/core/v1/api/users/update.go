package apicontract

import (
	"time"

	coreapicontract "github.com/HiIamJeff67/notezy-backend/contracts/core/v1/api"
)

type UpdateMeRequestDto struct {
	coreapicontract.RequestDto[
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
