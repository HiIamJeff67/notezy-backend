package schemas

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/datatypes"

	enums "notezy-backend/app/models/schemas/enums"
	types "notezy-backend/shared/types"
)

type RoutineTask struct {
	Id              uuid.UUID                `json:"id" gorm:"column:id; type:uuid; primaryKey; default:gen_random_uuid();"`
	StationId       uuid.UUID                `json:"stationId" gorm:"column:station_id; type:uuid; not null;"`
	Purpose         enums.RoutineTaskPurpose `json:"purpose" gorm:"column:purpose; type:\"RoutineTaskPurpose\"; not null; default:'CreateBlockPack';"`
	Payload         datatypes.JSON           `json:"payload" gorm:"column:payload; type:jsonb; not null; default:'{}'; check:routine_task_check_payload_size,octet_length(payload::text) <= 2048;"`
	Priority        int32                    `json:"priority" gorm:"column:priority; type:integer; not null; default:0; check:routine_task_check_priority_non_negative,priority >= 0;"`
	Status          enums.RoutineTaskStatus  `json:"status" gorm:"column:status; type:\"RoutineTaskStatus\"; not null; default:'Idle';"`
	Attempts        int32                    `json:"attempts" gorm:"column:attempts; type:integer; not null; default:0; check:routine_task_check_attempts_non_negative,attempts >= 0;"`
	MaxAttempts     int32                    `json:"maxAttempts" gorm:"column:max_attempts; type:integer; not null; default:1; check:routine_task_check_max_attempts_non_negative,max_attempts > 0;"`
	ScheduledAt     time.Time                `json:"scheduledAt" gorm:"column:scheduled_at; type:timestamptz; not null; default:NOW();"`
	ActualStartedAt *time.Time               `json:"actualStartedAt" gorm:"column:actual_started_at; type:timestamptz; default:null;"`
	ActualEndedAt   *time.Time               `json:"actualEndedAt" gorm:"column:actual_ended_at; type:timestamptz; default:null;"`
	UpdatedAt       time.Time                `json:"updatedAt" gorm:"column:updated_at; type:timestamptz; not null; autoUpdateTime:true;"`
	CreatedAt       time.Time                `json:"createdAt" gorm:"column:created_at; type:timestamptz; not null; autoCreateTime:true;"`

	// relations
	Station         Station           `json:"station" gorm:"foreignKey:StationId; references:Id; constraint:OnUpdate:CASCADE, OnDelete:CASCADE;"`
	RoutinesToTasks []RoutinesToTasks `json:"task" gorm:"foreignKey:TaskId; references:Id; constraint:OnUpdate:CASCADE, OnDelete:CASCADE;"`
}

// RoutineTask Table Name
func (RoutineTask) TableName() string {
	return types.TableName_RoutineTaskTable.String()
}

// RoutineTask Table Relations
type RoutineTaskRelation types.RelationName

const (
	RoutineTaskRelation_Station         RoutineTaskRelation = "Station"
	RoutineTaskRelation_RoutinesToTasks RoutineTaskRelation = "RoutinesToTasks"
)
