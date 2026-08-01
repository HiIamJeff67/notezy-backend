package util

import "time"

func IsTimeWithin(t1 time.Time, t2 time.Time, delta time.Duration) bool {
	difference := t1.Sub(t2)
	if difference < 0 {
		difference = -difference
	}

	return difference <= delta
}
