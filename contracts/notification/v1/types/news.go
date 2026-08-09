package notificationtypescontract

const TemplateKey_News = "news"

type NewsPayload struct {
	Title     string `json:"title" validate:"required,max=200,isnewscontent"`
	Summary   string `json:"summary" validate:"required,max=500,isnewscontent"`
	Body      string `json:"body" validate:"required,max=20000,isnewscontent"`
	ActionUrl string `json:"actionUrl,omitempty" validate:"omitempty,isurl"`
}
