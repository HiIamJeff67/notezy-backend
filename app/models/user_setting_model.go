package models

import (
	"notezy-backend/app/exceptions"
	"notezy-backend/app/util"
	"notezy-backend/global"
	"time"

	uuid "github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

/* ============================== Schema ============================== */
type UserSetting struct {
	Id                 uuid.UUID `json:"id" gorm:"column:id; type:uuid; primaryKey; default:gen_random_uuid();"`
	UserId             uuid.UUID `json:"userId" gorm:"column:user_id; type:uuid; not null; unique;"`
	Theme              Theme     `json:"theme" gorm:"column:theme; type:Theme; not null; default:'System';"`
	Language           Language  `json:"language" gorm:"column:language; type:Language; not null; default:'English';"`
	GeneralSettingCode int       `json:"generalSettingCode" gorm:"column:general_setting_code; type:integer; not null; default:0;"`
	PrivacySettingCode int       `json:"privacySettingCode" gorm:"column:privacy_setting_code; type:integer; not null; default:0;"`
	UpdatedAt          time.Time `json:"updatedAt" gorm:"column:updated_at; type:timestamptz; not null; autoUpdateTime:true;"`
}

func (UserSetting) TableName() string {
	return string(global.ValidTableName_UserSettingTable)
}

/* ============================== Input ============================== */
type CreateUserSettingInput struct {
	Theme              *Theme    `json:"theme" validate:"istheme" gorm:"column:theme;"`
	Language           *Language `json:"language" validate:"islanguage" gorm:"column:language;"`
	GeneralSettingCode *int      `json:"generalSettingCode" validate:"min=0,max=999999999" gorm:"column:general_setting_code;"`
	PrivacySettingCode *int      `json:"privacySettingCode" validate:"min=0,max=999999999" gorm:"column:privacy_setting_code;"`
}
type UpdateUserSettingInput struct {
	Theme              *Theme    `json:"theme" validate:"istheme" gorm:"column:theme;"`
	Language           *Language `json:"language" validate:"islanguage" gorm:"column:language;"`
	GeneralSettingCode *int      `json:"generalSettingCode" validate:"min=0,max=999999999" gorm:"column:general_setting_code;"`
	PrivacySettingCode *int      `json:"privacySettingCode" validate:"min=0,max=999999999" gorm:"column:privacy_setting_code;"`
}

/* ============================== Methods ============================== */
func GetUserSettingByUserId(db *gorm.DB, userId uuid.UUID) (*UserSetting, *exceptions.Exception) {
	if db == nil {
		db = NotezyDB
	}

	userSetting := UserSetting{}
	result := db.Table(UserSetting{}.TableName()).Where("user_id = ?", userId).First(&userSetting)
	if err := result.Error; err != nil {
		return nil, exceptions.UserSetting.NotFound().WithError(err)
	}

	return &userSetting, nil
}

func CreateUserSettingByUserId(db *gorm.DB, userId uuid.UUID, input CreateUserSettingInput) (*uuid.UUID, *exceptions.Exception) {
	if db == nil {
		db = NotezyDB
	}

	var newUserSetting UserSetting
	newUserSetting.UserId = userId
	util.CopyNonNilFields(&newUserSetting, input)
	result := db.Table(UserSetting{}.TableName()).
		Clauses(clause.Returning{Columns: []clause.Column{
			{Name: "id"},
		}}).
		Create(&newUserSetting)
	if err := result.Error; err != nil {
		return nil, exceptions.UserSetting.FailedToCreate().WithError(err)
	}

	return &newUserSetting.Id, nil
}

func UpdateUserSettingByUserId(db *gorm.DB, userId uuid.UUID, input UpdateUserSettingInput) (*UserSetting, *exceptions.Exception) {
	if db == nil {
		db = NotezyDB
	}

	var updatedUserSetting UserSetting
	updatedUserSetting.UserId = userId
	util.CopyNonNilFields(&updatedUserSetting, input)
	result := db.Table(UserSetting{}.TableName()).
		Where("user_id = ?", userId).
		Clauses(clause.Returning{}).
		Create(&updatedUserSetting)
	if err := result.Error; err != nil {
		return nil, exceptions.UserSetting.FailedToUpdate().WithError(err)
	}

	return &updatedUserSetting, nil
}
