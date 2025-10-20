package dtos

import (
	"time"

	"github.com/google/uuid"
)

/* ============================== Request DTO ============================== */
// make sure do NOT use the access token or refresh token as the request dto

type RegisterReqDto struct {
	NotezyRequest[
		struct {
			UserAgent string `json:"userAgent" validate:"required,isuseragent"`
		},
		any,
		struct {
			Name     string `json:"name" validate:"required,min=6,max=16,alphaandnum"`
			Email    string `json:"email" validate:"required,email"`
			Password string `json:"password" validate:"required,min=8,max=1024,isstrongpassword"`
		},
		any,
	]
}

type LoginReqDto struct {
	NotezyRequest[
		struct {
			UserAgent string `json:"userAgent" validate:"required,isuseragent"`
		},
		any,
		struct {
			Account  string `json:"account" validate:"required,account"`
			Password string `json:"password" validate:"required"` // don't validate other additions while login
		},
		any,
	]
}

type LogoutReqDto struct {
	NotezyRequest[
		struct {
			UserAgent string `json:"userAgent" validate:"required,isuseragent"`
		},
		struct {
			UserId uuid.UUID // extracted from the access token of AuthMiddleware
		},
		any,
		any,
	]
}

type SendAuthCodeReqDto struct {
	NotezyRequest[
		struct {
			UserAgent string `json:"userAgent" validate:"required,isuseragent"`
		},
		any,
		struct {
			Email string `json:"email" validate:"required,email"`
		},
		any,
	]
}

type ValidateEmailReqDto struct {
	NotezyRequest[
		struct {
			UserAgent string `json:"userAgent" validate:"required,isuseragent"`
		},
		struct {
			UserId uuid.UUID // extracted from the access token of AuthMiddleware()
		},
		struct {
			AuthCode string `json:"authCode" validate:"required,isnumberstring,len=6"`
		},
		any,
	]
}

type ResetEmailReqDto struct {
	NotezyRequest[
		struct {
			UserAgent string `json:"userAgent" validate:"required,isuseragent"`
		},
		struct {
			UserId uuid.UUID // extracted from the access token of AuthMiddleware()
		},
		struct {
			NewEmail string `json:"newEmail" validate:"required,email"`
			AuthCode string `json:"authCode" validate:"required,isnumberstring,len=6"`
		},
		any,
	]
}

type ForgetPasswordReqDto struct {
	NotezyRequest[
		struct {
			UserAgent string `json:"userAgent" validate:"required,isuseragent"`
		},
		any,
		struct {
			Account     string `json:"account" validate:"required,account"`
			NewPassword string `json:"newPassword" validation:"required,min=8,max=1024,isstrongpassword"`
			AuthCode    string `json:"authCode" validate:"required,isnumberstring,len=6"`
		},
		any,
	]
}

type DeleteMeReqDto struct {
	NotezyRequest[
		struct {
			UserAgent string `json:"userAgent" validate:"required,isuseragent"`
		},
		struct {
			UserId uuid.UUID // extracted from the access token of AuthMiddleware()
		},
		struct {
			AuthCode time.Time `json:"authCode" validate:"required,isnumberstring,len=6"`
		},
		any,
	]
}

/* ============================== Response DTO ============================== */
type RegisterResDto struct {
	AccessToken  string    `json:"accessToken"`
	RefreshToken string    `json:"refreshToken"`
	CSRFToken    string    `json:"csrfToken"`
	CreatedAt    time.Time `json:"createdAt"`
}

type LoginResDto struct {
	AccessToken  string    `json:"accessToken"`
	RefreshToken string    `json:"refreshToken"`
	CSRFToken    string    `json:"csrfToken"`
	UpdatedAt    time.Time `json:"updatedAt"`
}

type LogoutResDto struct {
	UpdatedAt time.Time `json:"updatedAt"`
}

type SendAuthCodeResDto struct {
	AuthCodeExpiredAt  time.Time `json:"authCodeExpiredAt"`
	BlockAuthCodeUntil time.Time `json:"blockAuthCodeUntil"`
	UpdatedAt          time.Time `json:"updatedAt"`
}

type ValidateEmailResDto struct {
	UpdatedAt time.Time `json:"updatedAt"`
}

type ResetEmailResDto struct {
	UpdatedAt time.Time `json:"updatedAt"`
}

type ForgetPasswordResDto struct {
	UpdatedAt time.Time `json:"updatedAt"`
}

type DeleteMeResDto struct {
	DeletedAt time.Time `json:"deletedAt"`
}
