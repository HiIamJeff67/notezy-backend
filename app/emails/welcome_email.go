package emails

import (
	exceptions "notezy-backend/app/exceptions"
	types "notezy-backend/app/shared/types"
	util "notezy-backend/app/util"
)

const (
	WelcomeEmailSubject = "Welcome and thanks for registration in the Notezy application!"
)

var _welcomeEmailRenderer = &HTMLEmailRenderer{
	TemplatePath: "templates/welcome_email_template.html",
	DataMap:      map[string]any{},
}

var _welcomeEmailSender = &EmailSender{
	Host:     "smtp.example.com",
	Port:     587,
	UserName: util.GetEnv("NOTEZY_OFFICIAL_GMAIL", ""),
	Password: util.GetEnv("NOTEZY_OFFICIAL_GOOGLE_APPLICATION_PASSWORD", ""),
	From:     util.GetEnv("NOTEZY_OFFICIAL_GMAIL", ""),
}

func SendWelcomeEmail(to string) *exceptions.Exception {
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
