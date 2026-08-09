package notificationtypescontract

const TemplateKey_Important = "important"

type ImportantPayload struct {
	Title     string `json:"title" validate:"required,max=200,isimportantcontent"`
	Message   string `json:"message" validate:"required,max=10000,isimportantcontent"`
	ActionUrl string `json:"actionUrl,omitempty" validate:"omitempty,isurl"`
}
