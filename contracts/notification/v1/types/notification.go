package notificationtypescontract

type NotificationMetadata struct {
	Type            string `json:"type" validate:"required,isnotificationtype"`
	Priority        string `json:"priority" validate:"required,isnotificationpriority"`
	TemplateVersion int    `json:"templateVersion" validate:"required,min=1"`
}
