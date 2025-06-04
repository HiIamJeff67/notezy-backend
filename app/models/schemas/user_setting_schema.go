package schemas

import (
	"notezy-backend/app/models/enums"
	"notezy-backend/global"
	"time"

	"github.com/google/uuid"
)

type UserSetting struct {
	Id                 uuid.UUID      `json:"id" gorm:"column:id; type:uuid; primaryKey; default:gen_random_uuid();"`
	UserId             uuid.UUID      `json:"userId" gorm:"column:user_id; type:uuid; not null; unique;"`
	Theme              enums.Theme    `json:"theme" gorm:"column:theme; type:Theme; not null; default:'System';"`
	Language           enums.Language `json:"language" gorm:"column:language; type:Language; not null; default:'English';"`
	GeneralSettingCode int            `json:"generalSettingCode" gorm:"column:general_setting_code; type:integer; not null; default:0;"`
	PrivacySettingCode int            `json:"privacySettingCode" gorm:"column:privacy_setting_code; type:integer; not null; default:0;"`
	UpdatedAt          time.Time      `json:"updatedAt" gorm:"column:updated_at; type:timestamptz; not null; autoUpdateTime:true;"`
}

func (UserSetting) TableName() string {
	return global.ValidTableName_UserSettingTable.String()
}
