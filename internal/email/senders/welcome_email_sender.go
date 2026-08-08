package senders

import (
	"context"

	emaileventscontract "github.com/HiIamJeff67/notezy-backend/contracts/email/v1/events"
	exceptions "github.com/HiIamJeff67/notezy-backend/contracts/types/exceptions"

	"github.com/HiIamJeff67/notezy-backend/internal/email/renderers"
	emailtypes "github.com/HiIamJeff67/notezy-backend/internal/email/types"
)

const welcomeEmailSubject = "Welcome to Notezy - Thanks for the Registration"

type WelcomeEmailSenderInterface interface {
	Send(context.Context, emaileventscontract.SendWelcomeEmailRequestDto) *exceptions.Exception
	SendAsync(context.Context, emaileventscontract.SendWelcomeEmailRequestDto) *exceptions.Exception
}

type WelcomeEmailSender struct {
	renderer    renderers.RendererInterface
	enqueueFunc emailtypes.EnqueueFunc
}

func NewWelcomeEmailSender(renderer renderers.RendererInterface, enqueueFunc emailtypes.EnqueueFunc) WelcomeEmailSenderInterface {
	return &WelcomeEmailSender{renderer: renderer, enqueueFunc: enqueueFunc}
}

func (s *WelcomeEmailSender) Send(
	_ context.Context,
	request emaileventscontract.SendWelcomeEmailRequestDto,
) *exceptions.Exception {
	body, exception := s.renderer.Render(map[string]any{
		"UserName": request.UserName,
		"Email":    request.To,
		"Status":   request.Status,
	})
	if exception != nil {
		return exception
	}

	return s.enqueueFunc(
		emailtypes.EmailObject{
			To:               request.To,
			Subject:          welcomeEmailSubject,
			Body:             body,
			EmailContentType: s.renderer.ContentType(),
		},
		emailtypes.EmailTaskType_Welcome,
		1,
		3,
	)
}

func (s *WelcomeEmailSender) SendAsync(
	ctx context.Context,
	request emaileventscontract.SendWelcomeEmailRequestDto,
) *exceptions.Exception {
	return s.Send(ctx, request)
}

var _ WelcomeEmailSenderInterface = (*WelcomeEmailSender)(nil)
