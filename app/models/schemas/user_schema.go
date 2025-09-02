package schemas

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"

	gqlmodels "notezy-backend/app/graphql/models"
	enums "notezy-backend/app/models/schemas/enums"
	types "notezy-backend/shared/types"
)

type User struct {
	Id              uuid.UUID        `json:"id" gorm:"column:id; type:uuid; primaryKey; default:gen_random_uuid();"`
	PublicId        string           `json:"publicId" gorm:"column:public_id; unique; not null; default:'';"`
	Name            string           `json:"name" gorm:"column:name; unique; not null; size:16;"`        // validate:"required,min=6,max=16,alphaandnum"
	DisplayName     string           `json:"displayName" gorm:"column:display_name; not null; size:32;"` // validate:"required,min=6,max=32,alphaandnum"
	Email           string           `json:"email" gorm:"column:email; unique; not null;"`               // validate:"required,email"
	Password        string           `json:"password" gorm:"column:password; not null; size:1024;"`      // validate:"required,min=8,max=1024"      // since we store the hashed password which is quite long
	RefreshToken    string           `json:"refreshToken" gorm:"column:refresh_token; not null;"`        // validate:"omitnil"
	LoginCount      int32            `json:"loginCount" gorm:"column:login_count; type:integer; not null; default:0;"`
	BlockLoginUntil time.Time        `json:"blockLoginUntil" gorm:"column:block_login_until; type:timestamptz; not null;"`
	UserAgent       string           `json:"userAgent" gorm:"column:user_agent; not null;"`                      // validate:"required,isuseragent"
	Role            enums.UserRole   `json:"role" gorm:"column:role; type:UserRole; not null; default:'Guest';"` // validate:"omitnil,isrole"
	Plan            enums.UserPlan   `json:"plan" gorm:"column:plan; type:UserPlan; not null; default:'Free';"`  // validate:"omitnil,isplan"
	PrevStatus      enums.UserStatus `json:"prevStatus" gorm:"column:prev_status; type:UserStatus; not null; default:'Online';"`
	Status          enums.UserStatus `json:"status" gorm:"column:status; type:UserStatus; not null; default:'Online';"` // validate:"omitnil,isstatus"
	UpdatedAt       time.Time        `json:"updatedAt" gorm:"column:updated_at; type:timestamptz; not null; autoUpdateTime:true;"`
	CreatedAt       time.Time        `json:"createdAt" gorm:"column:created_at; type:timestamptz; not null; autoCreateTime:true;"`

	// relations
	UserInfo       UserInfo         `json:"userInfo" gorm:"foreignKey:UserId; references:Id; constraint:OnUpdate:CASCADE, OnDelete:CASCADE;"`
	UserAccount    UserAccount      `json:"userAccount" gorm:"foreignKey:UserId; references:Id; constraint:OnUpdate:CASCADE, OnDelete:CASCADE;"`
	UserSetting    UserSetting      `json:"userSetting" gorm:"foreignKey:UserId; references:Id; constraint:OnUpdate:CASCADE, OnDelete:CASCADE;"`
	Themes         []Theme          `json:"themes" gorm:"foreignKey:AuthorId; references:Id; constraint:OnUpdate:CASCADE, OnDelete:CASCADE;"`
	UsersToBadges  []UsersToBadges  `json:"usersToBadges" gorm:"foreignKey:UserId; references:Id; constraint:OnUpdate:CASCADE, OnDelete:CASCADE;"`
	UsersToShelves []UsersToShelves `json:"usersToShelves" gorm:"foreignKey:UserId; references:Id; constraint:OnUpdate:CASCADE, OnDelete:CASCADE;"`
}

// User Table Name
func (User) TableName() string {
	return types.ValidTableName_UserTable.String()
}

// User Table Relations
type UserRelation types.ValidTableName

const (
	UserRelation_UserInfo       UserRelation = "UserInfo"
	UserRelation_UserAccount    UserRelation = "UserAccount"
	UserRelation_UserSetting    UserRelation = "UserSetting"
	UserRelation_Themes         UserRelation = "Themes"
	UserRelation_UsersToBadges  UserRelation = "UsersToBadges"
	UserRelation_UsersToShelves UserRelation = "UsersToShelves"
)

/* ============================== Relative Type Conversions ============================== */

func (u *User) ToPublicUser() *gqlmodels.PublicUser {
	return &gqlmodels.PublicUser{
		PublicID:    u.PublicId,
		Name:        u.Name,
		DisplayName: u.DisplayName,
		Role:        u.Role,
		Plan:        u.Plan,
		Status:      u.Status,
		CreatedAt:   u.CreatedAt,
		UserInfo:    &gqlmodels.PublicUserInfo{},
		Badges:      make([]*gqlmodels.PublicBadge, 0),
		Themes:      make([]*gqlmodels.PublicTheme, 0),
	}
}

/* ============================== Trigger Hook ============================== */

func (u *User) BeforeCreate(db *gorm.DB) error {
	if u.BlockLoginUntil.IsZero() {
		// just to make the new user can login
		u.BlockLoginUntil = time.Now().Add(-10 * time.Minute)
	}
	if u.PublicId == "" {
		u.PublicId = uuid.NewString()
	}
	return nil
}

func (u *User) BeforeUpdate(db *gorm.DB) error {
	if db.Statement.Changed("Status") {
		u.PrevStatus = u.Status
	}
	return nil
}
