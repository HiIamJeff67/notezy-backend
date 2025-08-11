package util

import (
	"fmt"
	"math/rand"
	"time"
)

// Generate a snow flake id.
// This is ensured to generate unique id, since the function is executing only once at the same time.
// If there's mutiple systems executing this function at the same time, we should make different machines using different machinIDs
func GenerateSnowflakeID() string {
	timestamp := time.Now().UnixMilli()
	machineID := 1
	sequence := rand.Intn(4096)

	return fmt.Sprintf("%d%03d%04d", timestamp, machineID, sequence)
}
