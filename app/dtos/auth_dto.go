package dtos

import (
	"notezy-backend/app/models"
	"time"
)

type RegisterReqDto struct {
	CreateUserInputData        models.CreateUserInput
	CreateUserInfoInputData    models.CreateUserInfoInput
	CreateUserAccountInputData models.CreateUserAccountInput
	CreateUserSettingInputData models.CreateUserSettingInput
}

type RegisterResDto struct {
	AccessToken string
	CreatedAt   time.Time
}
