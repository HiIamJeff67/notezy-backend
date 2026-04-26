package dtos

import (
	"time"

	"github.com/google/uuid"

	enums "notezy-backend/app/models/schemas/enums"
)

/* ============================== Request DTO ============================== */

type GetUserDataReqDto struct {
	NotezyRequest[
		struct {
			UserAgent string `json:"userAgent" validate:"required,isuseragent"`
		},
		struct {
			UserId   uuid.UUID // extracted from the access token of AuthMiddleware()
			UserName string    // extracted from the access token of AuthMiddleware()
		},
		any,
		any,
	]
}

type GetMeReqDto struct {
	NotezyRequest[
		struct {
			UserAgent string `json:"userAgent" validate:"required,isuseragent"`
		},
		struct {
			UserId uuid.UUID // extracted from the access token of AuthMiddleware()
		},
		any,
		any,
	]
}

type UpdateMeReqDto struct {
	NotezyRequest[
		struct {
			UserAgent string `json:"userAgent" validate:"required,isuseragent"`
		},
		struct {
			UserId   uuid.UUID // extracted from the access token of AuthMiddleware()
			UserName string    // extracted from the access token of AuthMiddleware()
		},
		struct {
			PartialUpdateDto[struct {
				DisplayName *string           `json:"displayName" validate:"omitnil,min=6,max=32,alphaandnum"`
				Status      *enums.UserStatus `json:"status" validate:"omitnil,isstatus"`
			}]
		},
		any,
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
		any,
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
		any,
	]
}

/* ============================== Response DTO ============================== */

type GetUserDataResDto = struct {
	PublicId           string           `json:"publicId"`           // user
	Name               string           `json:"name"`               // user
	DisplayName        string           `json:"displayName"`        // user
	Email              string           `json:"email"`              // user
	Role               enums.UserRole   `json:"role"`               // user
	Plan               enums.UserPlan   `json:"plan"`               // user
	Status             enums.UserStatus `json:"status"`             // user
	AvatarURL          string           `json:"avatarURL"`          // user info
	Language           enums.Language   `json:"language"`           // user setting
	GeneralSettingCode int64            `json:"generalSettingCode"` // user setting
	PrivacySettingCode int64            `json:"privacySettingCode"` // user setting
	CreatedAt          time.Time        `json:"createdAt"`          // user
	UpdatedAt          time.Time        `json:"updatedAt"`          // user
}

type GetMeResDto struct {
	PublicId    string           `json:"publicId"`
	Name        string           `json:"name"`
	DisplayName string           `json:"displayName"`
	Email       string           `json:"email"`
	Role        enums.UserRole   `json:"role"`
	Plan        enums.UserPlan   `json:"plan"`
	Status      enums.UserStatus `json:"status"`
	CreatedAt   time.Time        `json:"createdAt"`
	UpdatedAt   time.Time        `json:"updatedAt"`
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
