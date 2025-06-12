package dtos

import (
	"notezy-backend/app/caches"
	"notezy-backend/app/models/enums"
	"time"
)

type FindMeReqDto struct {
	AccessToken string
}

type UpdateMeReqDto struct {
	AccessToken string
	DisplayName string
	Email       string
	Password    string
	Status      enums.UserStatus
}

type FindMeResDto = caches.UserDataCache

type UpdateMeResDto struct {
	UpdatedAt time.Time
}
