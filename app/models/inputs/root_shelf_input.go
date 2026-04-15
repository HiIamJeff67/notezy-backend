package inputs

import (
	"time"

	"github.com/google/uuid"
)

type CreateRootShelfInput struct {
	Name           string     `json:"name" gorm:"column:name;"`
	LastAnalyzedAt *time.Time `json:"lastAnalyzedAt" gorm:"column:last_analyzed_at;"`
}

type UpdateRootShelfInput struct {
	Name           *string    `json:"name" gorm:"column:name;"`
	SubShelfCount  *int32     `json:"subShelfCount" gorm:"column:sub_shelf_count;"`
	ItemCount      *int32     `json:"itemCount" gorm:"column:item_count;"`
	LastAnalyzedAt *time.Time `json:"lastAnalyzedAt" gorm:"column:last_analyzed_at;"`
}

type PartialUpdateRootShelfInput = PartialUpdateInput[UpdateRootShelfInput]

type BulkUpdateRootShelfInput struct {
	Id                 uuid.UUID                                `json:"id" gorm:"column:id;"`
	PartialUpdateInput PartialUpdateInput[UpdateRootShelfInput] `json:"partialUpdateInput"`
}
