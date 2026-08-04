package routinestypes

import (
	"github.com/google/uuid"

	enumcontract "github.com/HiIamJeff67/notezy-backend/contracts/types/enums"
)

type LinkedRoutineAndItem struct {
	RoutineId uuid.UUID             `json:"routineId" validate:"required"`
	ItemId    uuid.UUID             `json:"itemId" validate:"required"`
	ItemType  enumcontract.ItemType `json:"itemType" validate:"required,isitemtype"`
}
