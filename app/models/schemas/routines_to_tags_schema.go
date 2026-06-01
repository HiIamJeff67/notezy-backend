package schemas

import (
	"time"

	"github.com/google/uuid"

	types "github.com/HiIamJeff67/notezy-backend/shared/types"
)

type RoutinesToTags struct {
	RoutineId uuid.UUID `json:"routineId" gorm:"routine_id; type:uuid; primaryKey;"`
	TagId     uuid.UUID `json:"tagId" gorm:"tag_id; type:uuid; primaryKey;"`
	CreatedAt time.Time `json:"createdAt" gorm:"created_at; type:timestamptz; not null; autoCreateTime:true;"`

	// relations
	Routine Routine    `json:"routine" gorm:"foreignKey:RoutineId; reference:Id; constraint:OnUpdate:CASCADE, OnDelete:CASCADE;"`
	Tag     RoutineTag `json:"tag" gorm:"foreignKey:TagId; references:Id; constraint:OnUpdate:CASCADE, OnDelete:CASCADE;"`
}

// RoutinesToTags Table Name
func (RoutinesToTags) TableName() string {
	return types.TableName_RoutinesToTagsTable.String()
}

// RoutinesToTags Table Relations
type RoutinesToTagsRelation types.RelationName

const (
	RoutinesToTagsRelation_Routine RoutinesToTagsRelation = "Routine"
	RoutinesToTagsRelation_Tag     RoutinesToTagsRelation = "Tag"
)
