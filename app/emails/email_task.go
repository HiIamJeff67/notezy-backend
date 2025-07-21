package emails

import (
	"time"

	types "notezy-backend/shared/types"
)

type EmailTaskType string

const (
	EmailTaskType_Undefined  EmailTaskType = "Undefined"
	EmailTaskType_Welcome    EmailTaskType = "EmailTaskType_Welcome"
	EmailTaskType_Validation EmailTaskType = "EmailTaskType_Validation"
	EmailTaskType_Security   EmailTaskType = "EmailTaskType_Security"
	EmailTaskType_News       EmailTaskType = "EmailTaskType_News"
)

type EmailObject struct {
	To          string `json:"to"`
	Subject     string `json:"subject"`
	Body        string `json:"body"`
	ContentType types.ContentType
}

type EmailTask struct {
	ID         string        `json:"id"`
	Type       EmailTaskType `json:"type"`
	Object     EmailObject   `json:"object"`
	CreatedAt  time.Time     `json:"createdAt"`
	Retries    int           `json:"retries"`
	MaxRetries int           `json:"maxRetries"`
	Priority   int           `json:"priority"` // the higher priotiy, the much more urgent
}
