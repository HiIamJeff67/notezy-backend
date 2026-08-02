package usersdto

import (
	"time"

	"github.com/google/uuid"

	apiv1 "github.com/HiIamJeff67/notezy-backend/contracts/gateway/v1/api"
)

type GetUserDataRequestDto struct {
	apiv1.RequestDto[
		struct {
			UserAgent string `json:"userAgent" validate:"required,isuseragent"`
		},
		struct{},
		struct{},
		struct{},
	]
}

type GetUserDataResponseDto struct {
	PublicId           uuid.UUID `json:"publicId"`
	Name               string    `json:"name"`
	DisplayName        string    `json:"displayName"`
	Email              string    `json:"email"`
	Role               string    `json:"role"`
	Plan               string    `json:"plan"`
	Status             string    `json:"status"`
	AvatarURL          string    `json:"avatarURL"`
	Language           string    `json:"language"`
	GeneralSettingCode int64     `json:"generalSettingCode"`
	PrivacySettingCode int64     `json:"privacySettingCode"`
	CreatedAt          time.Time `json:"createdAt"`
	UpdatedAt          time.Time `json:"updatedAt"`
}

type GetMeRequestDto struct {
	apiv1.RequestDto[
		struct {
			UserAgent string `json:"userAgent" validate:"required,isuseragent"`
		},
		struct{},
		struct{},
		struct{},
	]
}

type GetMeResponseDto struct {
	PublicId    uuid.UUID `json:"publicId"`
	Name        string    `json:"name"`
	DisplayName string    `json:"displayName"`
	Email       string    `json:"email"`
	Role        string    `json:"role"`
	Plan        string    `json:"plan"`
	Status      string    `json:"status"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}
