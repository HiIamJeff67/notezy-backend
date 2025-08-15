package schemas

import (
	"time"

	"github.com/google/uuid"

	enums "notezy-backend/app/models/schemas/enums"
	shared "notezy-backend/shared"
)

type UsersToShelves struct {
	UserId     uuid.UUID                     `json:"userId" gorm:"column:user_id; type:uuid; primaryKey;"`
	ShelfId    uuid.UUID                     `json:"shelfId" gorm:"column:shelf_id; type:uuid; primaryKey;"`
	Permission enums.AccessControlPermission `json:"permission" gorm:"column:permission; type:AccessControlPermission; not null; default:'Read';"`
	UpdatedAt  time.Time                     `json:"updatedAt" gorm:"column:updated_at; tpye:timestamptz; not null; autoUpdateTime:true;"`
	CreatedAt  time.Time                     `json:"createdAt" gorm:"column:created_at; type:timestamptz; not null; autoCreateTime:true;"`

	// relations
	User  User  `gorm:"foreignKey:UserId; reference:Id; constraint:OnUpdate:CASCADE, OnDelete:CASCADE;"`
	Shelf Shelf `gorm:"foreignKey:ShelfId; reference:Id; constraint:OnUpdate:CASCADE, OnDelete:CASCADE;"`
}

// Users To Shelves Table Name
func (UsersToShelves) TableName() string {
	return shared.ValidTableName_UsersToShelvesTable.String()
}

// Users To Shelves Table Relations
type UsersToShelvesRelation shared.ValidTableName

const (
	UsersToShelvesRelation_User  UsersToShelvesRelation = "User"
	UsersToShelvesRelation_Shelf UsersToShelvesRelation = "Shelf"
)
