package resolvers

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"gorm.io/gorm"

	dtos "github.com/HiIamJeff67/notezy-backend/app/dtos"
	exceptions "github.com/HiIamJeff67/notezy-backend/app/exceptions"
	inputs "github.com/HiIamJeff67/notezy-backend/app/models/inputs"
	repositories "github.com/HiIamJeff67/notezy-backend/app/models/repositories"
	schemas "github.com/HiIamJeff67/notezy-backend/app/models/schemas"
	enums "github.com/HiIamJeff67/notezy-backend/app/models/schemas/enums"
	options "github.com/HiIamJeff67/notezy-backend/app/options"
	types "github.com/HiIamJeff67/notezy-backend/shared/types"
)

type BlockPatternResolverInterface interface {
	Resolve(ctx context.Context, ownerId uuid.UUID, pattern dtos.RoutineTaskPattern) (map[string]string, *exceptions.Exception)
}

type BlockPatternResolver struct {
	db              *gorm.DB
	blockRepository repositories.BlockRepositoryInterface
}

func NewBlockPatternResolver(db *gorm.DB, blockRepository repositories.BlockRepositoryInterface) BlockPatternResolverInterface {
	return BlockPatternResolver{
		db:              db,
		blockRepository: blockRepository,
	}
}

func (r BlockPatternResolver) Resolve(
	ctx context.Context,
	ownerId uuid.UUID,
	pattern dtos.RoutineTaskPattern,
) (map[string]string, *exceptions.Exception) {
	values := map[string]string{}
	checkInputs := make([]inputs.BulkCheckBlockPermissionInput, 0)
	keysByBlockId := map[uuid.UUID][]string{}

	for key, binding := range pattern {
		if binding.Source != PatternSourceBlockText {
			continue
		}
		if binding.BlockId == nil || *binding.BlockId == uuid.Nil {
			return nil, exceptions.RoutineTask.InvalidDto().
				WithOrigin(fmt.Errorf("pattern.%s.blockId is required", key))
		}
		keysByBlockId[*binding.BlockId] = append(keysByBlockId[*binding.BlockId], key)
	}
	for blockId := range keysByBlockId {
		checkInputs = append(checkInputs, inputs.BulkCheckBlockPermissionInput{
			UserId: ownerId,
			Id:     blockId,
		})
	}
	if len(checkInputs) == 0 {
		return values, nil
	}
	if r.db == nil || r.blockRepository == nil {
		return nil, exceptions.RoutineTask.InvalidDto().
			WithOrigin(fmt.Errorf("block pattern source is not available"))
	}

	allowedPermissions := []enums.AccessControlPermission{
		enums.AccessControlPermission_Owner,
		enums.AccessControlPermission_Admin,
		enums.AccessControlPermission_Write,
		enums.AccessControlPermission_Read,
	}

	successes, blocks, exception := r.blockRepository.BulkCheckPermissionsAndGetManyByIds(
		checkInputs,
		nil,
		allowedPermissions,
		options.WithDB(r.db.WithContext(ctx)),
		options.WithOnlyDeleted(types.Ternary_Negative),
	)
	if exception != nil {
		return nil, exception
	}

	blocksById := make(map[uuid.UUID]schemas.Block, len(blocks))
	for _, block := range blocks {
		blocksById[block.Id] = block
	}
	for index, success := range successes {
		if !success {
			return nil, exceptions.Block.NoPermission("read pattern source block")
		}
		block := blocksById[checkInputs[index].Id]
		var content any
		if err := json.Unmarshal(block.Content, &content); err != nil {
			return nil, exceptions.RoutineTask.InvalidDto().WithOrigin(err)
		}
		parts := make([]string, 0)
		var walk func(any)
		walk = func(current any) {
			switch typed := current.(type) {
			case []any:
				for _, item := range typed {
					walk(item)
				}
			case map[string]any:
				if text, ok := typed["text"].(string); ok {
					parts = append(parts, text)
				}
				for _, value := range typed {
					walk(value)
				}
			}
		}
		walk(content)
		text := strings.Join(parts, "")
		for _, key := range keysByBlockId[block.Id] {
			values[key] = text
		}
	}

	return values, nil
}
