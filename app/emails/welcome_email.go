package emails

import (
	exceptions "notezy-backend/app/exceptions"
	util "notezy-backend/app/util"
	types "notezy-backend/shared/types"
)

const (
	WelcomeEmailSubject = "Welcome to Notezy - The Account Registration Was Successfully Done"
)

var _welcomeEmailRenderer = &HTMLEmailRenderer{
	TemplatePath: "app/emails/templates/welcome_email_template.html",
	DataMap:      map[string]any{},
}

var _welcomeEmailSender = &EmailSender{
	Host:     util.GetEnv("SMTP_HOST", "smtp.gmail.com"),
	Port:     util.GetIntEnv("SMTP_PORT", 587),
	UserName: util.GetEnv("NOTEZY_OFFICIAL_GMAIL", ""),
	Password: util.GetEnv("NOTEZY_OFFICIAL_GOOGLE_APPLICATION_PASSWORD", ""),
	From:     util.GetEnv("NOTEZY_OFFICIAL_NAME", "") + "<" + util.GetEnv("NOTEZY_OFFICIAL_GMAIL", "") + ">",
}

func SendWelcomeEmail(to string, name string, status string) *exceptions.Exception {
	_welcomeEmailRenderer.DataMap = map[string]any{
		"UserName": name,
		"Email":    to,
		"Status":   status,
	}
	body, exception := _welcomeEmailRenderer.Render()
	if exception != nil {
		return exception
	}

	exception = _welcomeEmailSender.Send(to, WelcomeEmailSubject, body, types.ContentType_HTML)
	if exception != nil {
		return exception
	}

	return nil
}
