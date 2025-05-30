package global

import "time"

const (
	GeneralTimeoutDuration = 10 * time.Second
)

/* ============================== Auth Timeout ============================== */
const (
	RegisterTimeoutDuration = 10 * time.Second
	LoginTimeoutDuration    = 5 * time.Second
)
