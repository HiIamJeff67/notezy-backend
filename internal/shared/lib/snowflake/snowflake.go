package snowflake

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

var (
	sequence       int64
	lastNanosecond int64
	mu             sync.Mutex
)

func GenerateRepeatableID() string {
	timestamp := time.Now().UnixMilli()
	sequence := rand.Intn(maxSequence)

	return fmt.Sprintf("%d%03d%04d", timestamp, machineID, sequence)
}

func GenerateUniqueID() string {
	mu.Lock()
	defer mu.Unlock()

	nowNanosecond := time.Now().UnixNano()
	if nowNanosecond == lastNanosecond {
		sequence++
		if sequence > maxSequence-1 {
			for nowNanosecond <= lastNanosecond {
				nowNanosecond = time.Now().UnixNano()
			}
			sequence = 0
		}
	}

	lastNanosecond = nowNanosecond
	return fmt.Sprintf("%d%03d%04d", nowNanosecond, machineID, sequence)
}
