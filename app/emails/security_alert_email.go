package emails

import (
	"notezy-backend/app/exceptions"
	"notezy-backend/shared/types"
	"time"
)

const (
	SecurityAlertEmailSubject = "Security Alert - Some Suspicious Actions Detected on Your Account"
)

var _securityAlertEmailRenderer = &HTMLEmailRenderer{
	TemplatePath: "app/emails/templates/security_alert_email_template.html",
	DataMap:      map[string]any{},
}

func SyncSendSecurityAlertEmail(
	to string,
	userName string,
	status string,
	alertType string,
	reason string,
	timeOfOccurrence time.Time,
	otherDetails string,
) *exceptions.Exception {
	_securityAlertEmailRenderer.DataMap = map[string]any{
		"UserName":         userName,
		"Status":           status,
		"AlertType":        alertType,
		"Reason":           reason,
		"TimeOfOccurrence": timeOfOccurrence,
		"OtherDetails":     otherDetails,
	}

	body, exception := _securityAlertEmailRenderer.Render()
	if exception != nil {
		return exception
	}

	emailObject := EmailObject{
		To:          to,
		Subject:     SecurityAlertEmailSubject,
		Body:        body,
		ContentType: types.ContentType_HTML,
	}

	exception = NotezyEmailWorkerManager.Enqueue(emailObject, EmailTaskType_Security, 3, 5)
	if exception != nil {
		return exception
	}

	return nil
}
