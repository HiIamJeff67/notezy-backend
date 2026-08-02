package emailcontract

import "time"

type SendValidationEmailRequestDto struct {
	To        string    `json:"to" validate:"required,email"`
	UserName  string    `json:"userName" validate:"required"`
	AuthCode  string    `json:"authCode" validate:"required"`
	UserAgent string    `json:"userAgent" validate:"required"`
	ExpiredAt time.Time `json:"expiredAt" validate:"required"`
}
