package dtos

import (
	"time"

	"github.com/google/uuid"
)

/* ============================== Request DTO ============================== */
// make sure do NOT use the access token or refresh token as the request dto

type RegisterReqDto struct {
	Name      string `json:"name" validate:"required,min=6,max=16,alphaandnum"`
	Email     string `json:"email" validate:"required,email"`
	Password  string `json:"password" validate:"required,min=8,max=1024,isstrongpassword"`
	UserAgent string `json:"userAgent" validate:"required,isuseragent"`
}

type LoginReqDto struct {
	Account   string `json:"account" validate:"required,account"`
	Password  string `json:"password" validate:"required"` // don't validate other additions while login
	UserAgent string `json:"userAgent" validate:"required,isuseragent"`
}

type LogoutReqDto struct {
	UserId uuid.UUID // extracted from the access token of AuthMiddleware
}

type SendAuthCodeReqDto struct {
	Email     string `json:"email" validate:"required,email"`
	UserAgent string `json:"userAgent" validate:"required,isuseragent"`
}

type ValidateEmailReqDto struct {
	UserId   uuid.UUID // extracted from the access token of AuthMiddleware()
	AuthCode string    `json:"authCode" validate:"required,isnumberstring,len=6"`
}

type ResetEmailReqDto struct {
	UserId   uuid.UUID // extracted from the access token of AuthMiddleware()
	NewEmail string    `json:"newEmail" validate:"required,email"`
	AuthCode string    `json:"authCode" validate:"required,isnumberstring,len=6"`
}

type ForgetPasswordReqDto struct {
	Account     string `json:"account" validate:"required,account"`
	NewPassword string `json:"newPassword" validation:"required,min=8,max=1024,isstrongpassword"`
	AuthCode    string `json:"authCode" validate:"required,isnumberstring,len=6"`
	UserAgent   string `json:"userAgent" validate:"required,isuseragent"`
}

type DeleteMeReqDto struct {
	UserId   uuid.UUID // extracted from the access token of AuthMiddleware()
	AuthCode time.Time `json:"authCode" validate:"required,isnumberstring,len=6"`
}

/* ============================== Response DTO ============================== */
type RegisterResDto struct {
	AccessToken  string    `json:"accessToken"`
	RefreshToken string    `json:"refreshToken"`
	CreatedAt    time.Time `json:"createdAt"`
}

type LoginResDto struct {
	AccessToken  string    `json:"accessToken"`
	RefreshToken string    `json:"refreshToken"`
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
