package emails

import (
	"gopkg.in/gomail.v2"

	exceptions "notezy-backend/app/exceptions"
	types "notezy-backend/app/shared/types"
)

type EmailSender struct {
	Host     string
	Port     int
	UserName string
	Password string
	From     string
}

func (s *EmailSender) Send(to string, subject string, body string, contentType types.ContentType) *exceptions.Exception {
	if !contentType.IsValidEnum() {
		return exceptions.Email.InvalidContentType(contentType)
	}

	contentTypeString := contentType.String()

	m := gomail.NewMessage()
	m.SetHeader("From", s.From)
	m.SetHeader("To", to)
	m.SetHeader("Subject", subject)
	m.SetBody(contentTypeString, body)

	d := gomail.NewDialer(s.Host, s.Port, s.UserName, s.Password)
	if err := d.DialAndSend(m); err != nil {
		return exceptions.Email.FailedToSendEmailWithSubject(subject).WithError(err)
	}
	return nil
}
