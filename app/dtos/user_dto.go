package dtos

import (
	"time"

	"github.com/google/uuid"

	caches "notezy-backend/app/caches"
	enums "notezy-backend/app/models/schemas/enums"
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
	UserId uuid.UUID      // extracted from the access token of AuthMiddleware()
	Role   enums.UserRole `json:"role" validate:"required,isrole"`
}

type UpdatePlanReqDto struct {
	UserId uuid.UUID      // extracted from the access token of AuthMiddleware()
	Plan   enums.UserPlan `json:"plan" validate:"required,isplan"`
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
