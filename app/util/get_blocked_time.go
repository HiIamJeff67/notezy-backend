package util

import (
	"time"

	"notezy-backend/app/exceptions"
	constants "notezy-backend/shared/constants"
)

/* ============================== Block Control of Login ============================== */

var loginCountToBlockDurationMap = map[int32]time.Duration{
	3:  5 * time.Minute,
	5:  15 * time.Minute,
	7:  30 * time.Minute,
	10: 1 * time.Hour,
	15: 6 * time.Hour,
	20: 24 * time.Hour,
	30: 7 * 24 * time.Hour,
}

func GetLoginBlockedUntilByLoginCount(loginCount int32) (*time.Time, *exceptions.Exception) {
	if loginCount < 0 {
		return nil, exceptions.Util.InvalidLoginCount(loginCount)
	}

	var blockDuration *time.Duration = nil

	for count, duration := range loginCountToBlockDurationMap {
		if loginCount >= count {
			blockDuration = &duration
		}
	}

	if blockDuration == nil {
		return nil, nil
	}

	result := time.Now().Add(*blockDuration)
	return &result, nil
}

func ShouldBlockLogin(loginCount int32) bool {
	for count := range loginCountToBlockDurationMap {
		if loginCount >= count {
			return true
		}
	}
	return false
}

func GetNextBlockThreshold(loginCount int32) int32 {
	nextThreshold := int32(constants.MAX_INT32)

	for count := range loginCountToBlockDurationMap {
		if count > loginCount && count < nextThreshold {
			nextThreshold = count
		}
	}

	if nextThreshold == constants.MAX_INT32 {
		return -1
	}
	return nextThreshold
}

/* ============================== Block Control of Auth Code ============================== */

const (
	authCodeBlockDuration = 60 * time.Second
)

func GetAuthCodeBlockUntil() time.Time {
	return time.Now().Add(authCodeBlockDuration)
}
