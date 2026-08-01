package schemas

import (
	"time"

	"github.com/google/uuid"

	types "github.com/HiIamJeff67/notezy-backend/internal/shared/types"
)

type SubShelf struct {
	Id             uuid.UUID       `json:"id" gorm:"column:id; type:uuid; primary; default:gen_random_uuid();"`
	Name           string          `json:"name" gorm:"column:name; size:128; not null; default:'undefined';"` // Previous unique-name constraint: uniqueIndex:sub_shelf_idx_name_root_shelf_id_path,where:deleted_at IS NULL
	RootShelfId    uuid.UUID       `json:"rootShelfId" gorm:"column:root_shelf_id; type:uuid; not null;"`     // Previous unique-name constraint: uniqueIndex:sub_shelf_idx_name_root_shelf_id_path,where:deleted_at IS NULL
	PrevSubShelfId *uuid.UUID      `json:"prevSubShelfId" gorm:"column:prev_sub_shelf_id; type:uuid; check:sub_shelf_check_prev_sub_shelf_id,prev_sub_shelf_id != id;"`
	Path           types.UUIDArray `json:"path" gorm:"column:path; type:uuid[]; not null; default:'{}'; check:sub_shelf_check_path_length,cardinality(path) >= 0 AND cardinality(path) <= 100;"` // Previous unique-name constraint: uniqueIndex:sub_shelf_idx_name_root_shelf_id_path,where:deleted_at IS NULL
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
// 2. Duplicate names are allowed. Previously, Name + RootShelfId + Path was unique while DeletedAt was NULL.
// 3. The maximum length of SubShelf.Path should be less than or equal to 100.

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
