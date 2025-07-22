package emails

import (
	"time"

	exceptions "notezy-backend/app/exceptions"
	types "notezy-backend/shared/types"
)

const (
	ValidationEmailSubject = "Verify Your Identity - Notezy Authentication Code"
)

var _validationEmailRenderer = &HTMLEmailRenderer{
	TemplatePath: "app/emails/templates/validation_email_template.html",
	DataMap:      map[string]any{},
}

func SyncSendValidationEmail(
	to string,
	userName string,
	authCode string,
	userAgent string,
	expiredAt time.Time,
) *exceptions.Exception {
	remainingMinutes := int(time.Until(expiredAt).Minutes())

	_validationEmailRenderer.DataMap = map[string]any{
		"UserName":      userName,
		"Email":         to,
		"AuthCode":      authCode,
		"UserAgent":     userAgent,
		"ExpiryMinutes": remainingMinutes,
		"RequestTime":   time.Now().Format("2006-01-02 15:04:05 MST"),
	}

	body, exception := _validationEmailRenderer.Render()
	if exception != nil {
		return exception
	}

	emailObject := EmailObject{
		To:          to,
		Subject:     ValidationEmailSubject,
		Body:        body,
		ContentType: types.ContentType_HTML,
	}

	exception = NotezyEmailWorkerManager.Enqueue(emailObject, EmailTaskType_Validation, 3, 2)
	if exception != nil {
		return exception
	}

	return nil
}
