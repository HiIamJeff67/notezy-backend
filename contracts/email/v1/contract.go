package emailcontract

import "time"

type SendEmailResponseDto struct {
	QueuedAt time.Time `json:"queuedAt"`
}
