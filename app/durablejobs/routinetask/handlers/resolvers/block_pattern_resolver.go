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
	ResolveMany(ctx context.Context, ownerIds []uuid.UUID, patterns []dtos.RoutineTaskPattern) ([]map[string]string, []bool, *exceptions.Exception)
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
	values, successes, exception := r.ResolveMany(ctx, []uuid.UUID{ownerId}, []dtos.RoutineTaskPattern{pattern})
	if exception != nil {
		return nil, exception
	}
	if len(successes) == 0 || !successes[0] {
		return nil, exceptions.RoutineTask.InvalidDto()
	}
	return values[0], nil
}

func (r BlockPatternResolver) ResolveMany(
	ctx context.Context,
	ownerIds []uuid.UUID,
	patterns []dtos.RoutineTaskPattern,
) ([]map[string]string, []bool, *exceptions.Exception) {
	values := make([]map[string]string, len(patterns))
	taskSuccesses := make([]bool, len(patterns))
	for index := range patterns {
		values[index] = map[string]string{}
		taskSuccesses[index] = true
	}
	if len(ownerIds) != len(patterns) {
		return nil, nil, exceptions.RoutineTask.InvalidDto().
			WithOrigin(fmt.Errorf("ownerIds and patterns length mismatch"))
	}

	checkInputs := make([]inputs.BulkCheckBlockPermissionInput, 0)
	keysByUserAndBlockId := map[[2]uuid.UUID][]struct {
		taskIndex int
		key       string
	}{}

	for patternIndex, pattern := range patterns {
		for key, binding := range pattern {
			if binding.Source != PatternSourceBlockText {
				continue
			}
			if binding.BlockId == nil || *binding.BlockId == uuid.Nil {
				taskSuccesses[patternIndex] = false
				continue
			}
			mapKey := [2]uuid.UUID{ownerIds[patternIndex], *binding.BlockId}
			if _, exists := keysByUserAndBlockId[mapKey]; !exists {
				checkInputs = append(checkInputs, inputs.BulkCheckBlockPermissionInput{
					UserId: ownerIds[patternIndex],
					Id:     *binding.BlockId,
				})
			}
			keysByUserAndBlockId[mapKey] = append(keysByUserAndBlockId[mapKey], struct {
				taskIndex int
				key       string
			}{taskIndex: patternIndex, key: key})
		}
	}
	if len(checkInputs) == 0 {
		return values, taskSuccesses, nil
	}
	if r.db == nil || r.blockRepository == nil {
		return nil, nil, exceptions.RoutineTask.InvalidDto().
			WithOrigin(fmt.Errorf("block pattern source is not available"))
	}

	allowedPermissions := []enums.AccessControlPermission{
		enums.AccessControlPermission_Owner,
		enums.AccessControlPermission_Admin,
		enums.AccessControlPermission_Write,
		enums.AccessControlPermission_Read,
	}

	permissionSuccesses, blocks, exception := r.blockRepository.BulkCheckPermissionsAndGetManyByIds(
		checkInputs,
		nil,
		allowedPermissions,
		options.WithDB(r.db.WithContext(ctx)),
		options.WithOnlyDeleted(types.Ternary_Negative),
	)
	if exception != nil {
		return nil, nil, exception
	}

	blocksById := make(map[uuid.UUID]schemas.Block, len(blocks))
	for _, block := range blocks {
		blocksById[block.Id] = block
	}
	for index, success := range permissionSuccesses {
		if !success {
			for _, request := range keysByUserAndBlockId[[2]uuid.UUID{checkInputs[index].UserId, checkInputs[index].Id}] {
				taskSuccesses[request.taskIndex] = false
			}
			continue
		}
		block := blocksById[checkInputs[index].Id]
		var content any
		if err := json.Unmarshal(block.Content, &content); err != nil {
			return nil, nil, exceptions.RoutineTask.InvalidDto().WithOrigin(err)
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
		for _, request := range keysByUserAndBlockId[[2]uuid.UUID{checkInputs[index].UserId, block.Id}] {
			values[request.taskIndex][request.key] = text
		}
	}

	return values, taskSuccesses, nil
}
