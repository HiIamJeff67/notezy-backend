package services

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/jinzhu/copier"
	"gorm.io/datatypes"
	"gorm.io/gorm"

	adapters "github.com/HiIamJeff67/notezy-backend/app/adapters"
	dtos "github.com/HiIamJeff67/notezy-backend/app/dtos"
	exceptions "github.com/HiIamJeff67/notezy-backend/app/exceptions"
	gqlmodels "github.com/HiIamJeff67/notezy-backend/app/graphql/models"
	inputs "github.com/HiIamJeff67/notezy-backend/app/models/inputs"
	repositories "github.com/HiIamJeff67/notezy-backend/app/models/repositories"
	schemas "github.com/HiIamJeff67/notezy-backend/app/models/schemas"
	enums "github.com/HiIamJeff67/notezy-backend/app/models/schemas/enums"
	scopes "github.com/HiIamJeff67/notezy-backend/app/models/scopes"
	options "github.com/HiIamJeff67/notezy-backend/app/options"
	util "github.com/HiIamJeff67/notezy-backend/app/util"
	validation "github.com/HiIamJeff67/notezy-backend/app/validation"
	constants "github.com/HiIamJeff67/notezy-backend/shared/constants"
	searchcursor "github.com/HiIamJeff67/notezy-backend/shared/lib/searchcursor"
	types "github.com/HiIamJeff67/notezy-backend/shared/types"
)

type BlockServiceInterface interface {
	GetMyBlockById(ctx context.Context, reqDto *dtos.GetMyBlockByIdReqDto) (*dtos.GetMyBlockByIdResDto, *exceptions.Exception)
	GetMyBlocksByIds(ctx context.Context, reqDto *dtos.GetMyBlocksByIdsReqDto) (*dtos.GetMyBlocksByIdsResDto, *exceptions.Exception)
	GetMyBlocksByBlockPackId(ctx context.Context, reqDto *dtos.GetMyBlocksByBlockPackIdReqDto) (*dtos.GetMyBlocksByBlockPackIdResDto, *exceptions.Exception)
	GetAllMyBlocks(ctx context.Context, reqDto *dtos.GetAllMyBlocksReqDto) (*dtos.GetAllMyBlocksResDto, *exceptions.Exception)
	AppendBlock(ctx context.Context, reqDto *dtos.AppendBlockReqDto) (*dtos.AppendBlockResDto, *exceptions.Exception)
	AppendBlocks(ctx context.Context, reqDto *dtos.AppendBlocksReqDto) (*dtos.AppendBlocksResDto, *exceptions.Exception)
	InsertBlock(ctx context.Context, reqDto *dtos.InsertBlockReqDto) (*dtos.InsertBlockResDto, *exceptions.Exception)
	InsertBlocks(ctx context.Context, reqDto *dtos.InsertBlocksReqDto) (*dtos.InsertBlocksResDto, *exceptions.Exception)
	UpdateMyBlockById(ctx context.Context, reqDto *dtos.UpdateMyBlockByIdReqDto) (*dtos.UpdateMyBlockByIdResDto, *exceptions.Exception)
	UpdateMyBlocksByIds(ctx context.Context, reqDto *dtos.UpdateMyBlocksByIdsReqDto) (*dtos.UpdateMyBlocksByIdsResDto, *exceptions.Exception)
	RestoreMyBlockById(ctx context.Context, reqDto *dtos.RestoreMyBlockByIdReqDto) (*dtos.RestoreMyBlockByIdResDto, *exceptions.Exception)
	RestoreMyBlocksByIds(ctx context.Context, reqDto *dtos.RestoreMyBlocksByIdsReqDto) (*dtos.RestoreMyBlocksByIdsResDto, *exceptions.Exception)
	DeleteMyBlockById(ctx context.Context, reqDto *dtos.DeleteMyBlockByIdReqDto) (*dtos.DeleteMyBlockByIdResDto, *exceptions.Exception)
	DeleteMyBlocksByIds(ctx context.Context, reqDto *dtos.DeleteMyBlocksByIdsReqDto) (*dtos.DeleteMyBlocksByIdsResDto, *exceptions.Exception)

	SearchPrivateBlocks(ctx context.Context, userId uuid.UUID, gqlInput gqlmodels.SearchBlockInput) (*gqlmodels.SearchBlockConnection, *exceptions.Exception)
}

type BlockService struct {
	db                   *gorm.DB
	blockScope           scopes.BlockScopeInterface
	blockPackScope       scopes.BlockPackScopeInterface
	subShelfScope        scopes.SubShelfScopeInterface
	blockPackRepository  repositories.BlockPackRepositoryInterface
	blockRepository      repositories.BlockRepositoryInterface
	editableBlockAdapter adapters.EditableBlockAdapterInterface
}

func NewBlockService(
	db *gorm.DB,
	blockScope scopes.BlockScopeInterface,
	blockPackScope scopes.BlockPackScopeInterface,
	subShelfScope scopes.SubShelfScopeInterface,
	blockPackRepository repositories.BlockPackRepositoryInterface,
	blockRepository repositories.BlockRepositoryInterface,
	editableBlockAdapter adapters.EditableBlockAdapterInterface,
) BlockServiceInterface {
	return &BlockService{
		db:                   db,
		blockScope:           blockScope,
		blockPackScope:       blockPackScope,
		subShelfScope:        subShelfScope,
		blockPackRepository:  blockPackRepository,
		blockRepository:      blockRepository,
		editableBlockAdapter: editableBlockAdapter,
	}
}

func (s *BlockService) GetMyBlockById(
	ctx context.Context, reqDto *dtos.GetMyBlockByIdReqDto,
) (*dtos.GetMyBlockByIdResDto, *exceptions.Exception) {
	if err := validation.Validator.Struct(reqDto); err != nil {
		return nil, exceptions.Block.InvalidDto().WithOrigin(err)
	}

	block, exception := s.blockRepository.GetOneById(reqDto.Param.BlockId, reqDto.ContextFields.UserId, nil, options.WithDB(s.db.WithContext(ctx)), options.WithOnlyDeleted(types.Ternary_Negative))
	if exception != nil {
		return nil, exception
	}

	res := dtos.GetMyBlockByIdResDto{
		Id:            block.Id,
		BlockPackId:   block.BlockPackId,
		ParentBlockId: block.ParentBlockId,
		PrevBlockId:   block.PrevBlockId,
		NextBlockId:   block.NextBlockId,
		Type:          block.Type,
		Props:         block.Props,
		Content:       block.Content,
		DeletedAt:     block.DeletedAt,
		UpdatedAt:     block.UpdatedAt,
		CreatedAt:     block.CreatedAt,
	}

	return &res, nil
}

func (s *BlockService) GetMyBlocksByIds(
	ctx context.Context, reqDto *dtos.GetMyBlocksByIdsReqDto,
) (*dtos.GetMyBlocksByIdsResDto, *exceptions.Exception) {
	if err := validation.Validator.Struct(reqDto); err != nil {
		return nil, exceptions.Block.InvalidDto().WithOrigin(err)
	}

	blocks, exception := s.blockRepository.CheckPermissionsAndGetManyByIds(
		reqDto.Param.BlockIds,
		reqDto.ContextFields.UserId,
		nil,
		[]enums.AccessControlPermission{
			enums.AccessControlPermission_Owner,
			enums.AccessControlPermission_Admin,
			enums.AccessControlPermission_Write,
			enums.AccessControlPermission_Read,
		},
		options.WithDB(s.db.WithContext(ctx)),
		options.WithOnlyDeleted(types.Ternary_Negative),
	)
	if exception != nil {
		return nil, exception
	}

	res := make(dtos.GetMyBlocksByIdsResDto, len(blocks))
	for index, block := range blocks {
		res[index] = dtos.GetMyBlockByIdResDto{
			Id:            block.Id,
			BlockPackId:   block.BlockPackId,
			ParentBlockId: block.ParentBlockId,
			PrevBlockId:   block.PrevBlockId,
			NextBlockId:   block.NextBlockId,
			Type:          block.Type,
			Props:         block.Props,
			Content:       block.Content,
			DeletedAt:     block.DeletedAt,
			UpdatedAt:     block.UpdatedAt,
			CreatedAt:     block.CreatedAt,
		}
	}

	return &res, nil
}

func (s *BlockService) GetMyBlocksByBlockPackId(
	ctx context.Context, reqDto *dtos.GetMyBlocksByBlockPackIdReqDto,
) (*dtos.GetMyBlocksByBlockPackIdResDto, *exceptions.Exception) {
	if err := validation.Validator.Struct(reqDto); err != nil {
		return nil, exceptions.Block.InvalidDto().WithOrigin(err)
	}

	db := s.db.WithContext(ctx)

	if !s.blockPackRepository.HasPermission(
		reqDto.Param.BlockPackId,
		reqDto.ContextFields.UserId,
		[]enums.AccessControlPermission{
			enums.AccessControlPermission_Owner,
			enums.AccessControlPermission_Admin,
			enums.AccessControlPermission_Write,
			enums.AccessControlPermission_Read,
		},
		options.WithDB(db),
		options.WithOnlyDeleted(types.Ternary_Negative),
	) {
		return nil, exceptions.Block.NoPermission("get the block pack of blocks")
	}

	var blocks []schemas.Block
	if err := db.Model(&schemas.Block{}).
		Where("block_pack_id = ?", reqDto.Param.BlockPackId).
		Scopes(s.blockScope.FilterOnlyDeleted(types.Ternary_Negative)).
		Order("created_at ASC").
		Find(&blocks).Error; err != nil {
		return nil, exceptions.Block.NotFound().WithOrigin(err)
	}

	res := make(dtos.GetMyBlocksByBlockPackIdResDto, len(blocks))
	for index, block := range blocks {
		res[index] = dtos.GetMyBlockByIdResDto{
			Id:            block.Id,
			BlockPackId:   block.BlockPackId,
			ParentBlockId: block.ParentBlockId,
			PrevBlockId:   block.PrevBlockId,
			NextBlockId:   block.NextBlockId,
			Type:          block.Type,
			Props:         block.Props,
			Content:       block.Content,
			DeletedAt:     block.DeletedAt,
			UpdatedAt:     block.UpdatedAt,
			CreatedAt:     block.CreatedAt,
		}
	}

	return &res, nil
}

func (s *BlockService) GetAllMyBlocks(
	ctx context.Context, reqDto *dtos.GetAllMyBlocksReqDto,
) (*dtos.GetAllMyBlocksResDto, *exceptions.Exception) {
	if err := validation.Validator.Struct(reqDto); err != nil {
		return nil, exceptions.Block.InvalidDto().WithOrigin(err)
	}

	var blocks []schemas.Block
	if err := s.db.WithContext(ctx).Model(&schemas.Block{}).
		Select(`"BlockTable".*`).
		Joins(`INNER JOIN "BlockPackTable" bp ON bp.id = "BlockTable".block_pack_id`).
		Joins(`INNER JOIN "SubShelfTable" ss ON ss.id = bp.parent_sub_shelf_id`).
		Joins(`INNER JOIN "UsersToShelvesTable" uts ON uts.root_shelf_id = ss.root_shelf_id`).
		Where("uts.user_id = ? AND uts.permission IN ?", reqDto.ContextFields.UserId, []enums.AccessControlPermission{
			enums.AccessControlPermission_Owner,
			enums.AccessControlPermission_Admin,
			enums.AccessControlPermission_Write,
			enums.AccessControlPermission_Read,
		}).
		Scopes(s.blockScope.FilterOnlyDeleted(types.Ternary_Negative)).
		Find(&blocks).Error; err != nil {
		return nil, exceptions.Block.NotFound().WithOrigin(err)
	}

	res := make(dtos.GetAllMyBlocksResDto, len(blocks))
	for index, block := range blocks {
		res[index] = dtos.GetMyBlockByIdResDto{
			Id:            block.Id,
			BlockPackId:   block.BlockPackId,
			ParentBlockId: block.ParentBlockId,
			PrevBlockId:   block.PrevBlockId,
			NextBlockId:   block.NextBlockId,
			Type:          block.Type,
			Props:         block.Props,
			Content:       block.Content,
			DeletedAt:     block.DeletedAt,
			UpdatedAt:     block.UpdatedAt,
			CreatedAt:     block.CreatedAt,
		}
	}

	return &res, nil
}

func (s *BlockService) AppendBlock(
	ctx context.Context, reqDto *dtos.AppendBlockReqDto,
) (*dtos.AppendBlockResDto, *exceptions.Exception) {
	if err := validation.Validator.Struct(reqDto); err != nil {
		return nil, exceptions.Block.InvalidDto().WithOrigin(err)
	}

	rawBlocks, _, exception := s.editableBlockAdapter.FlattenToRaw(&reqDto.Body.ArborizedEditableBlock)
	if exception != nil {
		return nil, exception
	}

	tx := s.db.WithContext(ctx).Begin()

	if !s.blockPackRepository.HasPermission(
		reqDto.Body.BlockPackId,
		reqDto.ContextFields.UserId,
		[]enums.AccessControlPermission{
			enums.AccessControlPermission_Owner,
			enums.AccessControlPermission_Admin,
			enums.AccessControlPermission_Write,
		},
		options.WithTransactionDB(tx),
		options.WithOnlyDeleted(types.Ternary_Negative),
	) {
		tx.Rollback()
		return nil, exceptions.Block.NoPermission("append blocks to the block pack")
	}

	lockingStrength := options.LockingStrengthUpdate
	if err := tx.Scopes(scopes.Locking(&lockingStrength)).
		Where(`"BlockPackTable".id = ? AND "BlockPackTable".deleted_at IS NULL`, reqDto.Body.BlockPackId).
		First(&schemas.BlockPack{}).Error; err != nil {
		tx.Rollback()
		return nil, exceptions.BlockPack.NotFound().WithOrigin(err)
	}

	var tail schemas.Block
	var prevBlockId *uuid.UUID
	if err := tx.Model(&schemas.Block{}).
		Where("block_pack_id = ? AND deleted_at IS NULL AND parent_block_id IS NULL AND next_block_id IS NULL", reqDto.Body.BlockPackId).
		First(&tail).Error; err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		tx.Rollback()
		return nil, exceptions.Block.NotFound().WithOrigin(err)
	}
	if tail.Id != uuid.Nil {
		prevBlockId = &tail.Id
	}

	rootId := rawBlocks[0].Id
	rawBlocks[0].ParentBlockId = nil
	rawBlocks[0].PrevBlockId = prevBlockId
	rawBlocks[0].NextBlockId = nil

	deletedAt := time.Now()
	newBlocks := make([]schemas.Block, len(rawBlocks))
	blockIds := make([]uuid.UUID, len(rawBlocks))
	for index, rawBlock := range rawBlocks {
		blockIds[index] = rawBlock.Id
		newBlocks[index] = schemas.Block{
			Id:            rawBlock.Id,
			BlockPackId:   reqDto.Body.BlockPackId,
			ParentBlockId: rawBlock.ParentBlockId,
			PrevBlockId:   rawBlock.PrevBlockId,
			NextBlockId:   rawBlock.NextBlockId,
			Type:          rawBlock.Type,
			Props:         rawBlock.Props,
			Content:       rawBlock.Content,
			DeletedAt:     &deletedAt,
		}
	}

	if err := tx.CreateInBatches(&newBlocks, constants.MaxBatchCreateBlockSize).Error; err != nil {
		tx.Rollback()
		return nil, exceptions.Block.FailedToCreate().WithOrigin(err)
	}

	if prevBlockId != nil {
		if err := tx.Model(&schemas.Block{}).
			Where("id = ?", *prevBlockId).
			Update("next_block_id", rootId).Error; err != nil {
			tx.Rollback()
			return nil, exceptions.Block.FailedToUpdate().WithOrigin(err)
		}
	}

	if err := tx.Model(&schemas.Block{}).
		Where("id IN ?", blockIds).
		Update("deleted_at", nil).Error; err != nil {
		tx.Rollback()
		return nil, exceptions.Block.FailedToUpdate().WithOrigin(err)
	}

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return nil, exceptions.Block.FailedToCommitTransaction().WithOrigin(err)
	}

	return &dtos.AppendBlockResDto{BlockPackId: reqDto.Body.BlockPackId, BlockIds: blockIds, CreatedAt: time.Now()}, nil
}

func (s *BlockService) AppendBlocks(
	ctx context.Context, reqDto *dtos.AppendBlocksReqDto,
) (*dtos.AppendBlocksResDto, *exceptions.Exception) {
	if err := validation.Validator.Struct(reqDto); err != nil {
		return nil, exceptions.Block.InvalidDto().WithOrigin(err)
	}

	tx := s.db.WithContext(ctx).Begin()

	res := dtos.AppendBlocksResDto{CreatedAt: time.Now()}

	type preparedAppendBlock struct {
		Index       int
		BlockPackId uuid.UUID
		RawBlocks   []dtos.RawFlattenedEditableBlock
		BlockIds    []uuid.UUID
	}

	preparedBlocks := make([]preparedAppendBlock, 0, len(reqDto.Body.AppendedBlocks))
	checkInputs := make([]inputs.BulkCheckBlockPackPermissionInput, 0, len(reqDto.Body.AppendedBlocks))
	for index, appendedBlock := range reqDto.Body.AppendedBlocks {
		rawBlocks, _, exception := s.editableBlockAdapter.FlattenToRaw(&appendedBlock.ArborizedEditableBlock)
		if exception != nil {
			res.FailedIndexes = append(res.FailedIndexes, index)
			continue
		}

		blockIds := make([]uuid.UUID, len(rawBlocks))
		for rawBlockIndex, rawBlock := range rawBlocks {
			blockIds[rawBlockIndex] = rawBlock.Id
		}

		preparedBlocks = append(preparedBlocks, preparedAppendBlock{
			Index:       index,
			BlockPackId: appendedBlock.BlockPackId,
			RawBlocks:   rawBlocks,
			BlockIds:    blockIds,
		})
		checkInputs = append(checkInputs, inputs.BulkCheckBlockPackPermissionInput{
			UserId: reqDto.ContextFields.UserId,
			Id:     appendedBlock.BlockPackId,
		})
	}

	if len(preparedBlocks) == 0 {
		tx.Rollback()
		res.IsAllSuccess = len(res.FailedIndexes) == 0
		return &res, nil
	}

	validBlockPacks, _, exception := s.blockPackRepository.BulkCheckPermissionsAndGetManyByIds(
		checkInputs,
		nil,
		[]enums.AccessControlPermission{
			enums.AccessControlPermission_Owner,
			enums.AccessControlPermission_Admin,
			enums.AccessControlPermission_Write,
		},
		options.WithTransactionDB(tx),
		options.WithOnlyDeleted(types.Ternary_Negative),
		options.WithLockingStrength(options.LockingStrengthUpdate),
	)
	if exception != nil {
		tx.Rollback()
		return nil, exception
	}

	blockPackIds := make([]uuid.UUID, 0, len(preparedBlocks))
	for index, preparedBlock := range preparedBlocks {
		if !validBlockPacks[index] {
			res.FailedIndexes = append(res.FailedIndexes, preparedBlock.Index)
			continue
		}
		blockPackIds = append(blockPackIds, preparedBlock.BlockPackId)
	}

	var tails []schemas.Block
	if err := tx.Model(&schemas.Block{}).
		Where("block_pack_id IN ? AND parent_block_id IS NULL AND next_block_id IS NULL AND deleted_at IS NULL", blockPackIds).
		Find(&tails).Error; err != nil {
		tx.Rollback()
		return nil, exceptions.Block.NotFound().WithOrigin(err)
	}

	tailByBlockPackId := make(map[uuid.UUID]*uuid.UUID, len(tails))
	for _, tail := range tails {
		tailId := tail.Id
		tailByBlockPackId[tail.BlockPackId] = &tailId
	}

	deletedAt := time.Now()
	newBlocks := make([]schemas.Block, 0)
	neighborUpdates := make(map[uuid.UUID]uuid.UUID)
	successItems := make([]preparedAppendBlock, 0, len(preparedBlocks))
	for index, preparedBlock := range preparedBlocks {
		if !validBlockPacks[index] {
			continue
		}

		rootId := preparedBlock.RawBlocks[0].Id
		prevBlockId := tailByBlockPackId[preparedBlock.BlockPackId]
		preparedBlock.RawBlocks[0].ParentBlockId = nil
		preparedBlock.RawBlocks[0].PrevBlockId = prevBlockId
		preparedBlock.RawBlocks[0].NextBlockId = nil

		if prevBlockId != nil {
			neighborUpdates[*prevBlockId] = rootId
		}
		tailByBlockPackId[preparedBlock.BlockPackId] = &rootId

		for _, rawBlock := range preparedBlock.RawBlocks {
			newBlocks = append(newBlocks, schemas.Block{
				Id:            rawBlock.Id,
				BlockPackId:   preparedBlock.BlockPackId,
				ParentBlockId: rawBlock.ParentBlockId,
				PrevBlockId:   rawBlock.PrevBlockId,
				NextBlockId:   rawBlock.NextBlockId,
				Type:          rawBlock.Type,
				Props:         rawBlock.Props,
				Content:       rawBlock.Content,
				DeletedAt:     &deletedAt,
			})
		}
		successItems = append(successItems, preparedBlock)
	}

	if len(newBlocks) > 0 {
		if err := tx.CreateInBatches(&newBlocks, constants.MaxBatchCreateBlockSize).Error; err != nil {
			tx.Rollback()
			return nil, exceptions.Block.FailedToCreate().WithOrigin(err)
		}
	}

	if len(neighborUpdates) > 0 {
		valuePlaceholders := make([]string, 0, len(neighborUpdates))
		valueArgs := make([]any, 0, len(neighborUpdates)*2)
		for prevBlockId, nextBlockId := range neighborUpdates {
			valuePlaceholders = append(valuePlaceholders, "(?::uuid, ?::uuid)")
			valueArgs = append(valueArgs, prevBlockId, nextBlockId)
		}

		if err := tx.Exec(fmt.Sprintf(`
			UPDATE "BlockTable" AS b
			SET next_block_id = v.next_block_id
			FROM (VALUES %s) AS v(id, next_block_id)
			WHERE b.id = v.id::uuid
		`, strings.Join(valuePlaceholders, ",")), valueArgs...).Error; err != nil {
			tx.Rollback()
			return nil, exceptions.Block.FailedToUpdate().WithOrigin(err)
		}
	}

	newBlockIds := make([]uuid.UUID, 0, len(newBlocks))
	for _, newBlock := range newBlocks {
		newBlockIds = append(newBlockIds, newBlock.Id)
	}

	if len(newBlockIds) > 0 {
		if err := tx.Model(&schemas.Block{}).
			Where("id IN ?", newBlockIds).
			Update("deleted_at", nil).Error; err != nil {
			tx.Rollback()
			return nil, exceptions.Block.FailedToUpdate().WithOrigin(err)
		}
	}

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return nil, exceptions.Block.FailedToCommitTransaction().WithOrigin(err)
	}

	for _, successItem := range successItems {
		res.SuccessIndexes = append(res.SuccessIndexes, successItem.Index)
		res.SuccessBlockPackAppendItems = append(res.SuccessBlockPackAppendItems, struct {
			BlockPackId uuid.UUID   `json:"blockPackId"`
			BlockIds    []uuid.UUID `json:"blockIds"`
		}{BlockPackId: successItem.BlockPackId, BlockIds: successItem.BlockIds})
	}

	res.IsAllSuccess = len(res.FailedIndexes) == 0
	return &res, nil
}

func (s *BlockService) InsertBlock(
	ctx context.Context, reqDto *dtos.InsertBlockReqDto,
) (*dtos.InsertBlockResDto, *exceptions.Exception) {
	if err := validation.Validator.Struct(reqDto); err != nil {
		return nil, exceptions.Block.InvalidDto().WithOrigin(err)
	}

	rawBlocks, _, exception := s.editableBlockAdapter.FlattenToRaw(&reqDto.Body.ArborizedEditableBlock)
	if exception != nil {
		return nil, exception
	}

	tx := s.db.WithContext(ctx).Begin()

	if !s.blockPackRepository.HasPermission(
		reqDto.Body.BlockPackId,
		reqDto.ContextFields.UserId,
		[]enums.AccessControlPermission{
			enums.AccessControlPermission_Owner,
			enums.AccessControlPermission_Admin,
			enums.AccessControlPermission_Write,
		},
		options.WithTransactionDB(tx),
		options.WithOnlyDeleted(types.Ternary_Negative),
	) {
		tx.Rollback()
		return nil, exceptions.Block.NoPermission("insert blocks to the block pack")
	}

	lockingStrength := options.LockingStrengthUpdate
	if err := tx.Scopes(scopes.Locking(&lockingStrength)).
		Where(`"BlockPackTable".id = ? AND "BlockPackTable".deleted_at IS NULL`, reqDto.Body.BlockPackId).
		First(&schemas.BlockPack{}).Error; err != nil {
		tx.Rollback()
		return nil, exceptions.BlockPack.NotFound().WithOrigin(err)
	}

	if reqDto.Body.ParentBlockId != nil {
		var parentBlockCount int64
		if err := tx.Model(&schemas.Block{}).
			Where("id = ? AND block_pack_id = ? AND deleted_at IS NULL", *reqDto.Body.ParentBlockId, reqDto.Body.BlockPackId).
			Count(&parentBlockCount).Error; err != nil {
			tx.Rollback()
			return nil, exceptions.Block.NotFound().WithOrigin(err)
		}
		if parentBlockCount == 0 {
			tx.Rollback()
			return nil, exceptions.Block.NotFound()
		}
	}

	rootId := rawBlocks[0].Id
	rawBlocks[0].ParentBlockId = reqDto.Body.ParentBlockId
	rawBlocks[0].PrevBlockId = reqDto.Body.PrevBlockId
	rawBlocks[0].NextBlockId = nil

	var nextBlockId *uuid.UUID
	if reqDto.Body.PrevBlockId != nil {
		var prevBlock schemas.Block
		query := tx.Model(&schemas.Block{}).
			Where("block_pack_id = ? AND deleted_at IS NULL", reqDto.Body.BlockPackId)
		if reqDto.Body.ParentBlockId == nil {
			query = query.Where("parent_block_id IS NULL")
		} else {
			query = query.Where("parent_block_id = ?", *reqDto.Body.ParentBlockId)
		}
		if err := query.
			Where("id = ?", *reqDto.Body.PrevBlockId).
			First(&prevBlock).Error; err != nil {
			tx.Rollback()
			return nil, exceptions.Block.NotFound().WithOrigin(err)
		}
		nextBlockId = prevBlock.NextBlockId
	} else {
		var head schemas.Block
		query := tx.Model(&schemas.Block{}).
			Where("block_pack_id = ? AND deleted_at IS NULL", reqDto.Body.BlockPackId)
		if reqDto.Body.ParentBlockId == nil {
			query = query.Where("parent_block_id IS NULL")
		} else {
			query = query.Where("parent_block_id = ?", *reqDto.Body.ParentBlockId)
		}
		if err := query.
			Where("prev_block_id IS NULL").
			First(&head).Error; err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			tx.Rollback()
			return nil, exceptions.Block.NotFound().WithOrigin(err)
		}
		if head.Id != uuid.Nil {
			nextBlockId = &head.Id
		}
	}

	rawBlocks[0].NextBlockId = nextBlockId

	deletedAt := time.Now()
	newBlocks := make([]schemas.Block, len(rawBlocks))
	blockIds := make([]uuid.UUID, len(rawBlocks))
	for index, rawBlock := range rawBlocks {
		blockIds[index] = rawBlock.Id
		newBlocks[index] = schemas.Block{
			Id:            rawBlock.Id,
			BlockPackId:   reqDto.Body.BlockPackId,
			ParentBlockId: rawBlock.ParentBlockId,
			PrevBlockId:   rawBlock.PrevBlockId,
			NextBlockId:   rawBlock.NextBlockId,
			Type:          rawBlock.Type,
			Props:         rawBlock.Props,
			Content:       rawBlock.Content,
			DeletedAt:     &deletedAt,
		}
	}

	if err := tx.CreateInBatches(&newBlocks, constants.MaxBatchCreateBlockSize).Error; err != nil {
		tx.Rollback()
		return nil, exceptions.Block.FailedToCreate().WithOrigin(err)
	}

	if reqDto.Body.PrevBlockId != nil {
		if err := tx.Model(&schemas.Block{}).
			Where("id = ?", *reqDto.Body.PrevBlockId).
			Update("next_block_id", rootId).Error; err != nil {
			tx.Rollback()
			return nil, exceptions.Block.FailedToUpdate().WithOrigin(err)
		}
	}

	if nextBlockId != nil {
		if err := tx.Model(&schemas.Block{}).
			Where("id = ?", *nextBlockId).
			Update("prev_block_id", rootId).Error; err != nil {
			tx.Rollback()
			return nil, exceptions.Block.FailedToUpdate().WithOrigin(err)
		}
	}

	if err := tx.Model(&schemas.Block{}).
		Where("id IN ?", blockIds).
		Update("deleted_at", nil).Error; err != nil {
		tx.Rollback()
		return nil, exceptions.Block.FailedToUpdate().WithOrigin(err)
	}

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return nil, exceptions.Block.FailedToCommitTransaction().WithOrigin(err)
	}

	return &dtos.InsertBlockResDto{CreatedAt: time.Now()}, nil
}

func (s *BlockService) InsertBlocks(
	ctx context.Context, reqDto *dtos.InsertBlocksReqDto,
) (*dtos.InsertBlocksResDto, *exceptions.Exception) {
	if err := validation.Validator.Struct(reqDto); err != nil {
		return nil, exceptions.Block.InvalidDto().WithOrigin(err)
	}

	tx := s.db.WithContext(ctx).Begin()

	res := dtos.InsertBlocksResDto{CreatedAt: time.Now()}

	type preparedInsertBlock struct {
		Index         int
		BlockPackId   uuid.UUID
		ParentBlockId *uuid.UUID
		PrevBlockId   *uuid.UUID
		RawBlocks     []dtos.RawFlattenedEditableBlock
		BlockIds      []uuid.UUID
	}

	preparedBlocks := make([]preparedInsertBlock, 0, len(reqDto.Body.InsertedBlocks))
	checkInputs := make([]inputs.BulkCheckBlockPackPermissionInput, 0, len(reqDto.Body.InsertedBlocks))
	parentBlockIds := make([]uuid.UUID, 0, len(reqDto.Body.InsertedBlocks))
	prevBlockIds := make([]uuid.UUID, 0, len(reqDto.Body.InsertedBlocks))
	for index, insertedBlock := range reqDto.Body.InsertedBlocks {
		rawBlocks, _, exception := s.editableBlockAdapter.FlattenToRaw(&insertedBlock.ArborizedEditableBlock)
		if exception != nil {
			res.FailedIndexes = append(res.FailedIndexes, index)
			continue
		}

		blockIds := make([]uuid.UUID, len(rawBlocks))
		for rawBlockIndex, rawBlock := range rawBlocks {
			blockIds[rawBlockIndex] = rawBlock.Id
		}
		if insertedBlock.ParentBlockId != nil {
			parentBlockIds = append(parentBlockIds, *insertedBlock.ParentBlockId)
		}
		if insertedBlock.PrevBlockId != nil {
			prevBlockIds = append(prevBlockIds, *insertedBlock.PrevBlockId)
		}

		preparedBlocks = append(preparedBlocks, preparedInsertBlock{
			Index:         index,
			BlockPackId:   insertedBlock.BlockPackId,
			ParentBlockId: insertedBlock.ParentBlockId,
			PrevBlockId:   insertedBlock.PrevBlockId,
			RawBlocks:     rawBlocks,
			BlockIds:      blockIds,
		})
		checkInputs = append(checkInputs, inputs.BulkCheckBlockPackPermissionInput{
			UserId: reqDto.ContextFields.UserId,
			Id:     insertedBlock.BlockPackId,
		})
	}

	if len(preparedBlocks) == 0 {
		tx.Rollback()
		res.IsAllSuccess = len(res.FailedIndexes) == 0
		return &res, nil
	}

	validBlockPacks, _, exception := s.blockPackRepository.BulkCheckPermissionsAndGetManyByIds(
		checkInputs,
		nil,
		[]enums.AccessControlPermission{
			enums.AccessControlPermission_Owner,
			enums.AccessControlPermission_Admin,
			enums.AccessControlPermission_Write,
		},
		options.WithTransactionDB(tx),
		options.WithOnlyDeleted(types.Ternary_Negative),
		options.WithLockingStrength(options.LockingStrengthUpdate),
	)
	if exception != nil {
		tx.Rollback()
		return nil, exception
	}

	var parentBlocks []schemas.Block
	if len(parentBlockIds) > 0 {
		if err := tx.Model(&schemas.Block{}).
			Where("id IN ? AND deleted_at IS NULL", parentBlockIds).
			Find(&parentBlocks).Error; err != nil {
			tx.Rollback()
			return nil, exceptions.Block.NotFound().WithOrigin(err)
		}
	}
	parentBlockById := make(map[uuid.UUID]schemas.Block, len(parentBlocks))
	for _, parentBlock := range parentBlocks {
		parentBlockById[parentBlock.Id] = parentBlock
	}

	var prevBlocks []schemas.Block
	if len(prevBlockIds) > 0 {
		if err := tx.Model(&schemas.Block{}).
			Where("id IN ? AND deleted_at IS NULL", prevBlockIds).
			Find(&prevBlocks).Error; err != nil {
			tx.Rollback()
			return nil, exceptions.Block.NotFound().WithOrigin(err)
		}
	}
	prevBlockById := make(map[uuid.UUID]schemas.Block, len(prevBlocks))
	for _, prevBlock := range prevBlocks {
		prevBlockById[prevBlock.Id] = prevBlock
	}

	blockPackIds := make([]uuid.UUID, 0, len(preparedBlocks))
	for index, preparedBlock := range preparedBlocks {
		if validBlockPacks[index] {
			blockPackIds = append(blockPackIds, preparedBlock.BlockPackId)
		}
	}

	var heads []schemas.Block
	if len(blockPackIds) > 0 {
		if err := tx.Model(&schemas.Block{}).
			Where("block_pack_id IN ? AND prev_block_id IS NULL AND deleted_at IS NULL", blockPackIds).
			Find(&heads).Error; err != nil {
			tx.Rollback()
			return nil, exceptions.Block.NotFound().WithOrigin(err)
		}
	}

	type siblingScope struct {
		BlockPackId   uuid.UUID
		ParentBlockId uuid.UUID
	}

	type insertedRootBlock struct {
		BlockPackId   uuid.UUID
		ParentBlockId *uuid.UUID
		NextBlockId   *uuid.UUID
		NewBlockIndex int
	}

	headByScope := make(map[siblingScope]*uuid.UUID, len(heads))
	for _, head := range heads {
		parentBlockId := uuid.Nil
		if head.ParentBlockId != nil {
			parentBlockId = *head.ParentBlockId
		}

		headId := head.Id
		headByScope[siblingScope{BlockPackId: head.BlockPackId, ParentBlockId: parentBlockId}] = &headId
	}

	deletedAt := time.Now()
	newBlocks := make([]schemas.Block, 0)
	prevNeighborUpdates := make(map[uuid.UUID]uuid.UUID)
	nextNeighborUpdates := make(map[uuid.UUID]uuid.UUID)
	successItems := make([]preparedInsertBlock, 0, len(preparedBlocks))
	usedInsertionPoints := make(map[string]bool)
	insertedRootBlockById := make(map[uuid.UUID]insertedRootBlock, len(preparedBlocks))
	for index, preparedBlock := range preparedBlocks {
		if !validBlockPacks[index] {
			res.FailedIndexes = append(res.FailedIndexes, preparedBlock.Index)
			continue
		}

		parentBlockId := uuid.Nil
		if preparedBlock.ParentBlockId != nil {
			parentBlockId = *preparedBlock.ParentBlockId
			parentBlock, exists := parentBlockById[parentBlockId]
			if !exists || parentBlock.BlockPackId != preparedBlock.BlockPackId {
				res.FailedIndexes = append(res.FailedIndexes, preparedBlock.Index)
				continue
			}
		}

		insertionPoint := fmt.Sprintf("head:%s:%s", preparedBlock.BlockPackId, parentBlockId)
		if preparedBlock.PrevBlockId != nil {
			insertionPoint = fmt.Sprintf("prev:%s", *preparedBlock.PrevBlockId)
		}
		if usedInsertionPoints[insertionPoint] {
			res.FailedIndexes = append(res.FailedIndexes, preparedBlock.Index)
			continue
		}
		usedInsertionPoints[insertionPoint] = true

		var nextBlockId *uuid.UUID
		if preparedBlock.PrevBlockId != nil {
			if prevBlock, exists := prevBlockById[*preparedBlock.PrevBlockId]; exists {
				if prevBlock.BlockPackId != preparedBlock.BlockPackId {
					res.FailedIndexes = append(res.FailedIndexes, preparedBlock.Index)
					continue
				}
				if preparedBlock.ParentBlockId == nil && prevBlock.ParentBlockId != nil ||
					preparedBlock.ParentBlockId != nil && (prevBlock.ParentBlockId == nil || *prevBlock.ParentBlockId != *preparedBlock.ParentBlockId) {
					res.FailedIndexes = append(res.FailedIndexes, preparedBlock.Index)
					continue
				}

				nextBlockId = prevBlock.NextBlockId
				prevNeighborUpdates[*preparedBlock.PrevBlockId] = preparedBlock.RawBlocks[0].Id
			} else if insertedRootBlock, exists := insertedRootBlockById[*preparedBlock.PrevBlockId]; exists {
				if insertedRootBlock.BlockPackId != preparedBlock.BlockPackId {
					res.FailedIndexes = append(res.FailedIndexes, preparedBlock.Index)
					continue
				}
				if preparedBlock.ParentBlockId == nil && insertedRootBlock.ParentBlockId != nil ||
					preparedBlock.ParentBlockId != nil && (insertedRootBlock.ParentBlockId == nil || *insertedRootBlock.ParentBlockId != *preparedBlock.ParentBlockId) {
					res.FailedIndexes = append(res.FailedIndexes, preparedBlock.Index)
					continue
				}

				nextBlockId = insertedRootBlock.NextBlockId
				nextRootId := preparedBlock.RawBlocks[0].Id
				newBlocks[insertedRootBlock.NewBlockIndex].NextBlockId = &nextRootId
				insertedRootBlock.NextBlockId = &nextRootId
				insertedRootBlockById[*preparedBlock.PrevBlockId] = insertedRootBlock
			} else {
				res.FailedIndexes = append(res.FailedIndexes, preparedBlock.Index)
				continue
			}
		} else {
			nextBlockId = headByScope[siblingScope{BlockPackId: preparedBlock.BlockPackId, ParentBlockId: parentBlockId}]
		}

		rootId := preparedBlock.RawBlocks[0].Id
		preparedBlock.RawBlocks[0].ParentBlockId = preparedBlock.ParentBlockId
		preparedBlock.RawBlocks[0].PrevBlockId = preparedBlock.PrevBlockId
		preparedBlock.RawBlocks[0].NextBlockId = nextBlockId

		if nextBlockId != nil {
			nextNeighborUpdates[*nextBlockId] = rootId
		}

		rootNewBlockIndex := len(newBlocks)
		for _, rawBlock := range preparedBlock.RawBlocks {
			newBlocks = append(newBlocks, schemas.Block{
				Id:            rawBlock.Id,
				BlockPackId:   preparedBlock.BlockPackId,
				ParentBlockId: rawBlock.ParentBlockId,
				PrevBlockId:   rawBlock.PrevBlockId,
				NextBlockId:   rawBlock.NextBlockId,
				Type:          rawBlock.Type,
				Props:         rawBlock.Props,
				Content:       rawBlock.Content,
				DeletedAt:     &deletedAt,
			})
		}
		insertedRootBlockById[rootId] = insertedRootBlock{
			BlockPackId:   preparedBlock.BlockPackId,
			ParentBlockId: preparedBlock.ParentBlockId,
			NextBlockId:   nextBlockId,
			NewBlockIndex: rootNewBlockIndex,
		}
		successItems = append(successItems, preparedBlock)
	}

	if len(newBlocks) > 0 {
		if err := tx.CreateInBatches(&newBlocks, constants.MaxBatchCreateBlockSize).Error; err != nil {
			tx.Rollback()
			return nil, exceptions.Block.FailedToCreate().WithOrigin(err)
		}
	}

	if len(prevNeighborUpdates) > 0 {
		valuePlaceholders := make([]string, 0, len(prevNeighborUpdates))
		valueArgs := make([]any, 0, len(prevNeighborUpdates)*2)
		for prevBlockId, nextBlockId := range prevNeighborUpdates {
			valuePlaceholders = append(valuePlaceholders, "(?::uuid, ?::uuid)")
			valueArgs = append(valueArgs, prevBlockId, nextBlockId)
		}

		if err := tx.Exec(fmt.Sprintf(`
			UPDATE "BlockTable" AS b
			SET next_block_id = v.next_block_id
			FROM (VALUES %s) AS v(id, next_block_id)
			WHERE b.id = v.id::uuid
		`, strings.Join(valuePlaceholders, ",")), valueArgs...).Error; err != nil {
			tx.Rollback()
			return nil, exceptions.Block.FailedToUpdate().WithOrigin(err)
		}
	}

	if len(nextNeighborUpdates) > 0 {
		valuePlaceholders := make([]string, 0, len(nextNeighborUpdates))
		valueArgs := make([]any, 0, len(nextNeighborUpdates)*2)
		for nextBlockId, prevBlockId := range nextNeighborUpdates {
			valuePlaceholders = append(valuePlaceholders, "(?::uuid, ?::uuid)")
			valueArgs = append(valueArgs, nextBlockId, prevBlockId)
		}

		if err := tx.Exec(fmt.Sprintf(`
			UPDATE "BlockTable" AS b
			SET prev_block_id = v.prev_block_id
			FROM (VALUES %s) AS v(id, prev_block_id)
			WHERE b.id = v.id::uuid
		`, strings.Join(valuePlaceholders, ",")), valueArgs...).Error; err != nil {
			tx.Rollback()
			return nil, exceptions.Block.FailedToUpdate().WithOrigin(err)
		}
	}

	newBlockIds := make([]uuid.UUID, 0, len(newBlocks))
	for _, newBlock := range newBlocks {
		newBlockIds = append(newBlockIds, newBlock.Id)
	}

	if len(newBlockIds) > 0 {
		if err := tx.Model(&schemas.Block{}).
			Where("id IN ?", newBlockIds).
			Update("deleted_at", nil).Error; err != nil {
			tx.Rollback()
			return nil, exceptions.Block.FailedToUpdate().WithOrigin(err)
		}
	}

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return nil, exceptions.Block.FailedToCommitTransaction().WithOrigin(err)
	}

	for _, successItem := range successItems {
		res.SuccessIndexes = append(res.SuccessIndexes, successItem.Index)
		res.SuccessBlockPackAndBlockIds = append(res.SuccessBlockPackAndBlockIds, struct {
			BlockPackId uuid.UUID   `json:"blockPackId"`
			BlockIds    []uuid.UUID `json:"blockIds"`
		}{BlockPackId: successItem.BlockPackId, BlockIds: successItem.BlockIds})
	}

	res.IsAllSuccess = len(res.FailedIndexes) == 0
	return &res, nil
}

func (s *BlockService) UpdateMyBlockById(
	ctx context.Context, reqDto *dtos.UpdateMyBlockByIdReqDto,
) (*dtos.UpdateMyBlockByIdResDto, *exceptions.Exception) {
	if err := validation.Validator.Struct(reqDto); err != nil {
		return nil, exceptions.Block.InvalidDto().WithOrigin(err)
	}

	tx := s.db.WithContext(ctx).Begin()

	var block schemas.Block
	if err := tx.Model(&schemas.Block{}).
		Scopes(s.blockScope.PassPermissionCheck(reqDto.Body.BlockId, reqDto.ContextFields.UserId, []enums.AccessControlPermission{
			enums.AccessControlPermission_Owner,
			enums.AccessControlPermission_Admin,
			enums.AccessControlPermission_Write,
		})).
		Scopes(s.blockScope.FilterOnlyDeleted(types.Ternary_Negative)).
		First(&block).Error; err != nil {
		tx.Rollback()
		return nil, exceptions.Block.NotFound().WithOrigin(err)
	}

	targetBlockPackId := block.BlockPackId
	if reqDto.Body.Values.BlockPackId != nil {
		targetBlockPackId = *reqDto.Body.Values.BlockPackId
	}

	if targetBlockPackId == uuid.Nil {
		tx.Rollback()
		return nil, exceptions.Block.InvalidDto()
	}

	targetParentBlockId := block.ParentBlockId
	if reqDto.Body.Values.ParentBlockId != nil || (reqDto.Body.SetNull != nil && (*reqDto.Body.SetNull)["ParentBlockId"]) {
		targetParentBlockId = reqDto.Body.Values.ParentBlockId
	}

	targetPrevBlockId := block.PrevBlockId
	if reqDto.Body.Values.PrevBlockId != nil || (reqDto.Body.SetNull != nil && (*reqDto.Body.SetNull)["PrevBlockId"]) {
		targetPrevBlockId = reqDto.Body.Values.PrevBlockId
	}

	if targetBlockPackId != block.BlockPackId || targetParentBlockId != block.ParentBlockId || targetPrevBlockId != block.PrevBlockId {
		if !s.blockPackRepository.HasPermission(
			targetBlockPackId,
			reqDto.ContextFields.UserId,
			[]enums.AccessControlPermission{
				enums.AccessControlPermission_Owner,
				enums.AccessControlPermission_Admin,
				enums.AccessControlPermission_Write,
			},
			options.WithTransactionDB(tx),
			options.WithOnlyDeleted(types.Ternary_Negative),
		) {
			tx.Rollback()
			return nil, exceptions.Block.NoPermission("move block to target block pack")
		}

		if block.BlockPackId == uuid.Nil {
			tx.Rollback()
			return nil, exceptions.BlockPack.NotFound()
		}

		lockingStrength := options.LockingStrengthUpdate
		if err := tx.Scopes(scopes.Locking(&lockingStrength)).
			Where(`"BlockPackTable".id = ? AND "BlockPackTable".deleted_at IS NULL`, block.BlockPackId).
			First(&schemas.BlockPack{}).Error; err != nil {
			tx.Rollback()
			return nil, exceptions.BlockPack.NotFound().WithOrigin(err)
		}

		if targetBlockPackId != block.BlockPackId {
			if targetBlockPackId == uuid.Nil {
				tx.Rollback()
				return nil, exceptions.BlockPack.NotFound()
			}
			if err := tx.Scopes(scopes.Locking(&lockingStrength)).
				Where(`"BlockPackTable".id = ? AND "BlockPackTable".deleted_at IS NULL`, targetBlockPackId).
				First(&schemas.BlockPack{}).Error; err != nil {
				tx.Rollback()
				return nil, exceptions.BlockPack.NotFound().WithOrigin(err)
			}
		}

		if targetParentBlockId != nil {
			var parentBlockCount int64
			if err := tx.Model(&schemas.Block{}).
				Where("id = ? AND block_pack_id = ? AND deleted_at IS NULL", *targetParentBlockId, targetBlockPackId).
				Count(&parentBlockCount).Error; err != nil {
				tx.Rollback()
				return nil, exceptions.Block.NotFound().WithOrigin(err)
			}
			if parentBlockCount == 0 {
				tx.Rollback()
				return nil, exceptions.Block.NotFound()
			}
		}

		var descendantRows []struct {
			Id uuid.UUID `gorm:"column:id"`
		}
		if err := tx.Raw(`
			WITH RECURSIVE descendants AS (
				SELECT id FROM "BlockTable" WHERE id = ?
				UNION ALL
				SELECT b.id FROM "BlockTable" b
				INNER JOIN descendants d ON b.parent_block_id = d.id
			)
			SELECT id FROM descendants
		`, block.Id).Scan(&descendantRows).Error; err != nil {
			tx.Rollback()
			return nil, exceptions.Block.NotFound().WithOrigin(err)
		}

		descendantIds := make([]uuid.UUID, len(descendantRows))
		for index, row := range descendantRows {
			descendantIds[index] = row.Id
		}

		if targetParentBlockId != nil {
			for _, descendantId := range descendantIds {
				if descendantId == *targetParentBlockId {
					tx.Rollback()
					return nil, exceptions.Block.InvalidDto()
				}
			}
		}

		if targetPrevBlockId != nil {
			for _, descendantId := range descendantIds {
				if descendantId == *targetPrevBlockId {
					tx.Rollback()
					return nil, exceptions.Block.InvalidDto()
				}
			}
		}

		if err := tx.Model(&schemas.Block{}).
			Where("id = ?", block.Id).
			Updates(map[string]any{"deleted_at": time.Now(), "prev_block_id": nil, "next_block_id": nil}).Error; err != nil {
			tx.Rollback()
			return nil, exceptions.Block.FailedToUpdate().WithOrigin(err)
		}

		if block.PrevBlockId != nil {
			if err := tx.Model(&schemas.Block{}).
				Where("id = ?", *block.PrevBlockId).
				Update("next_block_id", block.NextBlockId).Error; err != nil {
				tx.Rollback()
				return nil, exceptions.Block.FailedToUpdate().WithOrigin(err)
			}
		}

		if block.NextBlockId != nil {
			if err := tx.Model(&schemas.Block{}).
				Where("id = ?", *block.NextBlockId).
				Update("prev_block_id", block.PrevBlockId).Error; err != nil {
				tx.Rollback()
				return nil, exceptions.Block.FailedToUpdate().WithOrigin(err)
			}
		}

		if len(descendantIds) > 1 {
			if err := tx.Model(&schemas.Block{}).
				Where("id IN ? AND id <> ?", descendantIds, block.Id).
				Update("block_pack_id", targetBlockPackId).Error; err != nil {
				tx.Rollback()
				return nil, exceptions.Block.FailedToUpdate().WithOrigin(err)
			}
		}

		if err := tx.Model(&schemas.Block{}).
			Where("id = ?", block.Id).
			Updates(map[string]any{
				"block_pack_id":   targetBlockPackId,
				"parent_block_id": targetParentBlockId,
				"prev_block_id":   targetPrevBlockId,
				"next_block_id":   nil,
				"updated_at":      time.Now(),
			}).Error; err != nil {
			tx.Rollback()
			return nil, exceptions.Block.FailedToUpdate().WithOrigin(err)
		}

		var nextBlockId *uuid.UUID
		if targetPrevBlockId != nil {
			var prevBlock schemas.Block
			query := tx.Model(&schemas.Block{}).
				Where("block_pack_id = ? AND deleted_at IS NULL", targetBlockPackId)
			if targetParentBlockId == nil {
				query = query.Where("parent_block_id IS NULL")
			} else {
				query = query.Where("parent_block_id = ?", *targetParentBlockId)
			}
			if err := query.
				Where("id = ?", *targetPrevBlockId).
				First(&prevBlock).Error; err != nil {
				tx.Rollback()
				return nil, exceptions.Block.NotFound().WithOrigin(err)
			}
			nextBlockId = prevBlock.NextBlockId
		} else {
			var head schemas.Block
			query := tx.Model(&schemas.Block{}).
				Where("block_pack_id = ? AND deleted_at IS NULL", targetBlockPackId)
			if targetParentBlockId == nil {
				query = query.Where("parent_block_id IS NULL")
			} else {
				query = query.Where("parent_block_id = ?", *targetParentBlockId)
			}
			if err := query.
				Where("prev_block_id IS NULL").
				First(&head).Error; err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
				tx.Rollback()
				return nil, exceptions.Block.NotFound().WithOrigin(err)
			}
			if head.Id != uuid.Nil {
				nextBlockId = &head.Id
			}
		}

		if targetPrevBlockId != nil {
			if err := tx.Model(&schemas.Block{}).
				Where("id = ?", *targetPrevBlockId).
				Update("next_block_id", block.Id).Error; err != nil {
				tx.Rollback()
				return nil, exceptions.Block.FailedToUpdate().WithOrigin(err)
			}
		}

		if nextBlockId != nil {
			if err := tx.Model(&schemas.Block{}).
				Where("id = ?", *nextBlockId).
				Update("prev_block_id", block.Id).Error; err != nil {
				tx.Rollback()
				return nil, exceptions.Block.FailedToUpdate().WithOrigin(err)
			}
		}

		if err := tx.Model(&schemas.Block{}).
			Where("id = ?", block.Id).
			Updates(map[string]any{"next_block_id": nextBlockId, "deleted_at": nil}).Error; err != nil {
			tx.Rollback()
			return nil, exceptions.Block.FailedToUpdate().WithOrigin(err)
		}
	}

	updates := map[string]any{}
	if reqDto.Body.Values.Type != nil {
		updates["type"] = *reqDto.Body.Values.Type
	}
	if reqDto.Body.Values.Props != nil {
		updates["props"] = datatypes.JSON(*reqDto.Body.Values.Props)
	}
	if reqDto.Body.Values.Content != nil {
		updates["content"] = datatypes.JSON(*reqDto.Body.Values.Content)
	}

	if len(updates) > 0 {
		updates["updated_at"] = time.Now()
		if err := tx.Model(&schemas.Block{}).
			Where("id = ?", reqDto.Body.BlockId).
			Updates(updates).Error; err != nil {
			tx.Rollback()
			return nil, exceptions.Block.FailedToUpdate().WithOrigin(err)
		}
	}

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return nil, exceptions.Block.FailedToCommitTransaction().WithOrigin(err)
	}

	return &dtos.UpdateMyBlockByIdResDto{UpdatedAt: time.Now()}, nil
}

func (s *BlockService) UpdateMyBlocksByIds(
	ctx context.Context, reqDto *dtos.UpdateMyBlocksByIdsReqDto,
) (*dtos.UpdateMyBlocksByIdsResDto, *exceptions.Exception) {
	if err := validation.Validator.Struct(reqDto); err != nil {
		return nil, exceptions.Block.InvalidDto().WithOrigin(err)
	}

	res := dtos.UpdateMyBlocksByIdsResDto{UpdatedAt: time.Now()}
	hasStructuralUpdate := false
	for _, updatedBlock := range reqDto.Body.UpdatedBlocks {
		if updatedBlock.Values.BlockPackId != nil ||
			updatedBlock.Values.ParentBlockId != nil ||
			updatedBlock.Values.PrevBlockId != nil ||
			util.CheckSetNull(updatedBlock.SetNull, "ParentBlockId") ||
			util.CheckSetNull(updatedBlock.SetNull, "PrevBlockId") {
			hasStructuralUpdate = true
			break
		}
	}

	if !hasStructuralUpdate {
		db := s.db.WithContext(ctx)

		bulkInputs := make([]inputs.BulkUpdateBlockInput, 0, len(reqDto.Body.UpdatedBlocks))
		indexes := make([]int, 0, len(reqDto.Body.UpdatedBlocks))
		for index, updatedBlock := range reqDto.Body.UpdatedBlocks {
			if updatedBlock.Values.NextBlockId != nil || util.CheckSetNull(updatedBlock.SetNull, "NextBlockId") {
				res.FailedIndexes = append(res.FailedIndexes, index)
				continue
			}

			var props *datatypes.JSON
			if updatedBlock.Values.Props != nil {
				value := datatypes.JSON(*updatedBlock.Values.Props)
				props = &value
			}

			var content *datatypes.JSON
			if updatedBlock.Values.Content != nil {
				value := datatypes.JSON(*updatedBlock.Values.Content)
				content = &value
			}

			bulkInputs = append(bulkInputs, inputs.BulkUpdateBlockInput{
				UserId: reqDto.ContextFields.UserId,
				Id:     updatedBlock.BlockId,
				PartialUpdateInput: inputs.PartialUpdateBlockInput{
					Values: inputs.UpdateBlockInput{
						Type:    updatedBlock.Values.Type,
						Props:   props,
						Content: content,
					},
				},
			})
			indexes = append(indexes, index)
		}

		if len(bulkInputs) == 0 {
			res.IsAllSuccess = len(res.FailedIndexes) == 0
			return &res, nil
		}

		successes, exception := s.blockRepository.BulkUpdateMany(
			bulkInputs,
			options.WithDB(db),
			options.WithOnlyDeleted(types.Ternary_Negative),
		)
		if exception != nil {
			return nil, exception
		}

		for index, success := range successes {
			if !success {
				res.FailedIndexes = append(res.FailedIndexes, indexes[index])
				continue
			}

			res.SuccessIndexes = append(res.SuccessIndexes, indexes[index])
		}

		res.IsAllSuccess = len(res.FailedIndexes) == 0
		return &res, nil
	}

	tx := s.db.WithContext(ctx).Begin()

	type siblingKey struct {
		BlockPackId      uuid.UUID
		HasParentBlockId bool
		ParentBlockId    uuid.UUID
	}

	type structuralUpdateInput struct {
		Index        int
		BlockId      uuid.UUID
		Values       inputs.UpdateBlockInput
		SetNull      *map[string]bool
		HasStructure bool
	}

	updateInputs := make([]structuralUpdateInput, 0, len(reqDto.Body.UpdatedBlocks))
	updateInputIndexesByBlockId := make(map[uuid.UUID]int, len(reqDto.Body.UpdatedBlocks))
	candidateBlockIds := make([]uuid.UUID, 0, len(reqDto.Body.UpdatedBlocks))
	for index, updatedBlock := range reqDto.Body.UpdatedBlocks {
		if updatedBlock.Values.NextBlockId != nil || util.CheckSetNull(updatedBlock.SetNull, "NextBlockId") {
			res.FailedIndexes = append(res.FailedIndexes, index)
			continue
		}

		if _, exist := updateInputIndexesByBlockId[updatedBlock.BlockId]; exist {
			res.FailedIndexes = append(res.FailedIndexes, index)
			continue
		}
		updateInputIndexesByBlockId[updatedBlock.BlockId] = index

		var props *datatypes.JSON
		if updatedBlock.Values.Props != nil {
			value := datatypes.JSON(*updatedBlock.Values.Props)
			props = &value
		}

		var content *datatypes.JSON
		if updatedBlock.Values.Content != nil {
			value := datatypes.JSON(*updatedBlock.Values.Content)
			content = &value
		}

		updateInputs = append(updateInputs, structuralUpdateInput{
			Index:   index,
			BlockId: updatedBlock.BlockId,
			Values: inputs.UpdateBlockInput{
				BlockPackId:   updatedBlock.Values.BlockPackId,
				ParentBlockId: updatedBlock.Values.ParentBlockId,
				PrevBlockId:   updatedBlock.Values.PrevBlockId,
				Type:          updatedBlock.Values.Type,
				Props:         props,
				Content:       content,
			},
			SetNull: updatedBlock.SetNull,
			HasStructure: updatedBlock.Values.BlockPackId != nil ||
				updatedBlock.Values.ParentBlockId != nil ||
				updatedBlock.Values.PrevBlockId != nil ||
				util.CheckSetNull(updatedBlock.SetNull, "ParentBlockId") ||
				util.CheckSetNull(updatedBlock.SetNull, "PrevBlockId"),
		})
		candidateBlockIds = append(candidateBlockIds, updatedBlock.BlockId)
	}

	if len(updateInputs) == 0 {
		tx.Rollback()
		res.IsAllSuccess = len(res.FailedIndexes) == 0
		return &res, nil
	}

	lockingStrength := options.LockingStrengthUpdate
	allowedPermissions := []enums.AccessControlPermission{
		enums.AccessControlPermission_Owner,
		enums.AccessControlPermission_Admin,
		enums.AccessControlPermission_Write,
	}

	var permittedBlocks []schemas.Block
	if err := tx.Model(&schemas.Block{}).
		Scopes(s.blockScope.PassPermissionChecks(candidateBlockIds, reqDto.ContextFields.UserId, allowedPermissions)).
		Scopes(s.blockScope.FilterOnlyDeleted(types.Ternary_Negative)).
		Scopes(scopes.Locking(&lockingStrength)).
		Find(&permittedBlocks).Error; err != nil {
		tx.Rollback()
		return nil, exceptions.Block.NotFound().WithOrigin(err)
	}

	permittedBlocksById := make(map[uuid.UUID]schemas.Block, len(permittedBlocks))
	affectedBlockPackIdsMap := make(map[uuid.UUID]bool, len(permittedBlocks))
	for _, permittedBlock := range permittedBlocks {
		permittedBlocksById[permittedBlock.Id] = permittedBlock
		affectedBlockPackIdsMap[permittedBlock.BlockPackId] = true
	}

	for _, updateInput := range updateInputs {
		block, exist := permittedBlocksById[updateInput.BlockId]
		if !exist {
			continue
		}
		if updateInput.Values.BlockPackId != nil {
			affectedBlockPackIdsMap[*updateInput.Values.BlockPackId] = true
		} else {
			affectedBlockPackIdsMap[block.BlockPackId] = true
		}
	}

	affectedBlockPackIds := make([]uuid.UUID, 0, len(affectedBlockPackIdsMap))
	for blockPackId := range affectedBlockPackIdsMap {
		affectedBlockPackIds = append(affectedBlockPackIds, blockPackId)
	}

	var permittedBlockPacks []schemas.BlockPack
	if err := tx.Model(&schemas.BlockPack{}).
		Scopes(s.blockPackScope.PassPermissionChecks(affectedBlockPackIds, reqDto.ContextFields.UserId, allowedPermissions)).
		Scopes(s.blockPackScope.FilterOnlyDeleted(types.Ternary_Negative)).
		Scopes(scopes.Locking(&lockingStrength)).
		Find(&permittedBlockPacks).Error; err != nil {
		tx.Rollback()
		return nil, exceptions.BlockPack.NotFound().WithOrigin(err)
	}

	permittedBlockPackIds := make(map[uuid.UUID]bool, len(permittedBlockPacks))
	for _, permittedBlockPack := range permittedBlockPacks {
		permittedBlockPackIds[permittedBlockPack.Id] = true
	}

	var blocks []schemas.Block
	if err := tx.Model(&schemas.Block{}).
		Where("block_pack_id IN ? AND deleted_at IS NULL", affectedBlockPackIds).
		Scopes(scopes.Locking(&lockingStrength)).
		Find(&blocks).Error; err != nil {
		tx.Rollback()
		return nil, exceptions.Block.NotFound().WithOrigin(err)
	}

	originalBlocksById := make(map[uuid.UUID]schemas.Block, len(blocks))
	currentBlocksById := make(map[uuid.UUID]schemas.Block, len(blocks))
	siblingOrders := make(map[siblingKey][]uuid.UUID, len(blocks))
	unorderedSiblingIds := make(map[siblingKey]map[uuid.UUID]bool, len(blocks))
	for _, block := range blocks {
		originalBlocksById[block.Id] = block
		currentBlocksById[block.Id] = block

		key := siblingKey{BlockPackId: block.BlockPackId}
		if block.ParentBlockId != nil {
			key.HasParentBlockId = true
			key.ParentBlockId = *block.ParentBlockId
		}
		if unorderedSiblingIds[key] == nil {
			unorderedSiblingIds[key] = map[uuid.UUID]bool{}
		}
		unorderedSiblingIds[key][block.Id] = true
	}

	affectedSiblingKeys := make(map[siblingKey]bool, len(updateInputs)*2)
	for _, updateInput := range updateInputs {
		if !updateInput.HasStructure {
			continue
		}

		block, exist := currentBlocksById[updateInput.BlockId]
		if !exist {
			continue
		}

		oldKey := siblingKey{BlockPackId: block.BlockPackId}
		if block.ParentBlockId != nil {
			oldKey.HasParentBlockId = true
			oldKey.ParentBlockId = *block.ParentBlockId
		}
		affectedSiblingKeys[oldKey] = true

		processedBlock := schemas.Block{}
		if err := copier.Copy(&processedBlock, &block); err != nil {
			tx.Rollback()
			return nil, exceptions.Block.InvalidDto().WithOrigin(err)
		}
		processedBlock, err := util.PartialUpdatePreprocess(updateInput.Values, updateInput.SetNull, processedBlock)
		if err != nil {
			tx.Rollback()
			return nil, exceptions.Block.InvalidDto().WithOrigin(err)
		}

		newKey := siblingKey{BlockPackId: processedBlock.BlockPackId}
		if processedBlock.ParentBlockId != nil {
			newKey.HasParentBlockId = true
			newKey.ParentBlockId = *processedBlock.ParentBlockId
		}
		affectedSiblingKeys[newKey] = true
	}

	for key, siblingIds := range unorderedSiblingIds {
		if !affectedSiblingKeys[key] {
			continue
		}

		var headId uuid.UUID
		for siblingId := range siblingIds {
			block := currentBlocksById[siblingId]
			if block.PrevBlockId == nil {
				if headId != uuid.Nil {
					tx.Rollback()
					return nil, exceptions.Block.InvalidDto("duplicate sibling head detected")
				}
				headId = block.Id
			}
		}

		if headId == uuid.Nil && len(siblingIds) > 0 {
			for siblingId := range siblingIds {
				block := currentBlocksById[siblingId]
				if block.PrevBlockId != nil && !siblingIds[*block.PrevBlockId] {
					headId = block.Id
					break
				}
			}
		}

		if headId == uuid.Nil && len(siblingIds) > 0 {
			tx.Rollback()
			return nil, exceptions.Block.InvalidDto("missing sibling head detected")
		}

		visitedSiblingIds := make(map[uuid.UUID]bool, len(siblingIds))
		currentSiblingId := headId
		for currentSiblingId != uuid.Nil {
			if visitedSiblingIds[currentSiblingId] || !siblingIds[currentSiblingId] {
				break
			}

			visitedSiblingIds[currentSiblingId] = true
			siblingOrders[key] = append(siblingOrders[key], currentSiblingId)

			block := currentBlocksById[currentSiblingId]
			if block.NextBlockId == nil {
				currentSiblingId = uuid.Nil
			} else {
				currentSiblingId = *block.NextBlockId
			}
		}

		if len(visitedSiblingIds) != len(siblingIds) {
			for siblingId := range siblingIds {
				if !visitedSiblingIds[siblingId] {
					siblingOrders[key] = append(siblingOrders[key], siblingId)
				}
			}
		}
	}

	successInputIndexes := make(map[int]bool, len(updateInputs))
	for _, updateInput := range updateInputs {
		block, exist := currentBlocksById[updateInput.BlockId]
		if !exist {
			res.FailedIndexes = append(res.FailedIndexes, updateInput.Index)
			continue
		}

		processedBlock := schemas.Block{}
		if err := copier.Copy(&processedBlock, &block); err != nil {
			tx.Rollback()
			return nil, exceptions.Block.InvalidDto().WithOrigin(err)
		}
		processedBlock, err := util.PartialUpdatePreprocess(updateInput.Values, updateInput.SetNull, processedBlock)
		if err != nil {
			tx.Rollback()
			return nil, exceptions.Block.InvalidDto().WithOrigin(err)
		}

		if processedBlock.BlockPackId == uuid.Nil || !permittedBlockPackIds[processedBlock.BlockPackId] {
			res.FailedIndexes = append(res.FailedIndexes, updateInput.Index)
			continue
		}

		if processedBlock.ParentBlockId != nil {
			parentBlock, exist := currentBlocksById[*processedBlock.ParentBlockId]
			if !exist || parentBlock.BlockPackId != processedBlock.BlockPackId || parentBlock.Id == processedBlock.Id {
				res.FailedIndexes = append(res.FailedIndexes, updateInput.Index)
				continue
			}
		}

		if processedBlock.PrevBlockId != nil {
			prevBlock, exist := currentBlocksById[*processedBlock.PrevBlockId]
			if !exist || prevBlock.BlockPackId != processedBlock.BlockPackId || prevBlock.Id == processedBlock.Id {
				res.FailedIndexes = append(res.FailedIndexes, updateInput.Index)
				continue
			}
			if (prevBlock.ParentBlockId == nil) != (processedBlock.ParentBlockId == nil) {
				res.FailedIndexes = append(res.FailedIndexes, updateInput.Index)
				continue
			}
			if prevBlock.ParentBlockId != nil && processedBlock.ParentBlockId != nil && *prevBlock.ParentBlockId != *processedBlock.ParentBlockId {
				res.FailedIndexes = append(res.FailedIndexes, updateInput.Index)
				continue
			}
		}

		descendantIds := map[uuid.UUID]bool{}
		descendantFrontier := []uuid.UUID{processedBlock.Id}
		for len(descendantFrontier) > 0 {
			parentId := descendantFrontier[0]
			descendantFrontier = descendantFrontier[1:]

			for _, possibleDescendant := range currentBlocksById {
				if possibleDescendant.ParentBlockId != nil && *possibleDescendant.ParentBlockId == parentId {
					if descendantIds[possibleDescendant.Id] {
						continue
					}
					descendantIds[possibleDescendant.Id] = true
					descendantFrontier = append(descendantFrontier, possibleDescendant.Id)
				}
			}
		}

		if processedBlock.ParentBlockId != nil && descendantIds[*processedBlock.ParentBlockId] {
			res.FailedIndexes = append(res.FailedIndexes, updateInput.Index)
			continue
		}
		if processedBlock.PrevBlockId != nil && descendantIds[*processedBlock.PrevBlockId] {
			res.FailedIndexes = append(res.FailedIndexes, updateInput.Index)
			continue
		}

		if updateInput.HasStructure {
			oldKey := siblingKey{BlockPackId: block.BlockPackId}
			if block.ParentBlockId != nil {
				oldKey.HasParentBlockId = true
				oldKey.ParentBlockId = *block.ParentBlockId
			}

			oldOrder := siblingOrders[oldKey]
			removedIndex := -1
			for index, blockId := range oldOrder {
				if blockId == block.Id {
					removedIndex = index
					break
				}
			}
			if removedIndex == -1 {
				res.FailedIndexes = append(res.FailedIndexes, updateInput.Index)
				continue
			}
			siblingOrders[oldKey] = append(oldOrder[:removedIndex], oldOrder[removedIndex+1:]...)

			newKey := siblingKey{BlockPackId: processedBlock.BlockPackId}
			if processedBlock.ParentBlockId != nil {
				newKey.HasParentBlockId = true
				newKey.ParentBlockId = *processedBlock.ParentBlockId
			}

			newOrder := siblingOrders[newKey]
			if processedBlock.PrevBlockId == nil {
				siblingOrders[newKey] = append([]uuid.UUID{processedBlock.Id}, newOrder...)
			} else {
				insertedIndex := -1
				for index, blockId := range newOrder {
					if blockId == *processedBlock.PrevBlockId {
						insertedIndex = index
						break
					}
				}
				if insertedIndex == -1 {
					siblingOrders[oldKey] = oldOrder
					res.FailedIndexes = append(res.FailedIndexes, updateInput.Index)
					continue
				}

				newOrder = append(newOrder[:insertedIndex+1], append([]uuid.UUID{processedBlock.Id}, newOrder[insertedIndex+1:]...)...)
				siblingOrders[newKey] = newOrder
			}
		}

		currentBlocksById[processedBlock.Id] = processedBlock
		successInputIndexes[updateInput.Index] = true
	}

	for iteration := 0; iteration < len(currentBlocksById); iteration++ {
		hasChanged := false
		for blockId, block := range currentBlocksById {
			if block.ParentBlockId == nil {
				continue
			}

			parentBlock, exist := currentBlocksById[*block.ParentBlockId]
			if !exist {
				tx.Rollback()
				return nil, exceptions.Block.InvalidDto("parent block not found")
			}

			if block.BlockPackId != parentBlock.BlockPackId {
				block.BlockPackId = parentBlock.BlockPackId
				currentBlocksById[blockId] = block
				hasChanged = true
			}
		}

		if !hasChanged {
			break
		}
	}

	finalSiblingOrders := make(map[siblingKey][]uuid.UUID, len(siblingOrders))
	for _, siblingOrder := range siblingOrders {
		for _, blockId := range siblingOrder {
			block := currentBlocksById[blockId]

			key := siblingKey{BlockPackId: block.BlockPackId}
			if block.ParentBlockId != nil {
				key.HasParentBlockId = true
				key.ParentBlockId = *block.ParentBlockId
			}
			finalSiblingOrders[key] = append(finalSiblingOrders[key], blockId)
		}
	}

	for _, siblingOrder := range finalSiblingOrders {
		for index, blockId := range siblingOrder {
			block := currentBlocksById[blockId]

			block.PrevBlockId = nil
			if index > 0 {
				prevBlockId := siblingOrder[index-1]
				block.PrevBlockId = &prevBlockId
			}

			block.NextBlockId = nil
			if index < len(siblingOrder)-1 {
				nextBlockId := siblingOrder[index+1]
				block.NextBlockId = &nextBlockId
			}

			currentBlocksById[blockId] = block
		}
	}

	valuePlaceholders := make([]string, 0, len(currentBlocksById))
	valueArgs := make([]any, 0, len(currentBlocksById)*9)
	blockIdsForRelink := make([]uuid.UUID, 0, len(currentBlocksById))
	for blockId, block := range currentBlocksById {
		originalBlock := originalBlocksById[blockId]
		hasChanged := block.BlockPackId != originalBlock.BlockPackId ||
			block.Type != originalBlock.Type ||
			string(block.Props) != string(originalBlock.Props) ||
			string(block.Content) != string(originalBlock.Content)

		if !hasChanged && (block.ParentBlockId == nil) != (originalBlock.ParentBlockId == nil) {
			hasChanged = true
		}
		if !hasChanged && block.ParentBlockId != nil && originalBlock.ParentBlockId != nil && *block.ParentBlockId != *originalBlock.ParentBlockId {
			hasChanged = true
		}
		if !hasChanged && (block.PrevBlockId == nil) != (originalBlock.PrevBlockId == nil) {
			hasChanged = true
		}
		if !hasChanged && block.PrevBlockId != nil && originalBlock.PrevBlockId != nil && *block.PrevBlockId != *originalBlock.PrevBlockId {
			hasChanged = true
		}
		if !hasChanged && (block.NextBlockId == nil) != (originalBlock.NextBlockId == nil) {
			hasChanged = true
		}
		if !hasChanged && block.NextBlockId != nil && originalBlock.NextBlockId != nil && *block.NextBlockId != *originalBlock.NextBlockId {
			hasChanged = true
		}

		var parentBlockIdArg any
		if block.ParentBlockId != nil {
			parentBlockIdArg = *block.ParentBlockId
		}

		var prevBlockIdArg any
		if block.PrevBlockId != nil {
			prevBlockIdArg = *block.PrevBlockId
		}

		var nextBlockIdArg any
		if block.NextBlockId != nil {
			nextBlockIdArg = *block.NextBlockId
		}

		blockIdsForRelink = append(blockIdsForRelink, block.Id)
		valuePlaceholders = append(valuePlaceholders, "(?::uuid, ?::uuid, ?::uuid, ?::uuid, ?::uuid, ?::\"BlockType\", ?::jsonb, ?::jsonb, ?::boolean)")
		valueArgs = append(valueArgs,
			block.Id,
			block.BlockPackId,
			parentBlockIdArg,
			prevBlockIdArg,
			nextBlockIdArg,
			block.Type,
			block.Props,
			block.Content,
			hasChanged,
		)
	}

	if len(valuePlaceholders) > 0 && len(successInputIndexes) > 0 {
		if err := tx.Model(&schemas.Block{}).
			Where("id IN ?", blockIdsForRelink).
			Update("deleted_at", time.Now()).Error; err != nil {
			tx.Rollback()
			return nil, exceptions.Block.FailedToUpdate().WithOrigin(err)
		}

		if err := tx.Exec(fmt.Sprintf(`
			UPDATE "BlockTable" AS b
			SET
				block_pack_id = v.block_pack_id,
				parent_block_id = v.parent_block_id,
				prev_block_id = v.prev_block_id,
				next_block_id = v.next_block_id,
				type = v.type,
				props = v.props,
				content = v.content,
				updated_at = CASE
					WHEN v.should_touch_updated_at THEN NOW()
					ELSE b.updated_at
				END
			FROM (VALUES %s) AS v(id, block_pack_id, parent_block_id, prev_block_id, next_block_id, type, props, content, should_touch_updated_at)
			WHERE b.id = v.id::uuid
		`, strings.Join(valuePlaceholders, ",")), valueArgs...).Error; err != nil {
			tx.Rollback()
			return nil, exceptions.Block.FailedToUpdate().WithOrigin(err)
		}

		if err := tx.Model(&schemas.Block{}).
			Where("id IN ?", blockIdsForRelink).
			Update("deleted_at", nil).Error; err != nil {
			tx.Rollback()
			return nil, exceptions.Block.FailedToUpdate().WithOrigin(err)
		}
	}

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return nil, exceptions.Block.FailedToCommitTransaction().WithOrigin(err)
	}

	for _, updateInput := range updateInputs {
		if !successInputIndexes[updateInput.Index] {
			continue
		}

		block := currentBlocksById[updateInput.BlockId]
		res.SuccessIndexes = append(res.SuccessIndexes, updateInput.Index)
		res.SuccessBlockPackAndBlockIds = append(res.SuccessBlockPackAndBlockIds, struct {
			BlockPackId uuid.UUID   `json:"blockPackId"`
			BlockIds    []uuid.UUID `json:"blockIds"`
		}{BlockPackId: block.BlockPackId, BlockIds: []uuid.UUID{block.Id}})
	}

	res.IsAllSuccess = len(res.FailedIndexes) == 0
	return &res, nil
}

func (s *BlockService) RestoreMyBlockById(
	ctx context.Context, reqDto *dtos.RestoreMyBlockByIdReqDto,
) (*dtos.RestoreMyBlockByIdResDto, *exceptions.Exception) {
	if err := validation.Validator.Struct(reqDto); err != nil {
		return nil, exceptions.Block.InvalidDto().WithOrigin(err)
	}

	tx := s.db.WithContext(ctx).Begin()

	var block schemas.Block
	if err := tx.Model(&schemas.Block{}).
		Scopes(s.blockScope.PassPermissionCheck(reqDto.Body.BlockId, reqDto.ContextFields.UserId, []enums.AccessControlPermission{
			enums.AccessControlPermission_Owner,
			enums.AccessControlPermission_Admin,
			enums.AccessControlPermission_Write,
		})).
		Scopes(s.blockScope.FilterOnlyDeleted(types.Ternary_Positive)).
		First(&block).Error; err != nil {
		tx.Rollback()
		return nil, exceptions.Block.NotFound().WithOrigin(err)
	}

	if block.BlockPackId == uuid.Nil {
		tx.Rollback()
		return nil, exceptions.BlockPack.NotFound()
	}

	lockingStrength := options.LockingStrengthUpdate
	if err := tx.Scopes(scopes.Locking(&lockingStrength)).
		Where(`"BlockPackTable".id = ? AND "BlockPackTable".deleted_at IS NULL`, block.BlockPackId).
		First(&schemas.BlockPack{}).Error; err != nil {
		tx.Rollback()
		return nil, exceptions.BlockPack.NotFound().WithOrigin(err)
	}

	var tail schemas.Block
	query := tx.Model(&schemas.Block{}).
		Where("block_pack_id = ? AND deleted_at IS NULL", block.BlockPackId)
	if block.ParentBlockId == nil {
		query = query.Where("parent_block_id IS NULL")
	} else {
		query = query.Where("parent_block_id = ?", *block.ParentBlockId)
	}
	if err := query.
		Where("next_block_id IS NULL").
		First(&tail).Error; err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		tx.Rollback()
		return nil, exceptions.Block.NotFound().WithOrigin(err)
	}

	var prevBlockId *uuid.UUID
	if tail.Id != uuid.Nil {
		prevBlockId = &tail.Id
	}

	var descendantRows []struct {
		Id uuid.UUID `gorm:"column:id"`
	}
	if err := tx.Raw(`
		WITH RECURSIVE descendants AS (
			SELECT id FROM "BlockTable" WHERE id = ?
			UNION ALL
			SELECT b.id FROM "BlockTable" b
			INNER JOIN descendants d ON b.parent_block_id = d.id
		)
		SELECT id FROM descendants
	`, block.Id).Scan(&descendantRows).Error; err != nil {
		tx.Rollback()
		return nil, exceptions.Block.NotFound().WithOrigin(err)
	}

	ids := make([]uuid.UUID, len(descendantRows))
	for index, row := range descendantRows {
		ids[index] = row.Id
	}

	if err := tx.Model(&schemas.Block{}).
		Where("id IN ? AND id <> ?", ids, block.Id).
		Update("deleted_at", nil).Error; err != nil {
		tx.Rollback()
		return nil, exceptions.Block.FailedToUpdate().WithOrigin(err)
	}

	if err := tx.Model(&schemas.Block{}).
		Where("id = ?", block.Id).
		Updates(map[string]any{"prev_block_id": prevBlockId, "next_block_id": nil}).Error; err != nil {
		tx.Rollback()
		return nil, exceptions.Block.FailedToUpdate().WithOrigin(err)
	}

	if prevBlockId != nil {
		if err := tx.Model(&schemas.Block{}).
			Where("id = ?", *prevBlockId).
			Update("next_block_id", block.Id).Error; err != nil {
			tx.Rollback()
			return nil, exceptions.Block.FailedToUpdate().WithOrigin(err)
		}
	}

	if err := tx.Model(&schemas.Block{}).
		Where("id = ?", block.Id).
		Update("deleted_at", nil).Error; err != nil {
		tx.Rollback()
		return nil, exceptions.Block.FailedToUpdate().WithOrigin(err)
	}

	if err := tx.
		Where("id = ?", block.Id).
		First(&block).Error; err != nil {
		tx.Rollback()
		return nil, exceptions.Block.NotFound().WithOrigin(err)
	}

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return nil, exceptions.Block.FailedToCommitTransaction().WithOrigin(err)
	}

	res := dtos.RestoreMyBlockByIdResDto{
		Id:            block.Id,
		BlockPackId:   block.BlockPackId,
		ParentBlockId: block.ParentBlockId,
		PrevBlockId:   block.PrevBlockId,
		NextBlockId:   block.NextBlockId,
		Type:          block.Type,
		Props:         block.Props,
		Content:       block.Content,
		DeletedAt:     block.DeletedAt,
		UpdatedAt:     block.UpdatedAt,
		CreatedAt:     block.CreatedAt,
	}

	return &res, nil
}

func (s *BlockService) RestoreMyBlocksByIds(
	ctx context.Context, reqDto *dtos.RestoreMyBlocksByIdsReqDto,
) (*dtos.RestoreMyBlocksByIdsResDto, *exceptions.Exception) {
	if err := validation.Validator.Struct(reqDto); err != nil {
		return nil, exceptions.Block.InvalidDto().WithOrigin(err)
	}

	tx := s.db.WithContext(ctx).Begin()

	blocks, exception := s.blockRepository.CheckPermissionsAndGetManyByIds(
		reqDto.Body.BlockIds,
		reqDto.ContextFields.UserId,
		nil,
		[]enums.AccessControlPermission{
			enums.AccessControlPermission_Owner,
			enums.AccessControlPermission_Admin,
			enums.AccessControlPermission_Write,
		},
		options.WithTransactionDB(tx),
		options.WithOnlyDeleted(types.Ternary_Positive),
		options.WithLockingStrength(options.LockingStrengthUpdate),
	)
	if exception != nil {
		tx.Rollback()
		return nil, exception
	}

	blockPackIds := make([]uuid.UUID, 0, len(blocks))
	for _, block := range blocks {
		blockPackIds = append(blockPackIds, block.BlockPackId)
	}

	var tails []schemas.Block
	if err := tx.Model(&schemas.Block{}).
		Where("block_pack_id IN ? AND next_block_id IS NULL AND deleted_at IS NULL", blockPackIds).
		Find(&tails).Error; err != nil {
		tx.Rollback()
		return nil, exceptions.Block.NotFound().WithOrigin(err)
	}

	type siblingScope struct {
		BlockPackId   uuid.UUID
		ParentBlockId uuid.UUID
	}

	tailByScope := make(map[siblingScope]*uuid.UUID, len(tails))
	for _, tail := range tails {
		parentBlockId := uuid.Nil
		if tail.ParentBlockId != nil {
			parentBlockId = *tail.ParentBlockId
		}

		tailId := tail.Id
		tailByScope[siblingScope{BlockPackId: tail.BlockPackId, ParentBlockId: parentBlockId}] = &tailId
	}

	valuePlaceholders := make([]string, 0, len(blocks))
	valueArgs := make([]any, 0, len(blocks)*3)
	neighborUpdates := make(map[uuid.UUID]uuid.UUID)
	blockIds := make([]uuid.UUID, len(blocks))
	for index, block := range blocks {
		parentBlockId := uuid.Nil
		if block.ParentBlockId != nil {
			parentBlockId = *block.ParentBlockId
		}

		scope := siblingScope{BlockPackId: block.BlockPackId, ParentBlockId: parentBlockId}
		prevBlockId := tailByScope[scope]
		if prevBlockId != nil {
			neighborUpdates[*prevBlockId] = block.Id
		}

		var prevBlockIdArg any
		if prevBlockId != nil {
			prevBlockIdArg = *prevBlockId
		}

		tailByScope[scope] = &block.Id
		blockIds[index] = block.Id
		valuePlaceholders = append(valuePlaceholders, "(?::uuid, ?::uuid, NULL::uuid)")
		valueArgs = append(valueArgs, block.Id, prevBlockIdArg)
	}

	if err := tx.Exec(`
		WITH RECURSIVE descendants AS (
			SELECT id FROM "BlockTable" WHERE id IN ?
			UNION ALL
			SELECT b.id FROM "BlockTable" b
			INNER JOIN descendants d ON b.parent_block_id = d.id
		)
		UPDATE "BlockTable"
		SET deleted_at = NULL
		WHERE id IN (SELECT id FROM descendants)
			AND id NOT IN ?
	`, blockIds, blockIds).Error; err != nil {
		tx.Rollback()
		return nil, exceptions.Block.FailedToUpdate().WithOrigin(err)
	}

	if err := tx.Exec(fmt.Sprintf(`
		WITH target(id, prev_block_id, next_block_id) AS (
			VALUES %s
		)
		UPDATE "BlockTable" AS b
		SET
			prev_block_id = t.prev_block_id::uuid,
			next_block_id = t.next_block_id::uuid
		FROM target AS t
		WHERE b.id = t.id::uuid
	`, strings.Join(valuePlaceholders, ",")), valueArgs...).Error; err != nil {
		tx.Rollback()
		return nil, exceptions.Block.FailedToUpdate().WithOrigin(err)
	}

	if len(neighborUpdates) > 0 {
		neighborPlaceholders := make([]string, 0, len(neighborUpdates))
		neighborArgs := make([]any, 0, len(neighborUpdates)*2)
		for prevBlockId, nextBlockId := range neighborUpdates {
			neighborPlaceholders = append(neighborPlaceholders, "(?::uuid, ?::uuid)")
			neighborArgs = append(neighborArgs, prevBlockId, nextBlockId)
		}

		if err := tx.Exec(fmt.Sprintf(`
			UPDATE "BlockTable" AS b
			SET next_block_id = v.next_block_id
			FROM (VALUES %s) AS v(id, next_block_id)
			WHERE b.id = v.id::uuid
		`, strings.Join(neighborPlaceholders, ",")), neighborArgs...).Error; err != nil {
			tx.Rollback()
			return nil, exceptions.Block.FailedToUpdate().WithOrigin(err)
		}
	}

	if err := tx.Model(&schemas.Block{}).
		Where("id IN ?", blockIds).
		Update("deleted_at", nil).Error; err != nil {
		tx.Rollback()
		return nil, exceptions.Block.FailedToUpdate().WithOrigin(err)
	}

	var restoredBlocks []schemas.Block
	if err := tx.Model(&schemas.Block{}).
		Where("id IN ?", blockIds).
		Find(&restoredBlocks).Error; err != nil {
		tx.Rollback()
		return nil, exceptions.Block.NotFound().WithOrigin(err)
	}

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return nil, exceptions.Block.FailedToCommitTransaction().WithOrigin(err)
	}

	res := make(dtos.RestoreMyBlocksByIdsResDto, len(restoredBlocks))
	for index, block := range restoredBlocks {
		res[index] = dtos.RestoreMyBlockByIdResDto{
			Id:            block.Id,
			BlockPackId:   block.BlockPackId,
			ParentBlockId: block.ParentBlockId,
			PrevBlockId:   block.PrevBlockId,
			NextBlockId:   block.NextBlockId,
			Type:          block.Type,
			Props:         block.Props,
			Content:       block.Content,
			DeletedAt:     block.DeletedAt,
			UpdatedAt:     block.UpdatedAt,
			CreatedAt:     block.CreatedAt,
		}
	}

	return &res, nil
}

func (s *BlockService) DeleteMyBlockById(
	ctx context.Context, reqDto *dtos.DeleteMyBlockByIdReqDto,
) (*dtos.DeleteMyBlockByIdResDto, *exceptions.Exception) {
	if err := validation.Validator.Struct(reqDto); err != nil {
		return nil, exceptions.Block.InvalidDto().WithOrigin(err)
	}

	tx := s.db.WithContext(ctx).Begin()

	var block schemas.Block
	if err := tx.Model(&schemas.Block{}).
		Scopes(s.blockScope.PassPermissionCheck(reqDto.Body.BlockId, reqDto.ContextFields.UserId, []enums.AccessControlPermission{
			enums.AccessControlPermission_Owner,
			enums.AccessControlPermission_Admin,
			enums.AccessControlPermission_Write,
		})).
		Scopes(s.blockScope.FilterOnlyDeleted(types.Ternary_Negative)).
		First(&block).Error; err != nil {
		tx.Rollback()
		return nil, exceptions.Block.NotFound().WithOrigin(err)
	}

	if block.BlockPackId == uuid.Nil {
		tx.Rollback()
		return nil, exceptions.BlockPack.NotFound()
	}

	lockingStrength := options.LockingStrengthUpdate
	if err := tx.Scopes(scopes.Locking(&lockingStrength)).
		Where(`"BlockPackTable".id = ? AND "BlockPackTable".deleted_at IS NULL`, block.BlockPackId).
		First(&schemas.BlockPack{}).Error; err != nil {
		tx.Rollback()
		return nil, exceptions.BlockPack.NotFound().WithOrigin(err)
	}

	var descendantRows []struct {
		Id uuid.UUID `gorm:"column:id"`
	}
	if err := tx.Raw(`
		WITH RECURSIVE descendants AS (
			SELECT id FROM "BlockTable" WHERE id = ?
			UNION ALL
			SELECT b.id FROM "BlockTable" b
			INNER JOIN descendants d ON b.parent_block_id = d.id
		)
		SELECT id FROM descendants
	`, block.Id).Scan(&descendantRows).Error; err != nil {
		tx.Rollback()
		return nil, exceptions.Block.NotFound().WithOrigin(err)
	}

	ids := make([]uuid.UUID, len(descendantRows))
	for index, row := range descendantRows {
		ids[index] = row.Id
	}

	if err := tx.Model(&schemas.Block{}).
		Where("id = ?", block.Id).
		Updates(map[string]any{"deleted_at": time.Now(), "prev_block_id": nil, "next_block_id": nil}).Error; err != nil {
		tx.Rollback()
		return nil, exceptions.Block.FailedToUpdate().WithOrigin(err)
	}

	if block.PrevBlockId != nil {
		if err := tx.Model(&schemas.Block{}).
			Where("id = ?", *block.PrevBlockId).
			Update("next_block_id", block.NextBlockId).Error; err != nil {
			tx.Rollback()
			return nil, exceptions.Block.FailedToUpdate().WithOrigin(err)
		}
	}

	if block.NextBlockId != nil {
		if err := tx.Model(&schemas.Block{}).
			Where("id = ?", *block.NextBlockId).
			Update("prev_block_id", block.PrevBlockId).Error; err != nil {
			tx.Rollback()
			return nil, exceptions.Block.FailedToUpdate().WithOrigin(err)
		}
	}

	if err := tx.Model(&schemas.Block{}).
		Where("id IN ? AND id <> ?", ids, block.Id).
		Update("deleted_at", time.Now()).Error; err != nil {
		tx.Rollback()
		return nil, exceptions.Block.FailedToUpdate().WithOrigin(err)
	}

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return nil, exceptions.Block.FailedToCommitTransaction().WithOrigin(err)
	}

	return &dtos.DeleteMyBlockByIdResDto{DeletedAt: time.Now()}, nil
}

func (s *BlockService) DeleteMyBlocksByIds(
	ctx context.Context, reqDto *dtos.DeleteMyBlocksByIdsReqDto,
) (*dtos.DeleteMyBlocksByIdsResDto, *exceptions.Exception) {
	if err := validation.Validator.Struct(reqDto); err != nil {
		return nil, exceptions.Block.InvalidDto().WithOrigin(err)
	}

	tx := s.db.WithContext(ctx).Begin()

	blocks, exception := s.blockRepository.CheckPermissionsAndGetManyByIds(
		reqDto.Body.BlockIds,
		reqDto.ContextFields.UserId,
		nil,
		[]enums.AccessControlPermission{
			enums.AccessControlPermission_Owner,
			enums.AccessControlPermission_Admin,
			enums.AccessControlPermission_Write,
		},
		options.WithTransactionDB(tx),
		options.WithOnlyDeleted(types.Ternary_Negative),
		options.WithLockingStrength(options.LockingStrengthUpdate),
	)
	if exception != nil {
		tx.Rollback()
		return nil, exception
	}

	blockIds := make([]uuid.UUID, len(blocks))
	valuePlaceholders := make([]string, 0, len(blocks))
	valueArgs := make([]any, 0, len(blocks)*3)
	for index, block := range blocks {
		blockIds[index] = block.Id
		valuePlaceholders = append(valuePlaceholders, "(?::uuid, ?::uuid, ?::uuid)")
		valueArgs = append(valueArgs, block.Id, block.PrevBlockId, block.NextBlockId)
	}

	if err := tx.Exec(fmt.Sprintf(`
		WITH target(id, prev_block_id, next_block_id) AS (
			VALUES %s
		)
		UPDATE "BlockTable" AS b
		SET
			deleted_at = NOW(),
			prev_block_id = NULL,
			next_block_id = NULL
		FROM target AS t
		WHERE b.id = t.id::uuid
	`, strings.Join(valuePlaceholders, ",")), valueArgs...).Error; err != nil {
		tx.Rollback()
		return nil, exceptions.Block.FailedToUpdate().WithOrigin(err)
	}

	if err := tx.Exec(fmt.Sprintf(`
		WITH target(id, prev_block_id, next_block_id) AS (
			VALUES %s
		)
		UPDATE "BlockTable" AS b
		SET next_block_id = t.next_block_id::uuid
		FROM target AS t
		WHERE b.id = t.prev_block_id::uuid
	`, strings.Join(valuePlaceholders, ",")), valueArgs...).Error; err != nil {
		tx.Rollback()
		return nil, exceptions.Block.FailedToUpdate().WithOrigin(err)
	}

	if err := tx.Exec(fmt.Sprintf(`
		WITH target(id, prev_block_id, next_block_id) AS (
			VALUES %s
		)
		UPDATE "BlockTable" AS b
		SET prev_block_id = t.prev_block_id::uuid
		FROM target AS t
		WHERE b.id = t.next_block_id::uuid
	`, strings.Join(valuePlaceholders, ",")), valueArgs...).Error; err != nil {
		tx.Rollback()
		return nil, exceptions.Block.FailedToUpdate().WithOrigin(err)
	}

	if err := tx.Exec(`
		WITH RECURSIVE descendants AS (
			SELECT id FROM "BlockTable" WHERE id IN ?
			UNION ALL
			SELECT b.id FROM "BlockTable" b
			INNER JOIN descendants d ON b.parent_block_id = d.id
		)
		UPDATE "BlockTable"
		SET deleted_at = NOW()
		WHERE id IN (SELECT id FROM descendants)
			AND id NOT IN ?
	`, blockIds, blockIds).Error; err != nil {
		tx.Rollback()
		return nil, exceptions.Block.FailedToUpdate().WithOrigin(err)
	}

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return nil, exceptions.Block.FailedToCommitTransaction().WithOrigin(err)
	}

	return &dtos.DeleteMyBlocksByIdsResDto{DeletedAt: time.Now()}, nil
}

func (s *BlockService) SearchPrivateBlocks(
	ctx context.Context, userId uuid.UUID, gqlInput gqlmodels.SearchBlockInput,
) (*gqlmodels.SearchBlockConnection, *exceptions.Exception) {
	startTime := time.Now()

	db := s.db.WithContext(ctx)

	query := db.Model(&schemas.Block{}).
		Select(`"BlockTable".*`).
		Joins(`INNER JOIN "BlockPackTable" ON "BlockPackTable".id = "BlockTable".block_pack_id`).
		Joins(`INNER JOIN "SubShelfTable" ON "SubShelfTable".id = "BlockPackTable".parent_sub_shelf_id`).
		Joins(`INNER JOIN "UsersToShelvesTable" uts ON uts.root_shelf_id = "SubShelfTable".root_shelf_id`).
		Where("uts.user_id = ? AND uts.permission IN ?", userId, []enums.AccessControlPermission{
			enums.AccessControlPermission_Owner,
			enums.AccessControlPermission_Admin,
			enums.AccessControlPermission_Write,
			enums.AccessControlPermission_Read,
		}).
		Scopes(s.blockScope.FilterOnlyDeleted(types.Ternary_Negative)).
		Scopes(s.blockPackScope.FilterOnlyDeleted(types.Ternary_Negative)).
		Scopes(s.subShelfScope.FilterOnlyDeleted(types.Ternary_Negative)).
		Scopes(s.blockScope.IncludePreloads([]schemas.BlockRelation{schemas.BlockRelation_Children}))

	if len(strings.ReplaceAll(gqlInput.Query, " ", "")) > 0 {
		pattern := "%" + gqlInput.Query + "%"
		query = query.Where(`"BlockTable".content::text ILIKE ? OR "BlockTable".props::text ILIKE ? OR "BlockTable".type::text ILIKE ?`, pattern, pattern, pattern)
	}

	if gqlInput.After != nil && len(strings.ReplaceAll(*gqlInput.After, " ", "")) > 0 {
		searchCursor, err := searchcursor.Decode[gqlmodels.SearchBlockCursorFields](*gqlInput.After)
		if err != nil {
			return nil, exceptions.Search.FailedToDecode().WithOrigin(err)
		}

		query = query.Where(`"BlockTable".id > ?`, searchCursor.Fields.ID)
	}

	if gqlInput.SortBy != nil && gqlInput.SortOrder != nil {
		cending := gqlmodels.SearchSortOrderAsc.String()
		if *gqlInput.SortOrder == gqlmodels.SearchSortOrderDesc {
			cending = gqlmodels.SearchSortOrderDesc.String()
		}

		switch *gqlInput.SortBy {
		case gqlmodels.SearchBlockSortByType:
			query = query.Order(`"BlockTable".type ` + cending).Order(`"BlockTable".updated_at ` + cending).Order(`"BlockTable".created_at ` + cending)
		case gqlmodels.SearchBlockSortByLastUpdate:
			query = query.Order(`"BlockTable".updated_at ` + cending).Order(`"BlockTable".type ` + cending).Order(`"BlockTable".created_at ` + cending)
		case gqlmodels.SearchBlockSortByCreatedAt:
			query = query.Order(`"BlockTable".created_at ` + cending).Order(`"BlockTable".type ` + cending).Order(`"BlockTable".updated_at ` + cending)
		default:
			query = query.Order(`"BlockTable".type ` + cending).Order(`"BlockTable".updated_at ` + cending).Order(`"BlockTable".created_at ` + cending)
		}
	}

	limit := constants.DefaultSearchLimit
	if gqlInput.First != nil && *gqlInput.First > 0 {
		limit = int(*gqlInput.First)
	}
	limit = min(limit, constants.MaxSearchLimit)
	query = query.Limit(limit + 1)

	var blocks []schemas.Block
	if err := query.Find(&blocks).Error; err != nil {
		return nil, exceptions.Block.NotFound().WithOrigin(err)
	}

	hasNextPage := len(blocks) > limit
	searchEdges := make([]*gqlmodels.SearchBlockEdge, len(blocks))
	for index, block := range blocks {
		searchCursor := searchcursor.SearchCursor[gqlmodels.SearchBlockCursorFields]{Fields: gqlmodels.SearchBlockCursorFields{ID: block.Id}}
		encodedSearchCursor, err := searchCursor.Encode()
		if err != nil {
			return nil, exceptions.Search.FailedToEncode().WithOrigin(err)
		}
		if encodedSearchCursor == nil {
			return nil, exceptions.Search.FailedToUnmarshalSearchCursor()
		}

		searchEdges[index] = &gqlmodels.SearchBlockEdge{EncodedSearchCursor: *encodedSearchCursor, Node: block.ToPrivateBlock()}
	}

	searchPageInfo := &gqlmodels.SearchPageInfo{
		HasNextPage:     hasNextPage,
		HasPreviousPage: gqlInput.After != nil && len(strings.ReplaceAll(*gqlInput.After, " ", "")) > 0,
	}

	if hasNextPage {
		searchEdges = searchEdges[:limit]
	}

	if len(searchEdges) > 0 {
		searchPageInfo.StartEncodedSearchCursor = &searchEdges[0].EncodedSearchCursor
		searchPageInfo.EndEncodedSearchCursor = &searchEdges[len(searchEdges)-1].EncodedSearchCursor
	}

	searchTime := float64(time.Since(startTime).Nanoseconds()) / 1e6

	return &gqlmodels.SearchBlockConnection{
		SearchEdges:    searchEdges,
		SearchPageInfo: searchPageInfo,
		TotalCount:     int32(len(searchEdges)),
		SearchTime:     searchTime,
	}, nil
}
