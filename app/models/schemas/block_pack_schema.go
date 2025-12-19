package schemas

import (
	"time"

	"github.com/google/uuid"

	enums "notezy-backend/app/models/schemas/enums"
	types "notezy-backend/shared/types"
)

type BlockPack struct {
	Id                  uuid.UUID                     `json:"id" gorm:"column:id; type:uuid; primaryKey; not null; default:gen_random_uuid();"`
	ParentSubShelfId    uuid.UUID                     `json:"parentSubShelfId" gorm:"column:parent_sub_shelf_id; type:uuid; not null; uniqueIndex:block_pack_idx_parent_sub_shelf_id_name_deleted_at;"`
	FinalBlockGroupId   *uuid.UUID                    `json:"finalBlockGroupId" gorm:"column:final_block_group_id; type:uuid; default:null;"`
	Name                string                        `json:"name" gorm:"column:name; size:128; not null; default:'undefined'; uniqueIndex:block_pack_idx_parent_sub_shelf_id_name_deleted_at;"`
	Icon                *enums.SupportedBlockPackIcon `json:"icon" gorm:"column:icon; type:\"SupportedBlockPackIcon\"; default:null;"`
	HeaderBackgroundURL *string                       `json:"headerBackgroundURL" gorm:"column:header_background_url; default:null;"`
	BlockCount          int32                         `json:"blockCount" gorm:"block_count; type:integer; not null; default:0; check:check_max_block_count,block_count <= 100;"`
	DeletedAt           *time.Time                    `json:"deletedAt" gorm:"column:deleted_at; type:timestamptz; default:null; uniqueIndex:block_pack_idx_parent_sub_shelf_id_name_deleted_at;"`
	UpdatedAt           time.Time                     `json:"updatedAt" gorm:"column:updated_at; type:timestamptz; not null; autoUpdateTime:true;"`
	CreatedAt           time.Time                     `json:"createdAt" gorm:"column:created_at; type:timestamptz; not null; autoCreateTime:true;"`

	// relations
	ParentSubShelf  SubShelf     `json:"parentSubShelf" gorm:"foreignKey:ParentSubShelfId; references:Id; constraint:OnUpdate:CASCADE, OnDelete:CASCADE;"`
	BlockGroups     []BlockGroup `json:"blockGroups" gorm:"foreignKey:BlockPackId; references:Id; constraint:OnUpdate:CASCADE, OnDelete:CASCADE;"`
	FinalBlockGroup *BlockGroup  `json:"finalBlockGroups" gorm:"foreignKey:FinalBlockGroupId; references:Id; constraint:OnUpdate:CASCADE, OnDelete:SET NULL;"`
}

// BlockPack Table Name
func (BlockPack) TableName() string {
	return types.TableName_BlockPackTable.String()
}

// BlockPack Table Relations
type BlockPackRelation types.RelationName

const (
	BlockPackRelation_ParentSubShelf   BlockPackRelation = "ParentSubShelf"
	BlockPackRelation_BlockGroups      BlockPackRelation = "BlockGroups"
	BlockPackRelation_FinalBlockGroups BlockPackRelation = "FinalBlockGroup"
)
