package schemas

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"

	enums "notezy-backend/app/models/schemas/enums"
	types "notezy-backend/shared/types"
)

type UserAccount struct {
	Id                  uuid.UUID          `json:"id" gorm:"column:id; type:uuid; primaryKey; not null; default:gen_random_uuid();"`
	UserId              uuid.UUID          `json:"userId" gorm:"column:user_id; type:uuid; not null; unique;"`
	AuthCode            string             `json:"authCode" gorm:"column:auth_code; not null;"`                     // validate:"required,isnumberstring,len=6"
	AuthCodeExpiredAt   time.Time          `json:"authCodeExpiredAt" gorm:"column:auth_code_expired_at; not null;"` // the exact time when authCode expires
	BlockAuthCodeUntil  time.Time          `json:"blockAuthCodeUntil" gorm:"column:block_auth_code_until; type:timestamptz; not null;"`
	CountryCode         *enums.CountryCode `json:"countryCode" gorm:"column:country_code; type:\"CountryCode\";"` // validate:"omitnil,iscountrycode"
	PhoneNumber         *string            `json:"phoneNumber" gorm:"column:phone_number; unique;"`               // validate:"omitnil,max=0,max=15,isnumberstring"
	GoogleCredential    *string            `json:"googleCredential" gorm:"column:google_credential; unique;"`     // validate:"omitnil"
	DiscordCredential   *string            `json:"discordCredential" gorm:"column:discord_credential; unique;"`   // validate:"omitnil"
	RootShelfCount      int32              `json:"rootShelfCount" gorm:"column:root_shelf_count; type:integer; not null; default:0;"`
	BlockPackCount      int32              `json:"blockPackCount" gorm:"column:block_pack_count; type:integer; not null; default:0;"`
	BlockCount          int32              `json:"blockCount" gorm:"column:block_count; type:integer; not null; default:0;"`
	MaterialCount       int32              `json:"materialCount" gorm:"column:material_count; type:integer; not null; default:0;"`
	WorkflowCount       int32              `json:"workflowCount" gorm:"column:workflow_count; type:integer; not null; default:0;"`
	AdditionalItemCount int32              `json:"additionalItemCount" gorm:"column:additional_item_count; type:integer; not null; default:0;"`
	UpdatedAt           time.Time          `json:"updatedAt" gorm:"column:updated_at; type:timestamptz; not null; autoUpdateTime:true;"`
}

// User Account Table Name
func (UserAccount) TableName() string {
	return types.TableName_UserAccountTable.String()
}

/* ============================== Relative Type Conversions ============================== */
// note that there's no type like PublicUserAccount,
// since the userAccount shouldn't be public

/* ============================== Trigger Hook ============================== */

func (ua *UserAccount) BeforeCreate(db *gorm.DB) error {
	if ua.BlockAuthCodeUntil.IsZero() {
		ua.BlockAuthCodeUntil = time.Now().Add(-10 * time.Minute)
	}
	return nil
}

func (ua *UserAccount) BeforeUpdate(db *gorm.DB) error {
	if ua.AuthCode != "" && ua.BlockAuthCodeUntil.After(time.Now()) {
		return fmt.Errorf("cannot send auth code until: %v", ua.BlockAuthCodeUntil)
	}
	return nil
}
