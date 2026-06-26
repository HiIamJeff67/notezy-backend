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

func (a *RoutineTaskPayloadAdapter) Parse(
	purpose enums.RoutineTaskPurpose,
	payload datatypes.JSON,
) *exceptions.Exception {
	if len(bytes.TrimSpace(payload)) == 0 {
		return exceptions.RoutineTask.InvalidDto().WithOrigin(fmt.Errorf("payload is required"))
	}

	switch purpose {
	case enums.RoutineTaskPurpose_CreateBlockPack:
		return a.parseCreateBlockPackPayload(payload)

	case enums.RoutineTaskPurpose_DeleteBlockPack:
		var parsedPayload dtos.DeleteBlockPackRoutineTaskPayload
		return a.unmarshalAndValidate(payload, &parsedPayload)

	case enums.RoutineTaskPurpose_CreateBlock:
		var parsedPayload dtos.CreateBlockRoutineTaskPayload
		if exception := a.unmarshalAndValidate(payload, &parsedPayload); exception != nil {
			return exception
		}
		if exception := a.validateArborizedEditableBlock(&parsedPayload.ArborizedEditableBlock); exception != nil {
			return exception
		}
		return nil

	case enums.RoutineTaskPurpose_UpdateBlock:
		var parsedPayload dtos.UpdateBlockRoutineTaskPayload
		if exception := a.unmarshalAndValidate(payload, &parsedPayload); exception != nil {
			return exception
		}
		if exception := a.validateArborizedEditableBlock(parsedPayload.ArborizedEditableBlock); exception != nil {
			return exception
		}
		return nil

	case enums.RoutineTaskPurpose_DeleteBlock:
		var parsedPayload dtos.DeleteBlockRoutineTaskPayload
		return a.unmarshalAndValidate(payload, &parsedPayload)

	default:
		return exceptions.RoutineTask.InvalidDto().
			WithOrigin(fmt.Errorf("unsupported routine task purpose: %s", purpose))
	}
}

func (a *RoutineTaskPayloadAdapter) parseCreateBlockPackPayload(
	payload datatypes.JSON,
) *exceptions.Exception {
	var parsedPayload dtos.CreateBlockPackRoutineTaskPayload
	if exception := a.unmarshalAndValidate(payload, &parsedPayload); exception != nil {
		return exception
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
}

func (a *RoutineTaskPayloadAdapter) unmarshalAndValidate(
	payload datatypes.JSON,
	target any,
) *exceptions.Exception {
	if err := json.Unmarshal(payload, target); err != nil {
		return exceptions.RoutineTask.InvalidDto().WithOrigin(err)
	}
	if err := validation.Validator.Struct(target); err != nil {
		return exceptions.RoutineTask.InvalidDto().WithOrigin(err)
	}
	return nil
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
	return nil
}
