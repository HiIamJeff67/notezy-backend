package util

import "time"

var loginCountToBlockDurationMap = map[int32]time.Duration{
	3:  5 * time.Minute,
	5:  15 * time.Minute,
	7:  30 * time.Minute,
	10: 1 * time.Hour,
	15: 6 * time.Hour,
	20: 24 * time.Hour,
	30: 7 * 24 * time.Hour,
}

func GetLoginBlockedUntilByLoginCount(loginCount int32) *time.Time {
	var blockDuration time.Duration
	found := false

	for count, duration := range loginCountToBlockDurationMap {
		if loginCount >= count {
			blockDuration = duration
			found = true
		}
	}

	if !found {
		return nil
	}

	blockUntil := time.Now().Add(blockDuration)
	return &blockUntil
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
	nextThreshold := int32(1000) // 設定一個很大的預設值

	for count := range loginCountToBlockDurationMap {
		if count > loginCount && count < nextThreshold {
			nextThreshold = count
		}
	}

	if nextThreshold == 1000 {
		return -1
	}
	return nextThreshold
}
