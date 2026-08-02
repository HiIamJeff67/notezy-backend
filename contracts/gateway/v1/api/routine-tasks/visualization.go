package routinetasksdto

import (
	"encoding/json"
	"time"

	apiv1 "github.com/HiIamJeff67/notezy-backend/contracts/gateway/v1/api"
	sharedtypes "github.com/HiIamJeff67/notezy-backend/internal/shared/types"
)

type RoutineTaskCountDatum struct {
	Id    string          `json:"id"`
	X     string          `json:"x"`
	Value int64           `json:"value"`
	Meta  json.RawMessage `json:"meta"`
}
type RoutineTaskCountResponseDto struct {
	Data []RoutineTaskCountDatum `json:"data"`
}
type VisualizeMyRoutineTaskStatusCountRequestDto struct {
	apiv1.RequestDto[
		struct {
			UserAgent string `json:"userAgent" validate:"required,isuseragent"`
		},
		struct{},
		struct {
			Permission sharedtypes.AccessControlPermission `json:"permission" validate:"isaccesscontrolpermission,required"`
		},
		struct{},
	]
}
type VisualizeMyRoutineTaskStatusCountResponseDto = RoutineTaskCountResponseDto
type VisualizeMyRoutineTaskPurposeCountRequestDto = VisualizeMyRoutineTaskStatusCountRequestDto
type VisualizeMyRoutineTaskPurposeCountResponseDto = RoutineTaskCountResponseDto
type VisualizeMyRoutineTaskScheduledAtCountRequestDto struct {
	apiv1.RequestDto[
		struct {
			UserAgent string `json:"userAgent" validate:"required,isuseragent"`
		},
		struct{},
		struct {
			Permission          sharedtypes.AccessControlPermission `json:"permission" validate:"isaccesscontrolpermission,required"`
			TimeHourUnit        int                                 `json:"timeHourUnit" validate:"required,min=1"`
			QueryRangeStartedAt time.Time                           `json:"queryRangeStartedAt" validate:"required"`
			QueryRangeEndedAt   time.Time                           `json:"queryRangeEndedAt" validate:"required"`
		},
		struct{},
	]
}
type VisualizeMyRoutineTaskScheduledAtCountResponseDto = RoutineTaskCountResponseDto
type VisualizeMyRoutineTaskActualStartedAtCountRequestDto = VisualizeMyRoutineTaskScheduledAtCountRequestDto
type VisualizeMyRoutineTaskActualStartedAtCountResponseDto = RoutineTaskCountResponseDto
type VisualizeMyRoutineTaskActualEndedAtCountRequestDto = VisualizeMyRoutineTaskScheduledAtCountRequestDto
type VisualizeMyRoutineTaskActualEndedAtCountResponseDto = RoutineTaskCountResponseDto
