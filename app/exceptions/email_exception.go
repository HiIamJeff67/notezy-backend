package exceptions

import (
	"fmt"
	"net/http"

	types "notezy-backend/shared/types"
)

const (
	_ExceptionBaseCode_Email ExceptionCode = EmailExceptionSubDomainCode * ExceptionSubDomainCodeShiftAmount

	EmailExceptionSubDomainCode ExceptionCode   = 5
	ExceptionBaseCode_Email     ExceptionCode   = _ExceptionBaseCode_Email + ReservedExceptionCode
	ExceptionPrefix_Email       ExceptionPrefix = "Email"
)

type EmailExceptionDomain struct {
	BaseCode ExceptionCode
	Prefix   ExceptionPrefix
	APIExceptionDomain
}

var Email = &EmailExceptionDomain{
	BaseCode: ExceptionBaseCode_Email,
	Prefix:   ExceptionPrefix_Email,
	APIExceptionDomain: APIExceptionDomain{
		_BaseCode: _ExceptionBaseCode_Email,
		_Prefix:   ExceptionPrefix_Email,
	},
}

func (d *EmailExceptionDomain) FailedToSendEmailWithSubject(subject string) *Exception {
	return &Exception{
		Code:           d.BaseCode + 1,
		Reason:         "FailedToSendEmailWithSubject",
		Prefix:         d.Prefix,
		Message:        fmt.Sprintf("Failed to send the email with subject of %s", subject),
		HTTPStatusCode: http.StatusInternalServerError,
		LastStackFrame: &GetStackTrace(2, 1)[0],
	}
}

func (d *EmailExceptionDomain) InvalidContentType(contentType types.ContentType) *Exception {
	return &Exception{
		Code:           d.BaseCode + 2,
		Reason:         "InvalidContentType",
		Prefix:         d.Prefix,
		Message:        fmt.Sprintf("The given content type of %v is not a valid content type", contentType),
		HTTPStatusCode: http.StatusInternalServerError,
		LastStackFrame: &GetStackTrace(2, 1)[0],
	}
}

func (d *EmailExceptionDomain) FailedToReadTemplateFileWithPath(templateFilePath string) *Exception {
	return &Exception{
		Code:           d.BaseCode + 3,
		Reason:         "FailedToReadTemplateFileWithPath",
		Prefix:         d.Prefix,
		Message:        fmt.Sprintf("Failed to read the email template file from %s", templateFilePath),
		HTTPStatusCode: http.StatusInternalServerError,
		LastStackFrame: &GetStackTrace(2, 1)[0],
	}
}

func (d *EmailExceptionDomain) FailedToParseTemplateWithDataMap(dataMap map[string]any) *Exception {
	return &Exception{
		Code:           d.BaseCode + 4,
		Reason:         "FailedToParseTemplateWithDataMap",
		Prefix:         d.Prefix,
		Message:        fmt.Sprintf("Failed to parse the email template with %v", dataMap),
		HTTPStatusCode: http.StatusInternalServerError,
		LastStackFrame: &GetStackTrace(2, 1)[0],
	}
}

func (d *EmailExceptionDomain) FailedToRenderTemplate() *Exception {
	return &Exception{
		Code:           d.BaseCode + 5,
		Reason:         "FailedToRenderTemplate",
		Prefix:         d.Prefix,
		Message:        "Failed to render the template",
		HTTPStatusCode: http.StatusInternalServerError,
		LastStackFrame: &GetStackTrace(2, 1)[0],
	}
}

func (d *EmailExceptionDomain) TemplateFileTypeAndContentTypeNotMatch(templateFileType string, contentType types.ContentType) *Exception {
	return &Exception{
		Code:           d.BaseCode + 6,
		Reason:         "TemplateFileTypeAndContentTypeNotMatch",
		Prefix:         d.Prefix,
		Message:        fmt.Sprintf("The type of the template file of %s is not match with the content type of %v", templateFileType, contentType),
		HTTPStatusCode: http.StatusInternalServerError,
		LastStackFrame: &GetStackTrace(2, 1)[0],
	}
}

/* ============================== For EmailWorker Routine ============================== */

func (d *EmailExceptionDomain) FailedToSendEmailByWorkers(workerId int, numOfRetries int, maxRetries int) *Exception {
	return &Exception{
		Code:           d.BaseCode + 101,
		Reason:         "FailedToSendEmailByWorkers",
		Prefix:         d.Prefix,
		Message:        fmt.Sprintf("Worker %d failed to send email (attempt %d/%d)", workerId, numOfRetries, maxRetries),
		HTTPStatusCode: http.StatusInternalServerError,
		LastStackFrame: &GetStackTrace(2, 1)[0],
	}
}

func (d *EmailExceptionDomain) FailedToEnqueueTaskToEmailWorkerManager() *Exception {
	return &Exception{
		Code:           d.BaseCode + 102,
		Reason:         "FailedToEnqueueTaskToEmailWorkerManager",
		Prefix:         d.Prefix,
		Message:        "Failed to enqueue the given task to email worker manager",
		HTTPStatusCode: http.StatusInternalServerError,
		LastStackFrame: &GetStackTrace(2, 1)[0],
	}
}
