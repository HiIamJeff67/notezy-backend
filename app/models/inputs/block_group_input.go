package inputs

import "github.com/google/uuid"

type CreateBlockGroupInput struct {
	PrevBlockGroupId *uuid.UUID `json:"prevBlockGroupId" gorm:"column:prev_block_group_id;"`
}

type UpdateBlockGroupInput struct {
	PrevBlockGroupId *uuid.UUID `json:"prevBlockGroupId" gorm:"column:prev_block_group_id;"`
	Size             *int64     `json:"size" gorm:"column:size;"`
}

type PartialUpdateBlockGroupInput = PartialUpdateInput[UpdateBlockGroupInput]
