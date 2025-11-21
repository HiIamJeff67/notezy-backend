package schemas

import (
	"time"

	"github.com/google/uuid"

	enums "notezy-backend/app/models/schemas/enums"
	types "notezy-backend/shared/types"
)

type UsersToShelves struct {
	UserId      uuid.UUID                     `json:"userId" gorm:"column:user_id; type:uuid; primaryKey;"`
	RootShelfId uuid.UUID                     `json:"rootShelfId" gorm:"column:root_shelf_id; type:uuid; primaryKey; uniqueIndex:idx_root_shelf_owner,where:permission = 'Owner';"`
	Permission  enums.AccessControlPermission `json:"permission" gorm:"column:permission; type:AccessControlPermission; not null; default:'Read';"`
	UpdatedAt   time.Time                     `json:"updatedAt" gorm:"column:updated_at; tpye:timestamptz; not null; autoUpdateTime:true;"`
	CreatedAt   time.Time                     `json:"createdAt" gorm:"column:created_at; type:timestamptz; not null; autoCreateTime:true;"`

	// relations
	User      User      `gorm:"foreignKey:UserId; reference:Id; constraint:OnUpdate:CASCADE, OnDelete:CASCADE;"`
	RootShelf RootShelf `gorm:"foreignKey:RootShelfId; reference:Id; constraint:OnUpdate:CASCADE, OnDelete:CASCADE;"`
}

// Users To Shelves Table Name
func (UsersToShelves) TableName() string {
	return types.TableName_UsersToShelvesTable.String()
}

// Users To Shelves Table Relations
type UsersToShelvesRelation types.RelationName

const (
	UsersToShelvesRelation_User      UsersToShelvesRelation = "User"
	UsersToShelvesRelation_RootShelf UsersToShelvesRelation = "RootShelf"
)
