package enums

import (
	"database/sql/driver"
	"fmt"
	"reflect"
	"slices"
)

type RoutineTaskPurpose string

const (
	RoutineTaskPurpose_CreateRootShelf RoutineTaskPurpose = "CreateRootShelf" // create a root shelf with nothing inside of it
	RoutineTaskPurpose_UpdateRootShelf RoutineTaskPurpose = "UpdateRootShelf" // update the columns of the given root shelf
	RoutineTaskPurpose_ResetRootShelf  RoutineTaskPurpose = "ResetRootShelf"  // reset the children of the root shelf
	RoutineTaskPurpose_CreateSubShelf  RoutineTaskPurpose = "CreateSubShelf"  // create a sub shelf with nothing inside of it
	RoutineTaskPurpose_UpdateSubShelf  RoutineTaskPurpose = "UpdateSubShelf"  // update the columns of the given sub shelf
	RoutineTaskPurpose_ResetSubShelf   RoutineTaskPurpose = "ResetSubShelf"   // reset the children of the given sub shelf
	RoutineTaskPurpose_CreateBlockPack RoutineTaskPurpose = "CreateBlockPack" // create a block pack with the given content within the routine task payload
	RoutineTaskPurpose_UpdateBlockPack RoutineTaskPurpose = "UpdateBlockPack" // update blocks in the block pack
	RoutineTaskPurpose_ResetBlockPack  RoutineTaskPurpose = "ResetBlockPack"  // reset the block pack to an empty block pack
	RoutineTaskPurpose_AppendBlock     RoutineTaskPurpose = "AppendBlock"     // create a block at the end of the given block pack with the given props and content within the routine task payload
	RoutineTaskPurpose_UpdateBlock     RoutineTaskPurpose = "UpdateBlock"     // update a block with the given props and content within the routine task payload
	RoutineTaskPurpose_ResetBlock      RoutineTaskPurpose = "ResetBlock"      // reset the block to a paragraph with empty props and content
	RoutineTaskPurpose_CreateRoutine   RoutineTaskPurpose = "CreateRoutine"   // create a routine with no links
	RoutineTaskPurpose_UpdateRoutine   RoutineTaskPurpose = "UpdateRoutine"   // update the columns of the given routine, excluded links to it
)

var AllRoutineTaskPurposes = []RoutineTaskPurpose{
	RoutineTaskPurpose_CreateRootShelf,
	RoutineTaskPurpose_UpdateRootShelf,
	RoutineTaskPurpose_ResetRootShelf,
	RoutineTaskPurpose_CreateSubShelf,
	RoutineTaskPurpose_UpdateSubShelf,
	RoutineTaskPurpose_ResetSubShelf,
	RoutineTaskPurpose_CreateBlockPack,
	RoutineTaskPurpose_UpdateBlockPack,
	RoutineTaskPurpose_ResetBlockPack,
	RoutineTaskPurpose_AppendBlock,
	RoutineTaskPurpose_UpdateBlock,
	RoutineTaskPurpose_ResetBlock,
	RoutineTaskPurpose_CreateRoutine,
	RoutineTaskPurpose_UpdateRoutine,
}

var AllRoutineTaskPurposeStrings = []string{
	string(RoutineTaskPurpose_CreateRootShelf),
	string(RoutineTaskPurpose_UpdateRootShelf),
	string(RoutineTaskPurpose_ResetRootShelf),
	string(RoutineTaskPurpose_CreateSubShelf),
	string(RoutineTaskPurpose_UpdateSubShelf),
	string(RoutineTaskPurpose_ResetSubShelf),
	string(RoutineTaskPurpose_CreateBlockPack),
	string(RoutineTaskPurpose_UpdateBlockPack),
	string(RoutineTaskPurpose_ResetBlockPack),
	string(RoutineTaskPurpose_AppendBlock),
	string(RoutineTaskPurpose_UpdateBlock),
	string(RoutineTaskPurpose_ResetBlock),
	string(RoutineTaskPurpose_CreateRoutine),
	string(RoutineTaskPurpose_UpdateRoutine),
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
