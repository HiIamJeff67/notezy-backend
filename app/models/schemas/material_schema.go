package schemas

import (
	"time"

	"github.com/google/uuid"

	enums "notezy-backend/app/models/schemas/enums"
	types "notezy-backend/shared/types"
)

type Material struct {
	Id               uuid.UUID                 `json:"id" gorm:"column:id; type:uuid; primaryKey; not null;"`
	ParentSubShelfId uuid.UUID                 `json:"parentSubShelfId" gorm:"column:parent_sub_shelf_id; type:uuid; not null; uniqueIndex:material_idx_parent_sub_shelf_id_name;"`
	Name             string                    `json:"name" gorm:"column:name; size:128; not null; default:'undefined'; uniqueIndex:material_idx_parent_sub_shelf_id_name;"`
	Type             enums.MaterialType        `json:"type" gorm:"column:type; type:MaterialType; not null; default:'Notebook';"`
	Size             int64                     `json:"size" gorm:"column:size; type:bigint; not null; default:0;"`
	ContentKey       string                    `json:"contentKey" gorm:"column:content_key; unique; not null;"`
	ContentType      enums.MaterialContentType `json:"contentType" gorm:"column:content_type; type:MaterialContentType; not null; default:'text/plain';"`
	ParseMediaType   string                    `json:"parseMediaType" gorm:"column:parse_media_type; not null; default:'utf-8';"`
	DeletedAt        *time.Time                `json:"deletedAt" gorm:"column:deleted_at; type:timestamptz; default:null;"`
	UpdatedAt        time.Time                 `json:"updatedAt" gorm:"column:updated_at; type:timestamptz; not null; autoUpdateTime:true;"`
	CreatedAt        time.Time                 `json:"createdAt" gorm:"column:created_at; type:timestamptz; not null; autoCreateTime:true;"`

	// relations
	ParentSubShelf SubShelf `json:"parentSubShelf" gorm:"foreignKey:ParentSubShelfId; references:Id; constraint:OnUpdate:CASCADE, OnDelete:CASCADE;"`
}

// Constraints:
// Material.ParentSubShelfId + Material.Name as the unique constraints (That means the material name should be unique in the current sub shelf layer)

// Material Table Name
func (Material) TableName() string {
	return types.ValidTableName_MaterialTable.String()
}

// Material Table Relations
type MaterialRelation types.ValidTableName

const (
	MaterialRelation_ParentSubShelf MaterialRelation = "ParentSubShelf"
)

/* ============================== Relative Type Conversion ============================== */
