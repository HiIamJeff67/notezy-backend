package senders

import (
	"context"
	"time"

	emaileventscontract "github.com/HiIamJeff67/notezy-backend/contracts/email/v1/events"
	exceptions "github.com/HiIamJeff67/notezy-backend/shared/exceptions"

	"github.com/HiIamJeff67/notezy-backend/internal/email/renderers"
	emailtypes "github.com/HiIamJeff67/notezy-backend/internal/email/types"
)

const validationEmailSubject = "Verify Your Identity - Notezy Authentication Code"

type ValidationEmailSenderInterface interface {
	Send(context.Context, emaileventscontract.SendValidationEmailRequestDto) *exceptions.Exception
	SendAsync(context.Context, emaileventscontract.SendValidationEmailRequestDto) *exceptions.Exception
}

type ValidationEmailSender struct {
	renderer    renderers.RendererInterface
	enqueueFunc emailtypes.EnqueueFunc
}

func NewValidationEmailSender(renderer renderers.RendererInterface, enqueueFunc emailtypes.EnqueueFunc) ValidationEmailSenderInterface {
	return &ValidationEmailSender{renderer: renderer, enqueueFunc: enqueueFunc}
}

func (s *ValidationEmailSender) Send(
	_ context.Context,
	request emaileventscontract.SendValidationEmailRequestDto,
) *exceptions.Exception {
	body, exception := s.renderer.Render(map[string]any{
		"UserName":      request.UserName,
		"Email":         request.To,
		"AuthCode":      request.AuthCode,
		"UserAgent":     request.UserAgent,
		"ExpiryMinutes": int(time.Until(request.ExpiredAt).Minutes()),
		"RequestTime":   time.Now().Format("2006-01-02 15:04:05 MST"),
	})
	if exception != nil {
		return exception
	}

	return s.enqueueFunc(
		emailtypes.EmailObject{
			To:               request.To,
			Subject:          validationEmailSubject,
			Body:             body,
			EmailContentType: s.renderer.ContentType(),
		},
		emailtypes.EmailTaskType_Validation,
		2,
		3,
	)
}

func (s *ValidationEmailSender) SendAsync(
	ctx context.Context,
	request emaileventscontract.SendValidationEmailRequestDto,
) *exceptions.Exception {
	return s.Send(ctx, request)
}

var _ ValidationEmailSenderInterface = (*ValidationEmailSender)(nil)
