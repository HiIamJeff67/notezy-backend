package notificationeventscontract

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

type NotificationCreatedData struct {
	NotificationId        uuid.UUID       `json:"notificationId"`
	RecipientUserPublicId uuid.UUID       `json:"recipientUserPublicId"`
	Type                  string          `json:"type"`
	Priority              string          `json:"priority"`
	TemplateKey           string          `json:"templateKey"`
	TemplateVersion       int             `json:"templateVersion"`
	Payload               json.RawMessage `json:"payload"`
	CreatedAt             time.Time       `json:"createdAt"`
	ExpiresAt             *time.Time      `json:"expiresAt,omitempty"`
}
