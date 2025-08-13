package emails

import (
	"gopkg.in/gomail.v2"

	exceptions "notezy-backend/app/exceptions"
	util "notezy-backend/app/util"
	types "notezy-backend/shared/types"
)

/* ============================== Initialization & Instance ============================== */

type EmailSender struct {
	Host     string
	Port     int
	UserName string
	Password string
	From     string
}

var NotezyEmailSender = &EmailSender{
	Host:     util.GetEnv("SMTP_HOST", "smtp.gmail.com"),
	Port:     util.GetIntEnv("SMTP_PORT", 587),
	UserName: util.GetEnv("NOTEZY_OFFICIAL_GMAIL", ""),
	Password: util.GetEnv("NOTEZY_OFFICIAL_GOOGLE_APPLICATION_PASSWORD", ""),
	From:     util.GetEnv("NOTEZY_OFFICIAL_NAME", "") + "<" + util.GetEnv("NOTEZY_OFFICIAL_GMAIL", "") + ">",
}

/* ============================== Methods ============================== */

func (s *EmailSender) SyncSend(to string, subject string, body string, contentType types.EmailContentType) *exceptions.Exception {
	if !contentType.IsValidEnum() {
		return exceptions.Email.InvalidEmailContentType(string(contentType))
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
