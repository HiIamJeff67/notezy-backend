package core

import (
	"time"

	"github.com/go-playground/validator/v10"

	exceptions "github.com/HiIamJeff67/notezy-backend/internal/exceptions"
)

type Sender struct {
	SendWelcomeEmail       func(to string, userName string, status string) *exceptions.Exception
	SendValidationEmail    func(to string, userName string, authCode string, userAgent string, expiredAt time.Time) *exceptions.Exception
	SendSecurityAlertEmail func(to string, userName string, status string, alertType string, reason string, timeOfOccurrence time.Time, otherDetails string) *exceptions.Exception
}

type Endpoint struct {
	sender    Sender
	validator *validator.Validate
}

func NewEndpoint(sender Sender, validator *validator.Validate) Endpoint {
	return Endpoint{
		sender:    sender,
		validator: validator,
	}
}
