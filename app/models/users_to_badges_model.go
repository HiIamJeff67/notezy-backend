package models

import (
	"notezy-backend/global"
	"time"

	uuid "github.com/jackc/pgx/pgtype/ext/satori-uuid"
)

/* ============================== Schema ============================== */
type UsersToBadges struct {
	UserId    uuid.UUID `json:"userId" gorm:"column:user_id ;primaryKey;"`
	BadgeId   uuid.UUID `json:"badgeId" gorm:"column:badge_id; primaryKey;"`
	CreatedAt time.Time `json:"createdAt" gorm:"column:created_at; type:timestamptz; not null; autoCreateTime:true;"`

	User  User  `gorm:"foreignKey:UserID; references:ID; constraint:OnUpdate:CASCADE, OnDelete:CASCADE;"`
	Badge Badge `gorm:"foreignKey:BadgeID; references:ID; constraint:OnUpdate:CASCADE, OnDelete:CASCADE;"`
}

func (UsersToBadges) TableName() string {
	return string(global.ValidTableName_UsersToBadgesTable)
}

/* ============================== Schema ============================== */

/* ============================== Input ============================== */

/* ============================== Methods ============================== */
