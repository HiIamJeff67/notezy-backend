package inputs

import (
	"github.com/google/uuid"
	"gorm.io/datatypes"

	enums "notezy-backend/app/models/schemas/enums"
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
	Type    *enums.BlockType `json:"type" gorm:"column:type;"`
	Props   *datatypes.JSON  `json:"props" gorm:"column:props;"`
	Content *datatypes.JSON  `json:"content" gorm:"column:content;"`
}

type PartialUpdateBlockInput = PartialUpdateInput[UpdateBlockInput]
