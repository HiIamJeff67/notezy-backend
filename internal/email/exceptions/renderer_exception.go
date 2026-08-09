package exceptions

import (
	"net/http"

	exceptions "github.com/HiIamJeff67/notezy-backend/contracts/types/exceptions"
)

type RendererException struct {
	EmailException
}

func NewRendererException(domain string) RendererException {
	return RendererException{EmailException: NewEmailException(domain)}
}

func (e RendererException) InvalidContentType() *exceptions.Exception {
	return exceptions.New("InvalidContentType", e.Domain, "RenderEmail", "The email content type is invalid", http.StatusBadRequest)
}

func (e RendererException) InvalidTemplate() *exceptions.Exception {
	return exceptions.New("InvalidTemplate", e.Domain, "RenderEmail", "The email template is invalid", http.StatusBadRequest)
}

func (e RendererException) TemplateReadFailed(cause error) *exceptions.Exception {
	return exceptions.New("TemplateReadFailed", e.Domain, "RenderEmail", "Failed to read the email template", http.StatusInternalServerError, true).WithOrigin(cause)
}

func (e RendererException) TemplateParseFailed(cause error) *exceptions.Exception {
	return exceptions.New("TemplateParseFailed", e.Domain, "RenderEmail", "Failed to parse the email template", http.StatusInternalServerError, true).WithOrigin(cause)
}

func (e RendererException) TemplateRenderFailed(cause error) *exceptions.Exception {
	return exceptions.New("TemplateRenderFailed", e.Domain, "RenderEmail", "Failed to render the email template", http.StatusInternalServerError, true).WithOrigin(cause)
}
