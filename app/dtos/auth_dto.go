package dtos

import (
	"time"

	"github.com/google/uuid"
)

/* ============================== Request DTO ============================== */
// make sure do NOT use the access token or refresh token as the request dto

type RegisterReqDto struct {
	Name      string `json:"name" validate:"required"`
	Email     string `json:"email" validate:"required"`
	Password  string `json:"password" validate:"required,min=8,max=32,isstrongpassword"`
	UserAgent string `json:"userAgent" validate:"required" gorm:"column:user_agent;"`
}

type LoginReqDto struct {
	Account   string `json:"account" validate:"required"`
	Password  string `json:"password" validate:"required"` // don't validate other additions while login
	UserAgent string `json:"userAgent" validate:"required" gorm:"column:user_agent;"`
}

type LogoutReqDto struct {
	UserId uuid.UUID // extracted from the access token of authMidddleware
}

type SendAuthCodeReqDto struct {
	Email string `json:"email" validate:"required,email"`
}

/* ============================== Response DTO ============================== */
type RegisterResDto struct {
	AccessToken  string    `json:"accessToken"`
	RefreshToken string    `json:"refreshToken"`
	CreatedAt    time.Time `json:"createdAt"`
}

type LoginResDto struct {
	AccessToken string    `json:"accessToken"`
	UpdatedAt   time.Time `json:"updatedAt"`
}

type LogoutResDto struct {
	UpdatedAt time.Time `json:"updatedAt"`
}

type SendAuthCodeResDto struct {
	AuthCodeExpiredAt time.Time `json:"authCodeExpiredAt"`
	UpdatedAt         time.Time `json:"updatedAt"`
}
