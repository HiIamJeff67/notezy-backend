package schemas

import (
	"time"

	"github.com/google/uuid"

	types "notezy-backend/shared/types"
)

type SubShelf struct {
	Id             uuid.UUID   `json:"id" gorm:"column:id; type:uuid; primary; default:gen_random_uuid();"`
	Name           string      `json:"name" gorm:"column:name; size:128; not null; default:'undefined'; uniqueIndex:sub_shelf_idx_name_root_shelf_id_path;"`
	RootShelfId    uuid.UUID   `json:"rootShelfId" gorm:"column:root_shelf_id; type:uuid; not null; uniqueIndex:sub_shelf_idx_name_root_shelf_id_path;"`
	PrevSubShelfId *uuid.UUID  `json:"prevSubShelfId" gorm:"column:prev_sub_shelf_id; type:uuid;"`
	Path           []uuid.UUID `json:"path" gorm:"column:path; type:uuid[]; not null; check:path_length_check,length(path) >= 0 AND length(path) <= 100; uniqueIndex:sub_shelf_idx_name_root_shelf_id_path;"`
	DeletedAt      time.Time   `json:"deletedAt" gorm:"column:deleted_at; type:timestamptz;"`
	UpdatedAt      time.Time   `json:"updatedAt" gorm:"column:updated_at; type:timestamptz; not null; autoUpdateTime:true;"`
	CreatedAt      time.Time   `json:"createdAt" gorm:"column:created_at; type:timestamptz; not null; autoCreateTime:true;"`

	// relations
	RootShelf      RootShelf  `json:"rootShelf" gorm:"foreignKey:RootShelfId; references:Id; constraint:OnUpdate:CASCADE, OnDelete:CASCADE;"`
	NextSubShelves []SubShelf `json:"subShelves" gorm:"foreignKey:PrevSubShelfId; references:Id; constraint:OnUpdate:CASCADE, OnDelete:CASCADE;"`
	Materials      []Material `json:"materials" gorm:"foreignKey:ParentSubShelfId; constraint:OnUpdate:CASCADE, OnDelete:CASCADE;"`
}

// contraints:
// 1. Name + ParentShelfId : In the current SubShelf, it should be unique one their names
// 2. RootShelfId + Path : In the current SubShelf, the path from the root to it should be unique

func (SubShelf) TableName() string {
	return types.ValidTableName_SubShelfTable.String()
}

// SubShelf Table Relations
type SubShelfRelation types.ValidTableName

const (
	SubShelfRelation_RootShelf      SubShelfRelation = "RootShelf"
	SubShelfRelation_NextSubShelves SubShelfRelation = "NextSubShelves"
	SubShelfRelation_Materials      SubShelfRelation = "Materials"
)
