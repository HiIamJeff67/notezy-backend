package services

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"

	coreeventscontract "github.com/HiIamJeff67/notezy-backend/contracts/core/v1/events"
	notificationscontract "github.com/HiIamJeff67/notezy-backend/contracts/notification/v1/api"
	notificationtypescontract "github.com/HiIamJeff67/notezy-backend/contracts/notification/v1/types"
	eventcontract "github.com/HiIamJeff67/notezy-backend/contracts/types/events"

	repositories "github.com/HiIamJeff67/notezy-backend/internal/notification/data/database/repositories"
	notificationexceptions "github.com/HiIamJeff67/notezy-backend/internal/notification/exceptions"
)

type NotificationService struct {
	repository repositories.NotificationRepository
	validator  *validator.Validate
}

func NewNotificationService(
	repository repositories.NotificationRepository,
	notificationValidator *validator.Validate,
) *NotificationService {
	return &NotificationService{
		repository: repository,
		validator:  notificationValidator,
	}
}

func (s *NotificationService) ConsumeRequested(
	ctx context.Context,
	event eventcontract.EventEnvelope[coreeventscontract.NotificationRequestedData],
) error {
	if event.EventType != coreeventscontract.EventType_NotificationRequested {
		return notificationexceptions.NewEventException("Notification").UnsupportedEventType()
	}
	if event.AggregateId != event.Data.RecipientUserPublicId {
		return notificationexceptions.NewEventException("Notification").AggregateRecipientMismatch()
	}
	if err := s.validator.Struct(notificationtypescontract.NotificationMetadata{
		Type:            string(event.Data.Type),
		Priority:        string(event.Data.Priority),
		TemplateVersion: event.Data.TemplateVersion,
	}); err != nil {
		return notificationexceptions.NewEventException("Notification").InvalidMetadata(err)
	}
	if event.Data.TemplateVersion != 1 {
		return notificationexceptions.NewEventException("Notification").UnsupportedTemplateVersion(
			fmt.Errorf("version: %d", event.Data.TemplateVersion),
		)
	}
	switch event.Data.Type {
	case coreeventscontract.NotificationType_News:
		if event.Data.TemplateKey != notificationtypescontract.TemplateKey_News {
			return notificationexceptions.NewEventException("Notification").InvalidNewsTemplateKey()
		}
		var payload notificationtypescontract.NewsPayload
		if err := json.Unmarshal(event.Data.Payload, &payload); err != nil {
			return notificationexceptions.NewPayloadException("Notification").PayloadDecodeFailed(err)
		}
		if err := s.validator.Struct(payload); err != nil {
			return notificationexceptions.NewPayloadException("Notification").InvalidNewsPayload(err)
		}
	case coreeventscontract.NotificationType_Warning:
		if event.Data.TemplateKey != notificationtypescontract.TemplateKey_Warning {
			return notificationexceptions.NewEventException("Notification").InvalidWarningTemplateKey()
		}
		var payload notificationtypescontract.WarningPayload
		if err := json.Unmarshal(event.Data.Payload, &payload); err != nil {
			return notificationexceptions.NewPayloadException("Notification").PayloadDecodeFailed(err)
		}
		if err := s.validator.Struct(payload); err != nil {
			return notificationexceptions.NewPayloadException("Notification").InvalidWarningPayload(err)
		}
	case coreeventscontract.NotificationType_Important:
		if event.Data.TemplateKey != notificationtypescontract.TemplateKey_Important {
			return notificationexceptions.NewEventException("Notification").InvalidImportantTemplateKey()
		}
		var payload notificationtypescontract.ImportantPayload
		if err := json.Unmarshal(event.Data.Payload, &payload); err != nil {
			return notificationexceptions.NewPayloadException("Notification").PayloadDecodeFailed(err)
		}
		if err := s.validator.Struct(payload); err != nil {
			return notificationexceptions.NewPayloadException("Notification").InvalidImportantPayload(err)
		}
	default:
		return notificationexceptions.NewEventException("Notification").UnsupportedType(
			fmt.Errorf("type: %q", event.Data.Type),
		)
	}

	if err := s.repository.CreateFromRequest(ctx, event); err != nil {
		return notificationexceptions.NewOperationException("Notification").CreateFailed(err)
	}
	return nil
}

func (s *NotificationService) List(
	ctx context.Context,
	request *notificationscontract.ListNotificationsRequestDto,
) (*notificationscontract.ListNotificationsResponseDto, error) {
	if request == nil || request.RecipientUserPublicId == uuid.Nil {
		return nil, notificationexceptions.NewRequestException("Notification").RecipientRequired()
	}
	if err := s.validator.Struct(request); err != nil {
		return nil, notificationexceptions.NewRequestException("Notification").InvalidListRequest(err)
	}
	limit := request.Limit
	if limit <= 0 || limit > 100 {
		limit = 50
	}
	notifications, err := s.repository.List(ctx, request.RecipientUserPublicId, request.Before, limit)
	if err != nil {
		return nil, notificationexceptions.NewOperationException("Notification").ListFailed(err)
	}
	response := &notificationscontract.ListNotificationsResponseDto{
		Items: make([]notificationscontract.NotificationResponseDto, len(notifications)),
	}
	for index, notification := range notifications {
		payload := map[string]any{}
		if len(notification.Payload) > 0 {
			if err := json.Unmarshal(notification.Payload, &payload); err != nil {
				return nil, notificationexceptions.NewPayloadException("Notification").ResponsePayloadDecodeFailed(err)
			}
		}
		response.Items[index] = notificationscontract.NotificationResponseDto{
			Id:                    notification.Id,
			RecipientUserPublicId: notification.RecipientUserPublicId,
			Type:                  notification.Type,
			Priority:              notification.Priority,
			TemplateKey:           notification.TemplateKey,
			TemplateVersion:       notification.TemplateVersion,
			Payload:               payload,
			CreatedAt:             notification.CreatedAt,
			ReadAt:                notification.ReadAt,
			DeletedAt:             notification.DeletedAt,
			ExpiresAt:             notification.ExpiresAt,
		}
	}
	if len(notifications) == limit && len(notifications) > 0 {
		response.NextBefore = &notifications[len(notifications)-1].CreatedAt
	}

	return response, nil
}

func (s *NotificationService) CountUnread(
	ctx context.Context,
	request *notificationscontract.CountUnreadNotificationsRequestDto,
) (*notificationscontract.CountUnreadNotificationsResponseDto, error) {
	if request == nil || request.RecipientUserPublicId == uuid.Nil {
		return nil, notificationexceptions.NewRequestException("Notification").RecipientRequired()
	}
	if err := s.validator.Struct(request); err != nil {
		return nil, notificationexceptions.NewRequestException("Notification").InvalidCountRequest(err)
	}
	count, err := s.repository.CountUnread(ctx, request.RecipientUserPublicId)
	if err != nil {
		return nil, notificationexceptions.NewOperationException("Notification").CountUnreadFailed(err)
	}

	return &notificationscontract.CountUnreadNotificationsResponseDto{Count: count}, nil
}

func (s *NotificationService) MarkRead(
	ctx context.Context,
	request *notificationscontract.MarkNotificationsReadRequestDto,
) (*notificationscontract.MarkNotificationsReadResponseDto, error) {
	if request == nil || request.RecipientUserPublicId == uuid.Nil {
		return nil, notificationexceptions.NewRequestException("Notification").RecipientRequired()
	}
	if err := s.validator.Struct(request); err != nil {
		return nil, notificationexceptions.NewRequestException("Notification").InvalidMarkReadRequest(err)
	}
	count, err := s.repository.MarkRead(ctx, request.RecipientUserPublicId, request.NotificationIds)
	if err != nil {
		return nil, notificationexceptions.NewOperationException("Notification").MarkReadFailed(err)
	}

	return &notificationscontract.MarkNotificationsReadResponseDto{UpdatedCount: count}, nil
}

func (s *NotificationService) SoftDelete(
	ctx context.Context,
	request *notificationscontract.DeleteNotificationsRequestDto,
) (*notificationscontract.DeleteNotificationsResponseDto, error) {
	if request == nil || request.RecipientUserPublicId == uuid.Nil {
		return nil, notificationexceptions.NewRequestException("Notification").RecipientRequired()
	}
	if err := s.validator.Struct(request); err != nil {
		return nil, notificationexceptions.NewRequestException("Notification").InvalidDeleteRequest(err)
	}
	count, err := s.repository.SoftDelete(ctx, request.RecipientUserPublicId, request.NotificationIds)
	if err != nil {
		return nil, notificationexceptions.NewOperationException("Notification").DeleteFailed(err)
	}

	return &notificationscontract.DeleteNotificationsResponseDto{DeletedCount: count}, nil
}

func (s *NotificationService) HardDelete(
	ctx context.Context,
	now time.Time,
	retention time.Duration,
) (int64, error) {
	count, err := s.repository.DeleteExpired(ctx, now, retention)
	if err != nil {
		return 0, notificationexceptions.NewOperationException("Notification").HardDeleteFailed(err)
	}
	return count, nil
}

func (s *NotificationService) DeleteForUser(
	ctx context.Context,
	userPublicId uuid.UUID,
) error {
	if userPublicId == uuid.Nil {
		return notificationexceptions.NewRequestException("Notification").UserRequired()
	}

	_, err := s.repository.DeleteForUser(ctx, userPublicId)
	if err != nil {
		return notificationexceptions.NewOperationException("Notification").DeleteForUserFailed(err)
	}
	return nil
}
