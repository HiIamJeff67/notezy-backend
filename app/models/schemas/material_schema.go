package schemas

import (
	"time"

	"github.com/google/uuid"

	enums "notezy-backend/app/models/schemas/enums"
	types "notezy-backend/shared/types"
)

/* ====================================================================================================
|| ! Note that we can modify or manipulate the Matertial instance without first access to the Shelf  ||
==================================================================================================== */

type Material struct {
	Id            uuid.UUID                 `json:"id" gorm:"column:id; type:uuid; primaryKey; not null;"`
	RootShelfId   uuid.UUID                 `json:"rootShelfId" gorm:"column:root_shelf_id; type:uuid; primaryKey; not null; index;"`
	ParentShelfId uuid.UUID                 `json:"parentShelfId" gorm:"column:parent_shelf_id; type:uuid; not null;"`
	Name          string                    `json:"name" gorm:"column:name; size:128; not null; default:'undefined';"`
	Type          enums.MaterialType        `json:"type" gorm:"column:type; type:MaterialType; not null; default:'Notebook';"`
	Size          int64                     `json:"size" gorm:"column:size; type:bigint; not null; default:0"`
	ContentKey    string                    `json:"contentKey" gorm:"column:content_key; not null; default:'';"`
	ContentType   enums.MaterialContentType `json:"contentType" gorm:"column:content_type; type:MaterialContentType; not null; default:'Text_Markdown';"`
	DeletedAt     *time.Time                `json:"deletedAt" gorm:"column:deleted_at; type:timestamptz;"`
	UpdatedAt     time.Time                 `json:"updatedAt" gorm:"column:updated_at; type:timestamptz; not null; autoUpdateTime:true;"`
	CreatedAt     time.Time                 `json:"createdAt" gorm:"column:created_at; type:timestamptz; not null; autoCreateTime:true;"`

	// relations
	RootShelf Shelf `json:"rootShelf" gorm:"foreignKey:RootShelfId; references:Id; constraint:OnUpdate:CASCADE, OnDelete:CASCADE;"`
}

// Material Table Name
func (Material) TableName() string {
	return types.ValidTableName_MaterialTable.String()
}

// Material Table Relations
type MaterialRelation string

const (
	MaterialRelation_Shelf MaterialRelation = "Shelf"
)

/* ============================== Relative Type Conversion ============================== */
