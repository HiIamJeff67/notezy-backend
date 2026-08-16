package core

import (
	"context"

	emaileventscontract "github.com/HiIamJeff67/notegic-backend/contracts/email/v1/events"

	emailsenders "github.com/HiIamJeff67/notegic-backend/internal/email/senders"
)

type SenderInterface interface {
	SendWelcomeEmail(context.Context, emaileventscontract.SendWelcomeEmailRequestDto) error
	SendValidationEmail(context.Context, emaileventscontract.SendValidationEmailRequestDto) error
	SendSecurityAlertEmail(context.Context, emaileventscontract.SendSecurityAlertEmailRequestDto) error
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
) error {
	return s.welcome.Send(ctx, request)
}

func (s *Sender) SendValidationEmail(
	ctx context.Context,
	request emaileventscontract.SendValidationEmailRequestDto,
) error {
	return s.validation.Send(ctx, request)
}

func (s *Sender) SendSecurityAlertEmail(
	ctx context.Context,
	request emaileventscontract.SendSecurityAlertEmailRequestDto,
) error {
	return s.securityAlert.Send(ctx, request)
}

var _ SenderInterface = (*Sender)(nil)
