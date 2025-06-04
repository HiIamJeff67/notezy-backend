package schemas

import (
	"notezy-backend/app/models/enums"
	"notezy-backend/global"
	"time"

	"github.com/google/uuid"
)

type UserInfo struct {
	Id                 uuid.UUID        `json:"id" gorm:"column:id; type:uuid; primaryKey; default:gen_random_uuid();"`
	UserId             uuid.UUID        `json:"userId" gorm:"column:user_id; type:uuid; not null; unique;"`
	CoverBackgroundURL string           `json:"coverBackgroundURL" gorm:"column:cover_background_url; not null; default:''"`
	AvatarURL          string           `json:"avatarURL" gorm:"column:avatar_url; not null; default:''"`
	Header             string           `json:"header" gorm:"column:header; not null; default:''; size:64;"`
	Introduction       string           `json:"introduction" gorm:"column:introduction; not null; default:''; size:256;"`
	Gender             enums.UserGender `json:"gender" gorm:"column:gender; type:UserGender; not null; default:'PreferNotToSay'"`
	Country            enums.Country    `json:"country" gorm:"column:country; type:Country; not null; default:'Default'"`
	BirthDate          time.Time        `json:"birthDate" gorm:"column:birth_date; type:timestamptz; not null; default:CURRENT_TIMESTAMP"`
	UpdatedAt          time.Time        `json:"updatedAt" gorm:"column:updated_at; type:timestamptz; not null; autoUpdateTime:true;"`
}

func (UserInfo) TableName() string {
	return global.ValidTableName_UserInfoTable.String()
}
