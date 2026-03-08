package constants

import "time"

/* ============================== AccessToken and RefreshToken relative constants ============================== */

const (
	ExpirationTimeOfAccessToken  = 30 * time.Minute
	ExpirationTimeOfRefreshToken = 14 * 24 * time.Hour
)

/* ============================== AuthCode relative constants ============================== */

const (
	ExpirationTimeOfAuthCode = 3 * time.Minute
)

const (
	MaxLengthOfAuthCode int = 6
	MaxAuthCode         int = 999999
)
