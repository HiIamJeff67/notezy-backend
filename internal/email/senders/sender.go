package senders

import (
	"context"
	"net/http"

	"gopkg.in/gomail.v2"

	exceptions "github.com/HiIamJeff67/notezy-backend/shared/exceptions"

	emailconfig "github.com/HiIamJeff67/notezy-backend/internal/email/configs"
	emailtypes "github.com/HiIamJeff67/notezy-backend/internal/email/types"
)

type EmailSenderInterface interface {
	Send(ctx context.Context, emailObject emailtypes.EmailObject) *exceptions.Exception
	SendAsync(ctx context.Context, emailObject emailtypes.EmailObject) *exceptions.Exception
}

type EmailSender struct {
	config emailconfig.SMTPConfig
}

func NewEmailSender(config emailconfig.SMTPConfig) EmailSenderInterface {
	return &EmailSender{config: config}
}

func (s *EmailSender) Send(_ context.Context, emailObject emailtypes.EmailObject) *exceptions.Exception {
	if !emailObject.EmailContentType.IsValidEnum() {
		return exceptions.New(
			"InvalidContentType",
			"Email",
			"Send",
			"The email content type is invalid",
			http.StatusBadRequest,
		)
	}

	message := gomail.NewMessage()
	message.SetHeader("From", s.config.From)
	message.SetHeader("To", emailObject.To)
	message.SetHeader("Subject", emailObject.Subject)
	message.SetBody(emailObject.EmailContentType.String(), emailObject.Body)

	dialer := gomail.NewDialer(
		s.config.Host,
		s.config.Port,
		s.config.UserName,
		s.config.Password,
	)
	if err := dialer.DialAndSend(message); err != nil {
		return exceptions.New(
			"DeliveryFailed",
			"Email",
			"Send",
			"Failed to send the email",
			http.StatusInternalServerError,
			true,
		).WithOrigin(err)
	}

	return nil
}

func (s *EmailSender) SendAsync(ctx context.Context, emailObject emailtypes.EmailObject) *exceptions.Exception {
	return s.Send(ctx, emailObject)
}
