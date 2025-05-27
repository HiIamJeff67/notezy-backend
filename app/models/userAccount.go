package models

import (
	"go-gorm-api/global"
	"time"

	uuid "github.com/jackc/pgx/pgtype/ext/satori-uuid"
)

/* ============================== Schema ============================== */
type UserAccount struct {
	Id					uuid.UUID		`json:"id" gorm:"column:id; type:uuid; primaryKey; not null; default:gen_random_uuid();"`
	UserId				uuid.UUID       `json:"userId" gorm:"column:user_id; type:uuid; not null; unique;"`
	CountryCode			*CountryCode	`json:"countryCound" gorm:"column:country_code; type:CountryCode; not null;"`
	PhoneNumber			*string			`json:"phoneNumber" gorm:"column:phone_number; unique;"`
	GoogleCredential	*string			`json:"googleCrendential" gorm:"column:google_credential; unique;"`
	DiscordCredential	*string			`json:"discordCrendential" gorm:"column:discord_credential; unique;"`
	UpdatedAt			time.Time		`json:"updatedAt" gorm:"column:updated_at; type:timestamptz; not null; autoUpdateTime:true;"`
}

func (UserAccount) TableName() string {
	return string(global.ValidTableName_UserAccountTable)
}
/* ============================== Schema ============================== */

/* ============================== Input ============================== */
type CreateUserAccountInput struct {
	CountryCode			*CountryCode	`json:"countryCode" gorm:"column:country_code;"`
	PhoneNumber			*string			`json:"phoneNumber" gorm:"column:phone_number;"`
	GoogleCredential	*string			`json:"googleCrendential" gorm:"column:google_credential;"`
	DiscordCredential	*string			`json:"discordCrendential" gorm:"column:discord_credential;"`
}

type UpdateUserAccountInput struct {
	CountryCode			*CountryCode	`json:"countryCode" gorm:"column:country_code;"`
	PhoneNumber			*string			`json:"phoneNumber" gorm:"column:phone_number;"`
	GoogleCredential	*string			`json:"googleCrendential" gorm:"column:google_credential;"`
	DiscordCredential	*string			`json:"discordCrendential" gorm:"column:discord_credential;"`
}
/* ============================== Input ============================== */

/* ============================== Methods ============================== */
func GetUserAccountByUserId(userId uuid.UUID) (userAccount UserAccount, err error) {
	result := NotezyDB.Table(UserAccount{}.TableName()).Where("user_id = ?", userId).First(&userAccount)
	if err = result.Error; err != nil {
		return UserAccount{}, err
	}
	return userAccount, nil
}

func GetAllUserAccount() (userAccounts []UserAccount, err error) {
	result := NotezyDB.Table(UserAccount{}.TableName()).Find(&userAccounts)
	if err = result.Error; err != nil {
		return []UserAccount{}, err
	}
	return userAccounts, nil
}

func CreateUserAccountByUserId(userId uuid.UUID, input CreateUserAccountInput) (newUserAccount UserAccount, err error) {
	newUserAccount = UserAccount{
		UserId: userId, 
		CountryCode: input.CountryCode, 
		PhoneNumber: input.PhoneNumber,
		GoogleCredential: input.GoogleCredential,
		DiscordCredential: input.DiscordCredential,
	}

	result := NotezyDB.Table(UserAccount{}.TableName()).Create(&newUserAccount)
	if err = result.Error; err != nil {
		return UserAccount{}, err
	}
	return newUserAccount, nil
}

func UpdateUserAccountByUserId(userId uuid.UUID, input UpdateUserAccountInput) (updatedUserAccount UserAccount, err error) {
	updatedUserAccount = UserAccount{
		CountryCode: input.CountryCode,
		PhoneNumber: input.PhoneNumber,
		GoogleCredential: input.GoogleCredential,
		DiscordCredential: input.DiscordCredential,
	}
	
	result := NotezyDB.Table(UserAccount{}.TableName()).Where("user_id = ?", userId).Updates(&input)
	if err := result.Error; err != nil {
		return UserAccount{}, err
	}
	return updatedUserAccount, nil
}

// We do not allow to just delete the userAccount, 
// instead, the userAccount is only deleted by deleting the user
// func DeleteUserAccount(userId uuid.UUID) (deletedUserAccount User, err error) {}
/* ============================== Methods ============================== */