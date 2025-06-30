package util

import "time"

func IsTimeWithinDelta(t1, t2 time.Time, delta time.Duration) bool {
	diff := t1.Sub(t2)
	if diff < 0 {
		diff = -diff
	}
	return diff <= delta
}
