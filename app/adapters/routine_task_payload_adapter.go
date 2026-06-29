package adapters

import (
	"bytes"
	"encoding/json"
	"fmt"

	"gorm.io/datatypes"

	dtos "github.com/HiIamJeff67/notezy-backend/app/dtos"
	exceptions "github.com/HiIamJeff67/notezy-backend/app/exceptions"
	enums "github.com/HiIamJeff67/notezy-backend/app/models/schemas/enums"
	validation "github.com/HiIamJeff67/notezy-backend/app/validation"
	concurrency "github.com/HiIamJeff67/notezy-backend/shared/lib/concurrency"
)

type RoutineTaskPayloadAdapterInterface interface {
	Parse(purpose enums.RoutineTaskPurpose, payload datatypes.JSON) *exceptions.Exception
}

type RoutineTaskPayloadAdapter struct {
	editableBlockAdapter EditableBlockAdapterInterface
}

/* ============================== Routine Task Payload Adapter Constructor ============================== */

func NewRoutineTaskPayloadAdapter(
	editableBlockAdapter EditableBlockAdapterInterface,
) RoutineTaskPayloadAdapterInterface {
	if editableBlockAdapter == nil {
		editableBlockAdapter = NewEditableBlockAdapter()
	}
	return &RoutineTaskPayloadAdapter{
		editableBlockAdapter: editableBlockAdapter,
	}
}

/* ============================== Routine Task Payload Parser ============================== */

func (a *RoutineTaskPayloadAdapter) Parse(
	purpose enums.RoutineTaskPurpose,
	payload datatypes.JSON,
) *exceptions.Exception {
	if len(bytes.TrimSpace(payload)) == 0 {
		return exceptions.RoutineTask.InvalidDto().WithOrigin(fmt.Errorf("payload is required"))
	}

	switch purpose {
	case enums.RoutineTaskPurpose_CreateRootShelf:
		var parsedPayload dtos.CreateRootShelfRoutineTaskPayload
		if err := json.Unmarshal(payload, &parsedPayload); err != nil {
			return exceptions.RoutineTask.InvalidDto().WithOrigin(err)
		}
		if err := validation.Validator.Struct(&parsedPayload); err != nil {
			return exceptions.RoutineTask.InvalidDto().WithOrigin(err)
		}
		return nil

	case enums.RoutineTaskPurpose_UpdateRootShelf:
		var parsedPayload dtos.UpdateRootShelfRoutineTaskPayload
		if err := json.Unmarshal(payload, &parsedPayload); err != nil {
			return exceptions.RoutineTask.InvalidDto().WithOrigin(err)
		}
		if err := validation.Validator.Struct(&parsedPayload); err != nil {
			return exceptions.RoutineTask.InvalidDto().WithOrigin(err)
		}
		return nil

	case enums.RoutineTaskPurpose_ResetRootShelf:
		var parsedPayload dtos.ResetRootShelfRoutineTaskPayload
		if err := json.Unmarshal(payload, &parsedPayload); err != nil {
			return exceptions.RoutineTask.InvalidDto().WithOrigin(err)
		}
		if err := validation.Validator.Struct(&parsedPayload); err != nil {
			return exceptions.RoutineTask.InvalidDto().WithOrigin(err)
		}
		return nil

	case enums.RoutineTaskPurpose_CreateSubShelf:
		var parsedPayload dtos.CreateSubShelfRoutineTaskPayload
		if err := json.Unmarshal(payload, &parsedPayload); err != nil {
			return exceptions.RoutineTask.InvalidDto().WithOrigin(err)
		}
		if err := validation.Validator.Struct(&parsedPayload); err != nil {
			return exceptions.RoutineTask.InvalidDto().WithOrigin(err)
		}
		return nil

	case enums.RoutineTaskPurpose_UpdateSubShelf:
		var parsedPayload dtos.UpdateSubShelfRoutineTaskPayload
		if err := json.Unmarshal(payload, &parsedPayload); err != nil {
			return exceptions.RoutineTask.InvalidDto().WithOrigin(err)
		}
		if err := validation.Validator.Struct(&parsedPayload); err != nil {
			return exceptions.RoutineTask.InvalidDto().WithOrigin(err)
		}
		return nil

	case enums.RoutineTaskPurpose_ResetSubShelf:
		var parsedPayload dtos.ResetSubShelfRoutineTaskPayload
		if err := json.Unmarshal(payload, &parsedPayload); err != nil {
			return exceptions.RoutineTask.InvalidDto().WithOrigin(err)
		}
		if err := validation.Validator.Struct(&parsedPayload); err != nil {
			return exceptions.RoutineTask.InvalidDto().WithOrigin(err)
		}
		return nil

	case enums.RoutineTaskPurpose_CreateBlockPack:
		var parsedPayload dtos.CreateBlockPackRoutineTaskPayload
		if err := json.Unmarshal(payload, &parsedPayload); err != nil {
			return exceptions.RoutineTask.InvalidDto().WithOrigin(err)
		}
		if err := validation.Validator.Struct(&parsedPayload); err != nil {
			return exceptions.RoutineTask.InvalidDto().WithOrigin(err)
		}

		validateBlockDto := make([]dtos.ArborizedEditableBlock, len(parsedPayload.Template.BlockGroups))
		for index, blockGroup := range parsedPayload.Template.BlockGroups {
			validateBlockDto[index] = blockGroup.ArborizedEditableBlock
		}

		validateBlockFunc := func(validateDto dtos.ArborizedEditableBlock) (bool, error) {
			if exception := a.validateArborizedEditableBlock(&validateDto); exception != nil {
				return false, exception.GetOrigin()
			}
			return true, nil
		}

		validateBlockResults := concurrency.Execute(
			validateBlockDto,
			min(10, max(len(validateBlockDto)/10, len(validateBlockDto)%10)),
			validateBlockFunc,
		)

		for _, validateBlockResult := range validateBlockResults {
			if validateBlockResult.Err != nil {
				return exceptions.RoutineTask.InvalidDto().
					WithOrigin(fmt.Errorf(
						"invalid template.blockGroups[%d].arborizedEditableBlock: %w",
						validateBlockResult.Index,
						validateBlockResult.Err,
					))
			}
		}
		return nil

	case enums.RoutineTaskPurpose_UpdateBlockPack:
		var parsedPayload dtos.UpdateBlockPackRoutineTaskPayload
		if err := json.Unmarshal(payload, &parsedPayload); err != nil {
			return exceptions.RoutineTask.InvalidDto().WithOrigin(err)
		}
		if err := validation.Validator.Struct(&parsedPayload); err != nil {
			return exceptions.RoutineTask.InvalidDto().WithOrigin(err)
		}

		for index, updatedBlock := range parsedPayload.UpdatedBlocks {
			if exception := a.validateArborizedEditableBlock(updatedBlock.ArborizedEditableBlock); exception != nil {
				return exceptions.RoutineTask.InvalidDto().
					WithOrigin(fmt.Errorf(
						"invalid updatedBlocks[%d].arborizedEditableBlock: %w",
						index,
						exception.GetOrigin(),
					))
			}
		}
		return nil

	case enums.RoutineTaskPurpose_ResetBlockPack:
		var parsedPayload dtos.ResetBlockPackRoutineTaskPayload
		if err := json.Unmarshal(payload, &parsedPayload); err != nil {
			return exceptions.RoutineTask.InvalidDto().WithOrigin(err)
		}
		if err := validation.Validator.Struct(&parsedPayload); err != nil {
			return exceptions.RoutineTask.InvalidDto().WithOrigin(err)
		}
		return nil

	case enums.RoutineTaskPurpose_AppendBlock:
		var parsedPayload dtos.AppendBlockRoutineTaskPayload
		if err := json.Unmarshal(payload, &parsedPayload); err != nil {
			return exceptions.RoutineTask.InvalidDto().WithOrigin(err)
		}
		if err := validation.Validator.Struct(&parsedPayload); err != nil {
			return exceptions.RoutineTask.InvalidDto().WithOrigin(err)
		}
		if exception := a.validateArborizedEditableBlock(&parsedPayload.ArborizedEditableBlock); exception != nil {
			return exception
		}
		return nil

	case enums.RoutineTaskPurpose_UpdateBlock:
		var parsedPayload dtos.UpdateBlockRoutineTaskPayload
		if err := json.Unmarshal(payload, &parsedPayload); err != nil {
			return exceptions.RoutineTask.InvalidDto().WithOrigin(err)
		}
		if err := validation.Validator.Struct(&parsedPayload); err != nil {
			return exceptions.RoutineTask.InvalidDto().WithOrigin(err)
		}
		if exception := a.validateArborizedEditableBlock(parsedPayload.ArborizedEditableBlock); exception != nil {
			return exception
		}
		return nil

	case enums.RoutineTaskPurpose_ResetBlock:
		var parsedPayload dtos.ResetBlockRoutineTaskPayload
		if err := json.Unmarshal(payload, &parsedPayload); err != nil {
			return exceptions.RoutineTask.InvalidDto().WithOrigin(err)
		}
		if err := validation.Validator.Struct(&parsedPayload); err != nil {
			return exceptions.RoutineTask.InvalidDto().WithOrigin(err)
		}
		return nil

	case enums.RoutineTaskPurpose_CreateRoutine:
		var parsedPayload dtos.CreateRoutineRoutineTaskPayload
		if err := json.Unmarshal(payload, &parsedPayload); err != nil {
			return exceptions.RoutineTask.InvalidDto().WithOrigin(err)
		}
		if err := validation.Validator.Struct(&parsedPayload); err != nil {
			return exceptions.RoutineTask.InvalidDto().WithOrigin(err)
		}
		return nil

	case enums.RoutineTaskPurpose_UpdateRoutine:
		var parsedPayload dtos.UpdateRoutineRoutineTaskPayload
		if err := json.Unmarshal(payload, &parsedPayload); err != nil {
			return exceptions.RoutineTask.InvalidDto().WithOrigin(err)
		}
		if err := validation.Validator.Struct(&parsedPayload); err != nil {
			return exceptions.RoutineTask.InvalidDto().WithOrigin(err)
		}
		return nil

	default:
		return exceptions.RoutineTask.InvalidDto().
			WithOrigin(fmt.Errorf("unsupported routine task purpose: %s", purpose))
	}
}

func (a *RoutineTaskPayloadAdapter) validateArborizedEditableBlock(
	arborizedEditableBlock *dtos.ArborizedEditableBlock,
) *exceptions.Exception {
	if arborizedEditableBlock == nil {
		return exceptions.RoutineTask.InvalidDto().WithOrigin(fmt.Errorf("arborizedEditableBlock is required"))
	}

	rawFlattenedBlocks, _, exception := a.editableBlockAdapter.FlattenToRaw(arborizedEditableBlock)
	if exception != nil {
		return exception
	}
	if len(rawFlattenedBlocks) == 0 {
		return exceptions.RoutineTask.InvalidDto().
			WithOrigin(fmt.Errorf("arborizedEditableBlock must contain at least one block"))
	}

	if len(arborizedEditableBlock.Children) > 0 {
		return exceptions.RoutineTask.InvalidDto().
			WithOrigin(fmt.Errorf("arborizedEditableBlock must not contain children for update operations"))
	}

	return nil
}
