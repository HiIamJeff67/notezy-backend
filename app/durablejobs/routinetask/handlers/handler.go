package handlers

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/google/uuid"

	adapters "github.com/HiIamJeff67/notezy-backend/app/adapters"
	dtos "github.com/HiIamJeff67/notezy-backend/app/dtos"
	exceptions "github.com/HiIamJeff67/notezy-backend/app/exceptions"
	schemas "github.com/HiIamJeff67/notezy-backend/app/models/schemas"
	validation "github.com/HiIamJeff67/notezy-backend/app/validation"
)

type PurposeHandlerFunc func(
	ctx context.Context,
	tasks []schemas.RoutineTask,
	taskIdToOwnerId map[uuid.UUID]uuid.UUID,
) ([]bool, *exceptions.Exception)

func decodePayload[T any](task schemas.RoutineTask) (*T, *exceptions.Exception) {
	var payload T
	if err := json.Unmarshal(task.Payload, &payload); err != nil {
		return nil, exceptions.RoutineTask.InvalidDto().WithOrigin(err)
	}
	if err := validation.Validator.Struct(payload); err != nil {
		return nil, exceptions.RoutineTask.InvalidDto().WithOrigin(err)
	}
	return &payload, nil
}

func flattenArborizedBlock(
	editableBlockAdapter adapters.EditableBlockAdapterInterface,
	blockGroupId uuid.UUID,
	arborizedEditableBlock *dtos.ArborizedEditableBlock,
) ([]schemas.Block, []uuid.UUID, int64, *exceptions.Exception) {
	rawFlattenedBlocks, totalSize, exception := editableBlockAdapter.FlattenToRaw(arborizedEditableBlock)
	if exception != nil {
		return nil, nil, 0, exception
	}
	if len(rawFlattenedBlocks) == 0 {
		return nil, nil, 0, exceptions.RoutineTask.InvalidDto().
			WithOrigin(fmt.Errorf("arborizedEditableBlock must contain at least one block"))
	}

	blocks := make([]schemas.Block, len(rawFlattenedBlocks))
	blockIds := make([]uuid.UUID, len(rawFlattenedBlocks))
	for index, rawFlattenedBlock := range rawFlattenedBlocks {
		blockType := rawFlattenedBlock.Type
		if rawFlattenedBlock.Id == uuid.Nil || !blockType.IsValidEnum() {
			return nil, nil, 0, exceptions.RoutineTask.InvalidDto().
				WithOrigin(fmt.Errorf("invalid arborizedEditableBlock at flattened index %d", index))
		}

		blockIds[index] = rawFlattenedBlock.Id
		blocks[index] = schemas.Block{
			Id:            rawFlattenedBlock.Id,
			ParentBlockId: rawFlattenedBlock.ParentBlockId,
			BlockGroupId:  blockGroupId,
			Type:          rawFlattenedBlock.Type,
			Props:         rawFlattenedBlock.Props,
			Content:       rawFlattenedBlock.Content,
		}
	}
	return blocks, blockIds, totalSize, nil
}
