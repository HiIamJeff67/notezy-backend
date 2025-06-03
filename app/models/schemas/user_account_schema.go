package schemas

import (
	"notezy-backend/app/models/enums"
	"notezy-backend/global"
	"time"

	"github.com/google/uuid"
)

type UserAccount struct {
	Id                uuid.UUID         `json:"id" gorm:"column:id; type:uuid; primaryKey; not null; default:gen_random_uuid();"`
	UserId            uuid.UUID         `json:"userId" gorm:"column:user_id; type:uuid; not null; unique;"`
	CountryCode       enums.CountryCode `json:"countryCound" gorm:"column:country_code; type:CountryCode; not null; default:'Default'"`
	PhoneNumber       string            `json:"phoneNumber" gorm:"column:phone_number; unique; not null; default:''"`
	GoogleCredential  string            `json:"googleCrendential" gorm:"column:google_credential; unique; not null; default:''"`
	DiscordCredential string            `json:"discordCrendential" gorm:"column:discord_credential; unique; not null; default:''"`
	UpdatedAt         time.Time         `json:"updatedAt" gorm:"column:updated_at; type:timestamptz; not null; autoUpdateTime:true;"`
}

func (UserAccount) TableName() string {
	return string(global.ValidTableName_UserAccountTable)
}
