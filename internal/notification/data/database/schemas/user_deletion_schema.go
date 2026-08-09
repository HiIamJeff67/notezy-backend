package schemas

import (
	"time"

	"github.com/google/uuid"
)

type UserDeletion struct {
	UserPublicId uuid.UUID `json:"userPublicId" gorm:"column:user_public_id;type:uuid;primaryKey;"`
	DeletedAt    time.Time `json:"deletedAt" gorm:"column:deleted_at;type:timestamptz;not null;"`
}

func (UserDeletion) TableName() string {
	return "UserDeletionTable"
}
