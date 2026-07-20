package inputs

import (
	"time"

	enums "github.com/HiIamJeff67/notezy-backend/app/models/schemas/enums"
	types "github.com/HiIamJeff67/notezy-backend/shared/types"
)

type UpdateUserDataCacheInput struct {
	DisplayName        *string
	Email              *string
	AccessToken        *string
	CSRFToken          *string
	Role               *enums.UserRole
	Plan               *enums.UserPlan
	Status             *enums.UserStatus
	AvatarURL          *string
	Language           *enums.Language
	GeneralSettingCode *int64
	PrivacySettingCode *int64
}

type CheckAndUpdateUserQuotaInput struct {
	Field        types.UserQuotaField
	ChangeAmount int32
	MaxLimit     int32
	ExpiresIn    time.Time
}

type BatchCheckAndUpdateUserQuotaInput struct {
	Identifier string
	Input      CheckAndUpdateUserQuotaInput
}
