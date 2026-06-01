package schemas

import (
	"time"

	enums "github.com/HiIamJeff67/notezy-backend/app/models/schemas/enums"
	"github.com/HiIamJeff67/notezy-backend/shared/types"

	"github.com/google/uuid"
)

type UsersToRoutineTags struct {
	UserId     uuid.UUID                     `json:"userId" gorm:"column:user_id; type:uuid; primaryKey;"`
	TagId      uuid.UUID                     `json:"tagId" gorm:"column:tag_id; type:uuid; primaryKey;"`
	Permission enums.AccessControlPermission `json:"permission" gorm:"column:permission; type:\"AccessControlPermission\"; not null; default:'Read';"`
	UpdatedAt  time.Time                     `json:"updatedAt" gorm:"column:updated_at; type:timestamptz; not null; autoUpdateTime:true;"`
	CreatedAt  time.Time                     `json:"createdAt" gorm:"column:created_at; type:timestamptz; not null; autoCreateTime:true;"`

	// relations
	User       User       `gorm:"foreignKey:UserId; reference:Id; constraint:OnUpdate:CASCADE, OnDelete:CASCADE;"`
	RoutineTag RoutineTag `gorm:"foreignKey:TagId; reference:Id; constraint:OnUpdate:CASCADE, OnDelete:CASCADE;"`
}

// UsersToRoutineTags Table Name
func (UsersToRoutineTags) TableName() string {
	return types.TableName_UsersToRoutineTagsTable.String()
}

// UsersToRoutineTags Table Relations
type UsersToRoutineTagsRelation = types.RelationName

const (
	UsersToRoutineTagsRelation_User       UsersToRoutineTagsRelation = "User"
	UsersToRoutineTagsRelation_RoutineTag UsersToRoutineTagsRelation = "RoutineTag"
)
