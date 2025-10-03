package schemas

import (
	"time"

	"github.com/google/uuid"

	gqlmodels "notezy-backend/app/graphql/models"
	types "notezy-backend/shared/types"
)

type SubShelf struct {
	Id             uuid.UUID       `json:"id" gorm:"column:id; type:uuid; primary; default:gen_random_uuid();"`
	Name           string          `json:"name" gorm:"column:name; size:128; not null; default:'undefined'; uniqueIndex:sub_shelf_idx_name_root_shelf_id_path_deleted_at;"`
	RootShelfId    uuid.UUID       `json:"rootShelfId" gorm:"column:root_shelf_id; type:uuid; not null; uniqueIndex:sub_shelf_idx_name_root_shelf_id_path_deleted_at;"`
	PrevSubShelfId *uuid.UUID      `json:"prevSubShelfId" gorm:"column:prev_sub_shelf_id; type:uuid; check:prev_sub_shelf_id_check,prev_sub_shelf_id != id;"`
	Path           types.UUIDArray `json:"path" gorm:"column:path; type:uuid[]; not null; default:'{}'; check:path_length_check,cardinality(path) >= 0 AND cardinality(path) <= 100; uniqueIndex:sub_shelf_idx_name_root_shelf_id_path_deleted_at;"`
	DeletedAt      *time.Time      `json:"deletedAt" gorm:"column:deleted_at; type:timestamptz; default:null; uniqueIndex:sub_shelf_idx_name_root_shelf_id_path_deleted_at;"`
	UpdatedAt      time.Time       `json:"updatedAt" gorm:"column:updated_at; type:timestamptz; not null; autoUpdateTime:true;"`
	CreatedAt      time.Time       `json:"createdAt" gorm:"column:created_at; type:timestamptz; not null; autoCreateTime:true;"`

	// relations
	RootShelf      RootShelf  `json:"rootShelf" gorm:"foreignKey:RootShelfId; references:Id; constraint:OnUpdate:CASCADE, OnDelete:CASCADE;"`
	NextSubShelves []SubShelf `json:"subShelves" gorm:"foreignKey:PrevSubShelfId; references:Id; constraint:OnUpdate:CASCADE, OnDelete:CASCADE;"`
	Materials      []Material `json:"materials" gorm:"foreignKey:ParentSubShelfId; constraint:OnUpdate:CASCADE, OnDelete:CASCADE;"`
}

// Contraints:
// 1. SubShelf.PrevSubShelfId != SubShelf.Id (Do NOT point to itself)
// 2. Use Name + RootShelfId + Path as the unique constraints (That means in some specific layers of a RootShelf, the name of a SubShelf should be unique)
// 	  Also add the DeletedAt as a part of the constraints, since the user may create couple folders with the same name, and delete it, and then create it again and again...
// 3. The maximum length of SubShelf.Path should be less than or equal to 100

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

/* ============================== Relative Type Conversion ============================== */

func (ss *SubShelf) ToPrivateSubShelf() *gqlmodels.PrivateSubShelf {
	return &gqlmodels.PrivateSubShelf{
		ID:             ss.Id,
		Name:           ss.Name,
		RootShelfID:    ss.RootShelfId,
		PrevSubShelfID: ss.PrevSubShelfId,
		Path:           ss.Path,
		DeletedAt:      ss.DeletedAt,
		UpdatedAt:      ss.UpdatedAt,
		CreatedAt:      ss.CreatedAt,
		RootShelf:      nil,
		NextSubShelves: make([]*gqlmodels.PrivateSubShelf, 0),
		Materials:      make([]*gqlmodels.PrivateMaterial, 0),
	}
}
