package enums

import (
	"database/sql/driver"
	"fmt"
	"reflect"
	"slices"
)

type RoutineTaskRecordErrorCode string

const (
	RoutineTaskRecordErrorCode_PermissionDenied  RoutineTaskRecordErrorCode = "PermissionDenied"
	RoutineTaskRecordErrorCode_PayloadInvalid    RoutineTaskRecordErrorCode = "PayloadInvalid"
	RoutineTaskRecordErrorCode_TargetNotFound    RoutineTaskRecordErrorCode = "TargetNotFound"
	RoutineTaskRecordErrorCode_PlanLimitExceeded RoutineTaskRecordErrorCode = "PlanLimitExceeded"
	RoutineTaskRecordErrorCode_HandlerFailed     RoutineTaskRecordErrorCode = "HandlerFailed"
	RoutineTaskRecordErrorCode_DatabaseError     RoutineTaskRecordErrorCode = "DatabaseError"
	RoutineTaskRecordErrorCode_Timeout           RoutineTaskRecordErrorCode = "Timeout"
	RoutineTaskRecordErrorCode_Canceled          RoutineTaskRecordErrorCode = "Canceled"
	RoutineTaskRecordErrorCode_Unknown           RoutineTaskRecordErrorCode = "Unknown"
)

var AllRoutineTaskRecordErrorCodes = []RoutineTaskRecordErrorCode{
	RoutineTaskRecordErrorCode_PermissionDenied,
	RoutineTaskRecordErrorCode_PayloadInvalid,
	RoutineTaskRecordErrorCode_TargetNotFound,
	RoutineTaskRecordErrorCode_PlanLimitExceeded,
	RoutineTaskRecordErrorCode_HandlerFailed,
	RoutineTaskRecordErrorCode_DatabaseError,
	RoutineTaskRecordErrorCode_Timeout,
	RoutineTaskRecordErrorCode_Canceled,
	RoutineTaskRecordErrorCode_Unknown,
}

var AllRoutineTaskRecordErrorCodeStrings = []string{
	string(RoutineTaskRecordErrorCode_PermissionDenied),
	string(RoutineTaskRecordErrorCode_PayloadInvalid),
	string(RoutineTaskRecordErrorCode_TargetNotFound),
	string(RoutineTaskRecordErrorCode_PlanLimitExceeded),
	string(RoutineTaskRecordErrorCode_HandlerFailed),
	string(RoutineTaskRecordErrorCode_DatabaseError),
	string(RoutineTaskRecordErrorCode_Timeout),
	string(RoutineTaskRecordErrorCode_Canceled),
	string(RoutineTaskRecordErrorCode_Unknown),
}

func (rtrec RoutineTaskRecordErrorCode) Name() string {
	return reflect.TypeOf(rtrec).Name()
}

func (rtrec *RoutineTaskRecordErrorCode) Scan(value any) error {
	switch v := value.(type) {
	case []byte:
		*rtrec = RoutineTaskRecordErrorCode(string(v))
		return nil
	case string:
		*rtrec = RoutineTaskRecordErrorCode(v)
		return nil
	}
	return scanError(value, rtrec)
}

func (rtrec RoutineTaskRecordErrorCode) Value() (driver.Value, error) {
	return string(rtrec), nil
}

func (rtrec RoutineTaskRecordErrorCode) String() string {
	return string(rtrec)
}

func (rtrec *RoutineTaskRecordErrorCode) IsValidEnum() bool {
	return slices.Contains(AllRoutineTaskRecordErrorCodes, *rtrec)
}

func ConvertStringToRoutineTaskRecordErrorCode(enumString string) (*RoutineTaskRecordErrorCode, error) {
	for _, routineTaskRecordErrorCode := range AllRoutineTaskRecordErrorCodes {
		if string(routineTaskRecordErrorCode) == enumString {
			return &routineTaskRecordErrorCode, nil
		}
	}
	return nil, fmt.Errorf("invalid routine task record error code: %s", enumString)
}
