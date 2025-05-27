package models

import (
	"go-gorm-api/global"
	"time"

	uuid "github.com/jackc/pgx/pgtype/ext/satori-uuid"
)

/* ============================== Schema ============================== */
type UserSetting struct {
	Id        			uuid.UUID  	`json:"id" gorm:"column:id; type:uuid; primaryKey; default:gen_random_uuid();"`
	UserId				uuid.UUID	`json:"userId" gorm:"column:user_id; type:uuid; not null; unique;"`
	Theme				Theme       `json:"theme" gorm:"column:theme; type:Theme; not null; default:'System';"`
	Language    		Language	`json:"language" gorm:"column:language; type:Language; not null; default:'English';"`
	GeneralSettingCode	int64		`json:"generalSettingCode" gorm:"column:general_setting_code; type:bigint; not null; default:0;"`
	PrivacySettingCode  int64		`json:"privacySettingCode" gorm:"column:privacy_setting_code; type:bigint; not null; default:0;"`
	UpdatedAt 			time.Time  	`json:"updatedAt" gorm:"column:updated_at; type:timestamptz; not null; autoUpdateTime:true;"`
}

func (UserSetting) TableName() string {
	return string(global.ValidTableName_UserSettingTable)
}
/* ============================== Schema ============================== */

/* ============================== Input ============================== */
type CreateUserSettingInput struct {
	Theme				Theme		`json:"theme" gorm:"column:theme;"`
	Language			Language	`json:"language" gorm:"column:language;"`
	GeneralSettingCode  int64		`json:"generalSettingCode" gorm:"column:general_setting_code;"`
	PrivacySettingCode  int64		`json:"privacySettingCode" gorm:"column:privacy_setting_code;"`
}
type UpdateUserSettingInput struct {
	Theme				*Theme		`json:"theme" gorm:"column:theme;"`
	Language			*Language	`json:"language" gorm:"column:language;"`
	GeneralSettingCode  *int64		`json:"generalSettingCode" gorm:"column:general_setting_code;"`
	PrivacySettingCode  *int64		`json:"privacySettingCode" gorm:"column:privacy_setting_code;"`
}
/* ============================== Input ============================== */

/* ============================== Methods ============================== */
func GetUserSettingByUserId(userId uuid.UUID) (userSetting UserSetting, err error) {
	result := NotezyDB.Table(UserSetting{}.TableName()).Where("user_id = ?", userId).First(&userSetting)
	if err = result.Error; err != nil {
		return UserSetting{}, err
	}
	return userSetting, nil
}

func CreateUserSettingByUserId(userId uuid.UUID, input CreateUserSettingInput) (newUserSetting UserSetting, err error) {
	newUserSetting = UserSetting{
		UserId: userId, 
		Theme: input.Theme, 
		Language: input.Language,
		GeneralSettingCode: input.GeneralSettingCode,
		PrivacySettingCode: input.PrivacySettingCode,
	}

	result := NotezyDB.Table(UserSetting{}.TableName()).Create(&newUserSetting)
	if err = result.Error; err != nil {
		return UserSetting{}, err
	}
	return newUserSetting, nil
}

func UpdateUserSettingByUserId(userId uuid.UUID, input UpdateUserSettingInput) (updatedUserSetting UserSetting, err error) {
	updatedUserSetting = UserSetting{
		UserId: userId, 
		Theme: *input.Theme, 
		Language: *input.Language,
		GeneralSettingCode: *input.GeneralSettingCode,
		PrivacySettingCode: *input.PrivacySettingCode,
	}

	result := NotezyDB.Table(UserSetting{}.TableName()).Create(&updatedUserSetting)
	if err = result.Error; err != nil {
		return UserSetting{}, err
	}
	return updatedUserSetting, nil
}
/* ============================== Methods ============================== */