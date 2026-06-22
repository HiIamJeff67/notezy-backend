package schemas

import (
	"time"

	"github.com/google/uuid"

	gqlmodels "github.com/HiIamJeff67/notezy-backend/app/graphql/models"
	enums "github.com/HiIamJeff67/notezy-backend/app/models/schemas/enums"
	types "github.com/HiIamJeff67/notezy-backend/shared/types"
)

type Station struct {
	Id                  uuid.UUID            `json:"id" gorm:"column:id; type:uuid; primaryKey; default:gen_random_uuid();"`
	OwnerId             uuid.UUID            `json:"ownerId" gorm:"column:owner_id; type:uuid; not null;"`
	Name                string               `json:"name" gorm:"column:name; size:128; unique; not null; default:'undefined';"`
	Description         string               `json:"description" gorm:"column:description; size:1024; not null; default:'';"`
	Icon                *enums.SupportedIcon `json:"icon" gorm:"column:icon; type:\"SupportedIcon\"; default:null;"`
	HeaderBackgroundURL *string              `json:"headerBackgroundURL" gorm:"column:header_background_url; default:null;"`
	RoutineCount        int64                `json:"routineCount" gorm:"column:routine_count; type:bigint; not null; default:0; check:station_check_max_routine_count,routine_count <= 500;"`
	RoutineTaskCount    int64                `json:"routineTaskCount" gorm:"column:routine_task_count; type:bigint; not null; default:0; check:station_check_max_routine_task_count,routine_task_count <= 100;"`
	DeletedAt           *time.Time           `json:"deletedAt" gorm:"column:deleted_at; type:timestamptz; default:null;"`
	UpdatedAt           time.Time            `json:"updatedAt" gorm:"column:updated_at; type:timestamptz; not null; autoUpdateTime:true;"`
	CreatedAt           time.Time            `json:"createdAt" gorm:"column:created_at; type:timestamptz; not null; autoCreateTime:true;"`

	// relations
	Owner           User              `json:"owner" gorm:"foreignKey:OwnerId; references:Id; constraint:OnUpdate:CASCADE, OnDelete:CASCADE;"`
	UsersToStations []UsersToStations `json:"usersToStations" gorm:"foreignKey:StationId; references:Id; constraint:OnUpdate:CASCADE, OnDelete:CASCADE;"`
	Routines        []Routine         `json:"routines" gorm:"foreignKey:StationId; references:Id; constraint:OnUpdate:CASCADE, OnDelete:CASCADE;"`
	RoutineTasks    []RoutineTask     `json:"routineTasks" gorm:"foreignKey:StationId; references:Id; constraint:OnUpdate:CASCADE, OnDelete:CASCADE;"`
}

// Station Table Name
func (Station) TableName() string {
	return types.TableName_StationTable.String()
}

// Station Table Relations
type StationRelation types.RelationName

const (
	StationRelation_Owner           StationRelation = "Owner"
	StationRelation_UsersToStations StationRelation = "UsersToStations"
	StationRelation_Routines        StationRelation = "Routines"
	StationRelation_RoutineTasks    StationRelation = "RoutineTasks"
)

/* ============================== Relative Type Conversion ============================== */

func (s *Station) ToPrivateStation(permission enums.AccessControlPermission) *gqlmodels.PrivateStation {
	return &gqlmodels.PrivateStation{
		ID:                  s.Id,
		Permission:          permission,
		Name:                s.Name,
		Description:         s.Description,
		Icon:                s.Icon,
		HeaderBackgroundURL: s.HeaderBackgroundURL,
		RoutineCount:        s.RoutineCount,
		RoutineTaskCount:    s.RoutineTaskCount,
		DeletedAt:           s.DeletedAt,
		UpdatedAt:           s.UpdatedAt,
		CreatedAt:           s.CreatedAt,
	}
}

func (s *Station) ToPrivateSearchableStation(permission enums.AccessControlPermission) *gqlmodels.PrivateSearchableStation {
	return &gqlmodels.PrivateSearchableStation{
		ID:                  s.Id,
		Permission:          permission,
		Name:                s.Name,
		Icon:                s.Icon,
		HeaderBackgroundURL: s.HeaderBackgroundURL,
		RoutineCount:        s.RoutineCount,
		RoutineTaskCount:    s.RoutineTaskCount,
		DeletedAt:           s.DeletedAt,
		UpdatedAt:           s.UpdatedAt,
		CreatedAt:           s.CreatedAt,
	}
}
