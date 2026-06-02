package schemas

import (
	"time"

	"github.com/google/uuid"

	"github.com/HiIamJeff67/notezy-backend/app/models/schemas/enums"
	types "github.com/HiIamJeff67/notezy-backend/shared/types"
)

type RoutinesToItems struct {
	RoutineId uuid.UUID      `json:"routineId" gorm:"column:routine_id; type:uuid; primaryKey;"`
	ItemId    uuid.UUID      `json:"itemId" gorm:"column:item_id; type:uuid; primaryKey;"`
	ItemType  enums.ItemType `json:"itemType" gorm:"column:item_type; type:\"ItemType\"; primaryKey;"`
	CreatedAt time.Time      `json:"createdAt" gorm:"column:created_at; type:timestamptz; not null; autoCreateTime:true;"`

	// relations
	Routine Routine `json:"routine" gorm:"foreignKey:RoutineId; references:Id; constraint:OnUpdate:CASCADE, OnDelete:CASCADE;"`
	Item    Item    `json:"item" gorm:"foreignKey:ItemId,ItemType; references:Id,Type; constraint:OnUpdate:CASCADE, OnDelete:CASCADE;"`
}

// RoutinesToItems Table Name
func (RoutinesToItems) TableName() string {
	return types.TableName_RoutinesToItemsTable.String()
}

// RoutinesToItems Table Relations
type RoutinesToItemsRelation types.RelationName

const (
	RoutinesToItemsRelation_Routine RoutinesToItemsRelation = "Routine"
	RoutinesToItemsRelation_Item    RoutinesToItemsRelation = "Item"
)
