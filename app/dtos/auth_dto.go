package dtos

import (
	"time"
)

/* ============================== Request DTO ============================== */

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
	AccessToken string `json:"accessToken" validate:"required"`
}

type HeartBeatReqDto struct {
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
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

type HeartBeatResDto struct {
	AccessToken string    `json:"accessToken"`
	UpdatedAt   time.Time `json:"updatedAt"`
}
