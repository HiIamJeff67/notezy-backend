package emaileventscontract

import (
	"time"

	"github.com/google/uuid"
)

type SendWelcomeEmailRequestDto struct {
	RequestId  uuid.UUID `json:"requestId"`
	Operation  string    `json:"operation"`
	OccurredAt time.Time `json:"occurredAt"`
	To         string    `json:"to" validate:"required,email"`
	UserName   string    `json:"userName" validate:"required"`
	Status     string    `json:"status" validate:"required"`
}
