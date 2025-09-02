package inputs

import (
	"time"

	"github.com/google/uuid"
)

type CreateShelfInput struct {
	Id               uuid.UUID  `json:"id" gorm:"column:id;"`
	Name             string     `json:"name" gorm:"column:name;"`
	EncodedStructure []byte     `json:"encodedStructure" gorm:"column:encoded_structure;"`
	LastAnalyzedAt   *time.Time `json:"lastAnalyzedAt" gorm:"column:last_analyzed_at;"`
}

type UpdateShelfInput struct {
	Name                     *string    `json:"name" gorm:"column:name;"`
	EncodedStructure         *[]byte    `json:"encodedStructure" gorm:"column:encoded_structure;"`
	EncodedStructureByteSize *int64     `json:"encodedStructureByteSize" gorm:"column:encoded_structure_byte_size;"`
	TotalShelfNodes          *int32     `json:"totalShelfNodes" gorm:"column:total_shelf_nodes;"`
	TotalMaterials           *int32     `json:"totalMaterials" gorm:"column:total_materials;"`
	MaxWidth                 *int32     `json:"maxWidth" gorm:"column:max_width;"`
	MaxDepth                 *int32     `json:"maxDepth" gorm:"column:max_depth;"`
	LastAnalyzedAt           *time.Time `json:"lastAnalyzedAt" gorm:"column:last_analyzed_at;"`
}

type PartialUpdateShelfInput = PartialUpdateInput[UpdateShelfInput]
