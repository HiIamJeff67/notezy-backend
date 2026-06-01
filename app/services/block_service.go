package services

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"time"

	"github.com/google/uuid"
	"gorm.io/datatypes"
	"gorm.io/gorm"

	adapters "github.com/HiIamJeff67/notezy-backend/app/adapters"
	dtos "github.com/HiIamJeff67/notezy-backend/app/dtos"
	exceptions "github.com/HiIamJeff67/notezy-backend/app/exceptions"
	inputs "github.com/HiIamJeff67/notezy-backend/app/models/inputs"
	repositories "github.com/HiIamJeff67/notezy-backend/app/models/repositories"
	schemas "github.com/HiIamJeff67/notezy-backend/app/models/schemas"
	enums "github.com/HiIamJeff67/notezy-backend/app/models/schemas/enums"
	options "github.com/HiIamJeff67/notezy-backend/app/options"
	validation "github.com/HiIamJeff67/notezy-backend/app/validation"
	constants "github.com/HiIamJeff67/notezy-backend/shared/constants"
	blocknote "github.com/HiIamJeff67/notezy-backend/shared/lib/blocknote"
	concurrency "github.com/HiIamJeff67/notezy-backend/shared/lib/concurrency"
	types "github.com/HiIamJeff67/notezy-backend/shared/types"
)

type BlockServiceInterface interface {
	GetMyBlockById(ctx context.Context, reqDto *dtos.GetMyBlockByIdReqDto) (*dtos.GetMyBlockByIdResDto, *exceptions.Exception)
	GetMyBlocksByIds(ctx context.Context, reqDto *dtos.GetMyBlocksByIdsReqDto) (*dtos.GetMyBlocksByIdsResDto, *exceptions.Exception)
	GetMyBlocksByBlockGroupId(ctx context.Context, reqDto *dtos.GetMyBlocksByBlockGroupIdReqDto) (*dtos.GetMyBlocksByBlockGroupIdResDto, *exceptions.Exception)
	GetMyBlocksByBlockGroupIds(ctx context.Context, reqDto *dtos.GetMyBlocksByBlockGroupIdsReqDto) (*dtos.GetMyBlocksByBlockGroupIdsResDto, *exceptions.Exception)
	GetMyBlocksByBlockPackId(ctx context.Context, reqDto *dtos.GetMyBlocksByBlockPackIdReqDto) (*dtos.GetMyBlocksByBlockPackIdResDto, *exceptions.Exception)
	GetAllMyBlocks(ctx context.Context, reqDto *dtos.GetAllMyBlocksReqDto) (*dtos.GetAllMyBlocksResDto, *exceptions.Exception)
	InsertBlock(ctx context.Context, reqDto *dtos.InsertBlockReqDto) (*dtos.InsertBlockResDto, *exceptions.Exception)
	InsertBlocks(ctx context.Context, reqDto *dtos.InsertBlocksReqDto) (*dtos.InsertBlocksResDto, *exceptions.Exception)
	UpdateMyBlockById(ctx context.Context, reqDto *dtos.UpdateMyBlockByIdReqDto) (*dtos.UpdateMyBlockByIdResDto, *exceptions.Exception)
	UpdateMyBlocksByIds(ctx context.Context, reqDto *dtos.UpdateMyBlocksByIdsReqDto) (*dtos.UpdateMyBlocksByIdsResDto, *exceptions.Exception)
	RestoreMyBlockById(ctx context.Context, reqDto *dtos.RestoreMyBlockByIdReqDto) (*dtos.RestoreMyBlockByIdResDto, *exceptions.Exception)
	RestoreMyBlocksByIds(ctx context.Context, reqDto *dtos.RestoreMyBlocksByIdsReqDto) (*dtos.RestoreMyBlocksByIdsResDto, *exceptions.Exception)
	DeleteMyBlockById(ctx context.Context, reqDto *dtos.DeleteMyBlockByIdReqDto) (*dtos.DeleteMyBlockByIdResDto, *exceptions.Exception)
	DeleteMyBlocksByIds(ctx context.Context, reqDto *dtos.DeleteMyBlocksByIdsReqDto) (*dtos.DeleteMyBlockPacksByIdsResDto, *exceptions.Exception)
}

type BlockService struct {
	db                   *gorm.DB
	blockPackRepository  repositories.BlockPackRepositoryInterface
	blockGroupRepository repositories.BlockGroupRepositoryInterface
	blockRepository      repositories.BlockRepositoryInterface
	editableBlockAdapter adapters.EditableBlockAdapterInterface
}

func NewBlockService(
	db *gorm.DB,
	blockPackRepository repositories.BlockPackRepositoryInterface,
	blockGroupRepository repositories.BlockGroupRepositoryInterface,
	blockRepository repositories.BlockRepositoryInterface,
	editableBlockAdapter adapters.EditableBlockAdapterInterface,
) BlockServiceInterface {
	return &BlockService{
		db:                   db,
		blockPackRepository:  blockPackRepository,
		blockGroupRepository: blockGroupRepository,
		blockRepository:      blockRepository,
		editableBlockAdapter: editableBlockAdapter,
	}
}

/* ============================== Service Methods for Block ============================== */

func (s *BlockService) GetMyBlockById(
	ctx context.Context, reqDto *dtos.GetMyBlockByIdReqDto,
) (*dtos.GetMyBlockByIdResDto, *exceptions.Exception) {
	if err := validation.Validator.Struct(reqDto); err != nil {
		return nil, exceptions.Block.InvalidDto().WithOrigin(err)
	}

	db := s.db.WithContext(ctx)

	block, exception := s.blockRepository.GetOneById(
		reqDto.Param.BlockId,
		reqDto.ContextFields.UserId,
		nil,
		options.WithDB(db),
		options.WithOnlyDeleted(types.Ternary_Negative),
	)
	if exception != nil {
		return nil, exception
	}

	return &dtos.GetMyBlockByIdResDto{
		Id:            block.Id,
		ParentBlockId: block.ParentBlockId,
		BlockGroupId:  block.BlockGroupId,
		Type:          block.Type,
		Props:         block.Props,
		Content:       block.Content,
		DeletedAt:     block.DeletedAt,
		UpdatedAt:     block.UpdatedAt,
		CreatedAt:     block.CreatedAt,
	}, nil
}

func (s *BlockService) GetMyBlocksByIds(
	ctx context.Context, reqDto *dtos.GetMyBlocksByIdsReqDto,
) (*dtos.GetMyBlocksByIdsResDto, *exceptions.Exception) {
	if err := validation.Validator.Struct(reqDto); err != nil {
		return nil, exceptions.Block.InvalidDto().WithOrigin(err)
	}

	db := s.db.WithContext(ctx)

	allowedPermissions := []enums.AccessControlPermission{
		enums.AccessControlPermission_Owner,
		enums.AccessControlPermission_Admin,
		enums.AccessControlPermission_Write,
		enums.AccessControlPermission_Read,
	}

	blocks, exception := s.blockRepository.CheckPermissionsAndGetManyByIds(
		reqDto.Param.BlockIds,
		reqDto.ContextFields.UserId,
		nil,
		allowedPermissions,
		options.WithDB(db),
		options.WithOnlyDeleted(types.Ternary_Negative),
	)
	if exception != nil {
		return nil, exception
	}

	var resDto dtos.GetMyBlocksByIdsResDto
	for _, block := range blocks {
		resDto = append(resDto, dtos.GetMyBlockByIdResDto{
			Id:            block.Id,
			ParentBlockId: block.ParentBlockId,
			BlockGroupId:  block.BlockGroupId,
			Type:          block.Type,
			Props:         block.Props,
			Content:       block.Content,
			DeletedAt:     block.DeletedAt,
			UpdatedAt:     block.UpdatedAt,
			CreatedAt:     block.CreatedAt,
		})
	}

	return &resDto, nil
}

func (s *BlockService) GetMyBlocksByBlockGroupId(
	ctx context.Context, reqDto *dtos.GetMyBlocksByBlockGroupIdReqDto,
) (*dtos.GetMyBlocksByBlockGroupIdResDto, *exceptions.Exception) {
	if err := validation.Validator.Struct(reqDto); err != nil {
		return nil, exceptions.Block.InvalidDto().WithOrigin(err)
	}

	db := s.db.WithContext(ctx)

	allowedPermissions := []enums.AccessControlPermission{
		enums.AccessControlPermission_Owner,
		enums.AccessControlPermission_Admin,
		enums.AccessControlPermission_Write,
		enums.AccessControlPermission_Read,
	}

	if !s.blockGroupRepository.HasPermission(
		reqDto.Param.BlockGroupId,
		reqDto.ContextFields.UserId,
		allowedPermissions,
		options.WithDB(db),
		options.WithOnlyDeleted(types.Ternary_Negative),
	) {
		return nil, exceptions.Block.NoPermission("get the block group of blocks")
	}

	var blocks []schemas.Block
	result := db.Model(&schemas.Block{}).
		Where("block_group_id = ?", reqDto.Param.BlockGroupId).
		Find(&blocks)
	if err := result.Error; err != nil || len(blocks) == 0 {
		if err := db.Commit().Error; err != nil {
			return nil, exceptions.Block.FailedToCommitTransaction().WithOrigin(err)
		}
		return &dtos.GetMyBlocksByBlockGroupIdResDto{
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

	rawArborizedBlock, _, exception := s.editableBlockAdapter.ArborizeRawToRaw(root, childrenMap)
	if exception != nil {
		return nil, exception
	}

	return &dtos.GetMyBlocksByBlockGroupIdResDto{
		RawArborizedEditableBlock: *rawArborizedBlock,
	}, nil
}

func (s *BlockService) GetMyBlocksByBlockGroupIds(
	ctx context.Context, reqDto *dtos.GetMyBlocksByBlockGroupIdsReqDto,
) (*dtos.GetMyBlocksByBlockGroupIdsResDto, *exceptions.Exception) {
	if err := validation.Validator.Struct(reqDto); err != nil {
		return nil, exceptions.Block.InvalidDto().WithOrigin(err)
	}

	db := s.db.WithContext(ctx)

	allowedPermissions := []enums.AccessControlPermission{
		enums.AccessControlPermission_Owner,
		enums.AccessControlPermission_Admin,
		enums.AccessControlPermission_Write,
		enums.AccessControlPermission_Read,
	}

	if s.blockGroupRepository.HavePermissions(
		reqDto.Param.BlockGroupIds,
		reqDto.ContextFields.UserId,
		allowedPermissions,
		options.WithDB(db),
		options.WithOnlyDeleted(types.Ternary_Negative),
	) {
		return nil, exceptions.Block.NoPermission("get the block groups of blocks")
	}

	var resDto dtos.GetMyBlocksByBlockGroupIdsResDto

	var flattenedBlocks []schemas.Block
	result := db.Model(&schemas.Block{}).
		Where("block_group_id IN ?", reqDto.Param.BlockGroupIds).
		Find(&flattenedBlocks)
	if err := result.Error; err != nil || len(flattenedBlocks) == 0 {
		return &resDto, nil
	}

	blockGroupToBlocksMap := make(map[uuid.UUID][]schemas.Block)
	for _, flattenedBlock := range flattenedBlocks {
		blockGroupToBlocksMap[flattenedBlock.BlockGroupId] = append(blockGroupToBlocksMap[flattenedBlock.BlockGroupId], flattenedBlock)
	}

	for index, blockGroupId := range reqDto.Param.BlockGroupIds {
		blocks, exist := blockGroupToBlocksMap[blockGroupId]
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

		rawArborizedBlock, _, exception := s.editableBlockAdapter.ArborizeRawToRaw(root, childrenMap)
		if exception != nil {
			return nil, exception
		}

		if rawArborizedBlock != nil {
			resDto[index].RawArborizedEditableBlock = *rawArborizedBlock
		}
	}

	return &resDto, nil
}

func (s *BlockService) GetMyBlocksByBlockPackId(
	ctx context.Context, reqDto *dtos.GetMyBlocksByBlockPackIdReqDto,
) (*dtos.GetMyBlocksByBlockPackIdResDto, *exceptions.Exception) {
	if err := validation.Validator.Struct(reqDto); err != nil {
		return nil, exceptions.Block.InvalidDto().WithOrigin(err)
	}

	db := s.db.WithContext(ctx)

	allowedPermissions := []enums.AccessControlPermission{
		enums.AccessControlPermission_Owner,
		enums.AccessControlPermission_Admin,
		enums.AccessControlPermission_Write,
		enums.AccessControlPermission_Read,
	}

	if !s.blockPackRepository.HasPermission(
		reqDto.Param.BlockPackId,
		reqDto.ContextFields.UserId,
		allowedPermissions,
		options.WithDB(db),
		options.WithOnlyDeleted(types.Ternary_Negative),
	) {
		return nil, exceptions.Block.NoPermission("get the block pack of blocks")
	}

	var blocks []schemas.Block
	result := db.Model(&schemas.Block{}).
		Joins("INNER JOIN \"BlockGroupTable\" bg ON \"BlockTable\".block_group_id = bg.id").
		Where("bg.block_pack_id = ?", reqDto.Param.BlockPackId).
		Find(&blocks)
	if err := result.Error; err != nil {
		return nil, exceptions.Block.NotFound().WithOrigin(err)
	}

	var resDto dtos.GetMyBlocksByBlockPackIdResDto
	for _, block := range blocks {
		resDto = append(resDto, dtos.GetMyBlockByIdResDto{
			Id:            block.Id,
			ParentBlockId: block.ParentBlockId,
			BlockGroupId:  block.BlockGroupId,
			Type:          block.Type,
			Props:         block.Props,
			Content:       block.Content,
			DeletedAt:     block.DeletedAt,
			UpdatedAt:     block.UpdatedAt,
			CreatedAt:     block.CreatedAt,
		})
	}

	return &resDto, nil
}

func (s *BlockService) GetAllMyBlocks(
	ctx context.Context, reqDto *dtos.GetAllMyBlocksReqDto,
) (*dtos.GetAllMyBlocksResDto, *exceptions.Exception) {
	if err := validation.Validator.Struct(reqDto); err != nil {
		return nil, exceptions.Block.InvalidDto().WithOrigin(err)
	}

	db := s.db.WithContext(ctx)

	// since we're getting blocks with the owner id of block group, there's no need to check the permission of the owner
	var resDto dtos.GetAllMyBlocksResDto
	result := db.Model(&schemas.Block{}).
		Joins("LEFT JOIN \"BlockGroupTable\" bg ON bg.id = block_group_id").
		Where("bg.owner_id = ?", reqDto.ContextFields.UserId).
		Find(&resDto)
	if err := result.Error; err != nil {
		return nil, exceptions.Block.NotFound()
	}

	return &resDto, nil
}

func (s *BlockService) InsertBlock(
	ctx context.Context, reqDto *dtos.InsertBlockReqDto,
) (*dtos.InsertBlockResDto, *exceptions.Exception) {
	if err := validation.Validator.Struct(reqDto); err != nil {
		return nil, exceptions.Block.InvalidDto().WithOrigin(err)
	}

	rawFlattenedBlocks, totalSize, exception := s.editableBlockAdapter.FlattenToRaw(&reqDto.Body.ArborizedEditableBlock)
	if exception != nil {
		return nil, exception
	}

	tx := s.db.WithContext(ctx).Begin()

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
		reqDto.Body.BlockGroupId,
		reqDto.ContextFields.UserId,
		input,
		options.WithTransactionDB(tx),
		options.WithBatchSize(constants.MaxBatchCreateBlockSize),
		options.WithOnlyDeleted(types.Ternary_Negative),
	)
	if exception != nil {
		tx.Rollback()
		return nil, exception
	}

	if totalSize > 0 {
		if _, exception = s.blockGroupRepository.IncrementSizeById(
			reqDto.Body.BlockGroupId,
			reqDto.ContextFields.UserId,
			totalSize,
			options.WithTransactionDB(tx),
			options.WithOnlyDeleted(types.Ternary_Negative),
		); exception != nil {
			tx.Rollback()
			return nil, exception
		}
	}

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return nil, exceptions.Block.FailedToCommitTransaction().WithOrigin(err)
	}

	return &dtos.InsertBlockResDto{
		CreatedAt: time.Now(),
	}, nil
}

func (s *BlockService) InsertBlocks(
	ctx context.Context, reqDto *dtos.InsertBlocksReqDto,
) (*dtos.InsertBlocksResDto, *exceptions.Exception) {
	if err := validation.Validator.Struct(reqDto); err != nil {
		return nil, exceptions.Block.InvalidDto().WithOrigin(err)
	}

	type ValidateBlockDto struct {
		ParentBlockId          *uuid.UUID
		ArborizedEditableBlock dtos.ArborizedEditableBlock
	}
	type FlattenedBlocksWithSize struct {
		Blocks    []dtos.RawFlattenedEditableBlock
		TotalSize int64
	}
	validateBlockDto := make([]ValidateBlockDto, len(reqDto.Body.InsertedBlocks))
	for index, insertedBlock := range reqDto.Body.InsertedBlocks {
		validateBlockDto[index] = ValidateBlockDto{
			ParentBlockId:          insertedBlock.ParentBlockId,
			ArborizedEditableBlock: insertedBlock.ArborizedEditableBlock,
		}
	}

	validateBlockFunc := func(validateBlockDto ValidateBlockDto) (FlattenedBlocksWithSize, error) {
		rawFlattenedBlocks, totalSize, exception := s.editableBlockAdapter.FlattenToRaw(&validateBlockDto.ArborizedEditableBlock)
		if exception != nil {
			return FlattenedBlocksWithSize{}, exception.GetOrigin()
		}

		if len(rawFlattenedBlocks) > 0 {
			rawFlattenedBlocks[0].ParentBlockId = validateBlockDto.ParentBlockId
		}
		return FlattenedBlocksWithSize{
			Blocks:    rawFlattenedBlocks,
			TotalSize: totalSize,
		}, nil
	}

	validateBlockResults := concurrency.Execute(
		validateBlockDto,
		20,
		validateBlockFunc,
	)

	resDto := dtos.InsertBlocksResDto{
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
	sizeDeltaByBlockGroupId := make(map[uuid.UUID]int64)
	for _, validateResult := range validateBlockResults {
		if validateResult.Err == nil {
			resDto.SuccessIndexes = append(resDto.SuccessIndexes, validateResult.Index)
			blockIds := make([]uuid.UUID, len(validateResult.Data.Blocks))
			createBlockInputs := make([]inputs.CreateBlockInput, len(validateResult.Data.Blocks))
			for index, rawFlattenedBlock := range validateResult.Data.Blocks {
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
				BlockGroupId: reqDto.Body.InsertedBlocks[validateResult.Index].BlockGroupId,
				BlockIds:     blockIds,
			})
			if validateResult.Data.TotalSize > 0 { // if we increase the size with 0 by using IncrementSizesByIds repository function, it will yield an exception
				sizeDeltaByBlockGroupId[reqDto.Body.InsertedBlocks[validateResult.Index].BlockGroupId] += validateResult.Data.TotalSize
			}
			createBlockGroupContentInput = append(createBlockGroupContentInput, inputs.CreateBlockGroupContentInput{
				BlockGroupId: reqDto.Body.InsertedBlocks[validateResult.Index].BlockGroupId,
				Blocks:       createBlockInputs,
			})
		} else {
			resDto.FailedIndexes = append(resDto.FailedIndexes, validateResult.Index)
			resDto.IsAllSuccess = false
		}
	}

	if len(createBlockGroupContentInput) == 0 {
		return nil, exceptions.BlockGroup.FailedToCreate().WithDetails("no valid block tree structure in any of the given block groups")
	}

	tx := s.db.WithContext(ctx).Begin()

	_, exception := s.blockRepository.CreateManyByBlockGroupIds(
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

	if len(sizeDeltaByBlockGroupId) > 0 { // if we increase the size with 0, it will yield an exception
		if exception = s.blockGroupRepository.IncrementSizesByIds(
			reqDto.ContextFields.UserId,
			sizeDeltaByBlockGroupId,
			options.WithTransactionDB(tx),
			options.WithOnlyDeleted(types.Ternary_Negative),
		); exception != nil {
			tx.Rollback()
			return nil, exception
		}
	}

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return nil, exceptions.Block.FailedToCommitTransaction().WithOrigin(err)
	}

	resDto.CreatedAt = time.Now()
	return &resDto, nil
}

func (s *BlockService) UpdateMyBlockById(
	ctx context.Context, reqDto *dtos.UpdateMyBlockByIdReqDto,
) (*dtos.UpdateMyBlockByIdResDto, *exceptions.Exception) {
	if err := validation.Validator.Struct(reqDto); err != nil {
		return nil, exceptions.Block.InvalidDto().WithOrigin(err)
	}

	tx := s.db.WithContext(ctx).Begin()

	allowedPermissions := []enums.AccessControlPermission{
		enums.AccessControlPermission_Owner,
		enums.AccessControlPermission_Admin,
		enums.AccessControlPermission_Write,
	}

	block, exception := s.blockRepository.CheckPermissionAndGetOneById(
		reqDto.Body.BlockId,
		reqDto.ContextFields.UserId,
		nil,
		allowedPermissions,
		options.WithTransactionDB(tx),
		options.WithOnlyDeleted(types.Ternary_Negative),
	)
	if exception != nil {
		tx.Rollback()
		return nil, exception
	}

	updateInput := inputs.PartialUpdateBlockInput{
		Values: inputs.UpdateBlockInput{
			ParentBlockId: reqDto.Body.Values.ParentBlockId,
			BlockGroupId:  reqDto.Body.Values.BlockGroupId,
			Type:          reqDto.Body.Values.Type,
			Props:         nil,
			Content:       nil,
		},
		SetNull: reqDto.Body.SetNull,
	}
	oldBlockGroupId := block.BlockGroupId
	newBlockGroupId := oldBlockGroupId
	if reqDto.Body.Values.BlockGroupId != nil {
		newBlockGroupId = *reqDto.Body.Values.BlockGroupId
	}

	oldSize := int64(len(block.Props)) + int64(len(block.Content))
	newProps := block.Props
	newContent := block.Content

	if reqDto.Body.Values.Props != nil {
		propsString := string(bytes.TrimSpace(*reqDto.Body.Values.Props))
		if propsString == "{}" || propsString == "" {
			emptyPropsJson := datatypes.JSON("{}")
			updateInput.Values.Props = &emptyPropsJson
			newProps = emptyPropsJson
		} else {
			_, err := blocknote.ParseProps(block.Type.String(), *reqDto.Body.Values.Props)
			if err != nil {
				tx.Rollback()
				return nil, exceptions.Block.InvalidDto().WithOrigin(err)
			}
			rawPropsJson := datatypes.JSON(*reqDto.Body.Values.Props)
			updateInput.Values.Props = &rawPropsJson
			newProps = rawPropsJson
		}
	}

	if reqDto.Body.Values.Content != nil {
		trimContent := bytes.TrimSpace(*reqDto.Body.Values.Content)
		trimContentString := string(trimContent)
		if trimContentString == "null" || trimContentString == "[]" || trimContentString == "" {
			emptyContentsJson := datatypes.JSON("[]")
			updateInput.Values.Content = &emptyContentsJson
			newContent = emptyContentsJson
		} else {
			switch trimContent[0] {
			case '[':
				var list blocknote.InlineContentList
				if err := json.Unmarshal(trimContent, &list); err != nil {
					tx.Rollback()
					return nil, exceptions.Block.InvalidDto().WithOrigin(err)
				}
				rawContentJson := datatypes.JSON(*reqDto.Body.Values.Content)
				updateInput.Values.Content = &rawContentJson
				newContent = rawContentJson
			case '{':
				var table blocknote.TableContent
				if err := json.Unmarshal(trimContent, &table); err != nil {
					tx.Rollback()
					return nil, exceptions.Block.InvalidDto().WithOrigin(err)
				}
				rawContentJson := datatypes.JSON(*reqDto.Body.Values.Content)
				updateInput.Values.Content = &rawContentJson
				newContent = rawContentJson
			default:
				tx.Rollback()
				return nil, exceptions.Block.InvalidDto().WithOrigin(errors.New("invalid content format: must be array or object"))
			}
		}
	}

	updatedBlock, exception := s.blockRepository.UpdateOneById(
		reqDto.Body.BlockId,
		reqDto.ContextFields.UserId,
		updateInput,
		options.WithTransactionDB(tx),
		options.WithOnlyDeleted(types.Ternary_Negative),
		options.WithSkipPermissionCheck(),
	)
	if exception != nil {
		tx.Rollback()
		return nil, exception
	}

	newSize := int64(len(newProps)) + int64(len(newContent))
	blockGroupIdToSizeDelta := make(map[uuid.UUID]int64)
	var affectedBlockGroupIds []uuid.UUID

	if newBlockGroupId == oldBlockGroupId {
		sizeDelta := newSize - oldSize
		if sizeDelta != 0 {
			blockGroupIdToSizeDelta[oldBlockGroupId] = sizeDelta
		}
	} else {
		blockGroupIdToSizeDelta[oldBlockGroupId] -= oldSize
		blockGroupIdToSizeDelta[newBlockGroupId] += newSize
		affectedBlockGroupIds = append(affectedBlockGroupIds, oldBlockGroupId)
	}

	if len(blockGroupIdToSizeDelta) > 0 {
		if exception := s.blockGroupRepository.IncrementSizesByIds(
			reqDto.ContextFields.UserId,
			blockGroupIdToSizeDelta,
			options.WithTransactionDB(tx),
			options.WithOnlyDeleted(types.Ternary_Negative),
			options.WithSkipPermissionCheck(),
		); exception != nil {
			tx.Rollback()
			return nil, exception
		}
	}

	exception = s.blockGroupRepository.CollectOrphanedBlockGroupsByIds(
		affectedBlockGroupIds,
		reqDto.ContextFields.UserId,
		allowedPermissions,
		options.WithTransactionDB(tx),
		options.WithOnlyDeleted(types.Ternary_Negative),
	)
	if exception != nil && !exceptions.CommonlyCompare(exception, exceptions.BlockGroup.NoChanges(), false) {
		tx.Rollback()
		return nil, exception
	}

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return nil, exceptions.Block.FailedToCommitTransaction().WithOrigin(err)
	}

	return &dtos.UpdateMyBlockByIdResDto{
		UpdatedAt: updatedBlock.UpdatedAt,
	}, nil
}

func (s *BlockService) UpdateMyBlocksByIds(
	ctx context.Context, reqDto *dtos.UpdateMyBlocksByIdsReqDto,
) (*dtos.UpdateMyBlocksByIdsResDto, *exceptions.Exception) {
	if err := validation.Validator.Struct(reqDto); err != nil {
		return nil, exceptions.Block.InvalidDto().WithOrigin(err)
	}

	tx := s.db.WithContext(ctx).Begin()

	allowedPemissions := []enums.AccessControlPermission{
		enums.AccessControlPermission_Owner,
		enums.AccessControlPermission_Admin,
		enums.AccessControlPermission_Write,
	}

	blockIds := make([]uuid.UUID, len(reqDto.Body.UpdatedBlocks))
	blockIdToUpdateDto := make(map[uuid.UUID]dtos.PartialUpdateDto[struct {
		ParentBlockId *uuid.UUID       `json:"parentBlockId" validate:"omitnil"`
		BlockGroupId  *uuid.UUID       `json:"blockGroupId" validate:"omitnil"`
		Type          *enums.BlockType `json:"type" validate:"omitnil,isblocktype"`
		Props         *json.RawMessage `json:"props"`
		Content       *json.RawMessage `json:"content"`
	}], len(reqDto.Body.UpdatedBlocks))
	for index, updatedBlock := range reqDto.Body.UpdatedBlocks {
		blockIds[index] = updatedBlock.BlockId
		blockIdToUpdateDto[updatedBlock.BlockId] = dtos.PartialUpdateDto[struct {
			ParentBlockId *uuid.UUID       `json:"parentBlockId" validate:"omitnil"`
			BlockGroupId  *uuid.UUID       `json:"blockGroupId" validate:"omitnil"`
			Type          *enums.BlockType `json:"type" validate:"omitnil,isblocktype"`
			Props         *json.RawMessage `json:"props"`
			Content       *json.RawMessage `json:"content"`
		}]{
			Values:  updatedBlock.Values,
			SetNull: updatedBlock.SetNull,
		}
	}
	blocks, exception := s.blockRepository.CheckPermissionsAndGetManyByIds(
		blockIds,
		reqDto.ContextFields.UserId,
		nil,
		allowedPemissions,
		options.WithTransactionDB(tx),
		options.WithOnlyDeleted(types.Ternary_Negative),
	)
	if exception != nil {
		tx.Rollback()
		return nil, exception
	}

	type ValidateBlockPropsAndContentDto struct {
		Id           uuid.UUID        `json:"id"`
		BlockGroupId uuid.UUID        `json:"blockGroupId"`
		Type         enums.BlockType  `json:"type"`
		Props        *json.RawMessage `json:"props"`
		Content      *json.RawMessage `json:"content"`
	}
	validateBlockPropsAndContentDto := make([]ValidateBlockPropsAndContentDto, len(blocks))
	affectedBlockGroupIdsMap := make(map[uuid.UUID]bool)
	var affectedBlockGroupIds []uuid.UUID
	blockIdToOriginalPropsSize := make(map[uuid.UUID]int64)
	blockIdToOriginalContentSize := make(map[uuid.UUID]int64)
	for index, block := range blocks {
		if blockIdToUpdateDto[block.Id].Values.Type == nil {
			validateBlockPropsAndContentDto[index] = ValidateBlockPropsAndContentDto{
				Id:           block.Id,
				BlockGroupId: block.BlockGroupId,
				Type:         block.Type,
				Props:        blockIdToUpdateDto[block.Id].Values.Props,
				Content:      blockIdToUpdateDto[block.Id].Values.Content,
			}
		} else {
			validateBlockPropsAndContentDto[index] = ValidateBlockPropsAndContentDto{
				Id:           block.Id,
				BlockGroupId: block.BlockGroupId,
				Type:         *blockIdToUpdateDto[block.Id].Values.Type,
				Props:        blockIdToUpdateDto[block.Id].Values.Props,
				Content:      blockIdToUpdateDto[block.Id].Values.Content,
			}
		}

		if blockIdToUpdateDto[block.Id].Values.BlockGroupId != nil &&
			(*blockIdToUpdateDto[block.Id].Values.BlockGroupId) != block.BlockGroupId {
			if !affectedBlockGroupIdsMap[block.BlockGroupId] {
				affectedBlockGroupIdsMap[block.BlockGroupId] = true
				affectedBlockGroupIds = append(affectedBlockGroupIds, block.BlockGroupId)
			}
		}
		blockIdToOriginalPropsSize[block.Id] = int64(len(block.Props))
		blockIdToOriginalContentSize[block.Id] = int64(len(block.Content))
	}
	validateBlockPropsAndContentFunc := func(validateBlockPropsAndContentDto ValidateBlockPropsAndContentDto) (inputs.BulkUpdateBlocksInput, error) {
		result := inputs.BulkUpdateBlocksInput{
			Id: validateBlockPropsAndContentDto.Id,
			PartialUpdateInput: inputs.PartialUpdateInput[inputs.UpdateBlockInput]{
				Values: inputs.UpdateBlockInput{},
			},
		}

		if validateBlockPropsAndContentDto.Props != nil {
			propsString := string(bytes.TrimSpace(*validateBlockPropsAndContentDto.Props))
			if propsString == "{}" || propsString == "" {
				emptyPropsJson := datatypes.JSON("{}")
				result.PartialUpdateInput.Values.Props = &emptyPropsJson
			} else {
				_, err := blocknote.ParseProps(validateBlockPropsAndContentDto.Type.String(), *validateBlockPropsAndContentDto.Props)
				if err != nil {
					return result, err
				}
				rawPropsJson := datatypes.JSON(*validateBlockPropsAndContentDto.Props)
				result.PartialUpdateInput.Values.Props = &rawPropsJson
			}
		}

		if validateBlockPropsAndContentDto.Content != nil {
			trimContent := bytes.TrimSpace(*validateBlockPropsAndContentDto.Content)
			trimContentString := string(trimContent)
			if trimContentString == "null" || trimContentString == "[]" || trimContentString == "" {
				emptyContentsJson := datatypes.JSON("[]")
				result.PartialUpdateInput.Values.Content = &emptyContentsJson
			} else {
				switch trimContent[0] {
				case '[':
					var list blocknote.InlineContentList
					if err := json.Unmarshal(trimContent, &list); err != nil {
						return result, err
					}
					rawContentJson := datatypes.JSON(*validateBlockPropsAndContentDto.Content)
					result.PartialUpdateInput.Values.Content = &rawContentJson
				case '{':
					var table blocknote.TableContent
					if err := json.Unmarshal(trimContent, &table); err != nil {
						return result, err
					}
					rawContentJson := datatypes.JSON(*validateBlockPropsAndContentDto.Content)
					result.PartialUpdateInput.Values.Content = &rawContentJson
				default:
					return result, errors.New("invalid content format: must be array or object")
				}
			}
		}

		return result, nil
	}

	validateBlocksPropsAndContentResult := concurrency.Execute(
		validateBlockPropsAndContentDto,
		min(10, max(len(validateBlockPropsAndContentDto)/10, len(validateBlockPropsAndContentDto)%10)),
		validateBlockPropsAndContentFunc,
	)

	var bulkUpdateBlocksInputs []inputs.BulkUpdateBlocksInput
	resDto := dtos.UpdateMyBlocksByIdsResDto{
		IsAllSuccess:   true,
		FailedIndexes:  []int{},
		SuccessIndexes: []int{},
		SuccessBlockGroupAndBlockIds: []struct {
			BlockGroupId uuid.UUID   `json:"blockGroupId"`
			BlockIds     []uuid.UUID `json:"blockIds"`
		}{},
		UpdatedAt: time.Now(),
	}
	successBlockGroupMap := make(map[uuid.UUID][]uuid.UUID)
	blockGroupIdToSizeDelta := make(map[uuid.UUID]int64)
	// note that the validateResult only consists of props and content
	// please not using the field of parent block id and block group id, and type
	for _, validateResult := range validateBlocksPropsAndContentResult {
		if validateResult.Err == nil {
			resDto.SuccessIndexes = append(resDto.SuccessIndexes, validateResult.Index)
			originalBlockGroupId := validateBlockPropsAndContentDto[validateResult.Index].BlockGroupId
			targetBlockGroupId := originalBlockGroupId
			if blockIdToUpdateDto[validateResult.Data.Id].Values.BlockGroupId != nil {
				targetBlockGroupId = *blockIdToUpdateDto[validateResult.Data.Id].Values.BlockGroupId
			}

			var sizeDelta int64 = 0
			var totalOriginalSize int64 = blockIdToOriginalPropsSize[validateResult.Data.Id] + blockIdToOriginalContentSize[validateResult.Data.Id]
			var totalTargetSize int64 = blockIdToOriginalPropsSize[validateResult.Data.Id] + blockIdToOriginalContentSize[validateResult.Data.Id]
			if validateResult.Data.PartialUpdateInput.Values.Props != nil {
				sizeDelta += (int64(len(*validateResult.Data.PartialUpdateInput.Values.Props)) - blockIdToOriginalPropsSize[validateResult.Data.Id])
				totalTargetSize += (int64(len(*validateResult.Data.PartialUpdateInput.Values.Props)) - blockIdToOriginalPropsSize[validateResult.Data.Id])
			}
			if validateResult.Data.PartialUpdateInput.Values.Content != nil {
				sizeDelta += (int64(len(*validateResult.Data.PartialUpdateInput.Values.Content)) - blockIdToOriginalContentSize[validateResult.Data.Id])
				totalTargetSize += (int64(len(*validateResult.Data.PartialUpdateInput.Values.Content)) - blockIdToOriginalContentSize[validateResult.Data.Id])
			}

			if targetBlockGroupId == originalBlockGroupId {
				if sizeDelta != 0 {
					blockGroupIdToSizeDelta[originalBlockGroupId] += sizeDelta
				}
			} else {
				blockGroupIdToSizeDelta[originalBlockGroupId] -= totalOriginalSize
				blockGroupIdToSizeDelta[targetBlockGroupId] += totalTargetSize
			}

			successBlockGroupMap[targetBlockGroupId] =
				append(successBlockGroupMap[targetBlockGroupId], validateResult.Data.Id)
			bulkUpdateBlocksInputs = append(bulkUpdateBlocksInputs, inputs.BulkUpdateBlocksInput{
				Id: validateResult.Data.Id,
				PartialUpdateInput: inputs.PartialUpdateInput[inputs.UpdateBlockInput]{
					Values: inputs.UpdateBlockInput{
						BlockGroupId:  &targetBlockGroupId,
						ParentBlockId: blockIdToUpdateDto[validateResult.Data.Id].Values.ParentBlockId,
						Type:          blockIdToUpdateDto[validateResult.Data.Id].Values.Type,
						Props:         validateResult.Data.PartialUpdateInput.Values.Props,
						Content:       validateResult.Data.PartialUpdateInput.Values.Content,
					},
					SetNull: blockIdToUpdateDto[validateResult.Data.Id].SetNull,
				},
			})
		} else {
			resDto.FailedIndexes = append(resDto.FailedIndexes, validateResult.Index)
			resDto.IsAllSuccess = false
		}
	}

	for blockGroupId, blockIds := range successBlockGroupMap {
		resDto.SuccessBlockGroupAndBlockIds = append(resDto.SuccessBlockGroupAndBlockIds, struct {
			BlockGroupId uuid.UUID   `json:"blockGroupId"`
			BlockIds     []uuid.UUID `json:"blockIds"`
		}{
			BlockGroupId: blockGroupId,
			BlockIds:     blockIds,
		})
	}

	exception = s.blockRepository.BulkUpdateManyByIds(
		reqDto.ContextFields.UserId,
		bulkUpdateBlocksInputs,
		options.WithTransactionDB(tx),
		options.WithOnlyDeleted(types.Ternary_Negative),
		options.WithSkipPermissionCheck(),
	)
	if exception != nil {
		tx.Rollback()
		return nil, exception
	}

	if len(blockGroupIdToSizeDelta) > 0 {
		if exception := s.blockGroupRepository.IncrementSizesByIds(
			reqDto.ContextFields.UserId,
			blockGroupIdToSizeDelta,
			options.WithTransactionDB(tx),
			options.WithOnlyDeleted(types.Ternary_Negative),
			options.WithSkipPermissionCheck(),
		); exception != nil {
			tx.Rollback()
			return nil, exception
		}
	}

	// the garbage collection of the orphaned block group which don't have any blocks
	exception = s.blockGroupRepository.CollectOrphanedBlockGroupsByIds(
		affectedBlockGroupIds,
		reqDto.ContextFields.UserId,
		allowedPemissions,
		options.WithTransactionDB(tx),
		options.WithOnlyDeleted(types.Ternary_Negative),
	)
	if exception != nil && !exceptions.CommonlyCompare(exception, exceptions.BlockGroup.NoChanges(), false) {
		tx.Rollback()
		return nil, exception
	}

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return nil, exceptions.Block.FailedToCommitTransaction().WithOrigin(err)
	}

	resDto.UpdatedAt = time.Now()
	return &resDto, nil
}

func (s *BlockService) RestoreMyBlockById(
	ctx context.Context, reqDto *dtos.RestoreMyBlockByIdReqDto,
) (*dtos.RestoreMyBlockByIdResDto, *exceptions.Exception) {
	if err := validation.Validator.Struct(reqDto); err != nil {
		return nil, exceptions.Block.InvalidDto().WithOrigin(err)
	}

	tx := s.db.WithContext(ctx).Begin()

	restoredBlock, exception := s.blockRepository.RestoreSoftDeletedOneById(
		reqDto.Body.BlockId,
		reqDto.ContextFields.UserId,
		options.WithTransactionDB(tx),
		options.WithOnlyDeleted(types.Ternary_Negative),
	)
	if exception != nil {
		tx.Rollback()
		return nil, exception
	}

	if _, exception := s.blockGroupRepository.RestoreSoftDeletedManyByIds(
		[]uuid.UUID{restoredBlock.BlockGroupId},
		reqDto.ContextFields.UserId,
		options.WithTransactionDB(tx),
		options.WithOnlyDeleted(types.Ternary_Positive),
	); exception != nil && !exceptions.CommonlyCompare(exception, exceptions.BlockGroup.NoChanges(), false) {
		tx.Rollback()
		return nil, exception
	}

	if _, exception := s.blockGroupRepository.IncrementSizeById(
		restoredBlock.BlockGroupId,
		reqDto.ContextFields.UserId,
		int64(len(restoredBlock.Props))+int64(len(restoredBlock.Content)),
		options.WithTransactionDB(tx),
		options.WithOnlyDeleted(types.Ternary_Negative),
		options.WithSkipPermissionCheck(),
	); exception != nil {
		tx.Rollback()
		return nil, exception
	}

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return nil, exceptions.Block.FailedToCommitTransaction().WithOrigin(err)
	}

	return &dtos.RestoreMyBlockByIdResDto{
		Id:            restoredBlock.Id,
		ParentBlockId: restoredBlock.ParentBlockId,
		BlockGroupId:  restoredBlock.BlockGroupId,
		Type:          restoredBlock.Type,
		Props:         restoredBlock.Props,
		Content:       restoredBlock.Content,
		DeletedAt:     restoredBlock.DeletedAt,
		UpdatedAt:     restoredBlock.UpdatedAt,
		CreatedAt:     restoredBlock.CreatedAt,
	}, nil
}

func (s *BlockService) RestoreMyBlocksByIds(
	ctx context.Context, reqDto *dtos.RestoreMyBlocksByIdsReqDto,
) (*dtos.RestoreMyBlocksByIdsResDto, *exceptions.Exception) {
	if err := validation.Validator.Struct(reqDto); err != nil {
		return nil, exceptions.Block.InvalidDto().WithOrigin(err)
	}

	tx := s.db.WithContext(ctx).Begin()

	restoredBlocks, exception := s.blockRepository.RestoreSoftDeletedManyByIds(
		reqDto.Body.BlockIds,
		reqDto.ContextFields.UserId,
		options.WithTransactionDB(tx),
		options.WithOnlyDeleted(types.Ternary_Negative),
	)
	if exception != nil {
		tx.Rollback()
		return nil, exception
	}

	blockGroupIdToSizeDelta := make(map[uuid.UUID]int64)
	for _, restoredBlock := range restoredBlocks {
		blockGroupIdToSizeDelta[restoredBlock.BlockGroupId] += (int64(len(restoredBlock.Props)) + int64(len(restoredBlock.Content)))
	}

	if len(blockGroupIdToSizeDelta) > 0 {
		restoreBlockGroupIds := make([]uuid.UUID, 0, len(blockGroupIdToSizeDelta))
		for blockGroupId := range blockGroupIdToSizeDelta {
			restoreBlockGroupIds = append(restoreBlockGroupIds, blockGroupId)
		}

		if _, exception := s.blockGroupRepository.RestoreSoftDeletedManyByIds(
			restoreBlockGroupIds,
			reqDto.ContextFields.UserId,
			options.WithTransactionDB(tx),
			options.WithOnlyDeleted(types.Ternary_Positive),
		); exception != nil && !exceptions.CommonlyCompare(exception, exceptions.BlockGroup.NoChanges(), false) {
			tx.Rollback()
			return nil, exception
		}

		if exception := s.blockGroupRepository.IncrementSizesByIds(
			reqDto.ContextFields.UserId,
			blockGroupIdToSizeDelta,
			options.WithTransactionDB(tx),
			options.WithOnlyDeleted(types.Ternary_Negative),
			options.WithSkipPermissionCheck(),
		); exception != nil {
			tx.Rollback()
			return nil, exception
		}
	}

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return nil, exceptions.Block.FailedToCommitTransaction().WithOrigin(err)
	}

	resDto := make(dtos.RestoreMyBlocksByIdsResDto, len(restoredBlocks))
	for index, restoredBlock := range restoredBlocks {
		resDto[index] = dtos.RestoreMyBlockByIdResDto{
			Id:            restoredBlock.Id,
			ParentBlockId: restoredBlock.ParentBlockId,
			BlockGroupId:  restoredBlock.BlockGroupId,
			Type:          restoredBlock.Type,
			Props:         restoredBlock.Props,
			Content:       restoredBlock.Content,
			DeletedAt:     restoredBlock.DeletedAt,
			UpdatedAt:     restoredBlock.UpdatedAt,
			CreatedAt:     restoredBlock.CreatedAt,
		}
	}

	return &resDto, nil
}

func (s *BlockService) DeleteMyBlockById(
	ctx context.Context, reqDto *dtos.DeleteMyBlockByIdReqDto,
) (*dtos.DeleteMyBlockByIdResDto, *exceptions.Exception) {
	if err := validation.Validator.Struct(reqDto); err != nil {
		return nil, exceptions.Block.InvalidDto().WithOrigin(err)
	}

	tx := s.db.WithContext(ctx).Begin()

	deletedBlock, exception := s.blockRepository.SoftDeleteOneById(
		reqDto.Body.BlockId,
		reqDto.ContextFields.UserId,
		options.WithTransactionDB(tx),
		options.WithOnlyDeleted(types.Ternary_Negative),
	)
	if exception != nil {
		tx.Rollback()
		return nil, exception
	}

	// the garbage collection of the block group of deleted block
	var remainingBlockCount int64 = 0
	result := tx.Model(&schemas.Block{}).
		Where("block_group_id IN ?", deletedBlock.BlockGroupId).
		Count(&remainingBlockCount)
	if err := result.Error; err != nil {
		tx.Rollback()
		return nil, exceptions.Block.NotFound().WithOrigin(err)
	}
	if remainingBlockCount == 0 {
		if exception := s.blockGroupRepository.SoftDeleteOneById(
			deletedBlock.BlockGroupId,
			reqDto.ContextFields.UserId,
			options.WithTransactionDB(tx),
			options.WithOnlyDeleted(types.Ternary_Negative),
			options.WithSkipPermissionCheck(),
		); exception != nil {
			tx.Rollback()
			return nil, exception
		}
	} else { // update the size of the involved block group
		if _, exception := s.blockGroupRepository.IncrementSizeById(
			deletedBlock.BlockGroupId,
			reqDto.ContextFields.UserId,
			-int64(len(deletedBlock.Props))-int64(len(deletedBlock.Content)),
			options.WithTransactionDB(tx),
			options.WithOnlyDeleted(types.Ternary_Negative),
			options.WithSkipPermissionCheck(),
		); exception != nil {
			tx.Rollback()
			return nil, exception
		}
	}

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return nil, exceptions.Block.FailedToCommitTransaction().WithOrigin(err)
	}

	return &dtos.DeleteMyBlockByIdResDto{
		DeletedAt: time.Now(),
	}, nil
}

func (s *BlockService) DeleteMyBlocksByIds(
	ctx context.Context, reqDto *dtos.DeleteMyBlocksByIdsReqDto,
) (*dtos.DeleteMyBlockPacksByIdsResDto, *exceptions.Exception) {
	if err := validation.Validator.Struct(reqDto); err != nil {
		return nil, exceptions.Block.InvalidDto().WithOrigin(err)
	}

	tx := s.db.WithContext(ctx).Begin()

	deletedBlocks, exception := s.blockRepository.SoftDeleteManyByIds(
		reqDto.Body.BlockIds,
		reqDto.ContextFields.UserId,
		options.WithTransactionDB(tx),
		options.WithOnlyDeleted(types.Ternary_Negative),
	)
	if exception != nil {
		tx.Rollback()
		return nil, exception
	}

	// the garbage collection of the orphaned block group which don't have any blocks
	affectedBlockGroupIdsMap := make(map[uuid.UUID]bool)
	var affectedBlockGroupIds []uuid.UUID
	blockGroupIdToSizeDelta := make(map[uuid.UUID]int64)
	for _, deletedBlock := range deletedBlocks {
		if !affectedBlockGroupIdsMap[deletedBlock.BlockGroupId] {
			affectedBlockGroupIdsMap[deletedBlock.BlockGroupId] = true
			affectedBlockGroupIds = append(affectedBlockGroupIds, deletedBlock.BlockGroupId)
		}
		blockGroupIdToSizeDelta[deletedBlock.BlockGroupId] -= (int64(len(deletedBlock.Props)) + int64(len(deletedBlock.Content)))
	}

	allowedPermissions := []enums.AccessControlPermission{
		enums.AccessControlPermission_Owner,
		enums.AccessControlPermission_Admin,
		enums.AccessControlPermission_Write,
	}
	if exception := s.blockGroupRepository.IncrementSizesByIds(
		reqDto.ContextFields.UserId,
		blockGroupIdToSizeDelta,
		options.WithTransactionDB(tx),
		options.WithOnlyDeleted(types.Ternary_Negative),
		options.WithSkipPermissionCheck(),
	); exception != nil && !exceptions.CommonlyCompare(exception, exceptions.BlockGroup.NoChanges(), false) {
		tx.Rollback()
		return nil, exception
	}

	exception = s.blockGroupRepository.CollectOrphanedBlockGroupsByIds(
		affectedBlockGroupIds,
		reqDto.ContextFields.UserId,
		allowedPermissions,
		options.WithTransactionDB(tx),
		options.WithOnlyDeleted(types.Ternary_Negative),
	)
	if exception != nil && !exceptions.CommonlyCompare(exception, exceptions.BlockGroup.NoChanges(), false) {
		tx.Rollback()
		return nil, exception
	}

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return nil, exceptions.Block.FailedToCommitTransaction().WithOrigin(err)
	}

	return &dtos.DeleteMyBlockPacksByIdsResDto{
		DeletedAt: time.Now(),
	}, nil
}
