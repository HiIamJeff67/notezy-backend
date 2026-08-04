package authdto

import (
	"time"

	coreapicontract "github.com/HiIamJeff67/notezy-backend/contracts/core/v1/api"
)

type ForgetPasswordRequestDto struct {
	coreapicontract.RequestDto[
		struct {
			UserAgent string `json:"userAgent" validate:"required,isuseragent"`
		},
		struct {
			Account     string `json:"account" validate:"required,isaccount"`
			NewPassword string `json:"newPassword" validate:"required,min=8,max=1024,isstrongpassword"`
			AuthCode    string `json:"authCode" validate:"required,isnumberstring,len=6"`
		},
		struct{},
		struct{},
	]
}

type ForgetPasswordResponseDto struct {
	UpdatedAt time.Time `json:"updatedAt"`
}
