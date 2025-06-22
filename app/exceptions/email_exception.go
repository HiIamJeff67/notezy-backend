package exceptions

import (
	"fmt"
	"net/http"

	types "notezy-backend/app/shared/types"
)

const (
	_ExceptionBaseCode_Email ExceptionCode = EmailExceptionSubDomainCode * ExceptionSubDomainCodeShiftAmount

	EmailExceptionSubDomainCode ExceptionCode   = 5
	ExceptionBaseCode_Email     ExceptionCode   = _ExceptionBaseCode_Email + ReservedExceptionCode
	ExceptionPrefix_Email       ExceptionPrefix = "Email"
)

const (
	ExceptionReason_FailedToSendEmail                      ExceptionReason = "Failed_To_Send_Email"
	ExceptionReason_InvalidContentType                     ExceptionReason = "Invalid_Content_Type"
	ExceptionReason_FailedToReadTemplateFileWithPath       ExceptionReason = "Failed_To_Read_Template_File_With_Path"
	ExceptionReason_FailedToParseTemplateWithDataMap       ExceptionReason = "Faile_To_Parse_Template_With_DataMap"
	ExceptionReason_FailedToRenderTemplate                 ExceptionReason = "Failed_To_Render_Template"
	ExceptionReason_TemplateFileTypeAndContentTypeNotMatch ExceptionReason = "Template_File_Type_And_Content_Type_Not_Match"
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
		Prefix:         d.Prefix,
		Reason:         ExceptionReason_FailedToSendEmail,
		Message:        fmt.Sprintf("Failed to send the email with subject of %s", subject),
		HTTPStatusCode: http.StatusInternalServerError,
	}
}

func (d *EmailExceptionDomain) InvalidContentType(contentType types.ContentType) *Exception {
	return &Exception{
		Code:           d.BaseCode + 2,
		Prefix:         d.Prefix,
		Reason:         ExceptionReason_InvalidContentType,
		Message:        fmt.Sprintf("The given content type of %v is not a valid content type", contentType),
		HTTPStatusCode: http.StatusInternalServerError,
	}
}

func (d *EmailExceptionDomain) FailedToReadTemplateFileWithPath(templateFilePath string) *Exception {
	return &Exception{
		Code:           d.BaseCode + 3,
		Prefix:         d.Prefix,
		Reason:         ExceptionReason_FailedToReadTemplateFileWithPath,
		Message:        fmt.Sprintf("Failed to read the email template file from %s", templateFilePath),
		HTTPStatusCode: http.StatusInternalServerError,
	}
}

func (d *EmailExceptionDomain) FailedToParseTemplateWithDataMap(dataMap map[string]any) *Exception {
	return &Exception{
		Code:           d.BaseCode + 4,
		Prefix:         d.Prefix,
		Reason:         ExceptionReason_FailedToParseTemplateWithDataMap,
		Message:        fmt.Sprintf("Failed to parse the email template with %v", dataMap),
		HTTPStatusCode: http.StatusInternalServerError,
	}
}

func (d *EmailExceptionDomain) FailedToRenderTemplate() *Exception {
	return &Exception{
		Code:           d.BaseCode + 5,
		Prefix:         d.Prefix,
		Reason:         ExceptionReason_FailedToRenderTemplate,
		Message:        "Failed to render the template",
		HTTPStatusCode: http.StatusInternalServerError,
	}
}

func (d *EmailExceptionDomain) TemplateFileTypeAndContentTypeNotMatch(templateFileType string, contentType types.ContentType) *Exception {
	return &Exception{
		Code:           d.BaseCode + 6,
		Prefix:         d.Prefix,
		Reason:         ExceptionReason_TemplateFileTypeAndContentTypeNotMatch,
		Message:        fmt.Sprintf("The type of the template file of %s is not match with the content type of %v", templateFileType, contentType),
		HTTPStatusCode: http.StatusInternalServerError,
	}
}
