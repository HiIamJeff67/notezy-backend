package schemas

import (
	"time"

	"github.com/google/uuid"

	shared "notezy-backend/shared"
)

type UsersToBadges struct {
	UserId    uuid.UUID `json:"userId" gorm:"column:user_id; type:uuid; primaryKey;"`
	BadgeId   uuid.UUID `json:"badgeId" gorm:"column:badge_id; type:uuid; primaryKey;"`
	CreatedAt time.Time `json:"createdAt" gorm:"column:created_at; type:timestamptz; not null; autoCreateTime:true;"`

	// relations
	User  User  `gorm:"foreignKey:UserId; references:Id; constraint:OnUpdate:CASCADE, OnDelete:CASCADE;"`
	Badge Badge `gorm:"foreignKey:BadgeId; references:Id; constraint:OnUpdate:CASCADE, OnDelete:CASCADE;"`
}

// Users To Badges Table Name
func (UsersToBadges) TableName() string {
	return shared.ValidTableName_UsersToBadgesTable.String()
}

// Users To Badges Table Relations
type UsersToBadgesRelation string

const (
	UsersToBadgesRelation_User                        = "User"
	UsersToBadgesRelation_Badge UsersToBadgesRelation = "Badge"
)
