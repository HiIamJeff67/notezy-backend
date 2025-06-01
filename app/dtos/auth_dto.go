package dtos

import (
	"time"
)

type RegisterReqDto struct {
	Name     string `json:"name" validate:"required"`
	Email    string `json:"email" validate:"required"`
	Password string `json:"password" validate:"required,min=8,max=32,isstrongpassword"`
}

type RegisterResDto struct {
	AccessToken string    `json:"accessToken"`
	CreatedAt   time.Time `json:"createdAt"`
}

type LoginReqDto struct {
	AccessToken  *string `json:"accessToken" validate:"omitempty"`
	RefreshToken *string `json:"refreshToken" validate:"omitempty"`
	Account      *string `json:"account" validate:"omitempty"`
	Password     *string `json:"password" validate:"omitempty,min=8,max=32,isstrongpassword"`
}

type LoginResDto struct {
	AccessToken string    `json:"accessToken"`
	CreatedAt   time.Time `json:"createdAt"`
}
