package inputs

type CreateShelfInput struct {
	Name string `json:"name" gorm:"column:name;"`
}

type UpdateShelfInput struct {
	Name                     *string `json:"name" gorm:"column:name;"`
	EncodedStructure         *[]byte `json:"encodedStructure" gorm:"column:encoded_structure;"`
	EncodedStructureByteSize *int64  `json:"encodedStructureByteSize" gorm:"column:encoded_structure_byte_size;"`
	TotalShelfNodes          *int32  `json:"totalShelfNodes" gorm:"column:total_shelf_nodes;"`
	TotalMaterials           *int32  `json:"totalMaterials" gorm:"column:total_materials;"`
	MaxWidth                 *int32  `json:"maxWidth" gorm:"column:max_width;"`
	MaxDepth                 *int32  `json:"maxDepth" gorm:"column:max_depth;"`
}

type PartialUpdateShelfInput = PartialUpdateInput[UpdateShelfInput]
