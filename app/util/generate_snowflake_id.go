package util

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

const (
	MaxSnowflakeSequence = 4096
)

var (
	machineID            = 1 // could be set for other machines
	sequence       int64 = 0
	lastNanosecond int64 = 0
	snowflakeMu    sync.Mutex
)

// Generate a repeatable snow flake id.
func GenerateRepeatableSnowflakeID() string {
	timestamp := time.Now().UnixMilli()
	machineID := 1
	sequence := rand.Intn(MaxSnowflakeSequence)

	return fmt.Sprintf("%d%03d%04d", timestamp, machineID, sequence)
}

// Generate a unique snow flake Id in every microseconds.
// The maximum length of a generated snow flake Id will not exceed 27 digits
func GenerateUniqueSnowflakeID() string {
	snowflakeMu.Lock()
	defer snowflakeMu.Unlock()

	nowNanosecond := time.Now().UnixNano()
	if nowNanosecond == lastNanosecond {
		sequence++
		if sequence > MaxSnowflakeSequence-1 {
			// sequence overflow, wait for the next nanosecond
			for nowNanosecond <= lastNanosecond {
				nowNanosecond = time.Now().UnixNano()
			}
			sequence = 0
			lastNanosecond = nowNanosecond
		}
	} else {
		sequence = 0
		lastNanosecond = nowNanosecond
	}

	return fmt.Sprintf("%d%03d%04d", nowNanosecond, machineID, sequence)
}
