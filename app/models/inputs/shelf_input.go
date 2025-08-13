package inputs

type CreateShelfInput struct {
	Name string `json:"name" gorm:"column:name;"`
}

type UpdateShelfInput struct {
	Name             *string `json:"name" gorm:"column:name;"`
	EncodedStructure *[]byte `json:"encodedStructure" gorm:"column:encoded_structure;"`
}

type PartialUpdateShelfInput = PartialUpdateInput[UpdateShelfInput]
