package schemas

import (
	"time"

	"github.com/google/uuid"

	shared "notezy-backend/shared"
)

type Shelf struct {
	Id               uuid.UUID `json:"id" gorm:"column:id; type:uuid; primaryKey; default:gen_random_uuid();"`
	OwnerId          uuid.UUID `json:"ownerId" gorm:"column:owner_id; type:uuid; not null; index:shelf_idx_owner_id_name,unique"`
	Name             string    `json:"name" gorm:"column:name; size:128; not null; default:'undefined'; index:shelf_idx_owner_id_name,unique"`
	EncodedStructure []byte    `json:"encodedStructure" gorm:"column:encoded_structure; type:bytea; default:'null';"`
	UpdatedAt        time.Time `json:"updatedAt" gorm:"column:updated_at; type:timestamptz; not null; autoUpdateTime:true;"`
	CreatedAt        time.Time `json:"createdAt" gorm:"column:created_at; type:timestamptz; not null; autoCreateTime:true;"`

	// relations
	Users []User `json:"users" gorm:"many2many:\"UsersToShelvesTable\"; foreignKey:Id; joinForeignKey:ShelfId; references:Id; joinReferences:UserId;"`
}

// Shelf Table Name
func (Shelf) TableName() string {
	return shared.ValidTableName_ShelfTable.String()
}

// Shelf Table Relations
type ShelfRelations shared.ValidTableName

const (
	ShelfRelations_Users ShelfRelations = "Users"
)
