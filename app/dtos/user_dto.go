package dtos

import (
	"notezy-backend/app/caches"
)

type FindMeReqDto struct {
	AccessToken string
}

type FindMeResDto = caches.UserDataCache
