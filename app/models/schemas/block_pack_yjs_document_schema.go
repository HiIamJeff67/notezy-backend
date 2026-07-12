package schemas

import (
	"time"

	"github.com/google/uuid"

	types "github.com/HiIamJeff67/notezy-backend/shared/types"
)

type BlockPackYjsDocument struct {
	Id                     uuid.UUID  `json:"id" gorm:"column:id; type:uuid; primaryKey; not null; default:gen_random_uuid();"`
	BlockPackId            uuid.UUID  `json:"blockPackId" gorm:"column:block_pack_id; type:uuid; not null; uniqueIndex:block_pack_yjs_document_idx_block_pack_id;"`
	Snapshot               []byte     `json:"snapshot" gorm:"column:snapshot; type:bytea; not null;"`
	StateVector            []byte     `json:"stateVector" gorm:"column:state_vector; type:bytea; not null;"`
	LastUpdateSequence     int64      `json:"lastUpdateSequence" gorm:"column:last_update_sequence; type:bigint; not null; default:0;"`
	CompactedUntilSequence int64      `json:"compactedUntilSequence" gorm:"column:compacted_until_sequence; type:bigint; not null; default:0;"`
	DeletedAt              *time.Time `json:"deletedAt" gorm:"column:deleted_at; type:timestamptz; default:null;"`
	UpdatedAt              time.Time  `json:"updatedAt" gorm:"column:updated_at; type:timestamptz; not null; autoUpdateTime:true;"`
	CreatedAt              time.Time  `json:"createdAt" gorm:"column:created_at; type:timestamptz; not null; autoCreateTime:true;"`

	// relations
	BlockPack BlockPack `json:"blockPack" gorm:"foreignKey:BlockPackId; references:Id; constraint:OnUpdate:CASCADE, OnDelete:CASCADE;"`
}

func (BlockPackYjsDocument) TableName() string {
	return types.TableName_BlockPackYjsDocumentTable.String()
}

type BlockPackYjsDocumentRelation types.RelationName

const (
	BlockPackYjsDocumentRelation_BlockPack BlockPackYjsDocumentRelation = "BlockPack"
)
