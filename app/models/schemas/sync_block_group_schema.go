package schemas

import (
	"time"

	"github.com/google/uuid"

	enums "notezy-backend/app/models/schemas/enums"
	types "notezy-backend/shared/types"
)

type SyncBlockGroup struct {
	Id         uuid.UUID                     `json:"id" gorm:"column:id; type:uuid; primaryKey; not null; default:gen_random_uuid();"`
	OwnerId    uuid.UUID                     `json:"ownerId" gorm:"column:owner_id; type:uuid; not null;"`
	Permission enums.AccessControlPermission `json:"permission" gorm:"column:permission; type:\"AccessControlPermission\"; not null; default:'Read';"`
	Size       int64                         `json:"size" gorm:"column:size; type:bigint; not null; default:0;"`
	UpdatedAt  time.Time                     `json:"updatedAt" gorm:"column:updated_at; type:timestamptz; not null; autoUpdateTime:true;"`
	CreatedAt  time.Time                     `json:"createdAt" gorm:"column:created_at; type:timestamptz; not null; autoCreateTime:true;"`

	// relations
	BlockGroups []BlockGroup `json:"blockGroups" gorm:"foreignKey:SyncBlockGroupId; references:Id; constraint:OnDelete:SET NULL;"`
	SyncBlocks  []SyncBlock  `json:"syncBlocks" gorm:"foreignKey:BlockGroupId; references:Id; constraint:OnUpdate:CASCADE, OnDelete:CASCADE;"`
	Owner       User         `json:"owner" gorm:"foreignKey:OwnerId; references:Id; constraint:OnUpdate:CASCADE, OnDelete:CASCADE;"`
}

// Block Sync Group Table Name
func (SyncBlockGroup) TableName() string {
	return types.TableName_SyncBlockGroupTable.String()
}

// Block Sync Group Relations
type SyncBlockGroupRelation types.RelationName

const (
	SyncBlockGroupRelation_BlockGroups SyncBlockGroupRelation = "BlockGroups"
	SyncBlockGroupRelation_SyncBlocks  SyncBlockGroupRelation = "SyncBlocks"
	SyncBlockGroupRelation_Owner       SyncBlockGroupRelation = "Owner"
)
