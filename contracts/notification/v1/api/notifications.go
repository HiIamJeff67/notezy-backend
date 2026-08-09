package notificationscontract

import (
	"time"

	"github.com/google/uuid"
)

const (
	ListMyNotificationsOperation        = "ListMyNotifications"
	CountMyUnreadNotificationsOperation = "CountMyUnreadNotifications"
	MarkMyNotificationsReadOperation    = "MarkMyNotificationsRead"
	DeleteMyNotificationsOperation      = "DeleteMyNotifications"
)

type NotificationResponseDto struct {
	Id                    uuid.UUID      `json:"id"`
	RecipientUserPublicId uuid.UUID      `json:"recipientUserPublicId"`
	Type                  string         `json:"type"`
	Priority              string         `json:"priority"`
	TemplateKey           string         `json:"templateKey"`
	TemplateVersion       int            `json:"templateVersion"`
	Payload               map[string]any `json:"payload"`
	CreatedAt             time.Time      `json:"createdAt"`
	ReadAt                *time.Time     `json:"readAt,omitempty"`
	DeletedAt             *time.Time     `json:"deletedAt,omitempty"`
	ExpiresAt             *time.Time     `json:"expiresAt,omitempty"`
}

type ListNotificationsRequestDto struct {
	RecipientUserPublicId uuid.UUID  `json:"recipientUserPublicId" validate:"required"`
	Before                *time.Time `json:"before,omitempty" validate:"omitempty"`
	Limit                 int        `json:"limit" validate:"omitempty,min=1,max=100"`
}

type ListNotificationsResponseDto struct {
	Items      []NotificationResponseDto `json:"items"`
	NextBefore *time.Time                `json:"nextBefore,omitempty"`
}

type CountUnreadNotificationsRequestDto struct {
	RecipientUserPublicId uuid.UUID `json:"recipientUserPublicId" validate:"required"`
}

type CountUnreadNotificationsResponseDto struct {
	Count int64 `json:"count"`
}

type MarkNotificationsReadRequestDto struct {
	RecipientUserPublicId uuid.UUID   `json:"recipientUserPublicId" validate:"required"`
	NotificationIds       []uuid.UUID `json:"notificationIds" validate:"omitempty,max=100,dive,required"`
}

type MarkNotificationsReadResponseDto struct {
	UpdatedCount int64 `json:"updatedCount"`
}

type DeleteNotificationsRequestDto struct {
	RecipientUserPublicId uuid.UUID   `json:"recipientUserPublicId" validate:"required"`
	NotificationIds       []uuid.UUID `json:"notificationIds" validate:"omitempty,max=100,dive,required"`
}

type DeleteNotificationsResponseDto struct {
	DeletedCount int64 `json:"deletedCount"`
}
