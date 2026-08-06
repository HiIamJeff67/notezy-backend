package emaileventscontract

import (
	"time"

	"github.com/google/uuid"
)

type SendSecurityAlertEmailRequestDto struct {
	RequestId        uuid.UUID `json:"requestId"`
	Operation        string    `json:"operation"`
	OccurredAt       time.Time `json:"occurredAt"`
	To               string    `json:"to" validate:"required,email"`
	UserName         string    `json:"userName" validate:"required"`
	Status           string    `json:"status" validate:"required"`
	AlertType        string    `json:"alertType" validate:"required"`
	Reason           string    `json:"reason" validate:"required"`
	TimeOfOccurrence time.Time `json:"timeOfOccurrence" validate:"required"`
	OtherDetails     string    `json:"otherDetails"`
}
