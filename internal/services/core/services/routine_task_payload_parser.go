package services

import (
	"fmt"
	"net/http"

	"gorm.io/datatypes"

	typescontract "github.com/HiIamJeff67/notezy-backend/contracts/types"
	exceptions "github.com/HiIamJeff67/notezy-backend/internal/exceptions"
	enums "github.com/HiIamJeff67/notezy-backend/internal/services/core/data/database/schemas/enums"
	editableblock "github.com/HiIamJeff67/notezy-backend/shared/editableblock"

	concurrency "github.com/HiIamJeff67/notezy-backend/shared/lib/concurrency"
	jsonpayload "github.com/HiIamJeff67/notezy-backend/shared/lib/jsonpayload"
)

func (s *RoutineTaskService) parseRoutineTaskPayload(
	purpose enums.RoutineTaskPurpose,
	payload datatypes.JSON,
) *exceptions.Exception {
	switch purpose {
	case enums.RoutineTaskPurpose_CreateRootShelf:
		var parsedPayload typescontract.CreateRootShelfRoutineTaskPayload
		if err := jsonpayload.Decode(payload, &parsedPayload); err != nil {
			return exceptions.New(
				"InvalidRoutineTaskPayload",
				"RoutineTask",
				"Parse",
				"Routine task payload is invalid",
				http.StatusBadRequest,
			).WithOrigin(err)
		}
		if err := s.validator.Struct(&parsedPayload); err != nil {
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
		var parsedPayload typescontract.UpdateRootShelfRoutineTaskPayload
		if err := jsonpayload.Decode(payload, &parsedPayload); err != nil {
			return exceptions.New(
				"InvalidRoutineTaskPayload",
				"RoutineTask",
				"Parse",
				"Routine task payload is invalid",
				http.StatusBadRequest,
			).WithOrigin(err)
		}
		if err := s.validator.Struct(&parsedPayload); err != nil {
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
		var parsedPayload typescontract.ResetRootShelfRoutineTaskPayload
		if err := jsonpayload.Decode(payload, &parsedPayload); err != nil {
			return exceptions.New(
				"InvalidRoutineTaskPayload",
				"RoutineTask",
				"Parse",
				"Routine task payload is invalid",
				http.StatusBadRequest,
			).WithOrigin(err)
		}
		if err := s.validator.Struct(&parsedPayload); err != nil {
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
		var parsedPayload typescontract.CreateSubShelfRoutineTaskPayload
		if err := jsonpayload.Decode(payload, &parsedPayload); err != nil {
			return exceptions.New(
				"InvalidRoutineTaskPayload",
				"RoutineTask",
				"Parse",
				"Routine task payload is invalid",
				http.StatusBadRequest,
			).WithOrigin(err)
		}
		if err := s.validator.Struct(&parsedPayload); err != nil {
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
		var parsedPayload typescontract.UpdateSubShelfRoutineTaskPayload
		if err := jsonpayload.Decode(payload, &parsedPayload); err != nil {
			return exceptions.New(
				"InvalidRoutineTaskPayload",
				"RoutineTask",
				"Parse",
				"Routine task payload is invalid",
				http.StatusBadRequest,
			).WithOrigin(err)
		}
		if err := s.validator.Struct(&parsedPayload); err != nil {
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
		var parsedPayload typescontract.ResetSubShelfRoutineTaskPayload
		if err := jsonpayload.Decode(payload, &parsedPayload); err != nil {
			return exceptions.New(
				"InvalidRoutineTaskPayload",
				"RoutineTask",
				"Parse",
				"Routine task payload is invalid",
				http.StatusBadRequest,
			).WithOrigin(err)
		}
		if err := s.validator.Struct(&parsedPayload); err != nil {
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
		var parsedPayload typescontract.CreateBlockPackRoutineTaskPayload
		if err := jsonpayload.Decode(payload, &parsedPayload); err != nil {
			return exceptions.New(
				"InvalidRoutineTaskPayload",
				"RoutineTask",
				"Parse",
				"Routine task payload is invalid",
				http.StatusBadRequest,
			).WithOrigin(err)
		}
		if err := s.validator.Struct(&parsedPayload); err != nil {
			return exceptions.New(
				"InvalidRoutineTaskPayload",
				"RoutineTask",
				"Parse",
				"Routine task payload is invalid",
				http.StatusBadRequest,
			).WithOrigin(err)
		}

		validateBlockDto := make([]typescontract.ArborizedEditableBlock, len(parsedPayload.Template.Blocks))
		for index, block := range parsedPayload.Template.Blocks {
			validateBlockDto[index] = block.ArborizedEditableBlock
		}

		validateBlockFunc := func(validateDto typescontract.ArborizedEditableBlock) (bool, error) {
			if exception := validateArborizedEditableBlock(&validateDto); exception != nil {
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
		var parsedPayload typescontract.UpdateBlockPackRoutineTaskPayload
		if err := jsonpayload.Decode(payload, &parsedPayload); err != nil {
			return exceptions.New(
				"InvalidRoutineTaskPayload",
				"RoutineTask",
				"Parse",
				"Routine task payload is invalid",
				http.StatusBadRequest,
			).WithOrigin(err)
		}
		if err := s.validator.Struct(&parsedPayload); err != nil {
			return exceptions.New(
				"InvalidRoutineTaskPayload",
				"RoutineTask",
				"Parse",
				"Routine task payload is invalid",
				http.StatusBadRequest,
			).WithOrigin(err)
		}

		for index, updatedBlock := range parsedPayload.UpdatedBlocks {
			if exception := validateArborizedEditableBlock(updatedBlock.ArborizedEditableBlock); exception != nil {
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
		var parsedPayload typescontract.ResetBlockPackRoutineTaskPayload
		if err := jsonpayload.Decode(payload, &parsedPayload); err != nil {
			return exceptions.New(
				"InvalidRoutineTaskPayload",
				"RoutineTask",
				"Parse",
				"Routine task payload is invalid",
				http.StatusBadRequest,
			).WithOrigin(err)
		}
		if err := s.validator.Struct(&parsedPayload); err != nil {
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
		var parsedPayload typescontract.AppendBlockRoutineTaskPayload
		if err := jsonpayload.Decode(payload, &parsedPayload); err != nil {
			return exceptions.New(
				"InvalidRoutineTaskPayload",
				"RoutineTask",
				"Parse",
				"Routine task payload is invalid",
				http.StatusBadRequest,
			).WithOrigin(err)
		}
		if err := s.validator.Struct(&parsedPayload); err != nil {
			return exceptions.New(
				"InvalidRoutineTaskPayload",
				"RoutineTask",
				"Parse",
				"Routine task payload is invalid",
				http.StatusBadRequest,
			).WithOrigin(err)
		}
		if exception := validateArborizedEditableBlock(&parsedPayload.ArborizedEditableBlock); exception != nil {
			return exception
		}
		return nil

	case enums.RoutineTaskPurpose_UpdateBlock:
		var parsedPayload typescontract.UpdateBlockRoutineTaskPayload
		if err := jsonpayload.Decode(payload, &parsedPayload); err != nil {
			return exceptions.New(
				"InvalidRoutineTaskPayload",
				"RoutineTask",
				"Parse",
				"Routine task payload is invalid",
				http.StatusBadRequest,
			).WithOrigin(err)
		}
		if err := s.validator.Struct(&parsedPayload); err != nil {
			return exceptions.New(
				"InvalidRoutineTaskPayload",
				"RoutineTask",
				"Parse",
				"Routine task payload is invalid",
				http.StatusBadRequest,
			).WithOrigin(err)
		}
		if exception := validateArborizedEditableBlock(parsedPayload.ArborizedEditableBlock); exception != nil {
			return exception
		}
		return nil

	case enums.RoutineTaskPurpose_ResetBlock:
		var parsedPayload typescontract.ResetBlockRoutineTaskPayload
		if err := jsonpayload.Decode(payload, &parsedPayload); err != nil {
			return exceptions.New(
				"InvalidRoutineTaskPayload",
				"RoutineTask",
				"Parse",
				"Routine task payload is invalid",
				http.StatusBadRequest,
			).WithOrigin(err)
		}
		if err := s.validator.Struct(&parsedPayload); err != nil {
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
		var parsedPayload typescontract.CreateRoutineRoutineTaskPayload
		if err := jsonpayload.Decode(payload, &parsedPayload); err != nil {
			return exceptions.New(
				"InvalidRoutineTaskPayload",
				"RoutineTask",
				"Parse",
				"Routine task payload is invalid",
				http.StatusBadRequest,
			).WithOrigin(err)
		}
		if err := s.validator.Struct(&parsedPayload); err != nil {
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
		var parsedPayload typescontract.UpdateRoutineRoutineTaskPayload
		if err := jsonpayload.Decode(payload, &parsedPayload); err != nil {
			return exceptions.New(
				"InvalidRoutineTaskPayload",
				"RoutineTask",
				"Parse",
				"Routine task payload is invalid",
				http.StatusBadRequest,
			).WithOrigin(err)
		}
		if err := s.validator.Struct(&parsedPayload); err != nil {
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

func validateArborizedEditableBlock(
	arborizedEditableBlock *typescontract.ArborizedEditableBlock,
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

	rawFlattenedBlocks, _, err := editableblock.FlattenEditableBlock(arborizedEditableBlock)
	if err != nil {
		return exceptions.New(
			"InvalidRoutineTaskPayload",
			"RoutineTask",
			"Parse",
			"Routine task payload is invalid",
			http.StatusBadRequest,
		).WithOrigin(err)
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
