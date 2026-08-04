package email

import (
	"net/http"

	"gopkg.in/gomail.v2"

	emailcontract "github.com/HiIamJeff67/notezy-backend/contracts/email/v1"
	exceptions "github.com/HiIamJeff67/notezy-backend/internal/exceptions"
	emailconfig "github.com/HiIamJeff67/notezy-backend/internal/services/email/config"
)

/* ============================== Initialization & Instance ============================== */

type EmailSender struct {
	Host     string
	Port     int
	UserName string
	Password string
	From     string
}

var NotezyEmailSender *EmailSender

func Initialize(config emailconfig.SMTPConfig) {
	NotezyEmailSender = &EmailSender{
		Host:     config.Host,
		Port:     config.Port,
		UserName: config.UserName,
		Password: config.Password,
		From:     config.From,
	}
	NotezyEmailWorkerManager = NewEmailWorkerManager(16, *NotezyEmailSender)
}

func (s *EmailSender) AsyncSend(to string, subject string, body string, contentType emailcontract.EmailContentType) *exceptions.Exception {
	if !contentType.IsValidEnum() {
		return exceptions.New(
			"InvalidContentType",
			"Email",
			"AsyncSend",
			"The email content type is invalid",
			http.StatusBadRequest,
		)
	}

	contentTypeString := contentType.String()

	m := gomail.NewMessage()
	m.SetHeader("From", s.From)
	m.SetHeader("To", to)
	m.SetHeader("Subject", subject)
	m.SetBody(contentTypeString, body)

	d := gomail.NewDialer(s.Host, s.Port, s.UserName, s.Password)
	if err := d.DialAndSend(m); err != nil {
		return exceptions.New(
			"DeliveryFailed",
			"Email",
			"AsyncSend",
			"Failed to send the email",
			http.StatusInternalServerError,
			true,
		).WithOrigin(err)
	}
	return nil
}
