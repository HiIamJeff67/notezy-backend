package inputs

import (
	"time"
)

type CreateRootShelfInput struct {
	Name           string     `json:"name" gorm:"column:name;"`
	LastAnalyzedAt *time.Time `json:"lastAnalyzedAt" gorm:"column:last_analyzed_at;"`
}

type UpdateRootShelfInput struct {
	Name           *string    `json:"name" gorm:"column:name;"`
	SubShelfCount  *int32     `json:"subShelfCount" gorm:"column:total_shelf_nodes;"`
	TotalMaterials *int32     `json:"totalMaterials" gorm:"column:total_materials;"`
	LastAnalyzedAt *time.Time `json:"lastAnalyzedAt" gorm:"column:last_analyzed_at;"`
}

type PartialUpdateRootShelfInput = PartialUpdateInput[UpdateRootShelfInput]
