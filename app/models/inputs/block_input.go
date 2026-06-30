package inputs

import (
	"github.com/google/uuid"
	"gorm.io/datatypes"

	enums "github.com/HiIamJeff67/notezy-backend/app/models/schemas/enums"
)

type CreateBlockInput struct {
	Id            uuid.UUID       `json:"id" gorm:"column:id;"`
	ParentBlockId *uuid.UUID      `json:"parentBlockId" gorm:"column:parent_block_id;"`
	Type          enums.BlockType `json:"type" gorm:"column:type;"`
	Props         datatypes.JSON  `json:"props" gorm:"column:props;"`
	Content       datatypes.JSON  `json:"content" gorm:"column:content;"`
}

type CreateBlockGroupContentInput struct {
	BlockGroupId uuid.UUID `json:"blockGroupId"`
	Blocks       []CreateBlockInput
}

type UpdateBlockInput struct {
	ParentBlockId *uuid.UUID       `json:"parentBlockId" gorm:"column:parent_block_id;"`
	BlockGroupId  *uuid.UUID       `json:"blockGroupId" gorm:"column:block_group_id;"`
	Type          *enums.BlockType `json:"type" gorm:"column:type;"`
	Props         *datatypes.JSON  `json:"props" gorm:"column:props;"`
	Content       *datatypes.JSON  `json:"content" gorm:"column:content;"`
}

type PartialUpdateBlockInput = PartialUpdateInput[UpdateBlockInput]

type UpdateBlockByIdInput struct {
	Id                 uuid.UUID                            `json:"id" gorm:"column:id;"`
	PartialUpdateInput PartialUpdateInput[UpdateBlockInput] `json:"partialUpdateInput"`
}

/* ============================== System Only Input ============================== */

type BulkCheckBlockPermissionInput struct {
	UserId uuid.UUID `json:"userId" gorm:"column:user_id;"`
	Id     uuid.UUID `json:"id" gorm:"column:id;"`
}

type BulkCreateBlockGroupContentInput struct {
	UserId       uuid.UUID `json:"userId" gorm:"column:user_id;"`
	BlockGroupId uuid.UUID `json:"blockGroupId" gorm:"column:block_group_id;"`
	Blocks       []CreateBlockInput
}

type BulkUpdateBlockInput struct {
	UserId             uuid.UUID                            `json:"userId" gorm:"column:user_id;"`
	Id                 uuid.UUID                            `json:"id" gorm:"column:id;"`
	PartialUpdateInput PartialUpdateInput[UpdateBlockInput] `json:"partialUpdateInput"`
}

type BulkDeleteBlockInput struct {
	UserId uuid.UUID `json:"userId" gorm:"column:user_id;"`
	Id     uuid.UUID `json:"id" gorm:"column:id;"`
}
