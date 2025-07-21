package emails

import (
	exceptions "notezy-backend/app/exceptions"
	types "notezy-backend/shared/types"
)

const (
	WelcomeEmailSubject = "Welcome to Notezy - The Account Registration Was Successfully Done"
)

var _welcomeEmailRenderer = &HTMLEmailRenderer{
	TemplatePath: "app/emails/templates/welcome_email_template.html",
	DataMap:      map[string]any{},
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

	emailObject := EmailObject{
		To:          to,
		Subject:     WelcomeEmailSubject,
		Body:        body,
		ContentType: types.ContentType_HTML,
	}

	exception = CommonEmailWorkerManager.Enqueue(emailObject, EmailTaskType_Welcome, 3, 1)
	if exception != nil {
		return exception
	}

	return nil
}
