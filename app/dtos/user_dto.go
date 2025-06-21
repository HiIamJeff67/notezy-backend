package dtos

import (
	"notezy-backend/app/caches"
	"notezy-backend/app/models/enums"
	"time"

	"github.com/google/uuid"
)

/* ============================== Request DTO ============================== */

type GetMeReqDto struct {
	UserId uuid.UUID // extracted from the access token of AuthMiddleware()
}

type UpdateMeReqDto struct {
	UserId uuid.UUID // extracted from the access token of AuthMiddleware()
	PartialUpdateDto[struct {
		DisplayName *string           `json:"displayName" validate:"omitempty"`
		Status      *enums.UserStatus `json:"status" validate:"omitempty"`
	}]
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

type GetMeResDto = caches.UserDataCache

type UpdateMeResDto struct {
	UpdatedAt time.Time
}

type UpdateRoleResDto struct {
	UpdatedAt time.Time `json:"updatedAt"`
}

type UpdatePlanResDto struct {
	UpdatedAt time.Time `json:"updatedAt"`
}
