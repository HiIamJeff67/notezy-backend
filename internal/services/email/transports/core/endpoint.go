package core

import (
	"time"

	exceptions "github.com/HiIamJeff67/notezy-backend/internal/exceptions"
	validation "github.com/HiIamJeff67/notezy-backend/internal/shared/validation"
)

type Sender struct {
	SendWelcomeEmail       func(to string, userName string, status string) *exceptions.Exception
	SendValidationEmail    func(to string, userName string, authCode string, userAgent string, expiredAt time.Time) *exceptions.Exception
	SendSecurityAlertEmail func(to string, userName string, status string, alertType string, reason string, timeOfOccurrence time.Time, otherDetails string) *exceptions.Exception
}

type Endpoint struct {
	sender Sender
}

var requestValidator = validation.New()

func NewEndpoint(sender Sender) Endpoint {
	return Endpoint{sender: sender}
}
