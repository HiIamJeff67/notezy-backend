package apicontract

import (
	coreapicontract "github.com/HiIamJeff67/notezy-backend/contracts/core/v1/api"
)

type RevokeMyAPIKeyRequestDto struct {
	coreapicontract.RequestDto[
		struct {
			UserAgent string `json:"userAgent" validate:"required,isuseragent"`
		},
		struct{},
		struct {
			PublicId string `json:"publicId" validate:"required,uuid4"`
		},
		struct{},
	]
}

type RevokeMyAPIKeyResponseDto struct {
	RevokedAt string `json:"revokedAt"`
}
