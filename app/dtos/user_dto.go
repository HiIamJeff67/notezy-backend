package dtos

import (
	"notezy-backend/app/caches"
	"notezy-backend/app/models/enums"
	"time"
)

/* ============================== Request DTO ============================== */

type FindMeReqDto struct {
	AccessToken string
}

type UpdateMeReqValues struct {
	DisplayName *string
	Status      *enums.UserStatus
}

type UpdateMeReqDto struct {
	PartialUpdateDto[UpdateMeReqValues]
	AccessToken string
}

/* ============================== Response DTO ============================== */

type FindMeResDto = caches.UserDataCache

type UpdateMeResDto struct {
	UpdatedAt time.Time
}
