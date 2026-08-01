package email

import (
	"bytes"
	"html/template"
	"net/http"
	"os"
	"strings"

	exceptions "github.com/HiIamJeff67/notezy-backend/internal/exceptions"
	stringutil "github.com/HiIamJeff67/notezy-backend/internal/shared/lib/stringutil"
)

/* ==================== HTML Email Renderer ==================== */
type HTMLEmailRenderer struct {
	TemplatePath string
	DataMap      map[string]any
}

func (r *HTMLEmailRenderer) Render() (string, *exceptions.Exception) {
	if templateFileType := strings.Split(r.TemplatePath, ".")[1]; !stringutil.IsStringIn(templateFileType, []string{"html"}) {
		return "", exceptions.New(
			"InvalidTemplate",
			"Email",
			"RenderHTML",
			"The email template type is invalid",
			http.StatusInternalServerError,
			true,
		)
	}
	templateBytes, err := os.ReadFile(r.TemplatePath)
	if err != nil {
		return "", exceptions.New(
			"TemplateReadFailed",
			"Email",
			"RenderHTML",
			"Failed to read the email template",
			http.StatusInternalServerError,
			true,
		).WithOrigin(err)
	}

	extractedTemplate, err := template.New("email").Parse(string(templateBytes))
	if err != nil {
		return "", exceptions.New(
			"TemplateParseFailed",
			"Email",
			"RenderHTML",
			"Failed to parse the email template",
			http.StatusInternalServerError,
			true,
		).WithOrigin(err)
	}

	var buffer bytes.Buffer
	if err = extractedTemplate.Execute(&buffer, r.DataMap); err != nil {
		return "", exceptions.New(
			"TemplateRenderFailed",
			"Email",
			"RenderHTML",
			"Failed to render the email template",
			http.StatusInternalServerError,
			true,
		).WithOrigin(err)
	}

	return buffer.String(), nil
}

/* ==================== Plain Text Email Renderer ==================== */
type PlainTextEmailRenderer struct {
	TemplatePath string
	DataMap      map[string]any
}

func (r *PlainTextEmailRenderer) Render() (string, *exceptions.Exception) {
	if templateFileType := strings.Split(r.TemplatePath, ".")[1]; stringutil.IsStringIn(templateFileType, []string{"txt", "log", "conf", "ini", "csv"}) {
		return "", exceptions.New(
			"InvalidTemplate",
			"Email",
			"RenderPlainText",
			"The email template type is invalid",
			http.StatusInternalServerError,
			true,
		)
	}
	templateBytes, err := os.ReadFile(r.TemplatePath)
	if err != nil {
		return "", exceptions.New(
			"TemplateReadFailed",
			"Email",
			"RenderPlainText",
			"Failed to read the email template",
			http.StatusInternalServerError,
			true,
		).WithOrigin(err)
	}

	extractedTemplate, err := template.New("email").Parse(string(templateBytes))
	if err != nil {
		return "", exceptions.New(
			"TemplateParseFailed",
			"Email",
			"RenderPlainText",
			"Failed to parse the email template",
			http.StatusInternalServerError,
			true,
		).WithOrigin(err)
	}

	var buffer bytes.Buffer
	if err = extractedTemplate.Execute(&buffer, r.DataMap); err != nil {
		return "", exceptions.New(
			"TemplateRenderFailed",
			"Email",
			"RenderPlainText",
			"Failed to render the email template",
			http.StatusInternalServerError,
			true,
		).WithOrigin(err)
	}

	return buffer.String(), nil
}

/* ==================== Markdown Email Renderer ==================== */
type MarkdownEmailRenderer struct {
	TemplatePath string
	DataMap      map[string]any
}

func (r *MarkdownEmailRenderer) Render() (string, *exceptions.Exception) {
	if templateFileType := strings.Split(r.TemplatePath, ".")[1]; stringutil.IsStringIn(templateFileType, []string{"md"}) {
		return "", exceptions.New(
			"InvalidTemplate",
			"Email",
			"RenderMarkdown",
			"The email template type is invalid",
			http.StatusInternalServerError,
			true,
		)
	}
	templateBytes, err := os.ReadFile(r.TemplatePath)
	if err != nil {
		return "", exceptions.New(
			"TemplateReadFailed",
			"Email",
			"RenderMarkdown",
			"Failed to read the email template",
			http.StatusInternalServerError,
			true,
		).WithOrigin(err)
	}

	extractedTemplate, err := template.New("email").Parse(string(templateBytes))
	if err != nil {
		return "", exceptions.New(
			"TemplateParseFailed",
			"Email",
			"RenderMarkdown",
			"Failed to parse the email template",
			http.StatusInternalServerError,
			true,
		).WithOrigin(err)
	}

	var buffer bytes.Buffer
	if err = extractedTemplate.Execute(&buffer, r.DataMap); err != nil {
		return "", exceptions.New(
			"TemplateRenderFailed",
			"Email",
			"RenderMarkdown",
			"Failed to render the email template",
			http.StatusInternalServerError,
			true,
		).WithOrigin(err)
	}

	return buffer.String(), nil
}
