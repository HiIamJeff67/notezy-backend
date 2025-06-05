package inputs

import (
	"notezy-backend/app/models/enums"
	"time"
)

type CreateUserAccountInput struct {
	AuthCode          string             `json:"authCode" validate:"required" gorm:"column:auth_code"`
	AuthCodeExpiredAt time.Time          `json:"authCodeExpiredAt" validate:"required" gorm:"column:auth_code_expired_at"`
	CountryCode       *enums.CountryCode `json:"countryCode" validate:"omitempty,iscountrycode" gorm:"column:country_code;"`
	PhoneNumber       *string            `json:"phoneNumber" validate:"omitempty,max=0,max=15" gorm:"column:phone_number;"`
	GoogleCredential  *string            `json:"googleCrendential" validate:"omitempty" gorm:"column:google_credential;"`
	DiscordCredential *string            `json:"discordCrendential" validate:"omitempty" gorm:"column:discord_credential;"`
}

type UpdateUserAccountInput struct {
	AuthCode          *string            `json:"authCode" validate:"omitempty" gorm:"column:auth_code"`
	CountryCode       *enums.CountryCode `json:"countryCode" validate:"omitempty,iscountrycode" gorm:"column:country_code;"`
	PhoneNumber       *string            `json:"phoneNumber" validate:"omitempty,max=0,max=15" gorm:"column:phone_number;"`
	GoogleCredential  *string            `json:"googleCrendential" validate:"omitempty" gorm:"column:google_credential;"`
	DiscordCredential *string            `json:"discordCrendential" validate:"omitempty" gorm:"column:discord_credential;"`
}
