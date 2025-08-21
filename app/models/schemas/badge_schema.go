package schemas

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"

	gqlmodels "notezy-backend/app/graphql/models"
	enums "notezy-backend/app/models/schemas/enums"
	types "notezy-backend/shared/types"
)

type Badge struct {
	Id          uuid.UUID       `json:"id" gorm:"column:id; type:uuid; primaryKey; default:gen_random_uuid()"`
	PublicId    string          `json:"publicId" gorm:"column:public_id; unique; not null; default:'';"`
	Title       string          `json:"title" gorm:"column:title; not null; size:64;"`
	Description string          `json:"description" gorm:"column:description; not null; size:256;"`
	Type        enums.BadgeType `json:"type" gorm:"column:type; type:BadgeType; not null; default:'Bronze';"`
	ImageURL    *string         `json:"imageURL" gorm:"column:image_url;"`
	CreatedAt   time.Time       `json:"createdAt" gorm:"column:created_at; type:timestamptz; not null; autoCreateTime:true;"`

	// relation
	UsersToBadges []UsersToBadges `json:"usersToBadges" gorm:"foreignKey:BadgeId;"`
}

// Badge Table Name
func (Badge) TableName() string {
	return types.ValidTableName_BadgeTable.String()
}

// Badge Table Relations
type BadgeRelation string

const (
	BadgeRelation_UsersToBadges BadgeRelation = "UsersToBadges"
)

/* ============================== Relative Type Conversions ============================== */

func (b *Badge) ToPublicBadge() *gqlmodels.PublicBadge {
	return &gqlmodels.PublicBadge{
		PublicID:    b.PublicId,
		Title:       b.Title,
		Description: b.Description,
		Type:        b.Type,
		ImageURL:    b.ImageURL,
		CreatedAt:   b.CreatedAt,
		Users:       []*gqlmodels.PublicUser{},
	}
}

/* ============================== Trigger Hook ============================== */

func (b *Badge) BeforeCreate(db *gorm.DB) error {
	if b.PublicId == "" {
		b.PublicId = uuid.NewString()
	}
	return nil
}
