package email

import (
	"net/http"
	"os"
	"strconv"

	"gopkg.in/gomail.v2"

	emailcontract "github.com/HiIamJeff67/notezy-backend/contracts/email/v1"
	exceptions "github.com/HiIamJeff67/notezy-backend/internal/exceptions"
)

/* ============================== Initialization & Instance ============================== */

type EmailSender struct {
	Host     string
	Port     int
	UserName string
	Password string
	From     string
}

var smtpPort, _ = strconv.Atoi(os.Getenv("SMTP_PORT"))

var NotezyEmailSender = &EmailSender{
	Host:     os.Getenv("SMTP_HOST"),
	Port:     smtpPort,
	UserName: os.Getenv("NOTEZY_OFFICIAL_GMAIL"),
	Password: os.Getenv("NOTEZY_OFFICIAL_GOOGLE_APPLICATION_PASSWORD"),
	From:     os.Getenv("NOTEZY_OFFICIAL_NAME") + "<" + os.Getenv("NOTEZY_OFFICIAL_GMAIL") + ">",
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
