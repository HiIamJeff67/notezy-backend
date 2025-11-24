package inputs

import (
	"github.com/google/uuid"
	"gorm.io/datatypes"

	enums "notezy-backend/app/models/schemas/enums"
)

type CreateBlockInput struct {
	PrevBlockId *uuid.UUID      `json:"prevBlockId" gorm:"column:prev_block_id;"`
	Type        enums.BlockType `json:"type" gorm:"column:type;"`
	Props       datatypes.JSON  `json:"props" gorm:"column:props;"`
	Content     datatypes.JSON  `json:"content" gorm:"column:content;"`
}

type UpdateBlockInput struct {
	Type    *enums.BlockType `json:"type" gorm:"column:type;"`
	Props   *datatypes.JSON  `json:"props" gorm:"column:props;"`
	Content *datatypes.JSON  `json:"content" gorm:"column:content;"`
}

type PartialUpdateBlockInput = PartialUpdateInput[UpdateBlockInput]
