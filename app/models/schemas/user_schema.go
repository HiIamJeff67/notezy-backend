package schemas

import (
	"notezy-backend/app/models/enums"
	"notezy-backend/global"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type User struct {
	Id           uuid.UUID        `json:"id" gorm:"column:id; type:uuid; primaryKey; default:gen_random_uuid();"`
	Name         string           `json:"name" gorm:"column:name; unique; not null; size:16;"`
	DisplayName  string           `json:"displayName" gorm:"column:display_name; not null; size:32;"`
	Email        string           `json:"email" gorm:"column:email; unique; not null;"`
	Password     string           `json:"password" gorm:"column:password; not null; size:1024;"` // since we store the hashed password which is quite long
	RefreshToken string           `json:"refreshToken" gorm:"column:refresh_token; not null; default:'';"`
	UserAgent    string           `json:"userAgent" gorm:"column:user_agent; not null;"`
	Role         enums.UserRole   `json:"role" gorm:"column:role; type:UserRole; not null; default:'Guest';"`
	Plan         enums.UserPlan   `json:"plan" gorm:"column:plan; type:UserPlan; not null; default:'Free';"`
	PrevStatus   enums.UserStatus `json:"prevStatus" gorm:"column:prev_status; type:UserStatus; not null; default:'Online';"`
	Status       enums.UserStatus `json:"status" gorm:"column:status; type:UserStatus; not null; default:'Online';"`
	UpdatedAt    time.Time        `json:"updatedAt" gorm:"column:updated_at; type:timestamptz; not null; autoUpdateTime:true;"`
	CreatedAt    time.Time        `json:"createdAt" gorm:"column:created_at; type:timestamptz; not null; autoCreateTime:true;"`

	// relation
	UserInfo    UserInfo    `json:"userInfo" gorm:"foreignKey:UserId; references:Id; constraint:OnUpdate:CASCADE, OnDelete:CASCADE;"`
	UserAccount UserAccount `json:"userAccount" gorm:"foreignKey:UserId; references:Id; constraint:OnUpdate:CASCADE, OnDelete:CASCADE;"`
	UserSetting UserSetting `json:"userSetting" gorm:"foreignKey:UserId; references:Id; constraint:OnUpdate:CASCADE, OnDelete:CASCADE;"`
	Badges      []Badge     `json:"badges" gorm:"-"` // many2many:\"UsersToBadgesTable\"; foreignKey:Id; joinForeignKey:UserId; references:Id; joinReferences:BadgeId;
}

// force gorm to use the given table name
func (User) TableName() string {
	return global.ValidTableName_UserTable.String()
}

//* ============================== Trigger Hook ============================== */

func (u *User) BeforeUpdate(db *gorm.DB) (err error) {
	if db.Statement.Changed("Status") {
		u.PrevStatus = u.Status
	}
	return nil
}
