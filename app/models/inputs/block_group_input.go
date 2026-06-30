package inputs

import "github.com/google/uuid"

type CreateBlockGroupInput struct {
	BlockGroupId     *uuid.UUID `json:"blockGroupId" gorm:"column:block_group_id;"` // optional
	PrevBlockGroupId *uuid.UUID `json:"prevBlockGroupId" gorm:"column:prev_block_group_id;"`
}

type CreateBlockGroupByBlockPackIdInput struct {
	BlockPackId      uuid.UUID  `json:"blockPackId" gorm:"column:block_pack_id;"`
	BlockGroupId     *uuid.UUID `json:"blockGroupId" gorm:"column:block_group_id;"` // optional
	PrevBlockGroupId *uuid.UUID `json:"prevBlockGroupId" gorm:"column:prev_block_group_id;"`
}

type UpdateBlockGroupInput struct {
	PrevBlockGroupId *uuid.UUID `json:"prevBlockGroupId" gorm:"column:prev_block_group_id;"`
	Size             *int64     `json:"size" gorm:"column:size;"`
	SizeDelta        *int64     `json:"sizeDelta"`
}

type PartialUpdateBlockGroupInput = PartialUpdateInput[UpdateBlockGroupInput]

type UpdateBlockGroupByIdInput struct {
	Id                 uuid.UUID                                 `json:"id" gorm:"column:id;"`
	PartialUpdateInput PartialUpdateInput[UpdateBlockGroupInput] `json:"partialUpdateInput"`
}

/* ============================== System Only Input ============================== */

type BulkCheckBlockGroupPermissionInput struct {
	UserId uuid.UUID `json:"userId" gorm:"column:user_id;"`
	Id     uuid.UUID `json:"id" gorm:"column:id;"`
}

type BulkCreateBlockGroupInput struct {
	UserId           uuid.UUID  `json:"userId" gorm:"column:user_id;"`
	BlockPackId      uuid.UUID  `json:"blockPackId" gorm:"column:block_pack_id;"`
	BlockGroupId     *uuid.UUID `json:"blockGroupId" gorm:"column:block_group_id;"`
	PrevBlockGroupId *uuid.UUID `json:"prevBlockGroupId" gorm:"column:prev_block_group_id;"`
}

type BulkUpdateBlockGroupInput struct {
	UserId             uuid.UUID                                 `json:"userId" gorm:"column:user_id;"`
	Id                 uuid.UUID                                 `json:"id" gorm:"column:id;"`
	PartialUpdateInput PartialUpdateInput[UpdateBlockGroupInput] `json:"partialUpdateInput"`
}

type BulkDeleteBlockGroupInput struct {
	UserId uuid.UUID `json:"userId" gorm:"column:user_id;"`
	Id     uuid.UUID `json:"id" gorm:"column:id;"`
}
