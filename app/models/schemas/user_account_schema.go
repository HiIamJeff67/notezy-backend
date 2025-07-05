package schemas

import (
	"time"

	"github.com/google/uuid"

	enums "notezy-backend/app/models/schemas/enums"
	shared "notezy-backend/shared"
)

type UserAccount struct {
	Id                uuid.UUID         `json:"id" gorm:"column:id; type:uuid; primaryKey; not null; default:gen_random_uuid();"`
	UserId            uuid.UUID         `json:"userId" gorm:"column:user_id; type:uuid; not null; unique;"`
	AuthCode          string            `json:"authCode" gorm:"column:auth_code; not null;"`
	AuthCodeExpiredAt time.Time         `json:"authCodeExpiredAt" gorm:"column:auth_code_expired_at; not null;"` // the exact time when authCode expires
	CountryCode       enums.CountryCode `json:"countryCode" gorm:"column:country_code; type:CountryCode;"`
	PhoneNumber       *string           `json:"phoneNumber" gorm:"column:phone_number; unique;"`
	GoogleCredential  *string           `json:"googleCredential" gorm:"column:google_credential; unique;"`
	DiscordCredential *string           `json:"discordCredential" gorm:"column:discord_credential; unique;"`
	UpdatedAt         time.Time         `json:"updatedAt" gorm:"column:updated_at; type:timestamptz; not null; autoUpdateTime:true;"`
}

func (UserAccount) TableName() string {
	return shared.ValidTableName_UserAccountTable.String()
}
