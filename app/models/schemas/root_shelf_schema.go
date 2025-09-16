package schemas

import (
	"time"

	"github.com/google/uuid"

	gqlmodels "notezy-backend/app/graphql/models"
	types "notezy-backend/shared/types"
)

type RootShelf struct {
	Id              uuid.UUID `json:"id" gorm:"column:id; type:uuid; primaryKey; default:gen_random_uuid();"`
	OwnerId         uuid.UUID `json:"ownerId" gorm:"column:owner_id; type:uuid; not null; index:shelf_idx_owner_id_name,unique"`
	Name            string    `json:"name" gorm:"column:name; size:128; not null; default:'undefined'; index:shelf_idx_owner_id_name,unique"`
	TotalShelfNodes int32     `json:"totalShelfNodes" gorm:"column:total_shelf_nodes; type:integer; not null; default:1;"`
	TotalMaterials  int32     `json:"totalMaterials" gorm:"column:total_materials; type:integer; not null; default:0;"`
	LastAnalyzedAt  time.Time `json:"lastAnalyzedAt" gorm:"column:last_analyzed_at; type:timestamptz; not null; default:NOW();"`
	DeletedAt       time.Time `json:"deletedAt" gorm:"column:deleted_at; type:timestamptz;"`
	UpdatedAt       time.Time `json:"updatedAt" gorm:"column:updated_at; type:timestamptz; not null; autoUpdateTime:true;"`
	CreatedAt       time.Time `json:"createdAt" gorm:"column:created_at; type:timestamptz; not null; autoCreateTime:true;"`

	// relations
	SubShelves     []SubShelf       `json:"subShelves" gorm:"foreignKey:RootShelfId; references:Id; constraint:OnUpdate:CASCADE, OnDelete:CASCADE;"`
	UsersToShelves []UsersToShelves `json:"usersToShelves" gorm:"foreignKey:ShelfId; constraint:OnUpdate:CASCADE, OnDelete:CASCADE;"`
}

// Shelf Table Name
func (RootShelf) TableName() string {
	return types.ValidTableName_RootShelfTable.String()
}

// Shelf Table Relations
type RootShelfRelation types.ValidTableName

const (
	RootShelfRelation_UsersToShelves RootShelfRelation = "UsersToShelves"
)

/* ============================== Relative Tyoe Conversion ============================== */

func (s *RootShelf) ToPrivateShelf() *gqlmodels.PrivateRootShelf {
	return &gqlmodels.PrivateRootShelf{
		ID:              s.Id,
		Name:            s.Name,
		TotalShelfNodes: s.TotalShelfNodes,
		TotalMaterials:  s.TotalMaterials,
		UpdatedAt:       s.UpdatedAt,
		CreatedAt:       s.CreatedAt,
		Owner:           make([]*gqlmodels.PublicUser, 0),
	}
}
