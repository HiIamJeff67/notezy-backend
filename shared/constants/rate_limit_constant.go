package constants

import "time"

const (
	RequestFrequencyExtraCapacity = 2
)

const (
	MinIntervalTimeOfLastRequest = time.Microsecond
)

const (
	SynchronizationToWindowDurationRatio = 10
	MinSynchornizationInterval           = time.Second
)
