package emails

import (
	"fmt"
	"time"

	exceptions "notezy-backend/app/exceptions"
	util "notezy-backend/app/util"
	types "notezy-backend/shared/types"
)

const (
	ValidationEmailSubject = "Verify Your Identity - Notezy Authentication Code"
)

var _validationEmailRenderer = &HTMLEmailRenderer{
	TemplatePath: "app/emails/templates/validation_email_template.html",
	DataMap:      map[string]any{},
}

var _validationEmailSender = &EmailSender{
	Host:     util.GetEnv("SMTP_HOST", "smtp.gmail.com"),
	Port:     util.GetIntEnv("SMTP_PORT", 587),
	UserName: util.GetEnv("NOTEZY_OFFICIAL_GMAIL", ""),
	Password: util.GetEnv("NOTEZY_OFFICIAL_GOOGLE_APPLICATION_PASSWORD", ""),
	From:     util.GetEnv("NOTEZY_OFFICIAL_NAME", "") + "<" + util.GetEnv("NOTEZY_OFFICIAL_GMAIL", "") + ">",
}

func SendValidationEmail(to string, name string, authCode string, userAgent string, expiredAt time.Time) *exceptions.Exception {
	remainingMinutes := int(time.Until(expiredAt).Minutes())

	_validationEmailRenderer.DataMap = map[string]any{
		"Name":          name,
		"Email":         to,
		"AuthCode":      authCode,
		"UserAgent":     userAgent,
		"ExpiryMinutes": remainingMinutes,
		"RequestTime":   time.Now().Format("2006-01-02 15:04:05 MST"),
	}

	fmt.Printf("Sending validation email to: %s with code: %s\n", to, authCode)

	body, exception := _validationEmailRenderer.Render()
	if exception != nil {
		return exception
	}

	exception = _validationEmailSender.Send(to, ValidationEmailSubject, body, types.ContentType_HTML)
	if exception != nil {
		return exception
	}

	fmt.Printf("Validation email sent successfully to: %s\n", to)
	return nil
}
