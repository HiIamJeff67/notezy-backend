package schemas

import (
	"time"

	"github.com/google/uuid"

	gqlmodels "notezy-backend/app/graphql/models"
	shared "notezy-backend/shared"
)

type Shelf struct {
	Id                       uuid.UUID `json:"id" gorm:"column:id; type:uuid; primaryKey; default:gen_random_uuid();"`
	OwnerId                  uuid.UUID `json:"ownerId" gorm:"column:owner_id; type:uuid; not null; index:shelf_idx_owner_id_name,unique"`
	Name                     string    `json:"name" gorm:"column:name; size:128; not null; default:'undefined'; index:shelf_idx_owner_id_name,unique"`
	EncodedStructure         []byte    `json:"encodedStructure" gorm:"column:encoded_structure; type:bytea; not null;"`
	EncodedStructureByteSize int64     `json:"encodedStructureByteSize" gorm:"column:encoded_structure_byte_size; type:bigint; not null; default:36;"`
	TotalShelfNodes          int32     `json:"totalShelfNodes" gorm:"column:total_shelf_nodes; type:integer; not null; default:1;"`
	TotalMaterials           int32     `json:"totalMaterials" gorm:"column:total_materials; type:integer; not null; default:0;"`
	MaxWidth                 int32     `json:"maxWidth" gorm:"column:max_width; type:integer; not null; default:1;"`
	MaxDepth                 int32     `json:"maxDepth" gorm:"column:max_depth; type:integer; not null; default:1;"`
	UpdatedAt                time.Time `json:"updatedAt" gorm:"column:updated_at; type:timestamptz; not null; autoUpdateTime:true;"`
	CreatedAt                time.Time `json:"createdAt" gorm:"column:created_at; type:timestamptz; not null; autoCreateTime:true;"`

	// relations
	Materials      []Material       `json:"materials" gorm:"foreignKey:RootShelfId;"`
	UsersToShelves []UsersToShelves `json:"usersToShelves" gorm:"foreignKey:ShelfId;"`
}

// Shelf Table Name
func (Shelf) TableName() string {
	return shared.ValidTableName_ShelfTable.String()
}

// Shelf Table Relations
type ShelfRelation shared.ValidTableName

const (
	ShelfRelation_Materials      ShelfRelation = "Materials"
	ShelfRelation_UsersToShelves ShelfRelation = "UsersToShelves"
)

/* ============================== Relative Tyoe Conversion ============================== */

func (s *Shelf) ToPrivateShelf() *gqlmodels.PrivateShelf {
	return &gqlmodels.PrivateShelf{
		ID:                       s.Id,
		Name:                     s.Name,
		EncodedStructure:         s.EncodedStructure,
		EncodedStructureByteSize: s.EncodedStructureByteSize,
		TotalShelfNodes:          s.TotalShelfNodes,
		TotalMaterials:           s.TotalMaterials,
		MaxWidth:                 s.MaxDepth,
		MaxDepth:                 s.MaxDepth,
		UpdatedAt:                s.UpdatedAt,
		CreatedAt:                s.CreatedAt,
		UsersToShelves:           make([]*gqlmodels.PrivateUsersToShelves, 0),
	}
}
