package schemas

import (
	"time"

	"github.com/google/uuid"

	types "notezy-backend/shared/types"
)

type BlockGroup struct {
	Id               uuid.UUID  `json:"id" gorm:"column:id; type:uuid; primaryKey; not null; default:gen_random_uuid();"`
	OwnerId          uuid.UUID  `json:"ownerId" gorm:"column:owner_id; type:uuid; not null;"`
	BlockPackId      uuid.UUID  `json:"blockPackId" gorm:"column:block_pack_id; type:uuid; not null; index:block_group_idx_name_block_pack_id_prev_block_group_id;"`
	PrevBlockGroupId *uuid.UUID `json:"prevBlockGroupId" gorm:"column:prev_block_group_id; type:uuid; default:null; index:block_group_idx_name_block_pack_id_prev_block_group_id;"`
	SyncBlockGroupId *uuid.UUID `json:"syncBlockGroupId" gorm:"sync_block_group_id; type:uuid; default:null;"`
	MegaByteSize     float64    `json:"megaByteSize" gorm:"column:mega_byte_size; type:double precision; not null; default:0;"`
	DeletedAt        *time.Time `json:"deletedAt" gorm:"column:deleted_at; type:timestamptz; default:null;"`
	UpdatedAt        time.Time  `json:"updatedAt" gorm:"column:updated_at; type:timestamptz; not null; autoUpdateTime:true;"`
	CreatedAt        time.Time  `json:"createdAt" gorm:"column:created_at; type:timestamptz; not null; autoCreateTime:true;"`

	// relations
	Blocks         []Block         `json:"block" gorm:"foreignKey:BlockGroupId; references:Id; constraint:OnUpdate:CASCADE, OnDelete:CASCADE;"`
	Owner          User            `json:"owner" gorm:"foreignKey:OwnerId; references:Id; constraint:OnUpdate:CASCADE, OnDelete:CASCADE;"`
	BlockPack      BlockPack       `json:"blockPack" gorm:"foreignKey:BlockPackId; references:Id; constraint:OnUpdate:CASCADE, OnDelete:CASCADE;"`
	PrevBlockGroup *BlockGroup     `json:"prevBlockGroup" gorm:"foreignKey:PrevBlockGroupId; references:Id; constraint:OnUpdate:CASCADE, OnDelete:CASCADE;"`
	NextBlockGroup *BlockGroup     `json:"nextBlockGroup" gorm:"foreignKey:PrevBlockGroupId; references:Id;"`
	SyncBlockGroup *SyncBlockGroup `json:"syncBlockGroup" gorm:"foreignKey:SyncBlockGroupId; references:Id; constraint:OnDelete:SET NULL;"`
}

// Root Block Table Name
func (BlockGroup) TableName() string {
	return types.TableName_BlockGroupTable.String()
}

// Root Block Table Relations
type BlockGroupRelation types.RelationName

const (
	BlockGroupRelation_Blocks         BlockGroupRelation = "Blocks"
	BlockGroupRelation_Owner          BlockGroupRelation = "Owner"
	BlockGroupRelation_BlockPack      BlockGroupRelation = "BlockPack"
	BlockGroupRelation_PrevBlockGroup BlockGroupRelation = "PrevBlockGroup"
	BlockGroupRelation_NextBlockGroup BlockGroupRelation = "NextBlockGroup"
	BlockGroupRelation_SyncBlockGroup BlockGroupRelation = "SyncBlockGroup"
)
