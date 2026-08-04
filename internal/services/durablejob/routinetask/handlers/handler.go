package handlers

import (
	"context"
	"fmt"
	"net/http"

	validator "github.com/go-playground/validator/v10"
	"github.com/google/uuid"

	typescontract "github.com/HiIamJeff67/notezy-backend/contracts/types"
	exceptions "github.com/HiIamJeff67/notezy-backend/internal/exceptions"
	schemas "github.com/HiIamJeff67/notezy-backend/internal/services/core/data/database/schemas"
	enums "github.com/HiIamJeff67/notezy-backend/internal/services/core/data/database/schemas/enums"
	editableblock "github.com/HiIamJeff67/notezy-backend/shared/editableblock"
	jsonpayload "github.com/HiIamJeff67/notezy-backend/shared/lib/jsonpayload"
)

type PurposeHandler struct {
	HandlerFunc        PurposeHandlerFunc
	AllowedPermissions []enums.AccessControlPermission
}

type PurposeHandlerFunc func(
	ctx context.Context,
	tasks []schemas.RoutineTask,
	taskIdToActorUserId map[uuid.UUID]uuid.UUID,
	allowedPermissions []enums.AccessControlPermission,
) ([]bool, *exceptions.Exception)

func decodePayload[T any](validator *validator.Validate, task schemas.RoutineTask) (*T, *exceptions.Exception) {
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
	if err := validator.Struct(payload); err != nil {
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
	blockPackId uuid.UUID,
	arborizedEditableBlock *typescontract.ArborizedEditableBlock,
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
	rawFlattenedBlocks, totalSize, err := editableblock.FlattenEditableBlock(arborizedEditableBlock)
	if err != nil {
		return nil, nil, 0, exceptions.New(
			"InvalidRoutineTaskPayload",
			"RoutineTask",
			"Resolve",
			"Routine task payload is invalid",
			http.StatusBadRequest,
		).WithOrigin(err)
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
		blockType := enums.BlockType(rawFlattenedBlock.Type)
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
