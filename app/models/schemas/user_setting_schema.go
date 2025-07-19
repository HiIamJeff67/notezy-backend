package schemas

import (
	"time"

	"github.com/google/uuid"

	enums "notezy-backend/app/models/schemas/enums"
	shared "notezy-backend/shared"
)

type UserSetting struct {
	Id                 uuid.UUID      `json:"id" gorm:"column:id; type:uuid; primaryKey; default:gen_random_uuid();"`
	UserId             uuid.UUID      `json:"userId" gorm:"column:user_id; type:uuid; not null; unique;"`
	Language           enums.Language `json:"language" gorm:"column:language; type:Language; not null; default:'English';"`
	GeneralSettingCode int64          `json:"generalSettingCode" gorm:"column:general_setting_code; type:bigint; not null; default:0;"`
	PrivacySettingCode int64          `json:"privacySettingCode" gorm:"column:privacy_setting_code; type:bigint; not null; default:0;"`
	UpdatedAt          time.Time      `json:"updatedAt" gorm:"column:updated_at; type:timestamptz; not null; autoUpdateTime:true;"`
}

func (UserSetting) TableName() string {
	return shared.ValidTableName_UserSettingTable.String()
}

/* ============================== Relative Type Conversion ============================== */
// note that there's no type like PublicUserSetting,
// since the userSetting shouldn't be public
