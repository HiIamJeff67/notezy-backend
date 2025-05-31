package dtos

import (
	"notezy-backend/app/models"
	"time"
)

type RegisterReqDto struct {
	CreateUserInputData        models.CreateUserInput        `json:"createUserInputData"`
	CreateUserInfoInputData    models.CreateUserInfoInput    `json:"createUserInfoInputData"`
	CreateUserAccountInputData models.CreateUserAccountInput `json:"createUserAccountInputData"`
	CreateUserSettingInputData models.CreateUserSettingInput `json:"createUserSettingInputData"`
}

type RegisterResDto struct {
	AccessToken string    `json:"accessToken"`
	CreatedAt   time.Time `json:"createdAt"`
}

type LoginReqDto struct {
	AccessToken  *string `json:"accessToken"`
	RefreshToken *string `json:"refreshToken"`
	Account      *string `json:"account"`
	Password     *string `json:"password" validate:"min:8, max:255, isstrongpassword"`
}

type LoginResDto struct {
	AccessToken string    `json:"accessToken"`
	CreatedAt   time.Time `json:"createdAt"`
}
