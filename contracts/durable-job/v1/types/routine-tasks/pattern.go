package routinetasktypes

import "github.com/google/uuid"

type RoutineTaskPatternBinding struct {
	Source      string     `json:"source" validate:"required,oneof=scheduledAt recordId shortRecordId routineTaskId blockText blockCheckboxCount"`
	BlockId     *uuid.UUID `json:"blockId" validate:"omitnil"`
	BlockPackId *uuid.UUID `json:"blockPackId" validate:"omitnil"`
	Checked     *bool      `json:"checked" validate:"omitnil"`
	Format      *string    `json:"format" validate:"omitnil,max=64"`
	Timezone    *string    `json:"timezone" validate:"omitnil,max=64,istimezone"`
}

type RoutineTaskPattern map[string]RoutineTaskPatternBinding
