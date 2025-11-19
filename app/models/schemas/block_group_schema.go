package schemas

import (
	"time"

	"github.com/google/uuid"

	enums "notezy-backend/app/models/schemas/enums"
	types "notezy-backend/shared/types"
)

type BlockGroup struct {
	Id               uuid.UUID                     `json:"id" gorm:"column:id; type:uuid; primaryKey; not null; default:gen_random_uuid();"`
	BlockPackId      uuid.UUID                     `json:"blockPackId" gorm:"column:block_pack_id; type:uuid; not null;"`
	PrevBlockGroupId *uuid.UUID                    `json:"prevBlockGroupId" gorm:"column:prev_block_group_id; type:uuid; default:null;"`
	SyncBlockGroupId *uuid.UUID                    `json:"syncBlockGroupId" gorm:"sync_block_group_id; type:uuid; default:null;"`
	OwnerId          uuid.UUID                     `json:"ownerId" gorm:"column:owner_id; type:uuid; not null;"`
	Permission       enums.AccessControlPermission `json:"permission" gorm:"column:permission; type:AccessControllPermission; not null; default:'Read';"`
	Size             int64                         `json:"size" gorm:"column:size; type:bigint; not null; default:0;"`
	DeletedAt        *time.Time                    `json:"deletedAt" gorm:"column:deleted_at; type:timestamptz; default:null;"`
	UpdatedAt        time.Time                     `json:"updatedAt" gorm:"column:updated_at; type:timestamptz; not null; autoUpdateTime:true;"`
	CreatedAt        time.Time                     `json:"createdAt" gorm:"column:created_at; type:timestamptz; not null; autoCreateTime:true;"`

	// relations
	Blocks         []Block         `json:"block" gorm:"foreignKey:BlockGroupId; reference:Id; constraint:OnUpdate:CASCADE, OnDelete:CASCADE;"`
	BlockPack      BlockPack       `json:"blockPack" gorm:"foreignKey:BlockPackId; reference:Id; constraint:OnUpdate:CASCADE, OnDelete:CASCADE;"`
	NextBlockGroup *BlockGroup     `json:"nextBlockGroup" gorm:"foreignKey:PrevBlockGroupId; reference:Id;"`
	SyncBlockGroup *SyncBlockGroup `json:"syncBlockGroup" gorm:"foreignKey:SyncBlockGroupId; reference:Id; constraint:OnDelete:SET NULL;"`
	Owner          User            `json:"owner" gorm:"foreignKey:OwnerId; reference:Id; constraint:OnUpdate:CASCADE, OnDelete:CASCADE;"`
}

// Root Block Table Name
func (BlockGroup) TableName() string {
	return types.ValidTableName_BlockGroupTable.String()
}

// Root Block Table Relations
type BlockGroupRelation types.ValidTableName

const (
	BlockGroupRelation_Blocks         BlockGroupRelation = "Blocks"
	BlockGroupRelation_BlockPack      BlockGroupRelation = "BlockPack"
	BlockGroupRelation_NextBlockGroup BlockGroupRelation = "NextBlockGroup"
	BlockGroupRelation_SyncBlockGroup BlockGroupRelation = "SyncBlockGroup"
	BlockGroupRelation_Owner          BlockGroupRelation = "Owner"
)
