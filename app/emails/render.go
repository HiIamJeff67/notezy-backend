package emails

import (
	"bytes"
	"html/template"
	"os"
	"strings"

	exceptions "notezy-backend/app/exceptions"
	util "notezy-backend/app/util"
	types "notezy-backend/shared/types"
)

/* ==================== HTML Email Renderer ==================== */
type HTMLEmailRenderer struct {
	TemplatePath string
	DataMap      map[string]any
}

func (r *HTMLEmailRenderer) Render() (string, *exceptions.Exception) {
	if templateFileType := strings.Split(r.TemplatePath, ".")[1]; !util.IsStringIn(templateFileType, []string{"html"}) {
		return "", exceptions.Email.TemplateFileTypeAndContentTypeNotMatch(templateFileType, types.ContentType_HTML)
	}
	templateBytes, err := os.ReadFile(r.TemplatePath)
	if err != nil {
		return "", exceptions.Email.FailedToReadTemplateFileWithPath(r.TemplatePath).WithError(err)
	}

	extractedTemplate, err := template.New("email").Parse(string(templateBytes))
	if err != nil {
		return "", exceptions.Email.FailedToParseTemplateWithDataMap(r.DataMap).WithError(err)
	}

	var buffer bytes.Buffer
	if err = extractedTemplate.Execute(&buffer, r.DataMap); err != nil {
		return "", exceptions.Email.FailedToRenderTemplate().WithError(err)
	}

	return buffer.String(), nil
}

/* ==================== Plain Text Email Renderer ==================== */
type PlainTextEmailRenderer struct {
	TemplatePath string
	DataMap      map[string]any
}

func (r *PlainTextEmailRenderer) Render() (string, *exceptions.Exception) {
	if templateFileType := strings.Split(r.TemplatePath, ".")[1]; util.IsStringIn(templateFileType, []string{"txt", "log", "conf", "ini", "csv"}) {
		return "", exceptions.Email.TemplateFileTypeAndContentTypeNotMatch(templateFileType, types.ContentType_PlainText)
	}
	templateBytes, err := os.ReadFile(r.TemplatePath)
	if err != nil {
		return "", exceptions.Email.FailedToReadTemplateFileWithPath(r.TemplatePath).WithError(err)
	}

	extractedTemplate, err := template.New("email").Parse(string(templateBytes))
	if err != nil {
		return "", exceptions.Email.FailedToParseTemplateWithDataMap(r.DataMap).WithError(err)
	}

	var buffer bytes.Buffer
	if err = extractedTemplate.Execute(&buffer, r.DataMap); err != nil {
		return "", exceptions.Email.FailedToRenderTemplate().WithError(err)
	}

	return buffer.String(), nil
}

/* ==================== Markdown Email Renderer ==================== */
type MarkdownEmailRenderer struct {
	TemplatePath string
	DataMap      map[string]any
}

func (r *MarkdownEmailRenderer) Render() (string, *exceptions.Exception) {
	if templateFileType := strings.Split(r.TemplatePath, ".")[1]; util.IsStringIn(templateFileType, []string{"md"}) {
		return "", exceptions.Email.TemplateFileTypeAndContentTypeNotMatch(templateFileType, types.ContentType_PlainText)
	}
	templateBytes, err := os.ReadFile(r.TemplatePath)
	if err != nil {
		return "", exceptions.Email.FailedToReadTemplateFileWithPath(r.TemplatePath).WithError(err)
	}

	extractedTemplate, err := template.New("email").Parse(string(templateBytes))
	if err != nil {
		return "", exceptions.Email.FailedToParseTemplateWithDataMap(r.DataMap).WithError(err)
	}

	var buffer bytes.Buffer
	if err = extractedTemplate.Execute(&buffer, r.DataMap); err != nil {
		return "", exceptions.Email.FailedToRenderTemplate().WithError(err)
	}

	return buffer.String(), nil
}
