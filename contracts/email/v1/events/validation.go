package emaileventscontract

import (
	"time"

	"github.com/google/uuid"
)

type SendValidationEmailRequestDto struct {
	RequestId  uuid.UUID `json:"requestId"`
	Operation  string    `json:"operation"`
	OccurredAt time.Time `json:"occurredAt"`
	To         string    `json:"to" validate:"required,email"`
	UserName   string    `json:"userName" validate:"required"`
	AuthCode   string    `json:"authCode" validate:"required"`
	UserAgent  string    `json:"userAgent" validate:"required"`
	ExpiredAt  time.Time `json:"expiredAt" validate:"required"`
}
