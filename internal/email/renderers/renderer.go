package renderers

import (
	"bytes"
	"html/template"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	emailcontract "github.com/HiIamJeff67/notezy-backend/contracts/email/v1"
	exceptions "github.com/HiIamJeff67/notezy-backend/shared/exceptions"

	emailconfig "github.com/HiIamJeff67/notezy-backend/internal/email/configs"
)

type RendererInterface interface {
	Render(data map[string]any) (string, *exceptions.Exception)
	ContentType() emailcontract.EmailContentType
}

type Renderer struct {
	config            emailconfig.RendererConfig
	expectedExtension string
	renderOperation   string
}

func NewRenderer(config emailconfig.RendererConfig) (RendererInterface, *exceptions.Exception) {
	switch config.ContentType {
	case emailcontract.EmailContentType_HTML:
		return &HTMLEmailRenderer{Renderer: newRenderer(config, "html", "RenderHTML")}, nil
	case emailcontract.EmailContentType_PlainText:
		return &PlainTextEmailRenderer{Renderer: newRenderer(config, "txt", "RenderPlainText")}, nil
	case emailcontract.EmailContentType_Markdown:
		return &MarkdownEmailRenderer{Renderer: newRenderer(config, "md", "RenderMarkdown")}, nil
	default:
		return nil, exceptions.New(
			"InvalidContentType",
			"Email",
			"CreateRenderer",
			"The email content type is invalid",
			http.StatusInternalServerError,
		)
	}
}

func newRenderer(
	config emailconfig.RendererConfig,
	expectedExtension string,
	operation string,
) Renderer {
	return Renderer{
		config:            config,
		expectedExtension: expectedExtension,
		renderOperation:   operation,
	}
}

func (r *Renderer) Render(data map[string]any) (string, *exceptions.Exception) {
	return renderTemplate(r.config, r.expectedExtension, data, r.renderOperation)
}

func (r *Renderer) ContentType() emailcontract.EmailContentType {
	return r.config.ContentType
}

func renderTemplate(
	config emailconfig.RendererConfig,
	expectedExtension string,
	data map[string]any,
	operation string,
) (string, *exceptions.Exception) {
	if filepath.Ext(config.TemplatePath) != "."+expectedExtension {
		return "", exceptions.New(
			"InvalidTemplate",
			"Email",
			operation,
			"The email template type is invalid",
			http.StatusInternalServerError,
		)
	}

	templateBytes, err := os.ReadFile(config.TemplatePath)
	if err != nil {
		return "", exceptions.New(
			"TemplateReadFailed",
			"Email",
			operation,
			"Failed to read the email template",
			http.StatusInternalServerError,
			true,
		).WithOrigin(err)
	}

	extractedTemplate, err := template.New(strings.TrimSuffix(filepath.Base(config.TemplatePath), filepath.Ext(config.TemplatePath))).Parse(string(templateBytes))
	if err != nil {
		return "", exceptions.New(
			"TemplateParseFailed",
			"Email",
			operation,
			"Failed to parse the email template",
			http.StatusInternalServerError,
			true,
		).WithOrigin(err)
	}

	var buffer bytes.Buffer
	if err := extractedTemplate.Execute(&buffer, data); err != nil {
		return "", exceptions.New(
			"TemplateRenderFailed",
			"Email",
			operation,
			"Failed to render the email template",
			http.StatusInternalServerError,
			true,
		).WithOrigin(err)
	}

	body := buffer.String()
	return body, nil
}
