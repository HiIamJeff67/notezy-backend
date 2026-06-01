package schemas

import (
	"time"

	"github.com/google/uuid"

	enums "github.com/HiIamJeff67/notezy-backend/app/models/schemas/enums"
	types "github.com/HiIamJeff67/notezy-backend/shared/types"
)

type RoutineTag struct {
	Id        uuid.UUID            `json:"id" gorm:"column:id; type:uuid; primaryKey; default:gen_random_uuid();"`
	Name      string               `json:"name" gorm:"column:name; size: 128; not null; default:'undefined';"`
	Color     string               `json:"color" gorm:"column:color; size:7; not null; default:'#FFFFFF'; check:routine_tag_check_color_hex_code,color ~ '^#[0-9A-Fa-f]{6}$';"`
	Icon      *enums.SupportedIcon `json:"icon" gorm:"column:icon; type:\"SupportedIcon\"; default:null;"`
	UpdatedAt time.Time            `json:"updatedAt" gorm:"column:updated_at; type:timestamptz; not null; autoUpdateTime:true;"`
	CreatedAt time.Time            `json:"createdAt" gorm:"column:created_at; type:timestamptz; not null; autoCreateTime:true;"`

	// relations
	UsersToRoutineTags []UsersToRoutineTags `json:"usersToRoutineTags" gorm:"foreignKey:TagId; references:Id; constraint:OnUpdate:CASCADE, OnDelete:CASCADE;"`
	RoutinesToTags     []RoutinesToTags     `json:"routinesToTags" gorm:"foreignKey:TagId; references:Id; constraint:OnUpdate:CASCADE, OnDelete:CASCADE;"`
}

// RoutineTag Table Name
func (RoutineTag) TableName() string {
	return types.TableName_RoutineTagTable.String()
}

// RoutineTag Table Relations
type RoutineTagRelation types.RelationName

const (
	RoutineTagRelation_UsersToRoutineTags RoutineTagRelation = "UsersToRoutineTags"
	RoutineTagRelation_RoutinesToTags     RoutineTagRelation = "RoutinesToTags"
)
