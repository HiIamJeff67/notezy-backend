package schemas

import (
	"time"

	"github.com/google/uuid"

	enums "notezy-backend/app/models/schemas/enums"
	shared "notezy-backend/shared"
)

/* ====================================================================================================
|| ! Note that we can modify or manipulate the Matertial instance without first access to the Shelf  ||
==================================================================================================== */

type Material struct {
	Id          uuid.UUID                 `json:"id" gorm:"column:id; type:uuid; primaryKey; default:gen_random_uuid();"`
	RootShelfId uuid.UUID                 `json:"rootShelfId" gorm:"column:root_shelf_id; type:uuid; not null; index;"`
	Name        string                    `json:"name" gorm:"column:name; size:128; not null; default:'undefined';"`
	Content     []byte                    `json:"content" gorm:"column:content; type:bytea; default:'null';"`
	Type        enums.MaterialType        `json:"type" gorm:"column:type; type:MaterialType; not null; default:'Notebook';"`
	ContentType enums.MaterialContentType `json:"contentType" gorm:"column:content_type; type:MaterialContentType; not null; default:'text/markdown';"`
	UpdatedAt   time.Time                 `json:"updatedAt" gorm:"column:updated_at; type:timestamptz; not null; autoUpdateTime:true;"`
	CreatedAt   time.Time                 `json:"createdAt" gorm:"column:created_at; type:timestamptz; not null; autoCreateTime:true;"`

	// relations
	RootShelf Shelf `json:"rootShelf" gorm:"foreignKey:RootShelfId; references:Id;"`
}

// Material Table Name
func (Material) TableName() string {
	return shared.ValidTableName_MaterialTable.String()
}

// Material Table Relations
type MaterialRelation string

const (
	MaterialRelation_Shelf MaterialRelation = "Shelf"
)

/* ============================== Relative Type Conversion ============================== */
