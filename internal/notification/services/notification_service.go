package services

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"

	coreeventscontract "github.com/HiIamJeff67/notezy-backend/contracts/core/v1/events"
	notificationscontract "github.com/HiIamJeff67/notezy-backend/contracts/notification/v1/api"
	notificationtypescontract "github.com/HiIamJeff67/notezy-backend/contracts/notification/v1/types"
	eventcontract "github.com/HiIamJeff67/notezy-backend/contracts/types/events"

	repositories "github.com/HiIamJeff67/notezy-backend/internal/notification/data/database/repositories"
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
		return errors.New("unsupported notification event type")
	}
	if event.AggregateId != event.Data.RecipientUserPublicId {
		return errors.New("notification aggregate does not match recipient")
	}
	if err := s.validator.Struct(notificationtypescontract.NotificationMetadata{
		Type:            string(event.Data.Type),
		Priority:        string(event.Data.Priority),
		TemplateVersion: event.Data.TemplateVersion,
	}); err != nil {
		return err
	}
	if event.Data.TemplateVersion != 1 {
		return fmt.Errorf("unsupported notification template version: %d", event.Data.TemplateVersion)
	}
	switch event.Data.Type {
	case coreeventscontract.NotificationType_News:
		if event.Data.TemplateKey != notificationtypescontract.TemplateKey_News {
			return errors.New("news notification template key is invalid")
		}
		var payload notificationtypescontract.NewsPayload
		if err := json.Unmarshal(event.Data.Payload, &payload); err != nil {
			return fmt.Errorf("decode news notification payload: %w", err)
		}
		if err := s.validator.Struct(payload); err != nil {
			return err
		}
	case coreeventscontract.NotificationType_Warning:
		if event.Data.TemplateKey != notificationtypescontract.TemplateKey_Warning {
			return errors.New("warning notification template key is invalid")
		}
		var payload notificationtypescontract.WarningPayload
		if err := json.Unmarshal(event.Data.Payload, &payload); err != nil {
			return fmt.Errorf("decode warning notification payload: %w", err)
		}
		if err := s.validator.Struct(payload); err != nil {
			return err
		}
	case coreeventscontract.NotificationType_Important:
		if event.Data.TemplateKey != notificationtypescontract.TemplateKey_Important {
			return errors.New("important notification template key is invalid")
		}
		var payload notificationtypescontract.ImportantPayload
		if err := json.Unmarshal(event.Data.Payload, &payload); err != nil {
			return fmt.Errorf("decode important notification payload: %w", err)
		}
		if err := s.validator.Struct(payload); err != nil {
			return err
		}
	default:
		return fmt.Errorf("unsupported notification type: %q", event.Data.Type)
	}

	return s.repository.CreateFromRequest(ctx, event)
}

func (s *NotificationService) List(
	ctx context.Context,
	request *notificationscontract.ListNotificationsRequestDto,
) (*notificationscontract.ListNotificationsResponseDto, error) {
	if request == nil || request.RecipientUserPublicId == uuid.Nil {
		return nil, errors.New("recipient user public ID is required")
	}
	if err := s.validator.Struct(request); err != nil {
		return nil, err
	}
	limit := request.Limit
	if limit <= 0 || limit > 100 {
		limit = 50
	}
	notifications, err := s.repository.List(ctx, request.RecipientUserPublicId, request.Before, limit)
	if err != nil {
		return nil, err
	}
	response := &notificationscontract.ListNotificationsResponseDto{
		Items: make([]notificationscontract.NotificationResponseDto, len(notifications)),
	}
	for index, notification := range notifications {
		payload := map[string]any{}
		if len(notification.Payload) > 0 {
			if err := json.Unmarshal(notification.Payload, &payload); err != nil {
				return nil, err
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
		return nil, errors.New("recipient user public ID is required")
	}
	if err := s.validator.Struct(request); err != nil {
		return nil, err
	}
	count, err := s.repository.CountUnread(ctx, request.RecipientUserPublicId)
	if err != nil {
		return nil, err
	}

	return &notificationscontract.CountUnreadNotificationsResponseDto{Count: count}, nil
}

func (s *NotificationService) MarkRead(
	ctx context.Context,
	request *notificationscontract.MarkNotificationsReadRequestDto,
) (*notificationscontract.MarkNotificationsReadResponseDto, error) {
	if request == nil || request.RecipientUserPublicId == uuid.Nil {
		return nil, errors.New("recipient user public ID is required")
	}
	if err := s.validator.Struct(request); err != nil {
		return nil, err
	}
	count, err := s.repository.MarkRead(ctx, request.RecipientUserPublicId, request.NotificationIds)
	if err != nil {
		return nil, err
	}

	return &notificationscontract.MarkNotificationsReadResponseDto{UpdatedCount: count}, nil
}

func (s *NotificationService) SoftDelete(
	ctx context.Context,
	request *notificationscontract.DeleteNotificationsRequestDto,
) (*notificationscontract.DeleteNotificationsResponseDto, error) {
	if request == nil || request.RecipientUserPublicId == uuid.Nil {
		return nil, errors.New("recipient user public ID is required")
	}
	if err := s.validator.Struct(request); err != nil {
		return nil, err
	}
	count, err := s.repository.SoftDelete(ctx, request.RecipientUserPublicId, request.NotificationIds)
	if err != nil {
		return nil, err
	}

	return &notificationscontract.DeleteNotificationsResponseDto{DeletedCount: count}, nil
}

func (s *NotificationService) HardDelete(
	ctx context.Context,
	now time.Time,
	retention time.Duration,
) (int64, error) {
	return s.repository.DeleteExpired(ctx, now, retention)
}

func (s *NotificationService) DeleteForUser(
	ctx context.Context,
	userPublicId uuid.UUID,
) error {
	if userPublicId == uuid.Nil {
		return errors.New("user public ID is required")
	}

	_, err := s.repository.DeleteForUser(ctx, userPublicId)
	return err
}
