package inputs

import "github.com/google/uuid"

type BulkCheckRoutineTaskRecordPermissionInput struct {
	Id     uuid.UUID `json:"id" gorm:"column:id;"`
	UserId uuid.UUID `json:"userId" gorm:"column:user_id;"`
}
