package dtos

import (
	"time"

	"github.com/google/uuid"

	caches "notezy-backend/app/caches"
	enums "notezy-backend/app/models/schemas/enums"
)

/* ============================== Request DTO ============================== */

type GetUserDataReqDto struct {
	NotezyRequest[
		any,
		struct {
			UserId uuid.UUID // extracted from the access token of AuthMiddleware()
		},
		any,
	]
}

type GetMeReqDto struct {
	NotezyRequest[
		any,
		struct {
			UserId uuid.UUID // extracted from the access token of AuthMiddleware()
		},
		any,
	]
}

type UpdateMeReqDto struct {
	NotezyRequest[
		any,
		struct {
			UserId uuid.UUID // extracted from the access token of AuthMiddleware()
		},
		struct {
			PartialUpdateDto[struct {
				DisplayName *string           `json:"displayName" validate:"omitnil,min=6,max=32,alphaandnum"`
				Status      *enums.UserStatus `json:"status" validate:"omitnil,isstatus"`
			}]
		},
	]
}

type UpdateRoleReqDto struct {
	NotezyRequest[
		any,
		struct {
			UserId uuid.UUID // extracted from the access token of AuthMiddleware()
		},
		struct {
			Role enums.UserRole `json:"role" validate:"required,isrole"`
		},
	]
}

type UpdatePlanReqDto struct {
	NotezyRequest[
		any,
		struct {
			UserId uuid.UUID // extracted from the access token of AuthMiddleware()
		},
		struct {
			Plan enums.UserPlan `json:"plan" validate:"required,isplan"`
		},
	]
}

/* ============================== Response DTO ============================== */

type GetUserDataResDto = caches.UserDataCache

type GetMeResDto struct {
	PublicId    string           `json:"publicId"`
	Name        string           `json:"name"`
	DisplayName string           `json:"displayName"`
	Email       string           `json:"email"`
	Role        enums.UserRole   `json:"role"`
	Plan        enums.UserPlan   `json:"plan"`
	Status      enums.UserStatus `json:"status"`
	CreatedAt   time.Time        `json:"createdAt"`
}

type UpdateMeResDto struct {
	UpdatedAt time.Time `json:"updatedAt"`
}

type UpdateRoleResDto struct {
	UpdatedAt time.Time `json:"updatedAt"`
}

type UpdatePlanResDto struct {
	UpdatedAt time.Time `json:"updatedAt"`
}
