package inputs

import (
	"time"

	enums "notezy-backend/app/models/schemas/enums"
)

type CreateUserAccountInput struct {
	AuthCode          string             `json:"authCode" validate:"required" gorm:"column:auth_code"`
	AuthCodeExpiredAt time.Time          `json:"authCodeExpiredAt" validate:"required" gorm:"column:auth_code_expired_at"`
	CountryCode       *enums.CountryCode `json:"countryCode" validate:"omitempty,iscountrycode" gorm:"column:country_code;"`
	PhoneNumber       *string            `json:"phoneNumber" validate:"omitempty,max=0,max=15" gorm:"column:phone_number;"`
	GoogleCredential  *string            `json:"googleCredential" validate:"omitempty" gorm:"column:google_credential;"`
	DiscordCredential *string            `json:"discordCredential" validate:"omitempty" gorm:"column:discord_credential;"`
}

type UpdateUserAccountInput struct {
	AuthCode           *string            `json:"authCode" validate:"omitempty" gorm:"column:auth_code"`
	AuthCodeExpiredAt  *time.Time         `json:"authCodeExpiredAt" validate:"omitempty" gorm:"column:auth_code_expired_at"`
	BlockAuthCodeUntil *time.Time         `json:"blockAuthCodeUntil" validate:"omitempty" gorm:"column:block_auth_code_until"`
	CountryCode        *enums.CountryCode `json:"countryCode" validate:"omitempty,iscountrycode" gorm:"column:country_code;"`
	PhoneNumber        *string            `json:"phoneNumber" validate:"omitempty,max=0,max=15" gorm:"column:phone_number;"`
	GoogleCredential   *string            `json:"googleCredential" validate:"omitempty" gorm:"column:google_credential;"`
	DiscordCredential  *string            `json:"discordCredential" validate:"omitempty" gorm:"column:discord_credential;"`
}

type PartialUpdateUserAccountInput = PartialUpdateInput[UpdateUserAccountInput]
