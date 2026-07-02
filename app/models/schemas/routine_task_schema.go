package schemas

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
	"gorm.io/datatypes"

	gqlmodels "github.com/HiIamJeff67/notezy-backend/app/graphql/models"
	enums "github.com/HiIamJeff67/notezy-backend/app/models/schemas/enums"
	types "github.com/HiIamJeff67/notezy-backend/shared/types"
)

type RoutineTask struct {
	Id                uuid.UUID                `json:"id" gorm:"column:id; type:uuid; primaryKey; default:gen_random_uuid();"`
	StationId         uuid.UUID                `json:"stationId" gorm:"column:station_id; type:uuid; not null;"`
	Title             string                   `json:"title" gorm:"column:title; size:128; not null; default:'undefined';"`
	Purpose           enums.RoutineTaskPurpose `json:"purpose" gorm:"column:purpose; type:\"RoutineTaskPurpose\"; not null; default:'CreateBlockPack';"`
	Payload           datatypes.JSON           `json:"payload" gorm:"column:payload; type:jsonb; not null; default:'{}'; check:routine_task_check_payload_size,octet_length(payload::text) <= 16777216;"`
	CostUnit          int64                    `json:"costUnit" gorm:"column:cost_unit; type:bigint; not null; default:0; check:routine_task_check_cost_unit_non_negative,cost_unit >= 0;"`
	Priority          int32                    `json:"priority" gorm:"column:priority; type:integer; not null; default:0; check:routine_task_check_priority_validation,priority >= 0 AND priority <= 100;"`
	Status            enums.RoutineTaskStatus  `json:"status" gorm:"column:status; type:\"RoutineTaskStatus\"; not null; default:'Idle';"`
	Attempts          int32                    `json:"attempts" gorm:"column:attempts; type:integer; not null; default:0; check:routine_task_check_attempts_non_negative,attempts >= 0;"`
	MaxAttempts       int32                    `json:"maxAttempts" gorm:"column:max_attempts; type:integer; not null; default:1; check:routine_task_check_max_attempts_non_negative,max_attempts > 0;"`
	Period            *enums.RoutinePeriod     `json:"period" gorm:"column:period; type:\"RoutinePeriod\"; default:null;"`
	ScheduledAt       time.Time                `json:"scheduledAt" gorm:"column:scheduled_at; type:timestamptz; not null; default:NOW();"`
	ActualStartedAt   *time.Time               `json:"actualStartedAt" gorm:"column:actual_started_at; type:timestamptz; default:null;"`
	ActualEndedAt     *time.Time               `json:"actualEndedAt" gorm:"column:actual_ended_at; type:timestamptz; default:null;"`
	UpdatedAt         time.Time                `json:"updatedAt" gorm:"column:updated_at; type:timestamptz; not null; autoUpdateTime:true;"`
	CreatedAt         time.Time                `json:"createdAt" gorm:"column:created_at; type:timestamptz; not null; autoCreateTime:true;"`
	RecordScheduledAt time.Time                `json:"-" gorm:"column:record_scheduled_at;->;-:migration"` // to store the scheduled at column temporary while claiming the routine tasks to execute so that we can insert routine task record with this column
	RecordId          uuid.UUID                `json:"-" gorm:"column:record_id;->;-:migration"`           // to store the latest routine task record id created by the claimer

	// relations
	Station         Station             `json:"station" gorm:"foreignKey:StationId; references:Id; constraint:OnUpdate:CASCADE, OnDelete:CASCADE;"`
	RoutinesToTasks []RoutinesToTasks   `json:"routinesToTasks" gorm:"foreignKey:TaskId; references:Id; constraint:OnUpdate:CASCADE, OnDelete:CASCADE;"`
	Records         []RoutineTaskRecord `json:"records" gorm:"foreignKey:RoutineTaskId; references:Id; constraint:OnUpdate:CASCADE, OnDelete:CASCADE;"`
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
	RoutineTaskRelation_Records         RoutineTaskRelation = "Records"
)

/* ============================== Relative Type Conversion ============================== */

func (rt *RoutineTask) ToPrivateRoutineTask() *gqlmodels.PrivateRoutineTask {
	return &gqlmodels.PrivateRoutineTask{
		ID:              rt.Id,
		StationID:       rt.StationId,
		Title:           rt.Title,
		Purpose:         rt.Purpose,
		Payload:         json.RawMessage(rt.Payload),
		CostUnit:        rt.CostUnit,
		Priority:        rt.Priority,
		Status:          rt.Status,
		Attempts:        rt.Attempts,
		MaxAttempts:     rt.MaxAttempts,
		Period:          rt.Period,
		ScheduledAt:     rt.ScheduledAt,
		ActualStartedAt: rt.ActualStartedAt,
		ActualEndedAt:   rt.ActualEndedAt,
		UpdatedAt:       rt.UpdatedAt,
		CreatedAt:       rt.CreatedAt,
	}
}
