package schemas

import (
	"time"

	"github.com/google/uuid"

	gqlmodels "notezy-backend/app/graphql/models"
	enums "notezy-backend/app/models/schemas/enums"
	shared "notezy-backend/shared"
)

type UserInfo struct {
	Id                 uuid.UUID        `json:"id" gorm:"column:id; type:uuid; primaryKey; default:gen_random_uuid();"`
	UserId             uuid.UUID        `json:"userId" gorm:"column:user_id; type:uuid; not null; unique;"`
	CoverBackgroundURL *string          `json:"coverBackgroundURL" gorm:"column:cover_background_url;"`
	AvatarURL          *string          `json:"avatarURL" gorm:"column:avatar_url;"`
	Header             *string          `json:"header" gorm:"column:header; size:64;"`
	Introduction       *string          `json:"introduction" gorm:"column:introduction; size:256;"`
	Gender             enums.UserGender `json:"gender" gorm:"column:gender; type:UserGender; not null; default:'PreferNotToSay'"`
	Country            *enums.Country   `json:"country" gorm:"column:country; type:Country;"`
	BirthDate          time.Time        `json:"birthDate" gorm:"column:birth_date; type:timestamptz; not null; default:CURRENT_TIMESTAMP"`
	UpdatedAt          time.Time        `json:"updatedAt" gorm:"column:updated_at; type:timestamptz; not null; autoUpdateTime:true;"`
}

func (UserInfo) TableName() string {
	return shared.ValidTableName_UserInfoTable.String()
}

/* ============================== Relative Type Conversions ============================== */

func (ui *UserInfo) ToPublicUserInfo() *gqlmodels.PublicUserInfo {
	return &gqlmodels.PublicUserInfo{
		CoverBackgroundURL: ui.CoverBackgroundURL,
		AvatarURL:          ui.AvatarURL,
		Header:             ui.Header,
		Introduction:       ui.Introduction,
		Gender:             ui.Gender,
		Country:            ui.Country,
		BirthDate:          ui.BirthDate,
	}
}
