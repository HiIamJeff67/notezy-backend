package constants

import "time"

const (
	ExtraCapacity = 2
)

const (
	MinIntervalTimeOfLastRequest = time.Microsecond
)

const (
	MaxDatabaseUpdateParameters = 65535
)

const (
	DefaultSearchLimit = 10
	MaxSearchLimit     = 100
)

const (
	CheckPointPerTraverse  = 1000000 // 10^6
	MaxTraverseTimeout     = 8 * time.Second
	MaxNumOfTraversedNodes = 1000000000 // 10^9
)
