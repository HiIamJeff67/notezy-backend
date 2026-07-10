package resolvers

import (
	"context"
	"fmt"
	"strconv"

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

type BlockPackPatternResolverInterface interface {
	Resolve(ctx context.Context, ownerId uuid.UUID, pattern dtos.RoutineTaskPattern) (map[string]string, *exceptions.Exception)
	ResolveMany(ctx context.Context, ownerIds []uuid.UUID, patterns []dtos.RoutineTaskPattern) ([]map[string]string, []bool, *exceptions.Exception)
}

type BlockPackPatternResolver struct {
	db                  *gorm.DB
	blockPackRepository repositories.BlockPackRepositoryInterface
}

func NewBlockPackPatternResolver(db *gorm.DB, blockPackRepository repositories.BlockPackRepositoryInterface) BlockPackPatternResolverInterface {
	return BlockPackPatternResolver{
		db:                  db,
		blockPackRepository: blockPackRepository,
	}
}

func (r BlockPackPatternResolver) Resolve(
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

func (r BlockPackPatternResolver) ResolveMany(
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

	checkInputs := make([]inputs.BulkCheckBlockPackPermissionInput, 0)
	keysByUserAndBlockPackId := map[[2]uuid.UUID][]struct {
		taskIndex int
		key       string
	}{}

	for patternIndex, pattern := range patterns {
		for key, binding := range pattern {
			if binding.Source != PatternSourceBlockCheckboxCount {
				continue
			}
			if binding.BlockPackId == nil || *binding.BlockPackId == uuid.Nil {
				taskSuccesses[patternIndex] = false
				continue
			}
			mapKey := [2]uuid.UUID{ownerIds[patternIndex], *binding.BlockPackId}
			if _, exists := keysByUserAndBlockPackId[mapKey]; !exists {
				checkInputs = append(checkInputs, inputs.BulkCheckBlockPackPermissionInput{
					UserId: ownerIds[patternIndex],
					Id:     *binding.BlockPackId,
				})
			}
			keysByUserAndBlockPackId[mapKey] = append(keysByUserAndBlockPackId[mapKey], struct {
				taskIndex int
				key       string
			}{taskIndex: patternIndex, key: key})
		}
	}
	if len(checkInputs) == 0 {
		return values, taskSuccesses, nil
	}
	if r.db == nil || r.blockPackRepository == nil {
		return nil, nil, exceptions.RoutineTask.InvalidDto().
			WithOrigin(fmt.Errorf("block pack pattern source is not available"))
	}

	allowedPermissions := []enums.AccessControlPermission{
		enums.AccessControlPermission_Owner,
		enums.AccessControlPermission_Admin,
		enums.AccessControlPermission_Write,
		enums.AccessControlPermission_Read,
	}

	permissionSuccesses, _, exception := r.blockPackRepository.BulkCheckPermissionsAndGetManyByIds(
		checkInputs,
		nil,
		allowedPermissions,
		options.WithDB(r.db.WithContext(ctx)),
		options.WithOnlyDeleted(types.Ternary_Negative),
	)
	if exception != nil {
		return nil, nil, exception
	}

	validBlockPackIds := make([]uuid.UUID, 0, len(checkInputs))
	validBlockPackIdSet := map[uuid.UUID]bool{}
	for index, success := range permissionSuccesses {
		if !success {
			for _, request := range keysByUserAndBlockPackId[[2]uuid.UUID{checkInputs[index].UserId, checkInputs[index].Id}] {
				taskSuccesses[request.taskIndex] = false
			}
			continue
		}
		blockPackId := checkInputs[index].Id
		validBlockPackIds = append(validBlockPackIds, blockPackId)
		validBlockPackIdSet[blockPackId] = true
	}
	if len(validBlockPackIds) == 0 {
		return values, taskSuccesses, nil
	}

	var rows []struct {
		BlockPackId uuid.UUID `gorm:"column:block_pack_id"`
		Checked     bool      `gorm:"column:checked"`
	}
	if err := r.db.WithContext(ctx).
		Model(&schemas.Block{}).
		Select(`block_pack_id, COALESCE((props->>'checked')::boolean, false) AS checked`).
		Where("block_pack_id IN ? AND type = ? AND deleted_at IS NULL", validBlockPackIds, enums.BlockType_CheckListItem).
		Find(&rows).Error; err != nil {
		return nil, nil, exceptions.Block.NotFound().WithOrigin(err)
	}

	totalByBlockPackId := map[uuid.UUID]int{}
	checkedByBlockPackId := map[uuid.UUID]int{}
	uncheckedByBlockPackId := map[uuid.UUID]int{}
	for _, row := range rows {
		totalByBlockPackId[row.BlockPackId]++
		if row.Checked {
			checkedByBlockPackId[row.BlockPackId]++
		} else {
			uncheckedByBlockPackId[row.BlockPackId]++
		}
	}

	for patternIndex, pattern := range patterns {
		for key, binding := range pattern {
			if binding.Source != PatternSourceBlockCheckboxCount || binding.BlockPackId == nil {
				continue
			}
			blockPackId := *binding.BlockPackId
			if !validBlockPackIdSet[blockPackId] {
				continue
			}

			count := totalByBlockPackId[blockPackId]
			if binding.Checked != nil {
				if *binding.Checked {
					count = checkedByBlockPackId[blockPackId]
				} else {
					count = uncheckedByBlockPackId[blockPackId]
				}
			}
			values[patternIndex][key] = strconv.Itoa(count)
		}
	}

	return values, taskSuccesses, nil
}
