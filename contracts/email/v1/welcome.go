package emailcontract

type SendWelcomeEmailRequestDto struct {
	To       string `json:"to" validate:"required,email"`
	UserName string `json:"userName" validate:"required"`
	Status   string `json:"status" validate:"required"`
}
