package senders

import (
	"context"

	emaileventscontract "github.com/HiIamJeff67/notezy-backend/contracts/email/v1/events"
	exceptions "github.com/HiIamJeff67/notezy-backend/contracts/types/exceptions"

	"github.com/HiIamJeff67/notezy-backend/internal/email/renderers"
	emailtypes "github.com/HiIamJeff67/notezy-backend/internal/email/types"
)

const securityAlertEmailSubject = "Security Alert - Some Suspicious Actions Detected on Your Account"

type SecurityAlertEmailSenderInterface interface {
	Send(context.Context, emaileventscontract.SendSecurityAlertEmailRequestDto) *exceptions.Exception
	SendAsync(context.Context, emaileventscontract.SendSecurityAlertEmailRequestDto) *exceptions.Exception
}

type SecurityAlertEmailSender struct {
	renderer    renderers.RendererInterface
	enqueueFunc emailtypes.EnqueueFunc
}

func NewSecurityAlertEmailSender(renderer renderers.RendererInterface, enqueueFunc emailtypes.EnqueueFunc) SecurityAlertEmailSenderInterface {
	return &SecurityAlertEmailSender{renderer: renderer, enqueueFunc: enqueueFunc}
}

func (s *SecurityAlertEmailSender) Send(
	_ context.Context,
	request emaileventscontract.SendSecurityAlertEmailRequestDto,
) *exceptions.Exception {
	body, exception := s.renderer.Render(map[string]any{
		"UserName":         request.UserName,
		"Status":           request.Status,
		"AlertType":        request.AlertType,
		"Reason":           request.Reason,
		"TimeOfOccurrence": request.TimeOfOccurrence,
		"OtherDetails":     request.OtherDetails,
	})
	if exception != nil {
		return exception
	}

	return s.enqueueFunc(
		emailtypes.EmailObject{
			To:               request.To,
			Subject:          securityAlertEmailSubject,
			Body:             body,
			EmailContentType: s.renderer.ContentType(),
		},
		emailtypes.EmailTaskType_Security,
		5,
		3,
	)
}

func (s *SecurityAlertEmailSender) SendAsync(
	ctx context.Context,
	request emaileventscontract.SendSecurityAlertEmailRequestDto,
) *exceptions.Exception {
	return s.Send(ctx, request)
}

var _ SecurityAlertEmailSenderInterface = (*SecurityAlertEmailSender)(nil)
