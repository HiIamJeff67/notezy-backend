package enums

import (
	"database/sql/driver"
	"fmt"
	"reflect"
	"slices"
)

type RoutineTaskPurpose string

const (
	RoutineTaskPurpose_CreateBlockPack RoutineTaskPurpose = "CreateBlockPack"
	RoutineTaskPurpose_DeleteBlockPack RoutineTaskPurpose = "DeleteBlockPack"
	RoutineTaskPurpose_CreateBlock     RoutineTaskPurpose = "CreateBlock"
	RoutineTaskPurpose_UpdateBlock     RoutineTaskPurpose = "UpdateBlock"
	RoutineTaskPurpose_DeleteBlock     RoutineTaskPurpose = "DeleteBlock"
)

var AllRoutineTaskPurposes = []RoutineTaskPurpose{
	RoutineTaskPurpose_CreateBlockPack,
	RoutineTaskPurpose_DeleteBlockPack,
	RoutineTaskPurpose_CreateBlock,
	RoutineTaskPurpose_UpdateBlock,
	RoutineTaskPurpose_DeleteBlock,
}

var AllRoutineTaskPurposeStrings = []string{
	string(RoutineTaskPurpose_CreateBlockPack),
	string(RoutineTaskPurpose_DeleteBlockPack),
	string(RoutineTaskPurpose_CreateBlock),
	string(RoutineTaskPurpose_UpdateBlock),
	string(RoutineTaskPurpose_DeleteBlock),
}

func (rtp RoutineTaskPurpose) Name() string {
	return reflect.TypeOf(rtp).Name()
}

func (rtp *RoutineTaskPurpose) Scan(value any) error {
	switch v := value.(type) {
	case []byte:
		*rtp = RoutineTaskPurpose(string(v))
		return nil
	case string:
		*rtp = RoutineTaskPurpose(v)
		return nil
	}
	return scanError(value, rtp)
}

func (rtp RoutineTaskPurpose) Value() (driver.Value, error) {
	return string(rtp), nil
}

func (rtp RoutineTaskPurpose) String() string {
	return string(rtp)
}

func (rtp *RoutineTaskPurpose) IsValidEnum() bool {
	return slices.Contains(AllRoutineTaskPurposes, *rtp)
}

func ConvertStringToRoutineTaskPurpose(enumString string) (*RoutineTaskPurpose, error) {
	for _, routineTaskPurpose := range AllRoutineTaskPurposes {
		if string(routineTaskPurpose) == enumString {
			return &routineTaskPurpose, nil
		}
	}
	return nil, fmt.Errorf("invalid routine task purpose: %s", enumString)
}
