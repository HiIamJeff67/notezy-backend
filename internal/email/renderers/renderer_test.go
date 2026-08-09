package renderers

import (
	"errors"
	"os"
	"path/filepath"
	"testing"

	emailcontract "github.com/HiIamJeff67/notezy-backend/contracts/email/v1"
	exceptions "github.com/HiIamJeff67/notezy-backend/contracts/types/exceptions"
	emailconfig "github.com/HiIamJeff67/notezy-backend/internal/email/configs"
)

func TestRendererRenderAndContentType(t *testing.T) {
	tests := []struct {
		name        string
		contentType emailcontract.EmailContentType
		extension   string
		wantType    emailcontract.EmailContentType
	}{
		{
			name:        "html",
			contentType: emailcontract.EmailContentType_HTML,
			extension:   ".html",
			wantType:    emailcontract.EmailContentType_HTML,
		},
		{
			name:        "plain text",
			contentType: emailcontract.EmailContentType_PlainText,
			extension:   ".txt",
			wantType:    emailcontract.EmailContentType_PlainText,
		},
		{
			name:        "markdown",
			contentType: emailcontract.EmailContentType_Markdown,
			extension:   ".md",
			wantType:    emailcontract.EmailContentType_Markdown,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			templatePath := filepath.Join(t.TempDir(), "message"+test.extension)
			if err := os.WriteFile(templatePath, []byte("Hello, {{.Name}}!"), 0o600); err != nil {
				t.Fatal(err)
			}

			renderer, exception := NewRenderer(emailconfig.RendererConfig{
				TemplatePath: templatePath,
				ContentType:  test.contentType,
			})
			if exception != nil {
				t.Fatalf("NewRenderer() exception = %v", exception)
			}
			if renderer.ContentType() != test.wantType {
				t.Fatalf("ContentType() = %q, want %q", renderer.ContentType(), test.wantType)
			}

			body, exception := renderer.Render(map[string]any{"Name": "Notezy"})
			if exception != nil {
				t.Fatalf("Render() exception = %v", exception)
			}
			if body != "Hello, Notezy!" {
				t.Fatalf("Render() = %q, want %q", body, "Hello, Notezy!")
			}
		})
	}
}

func TestNewRendererRejectsUnsupportedContentType(t *testing.T) {
	renderer, exception := NewRenderer(emailconfig.RendererConfig{
		ContentType: emailcontract.EmailContentType("application/octet-stream"),
	})
	if renderer != nil {
		t.Fatalf("NewRenderer() renderer = %#v, want nil", renderer)
	}
	if exception == nil {
		t.Fatal("NewRenderer() exception = nil, want an exception")
	}
	var emailException *exceptions.Exception
	if !errors.As(exception, &emailException) {
		t.Fatalf("exception type = %T, want *exceptions.Exception", exception)
	}
	if emailException.Reason != "InvalidContentType" {
		t.Fatalf("exception.Reason = %q, want %q", emailException.Reason, "InvalidContentType")
	}
}

func TestRendererRejectsInvalidTemplate(t *testing.T) {
	tests := []struct {
		name       string
		config     emailconfig.RendererConfig
		wantReason string
	}{
		{
			name: "wrong extension",
			config: emailconfig.RendererConfig{
				TemplatePath: filepath.Join(t.TempDir(), "message.txt"),
				ContentType:  emailcontract.EmailContentType_HTML,
			},
			wantReason: "InvalidTemplate",
		},
		{
			name: "missing file",
			config: emailconfig.RendererConfig{
				TemplatePath: filepath.Join(t.TempDir(), "missing.html"),
				ContentType:  emailcontract.EmailContentType_HTML,
			},
			wantReason: "TemplateReadFailed",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			renderer, exception := NewRenderer(test.config)
			if exception != nil {
				t.Fatalf("NewRenderer() exception = %v", exception)
			}

			_, exception = renderer.Render(nil)
			if exception == nil {
				t.Fatal("Render() exception = nil, want an exception")
			}
			var emailException *exceptions.Exception
			if !errors.As(exception, &emailException) {
				t.Fatalf("exception type = %T, want *exceptions.Exception", exception)
			}
			if emailException.Reason != test.wantReason {
				t.Fatalf("exception.Reason = %q, want %q", emailException.Reason, test.wantReason)
			}
		})
	}
}

func TestRendererRejectsMalformedTemplate(t *testing.T) {
	templatePath := filepath.Join(t.TempDir(), "message.html")
	if err := os.WriteFile(templatePath, []byte("{{"), 0o600); err != nil {
		t.Fatal(err)
	}

	renderer, exception := NewRenderer(emailconfig.RendererConfig{
		TemplatePath: templatePath,
		ContentType:  emailcontract.EmailContentType_HTML,
	})
	if exception != nil {
		t.Fatalf("NewRenderer() exception = %v", exception)
	}

	_, exception = renderer.Render(nil)
	if exception == nil {
		t.Fatal("Render() exception = nil, want an exception")
	}
	var emailException *exceptions.Exception
	if !errors.As(exception, &emailException) {
		t.Fatalf("exception type = %T, want *exceptions.Exception", exception)
	}
	if emailException.Reason != "TemplateParseFailed" {
		t.Fatalf("exception.Reason = %q, want %q", emailException.Reason, "TemplateParseFailed")
	}
}
