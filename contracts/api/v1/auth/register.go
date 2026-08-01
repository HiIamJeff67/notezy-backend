package authdto

import (
	"time"

	"github.com/google/uuid"

	apiv1 "github.com/HiIamJeff67/notezy-backend/contracts/api/v1"
)

type RegisterRequestDto struct {
	apiv1.RequestDto[
		struct {
			UserAgent string `json:"userAgent" validate:"required,isuseragent"`
		},
		struct {
			Name     string `json:"name" validate:"required,min=6,max=32,alphaandnum"`
			Email    string `json:"email" validate:"required,email"`
			Password string `json:"password" validate:"required,min=8,max=1024,isstrongpassword"`
		},
		struct{},
		struct{},
	]
}

type RegisterResponseDto struct {
	PublicId     uuid.UUID `json:"publicId"`
	Name         string    `json:"name"`
	DisplayName  string    `json:"displayName"`
	Email        string    `json:"email"`
	AccessToken  string    `json:"accessToken"`
	RefreshToken string    `json:"refreshToken"`
	CSRFToken    string    `json:"csrfToken"`
	CreatedAt    time.Time `json:"createdAt"`
}

type RegisterViaGoogleRequestDto struct {
	apiv1.RequestDto[
		struct {
			UserAgent string `json:"userAgent" validate:"required,isuseragent"`
		},
		struct {
			AuthorizationCode string `json:"authorizationCode" validate:"required"`
		},
		struct{},
		struct{},
	]
}

type RegisterViaGoogleResponseDto struct {
	PublicId     uuid.UUID `json:"publicId"`
	Name         string    `json:"name"`
	DisplayName  string    `json:"displayName"`
	Email        string    `json:"email"`
	AccessToken  string    `json:"accessToken"`
	RefreshToken string    `json:"refreshToken"`
	CSRFToken    string    `json:"csrfToken"`
	CreatedAt    time.Time `json:"createdAt"`
}
