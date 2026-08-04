package schemas

import (
	"time"

	"github.com/google/uuid"

	enums "github.com/HiIamJeff67/notezy-backend/internal/services/durablejob/data/schemas/enums"
	types "github.com/HiIamJeff67/notezy-backend/shared/types"
)

type UserInfo struct {
	Id                 uuid.UUID        `json:"id" gorm:"column:id; type:uuid; primaryKey; default:gen_random_uuid();"`
	UserId             uuid.UUID        `json:"userId" gorm:"column:user_id; type:uuid; not null; unique;"`
	CoverBackgroundURL *string          `json:"coverBackgroundURL" gorm:"column:cover_background_url; default:null;"`                 // validate:"omitnil,isimageurl"
	AvatarURL          *string          `json:"avatarURL" gorm:"column:avatar_url; default:null;"`                                    // validate:"omitnil,isimageurl"
	Header             *string          `json:"header" gorm:"column:header; size:64;"`                                                // validate:"omitnil,min=0,max=64"
	Introduction       *string          `json:"introduction" gorm:"column:introduction; size:256;"`                                   // validate:"omitnil,min=0,max=256"
	Gender             enums.UserGender `json:"gender" gorm:"column:gender; type:\"UserGender\"; not null; default:'PreferNotToSay'"` // validate:"omitnil,isgender"
	Country            *enums.Country   `json:"country" gorm:"column:country; type:\"Country\";"`                                     // validate:"omitnil,iscountry"
	BirthDate          time.Time        `json:"birthDate" gorm:"column:birth_date; type:timestamptz; not null; default:CURRENT_TIMESTAMP"`
	UpdatedAt          time.Time        `json:"updatedAt" gorm:"column:updated_at; type:timestamptz; not null; autoUpdateTime:true;"`
}

// User Info Table Name
func (UserInfo) TableName() string {
	return types.TableName_UserInfoTable.String()
}
