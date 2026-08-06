package core

import (
	"context"

	emaileventscontract "github.com/HiIamJeff67/notezy-backend/contracts/email/v1/events"
	exceptions "github.com/HiIamJeff67/notezy-backend/shared/exceptions"

	emailsenders "github.com/HiIamJeff67/notezy-backend/internal/email/senders"
)

type SenderInterface interface {
	SendWelcomeEmail(context.Context, emaileventscontract.SendWelcomeEmailRequestDto) *exceptions.Exception
	SendValidationEmail(context.Context, emaileventscontract.SendValidationEmailRequestDto) *exceptions.Exception
	SendSecurityAlertEmail(context.Context, emaileventscontract.SendSecurityAlertEmailRequestDto) *exceptions.Exception
}

type Sender struct {
	welcome       emailsenders.WelcomeEmailSenderInterface
	validation    emailsenders.ValidationEmailSenderInterface
	securityAlert emailsenders.SecurityAlertEmailSenderInterface
}

func NewSender(
	welcome emailsenders.WelcomeEmailSenderInterface,
	validation emailsenders.ValidationEmailSenderInterface,
	securityAlert emailsenders.SecurityAlertEmailSenderInterface,
) SenderInterface {
	return &Sender{
		welcome:       welcome,
		validation:    validation,
		securityAlert: securityAlert,
	}
}

func (s *Sender) SendWelcomeEmail(
	ctx context.Context,
	request emaileventscontract.SendWelcomeEmailRequestDto,
) *exceptions.Exception {
	return s.welcome.Send(ctx, request)
}

func (s *Sender) SendValidationEmail(
	ctx context.Context,
	request emaileventscontract.SendValidationEmailRequestDto,
) *exceptions.Exception {
	return s.validation.Send(ctx, request)
}

func (s *Sender) SendSecurityAlertEmail(
	ctx context.Context,
	request emaileventscontract.SendSecurityAlertEmailRequestDto,
) *exceptions.Exception {
	return s.securityAlert.Send(ctx, request)
}

var _ SenderInterface = (*Sender)(nil)
