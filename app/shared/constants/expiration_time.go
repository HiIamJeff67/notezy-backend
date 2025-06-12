package constants

import "time"

const (
	ExpirationTimeOfAccessToken  = 30 * time.Minute
	ExpirationTimeOfRefreshToken = 14 * 24 * time.Hour
)

const (
	ExpirationTimeOfAuthCode = 60 * time.Second
)
