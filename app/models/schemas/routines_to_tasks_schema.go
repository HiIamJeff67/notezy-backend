package schemas

import (
	"time"

	"github.com/google/uuid"

	types "notezy-backend/shared/types"
)

type RoutinesToTasks struct {
	RoutineId uuid.UUID `json:"routineId" gorm:"column:routine_id; type:uuid; primaryKey;"`
	TaskId    uuid.UUID `json:"taskId" gorm:"column:task_id; type:uuid; primaryKey;"`
	CreatedAt time.Time `json:"createdAt" gorm:"column:created_at; type:timestamptz; not null; autoCreateTime:true;"`

	// relations
	Routine Routine     `json:"routine" gorm:"foreignKey:RoutineId; references:Id; constraint:OnUpdate:CASCADE, OnDelete:CASCADE;"`
	Task    RoutineTask `json:"task" gorm:"foreignKey:TaskId; references:Id; constraint:OnUpdate:CASCADE, OnDelete:CASCADE;"`
}

// RoutinesToTasks Table Name
func (RoutinesToTasks) TableName() string {
	return types.TableName_RoutinesToTasksTable.String()
}

// RoutinesToTasks Table Relations
type RoutinesToTasksRelation types.RelationName

const (
	RoutinesToTasksRelation_Routine RoutinesToTasksRelation = "Routine"
	RoutinesToTasksRelation_Task    RoutinesToTasksRelation = "Task"
)
