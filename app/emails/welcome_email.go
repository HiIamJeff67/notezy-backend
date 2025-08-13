package emails

import (
	exceptions "notezy-backend/app/exceptions"
	types "notezy-backend/shared/types"
)

const (
	WelcomeEmailSubject = "Welcome to Notezy - Thanks for the Registration"
)

var _welcomeEmailRenderer = &HTMLEmailRenderer{
	TemplatePath: "app/emails/templates/welcome_email_template.html",
	DataMap:      map[string]any{},
}

func SyncSendWelcomeEmail(
	to string,
	userName string,
	status string,
) *exceptions.Exception {
	_welcomeEmailRenderer.DataMap = map[string]any{
		"UserName": userName,
		"Email":    to,
		"Status":   status,
	}
	body, exception := _welcomeEmailRenderer.Render()
	if exception != nil {
		return exception
	}

	emailObject := EmailObject{
		To:               to,
		Subject:          WelcomeEmailSubject,
		Body:             body,
		EmailContentType: types.EmailContentType_HTML,
	}

	exception = NotezyEmailWorkerManager.Enqueue(emailObject, EmailTaskType_Welcome, 3, 1)
	if exception != nil {
		return exception
	}

	return nil
}
