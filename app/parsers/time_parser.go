package parsers

import (
	"notezy-backend/app/exceptions"
	"time"
)

var timeFormates = []string{
	time.RFC3339,
	time.RFC3339Nano,
	"2006-01-02",
}

func ParseStringToTime(timeString string) (time.Time, *exceptions.Exception) {
	for _, format := range timeFormates {
		if t, err := time.Parse(format, timeString); err == nil {
			return t, nil
		}
	}

	return time.Time{}, exceptions.Parser.FailedToParseFromStringToTime(timeString)
}
