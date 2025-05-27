package models

import (
	"go-gorm-api/global"
	"time"

	uuid "github.com/jackc/pgx/pgtype/ext/satori-uuid"
)

/* ============================== Schema ============================== */
type UserInfo struct {
	Id					uuid.UUID		`json:"id" gorm:"column:id; type:uuid; primaryKey; default:gen_random_uuid();"`
	UserId				uuid.UUID       `json:"userId" gorm:"column:user_id; type:uuid; not null; unique;"`
	CoverBackgroundURL  *string			`json:"coverBackgroundURL" gorm:"column:cover_background_url;"`
	AvatarURL			*string			`json:"avatarURL" gorm:"column:avatar_url;"`
	Header				*string			`json:"header" gorm:"column:header; not null; default:''; size:64;"`
	Introduction		*string			`json:"introduction" gorm:"column:introduction; not null; default:''; size:256;"`
	Gender				UserGender		`json:"gender" gorm:"column:gender; type:UserGender; not null; default:'PreferNotToSay'"`
	Country				Country			`json:"country" gorm:"column:country; type:Country; not null; default:'UnitedStatusOfAmerica'"`
	BirthDate			*time.Time		`json:"birthDate" gorm:"column:birth_date; type:timestamptz;"`
	UpdatedAt			time.Time		`json:"updatedAt" gorm:"column:updated_at; type:timestamptz; not null; autoUpdateTime:true;"`
}

func (UserInfo) TableName() string {
	return string(global.ValidTableName_UserInfoTable)
}
/* ============================== Schema ============================== */

/* ============================== Input ============================== */
type CreateUserInfoInput struct {
	CoverBackgroundURL  *string			`json:"coverBackgroundURL" gorm:"column:cover_background_url;"`
	AvatarURL			*string			`json:"avatarURL" gorm:"column:avatar_url;"`
	Header				*string			`json:"header" gorm:"column:header;"`
	Introduction		*string			`json:"introduction" gorm:"column:introduction;"`
	Gender				UserGender		`json:"gender" gorm:"column:gender;"`
	Country				Country			`json:"country" gorm:"column:country;"`
	BirthDate			*time.Time		`json:"birthDate" gorm:"column:birth_date;"`
}

type UpdateUserInfoInput struct {
	CoverBackgroundURL  *string			`json:"coverBackgroundURL" gorm:"column:cover_background_url;"`
	AvatarURL			*string			`json:"avatarURL" gorm:"column:avatar_url;"`
	Header				*string			`json:"header" gorm:"column:header;"`
	Introduction		*string			`json:"introduction" gorm:"column:introduction;"`
	Gender				UserGender		`json:"gender" gorm:"column:gender;"`
	Country				Country			`json:"country" gorm:"column:country;"`
	BirthDate			*time.Time		`json:"birthDate" gorm:"column:birth_date;"`
}
/* ============================== Input ============================== */

/* ============================== Methods ============================== */
func GetUserInfoByUserId(userId uuid.UUID) (userInfo UserInfo, err error) {
	result := NotezyDB.Table(UserInfo{}.TableName()).Where("user_id = ?", userId).First(&userInfo)
	if err = result.Error; err != nil {
		return UserInfo{}, err
	}
	return userInfo, nil
}

func CreateUserInfoByUserId(userId uuid.UUID, input CreateUserInfoInput) (newUserInfo UserInfo, err error) {
	newUserInfo = UserInfo{
		UserId: userId, 
		CoverBackgroundURL: input.CoverBackgroundURL,
		AvatarURL: input.AvatarURL,
		Header: input.Header,
		Introduction: input.Introduction,
		Gender: input.Gender, 
		Country: input.Country,
		BirthDate: input.BirthDate,
	}

	result := NotezyDB.Table(UserInfo{}.TableName()).Create(&newUserInfo)
	if err = result.Error; err != nil {
		return UserInfo{}, err
	}
	return newUserInfo, nil
}

func UpdateUserInfoByUserId(userId uuid.UUID, input UpdateUserInfoInput) (updatedUserInfo UserInfo, err error) {
	updatedUserInfo = UserInfo{
		UserId: userId, 
		CoverBackgroundURL: input.CoverBackgroundURL,
		AvatarURL: input.AvatarURL,
		Header: input.Header,
		Introduction: input.Introduction,
		Gender: input.Gender, 
		Country: input.Country,
		BirthDate: input.BirthDate,
	}

	result := NotezyDB.Table(UserInfo{}.TableName()).Updates(&updatedUserInfo)
	if err = result.Error; err != nil {
		return UserInfo{}, err
	}
	return updatedUserInfo, nil
}
/* ============================== Methods ============================== */