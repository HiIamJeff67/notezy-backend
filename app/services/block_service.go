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

	adapters "notezy-backend/app/adapters"
	dtos "notezy-backend/app/dtos"
	exceptions "notezy-backend/app/exceptions"
	concurrency "notezy-backend/app/lib/concurrency"
	inputs "notezy-backend/app/models/inputs"
	repositories "notezy-backend/app/models/repositories"
	schemas "notezy-backend/app/models/schemas"
	enums "notezy-backend/app/models/schemas/enums"
	blockgroupsql "notezy-backend/app/models/sqls/block_group"
	options "notezy-backend/app/options"
	validation "notezy-backend/app/validation"
	constants "notezy-backend/shared/constants"
	blocknote "notezy-backend/shared/lib/blocknote"
	types "notezy-backend/shared/types"
)

/* ============================== Interface & Instance ============================== */

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

/* ============================== Implementations ============================== */

func (s *BlockService) GetMyBlockById(
	ctx context.Context, reqDto *dtos.GetMyBlockByIdReqDto,
) (*dtos.GetMyBlockByIdResDto, *exceptions.Exception) {
	if err := validation.Validator.Struct(reqDto); err != nil {
		return nil, exceptions.Block.InvalidDto().WithError(err)
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
		return nil, exceptions.Block.InvalidDto().WithError(err)
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
		return nil, exceptions.Block.InvalidDto().WithError(err)
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

	rawArborizedBlock, exception := s.editableBlockAdapter.ArborizeRawToRaw(root, childrenMap)
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
		return nil, exceptions.Block.InvalidDto().WithError(err)
	}

	db := s.db.WithContext(ctx)

	allowedPermissions := []enums.AccessControlPermission{
		enums.AccessControlPermission_Owner,
		enums.AccessControlPermission_Admin,
		enums.AccessControlPermission_Write,
		enums.AccessControlPermission_Read,
	}

	if s.blockGroupRepository.HasPermissions(
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

func (s *BlockService) GetMyBlocksByBlockPackId(
	ctx context.Context, reqDto *dtos.GetMyBlocksByBlockPackIdReqDto,
) (*dtos.GetMyBlocksByBlockPackIdResDto, *exceptions.Exception) {
	if err := validation.Validator.Struct(reqDto); err != nil {
		return nil, exceptions.Block.InvalidDto().WithError(err)
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
		return nil, exceptions.Block.NotFound().WithError(err)
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
		return nil, exceptions.Block.InvalidDto().WithError(err)
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
		return nil, exceptions.Block.InvalidDto().WithError(err)
	}

	db := s.db.WithContext(ctx)

	rawFlattenedBlocks, exception := s.editableBlockAdapter.FlattenToRaw(&reqDto.Body.ArborizedEditableBlock)
	if exception != nil {
		return nil, exception
	}

	input := make([]inputs.CreateBlockInput, len(rawFlattenedBlocks))
	for index, rawFlarawFlattenedBlock := range rawFlattenedBlocks {
		input[index] = inputs.CreateBlockInput{
			Id:            rawFlarawFlattenedBlock.Id,
			ParentBlockId: rawFlarawFlattenedBlock.ParentBlockId,
			Type:          rawFlarawFlattenedBlock.Type,
			Props:         rawFlarawFlattenedBlock.Props,
			Content:       rawFlarawFlattenedBlock.Content,
		}
	}

	_, exception = s.blockRepository.CreateManyByBlockGroupId(
		reqDto.Body.BlockGroupId,
		reqDto.ContextFields.UserId,
		input,
		options.WithDB(db),
		options.WithBatchSize(constants.MaxBatchCreateBlockSize),
		options.WithOnlyDeleted(types.Ternary_Negative),
	)
	if exception != nil {
		return nil, exception
	}

	return &dtos.InsertBlockResDto{
		CreatedAt: time.Now(),
	}, nil
}

func (s *BlockService) InsertBlocks(
	ctx context.Context, reqDto *dtos.InsertBlocksReqDto,
) (*dtos.InsertBlocksResDto, *exceptions.Exception) {
	if err := validation.Validator.Struct(reqDto); err != nil {
		return nil, exceptions.Block.InvalidDto().WithError(err)
	}

	db := s.db.WithContext(ctx)

	type ValidateBlockDto struct {
		ParentBlockId          *uuid.UUID
		ArborizedEditableBlock dtos.ArborizedEditableBlock
	}
	validateBlockDto := make([]ValidateBlockDto, len(reqDto.Body.InsertedBlocks))
	for index, insertedBlock := range reqDto.Body.InsertedBlocks {
		validateBlockDto[index] = ValidateBlockDto{
			ParentBlockId:          insertedBlock.ParentBlockId,
			ArborizedEditableBlock: insertedBlock.ArborizedEditableBlock,
		}
	}

	validateBlockFunc := func(validateBlockDto ValidateBlockDto) ([]dtos.RawFlattenedEditableBlock, error) {
		rawFlattenedBlocks, exception := s.editableBlockAdapter.FlattenToRaw(&validateBlockDto.ArborizedEditableBlock)
		if exception != nil {
			return rawFlattenedBlocks, exception.ToError()
		}

		if len(rawFlattenedBlocks) > 0 {
			rawFlattenedBlocks[0].ParentBlockId = validateBlockDto.ParentBlockId
		}
		return rawFlattenedBlocks, nil
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
	for _, validateResult := range validateBlockResults {
		if validateResult.Err == nil {
			resDto.SuccessIndexes = append(resDto.SuccessIndexes, validateResult.Index)
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
				BlockGroupId: reqDto.Body.InsertedBlocks[validateResult.Index].BlockGroupId,
				BlockIds:     blockIds,
			})
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

	_, exception := s.blockRepository.CreateManyByBlockGroupIds(
		reqDto.ContextFields.UserId,
		createBlockGroupContentInput,
		options.WithDB(db),
		options.WithBatchSize(constants.MaxBatchCreateBlockSize),
		options.WithOnlyDeleted(types.Ternary_Negative),
		options.WithSkipPermissionCheck(),
	)
	if exception != nil {
		return nil, exception
	}

	resDto.CreatedAt = time.Now()
	return &resDto, nil
}

func (s *BlockService) UpdateMyBlockById(
	ctx context.Context, reqDto *dtos.UpdateMyBlockByIdReqDto,
) (*dtos.UpdateMyBlockByIdResDto, *exceptions.Exception) {
	if err := validation.Validator.Struct(reqDto); err != nil {
		return nil, exceptions.Block.InvalidDto().WithError(err)
	}

	db := s.db.WithContext(ctx)

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
		options.WithDB(db),
		options.WithOnlyDeleted(types.Ternary_Negative),
	)
	if exception != nil {
		return nil, exception
	}

	updateInput := inputs.PartialUpdateBlockInput{
		Values: inputs.UpdateBlockInput{
			ParentBlockId: reqDto.Body.Values.ParentBlockId,
			BlockGroupId:  reqDto.Body.Values.BlockGroupId,
			Props:         nil,
			Content:       nil,
		},
		SetNull: reqDto.Body.SetNull,
	}

	if reqDto.Body.Values.Props != nil {
		propsString := string(bytes.TrimSpace(*reqDto.Body.Values.Props))
		if propsString == "{}" || propsString == "" {
			emptyPropsJson := datatypes.JSON("{}")
			updateInput.Values.Props = &emptyPropsJson
		} else {
			_, err := blocknote.ParseProps(block.Type.String(), *reqDto.Body.Values.Props)
			if err != nil {
				return nil, exceptions.Block.InvalidDto().WithError(err)
			}
			rawPropsJson := datatypes.JSON(*reqDto.Body.Values.Props)
			updateInput.Values.Props = &rawPropsJson
		}
	}

	if reqDto.Body.Values.Content != nil {
		trimContent := bytes.TrimSpace(*reqDto.Body.Values.Content)
		trimContentString := string(trimContent)
		if trimContentString == "null" || trimContentString == "[]" || trimContentString == "" {
			emptyContentsJson := datatypes.JSON("[]")
			updateInput.Values.Content = &emptyContentsJson
		} else {
			switch trimContent[0] {
			case '[':
				var list blocknote.InlineContentList
				if err := json.Unmarshal(trimContent, &list); err != nil {
					return nil, exceptions.Block.InvalidDto().WithError(err)
				}
				rawContentJson := datatypes.JSON(*reqDto.Body.Values.Content)
				updateInput.Values.Content = &rawContentJson
			case '{':
				var table blocknote.TableContent
				if err := json.Unmarshal(trimContent, &table); err != nil {
					return nil, exceptions.Block.InvalidDto().WithError(err)
				}
				rawContentJson := datatypes.JSON(*reqDto.Body.Values.Content)
				updateInput.Values.Content = &rawContentJson
			default:
				return nil, exceptions.Block.InvalidDto().WithError(errors.New("invalid content format: must be array or object"))
			}
		}
	}

	updatedBlock, exception := s.blockRepository.UpdateOneById(
		reqDto.Body.BlockId,
		reqDto.ContextFields.UserId,
		updateInput,
		options.WithDB(db),
		options.WithOnlyDeleted(types.Ternary_Negative),
		options.WithSkipPermissionCheck(),
	)
	if exception != nil {
		return nil, exception
	}

	return &dtos.UpdateMyBlockByIdResDto{
		UpdatedAt: updatedBlock.UpdatedAt,
	}, nil
}

func (s *BlockService) UpdateMyBlocksByIds(
	ctx context.Context, reqDto *dtos.UpdateMyBlocksByIdsReqDto,
) (*dtos.UpdateMyBlocksByIdsResDto, *exceptions.Exception) {
	if err := validation.Validator.Struct(reqDto); err != nil {
		return nil, exceptions.Block.InvalidDto().WithError(err)
	}

	db := s.db.WithContext(ctx)

	allowedPemissions := []enums.AccessControlPermission{
		enums.AccessControlPermission_Owner,
		enums.AccessControlPermission_Admin,
		enums.AccessControlPermission_Write,
	}

	blockIds := make([]uuid.UUID, len(reqDto.Body.UpdatedBlocks))
	blockIdToUpdateDto := make(map[uuid.UUID]dtos.PartialUpdateDto[struct {
		ParentBlockId *uuid.UUID       `json:"parentBlockId" validate:"omitnil"`
		BlockGroupId  *uuid.UUID       `json:"blockGroupId" validate:"omitnil"`
		Props         *json.RawMessage `json:"props"`
		Content       *json.RawMessage `json:"content"`
	}], len(reqDto.Body.UpdatedBlocks))
	for index, updatedBlock := range reqDto.Body.UpdatedBlocks {
		blockIds[index] = updatedBlock.BlockId
		blockIdToUpdateDto[updatedBlock.BlockId] = dtos.PartialUpdateDto[struct {
			ParentBlockId *uuid.UUID       `json:"parentBlockId" validate:"omitnil"`
			BlockGroupId  *uuid.UUID       `json:"blockGroupId" validate:"omitnil"`
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
		options.WithDB(db),
		options.WithOnlyDeleted(types.Ternary_Negative),
	)
	if exception != nil {
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
	for index, block := range blocks {
		validateBlockPropsAndContentDto[index] = ValidateBlockPropsAndContentDto{
			Id:           block.Id,
			BlockGroupId: block.BlockGroupId,
			Type:         block.Type,
			Props:        blockIdToUpdateDto[block.Id].Values.Props,
			Content:      blockIdToUpdateDto[block.Id].Values.Content,
		}
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

	var bulkUpdateBlocksInputs inputs.BulkUpdateBlocksInputs
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
	for _, validateResult := range validateBlocksPropsAndContentResult {
		if validateResult.Err == nil {
			resDto.SuccessIndexes = append(resDto.SuccessIndexes, validateResult.Index)
			successBlockGroupMap[validateBlockPropsAndContentDto[validateResult.Index].BlockGroupId] = append(successBlockGroupMap[validateBlockPropsAndContentDto[validateResult.Index].BlockGroupId], validateResult.Data.Id)
			bulkUpdateBlocksInputs = append(bulkUpdateBlocksInputs, inputs.BulkUpdateBlocksInput{
				Id: validateResult.Data.Id,
				PartialUpdateInput: inputs.PartialUpdateInput[inputs.UpdateBlockInput]{
					Values: inputs.UpdateBlockInput{
						ParentBlockId: blockIdToUpdateDto[validateResult.Data.Id].Values.ParentBlockId,
						BlockGroupId:  blockIdToUpdateDto[validateResult.Data.Id].Values.BlockGroupId,
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
		options.WithDB(db),
		options.WithOnlyDeleted(types.Ternary_Negative),
		options.WithSkipPermissionCheck(),
	)
	if exception != nil {
		return nil, exception
	}

	resDto.UpdatedAt = time.Now()
	return &resDto, nil
}

func (s *BlockService) RestoreMyBlockById(
	ctx context.Context, reqDto *dtos.RestoreMyBlockByIdReqDto,
) (*dtos.RestoreMyBlockByIdResDto, *exceptions.Exception) {
	if err := validation.Validator.Struct(reqDto); err != nil {
		return nil, exceptions.Block.InvalidDto().WithError(err)
	}

	db := s.db.WithContext(ctx)

	restoredBlock, exception := s.blockRepository.RestoreSoftDeletedOneById(
		reqDto.Body.BlockId,
		reqDto.ContextFields.UserId,
		options.WithDB(db),
		options.WithOnlyDeleted(types.Ternary_Negative),
	)
	if exception != nil {
		return nil, exception
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
		return nil, exceptions.Block.InvalidDto().WithError(err)
	}

	db := s.db.WithContext(ctx)

	restoredBlocks, exception := s.blockRepository.RestoreSoftDeletedManyByIds(
		reqDto.Body.BlockIds,
		reqDto.ContextFields.UserId,
		options.WithDB(db),
		options.WithOnlyDeleted(types.Ternary_Negative),
	)
	if exception != nil {
		return nil, exception
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
		return nil, exceptions.Block.InvalidDto().WithError(err)
	}

	tx := s.db.WithContext(ctx).Begin()

	deletedBlock, exception := s.blockRepository.SoftDeleteOneById(
		reqDto.Body.BlockId,
		reqDto.ContextFields.UserId,
		options.WithDB(tx),
		options.WithOnlyDeleted(types.Ternary_Negative),
	)
	if exception != nil {
		return nil, exception
	}

	// the garbage collection of the block group of deleted block
	var remainingBlockCount int64 = 0
	result := tx.Model(&schemas.Block{}).
		Where("block_group_id IN ?", deletedBlock.BlockGroupId).
		Count(&remainingBlockCount)
	if err := result.Error; err != nil {
		tx.Rollback()
		return nil, exceptions.Block.NotFound().WithError(err)
	}
	if remainingBlockCount == 0 {
		if exception := s.blockGroupRepository.SoftDeleteOneById(
			deletedBlock.BlockGroupId,
			reqDto.ContextFields.UserId,
			options.WithDB(tx),
			options.WithOnlyDeleted(types.Ternary_Negative),
			options.WithSkipPermissionCheck(),
		); exception != nil {
			tx.Rollback()
			return nil, exception
		}
	}

	if err := tx.Commit().Error; err != nil {
		return nil, exceptions.Block.FailedToCommitTransaction().WithError(err)
	}

	return &dtos.DeleteMyBlockByIdResDto{
		DeletedAt: time.Now(),
	}, nil
}

func (s *BlockService) DeleteMyBlocksByIds(
	ctx context.Context, reqDto *dtos.DeleteMyBlocksByIdsReqDto,
) (*dtos.DeleteMyBlockPacksByIdsResDto, *exceptions.Exception) {
	if err := validation.Validator.Struct(reqDto); err != nil {
		return nil, exceptions.Block.InvalidDto().WithError(err)
	}

	tx := s.db.WithContext(ctx).Begin()

	deletedBlocks, exception := s.blockRepository.SoftDeleteManyByIds(
		reqDto.Body.BlockIds,
		reqDto.ContextFields.UserId,
		options.WithDB(tx),
		options.WithOnlyDeleted(types.Ternary_Negative),
	)
	if exception != nil {
		return nil, exception
	}

	// the garbage collection of the block group of deleted block
	affectedGroupIdsMap := make(map[uuid.UUID]bool)
	var affectedGroupIds []uuid.UUID
	for _, block := range deletedBlocks {
		if !affectedGroupIdsMap[block.BlockGroupId] {
			affectedGroupIdsMap[block.BlockGroupId] = true
			affectedGroupIds = append(affectedGroupIds, block.BlockGroupId)
		}
	}
	if len(affectedGroupIds) > 0 {
		result := tx.Exec(blockgroupsql.CollectGarbageBlockGroupByIdsSQL, affectedGroupIds)
		if err := result.Error; err != nil {
			tx.Rollback()
			return nil, exceptions.BlockGroup.FailedToDelete().WithError(err)
		}
	}

	if err := tx.Commit().Error; err != nil {
		return nil, exceptions.Block.FailedToCommitTransaction().WithError(err)
	}

	return &dtos.DeleteMyBlockPacksByIdsResDto{
		DeletedAt: time.Now(),
	}, nil
}
