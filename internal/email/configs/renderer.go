package config

import emailcontract "github.com/HiIamJeff67/notezy-backend/contracts/email/v1"

type RendererConfig struct {
	TemplatePath string
	ContentType  emailcontract.EmailContentType
}

type RendererConfigs struct {
	Welcome       RendererConfig
	Validation    RendererConfig
	SecurityAlert RendererConfig
}

func loadRendererConfigs() RendererConfigs {
	return RendererConfigs{
		Welcome: RendererConfig{
			TemplatePath: "internal/email/templates/welcome_email_template.html",
			ContentType:  emailcontract.EmailContentType_HTML,
		},
		Validation: RendererConfig{
			TemplatePath: "internal/email/templates/validation_email_template.html",
			ContentType:  emailcontract.EmailContentType_HTML,
		},
		SecurityAlert: RendererConfig{
			TemplatePath: "internal/email/templates/security_alert_email_template.html",
			ContentType:  emailcontract.EmailContentType_HTML,
		},
	}
}
