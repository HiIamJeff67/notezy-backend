package notificationscontract

import (
	"time"

	"github.com/google/uuid"
)

const (
	SearchPrivateNotificationsOperation = "SearchPrivateNotifications"
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

type SearchPrivateNotificationsRequestDto struct {
	RecipientUserPublicId uuid.UUID `json:"recipientUserPublicId" validate:"required"`
	After                 *string   `json:"after,omitempty" validate:"omitempty"`
	First                 int       `json:"first" validate:"omitempty,min=1,max=100"`
}

type SearchPrivateNotificationsResponseDto struct {
	SearchEdges    []SearchPrivateNotificationEdge `json:"searchEdges"`
	SearchPageInfo SearchNotificationPageInfo      `json:"searchPageInfo"`
	TotalCount     int32                           `json:"totalCount"`
	SearchTime     float64                         `json:"searchTime"`
}

type SearchPrivateNotificationEdge struct {
	EncodedSearchCursor string                  `json:"encodedSearchCursor"`
	Node                NotificationResponseDto `json:"node"`
}

type SearchNotificationPageInfo struct {
	HasNextPage              bool    `json:"hasNextPage"`
	HasPreviousPage          bool    `json:"hasPreviousPage"`
	StartEncodedSearchCursor *string `json:"startEncodedSearchCursor,omitempty"`
	EndEncodedSearchCursor   *string `json:"endEncodedSearchCursor,omitempty"`
}

type SearchNotificationCursorFields struct {
	CreatedAt time.Time `json:"createdAt"`
	Id        uuid.UUID `json:"id"`
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
