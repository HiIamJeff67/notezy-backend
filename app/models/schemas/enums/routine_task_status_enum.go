package enums

import (
	"database/sql/driver"
	"fmt"
	"reflect"
	"slices"
)

type RoutineTaskStatus string

const (
	RoutineTaskStatus_Idle    RoutineTaskStatus = "Idle"
	RoutineTaskStatus_Waiting RoutineTaskStatus = "Waiting" // include scheduling, but we don't need to present to the client
	RoutineTaskStatus_Running RoutineTaskStatus = "Running"
	RoutineTaskStatus_Pause   RoutineTaskStatus = "Pause"
	RoutineTaskStatus_Cancel  RoutineTaskStatus = "Cancel"
	RoutineTaskStatus_Success RoutineTaskStatus = "Success"
	RoutineTaskStatus_Fail    RoutineTaskStatus = "Fail"
)

var AllRoutineTaskStatuses = []RoutineTaskStatus{
	RoutineTaskStatus_Idle,
	RoutineTaskStatus_Waiting,
	RoutineTaskStatus_Running,
	RoutineTaskStatus_Pause,
	RoutineTaskStatus_Cancel,
	RoutineTaskStatus_Success,
	RoutineTaskStatus_Fail,
}

var AllRoutineTaskStatusStrings = []string{
	string(RoutineTaskStatus_Idle),
	string(RoutineTaskStatus_Waiting),
	string(RoutineTaskStatus_Running),
	string(RoutineTaskStatus_Pause),
	string(RoutineTaskStatus_Cancel),
	string(RoutineTaskStatus_Success),
	string(RoutineTaskStatus_Fail),
}

func (rts RoutineTaskStatus) Name() string {
	return reflect.TypeOf(rts).Name()
}

func (rts *RoutineTaskStatus) Scan(value any) error {
	switch v := value.(type) {
	case []byte:
		*rts = RoutineTaskStatus(string(v))
		return nil
	case string:
		*rts = RoutineTaskStatus(v)
		return nil
	}
	return scanError(value, rts)
}

func (rts RoutineTaskStatus) Value() (driver.Value, error) {
	return string(rts), nil
}

func (rts RoutineTaskStatus) String() string {
	return string(rts)
}

func (rts *RoutineTaskStatus) IsValidEnum() bool {
	return slices.Contains(AllRoutineTaskStatuses, *rts)
}

func ConvertStringToRoutineTaskStatus(enumString string) (*RoutineTaskStatus, error) {
	for _, routineTaskStatus := range AllRoutineTaskStatuses {
		if string(routineTaskStatus) == enumString {
			return &routineTaskStatus, nil
		}
	}
	return nil, fmt.Errorf("invalid routine task status: %s", enumString)
}
