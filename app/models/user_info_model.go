package models

import (
	"time"

	uuid "github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	exceptions "notezy-backend/app/exceptions"
	util "notezy-backend/app/util"
	global "notezy-backend/global"
)

/* ============================== Schema ============================== */
type UserInfo struct {
	Id                 uuid.UUID  `json:"id" gorm:"column:id; type:uuid; primaryKey; default:gen_random_uuid();"`
	UserId             uuid.UUID  `json:"userId" gorm:"column:user_id; type:uuid; not null; unique;"`
	CoverBackgroundURL string     `json:"coverBackgroundURL" gorm:"column:cover_background_url; not null; default:''"`
	AvatarURL          string     `json:"avatarURL" gorm:"column:avatar_url; not null; default:''"`
	Header             string     `json:"header" gorm:"column:header; not null; default:''; size:64;"`
	Introduction       string     `json:"introduction" gorm:"column:introduction; not null; default:''; size:256;"`
	Gender             UserGender `json:"gender" gorm:"column:gender; type:UserGender; not null; default:'PreferNotToSay'"`
	Country            Country    `json:"country" gorm:"column:country; type:Country; not null; default:'Default'"`
	BirthDate          time.Time  `json:"birthDate" gorm:"column:birth_date; type:timestamptz; not null; default:CURRENT_TIMESTAMP"`
	UpdatedAt          time.Time  `json:"updatedAt" gorm:"column:updated_at; type:timestamptz; not null; autoUpdateTime:true;"`
}

func (UserInfo) TableName() string {
	return string(global.ValidTableName_UserInfoTable)
}

/* ============================== Input ============================== */
type CreateUserInfoInput struct {
	CoverBackgroundURL *string     `json:"coverBackgroundURL" validate:"omitempty" gorm:"column:cover_background_url;"`
	AvatarURL          *string     `json:"avatarURL" validate:"omitempty" gorm:"column:avatar_url;"`
	Header             *string     `json:"header" validate:"omitempty,min=0,max=64" gorm:"column:header;"`
	Introduction       *string     `json:"introduction" validate:"omitempty,min=0,max=256" gorm:"column:introduction;"`
	Gender             *UserGender `json:"gender" validate:"omitempty,isgender" gorm:"column:gender;"`
	Country            *Country    `json:"country" validate:"omitempty,iscountry" gorm:"column:country;"`
	BirthDate          *time.Time  `json:"birthDate" validate:"omitempty" gorm:"column:birth_date;"`
}

type UpdateUserInfoInput struct {
	CoverBackgroundURL *string     `json:"coverBackgroundURL" validate:"omitempty" gorm:"column:cover_background_url;"`
	AvatarURL          *string     `json:"avatarURL" validate:"omitempty" gorm:"column:avatar_url;"`
	Header             *string     `json:"header" validate:"omitempty,min=0,max=64" gorm:"column:header;"`
	Introduction       *string     `json:"introduction" validate:"omitempty,min=0,max=256" gorm:"column:introduction;"`
	Gender             *UserGender `json:"gender" validate:"omitempty,isgender" gorm:"column:gender;"`
	Country            *Country    `json:"country" validate:"omitempty,iscountry" gorm:"column:country;"`
	BirthDate          *time.Time  `json:"birthDate" validate:"omitempty" gorm:"column:birth_date;"`
}

/* ============================== Methods ============================== */
func GetUserInfoByUserId(db *gorm.DB, userId uuid.UUID) (*UserInfo, *exceptions.Exception) {
	if db == nil {
		db = NotezyDB
	}

	userInfo := UserInfo{}
	result := db.Table(UserInfo{}.TableName()).Where("user_id = ?", userId).First(&userInfo)
	if err := result.Error; err != nil {
		return nil, exceptions.UserInfo.NotFound().WithError(err)
	}

	return &userInfo, nil
}

func CreateUserInfoByUserId(db *gorm.DB, userId uuid.UUID, input CreateUserInfoInput) (*uuid.UUID, *exceptions.Exception) {
	if db == nil {
		db = NotezyDB
	}

	if err := Validator.Struct(input); err != nil {
		return nil, exceptions.UserInfo.InvalidInput().WithError(err)
	}

	var newUserInfo UserInfo
	newUserInfo.UserId = userId
	util.CopyNonNilFields(&newUserInfo, input)

	result := db.Table(UserInfo{}.TableName()).
		Clauses(clause.Returning{Columns: []clause.Column{
			{Name: "id"},
		}}).
		Create(&newUserInfo)
	if err := result.Error; err != nil {
		return nil, exceptions.UserInfo.FailedToCreate().WithError(err)
	}

	return &newUserInfo.Id, nil
}

func UpdateUserInfoByUserId(db *gorm.DB, userId uuid.UUID, input UpdateUserInfoInput) (*UserInfo, *exceptions.Exception) {
	if db == nil {
		db = NotezyDB
	}

	if err := Validator.Struct(input); err != nil {
		return nil, exceptions.UserInfo.InvalidInput().WithError(err)
	}

	var updatedUserInfo UserInfo
	util.CopyNonNilFields(&updatedUserInfo, input)
	result := db.Table(UserInfo{}.TableName()).
		Where("user_id = ?", userId).
		Clauses(clause.Returning{}).
		Updates(&updatedUserInfo)
	if err := result.Error; err != nil {
		return nil, exceptions.UserInfo.FailedToUpdate().WithError(err)
	}

	return &updatedUserInfo, nil
}
