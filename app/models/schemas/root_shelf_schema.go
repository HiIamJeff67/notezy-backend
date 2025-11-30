package schemas

import (
	"time"

	"github.com/google/uuid"

	gqlmodels "notezy-backend/app/graphql/models"
	types "notezy-backend/shared/types"
)

type RootShelf struct {
	Id             uuid.UUID  `json:"id" gorm:"column:id; type:uuid; primaryKey; default:gen_random_uuid();"`
	OwnerId        uuid.UUID  `json:"ownerId" gorm:"column:owner_id; type:uuid; not null; index:shelf_idx_owner_id_name,unique"`
	Name           string     `json:"name" gorm:"column:name; size:128; not null; default:'undefined'; index:shelf_idx_owner_id_name,unique"`
	SubShelfCount  int32      `json:"subShelfCount" gorm:"column:sub_shelf_count; type:integer; not null; default:1;"`
	ItemCount      int32      `json:"itemCount" gorm:"column:item_count; type:integer; not null; default:0;"`
	LastAnalyzedAt time.Time  `json:"lastAnalyzedAt" gorm:"column:last_analyzed_at; type:timestamptz; not null; default:NOW();"`
	DeletedAt      *time.Time `json:"deletedAt" gorm:"column:deleted_at; type:timestamptz; default:null;"`
	UpdatedAt      time.Time  `json:"updatedAt" gorm:"column:updated_at; type:timestamptz; not null; autoUpdateTime:true;"`
	CreatedAt      time.Time  `json:"createdAt" gorm:"column:created_at; type:timestamptz; not null; autoCreateTime:true;"`

	// relations
	SubShelves     []SubShelf       `json:"subShelves" gorm:"foreignKey:RootShelfId; references:Id; constraint:OnUpdate:CASCADE, OnDelete:CASCADE;"`
	UsersToShelves []UsersToShelves `json:"usersToShelves" gorm:"foreignKey:RootShelfId; references:Id; constraint:OnUpdate:CASCADE, OnDelete:CASCADE;"`
}

// Shelf Table Name
func (RootShelf) TableName() string {
	return types.TableName_RootShelfTable.String()
}

// Shelf Table Relations
type RootShelfRelation types.RelationName

const (
	RootShelfRelation_SubShelves     RootShelfRelation = "SubShelves"
	RootShelfRelation_UsersToShelves RootShelfRelation = "UsersToShelves"
)

/* ============================== Relative Type Conversion ============================== */

func (rs *RootShelf) ToPrivateRootShelf() *gqlmodels.PrivateRootShelf {
	return &gqlmodels.PrivateRootShelf{
		ID:             rs.Id,
		Name:           rs.Name,
		SubShelfCount:  rs.SubShelfCount,
		ItemCount:      rs.ItemCount,
		LastAnalyzedAt: rs.LastAnalyzedAt,
		DeletedAt:      rs.DeletedAt,
		UpdatedAt:      rs.UpdatedAt,
		CreatedAt:      rs.CreatedAt,
		Owner:          make([]*gqlmodels.PublicUser, 0),
	}
}
