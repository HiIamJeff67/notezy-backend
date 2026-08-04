package email

import (
	emailcontract "github.com/HiIamJeff67/notezy-backend/contracts/email/v1"
	exceptions "github.com/HiIamJeff67/notezy-backend/internal/exceptions"
)

const (
	WelcomeEmailSubject = "Welcome to Notezy - Thanks for the Registration"
)

var _welcomeEmailRenderer = &HTMLEmailRenderer{
	TemplatePath: "internal/services/email/templates/welcome_email_template.html",
	DataMap:      map[string]any{},
}

func AsyncSendWelcomeEmail(
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
		EmailContentType: emailcontract.EmailContentType_HTML,
	}

	exception = NotezyEmailWorkerManager.Enqueue(emailObject, EmailTaskType_Welcome, 3, 1)
	if exception != nil {
		return exception
	}

	return nil
}
