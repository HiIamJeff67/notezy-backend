package schemas

import (
	"time"

	enums "notezy-backend/app/models/schemas/enums"
	types "notezy-backend/shared/types"
)

// This table is only mutatable by the admin, and accessable by both client user and admin.
// To declare the value or data of this table, you MUST use the seeding method under notezy-backend/app/models/seeds/
type PlanLimitation struct {
	Key                          enums.UserPlan `json:"key" gorm:"column:key; type:UserPlan; primaryKey;"`
	MaxRootShelfCount            int32          `json:"maxRootShelfCount" gorm:"column:max_root_shelf_count; type:integer; not null;"`
	MaxBlockPackCount            int32          `json:"maxBlockPackCount" gorm:"column:max_block_pack_count; type:integer; not null;"`
	MaxBlockCount                int32          `json:"maxBlockCount" gorm:"column:max_block_count; type:integer; not null;"`
	MaxSyncBlockCount            int32          `json:"maxSyncBlockCount" gorm:"column:max_sync_block_count; type:integer; not null;"`
	MaxMaterialCount             int32          `json:"maxMaterialCount" gorm:"column:max_material_count; type:integer; not null;"`
	MaxWorkflowCount             int32          `json:"maxWorkflowCount" gorm:"column:max_work_flow_count; type:integer; not null;"`
	MaxAdditionalItemCount       int32          `json:"maxAdditionalItemCount" gorm:"column:max_additional_item_count; type:integer; not null;"`
	MaxSubShelfCountPerRootShelf int32          `json:"maxSubShelfCountPerRootShelf" gorm:"column:max_sub_shelf_count_per_root_shelf; type:integer; not null;"`
	MaxItemCountPerRootShelf     int32          `json:"maxItemCountPerRootShelf" gorm:"column:max_item_count_per_root_shelf; type:integer; not null;"`
	MaxBlockCountPerBlockPack    int32          `json:"maxBlockCountPerBlockPack" gorm:"column:max_block_count_per_block_pack; type:integer; not null;"`
	MaxMaterialMegaByteSize      float64        `json:"maxMaterialMegaByteSize" gorm:"column:max_material_mega_byte_size; type:double precision; not null;"`
	UpdatedAt                    time.Time      `json:"updatedAt" gorm:"column:updated_at; type:timestamptz; not null; autoUpdateTime:true;"`
	CreatedAt                    time.Time      `json:"createdAt" gorm:"column:created_at; type:timestamptz; not null; autoCreateTime:true;"`
}

// Plan Limitation Table Name
func (PlanLimitation) TableName() string {
	return types.TableName_PlanLimitationTable.String()
}
