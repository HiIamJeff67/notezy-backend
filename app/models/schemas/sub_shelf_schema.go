package schemas

import (
	"time"

	"github.com/google/uuid"

	gqlmodels "github.com/HiIamJeff67/notezy-backend/app/graphql/models"
	types "github.com/HiIamJeff67/notezy-backend/shared/types"
)

type SubShelf struct {
	Id             uuid.UUID       `json:"id" gorm:"column:id; type:uuid; primary; default:gen_random_uuid();"`
	Name           string          `json:"name" gorm:"column:name; size:128; not null; default:'undefined'; uniqueIndex:sub_shelf_idx_name_root_shelf_id_path,where:deleted_at IS NULL;"`
	RootShelfId    uuid.UUID       `json:"rootShelfId" gorm:"column:root_shelf_id; type:uuid; not null; uniqueIndex:sub_shelf_idx_name_root_shelf_id_path,where:deleted_at IS NULL;"`
	PrevSubShelfId *uuid.UUID      `json:"prevSubShelfId" gorm:"column:prev_sub_shelf_id; type:uuid; check:sub_shelf_check_prev_sub_shelf_id,prev_sub_shelf_id != id;"`
	Path           types.UUIDArray `json:"path" gorm:"column:path; type:uuid[]; not null; default:'{}'; check:sub_shelf_check_path_length,cardinality(path) >= 0 AND cardinality(path) <= 100; uniqueIndex:sub_shelf_idx_name_root_shelf_id_path,where:deleted_at IS NULL;"`
	DeletedAt      *time.Time      `json:"deletedAt" gorm:"column:deleted_at; type:timestamptz; default:null;"`
	UpdatedAt      time.Time       `json:"updatedAt" gorm:"column:updated_at; type:timestamptz; not null; autoUpdateTime:true;"`
	CreatedAt      time.Time       `json:"createdAt" gorm:"column:created_at; type:timestamptz; not null; autoCreateTime:true;"`

	// relations
	RootShelf      RootShelf   `json:"rootShelf" gorm:"foreignKey:RootShelfId; references:Id; constraint:OnUpdate:CASCADE, OnDelete:CASCADE;"`
	NextSubShelves []SubShelf  `json:"subShelves" gorm:"foreignKey:PrevSubShelfId; references:Id; constraint:OnUpdate:CASCADE, OnDelete:CASCADE;"`
	Materials      []Material  `json:"materials" gorm:"foreignKey:ParentSubShelfId; constraint:OnUpdate:CASCADE, OnDelete:CASCADE;"`
	BlockPacks     []BlockPack `json:"blockSets" gorm:"foreignKey:ParentSubShelfId; references:Id; constraint:OnUpdate:CASCADE, OnDelete:CASCADE;"`
	Items          []Item      `json:"items" gorm:"foreignKey:ParentSubShelfId; references:Id; constraint:OnUpdate:CASCADE, OnDelete:CASCADE;"`
}

// Constraints:
// 1. SubShelf.PrevSubShelfId != SubShelf.Id (Do NOT point to itself)
// 2. Use Name + RootShelfId + Path as the unique constraints (That means in some specific layers of a RootShelf, the name of a SubShelf should be unique)
// 	  Also add the DeletedAt as a part of the constraints, since the user may create couple folders with the same name, and delete it, and then create it again and again...
// 3. The maximum length of SubShelf.Path should be less than or equal to 100

func (SubShelf) TableName() string {
	return types.TableName_SubShelfTable.String()
}

// SubShelf Table Relations
type SubShelfRelation types.RelationName

const (
	SubShelfRelation_RootShelf      SubShelfRelation = "RootShelf"
	SubShelfRelation_NextSubShelves SubShelfRelation = "NextSubShelves"
	SubShelfRelation_Materials      SubShelfRelation = "Materials"
	SubShelfRelation_BlockPacks     SubShelfRelation = "BlockPacks"
	SubShelfRelation_Items          SubShelfRelation = "Items"
)

/* ============================== Relative Type Conversion ============================== */

func (ss *SubShelf) ToPrivateSubShelf() *gqlmodels.PrivateSubShelf {
	nextSubShelfIds := make([]uuid.UUID, 0, len(ss.NextSubShelves))
	for _, nextSubShelf := range ss.NextSubShelves {
		nextSubShelfIds = append(nextSubShelfIds, nextSubShelf.Id)
	}

	itemIds := make([]uuid.UUID, 0, len(ss.Items))
	for _, item := range ss.Items {
		itemIds = append(itemIds, item.Id)
	}

	return &gqlmodels.PrivateSubShelf{
		ID:              ss.Id,
		Name:            ss.Name,
		RootShelfID:     ss.RootShelfId,
		PrevSubShelfID:  ss.PrevSubShelfId,
		Path:            ss.Path,
		DeletedAt:       ss.DeletedAt,
		UpdatedAt:       ss.UpdatedAt,
		CreatedAt:       ss.CreatedAt,
		NextSubShelfIds: nextSubShelfIds,
		ItemIds:         itemIds,
	}
}
