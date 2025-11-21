package schemas

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/datatypes"

	enums "notezy-backend/app/models/schemas/enums"
	types "notezy-backend/shared/types"
)

type SyncBlock struct {
	Id            uuid.UUID       `json:"id" gorm:"column:id; type:uuid; primaryKey; not null; default:gen_random_uuid();"`
	ParentBlockId *uuid.UUID      `json:"parentBlockId" gorm:"column:parent_block_id; type:uuid; check:check_parent_block_is_not_itself,parent_block_id != id;"`
	BlockGroupId  uuid.UUID       `json:"blockGroupId" gorm:"column:block_group_id; type:uuid; not null;"`
	Type          enums.BlockType `json:"type" gorm:"column:type; type:BlockType; not null; default:'paragraph';"`
	Props         datatypes.JSON  `json:"props" gorm:"column:props; type:jsonb; not null; default:'{}';"`
	Content       datatypes.JSON  `json:"content" gorm:"column:content; type:jsonb; default:'[]';"`
	UpdatedAt     time.Time       `json:"updatedAt" gorm:"column:updated_at; type:timestamptz; not null; autoUpdateTime:true;"`
	CreatedAt     time.Time       `json:"createdAt" gorm:"column:created_at; type:timestamptz; not null; autoCreateTime:true;"`

	Parent     *SyncBlock     `json:"parent" gorm:"foreignKey:ParentBlockId; references:Id; contraint:OnUpdate:CASCADE, OnDelete:CASCADE;"`
	Children   []SyncBlock    `json:"children" gorm:"foreignKey:ParentBlockId; references:Id; constraint:OnUpdate:CASCADE, OnDelete:CASCADE;"`
	BlockGroup SyncBlockGroup `json:"blockGroup" gorm:"foreignKey:BlockGroupId; references:Id; constraint:OnUpdate:CASCADE, OnDelete:CASCADE;"`
}

// Sync Block Table Name
func (SyncBlock) TableName() string {
	return types.TableName_SyncBlockTableName.String()
}

// Sync Block Relations
type SyncBlockRelation types.RelationName

const (
	SyncBlockRelation_Parent     SyncBlockRelation = "Parent"
	SyncBlockRelation_Children   SyncBlockRelation = "Children"
	SyncBlockRelation_BlockGroup SyncBlockRelation = "BlockGroup"
)
