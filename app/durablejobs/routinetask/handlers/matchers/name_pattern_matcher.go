package matchers

import (
	"fmt"
	"strings"
	"time"

	dtos "github.com/HiIamJeff67/notezy-backend/app/dtos"
	exceptions "github.com/HiIamJeff67/notezy-backend/app/exceptions"
	schemas "github.com/HiIamJeff67/notezy-backend/app/models/schemas"
)

type NamePatternMatcherInterface interface {
	Match(value string, pattern dtos.RoutineTaskNamePattern, task schemas.RoutineTask) (string, *exceptions.Exception)
}

type NamePatternMatcher struct{}

func NewNamePatternMatcher() NamePatternMatcherInterface {
	return NamePatternMatcher{}
}

func (m NamePatternMatcher) Match(
	value string,
	pattern dtos.RoutineTaskNamePattern,
	task schemas.RoutineTask,
) (string, *exceptions.Exception) {
	if len(pattern) == 0 {
		return value, nil
	}

	matchedValue := value
	for key, binding := range pattern {
		var resolved string
		switch binding.Source {
		case "scheduledAt":
			scheduledAt := task.RecordScheduledAt
			if scheduledAt.IsZero() {
				scheduledAt = task.ScheduledAt
			}

			if binding.Timezone != nil && *binding.Timezone != "" {
				location, err := time.LoadLocation(*binding.Timezone)
				if err != nil {
					return "", exceptions.RoutineTask.InvalidDto().WithOrigin(err)
				}
				scheduledAt = scheduledAt.In(location)
			}

			format := time.RFC3339
			if binding.Format != nil && *binding.Format != "" {
				format = *binding.Format
			}
			resolved = scheduledAt.Format(format)

		case "recordId":
			resolved = task.RecordId.String()

		case "shortRecordId":
			recordId := task.RecordId.String()
			if len(recordId) > 8 {
				recordId = recordId[:8]
			}
			resolved = recordId

		case "routineTaskId":
			resolved = task.Id.String()

		default:
			return "", exceptions.RoutineTask.InvalidDto().
				WithOrigin(fmt.Errorf("unsupported name pattern source: %s", binding.Source))
		}

		matchedValue = strings.ReplaceAll(matchedValue, "{{"+key+"}}", resolved)
	}

	return matchedValue, nil
}
