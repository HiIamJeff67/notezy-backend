package schemas

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/datatypes"

	enums "notezy-backend/app/models/schemas/enums"
	types "notezy-backend/shared/types"
)

type Block struct {
	Id            uuid.UUID       `json:"id" gorm:"column:id; type:uuid; primaryKey; not null; default:gen_random_uuid();"`
	ParentBlockId *uuid.UUID      `json:"parentBlockId" gorm:"column:parent_block_id; type:uuid; check:check_parent_block_id_is_not_itself,parent_block_id != id;"`
	BlockGroupId  uuid.UUID       `json:"blockGroupId" gorm:"column:block_group_id; type:uuid; not null;"`
	Type          enums.BlockType `json:"type" gorm:"column:type; type:BlockType; not null; default:'paragraph';"`
	Props         datatypes.JSON  `json:"props" gorm:"column:props; type:jsonb; not null; default:'{}'; check:check_props_size,octet_length(props::text) <= 10240;"`
	Content       datatypes.JSON  `json:"content" gorm:"column:content; type:jsonb; default:'[]'; check:check_content_size,octet_length(content::text) <= 10240;"`
	DeletedAt     *time.Time      `json:"deletedAt" gorm:"column:deleted_at; type:timestamptz; default:null;"`
	UpdatedAt     time.Time       `json:"updatedAt" gorm:"column:updated_at; type:timestamptz; not null; autoUpdateTime:true;"`
	CreatedAt     time.Time       `json:"createdAt" gorm:"column:created_at; type:timestamptz; not null; autoCreateTime:true;"`

	// relations
	Parent     *Block     `json:"parent" gorm:"foreignKey:ParentBlockId; references:Id; constraint:OnUpdate:CASCADE, OnDelete:CASCADE;"`
	Children   []Block    `json:"children" gorm:"foreignKey:ParentBlockId; references:Id; constraint:OnUpdate:CASCADE, OnDelete:CASCADE;"`
	BlockGroup BlockGroup `json:"blockGroup" gorm:"foreignKey:BlockGroupId; references:Id; constraint:OnUpdate:CASCADE, OnDelete:CASCADE;"`
}

// Root Block Table Name
func (Block) TableName() string {
	return types.TableName_BlockTable.String()
}

// Root Block Table Relations
type BlockRelation types.RelationName

const (
	BlockRelation_Parent     BlockRelation = "Parent"
	BlockRelation_Children   BlockRelation = "Children"
	BlockRelation_BlockGroup BlockRelation = "BlockGroup"
)
