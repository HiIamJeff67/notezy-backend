package dtos

import (
	"notezy-backend/app/caches"
)

type FindMeReqDto struct {
	AccessToken  *string
	RefreshToken *string
}

type FindMeResDto = caches.UserDataCache
