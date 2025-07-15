package schemas

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"

	gqlmodels "notezy-backend/app/graphql/models"
	enums "notezy-backend/app/models/schemas/enums"
	util "notezy-backend/app/util"
	shared "notezy-backend/shared"
)

type User struct {
	Id             uuid.UUID        `json:"id" gorm:"column:id; type:uuid; primaryKey; default:gen_random_uuid();"`
	PublicId       string           `json:"publicId" gorm:"column:public_id; unique; not null; default:'';"`
	Name           string           `json:"name" gorm:"column:name; unique; not null; size:16;"`
	DisplayName    string           `json:"displayName" gorm:"column:display_name; not null; size:32;"`
	Email          string           `json:"email" gorm:"column:email; unique; not null;"`
	Password       string           `json:"password" gorm:"column:password; not null; size:1024;"` // since we store the hashed password which is quite long
	RefreshToken   string           `json:"refreshToken" gorm:"column:refresh_token; not null;"`
	LoginCount     int32            `json:"loginCount" gorm:"column:login_count; type:integer not null; default:0;"`
	BlockLoginUtil time.Time        `json:"blockLoginUntil" gorm:"column:block_login_until; type:timestamptz; not null;"`
	UserAgent      string           `json:"userAgent" gorm:"column:user_agent; not null;"`
	Role           enums.UserRole   `json:"role" gorm:"column:role; type:UserRole; not null; default:'Guest';"`
	Plan           enums.UserPlan   `json:"plan" gorm:"column:plan; type:UserPlan; not null; default:'Free';"`
	PrevStatus     enums.UserStatus `json:"prevStatus" gorm:"column:prev_status; type:UserStatus; not null; default:'Online';"`
	Status         enums.UserStatus `json:"status" gorm:"column:status; type:UserStatus; not null; default:'Online';"`
	UpdatedAt      time.Time        `json:"updatedAt" gorm:"column:updated_at; type:timestamptz; not null; autoUpdateTime:true;"`
	CreatedAt      time.Time        `json:"createdAt" gorm:"column:created_at; type:timestamptz; not null; autoCreateTime:true;"`

	// relations
	UserInfo    UserInfo    `json:"userInfo" gorm:"foreignKey:UserId; references:Id; constraint:OnUpdate:CASCADE, OnDelete:CASCADE;"`
	UserAccount UserAccount `json:"userAccount" gorm:"foreignKey:UserId; references:Id; constraint:OnUpdate:CASCADE, OnDelete:CASCADE;"`
	UserSetting UserSetting `json:"userSetting" gorm:"foreignKey:UserId; references:Id; constraint:OnUpdate:CASCADE, OnDelete:CASCADE;"`
	Badges      []Badge     `json:"badges" gorm:"many2many:\"UsersToBadgesTable\"; foreignKey:Id; joinForeignKey:UserId; references:Id; joinReferences:BadgeId;"`
	Themes      []Theme     `json:"themes" gorm:"foreignKey:AuthorId;"`
}

// force gorm to use the given table name
func (User) TableName() string {
	return shared.ValidTableName_UserTable.String()
}

/* ============================== Relative Type Conversions ============================== */

func (u *User) ToPublicUser() *gqlmodels.PublicUser {
	return &gqlmodels.PublicUser{
		PublicID:    u.PublicId,
		Name:        u.Name,
		DisplayName: u.DisplayName,
		Email:       u.Email,
		Role:        u.Role,
		Plan:        u.Plan,
		Status:      u.Status,
		CreatedAt:   u.CreatedAt,
		UpdatedAt:   u.UpdatedAt,
		UserInfo:    &gqlmodels.PublicUserInfo{},
		Badges:      []*gqlmodels.PublicBadge{},
		Themes:      []*gqlmodels.PublicTheme{},
	}
}

/* ============================== Trigger Hook ============================== */

func (u *User) BeforeCreate(db *gorm.DB) error {
	if u.BlockLoginUtil.IsZero() {
		// just to make the new user can login
		u.BlockLoginUtil = time.Now().Add(-10 * time.Minute)
	}
	if u.PublicId == "" {
		u.PublicId = util.GenerateSnowflakeID()
	}
	return nil
}

func (u *User) BeforeUpdate(db *gorm.DB) error {
	if db.Statement.Changed("Status") {
		u.PrevStatus = u.Status
	}
	return nil
}
