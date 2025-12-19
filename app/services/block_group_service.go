package services

import (
	"context"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"

	adapters "notezy-backend/app/adapters"
	dtos "notezy-backend/app/dtos"
	exceptions "notezy-backend/app/exceptions"
	inputs "notezy-backend/app/models/inputs"
	repositories "notezy-backend/app/models/repositories"
	schemas "notezy-backend/app/models/schemas"
	enums "notezy-backend/app/models/schemas/enums"
	options "notezy-backend/app/options"
	validation "notezy-backend/app/validation"
	constants "notezy-backend/shared/constants"
	concurrency "notezy-backend/shared/lib/concurrency"
	types "notezy-backend/shared/types"
)

/* ============================== Interface & Instance ============================== */

type BlockGroupServiceInterface interface {
	GetMyBlockGroupById(ctx context.Context, reqDto *dtos.GetMyBlockGroupByIdReqDto) (*dtos.GetMyBlockGroupByIdResDto, *exceptions.Exception)
	GetMyBlockGroupAndItsBlocksById(ctx context.Context, reqDto *dtos.GetMyBlockGroupAndItsBlocksByIdReqDto) (*dtos.GetMyBlockGroupAndItsBlocksByIdResDto, *exceptions.Exception)
	GetMyBlockGroupsAndTheirBlocksByBlockPackId(ctx context.Context, reqDto *dtos.GetMyBlockGroupsAndTheirBlocksByBlockPackIdReqDto) (*dtos.GetMyBlockGroupsAndTheirBlocksByBlockPackIdResDto, *exceptions.Exception)
	GetMyBlockGroupsByPrevBlockGroupId(ctx context.Context, reqDto *dtos.GetMyBlockGroupsByPrevBlockGroupIdReqDto) (*dtos.GetMyBlockGroupsByPrevBlockGroupIdResDto, *exceptions.Exception)
	GetAllMyBlockGroupsByBlockPackId(ctx context.Context, reqDto *dtos.GetAllMyBlockGroupsByBlockPackIdReqDto) (*dtos.GetAllMyBlockGroupsByBlockPackIdResDto, *exceptions.Exception)
	InsertBlockGroupByBlockPackId(ctx context.Context, reqDto *dtos.CreateBlockGroupByBlockPackIdReqDto) (*dtos.CreateBlockGroupByBlockPackIdResDto, *exceptions.Exception)
	InsertBlockGroupAndItsBlocksByBlockPackId(ctx context.Context, reqDto *dtos.CreateBlockGroupAndItsBlocksByBlockPackIdReqDto) (*dtos.CreateBlockGroupAndItsBlocksByBlockPackIdResDto, *exceptions.Exception)
	InsertBlockGroupsAndTheirBlocksByBlockPackId(ctx context.Context, reqDto *dtos.CreateBlockGroupsAndTheirBlocksByBlockPackIdReqDto) (*dtos.CreateBlockGroupsAndTheirBlocksByBlockPackIdResDto, *exceptions.Exception)
	MoveMyBlockGroupsByIds(ctx context.Context, reqDto *dtos.MoveMyBlockGroupsByIdsReqDto) (*dtos.MoveMyBlockGroupsByIdsResDto, *exceptions.Exception)
	RestoreMyBlockGroupById(ctx context.Context, reqDto *dtos.RestoreMyBlockGroupByIdReqDto) (*dtos.RestoreMyBlockGroupByIdResDto, *exceptions.Exception)
	RestoreMyBlockGroupsByIds(ctx context.Context, reqDto *dtos.RestoreMyBlockGroupsByIdsReqDto) (*dtos.RestoreMyBlockGroupsByIdsResDto, *exceptions.Exception)
	DeleteMyBlockGroupById(ctx context.Context, reqDto *dtos.DeleteMyBlockGroupByIdReqDto) (*dtos.DeleteMyBlockGroupByIdResDto, *exceptions.Exception)
	DeleteMyBlockGroupsByIds(ctx context.Context, reqDto *dtos.DeleteMyBlockGroupsByIdsReqDto) (*dtos.DeleteMyBlockGroupsByIdsResDto, *exceptions.Exception)
}

type BlockGroupService struct {
	db                   *gorm.DB
	blockGroupRepository repositories.BlockGroupRepositoryInterface
	blockRepository      repositories.BlockRepositoryInterface
	editableBlockAdapter adapters.EditableBlockAdapterInterface
}

func NewBlockGroupService(
	db *gorm.DB,
	blockGroupRepository repositories.BlockGroupRepositoryInterface,
	blockRepository repositories.BlockRepositoryInterface,
	editableBlockAdapter adapters.EditableBlockAdapterInterface,
) BlockGroupServiceInterface {
	return &BlockGroupService{
		db:                   db,
		blockGroupRepository: blockGroupRepository,
		blockRepository:      blockRepository,
		editableBlockAdapter: editableBlockAdapter,
	}
}

/* ============================== Implementations ============================== */

func (s *BlockGroupService) GetMyBlockGroupById(
	ctx context.Context, reqDto *dtos.GetMyBlockGroupByIdReqDto,
) (*dtos.GetMyBlockGroupByIdResDto, *exceptions.Exception) {
	if err := validation.Validator.Struct(reqDto); err != nil {
		return nil, exceptions.BlockGroup.InvalidDto().WithError(err)
	}

	db := s.db.WithContext(ctx)

	blockGroup, exception := s.blockGroupRepository.GetOneById(
		reqDto.Param.BlockGroupId,
		reqDto.ContextFields.UserId,
		nil,
		options.WithDB(db),
		options.WithOnlyDeleted(types.Ternary_Negative),
	)
	if exception != nil {
		return nil, exception
	}

	return &dtos.GetMyBlockGroupByIdResDto{
		Id:               blockGroup.Id,
		BlockPackId:      blockGroup.BlockPackId,
		PrevBlockGroupId: blockGroup.PrevBlockGroupId,
		SyncBlockGroupId: blockGroup.SyncBlockGroupId,
		MegaByteSize:     blockGroup.MegaByteSize,
		DeletedAt:        blockGroup.DeletedAt,
		UpdatedAt:        blockGroup.UpdatedAt,
		CreatedAt:        blockGroup.CreatedAt,
	}, nil
}

func (s *BlockGroupService) GetMyBlockGroupAndItsBlocksById(
	ctx context.Context, reqDto *dtos.GetMyBlockGroupAndItsBlocksByIdReqDto,
) (*dtos.GetMyBlockGroupAndItsBlocksByIdResDto, *exceptions.Exception) {
	if err := validation.Validator.Struct(reqDto); err != nil {
		return nil, exceptions.BlockGroup.InvalidDto().WithError(err)
	}

	db := s.db.WithContext(ctx)

	allowedPermissions := []enums.AccessControlPermission{
		enums.AccessControlPermission_Owner,
		enums.AccessControlPermission_Admin,
		enums.AccessControlPermission_Write,
		enums.AccessControlPermission_Read,
	}

	blockGroup, exception := s.blockGroupRepository.CheckPermissionAndGetOneById(
		reqDto.Param.BlockGroupId,
		reqDto.ContextFields.UserId,
		nil,
		allowedPermissions,
		options.WithDB(db),
		options.WithOnlyDeleted(types.Ternary_Negative),
	)
	if exception != nil {
		return nil, exception
	}

	var blocks []schemas.Block
	result := s.db.Model(&schemas.Block{}).
		Where("block_group_id = ?", blockGroup.Id).
		Find(&blocks)
	if err := result.Error; err != nil || len(blocks) == 0 {
		// return the current node with empty editable block if we cannot find them
		return &dtos.GetMyBlockGroupAndItsBlocksByIdResDto{
			Id:                        blockGroup.Id,
			BlockPackId:               blockGroup.BlockPackId,
			PrevBlockGroupId:          blockGroup.PrevBlockGroupId,
			SyncBlockGroupId:          blockGroup.SyncBlockGroupId,
			DeletedAt:                 blockGroup.DeletedAt,
			UpdatedAt:                 blockGroup.UpdatedAt,
			CreatedAt:                 blockGroup.CreatedAt,
			RawArborizedEditableBlock: dtos.RawArborizedEditableBlock{},
		}, nil
	}

	var root *dtos.RawFlattenedEditableBlock = nil
	childrenMap := make(map[uuid.UUID][]dtos.RawFlattenedEditableBlock, len(blocks))
	for _, block := range blocks {
		if block.ParentBlockId == nil {
			if root != nil {
				// duplicate root block detected
				return nil, exceptions.BlockGroup.RepeatedRootBlockInBlockGroupDetected(blocks[0].BlockGroupId, block.Id)
			}

			root = &dtos.RawFlattenedEditableBlock{
				Id:            block.Id,
				ParentBlockId: nil,
				Type:          block.Type,
				Props:         block.Props,
				Content:       block.Content,
			}
		} else {
			childrenMap[*block.ParentBlockId] = append(childrenMap[*block.ParentBlockId], dtos.RawFlattenedEditableBlock{
				Id:            block.Id,
				ParentBlockId: block.ParentBlockId,
				Type:          block.Type,
				Props:         block.Props,
				Content:       block.Content,
			})
		}
	}

	rawArborizedBlock, exception := s.editableBlockAdapter.ArborizeRawToRaw(root, childrenMap)
	if exception != nil {
		return nil, exception
	}

	return &dtos.GetMyBlockGroupAndItsBlocksByIdResDto{
		Id:                        blockGroup.Id,
		BlockPackId:               blockGroup.BlockPackId,
		PrevBlockGroupId:          blockGroup.PrevBlockGroupId,
		SyncBlockGroupId:          blockGroup.SyncBlockGroupId,
		DeletedAt:                 blockGroup.DeletedAt,
		UpdatedAt:                 blockGroup.UpdatedAt,
		CreatedAt:                 blockGroup.CreatedAt,
		RawArborizedEditableBlock: *rawArborizedBlock,
	}, nil
}

func (s *BlockGroupService) GetMyBlockGroupsAndTheirBlocksByBlockPackId(
	ctx context.Context, reqDto *dtos.GetMyBlockGroupsAndTheirBlocksByBlockPackIdReqDto,
) (*dtos.GetMyBlockGroupsAndTheirBlocksByBlockPackIdResDto, *exceptions.Exception) {
	if err := validation.Validator.Struct(reqDto); err != nil {
		return nil, exceptions.BlockGroup.InvalidDto().WithError(err)
	}

	db := s.db.WithContext(ctx)

	allowedPermissions := []enums.AccessControlPermission{
		enums.AccessControlPermission_Owner,
		enums.AccessControlPermission_Admin,
		enums.AccessControlPermission_Write,
		enums.AccessControlPermission_Read,
	}

	blockGroups, exception := s.blockGroupRepository.CheckPermissionsAndGetManyByBlockPackId(
		reqDto.Param.BlockPackId,
		reqDto.ContextFields.UserId,
		nil,
		allowedPermissions,
		options.WithDB(db),
		options.WithOnlyDeleted(types.Ternary_Negative),
	)
	if exception != nil {
		return nil, exception
	}

	var resDto dtos.GetMyBlockGroupsAndTheirBlocksByBlockPackIdResDto

	blockGroupIds := make([]uuid.UUID, len(blockGroups))
	for index, blockGroup := range blockGroups {
		blockGroupIds[index] = blockGroup.Id
		resDto = append(resDto, dtos.GetMyBlockGroupAndItsBlocksByIdResDto{
			Id:                        blockGroup.Id,
			BlockPackId:               blockGroup.BlockPackId,
			PrevBlockGroupId:          blockGroup.PrevBlockGroupId,
			SyncBlockGroupId:          blockGroup.SyncBlockGroupId,
			DeletedAt:                 blockGroup.DeletedAt,
			UpdatedAt:                 blockGroup.UpdatedAt,
			CreatedAt:                 blockGroup.CreatedAt,
			RawArborizedEditableBlock: dtos.RawArborizedEditableBlock{},
		})
	}

	var flattenedBlocks []schemas.Block
	result := s.db.Model(&schemas.Block{}).
		Where("block_group_id IN ?", blockGroupIds).
		Find(&flattenedBlocks)
	if err := result.Error; err != nil || len(flattenedBlocks) == 0 {
		// return the current node with empty editable block if we cannot find them
		return &resDto, nil
	}

	blockGroupToBlocksMap := make(map[uuid.UUID][]schemas.Block)
	for _, flattenedBlock := range flattenedBlocks {
		blockGroupToBlocksMap[flattenedBlock.BlockGroupId] = append(blockGroupToBlocksMap[flattenedBlock.BlockGroupId], flattenedBlock)
	}

	for index, dto := range resDto {
		blocks, exist := blockGroupToBlocksMap[dto.Id]
		if !exist {
			// skip the block groups with no children blocks
			continue
		}

		var root *dtos.RawFlattenedEditableBlock = nil
		childrenMap := make(map[uuid.UUID][]dtos.RawFlattenedEditableBlock, len(blocks))
		for _, block := range blocks {
			if block.ParentBlockId == nil {
				if root != nil {
					// duplicate root block detected
					return nil, exceptions.BlockGroup.RepeatedRootBlockInBlockGroupDetected(blocks[0].BlockGroupId, block.Id)
				}

				root = &dtos.RawFlattenedEditableBlock{
					Id:            block.Id,
					ParentBlockId: nil,
					Type:          block.Type,
					Props:         block.Props,
					Content:       block.Content,
				}
			} else {
				childrenMap[*block.ParentBlockId] = append(childrenMap[*block.ParentBlockId], dtos.RawFlattenedEditableBlock{
					Id:            block.Id,
					ParentBlockId: block.ParentBlockId,
					Type:          block.Type,
					Props:         block.Props,
					Content:       block.Content,
				})
			}
		}

		rawArborizedBlock, exception := s.editableBlockAdapter.ArborizeRawToRaw(root, childrenMap)
		if exception != nil {
			return nil, exception
		}

		if rawArborizedBlock != nil {
			resDto[index].RawArborizedEditableBlock = *rawArborizedBlock
		}
	}

	return &resDto, nil
}

func (s *BlockGroupService) GetMyBlockGroupsByPrevBlockGroupId(
	ctx context.Context, reqDto *dtos.GetMyBlockGroupsByPrevBlockGroupIdReqDto,
) (*dtos.GetMyBlockGroupsByPrevBlockGroupIdResDto, *exceptions.Exception) {
	if err := validation.Validator.Struct(reqDto); err != nil {
		return nil, exceptions.BlockGroup.InvalidDto().WithError(err)
	}

	db := s.db.WithContext(ctx)

	allowedPermissions := []enums.AccessControlPermission{
		enums.AccessControlPermission_Owner,
		enums.AccessControlPermission_Admin,
		enums.AccessControlPermission_Write,
		enums.AccessControlPermission_Read,
	}

	var resDto dtos.GetMyBlockGroupsByPrevBlockGroupIdResDto

	subQuery := db.Model(&schemas.UsersToShelves{}).
		Select("1").
		Where("root_shelf_id = \"SubShelfTable\".root_shelf_id").
		Where("user_id = ? AND permission IN ?",
			reqDto.ContextFields.UserId, allowedPermissions,
		)
	result := db.Model(&schemas.BlockGroup{}).
		Joins("INNER JOIN \"BlockPackTable\" bp ON block_pack_id = bp.id").
		Joins("INNER JOIN \"SubShelfTable\" ss ON bp.parent_sub_shelf_id = ss.id").
		Where("prev_block_group_id = ? AND EXISTS (?) AND deleted_at IS NULL",
			reqDto.Param.PrevBlockGroupId, subQuery,
		).Find(&resDto)
	if err := result.Error; err != nil {
		return nil, exceptions.BlockGroup.NotFound().WithError(err)
	}

	return &resDto, nil
}

func (s *BlockGroupService) GetAllMyBlockGroupsByBlockPackId(
	ctx context.Context, reqDto *dtos.GetAllMyBlockGroupsByBlockPackIdReqDto,
) (*dtos.GetAllMyBlockGroupsByBlockPackIdResDto, *exceptions.Exception) {
	if err := validation.Validator.Struct(reqDto); err != nil {
		return nil, exceptions.BlockGroup.InvalidDto().WithError(err)
	}

	db := s.db.WithContext(ctx)

	allowedPermissions := []enums.AccessControlPermission{
		enums.AccessControlPermission_Owner,
		enums.AccessControlPermission_Admin,
		enums.AccessControlPermission_Write,
		enums.AccessControlPermission_Read,
	}

	var resDto dtos.GetAllMyBlockGroupsByBlockPackIdResDto

	subQuery := db.Model(&schemas.UsersToShelves{}).
		Select("1").
		Where("root_shelf_id = \"SubShelfTable\".root_shelf_id").
		Where("user_id = ? AND permission IN ?",
			reqDto.ContextFields.UserId, allowedPermissions,
		)
	result := db.Model(&schemas.BlockGroup{}).
		Joins("INNER JOIN \"BlockPackTable\" bp ON block_pack_id = bp.id").
		Joins("INNER JOIN \"SubShelfTable\" ss ON bp.parent_sub_shelf_id = ss.id").
		Where("block_pack_id = ? AND EXISTS (?) AND deleted_at IS NULL",
			reqDto.Param.BlockPackId, subQuery,
		).Find(&resDto)
	if err := result.Error; err != nil {
		return nil, exceptions.BlockGroup.NotFound().WithError(err)
	}

	return &resDto, nil
}

func (s *BlockGroupService) InsertBlockGroupByBlockPackId(
	ctx context.Context, reqDto *dtos.CreateBlockGroupByBlockPackIdReqDto,
) (*dtos.CreateBlockGroupByBlockPackIdResDto, *exceptions.Exception) {
	if err := validation.Validator.Struct(reqDto); err != nil {
		return nil, exceptions.BlockGroup.InvalidDto().WithError(err)
	}

	tx := s.db.WithContext(ctx)

	newBlockGroupId, exception := s.blockGroupRepository.InsertOneByBlockPackId(
		reqDto.Body.BlockPackId,
		reqDto.ContextFields.UserId,
		inputs.CreateBlockGroupInput{
			PrevBlockGroupId: reqDto.Body.PrevBlockGroupId,
		},
		options.WithDB(tx),
	)
	if exception != nil {
		return nil, exception
	}
	if newBlockGroupId == nil {
		return nil, exceptions.BlockGroup.FailedToCreate().WithDetails("got nil block group id")
	}

	return &dtos.CreateBlockGroupByBlockPackIdResDto{
		Id:        *newBlockGroupId,
		CreatedAt: time.Now(),
	}, nil
}

func (s *BlockGroupService) InsertBlockGroupAndItsBlocksByBlockPackId(
	ctx context.Context, reqDto *dtos.CreateBlockGroupAndItsBlocksByBlockPackIdReqDto,
) (*dtos.CreateBlockGroupAndItsBlocksByBlockPackIdResDto, *exceptions.Exception) {
	if err := validation.Validator.Struct(reqDto); err != nil {
		return nil, exceptions.BlockGroup.InvalidDto().WithError(err)
	}

	tx := s.db.WithContext(ctx).Begin()

	newBlockGroupId, exception := s.blockGroupRepository.InsertOneByBlockPackId(
		reqDto.Body.BlockPackId,
		reqDto.ContextFields.UserId,
		inputs.CreateBlockGroupInput{
			PrevBlockGroupId: reqDto.Body.PrevBlockGroupId,
		},
		options.WithDB(tx),
	)
	if exception != nil {
		tx.Rollback()
		return nil, exception
	}
	if newBlockGroupId == nil {
		tx.Rollback()
		return nil, exceptions.BlockGroup.FailedToCreate().WithDetails("got nil block group id")
	}

	rawFlattenedBlocks, exception := s.editableBlockAdapter.FlattenToRaw(&reqDto.Body.ArborizedEditableBlock)
	if exception != nil {
		return nil, exception
	}

	input := make([]inputs.CreateBlockInput, len(rawFlattenedBlocks))
	for index, rawFlattenedBlock := range rawFlattenedBlocks {
		input[index] = inputs.CreateBlockInput{
			Id:            rawFlattenedBlock.Id,
			ParentBlockId: rawFlattenedBlock.ParentBlockId,
			Type:          rawFlattenedBlock.Type,
			Props:         rawFlattenedBlock.Props,
			Content:       rawFlattenedBlock.Content,
		}
	}

	_, exception = s.blockRepository.CreateManyByBlockGroupId(
		*newBlockGroupId,
		reqDto.ContextFields.UserId,
		input,
		options.WithDB(tx),
		options.WithBatchSize(constants.MaxBatchCreateBlockSize),
		options.WithOnlyDeleted(types.Ternary_Negative),
		options.WithSkipPermissionCheck(),
	)
	if exception != nil {
		tx.Rollback()
		return nil, exception
	}

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return nil, exceptions.BlockGroup.FailedToCommitTransaction().WithError(err)
	}

	return &dtos.CreateBlockGroupAndItsBlocksByBlockPackIdResDto{
		Id:        *newBlockGroupId,
		CreatedAt: time.Now(),
	}, nil
}

// TODO: use concurrency to run seperate create operation using FlattenRaw
func (s *BlockGroupService) InsertBlockGroupsAndTheirBlocksByBlockPackId(
	ctx context.Context, reqDto *dtos.CreateBlockGroupsAndTheirBlocksByBlockPackIdReqDto,
) (*dtos.CreateBlockGroupsAndTheirBlocksByBlockPackIdResDto, *exceptions.Exception) {
	if err := validation.Validator.Struct(reqDto); err != nil {
		return nil, exceptions.BlockGroup.InvalidDto().WithError(err)
	}

	// try transaction here, but not passing the transaction into the concurrency part
	// this may cause problem, if it doesn't work, maybe we should try transaction by our hands
	tx := s.db.WithContext(ctx).Begin()

	createBlockGroupsInput := make([]inputs.CreateBlockGroupInput, len(reqDto.Body.BlockGroupContents))
	validateBlockDto := make([]dtos.ArborizedEditableBlock, len(reqDto.Body.BlockGroupContents))
	for index, blockGroupContent := range reqDto.Body.BlockGroupContents {
		createBlockGroupsInput[index] = inputs.CreateBlockGroupInput{
			PrevBlockGroupId: blockGroupContent.PrevBlockGroupId,
		}
		validateBlockDto[index] = blockGroupContent.ArborizedEditableBlock
	}

	// note that the order of the output newBlockGroupIds is the same as the order of reqDto.Body.BlockGroupContents
	newBlockGroupIds, exception := s.blockGroupRepository.InsertManyByBlockPackId(
		reqDto.Body.BlockPackId,
		reqDto.ContextFields.UserId,
		createBlockGroupsInput,
		options.WithDB(tx),
	)
	if exception != nil {
		tx.Rollback()
		return nil, exception
	}
	if len(newBlockGroupIds) == 0 {
		tx.Rollback()
		return nil, exceptions.BlockGroup.FailedToCreate().WithDetails("got nil block group id")
	}

	validateBlockFunc := func(validateDto dtos.ArborizedEditableBlock) ([]dtos.RawFlattenedEditableBlock, error) {
		data, exception := s.editableBlockAdapter.FlattenToRaw(&validateDto)
		return data, exception.ToError()
	}

	validateBlockResults := concurrency.BatchExecute(
		validateBlockDto,
		10,
		validateBlockFunc,
	)

	resDto := dtos.CreateBlockGroupsAndTheirBlocksByBlockPackIdResDto{
		IsAllSuccess:   true,
		FailedIndexes:  []int{},
		SuccessIndexes: []int{},
		SuccessBlockGroupAndBlockIds: []struct {
			BlockGroupId uuid.UUID
			BlockIds     []uuid.UUID
		}{},
		CreatedAt: time.Now(),
	}
	var createBlocksInput []inputs.CreateBlockGroupContentInput
	for _, validateResult := range validateBlockResults {
		if validateResult.Err != nil {
			resDto.SuccessIndexes = append(resDto.SuccessIndexes, validateResult.Index)
			// note that since the order of newBlockGroupIds is the same as the reqDto
			// and here the concurrency job worker will output a result with index in Result.Index
			// which provide us enough ability to reorder and align the result with the block group ids here
			blockIds := make([]uuid.UUID, len(validateResult.Data))
			createBlocksByBlockGroupInput := make([]inputs.CreateBlockInput, len(validateResult.Data))
			for index, rawFlattenedBlock := range validateResult.Data {
				blockIds[index] = rawFlattenedBlock.Id
				createBlocksByBlockGroupInput[index] = inputs.CreateBlockInput{
					Id:            rawFlattenedBlock.Id,
					ParentBlockId: rawFlattenedBlock.ParentBlockId,
					Type:          rawFlattenedBlock.Type,
					Props:         rawFlattenedBlock.Props,
					Content:       rawFlattenedBlock.Content,
				}
			}
			resDto.SuccessBlockGroupAndBlockIds = append(resDto.SuccessBlockGroupAndBlockIds, struct {
				BlockGroupId uuid.UUID
				BlockIds     []uuid.UUID
			}{
				BlockGroupId: newBlockGroupIds[validateResult.Index],
				BlockIds:     blockIds,
			})
			createBlocksInput = append(createBlocksInput, inputs.CreateBlockGroupContentInput{
				BlockGroupId: newBlockGroupIds[validateResult.Index],
				Blocks:       createBlocksByBlockGroupInput,
			})
		} else {
			resDto.FailedIndexes = append(resDto.FailedIndexes, validateResult.Index)
			resDto.IsAllSuccess = false
		}
	}

	if len(createBlocksInput) == 0 {
		tx.Rollback()
		return nil, exceptions.BlockGroup.FailedToCreate().WithDetails("no valid block tree structure in any of the given block groups")
	}

	_, exception = s.blockRepository.CreateManyByBlockGroupIds(
		reqDto.ContextFields.UserId,
		createBlocksInput,
		options.WithDB(tx),
		options.WithBatchSize(constants.MaxBatchCreateBlockSize),
		options.WithOnlyDeleted(types.Ternary_Negative),
		options.WithSkipPermissionCheck(),
	)
	if exception != nil {
		tx.Rollback()
		return nil, exception
	}

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return nil, exceptions.BlockGroup.FailedToCommitTransaction().WithError(err)
	}

	resDto.CreatedAt = time.Now()
	return &resDto, nil
}

func (s *BlockGroupService) MoveMyBlockGroupsToTheEndByIds(
	ctx context.Context,
) {

}

func (s *BlockGroupService) MoveMyBlockGroupsByIds(
	ctx context.Context, reqDto *dtos.MoveMyBlockGroupsByIdsReqDto,
) (*dtos.MoveMyBlockGroupsByIdsResDto, *exceptions.Exception) {
	if err := validation.Validator.Struct(reqDto); err != nil {
		return nil, exceptions.BlockGroup.InvalidDto().WithError(err)
	}
	if len(reqDto.Body.MovableBlockGroupIds) != len(reqDto.Body.MovablePrevBlockGroupIds) {
		return nil, exceptions.BlockGroup.InvalidDto("The length of movable block group ids is not equal to the length of movable previous block group id")
	}
	if len(reqDto.Body.MovablePrevBlockGroupIds) == 0 ||
		reqDto.Body.DestinationBlockGroupId == reqDto.Body.MovablePrevBlockGroupIds[0] {
		return nil, exceptions.BlockGroup.NoChanges()
	}

	numOfMovableBlockGroups := len(reqDto.Body.MovableBlockGroupIds)

	tx := s.db.WithContext(ctx).Begin()

	allowedPermissions := []enums.AccessControlPermission{
		enums.AccessControlPermission_Owner,
		enums.AccessControlPermission_Admin,
		enums.AccessControlPermission_Write,
	}

	movableBlockGroups, exception := s.blockGroupRepository.CheckPermissionsAndGetManyByIds(
		reqDto.Body.MovableBlockGroupIds,
		reqDto.ContextFields.UserId,
		nil,
		allowedPermissions,
		options.WithDB(tx),
		options.WithOnlyDeleted(types.Ternary_Negative),
	)
	if exception != nil {
		tx.Rollback()
		return nil, exception
	}

	if numOfMovableBlockGroups != len(movableBlockGroups) {
		tx.Rollback()
		return nil, exceptions.BlockGroup.InvalidDto("Invalid block groups detected")
	}

	movableBlockGroupsMap := make(map[uuid.UUID]schemas.BlockGroup, len(movableBlockGroups))
	for _, movableBlockGroup := range movableBlockGroups {
		movableBlockGroupsMap[movableBlockGroup.Id] = movableBlockGroup
	}

	// check the validation of the given block groups by:
	visited := make(map[*uuid.UUID]bool, len(movableBlockGroups)) // it is also gaurantee that there's only one block group with the its PrevBlockGroupId equals to nil
	for index, movableBlockGroupId := range reqDto.Body.MovableBlockGroupIds {
		// 1. check repeated PrevBlockGroupId for cycular block groups
		if visited[reqDto.Body.MovablePrevBlockGroupIds[index]] {
			tx.Rollback()
			return nil, exceptions.BlockGroup.InvalidDto("Detect cycle in the given block groups")
		}
		visited[reqDto.Body.MovablePrevBlockGroupIds[index]] = true

		// 2. check the given block group is also exist in the database
		blockGroup, exist := movableBlockGroupsMap[movableBlockGroupId]
		if !exist {
			tx.Rollback()
			return nil, exceptions.BlockGroup.NotFound("Cannot find the block group with id of %s", movableBlockGroupId.String())
		}

		// 3. check if the data is consistent on BlockGroupId and PrevBlockGroupId
		if blockGroup.PrevBlockGroupId != reqDto.Body.MovablePrevBlockGroupIds[index] {
			tx.Rollback()
			prevBlockGroupIdString := blockGroup.PrevBlockGroupId.String()
			if blockGroup.PrevBlockGroupId == nil {
				prevBlockGroupIdString = "null"
			}
			movablePrevBlockGroupIdString := reqDto.Body.MovablePrevBlockGroupIds[index].String()
			if reqDto.Body.MovablePrevBlockGroupIds[index] == nil {
				movablePrevBlockGroupIdString = "null"
			}
			return nil, exceptions.BlockGroup.InvalidDto(
				"The given block group has different previous block group id from the actual block group, expected %s, got %s",
				prevBlockGroupIdString,
				movablePrevBlockGroupIdString,
			)
		}

		// 4. check if all the BlockPackIds in the database are the same and equal to the given BlockPackId
		if blockGroup.BlockPackId != reqDto.Body.BlockPackId {
			tx.Rollback()
			return nil, exceptions.BlockGroup.InvalidDto("The given block groups are not all under the given block pack")
		}

		// 5. check to make sure the destination is not inside the movable block groups
		if reqDto.Body.DestinationBlockGroupId != nil &&
			movableBlockGroupId == *reqDto.Body.DestinationBlockGroupId &&
			index != numOfMovableBlockGroups-1 { // it is allowed to move the movable block groups to the next of the final movable block group
			tx.Rollback()
			return nil, exceptions.BlockGroup.InvalidDto("Cannot move the block groups to a position inside of themselves")
		}
	}

	// note that we have gaurantee that the below block groups MUST exist on the above procedure
	startBlockGroup := movableBlockGroupsMap[reqDto.Body.MovableBlockGroupIds[0]]
	endBlockGroup := movableBlockGroupsMap[reqDto.Body.MovableBlockGroupIds[numOfMovableBlockGroups-1]]
	// ignore the exception which indicate the collapsed block group is not exist
	collapsedBlockGroup, _ := s.blockGroupRepository.GetOneByPrevBlockGroupId(
		reqDto.Body.BlockPackId,
		&endBlockGroup.Id,
		reqDto.ContextFields.UserId,
		nil,
		options.WithDB(tx),
		options.WithOnlyDeleted(types.Ternary_Negative),
	)
	// ignore the exception which indicate the destination next block group is not exist
	destinationNextBlockGroup, _ := s.blockGroupRepository.GetOneByPrevBlockGroupId(
		reqDto.Body.BlockPackId,
		reqDto.Body.DestinationBlockGroupId,
		reqDto.ContextFields.UserId,
		nil,
		options.WithDB(tx),
		options.WithOnlyDeleted(types.Ternary_Negative),
	)

	// Sstep 1: connect the start block group to the destination block group
	if _, exception = s.blockGroupRepository.UpdateOneById(
		startBlockGroup.Id,
		reqDto.ContextFields.UserId,
		inputs.PartialUpdateBlockGroupInput{
			Values: inputs.UpdateBlockGroupInput{
				PrevBlockGroupId: reqDto.Body.DestinationBlockGroupId,
			},
			SetNull: nil,
		},
		options.WithDB(tx),
	); exception != nil {
		tx.Rollback()
		return nil, exception
	}

	// step 2: connect the destination next block group to the end block group
	if destinationNextBlockGroup != nil {
		// check if the destination next block group is the start block group which means we don't need to maintain here
		if destinationNextBlockGroup.Id != startBlockGroup.Id {
			if _, exception = s.blockGroupRepository.UpdateOneById(
				destinationNextBlockGroup.Id,
				reqDto.ContextFields.UserId,
				inputs.PartialUpdateBlockGroupInput{
					Values: inputs.UpdateBlockGroupInput{
						PrevBlockGroupId: &endBlockGroup.Id,
					},
					SetNull: nil,
				},
				options.WithDB(tx),
			); exception != nil {
				tx.Rollback()
				return nil, exception
			}
		}
	}

	// step 3: relink the collapsed block group to the old prev block group of the start block group
	if collapsedBlockGroup != nil {
		if _, exception = s.blockGroupRepository.UpdateOneById(
			collapsedBlockGroup.Id,
			reqDto.ContextFields.UserId,
			inputs.PartialUpdateBlockGroupInput{
				Values: inputs.UpdateBlockGroupInput{
					PrevBlockGroupId: startBlockGroup.PrevBlockGroupId,
				},
				SetNull: nil,
			},
			options.WithDB(tx),
		); exception != nil {
			tx.Rollback()
			return nil, exception
		}
	}

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return nil, exceptions.BlockGroup.FailedToCommitTransaction().WithError(err)
	}

	return &dtos.MoveMyBlockGroupsByIdsResDto{
		UpdatedAt: time.Now(),
	}, nil
}

func (s *BlockGroupService) RestoreMyBlockGroupById(
	ctx context.Context, reqDto *dtos.RestoreMyBlockGroupByIdReqDto,
) (*dtos.RestoreMyBlockGroupByIdResDto, *exceptions.Exception) {
	if err := validation.Validator.Struct(reqDto); err != nil {
		return nil, exceptions.BlockGroup.InvalidDto().WithError(err)
	}

	db := s.db.WithContext(ctx)

	if exception := s.blockGroupRepository.RestoreSoftDeletedOneById(
		reqDto.Body.BlockGroupId,
		reqDto.ContextFields.UserId,
		options.WithDB(db),
	); exception != nil {
		return nil, exception
	}

	return &dtos.RestoreMyBlockGroupByIdResDto{
		UpdatedAt: time.Now(),
	}, nil
}

func (s *BlockGroupService) RestoreMyBlockGroupsByIds(
	ctx context.Context, reqDto *dtos.RestoreMyBlockGroupsByIdsReqDto,
) (*dtos.RestoreMyBlockGroupsByIdsResDto, *exceptions.Exception) {
	if err := validation.Validator.Struct(reqDto); err != nil {
		return nil, exceptions.BlockGroup.InvalidDto().WithError(err)
	}

	db := s.db.WithContext(ctx)

	if exception := s.blockGroupRepository.RestoreSoftDeletedManyByIds(
		reqDto.Body.BlockGroupIds,
		reqDto.ContextFields.UserId,
		options.WithDB(db),
	); exception != nil {
		return nil, exception
	}

	return &dtos.RestoreMyBlockGroupsByIdsResDto{
		UpdatedAt: time.Now(),
	}, nil
}

func (s *BlockGroupService) DeleteMyBlockGroupById(
	ctx context.Context, reqDto *dtos.DeleteMyBlockGroupByIdReqDto,
) (*dtos.DeleteMyBlockGroupByIdResDto, *exceptions.Exception) {
	if err := validation.Validator.Struct(reqDto); err != nil {
		return nil, exceptions.BlockGroup.InvalidDto().WithError(err)
	}

	db := s.db.WithContext(ctx)

	if exception := s.blockGroupRepository.SoftDeleteOneById(
		reqDto.Body.BlockGroupId,
		reqDto.ContextFields.UserId,
		options.WithDB(db),
	); exception != nil {
		return nil, exception
	}

	return &dtos.DeleteMyBlockGroupByIdResDto{
		DeletedAt: time.Now(),
	}, nil
}

func (s *BlockGroupService) DeleteMyBlockGroupsByIds(
	ctx context.Context, reqDto *dtos.DeleteMyBlockGroupsByIdsReqDto,
) (*dtos.DeleteMyBlockGroupsByIdsResDto, *exceptions.Exception) {
	if err := validation.Validator.Struct(reqDto); err != nil {
		return nil, exceptions.BlockGroup.InvalidDto().WithError(err)
	}

	db := s.db.WithContext(ctx)

	if exception := s.blockGroupRepository.SoftDeleteManyByIds(
		reqDto.Body.BlockGroupIds,
		reqDto.ContextFields.UserId,
		options.WithDB(db),
	); exception != nil {
		return nil, exception
	}

	return &dtos.DeleteMyBlockGroupsByIdsResDto{
		DeletedAt: time.Now(),
	}, nil
}
