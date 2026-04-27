package inputs

import "github.com/google/uuid"

type CreateBlockGroupInput struct {
	BlockGroupId     *uuid.UUID `json:"blockGroupId" gorm:"column:block_group_id;"` // optional
	PrevBlockGroupId *uuid.UUID `json:"prevBlockGroupId" gorm:"column:prev_block_group_id;"`
}

type BulkCreateBlockGroupInput struct {
	BlockPackId      uuid.UUID  `json:"blockPackId" gorm:"column:block_pack_id;"`
	BlockGroupId     *uuid.UUID `json:"blockGroupId" gorm:"column:block_group_id;"` // optional
	PrevBlockGroupId *uuid.UUID `json:"prevBlockGroupId" gorm:"column:prev_block_group_id;"`
}

type UpdateBlockGroupInput struct {
	PrevBlockGroupId *uuid.UUID `json:"prevBlockGroupId" gorm:"column:prev_block_group_id;"`
	Size             *int64     `json:"size" gorm:"column:size;"`
}

type PartialUpdateBlockGroupInput = PartialUpdateInput[UpdateBlockGroupInput]
