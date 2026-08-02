package adapters

import (
	"fmt"
	"net/http"

	"gorm.io/datatypes"

	blocksdto "github.com/HiIamJeff67/notezy-backend/contracts/gateway/v1/api/blocks"
	exceptions "github.com/HiIamJeff67/notezy-backend/internal/exceptions"
	enums "github.com/HiIamJeff67/notezy-backend/internal/services/core/data/database/schemas/enums"
	validation "github.com/HiIamJeff67/notezy-backend/internal/services/core/validation"
	payloads "github.com/HiIamJeff67/notezy-backend/internal/services/durablejob/routinetask/payloads"
	concurrency "github.com/HiIamJeff67/notezy-backend/internal/shared/lib/concurrency"
	jsonpayload "github.com/HiIamJeff67/notezy-backend/internal/shared/parsers/jsonpayload"
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
	switch purpose {
	case enums.RoutineTaskPurpose_CreateRootShelf:
		var parsedPayload payloads.CreateRootShelfRoutineTaskPayload
		if err := jsonpayload.Decode(payload, &parsedPayload); err != nil {
			return exceptions.New(
				"InvalidRoutineTaskPayload",
				"RoutineTask",
				"Parse",
				"Routine task payload is invalid",
				http.StatusBadRequest,
			).WithOrigin(err)
		}
		if err := validation.Validator.Struct(&parsedPayload); err != nil {
			return exceptions.New(
				"InvalidRoutineTaskPayload",
				"RoutineTask",
				"Parse",
				"Routine task payload is invalid",
				http.StatusBadRequest,
			).WithOrigin(err)
		}
		return nil

	case enums.RoutineTaskPurpose_UpdateRootShelf:
		var parsedPayload payloads.UpdateRootShelfRoutineTaskPayload
		if err := jsonpayload.Decode(payload, &parsedPayload); err != nil {
			return exceptions.New(
				"InvalidRoutineTaskPayload",
				"RoutineTask",
				"Parse",
				"Routine task payload is invalid",
				http.StatusBadRequest,
			).WithOrigin(err)
		}
		if err := validation.Validator.Struct(&parsedPayload); err != nil {
			return exceptions.New(
				"InvalidRoutineTaskPayload",
				"RoutineTask",
				"Parse",
				"Routine task payload is invalid",
				http.StatusBadRequest,
			).WithOrigin(err)
		}
		return nil

	case enums.RoutineTaskPurpose_ResetRootShelf:
		var parsedPayload payloads.ResetRootShelfRoutineTaskPayload
		if err := jsonpayload.Decode(payload, &parsedPayload); err != nil {
			return exceptions.New(
				"InvalidRoutineTaskPayload",
				"RoutineTask",
				"Parse",
				"Routine task payload is invalid",
				http.StatusBadRequest,
			).WithOrigin(err)
		}
		if err := validation.Validator.Struct(&parsedPayload); err != nil {
			return exceptions.New(
				"InvalidRoutineTaskPayload",
				"RoutineTask",
				"Parse",
				"Routine task payload is invalid",
				http.StatusBadRequest,
			).WithOrigin(err)
		}
		return nil

	case enums.RoutineTaskPurpose_CreateSubShelf:
		var parsedPayload payloads.CreateSubShelfRoutineTaskPayload
		if err := jsonpayload.Decode(payload, &parsedPayload); err != nil {
			return exceptions.New(
				"InvalidRoutineTaskPayload",
				"RoutineTask",
				"Parse",
				"Routine task payload is invalid",
				http.StatusBadRequest,
			).WithOrigin(err)
		}
		if err := validation.Validator.Struct(&parsedPayload); err != nil {
			return exceptions.New(
				"InvalidRoutineTaskPayload",
				"RoutineTask",
				"Parse",
				"Routine task payload is invalid",
				http.StatusBadRequest,
			).WithOrigin(err)
		}
		return nil

	case enums.RoutineTaskPurpose_UpdateSubShelf:
		var parsedPayload payloads.UpdateSubShelfRoutineTaskPayload
		if err := jsonpayload.Decode(payload, &parsedPayload); err != nil {
			return exceptions.New(
				"InvalidRoutineTaskPayload",
				"RoutineTask",
				"Parse",
				"Routine task payload is invalid",
				http.StatusBadRequest,
			).WithOrigin(err)
		}
		if err := validation.Validator.Struct(&parsedPayload); err != nil {
			return exceptions.New(
				"InvalidRoutineTaskPayload",
				"RoutineTask",
				"Parse",
				"Routine task payload is invalid",
				http.StatusBadRequest,
			).WithOrigin(err)
		}
		return nil

	case enums.RoutineTaskPurpose_ResetSubShelf:
		var parsedPayload payloads.ResetSubShelfRoutineTaskPayload
		if err := jsonpayload.Decode(payload, &parsedPayload); err != nil {
			return exceptions.New(
				"InvalidRoutineTaskPayload",
				"RoutineTask",
				"Parse",
				"Routine task payload is invalid",
				http.StatusBadRequest,
			).WithOrigin(err)
		}
		if err := validation.Validator.Struct(&parsedPayload); err != nil {
			return exceptions.New(
				"InvalidRoutineTaskPayload",
				"RoutineTask",
				"Parse",
				"Routine task payload is invalid",
				http.StatusBadRequest,
			).WithOrigin(err)
		}
		return nil

	case enums.RoutineTaskPurpose_CreateBlockPack:
		var parsedPayload payloads.CreateBlockPackRoutineTaskPayload
		if err := jsonpayload.Decode(payload, &parsedPayload); err != nil {
			return exceptions.New(
				"InvalidRoutineTaskPayload",
				"RoutineTask",
				"Parse",
				"Routine task payload is invalid",
				http.StatusBadRequest,
			).WithOrigin(err)
		}
		if err := validation.Validator.Struct(&parsedPayload); err != nil {
			return exceptions.New(
				"InvalidRoutineTaskPayload",
				"RoutineTask",
				"Parse",
				"Routine task payload is invalid",
				http.StatusBadRequest,
			).WithOrigin(err)
		}

		validateBlockDto := make([]blocksdto.ArborizedEditableBlock, len(parsedPayload.Template.Blocks))
		for index, block := range parsedPayload.Template.Blocks {
			validateBlockDto[index] = block.ArborizedEditableBlock
		}

		validateBlockFunc := func(validateDto blocksdto.ArborizedEditableBlock) (bool, error) {
			if exception := a.validateArborizedEditableBlock(&validateDto); exception != nil {
				return false, exception
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
				return exceptions.New(
					"InvalidRoutineTaskPayload",
					"RoutineTask",
					"Parse",
					"Routine task payload is invalid",
					http.StatusBadRequest,
				).
					WithOrigin(fmt.Errorf(
						"invalid template.blocks[%d].arborizedEditableBlock: %w",
						validateBlockResult.Index,
						validateBlockResult.Err,
					))
			}
		}
		return nil

	case enums.RoutineTaskPurpose_UpdateBlockPack:
		var parsedPayload payloads.UpdateBlockPackRoutineTaskPayload
		if err := jsonpayload.Decode(payload, &parsedPayload); err != nil {
			return exceptions.New(
				"InvalidRoutineTaskPayload",
				"RoutineTask",
				"Parse",
				"Routine task payload is invalid",
				http.StatusBadRequest,
			).WithOrigin(err)
		}
		if err := validation.Validator.Struct(&parsedPayload); err != nil {
			return exceptions.New(
				"InvalidRoutineTaskPayload",
				"RoutineTask",
				"Parse",
				"Routine task payload is invalid",
				http.StatusBadRequest,
			).WithOrigin(err)
		}

		for index, updatedBlock := range parsedPayload.UpdatedBlocks {
			if exception := a.validateArborizedEditableBlock(updatedBlock.ArborizedEditableBlock); exception != nil {
				return exceptions.New(
					"InvalidRoutineTaskPayload",
					"RoutineTask",
					"Parse",
					"Routine task payload is invalid",
					http.StatusBadRequest,
				).
					WithOrigin(fmt.Errorf(
						"invalid updatedBlocks[%d].arborizedEditableBlock: %w",
						index,
						exception,
					))
			}
		}
		return nil

	case enums.RoutineTaskPurpose_ResetBlockPack:
		var parsedPayload payloads.ResetBlockPackRoutineTaskPayload
		if err := jsonpayload.Decode(payload, &parsedPayload); err != nil {
			return exceptions.New(
				"InvalidRoutineTaskPayload",
				"RoutineTask",
				"Parse",
				"Routine task payload is invalid",
				http.StatusBadRequest,
			).WithOrigin(err)
		}
		if err := validation.Validator.Struct(&parsedPayload); err != nil {
			return exceptions.New(
				"InvalidRoutineTaskPayload",
				"RoutineTask",
				"Parse",
				"Routine task payload is invalid",
				http.StatusBadRequest,
			).WithOrigin(err)
		}
		return nil

	case enums.RoutineTaskPurpose_AppendBlock:
		var parsedPayload payloads.AppendBlockRoutineTaskPayload
		if err := jsonpayload.Decode(payload, &parsedPayload); err != nil {
			return exceptions.New(
				"InvalidRoutineTaskPayload",
				"RoutineTask",
				"Parse",
				"Routine task payload is invalid",
				http.StatusBadRequest,
			).WithOrigin(err)
		}
		if err := validation.Validator.Struct(&parsedPayload); err != nil {
			return exceptions.New(
				"InvalidRoutineTaskPayload",
				"RoutineTask",
				"Parse",
				"Routine task payload is invalid",
				http.StatusBadRequest,
			).WithOrigin(err)
		}
		if exception := a.validateArborizedEditableBlock(&parsedPayload.ArborizedEditableBlock); exception != nil {
			return exception
		}
		return nil

	case enums.RoutineTaskPurpose_UpdateBlock:
		var parsedPayload payloads.UpdateBlockRoutineTaskPayload
		if err := jsonpayload.Decode(payload, &parsedPayload); err != nil {
			return exceptions.New(
				"InvalidRoutineTaskPayload",
				"RoutineTask",
				"Parse",
				"Routine task payload is invalid",
				http.StatusBadRequest,
			).WithOrigin(err)
		}
		if err := validation.Validator.Struct(&parsedPayload); err != nil {
			return exceptions.New(
				"InvalidRoutineTaskPayload",
				"RoutineTask",
				"Parse",
				"Routine task payload is invalid",
				http.StatusBadRequest,
			).WithOrigin(err)
		}
		if exception := a.validateArborizedEditableBlock(parsedPayload.ArborizedEditableBlock); exception != nil {
			return exception
		}
		return nil

	case enums.RoutineTaskPurpose_ResetBlock:
		var parsedPayload payloads.ResetBlockRoutineTaskPayload
		if err := jsonpayload.Decode(payload, &parsedPayload); err != nil {
			return exceptions.New(
				"InvalidRoutineTaskPayload",
				"RoutineTask",
				"Parse",
				"Routine task payload is invalid",
				http.StatusBadRequest,
			).WithOrigin(err)
		}
		if err := validation.Validator.Struct(&parsedPayload); err != nil {
			return exceptions.New(
				"InvalidRoutineTaskPayload",
				"RoutineTask",
				"Parse",
				"Routine task payload is invalid",
				http.StatusBadRequest,
			).WithOrigin(err)
		}
		return nil

	case enums.RoutineTaskPurpose_CreateRoutine:
		var parsedPayload payloads.CreateRoutineRoutineTaskPayload
		if err := jsonpayload.Decode(payload, &parsedPayload); err != nil {
			return exceptions.New(
				"InvalidRoutineTaskPayload",
				"RoutineTask",
				"Parse",
				"Routine task payload is invalid",
				http.StatusBadRequest,
			).WithOrigin(err)
		}
		if err := validation.Validator.Struct(&parsedPayload); err != nil {
			return exceptions.New(
				"InvalidRoutineTaskPayload",
				"RoutineTask",
				"Parse",
				"Routine task payload is invalid",
				http.StatusBadRequest,
			).WithOrigin(err)
		}
		return nil

	case enums.RoutineTaskPurpose_UpdateRoutine:
		var parsedPayload payloads.UpdateRoutineRoutineTaskPayload
		if err := jsonpayload.Decode(payload, &parsedPayload); err != nil {
			return exceptions.New(
				"InvalidRoutineTaskPayload",
				"RoutineTask",
				"Parse",
				"Routine task payload is invalid",
				http.StatusBadRequest,
			).WithOrigin(err)
		}
		if err := validation.Validator.Struct(&parsedPayload); err != nil {
			return exceptions.New(
				"InvalidRoutineTaskPayload",
				"RoutineTask",
				"Parse",
				"Routine task payload is invalid",
				http.StatusBadRequest,
			).WithOrigin(err)
		}
		return nil

	default:
		return exceptions.New(
			"InvalidRoutineTaskPayload",
			"RoutineTask",
			"Parse",
			"Routine task payload is invalid",
			http.StatusBadRequest,
		).
			WithOrigin(fmt.Errorf("unsupported routine task purpose: %s", purpose))
	}
}

func (a *RoutineTaskPayloadAdapter) validateArborizedEditableBlock(
	arborizedEditableBlock *blocksdto.ArborizedEditableBlock,
) *exceptions.Exception {
	if arborizedEditableBlock == nil {
		return exceptions.New(
			"InvalidRoutineTaskPayload",
			"RoutineTask",
			"Parse",
			"Routine task payload is invalid",
			http.StatusBadRequest,
		).WithOrigin(fmt.Errorf("arborizedEditableBlock is required"))
	}

	rawFlattenedBlocks, _, exception := a.editableBlockAdapter.FlattenToRaw(arborizedEditableBlock)
	if exception != nil {
		return exception
	}
	if len(rawFlattenedBlocks) == 0 {
		return exceptions.New(
			"InvalidRoutineTaskPayload",
			"RoutineTask",
			"Parse",
			"Routine task payload is invalid",
			http.StatusBadRequest,
		).
			WithOrigin(fmt.Errorf("arborizedEditableBlock must contain at least one block"))
	}

	if len(arborizedEditableBlock.Children) > 0 {
		return exceptions.New(
			"InvalidRoutineTaskPayload",
			"RoutineTask",
			"Parse",
			"Routine task payload is invalid",
			http.StatusBadRequest,
		).
			WithOrigin(fmt.Errorf("arborizedEditableBlock must not contain children for update operations"))
	}

	return nil
}
