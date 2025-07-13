package util

import (
	"fmt"
	"math/rand"
	"time"
)

func GenerateSnowflakeID() string {
	timestamp := time.Now().UnixMilli()
	machineID := 1
	sequence := rand.Intn(4096)

	return fmt.Sprintf("%d%03d%04d", timestamp, machineID, sequence)
}
