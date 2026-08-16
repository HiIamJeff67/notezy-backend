package config

import emailcontract "github.com/HiIamJeff67/notegic-backend/contracts/email/v1"

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
			TemplatePath: "templates/welcome_email_template.html",
			ContentType:  emailcontract.EmailContentType_HTML,
		},
		Validation: RendererConfig{
			TemplatePath: "templates/validation_email_template.html",
			ContentType:  emailcontract.EmailContentType_HTML,
		},
		SecurityAlert: RendererConfig{
			TemplatePath: "templates/security_alert_email_template.html",
			ContentType:  emailcontract.EmailContentType_HTML,
		},
	}
}
