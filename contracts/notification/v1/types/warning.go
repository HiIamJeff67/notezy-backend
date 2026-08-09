package notificationtypescontract

const TemplateKey_Warning = "warning"

type WarningPayload struct {
	Title   string         `json:"title" validate:"required,max=200,iswarningcontent"`
	Message string         `json:"message" validate:"required,max=10000,iswarningcontent"`
	Details map[string]any `json:"details,omitempty"`
}
