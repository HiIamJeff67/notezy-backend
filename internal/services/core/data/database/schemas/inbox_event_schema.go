package schemas

import (
	"time"

	"github.com/google/uuid"
)

type InboxEvent struct {
	EventId    uuid.UUID `json:"eventId" gorm:"column:event_id; type:uuid; primaryKey;"`
	ConsumedAt time.Time `json:"consumedAt" gorm:"column:consumed_at; type:timestamptz; not null; autoCreateTime:true;"`
}

func (InboxEvent) TableName() string {
	return "InboxEventTable"
}
