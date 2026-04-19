package services

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"

	adapters "notezy-backend/app/adapters"
	dtos "notezy-backend/app/dtos"
	exceptions "notezy-backend/app/exceptions"
	concurrency "notezy-backend/app/lib/concurrency"
	inputs "notezy-backend/app/models/inputs"
	repositories "notezy-backend/app/models/repositories"
	schemas "notezy-backend/app/models/schemas"
	enums "notezy-backend/app/models/schemas/enums"
	options "notezy-backend/app/options"
	validation "notezy-backend/app/validation"
	constants "notezy-backend/shared/constants"
	types "notezy-backend/shared/types"
)

type BlockGroupServiceInterface interface {
	GetMyBlockGroupById(ctx context.Context, reqDto *dtos.GetMyBlockGroupByIdReqDto) (*dtos.GetMyBlockGroupByIdResDto, *exceptions.Exception)
	GetMyBlockGroupAndItsBlocksById(ctx context.Context, reqDto *dtos.GetMyBlockGroupAndItsBlocksByIdReqDto) (*dtos.GetMyBlockGroupAndItsBlocksByIdResDto, *exceptions.Exception)
	GetMyBlockGroupsAndTheirBlocksByIds(ctx context.Context, reqDto *dtos.GetMyBlockGroupsAndTheirBlocksByIdsReqDto) (*dtos.GetMyBlockGroupsAndTheirBlocksByIdsResDto, *exceptions.Exception)
	GetMyBlockGroupsAndTheirBlocksByBlockPackId(ctx context.Context, reqDto *dtos.GetMyBlockGroupsAndTheirBlocksByBlockPackIdReqDto) (*dtos.GetMyBlockGroupsAndTheirBlocksByBlockPackIdResDto, *exceptions.Exception)
	GetMyBlockGroupsByPrevBlockGroupId(ctx context.Context, reqDto *dtos.GetMyBlockGroupsByPrevBlockGroupIdReqDto) (*dtos.GetMyBlockGroupsByPrevBlockGroupIdResDto, *exceptions.Exception)
	GetAllMyBlockGroupsByBlockPackId(ctx context.Context, reqDto *dtos.GetAllMyBlockGroupsByBlockPackIdReqDto) (*dtos.GetAllMyBlockGroupsByBlockPackIdResDto, *exceptions.Exception)
	InsertBlockGroupByBlockPackId(ctx context.Context, reqDto *dtos.InsertBlockGroupByBlockPackIdReqDto) (*dtos.InsertBlockGroupByBlockPackIdResDto, *exceptions.Exception)
	InsertBlockGroupAndItsBlocksByBlockPackId(ctx context.Context, reqDto *dtos.InsertBlockGroupAndItsBlocksByBlockPackIdReqDto) (*dtos.InsertBlockGroupAndItsBlocksByBlockPackIdResDto, *exceptions.Exception)
	InsertBlockGroupsAndTheirBlocksByBlockPackId(ctx context.Context, reqDto *dtos.InsertBlockGroupsAndTheirBlocksByBlockPackIdReqDto) (*dtos.InsertBlockGroupsAndTheirBlocksByBlockPackIdResDto, *exceptions.Exception)
	InsertSequentialBlockGroupsAndTheirBlocksByBlockPackId(ctx context.Context, reqDto *dtos.InsertSequentialBlockGroupsAndTheirBlocksByBlockPackIdReqDto) (*dtos.InsertSequentialBlockGroupsAndTheirBlocksByBlockPackIdResDto, *exceptions.Exception)
	MoveMyBlockGroupById(ctx context.Context, reqDto *dtos.MoveMyBlockGroupByIdReqDto) (*dtos.MoveMyBlockGroupByIdResDto, *exceptions.Exception)
	MoveMyBlockGroupsByIds(ctx context.Context, reqDto *dtos.MoveMyBlockGroupsByIdsReqDto) (*dtos.MoveMyBlockGroupsByIdsResDto, *exceptions.Exception)
	BatchMoveMyBlockGroupsByIds(ctx context.Context, reqDto *dtos.BatchMoveMyBlockGroupsByIdsReqDto) (*dtos.BatchMoveMyBlockGroupsByIdsResDto, *exceptions.Exception)
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

func (s *BlockGroupService) GetMyBlockGroupById(
	ctx context.Context, reqDto *dtos.GetMyBlockGroupByIdReqDto,
) (*dtos.GetMyBlockGroupByIdResDto, *exceptions.Exception) {
	if err := validation.Validator.Struct(reqDto); err != nil {
		return nil, exceptions.BlockGroup.InvalidDto().WithOrigin(err)
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
		Size:             blockGroup.Size,
		DeletedAt:        blockGroup.DeletedAt,
		UpdatedAt:        blockGroup.UpdatedAt,
		CreatedAt:        blockGroup.CreatedAt,
	}, nil
}

func (s *BlockGroupService) GetMyBlockGroupAndItsBlocksById(
	ctx context.Context, reqDto *dtos.GetMyBlockGroupAndItsBlocksByIdReqDto,
) (*dtos.GetMyBlockGroupAndItsBlocksByIdResDto, *exceptions.Exception) {
	if err := validation.Validator.Struct(reqDto); err != nil {
		return nil, exceptions.BlockGroup.InvalidDto().WithOrigin(err)
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
	result := db.Model(&schemas.Block{}).
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

	if rawArborizedBlock == nil {
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

func (s *BlockGroupService) GetMyBlockGroupsAndTheirBlocksByIds(
	ctx context.Context, reqDto *dtos.GetMyBlockGroupsAndTheirBlocksByIdsReqDto,
) (*dtos.GetMyBlockGroupsAndTheirBlocksByIdsResDto, *exceptions.Exception) {
	if err := validation.Validator.Struct(reqDto); err != nil {
		return nil, exceptions.BlockGroup.InvalidDto().WithOrigin(err)
	}

	db := s.db.WithContext(ctx)

	allowedPermissions := []enums.AccessControlPermission{
		enums.AccessControlPermission_Owner,
		enums.AccessControlPermission_Admin,
		enums.AccessControlPermission_Write,
		enums.AccessControlPermission_Read,
	}

	blockGroups, exception := s.blockGroupRepository.CheckPermissionsAndGetManyByIds(
		reqDto.Param.BlockGroupIds,
		reqDto.ContextFields.UserId,
		nil,
		allowedPermissions,
		options.WithDB(db),
		options.WithOnlyDeleted(types.Ternary_Negative),
	)
	if exception != nil {
		return nil, exception
	}

	var resDto dtos.GetMyBlockGroupsAndTheirBlocksByIdsResDto

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
	result := db.Model(&schemas.Block{}).
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
		} else {
			resDto[index].RawArborizedEditableBlock = dtos.RawArborizedEditableBlock{}
		}
	}

	return &resDto, nil
}

func (s *BlockGroupService) GetMyBlockGroupsAndTheirBlocksByBlockPackId(
	ctx context.Context, reqDto *dtos.GetMyBlockGroupsAndTheirBlocksByBlockPackIdReqDto,
) (*dtos.GetMyBlockGroupsAndTheirBlocksByBlockPackIdResDto, *exceptions.Exception) {
	if err := validation.Validator.Struct(reqDto); err != nil {
		return nil, exceptions.BlockGroup.InvalidDto().WithOrigin(err)
	}

	if len(reqDto.Param.BlockPackId) == 0 {
		return &dtos.GetMyBlockGroupsAndTheirBlocksByBlockPackIdResDto{}, nil
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
		if exceptions.CommonlyCompare(exceptions.BlockGroup.NotFound(), exception, false) {
			return &dtos.GetMyBlockGroupsAndTheirBlocksByBlockPackIdResDto{}, nil
		}
		return nil, exception
	}

	blockGroupIds := make([]uuid.UUID, len(blockGroups))
	nextBlockGroupMap := make(map[uuid.UUID]schemas.BlockGroup, len(blockGroups))
	var firstBlockGroup *schemas.BlockGroup = nil
	for index := range blockGroups {
		blockGroup := &blockGroups[index]
		blockGroupIds[index] = blockGroup.Id

		if blockGroup.PrevBlockGroupId == nil {
			if firstBlockGroup != nil {
				return nil, exceptions.BlockGroup.DuplicateBlockGroupsWithSamePrevBlockGroupId(reqDto.Param.BlockPackId)
			}
			firstBlockGroup = blockGroup
		} else {
			nextBlockGroupMap[*blockGroup.PrevBlockGroupId] = *blockGroup
		}
	}

	if firstBlockGroup == nil && len(blockGroups) > 0 {
		return nil, exceptions.BlockPack.NoRootBlockGroupInBlockPack(reqDto.Param.BlockPackId)
	}

	var resDto dtos.GetMyBlockGroupsAndTheirBlocksByBlockPackIdResDto
	currentBlockGroup := firstBlockGroup
	for currentBlockGroup != nil {
		resDto = append(resDto, dtos.GetMyBlockGroupAndItsBlocksByIdResDto{
			Id:                        currentBlockGroup.Id,
			BlockPackId:               currentBlockGroup.BlockPackId,
			PrevBlockGroupId:          currentBlockGroup.PrevBlockGroupId,
			SyncBlockGroupId:          currentBlockGroup.SyncBlockGroupId,
			DeletedAt:                 currentBlockGroup.DeletedAt,
			UpdatedAt:                 currentBlockGroup.UpdatedAt,
			CreatedAt:                 currentBlockGroup.CreatedAt,
			RawArborizedEditableBlock: dtos.RawArborizedEditableBlock{},
		})
		nextBlockGroup, exist := nextBlockGroupMap[currentBlockGroup.Id]
		if exist {
			currentBlockGroup = &nextBlockGroup
		} else {
			currentBlockGroup = nil
		}
	}

	if len(resDto) != len(blockGroups) {
		return nil, exceptions.BlockGroup.BrokenBlockGroupsLinkedListDetected(reqDto.Param.BlockPackId, blockGroupIds)
	}

	var flattenedBlocks []schemas.Block
	result := db.Model(&schemas.Block{}).
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
		} else {
			resDto[index].RawArborizedEditableBlock = dtos.RawArborizedEditableBlock{}
		}
	}

	return &resDto, nil
}

func (s *BlockGroupService) GetMyBlockGroupsByPrevBlockGroupId(
	ctx context.Context, reqDto *dtos.GetMyBlockGroupsByPrevBlockGroupIdReqDto,
) (*dtos.GetMyBlockGroupsByPrevBlockGroupIdResDto, *exceptions.Exception) {
	if err := validation.Validator.Struct(reqDto); err != nil {
		return nil, exceptions.BlockGroup.InvalidDto().WithOrigin(err)
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
		return nil, exceptions.BlockGroup.NotFound().WithOrigin(err)
	}

	return &resDto, nil
}

func (s *BlockGroupService) GetAllMyBlockGroupsByBlockPackId(
	ctx context.Context, reqDto *dtos.GetAllMyBlockGroupsByBlockPackIdReqDto,
) (*dtos.GetAllMyBlockGroupsByBlockPackIdResDto, *exceptions.Exception) {
	if err := validation.Validator.Struct(reqDto); err != nil {
		return nil, exceptions.BlockGroup.InvalidDto().WithOrigin(err)
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
		return nil, exceptions.BlockGroup.NotFound().WithOrigin(err)
	}

	return &resDto, nil
}

func (s *BlockGroupService) InsertBlockGroupByBlockPackId(
	ctx context.Context, reqDto *dtos.InsertBlockGroupByBlockPackIdReqDto,
) (*dtos.InsertBlockGroupByBlockPackIdResDto, *exceptions.Exception) {
	if err := validation.Validator.Struct(reqDto); err != nil {
		return nil, exceptions.BlockGroup.InvalidDto().WithOrigin(err)
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

	return &dtos.InsertBlockGroupByBlockPackIdResDto{
		Id:        *newBlockGroupId,
		CreatedAt: time.Now(),
	}, nil
}

func (s *BlockGroupService) InsertBlockGroupAndItsBlocksByBlockPackId(
	ctx context.Context, reqDto *dtos.InsertBlockGroupAndItsBlocksByBlockPackIdReqDto,
) (*dtos.InsertBlockGroupAndItsBlocksByBlockPackIdResDto, *exceptions.Exception) {
	if err := validation.Validator.Struct(reqDto); err != nil {
		return nil, exceptions.BlockGroup.InvalidDto().WithOrigin(err)
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
		return nil, exceptions.BlockGroup.FailedToCommitTransaction().WithOrigin(err)
	}

	return &dtos.InsertBlockGroupAndItsBlocksByBlockPackIdResDto{
		Id:        *newBlockGroupId,
		CreatedAt: time.Now(),
	}, nil
}

func (s *BlockGroupService) InsertBlockGroupsAndTheirBlocksByBlockPackId(
	ctx context.Context, reqDto *dtos.InsertBlockGroupsAndTheirBlocksByBlockPackIdReqDto,
) (*dtos.InsertBlockGroupsAndTheirBlocksByBlockPackIdResDto, *exceptions.Exception) {
	if err := validation.Validator.Struct(reqDto); err != nil {
		return nil, exceptions.BlockGroup.InvalidDto().WithOrigin(err)
	}

	// try transaction here, but not passing the transaction into the concurrency part
	// this may cause problem, if it doesn't work, maybe we should try transaction by our hands
	tx := s.db.WithContext(ctx).Begin()

	createBlockGroupsInput := make([]inputs.CreateBlockGroupInput, len(reqDto.Body.BlockGroupContents))
	validateBlockDto := make([]dtos.ArborizedEditableBlock, len(reqDto.Body.BlockGroupContents))
	for index, blockGroupContent := range reqDto.Body.BlockGroupContents {
		createBlockGroupsInput[index] = inputs.CreateBlockGroupInput{
			BlockGroupId:     blockGroupContent.BlockGroupId,
			PrevBlockGroupId: blockGroupContent.PrevBlockGroupId,
		}
		validateBlockDto[index] = blockGroupContent.ArborizedEditableBlock
	}

	// note that the order of the output newBlockGroupIds is the same as the order of reqDto.Body.BlockGroupContents
	newBlockGroupIds, exception := s.blockGroupRepository.InsertManyByBlockPackId(
		reqDto.Body.BlockPackId,
		reqDto.ContextFields.UserId,
		createBlockGroupsInput,
		options.WithTransactionDB(tx),
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
		rawFlattenedBlocks, exception := s.editableBlockAdapter.FlattenToRaw(&validateDto)
		if exception != nil {
			return rawFlattenedBlocks, exception.GetOrigin()
		}
		return rawFlattenedBlocks, nil
	}

	validateBlockResults := concurrency.Execute(
		validateBlockDto,
		min(10, max(len(validateBlockDto)/10, len(validateBlockDto)%10)),
		validateBlockFunc,
	)

	resDto := dtos.InsertBlockGroupsAndTheirBlocksByBlockPackIdResDto{
		IsAllSuccess:   true,
		FailedIndexes:  []int{},
		SuccessIndexes: []int{},
		SuccessBlockGroupAndBlockIds: []struct {
			BlockGroupId uuid.UUID   `json:"blockGroupId"`
			BlockIds     []uuid.UUID `json:"blockIds"`
		}{},
		CreatedAt: time.Now(),
	}
	var createBlockGroupContentInput []inputs.CreateBlockGroupContentInput
	for _, validateResult := range validateBlockResults {
		if validateResult.Err == nil {
			resDto.SuccessIndexes = append(resDto.SuccessIndexes, validateResult.Index)
			// note that since the order of newBlockGroupIds is the same as the reqDto
			// and here the concurrency job worker will output a result with index in Result.Index
			// which provide us enough ability to reorder and align the result with the block group ids here
			blockIds := make([]uuid.UUID, len(validateResult.Data))
			createBlockInputs := make([]inputs.CreateBlockInput, len(validateResult.Data))
			for index, rawFlattenedBlock := range validateResult.Data {
				blockIds[index] = rawFlattenedBlock.Id
				createBlockInputs[index] = inputs.CreateBlockInput{
					Id:            rawFlattenedBlock.Id,
					ParentBlockId: rawFlattenedBlock.ParentBlockId,
					Type:          rawFlattenedBlock.Type,
					Props:         rawFlattenedBlock.Props,
					Content:       rawFlattenedBlock.Content,
				}
			}
			resDto.SuccessBlockGroupAndBlockIds = append(resDto.SuccessBlockGroupAndBlockIds, struct {
				BlockGroupId uuid.UUID   `json:"blockGroupId"`
				BlockIds     []uuid.UUID `json:"blockIds"`
			}{
				BlockGroupId: newBlockGroupIds[validateResult.Index],
				BlockIds:     blockIds,
			})
			createBlockGroupContentInput = append(createBlockGroupContentInput, inputs.CreateBlockGroupContentInput{
				BlockGroupId: newBlockGroupIds[validateResult.Index],
				Blocks:       createBlockInputs,
			})
		} else {
			resDto.FailedIndexes = append(resDto.FailedIndexes, validateResult.Index)
			resDto.IsAllSuccess = false
		}
	}

	if len(createBlockGroupContentInput) == 0 {
		tx.Rollback()
		return nil, exceptions.BlockGroup.FailedToCreate().WithDetails("no valid block tree structure in any of the given block groups")
	}

	_, exception = s.blockRepository.CreateManyByBlockGroupIds(
		reqDto.ContextFields.UserId,
		createBlockGroupContentInput,
		options.WithTransactionDB(tx),
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
		return nil, exceptions.BlockGroup.FailedToCommitTransaction().WithOrigin(err)
	}

	resDto.CreatedAt = time.Now()
	return &resDto, nil
}

func (s *BlockGroupService) InsertSequentialBlockGroupsAndTheirBlocksByBlockPackId(
	ctx context.Context, reqDto *dtos.InsertSequentialBlockGroupsAndTheirBlocksByBlockPackIdReqDto,
) (*dtos.InsertSequentialBlockGroupsAndTheirBlocksByBlockPackIdResDto, *exceptions.Exception) {
	if err := validation.Validator.Struct(reqDto); err != nil {
		return nil, exceptions.BlockGroup.InvalidDto().WithOrigin(err)
	}

	// try transaction here, but not passing the transaction into the concurrency part
	// this may cause problem, if it doesn't work, maybe we should try transaction by our hands
	tx := s.db.WithContext(ctx).Begin()

	createBlockGroupsInput := make([]inputs.CreateBlockGroupInput, len(reqDto.Body.ArborizedEditableBlocks))
	var prevBlockGroupId *uuid.UUID = nil
	for index, _ := range reqDto.Body.ArborizedEditableBlocks {
		newBlockGroupId := uuid.New()
		createBlockGroupsInput[index] = inputs.CreateBlockGroupInput{
			BlockGroupId:     &newBlockGroupId,
			PrevBlockGroupId: prevBlockGroupId,
		}
		prevBlockGroupId = &newBlockGroupId
	}

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
		if exception != nil {
			return data, exception.GetOrigin()
		}
		return data, nil
	}

	validateBlockResults := concurrency.Execute(
		reqDto.Body.ArborizedEditableBlocks,
		min(10, max(len(reqDto.Body.ArborizedEditableBlocks)/10, len(reqDto.Body.ArborizedEditableBlocks)%10)),
		validateBlockFunc,
	)

	resDto := dtos.InsertSequentialBlockGroupsAndTheirBlocksByBlockPackIdResDto{
		IsAllSuccess:   true,
		FailedIndexes:  []int{},
		SuccessIndexes: []int{},
		SuccessBlockGroupAndBlockIds: []struct {
			BlockGroupId uuid.UUID   `json:"blockGroupId"`
			BlockIds     []uuid.UUID `json:"blockIds"`
		}{},
		CreatedAt: time.Now(),
	}

	var createBlocksInputs []inputs.CreateBlockGroupContentInput
	for _, validateResult := range validateBlockResults {
		if validateResult.Err == nil {
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
				BlockGroupId uuid.UUID   `json:"blockGroupId"`
				BlockIds     []uuid.UUID `json:"blockIds"`
			}{
				BlockGroupId: newBlockGroupIds[validateResult.Index],
				BlockIds:     blockIds,
			})
			createBlocksInputs = append(createBlocksInputs, inputs.CreateBlockGroupContentInput{
				BlockGroupId: newBlockGroupIds[validateResult.Index],
				Blocks:       createBlocksByBlockGroupInput,
			})
		} else {
			resDto.FailedIndexes = append(resDto.FailedIndexes, validateResult.Index)
			resDto.IsAllSuccess = false
		}
	}

	if len(createBlocksInputs) == 0 {
		tx.Rollback()
		return nil, exceptions.BlockGroup.FailedToCreate().WithDetails("no valid block tree structure in any of the given block groups")
	}

	_, exception = s.blockRepository.CreateManyByBlockGroupIds(
		reqDto.ContextFields.UserId,
		createBlocksInputs,
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
		return nil, exceptions.BlockGroup.FailedToCommitTransaction().WithOrigin(err)
	}

	resDto.CreatedAt = time.Now()
	return &resDto, nil
}

func (s *BlockGroupService) MoveMyBlockGroupById(
	ctx context.Context, reqDto *dtos.MoveMyBlockGroupByIdReqDto,
) (*dtos.MoveMyBlockGroupByIdResDto, *exceptions.Exception) {
	if err := validation.Validator.Struct(reqDto); err != nil {
		return nil, exceptions.BlockGroup.InvalidDto().WithOrigin(err)
	}

	if reqDto.Body.DestinationBlockGroupId == reqDto.Body.MovablePrevBlockGroupId {
		return nil, exceptions.BlockGroup.NoChanges()
	}

	if reqDto.Body.DestinationBlockGroupId != nil &&
		reqDto.Body.MovableBlockGroupId == *reqDto.Body.DestinationBlockGroupId {
		return nil, exceptions.BlockGroup.InvalidDto("Cannot move the block group to its own next position")
	}

	tx := s.db.WithContext(ctx).Begin()

	allowedPermissions := []enums.AccessControlPermission{
		enums.AccessControlPermission_Owner,
		enums.AccessControlPermission_Admin,
		enums.AccessControlPermission_Write,
	}

	movableBlockGroup, exception := s.blockGroupRepository.CheckPermissionAndGetOneById(
		reqDto.Body.MovableBlockGroupId,
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

	if movableBlockGroup.PrevBlockGroupId != reqDto.Body.MovablePrevBlockGroupId {
		tx.Rollback()
		prevBlockGroupIdString := "null"
		if movableBlockGroup.PrevBlockGroupId != nil {
			prevBlockGroupIdString = movableBlockGroup.PrevBlockGroupId.String()
		}
		movablePrevBlockGroupIdString := "null"
		if reqDto.Body.MovablePrevBlockGroupId != nil {
			movablePrevBlockGroupIdString = reqDto.Body.MovablePrevBlockGroupId.String()
		}
		return nil, exceptions.BlockGroup.InvalidDto(
			"The given block group has different previous block group id from the actual block group, expected %s, got %s",
			prevBlockGroupIdString,
			movablePrevBlockGroupIdString,
		)
	}

	if movableBlockGroup.BlockPackId != reqDto.Body.BlockPackId {
		tx.Rollback()
		return nil, exceptions.BlockGroup.InvalidDto("The given block group is not under the given block pack")
	}

	collapsedBlockGroup, _ := s.blockGroupRepository.GetOneByPrevBlockGroupId(
		reqDto.Body.BlockPackId,
		&movableBlockGroup.Id,
		reqDto.ContextFields.UserId,
		nil,
		options.WithDB(tx),
		options.WithOnlyDeleted(types.Ternary_Negative),
	)
	destinationNextBlockGroup, _ := s.blockGroupRepository.GetOneByPrevBlockGroupId(
		reqDto.Body.BlockPackId,
		reqDto.Body.DestinationBlockGroupId,
		reqDto.ContextFields.UserId,
		nil,
		options.WithDB(tx),
		options.WithOnlyDeleted(types.Ternary_Negative),
	)

	if _, exception = s.blockGroupRepository.UpdateOneById(
		movableBlockGroup.Id,
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

	if destinationNextBlockGroup != nil {
		if destinationNextBlockGroup.Id != movableBlockGroup.Id {
			if _, exception = s.blockGroupRepository.UpdateOneById(
				destinationNextBlockGroup.Id,
				reqDto.ContextFields.UserId,
				inputs.PartialUpdateBlockGroupInput{
					Values: inputs.UpdateBlockGroupInput{
						PrevBlockGroupId: &movableBlockGroup.Id,
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

	if collapsedBlockGroup != nil {
		if _, exception = s.blockGroupRepository.UpdateOneById(
			collapsedBlockGroup.Id,
			reqDto.ContextFields.UserId,
			inputs.PartialUpdateBlockGroupInput{
				Values: inputs.UpdateBlockGroupInput{
					PrevBlockGroupId: movableBlockGroup.PrevBlockGroupId,
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
		return nil, exceptions.BlockGroup.FailedToCommitTransaction().WithOrigin(err)
	}

	return &dtos.MoveMyBlockGroupByIdResDto{
		UpdatedAt: time.Now(),
	}, nil
}

func (s *BlockGroupService) MoveMyBlockGroupsByIds(
	ctx context.Context, reqDto *dtos.MoveMyBlockGroupsByIdsReqDto,
) (*dtos.MoveMyBlockGroupsByIdsResDto, *exceptions.Exception) {
	if err := validation.Validator.Struct(reqDto); err != nil {
		return nil, exceptions.BlockGroup.InvalidDto().WithOrigin(err)
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
	visited := make(map[*uuid.UUID]bool, len(movableBlockGroups)) // it is also guarantee that there's only one block group with the its PrevBlockGroupId equals to nil
	for index, movableBlockGroupId := range reqDto.Body.MovableBlockGroupIds {
		// 1. check repeated PrevBlockGroupId for circular block groups
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

	// note that we have guarantee that the below block groups MUST exist on the above procedure
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

	// Step 1: connect the start block group to the destination block group
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
		return nil, exceptions.BlockGroup.FailedToCommitTransaction().WithOrigin(err)
	}

	return &dtos.MoveMyBlockGroupsByIdsResDto{
		UpdatedAt: time.Now(),
	}, nil
}

func (s *BlockGroupService) BatchMoveMyBlockGroupsByIds(
	ctx context.Context, reqDto *dtos.BatchMoveMyBlockGroupsByIdsReqDto,
) (*dtos.BatchMoveMyBlockGroupsByIdsResDto, *exceptions.Exception) {
	if err := validation.Validator.Struct(reqDto); err != nil {
		return nil, exceptions.BlockGroup.InvalidDto().WithOrigin(err)
	}

	tx := s.db.WithContext(ctx).Begin()

	allowedPermissions := []enums.AccessControlPermission{
		enums.AccessControlPermission_Owner,
		enums.AccessControlPermission_Admin,
		enums.AccessControlPermission_Write,
	}

	blockPackIds := make([]uuid.UUID, len(reqDto.Body.MovedBlockGroups))
	movableBlockGroupIds := make([]uuid.UUID, len(reqDto.Body.MovedBlockGroups))
	var nonNullDestinationBlockGroupIds []uuid.UUID
	var nullDestinationBlockGroupBlockPackIds []uuid.UUID
	for index, movedBlockGroup := range reqDto.Body.MovedBlockGroups {
		blockPackIds[index] = movedBlockGroup.BlockPackId
		movableBlockGroupIds[index] = movedBlockGroup.MovableBlockGroupId
		if movedBlockGroup.DestinationBlockGroupId != nil {
			nonNullDestinationBlockGroupIds = append(nonNullDestinationBlockGroupIds, *movedBlockGroup.DestinationBlockGroupId)
		} else {
			nullDestinationBlockGroupBlockPackIds = append(nullDestinationBlockGroupBlockPackIds, movedBlockGroup.BlockPackId)
		}
	}

	validMovableBlockGroupIndexes := make(map[uuid.UUID]int)
	validMovableBlockGroups, exception := s.blockGroupRepository.CheckPermissionsAndGetManyByIds(
		movableBlockGroupIds,
		reqDto.ContextFields.UserId,
		[]schemas.BlockGroupRelation{
			schemas.BlockGroupRelation_NextBlockGroup,
		},
		allowedPermissions,
		options.WithTransactionDB(tx),
		options.WithOnlyDeleted(types.Ternary_Negative),
	)
	if exception != nil {
		tx.Rollback()
		return nil, exception
	}
	for index, validMovableBlockGroup := range validMovableBlockGroups {
		validMovableBlockGroupIndexes[validMovableBlockGroup.Id] = index
	}

	validNonNullDestinationBlockGroupIndexes := make(map[uuid.UUID]int)
	var validNonNullDestinationBlockGroups []schemas.BlockGroup
	if len(nonNullDestinationBlockGroupIds) > 0 {
		validNonNullDestinationBlockGroups, exception := s.blockGroupRepository.CheckPermissionsAndGetManyByIds(
			nonNullDestinationBlockGroupIds,
			reqDto.ContextFields.UserId,
			[]schemas.BlockGroupRelation{
				schemas.BlockGroupRelation_NextBlockGroup,
			},
			allowedPermissions,
			options.WithTransactionDB(tx),
			options.WithOnlyDeleted(types.Ternary_Negative),
		)
		if exception != nil {
			tx.Rollback()
			return nil, exception
		}
		for index, validDestinationBlockGroup := range validNonNullDestinationBlockGroups {
			validNonNullDestinationBlockGroupIndexes[validDestinationBlockGroup.Id] = index
		}
	}

	// get the first block groups of the block pack for the null destinations
	nullDestinationBlockPackIdToNextBlockGroupId := make(map[uuid.UUID]uuid.UUID)
	if len(nullDestinationBlockGroupBlockPackIds) > 0 {
		// the result is not validated by the permissions check yet,
		// but since we're getting the first block group by its prev block group id,
		// we can just check if the block pack is owned by the current user
		nullDestinationNextBlockGroups, exception := s.blockGroupRepository.GetManyByPrevBlockGroupIds(
			nullDestinationBlockGroupBlockPackIds,
			make([]*uuid.UUID, len(nullDestinationBlockGroupBlockPackIds)),
			reqDto.ContextFields.UserId,
			nil,
			options.WithTransactionDB(tx),
			options.WithOnlyDeleted(types.Ternary_Negative),
		)
		if exception != nil {
			tx.Rollback()
			return nil, exception
		}
		for _, nullDestinationBlockGroup := range nullDestinationNextBlockGroups {
			nullDestinationBlockPackIdToNextBlockGroupId[nullDestinationBlockGroup.BlockPackId] = nullDestinationBlockGroup.Id
		}
	}

	var nullUUID *uuid.UUID = nil
	var valuePlaceholders []string
	var valueArgs []interface{}
	for _, movedBlockGroups := range reqDto.Body.MovedBlockGroups {
		movableIndex, exist := validMovableBlockGroupIndexes[movedBlockGroups.MovableBlockGroupId]
		if !exist {
			continue
		}

		var movable schemas.BlockGroup = validMovableBlockGroups[movableIndex]
		var destination *schemas.BlockGroup = nil

		if movedBlockGroups.DestinationBlockGroupId != nil {
			destinationIndex, exist := validNonNullDestinationBlockGroupIndexes[*movedBlockGroups.DestinationBlockGroupId]
			if !exist || destinationIndex >= len(validNonNullDestinationBlockGroups) {
				continue
			}
			destination = &validNonNullDestinationBlockGroups[destinationIndex]
		} else if _, exist := nullDestinationBlockPackIdToNextBlockGroupId[movedBlockGroups.BlockPackId]; !exist {
			continue
		}

		if movable.BlockPackId != movedBlockGroups.BlockPackId ||
			(destination != nil && destination.BlockPackId != movedBlockGroups.BlockPackId) {
			continue
		}

		// Step 1 (Optional): Restore the block groups linked list on the next block group of the movable block group
		if movable.NextBlockGroup != nil { // only restore it when the next block group exist
			valuePlaceholders = append(valuePlaceholders, "(?::uuid, ?::uuid)")
			valueArgs = append(valueArgs,
				movable.NextBlockGroup.Id,
				movable.PrevBlockGroupId,
			)
		}

		// Step 2: Prepare for the insertion by linking the next block group of the destination block group to pointing to the inserted block group
		valuePlaceholders = append(valuePlaceholders, "(?::uuid, ?::uuid)")
		if destination != nil && destination.NextBlockGroup != nil {
			valueArgs = append(valueArgs,
				destination.NextBlockGroup.Id,
				movable.Id,
			)
		} else {
			valueArgs = append(valueArgs,
				nullDestinationBlockPackIdToNextBlockGroupId[movedBlockGroups.BlockPackId],
				movable.Id,
			)
		}

		// Step 3: Insert the movable block group to the destination by make it point to the destination block group
		valuePlaceholders = append(valuePlaceholders, "(?::uuid, ?::uuid)")
		if destination != nil {
			valueArgs = append(valueArgs,
				movable.Id,
				destination.Id,
			)
		} else {
			valueArgs = append(valueArgs,
				movable.Id,
				nullUUID,
			)
		}
	}

	valuePlaceholdersStr := strings.Join(valuePlaceholders, ",")

	// Pass 1 First set the deleted_at to now
	fakeDeleteSql := fmt.Sprintf(`
		UPDATE "BlockGroupTable" AS bg
		SET
			deleted_at = NOW()
		FROM (VALUES %s) AS v(id, prev_block_group_id)
		WHERE bg.id = v.id::uuid AND bg.deleted_at IS NULL
	`, valuePlaceholdersStr)
	result := tx.Exec(fakeDeleteSql, valueArgs...)
	if exception := exceptions.Cover(nil, []types.Pair[bool, *exceptions.Exception]{
		{First: result.Error != nil, Second: exceptions.BlockGroup.FailedToUpdate().WithOrigin(result.Error)},
		{First: result.RowsAffected == 0, Second: exceptions.BlockGroup.NoChanges()},
	}); exception != nil {
		tx.Rollback()
		return nil, exception
	}

	// Pass 2 Update the prev block group id to the new one
	sql := fmt.Sprintf(`
		UPDATE "BlockGroupTable" AS bg
		SET
			prev_block_group_id = v.prev_block_group_id::uuid, 
			updated_at = NOW()
		FROM (VALUES %s) AS v(id, prev_block_group_id)
		WHERE bg.id = v.id::uuid AND bg.deleted_at IS NOT NULL
	`, valuePlaceholdersStr) // update its deleted_at to now temporary to avoid the unique constraints
	result = tx.Exec(sql, valueArgs...)
	if exception := exceptions.Cover(nil, []types.Pair[bool, *exceptions.Exception]{
		{First: result.Error != nil, Second: exceptions.BlockGroup.FailedToUpdate().WithOrigin(result.Error)},
		{First: result.RowsAffected == 0, Second: exceptions.BlockGroup.NoChanges()},
	}); exception != nil {
		tx.Rollback()
		return nil, exception
	}

	// Pass 3 Restore the deleted_at back to null
	restoreSql := fmt.Sprintf(`
		UPDATE "BlockGroupTable" AS bg
		SET
			deleted_at = NULL
		FROM (VALUES %s) AS v(id, prev_block_group_id)
		WHERE bg.id = v.id::uuid AND bg.deleted_at IS NOT NULL
	`, valuePlaceholdersStr)
	result = tx.Exec(restoreSql, valueArgs...)
	if exception := exceptions.Cover(nil, []types.Pair[bool, *exceptions.Exception]{
		{First: result.Error != nil, Second: exceptions.BlockGroup.FailedToUpdate().WithOrigin(result.Error)},
		{First: result.RowsAffected == 0, Second: exceptions.BlockGroup.NoChanges()},
	}); exception != nil {
		tx.Rollback()
		return nil, exception
	}

	if err := tx.Commit().Error; err != nil {
		return nil, exceptions.BlockGroup.FailedToCommitTransaction().WithOrigin(err)
	}

	return &dtos.BatchMoveMyBlockGroupsByIdsResDto{
		UpdatedAt: time.Now(),
	}, nil
}

func (s *BlockGroupService) RestoreMyBlockGroupById(
	ctx context.Context, reqDto *dtos.RestoreMyBlockGroupByIdReqDto,
) (*dtos.RestoreMyBlockGroupByIdResDto, *exceptions.Exception) {
	if err := validation.Validator.Struct(reqDto); err != nil {
		return nil, exceptions.BlockGroup.InvalidDto().WithOrigin(err)
	}

	db := s.db.WithContext(ctx)

	restoredBlockGroup, exception := s.blockGroupRepository.RestoreSoftDeletedOneById(
		reqDto.Body.BlockGroupId,
		reqDto.ContextFields.UserId,
		options.WithDB(db),
	)
	if exception != nil {
		return nil, exception
	}

	return &dtos.RestoreMyBlockGroupByIdResDto{
		Id:               restoredBlockGroup.Id,
		BlockPackId:      restoredBlockGroup.BlockPackId,
		PrevBlockGroupId: restoredBlockGroup.PrevBlockGroupId,
		SyncBlockGroupId: restoredBlockGroup.SyncBlockGroupId,
		Size:             restoredBlockGroup.Size,
		DeletedAt:        restoredBlockGroup.DeletedAt,
		UpdatedAt:        restoredBlockGroup.UpdatedAt,
		CreatedAt:        restoredBlockGroup.CreatedAt,
	}, nil
}

func (s *BlockGroupService) RestoreMyBlockGroupsByIds(
	ctx context.Context, reqDto *dtos.RestoreMyBlockGroupsByIdsReqDto,
) (*dtos.RestoreMyBlockGroupsByIdsResDto, *exceptions.Exception) {
	if err := validation.Validator.Struct(reqDto); err != nil {
		return nil, exceptions.BlockGroup.InvalidDto().WithOrigin(err)
	}

	db := s.db.WithContext(ctx)

	restoredBlockGroups, exception := s.blockGroupRepository.RestoreSoftDeletedManyByIds(
		reqDto.Body.BlockGroupIds,
		reqDto.ContextFields.UserId,
		options.WithDB(db),
	)
	if exception != nil {
		return nil, exception
	}

	resDto := make(dtos.RestoreMyBlockGroupsByIdsResDto, len(restoredBlockGroups))
	for index, restoredBlockGroup := range restoredBlockGroups {
		resDto[index] = dtos.RestoreMyBlockGroupByIdResDto{
			Id:               restoredBlockGroup.Id,
			BlockPackId:      restoredBlockGroup.BlockPackId,
			PrevBlockGroupId: restoredBlockGroup.PrevBlockGroupId,
			SyncBlockGroupId: restoredBlockGroup.SyncBlockGroupId,
			Size:             restoredBlockGroup.Size,
			DeletedAt:        restoredBlockGroup.DeletedAt,
			UpdatedAt:        restoredBlockGroup.UpdatedAt,
			CreatedAt:        restoredBlockGroup.CreatedAt,
		}
	}

	return &resDto, nil
}

func (s *BlockGroupService) DeleteMyBlockGroupById(
	ctx context.Context, reqDto *dtos.DeleteMyBlockGroupByIdReqDto,
) (*dtos.DeleteMyBlockGroupByIdResDto, *exceptions.Exception) {
	if err := validation.Validator.Struct(reqDto); err != nil {
		return nil, exceptions.BlockGroup.InvalidDto().WithOrigin(err)
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
		return nil, exceptions.BlockGroup.InvalidDto().WithOrigin(err)
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
