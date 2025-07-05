package schemas

import (
	"time"

	"github.com/google/uuid"

	enums "notezy-backend/app/models/schemas/enums"
	shared "notezy-backend/shared"
)

type Badge struct {
	Id          uuid.UUID       `json:"id" gorm:"column:id; type:uuid; primaryKey; default:gen_random_uuid()"`
	Title       string          `json:"title" gorm:"column:title; not null; size:64;"`
	Description string          `json:"description" gorm:"column:description; not null; size:256;"`
	Type        enums.BadgeType `json:"type" gorm:"column:type; type:BadgeType; not null; default:'Bronze';"`
	ImageURL    *string         `json:"imageURL" gorm:"column:image_url;"`
	CreatedAt   time.Time       `json:"createdAt" gorm:"column:created_at; type:timestamptz; not null; autoCreateTime:true;"`

	// relation
	Users []User `json:"users" gorm:"-"` // many2many:\"UsersToBadgesTable\"; foreignKey:Id; joinForeignKey:BadgeId; references:Id; joinReferences:UserId;
}

func (Badge) TableName() string {
	return shared.ValidTableName_BadgeTable.String()
}
