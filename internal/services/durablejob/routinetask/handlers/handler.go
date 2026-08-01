package handlers

import (
	"context"
	"fmt"
	"net/http"

	"github.com/google/uuid"

	blocksdto "github.com/HiIamJeff67/notezy-backend/contracts/api/v1/blocks"
	adapters "github.com/HiIamJeff67/notezy-backend/internal/adapters"
	exceptions "github.com/HiIamJeff67/notezy-backend/internal/exceptions"
	schemas "github.com/HiIamJeff67/notezy-backend/internal/services/durablejob/data/schemas"
	enums "github.com/HiIamJeff67/notezy-backend/internal/services/durablejob/data/schemas/enums"
	validation "github.com/HiIamJeff67/notezy-backend/internal/services/durablejob/validation"
	jsonpayload "github.com/HiIamJeff67/notezy-backend/internal/shared/parsers/jsonpayload"
)

type PurposeHandlerFunc func(
	ctx context.Context,
	tasks []schemas.RoutineTask,
	taskIdToActorUserId map[uuid.UUID]uuid.UUID,
	allowedPermissions []enums.AccessControlPermission,
) ([]bool, *exceptions.Exception)

func decodePayload[T any](task schemas.RoutineTask) (*T, *exceptions.Exception) {
	var payload T
	if err := jsonpayload.Decode(task.Payload, &payload); err != nil {
		return nil, exceptions.New(
			"InvalidRoutineTaskPayload",
			"RoutineTask",
			"Resolve",
			"Routine task payload is invalid",
			http.StatusBadRequest,
		).WithOrigin(err)
	}
	if err := validation.Validator.Struct(payload); err != nil {
		return nil, exceptions.New(
			"InvalidRoutineTaskPayload",
			"RoutineTask",
			"Resolve",
			"Routine task payload is invalid",
			http.StatusBadRequest,
		).WithOrigin(err)
	}
	return &payload, nil
}

func flattenArborizedBlock(
	editableBlockAdapter adapters.EditableBlockAdapterInterface,
	blockPackId uuid.UUID,
	arborizedEditableBlock *blocksdto.ArborizedEditableBlock,
) ([]schemas.Block, []uuid.UUID, int64, *exceptions.Exception) {
	if blockPackId == uuid.Nil {
		return nil, nil, 0, exceptions.New(
			"InvalidRoutineTaskPayload",
			"RoutineTask",
			"Resolve",
			"Routine task payload is invalid",
			http.StatusBadRequest,
		).
			WithOrigin(fmt.Errorf("blockPackId is required"))
	}
	rawFlattenedBlocks, totalSize, exception := editableBlockAdapter.FlattenToRaw(arborizedEditableBlock)
	if exception != nil {
		return nil, nil, 0, exception
	}
	if len(rawFlattenedBlocks) == 0 {
		return nil, nil, 0, exceptions.New(
			"InvalidRoutineTaskPayload",
			"RoutineTask",
			"Resolve",
			"Routine task payload is invalid",
			http.StatusBadRequest,
		).
			WithOrigin(fmt.Errorf("arborizedEditableBlock must contain at least one block"))
	}

	blocks := make([]schemas.Block, len(rawFlattenedBlocks))
	blockIds := make([]uuid.UUID, len(rawFlattenedBlocks))
	for index, rawFlattenedBlock := range rawFlattenedBlocks {
		blockType := rawFlattenedBlock.Type
		if rawFlattenedBlock.Id == uuid.Nil || !blockType.IsValidEnum() {
			return nil, nil, 0, exceptions.New(
				"InvalidRoutineTaskPayload",
				"RoutineTask",
				"Resolve",
				"Routine task payload is invalid",
				http.StatusBadRequest,
			).
				WithOrigin(fmt.Errorf("invalid arborizedEditableBlock at flattened index %d", index))
		}

		blockIds[index] = rawFlattenedBlock.Id
		blocks[index] = schemas.Block{
			Id:            rawFlattenedBlock.Id,
			BlockPackId:   blockPackId,
			ParentBlockId: rawFlattenedBlock.ParentBlockId,
			PrevBlockId:   rawFlattenedBlock.PrevBlockId,
			NextBlockId:   rawFlattenedBlock.NextBlockId,
			Type:          enums.BlockType(rawFlattenedBlock.Type),
			Props:         rawFlattenedBlock.Props,
			Content:       rawFlattenedBlock.Content,
		}
	}
	return blocks, blockIds, totalSize, nil
}
