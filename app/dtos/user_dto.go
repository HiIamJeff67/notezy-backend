package dtos

import (
	"notezy-backend/app/caches"
	"notezy-backend/app/models/enums"
	"time"

	"github.com/google/uuid"
)

/* ============================== Request DTO ============================== */

type FindMeReqDto struct {
	Id uuid.UUID
}

type UpdateMeReqValues struct {
	DisplayName *string
	Status      *enums.UserStatus
}

type UpdateMeReqDto struct {
	PartialUpdateDto[UpdateMeReqValues]
	AccessToken string
}

type UpdateRoleReqDto struct {
	UserId  uuid.UUID      // extracted from the access token of AuthMiddleware()
	NewRole enums.UserRole `json:"role" validate:"required,isrole"`
}

type UpdatePlanReqDto struct {
	UserId  uuid.UUID      // extracted from the access token of AuthMiddleware()
	NewPlan enums.UserPlan `json:"plan" validate:"required,isplan"`
}

/* ============================== Response DTO ============================== */

type FindMeResDto = caches.UserDataCache

type UpdateMeResDto struct {
	UpdatedAt time.Time
}

type UpdateRoleResDto struct {
	UpdatedAt time.Time `json:"updatedAt"`
}

type UpdatePlanResDto struct {
	UpdatedAt time.Time `json:"updatedAt"`
}
