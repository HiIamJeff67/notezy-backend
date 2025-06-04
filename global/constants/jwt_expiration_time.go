package constants

import "time"

const (
	AccessTokenExpirationTime  = 30 * time.Minute
	RefreshTokenExpirationTime = 14 * 24 * time.Hour
)
