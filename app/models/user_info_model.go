package models

import (
	"notezy-backend/app/exceptions"
	"notezy-backend/global"
	"time"

	uuid "github.com/jackc/pgx/pgtype/ext/satori-uuid"
	"gorm.io/gorm"
)

/* ============================== Schema ============================== */
type UserInfo struct {
	Id                 uuid.UUID  `json:"id" gorm:"column:id; type:uuid; primaryKey; default:gen_random_uuid();"`
	UserId             uuid.UUID  `json:"userId" gorm:"column:user_id; type:uuid; not null; unique;"`
	CoverBackgroundURL *string    `json:"coverBackgroundURL" gorm:"column:cover_background_url;"`
	AvatarURL          *string    `json:"avatarURL" gorm:"column:avatar_url;"`
	Header             *string    `json:"header" gorm:"column:header; not null; default:''; size:64;"`
	Introduction       *string    `json:"introduction" gorm:"column:introduction; not null; default:''; size:256;"`
	Gender             UserGender `json:"gender" gorm:"column:gender; type:UserGender; not null; default:'PreferNotToSay'"`
	Country            Country    `json:"country" gorm:"column:country; type:Country; not null; default:'UnitedStatusOfAmerica'"`
	BirthDate          *time.Time `json:"birthDate" gorm:"column:birth_date; type:timestamptz;"`
	UpdatedAt          time.Time  `json:"updatedAt" gorm:"column:updated_at; type:timestamptz; not null; autoUpdateTime:true;"`
}

func (UserInfo) TableName() string {
	return string(global.ValidTableName_UserInfoTable)
}

/* ============================== Input ============================== */
type CreateUserInfoInput struct {
	CoverBackgroundURL *string     `json:"coverBackgroundURL" gorm:"column:cover_background_url;"`
	AvatarURL          *string     `json:"avatarURL" gorm:"column:avatar_url;"`
	Header             *string     `json:"header" validate:"min:0, max:64" gorm:"column:header;"`
	Introduction       *string     `json:"introduction" validate:"min:0, max:256" gorm:"column:introduction;"`
	Gender             *UserGender `json:"gender" validate:"isgender" gorm:"column:gender;"`
	Country            *Country    `json:"country" validate:"iscountry" gorm:"column:country;"`
	BirthDate          *time.Time  `json:"birthDate" gorm:"column:birth_date;"`
}

type UpdateUserInfoInput struct {
	CoverBackgroundURL *string     `json:"coverBackgroundURL" gorm:"column:cover_background_url;"`
	AvatarURL          *string     `json:"avatarURL" gorm:"column:avatar_url;"`
	Header             *string     `json:"header" validate:"min:0, max:64" gorm:"column:header;"`
	Introduction       *string     `json:"introduction" validate:"min:0, max:256" gorm:"column:introduction;"`
	Gender             *UserGender `json:"gender" validate:"isgender" gorm:"column:gender;"`
	Country            *Country    `json:"country" validate:"iscountry" gorm:"column:country;"`
	BirthDate          *time.Time  `json:"birthDate" gorm:"column:birth_date;"`
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

func CreateUserInfoByUserId(db *gorm.DB, userId uuid.UUID, input CreateUserInfoInput) (*UserInfo, *exceptions.Exception) {
	if db == nil {
		db = NotezyDB
	}

	if err := Validator.Struct(input); err != nil {
		return nil, exceptions.UserInfo.InvalidInput().WithError(err)
	}

	newUserInfo := UserInfo{
		UserId:             userId,
		CoverBackgroundURL: input.CoverBackgroundURL,
		AvatarURL:          input.AvatarURL,
		Header:             input.Header,
		Introduction:       input.Introduction,
		Gender:             UserGender_PreferNotToSay,
		Country:            Country_UnitedStatusOfAmerica,
		BirthDate:          input.BirthDate,
	}
	if input.Gender != nil {
		newUserInfo.Gender = *input.Gender
	}
	if input.Country != nil {
		newUserInfo.Country = *input.Country
	}

	result := db.Table(UserInfo{}.TableName()).Create(&newUserInfo)
	if err := result.Error; err != nil {
		return nil, exceptions.UserInfo.FailedToCreate().WithError(err)
	}
	return &newUserInfo, nil
}

func UpdateUserInfoByUserId(db *gorm.DB, userId uuid.UUID, input UpdateUserInfoInput) (*UserInfo, *exceptions.Exception) {
	if db == nil {
		db = NotezyDB
	}

	if err := Validator.Struct(input); err != nil {
		return nil, exceptions.UserInfo.InvalidInput().WithError(err)
	}

	updatedUserInfo := UserInfo{
		UserId:             userId,
		CoverBackgroundURL: input.CoverBackgroundURL,
		AvatarURL:          input.AvatarURL,
		Header:             input.Header,
		Introduction:       input.Introduction,
		BirthDate:          input.BirthDate,
	}
	if input.Gender != nil {
		updatedUserInfo.Gender = *input.Gender
	}
	if input.Country != nil {
		updatedUserInfo.Country = *input.Country
	}

	result := db.Table(UserInfo{}.TableName()).Updates(&updatedUserInfo)
	if err := result.Error; err != nil {
		return nil, exceptions.UserInfo.FailedToUpdate().WithError(err)
	}
	return &updatedUserInfo, nil
}
