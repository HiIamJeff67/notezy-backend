package services

import (
	"context"
	"strings"
	"time"

	"github.com/google/uuid"
	pg "github.com/lib/pq"
	"gorm.io/gorm"

	caches "github.com/HiIamJeff67/notezy-backend/app/caches"
	contexts "github.com/HiIamJeff67/notezy-backend/app/contexts"
	dtos "github.com/HiIamJeff67/notezy-backend/app/dtos"
	exceptions "github.com/HiIamJeff67/notezy-backend/app/exceptions"
	gqlmodels "github.com/HiIamJeff67/notezy-backend/app/graphql/models"
	inputs "github.com/HiIamJeff67/notezy-backend/app/models/inputs"
	repositories "github.com/HiIamJeff67/notezy-backend/app/models/repositories"
	schemas "github.com/HiIamJeff67/notezy-backend/app/models/schemas"
	scopes "github.com/HiIamJeff67/notezy-backend/app/models/scopes"
	blockpacksql "github.com/HiIamJeff67/notezy-backend/app/models/sqls/block_pack"
	logs "github.com/HiIamJeff67/notezy-backend/app/monitor/logs"
	options "github.com/HiIamJeff67/notezy-backend/app/options"
	validation "github.com/HiIamJeff67/notezy-backend/app/validation"
	constants "github.com/HiIamJeff67/notezy-backend/shared/constants"
	searchcursor "github.com/HiIamJeff67/notezy-backend/shared/lib/searchcursor"
	types "github.com/HiIamJeff67/notezy-backend/shared/types"
)

type BlockPackServiceInterface interface {
	GetMyBlockPackById(ctx context.Context, reqDto *dtos.GetMyBlockPackByIdReqDto) (*dtos.GetMyBlockPackByIdResDto, *exceptions.Exception)
	GetMyBlockPackAndItsParentById(ctx context.Context, reqDto *dtos.GetMyBlockPackAndItsParentByIdReqDto) (*dtos.GetMyBlockPackAndItsParentByIdResDto, *exceptions.Exception)
	GetMyBlockPacksByParentSubShelfId(ctx context.Context, reqDto *dtos.GetMyBlockPacksByParentSubShelfIdReqDto) (*dtos.GetMyBlockPacksByParentSubShelfIdResDto, *exceptions.Exception)
	GetAllMyBlockPacksByRootShelfId(ctx context.Context, reqDto *dtos.GetAllMyBlockPacksByRootShelfIdReqDto) (*dtos.GetAllMyBlockPacksByRootShelfIdResDto, *exceptions.Exception)
	CreateBlockPack(ctx context.Context, reqDto *dtos.CreateBlockPackReqDto) (*dtos.CreateBlockPackResDto, *exceptions.Exception)
	CreateBlockPacks(ctx context.Context, reqDto *dtos.CreateBlockPacksReqDto) (*dtos.CreateBlockPacksResDto, *exceptions.Exception)
	UpdateMyBlockPackById(ctx context.Context, reqDto *dtos.UpdateMyBlockPackByIdReqDto) (*dtos.UpdateMyBlockPackByIdResDto, *exceptions.Exception)
	UpdateMyBlockPacksByIds(ctx context.Context, reqDto *dtos.UpdateMyBlockPacksByIdsReqDto) (*dtos.UpdateMyBlockPacksByIdsResDto, *exceptions.Exception)
	MoveMyBlockPackByParentSubShelfId(ctx context.Context, reqDto *dtos.MoveMyBlockPackByParentSubShelfIdReqDto) (*dtos.MoveMyBlockPackByParentSubShelfIdResDto, *exceptions.Exception)
	MoveMyBlockPacksByParentSubShelfId(ctx context.Context, reqDto *dtos.MoveMyBlockPacksByParentSubShelfIdReqDto) (*dtos.MoveMyBlockPacksByParentSubShelfIdResDto, *exceptions.Exception)
	MoveMyBlockPacksByParentSubShelfIds(ctx context.Context, reqDto *dtos.MoveMyBlockPacksByParentSubShelfIdsReqDto) (*dtos.MoveMyBlockPacksByParentSubShelfIdsResDto, *exceptions.Exception)
	RestoreMyBlockPackById(ctx context.Context, reqDto *dtos.RestoreMyBlockPackByIdReqDto) (*dtos.RestoreMyBlockPackByIdResDto, *exceptions.Exception)
	RestoreMyBlockPacksByIds(ctx context.Context, reqDto *dtos.RestoreMyBlockPacksByIdsReqDto) (*dtos.RestoreMyBlockPacksByIdsResDto, *exceptions.Exception)
	DeleteMyBlockPackById(ctx context.Context, reqDto *dtos.DeleteMyBlockPackByIdReqDto) (*dtos.DeleteMyBlockPackByIdResDto, *exceptions.Exception)
	DeleteMyBlockPacksByIds(ctx context.Context, reqDto *dtos.DeleteMyBlockPacksByIdsReqDto) (*dtos.DeleteMyBlockPacksByIdsResDto, *exceptions.Exception)

	SearchPrivateBlockPacks(ctx context.Context, userId uuid.UUID, gqlInput gqlmodels.SearchBlockPackInput) (*gqlmodels.SearchBlockPackConnection, *exceptions.Exception)
}

type BlockPackService struct {
	db                  *gorm.DB
	blockPackScope      scopes.BlockPackScopeInterface
	subShelfRepository  repositories.SubShelfRepositoryInterface
	blockPackRepository repositories.BlockPackRepositoryInterface
	realtimeLeaseStore  *caches.RealtimeLeaseStore
}

func NewBlockPackService(
	db *gorm.DB,
	blockPackScope scopes.BlockPackScopeInterface,
	subShelfRepository repositories.SubShelfRepositoryInterface,
	blockPackRepository repositories.BlockPackRepositoryInterface,
	realtimeLeaseStore *caches.RealtimeLeaseStore,
) BlockPackServiceInterface {
	if realtimeLeaseStore == nil {
		realtimeLeaseStore = caches.NewRealtimeLeaseStore(caches.RedisClientMap)
	}
	return &BlockPackService{
		db:                  db,
		blockPackScope:      blockPackScope,
		subShelfRepository:  subShelfRepository,
		blockPackRepository: blockPackRepository,
		realtimeLeaseStore:  realtimeLeaseStore,
	}
}

func (s *BlockPackService) GetMyBlockPackById(
	ctx context.Context, reqDto *dtos.GetMyBlockPackByIdReqDto,
) (*dtos.GetMyBlockPackByIdResDto, *exceptions.Exception) {
	if err := validation.Validator.Struct(reqDto); err != nil {
		return nil, exceptions.BlockPack.InvalidDto().WithOrigin(err)
	}

	db := s.db.WithContext(ctx)

	allowedPermissions, exception := contexts.GetAllowedPermissions(ctx)
	if exception != nil {
		return nil, exception
	}

	onlyDeleted := types.Ternary_Neutral
	if reqDto.Param.IsDeleted != nil {
		if *reqDto.Param.IsDeleted {
			onlyDeleted = types.Ternary_Positive
		} else {
			onlyDeleted = types.Ternary_Negative
		}
	}

	blockPack, exception := s.blockPackRepository.CheckPermissionAndGetOneById(
		reqDto.Param.BlockPackId,
		reqDto.ContextFields.UserId,
		[]schemas.BlockPackRelation{schemas.BlockPackRelation_YjsDocument},
		allowedPermissions,
		options.WithDB(db),
		options.WithAllowedPermissions(allowedPermissions),
		options.WithOnlyDeleted(onlyDeleted),
	)
	if exception != nil {
		return nil, exception
	}

	resDto := &dtos.GetMyBlockPackByIdResDto{
		Id:                  blockPack.Id,
		ParentSubShelfId:    blockPack.ParentSubShelfId,
		Name:                blockPack.Name,
		Icon:                blockPack.Icon,
		HeaderBackgroundURL: blockPack.HeaderBackgroundURL,
		BlockCount:          blockPack.BlockCount,
		DeletedAt:           blockPack.DeletedAt,
		UpdatedAt:           blockPack.UpdatedAt,
		CreatedAt:           blockPack.CreatedAt,
	}
	if blockPack.YjsDocument != nil {
		resDto.LastUpdateSequence = blockPack.YjsDocument.LastUpdateSequence
		resDto.CompactedUntilSequence = blockPack.YjsDocument.CompactedUntilSequence
		resDto.ProjectedUntilSequence = blockPack.YjsDocument.ProjectedUntilSequence
		resDto.IsProjectionCurrent = blockPack.YjsDocument.ProjectedUntilSequence >= blockPack.YjsDocument.LastUpdateSequence
	}

	return resDto, nil
}

func (s *BlockPackService) GetMyBlockPackAndItsParentById(
	ctx context.Context, reqDto *dtos.GetMyBlockPackAndItsParentByIdReqDto,
) (*dtos.GetMyBlockPackAndItsParentByIdResDto, *exceptions.Exception) {
	if err := validation.Validator.Struct(reqDto); err != nil {
		return nil, exceptions.BlockPack.InvalidDto().WithOrigin(err)
	}

	db := s.db.WithContext(ctx)

	allowedPermissions, exception := contexts.GetAllowedPermissions(ctx)
	if exception != nil {
		return nil, exception
	}

	onlyDeleted := types.Ternary_Neutral
	if reqDto.Param.IsDeleted != nil {
		if *reqDto.Param.IsDeleted {
			onlyDeleted = types.Ternary_Positive
		} else {
			onlyDeleted = types.Ternary_Negative
		}
	}

	resDto := dtos.GetMyBlockPackAndItsParentByIdResDto{}
	err := db.Raw(blockpacksql.GetMyBlockPackAndItsParentByIdSQL,
		reqDto.Param.BlockPackId, reqDto.ContextFields.UserId, pg.Array(allowedPermissions), onlyDeleted,
	).Row().
		Scan(&resDto.Id,
			&resDto.Name,
			&resDto.Icon,
			&resDto.HeaderBackgroundURL,
			&resDto.BlockCount,
			&resDto.LastUpdateSequence,
			&resDto.CompactedUntilSequence,
			&resDto.ProjectedUntilSequence,
			&resDto.IsProjectionCurrent,
			&resDto.DeletedAt,
			&resDto.UpdatedAt,
			&resDto.CreatedAt,
			&resDto.RootShelfId,
			&resDto.Permission,
			&resDto.ParentSubShelfId,
			&resDto.ParentSubShelfName,
			&resDto.ParentSubShelfPrevSubShelfId,
			&resDto.ParentSubShelfPath,
			&resDto.ParentSubShelfDeletedAt,
			&resDto.ParentSubShelfUpdatedAt,
			&resDto.ParentSubShelfCreatedAt)
	if err != nil {
		return nil, exceptions.BlockPack.NotFound().WithOrigin(err)
	}

	return &resDto, nil
}

func (s *BlockPackService) GetMyBlockPacksByParentSubShelfId(
	ctx context.Context, reqDto *dtos.GetMyBlockPacksByParentSubShelfIdReqDto,
) (*dtos.GetMyBlockPacksByParentSubShelfIdResDto, *exceptions.Exception) {
	if err := validation.Validator.Struct(reqDto); err != nil {
		return nil, exceptions.BlockPack.InvalidDto().WithOrigin(err)
	}

	db := s.db.WithContext(ctx)

	allowedPermissions, exception := contexts.GetAllowedPermissions(ctx)
	if exception != nil {
		return nil, exception
	}

	onlyDeleted := types.Ternary_Neutral
	if reqDto.Param.AreDeleted != nil {
		if *reqDto.Param.AreDeleted {
			onlyDeleted = types.Ternary_Positive
		} else {
			onlyDeleted = types.Ternary_Negative
		}
	}

	resDto := dtos.GetMyBlockPacksByParentSubShelfIdResDto{}
	result := db.Model(&schemas.BlockPack{}).
		Select(`
			"BlockPackTable".*,
			COALESCE(ydoc.last_update_sequence, 0) AS last_update_sequence,
			COALESCE(ydoc.compacted_until_sequence, 0) AS compacted_until_sequence,
			COALESCE(ydoc.projected_until_sequence, -1) AS projected_until_sequence,
			COALESCE(ydoc.projected_until_sequence, -1) >= COALESCE(ydoc.last_update_sequence, 0) AS is_projection_current
		`).
		Joins(`LEFT JOIN "BlockPackYjsDocumentTable" ydoc ON ydoc.block_pack_id = "BlockPackTable".id AND ydoc.deleted_at IS NULL`).
		Joins(`LEFT JOIN "SubShelfTable" ss ON "BlockPackTable".parent_sub_shelf_id = ss.id`).
		Joins(`LEFT JOIN "UsersToShelvesTable" uts ON ss.root_shelf_id = uts.root_shelf_id`).
		Where("ss.id = ? AND uts.user_id = ? AND uts.permission IN ?",
			reqDto.Param.ParentSubShelfId,
			reqDto.ContextFields.UserId,
			allowedPermissions,
		).Scopes(scopes.NewBlockPackScope().FilterOnlyDeleted(onlyDeleted)).
		Order("name ASC").
		Limit(int(constants.MaxBlockPackOfSubShelf)).
		Scan(&resDto)
	if err := result.Error; err != nil {
		return nil, exceptions.BlockPack.NotFound().WithOrigin(err)
	}

	return &resDto, nil
}

func (s *BlockPackService) GetAllMyBlockPacksByRootShelfId(
	ctx context.Context, reqDto *dtos.GetAllMyBlockPacksByRootShelfIdReqDto,
) (*dtos.GetAllMyBlockPacksByRootShelfIdResDto, *exceptions.Exception) {
	if err := validation.Validator.Struct(reqDto); err != nil {
		return nil, exceptions.BlockPack.InvalidDto().WithOrigin(err)
	}

	db := s.db.WithContext(ctx)

	allowedPermissions, exception := contexts.GetAllowedPermissions(ctx)
	if exception != nil {
		return nil, exception
	}

	onlyDeleted := types.Ternary_Neutral
	if reqDto.Param.AreDeleted != nil {
		if *reqDto.Param.AreDeleted {
			onlyDeleted = types.Ternary_Positive
		} else {
			onlyDeleted = types.Ternary_Negative
		}
	}

	resDto := dtos.GetAllMyBlockPacksByRootShelfIdResDto{}
	result := db.Model(&schemas.BlockPack{}).
		Select(`
			"BlockPackTable".*,
			COALESCE(ydoc.last_update_sequence, 0) AS last_update_sequence,
			COALESCE(ydoc.compacted_until_sequence, 0) AS compacted_until_sequence,
			COALESCE(ydoc.projected_until_sequence, -1) AS projected_until_sequence,
			COALESCE(ydoc.projected_until_sequence, -1) >= COALESCE(ydoc.last_update_sequence, 0) AS is_projection_current
		`).
		Joins(`LEFT JOIN "BlockPackYjsDocumentTable" ydoc ON ydoc.block_pack_id = "BlockPackTable".id AND ydoc.deleted_at IS NULL`).
		Joins(`LEFT JOIN "SubShelfTable" ss ON "BlockPackTable".parent_sub_shelf_id = ss.id`).
		Joins(`LEFT JOIN "UsersToShelvesTable" uts ON ss.root_shelf_id = uts.root_shelf_id`).
		Where("ss.root_shelf_id = ? AND uts.user_id = ? AND uts.permission IN ?",
			reqDto.Param.RootShelfId, reqDto.ContextFields.UserId, allowedPermissions,
		).Scopes(scopes.NewBlockPackScope().FilterOnlyDeleted(onlyDeleted)).
		Limit(int(constants.MaxBlockPackOfRootShelf)).
		Order("name ASC").
		Scan(&resDto)
	if err := result.Error; err != nil {
		return nil, exceptions.BlockPack.NotFound().WithOrigin(err)
	}

	return &resDto, nil
}

func (s *BlockPackService) CreateBlockPack(
	ctx context.Context, reqDto *dtos.CreateBlockPackReqDto,
) (*dtos.CreateBlockPackResDto, *exceptions.Exception) {
	if err := validation.Validator.Struct(reqDto); err != nil {
		return nil, exceptions.BlockPack.InvalidDto().WithOrigin(err)
	}

	allowedPermissions, exception := contexts.GetAllowedPermissions(ctx)
	if exception != nil {
		return nil, exception
	}

	tx := s.db.WithContext(ctx).Begin()

	newBlockPackId, exception := s.blockPackRepository.CreateOneBySubShelfId(
		reqDto.Body.ParentSubShelfId,
		reqDto.ContextFields.UserId,
		inputs.CreateBlockPackInput{
			Id:                  reqDto.Body.Id,
			Name:                reqDto.Body.Name,
			Icon:                reqDto.Body.Icon,
			HeaderBackgroundURL: reqDto.Body.HeaderBackgroundURL,
		},
		options.WithTransactionDB(tx),
		options.WithAllowedPermissions(allowedPermissions),
		options.WithLockingStrength(options.LockingStrengthNoKeyUpdate),
	)
	if exception != nil {
		tx.Rollback()
		return nil, exception
	}

	document := schemas.BlockPackYjsDocument{BlockPackId: *newBlockPackId}
	if err := tx.Create(&document).Error; err != nil {
		tx.Rollback()
		return nil, exceptions.BlockPack.FailedToCreate().WithOrigin(err)
	}

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return nil, exceptions.BlockPack.FailedToCommitTransaction().WithOrigin(err)
	}

	return &dtos.CreateBlockPackResDto{
		Id:        *newBlockPackId,
		CreatedAt: time.Now(),
	}, nil
}

func (s *BlockPackService) CreateBlockPacks(
	ctx context.Context, reqDto *dtos.CreateBlockPacksReqDto,
) (*dtos.CreateBlockPacksResDto, *exceptions.Exception) {
	if err := validation.Validator.Struct(reqDto); err != nil {
		return nil, exceptions.BlockPack.InvalidDto().WithOrigin(err)
	}

	allowedPermissions, exception := contexts.GetAllowedPermissions(ctx)
	if exception != nil {
		return nil, exception
	}

	tx := s.db.WithContext(ctx).Begin()

	input := make([]inputs.CreateBlockPackBySubShelfIdInput, len(reqDto.Body.CreatedBlockPacks))
	for index, createdBlockPack := range reqDto.Body.CreatedBlockPacks {
		input[index] = inputs.CreateBlockPackBySubShelfIdInput{
			Id:                  createdBlockPack.Id,
			ParentSubShelfId:    createdBlockPack.ParentSubShelfId,
			Name:                createdBlockPack.Name,
			Icon:                createdBlockPack.Icon,
			HeaderBackgroundURL: createdBlockPack.HeaderBackgroundURL,
		}
	}
	newBlockPackIds, exception := s.blockPackRepository.CreateManyBySubShelfIds(
		reqDto.ContextFields.UserId,
		input,
		options.WithTransactionDB(tx),
		options.WithAllowedPermissions(allowedPermissions),
		options.WithLockingStrength(options.LockingStrengthNoKeyUpdate),
	)
	if exception != nil {
		tx.Rollback()
		return nil, exception
	}

	documents := make([]schemas.BlockPackYjsDocument, len(newBlockPackIds))
	for index, newBlockPackId := range newBlockPackIds {
		documents[index] = schemas.BlockPackYjsDocument{BlockPackId: newBlockPackId}
	}
	if err := tx.CreateInBatches(&documents, constants.MaxBatchCreateBlockSize).Error; err != nil {
		tx.Rollback()

		return nil, exceptions.BlockPack.FailedToCreate().WithOrigin(err)
	}

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()

		return nil, exceptions.BlockPack.FailedToCommitTransaction().WithOrigin(err)
	}

	return &dtos.CreateBlockPacksResDto{
		Ids:       newBlockPackIds,
		CreatedAt: time.Now(),
	}, nil
}

func (s *BlockPackService) UpdateMyBlockPackById(
	ctx context.Context, reqDto *dtos.UpdateMyBlockPackByIdReqDto,
) (*dtos.UpdateMyBlockPackByIdResDto, *exceptions.Exception) {
	if err := validation.Validator.Struct(reqDto); err != nil {
		return nil, exceptions.BlockPack.InvalidDto().WithOrigin(err)
	}

	db := s.db.WithContext(ctx)
	allowedPermissions, exception := contexts.GetAllowedPermissions(ctx)
	if exception != nil {
		return nil, exception
	}

	blockPack, exception := s.blockPackRepository.UpdateOneById(
		reqDto.Body.BlockPackId,
		reqDto.ContextFields.UserId,
		inputs.PartialUpdateBlockPackInput{
			Values: inputs.UpdateBlockPackInput{
				Name:                reqDto.Body.Values.Name,
				Icon:                reqDto.Body.Values.Icon,
				HeaderBackgroundURL: reqDto.Body.Values.HeaderBackgroundURL,
			},
			SetNull: reqDto.Body.SetNull,
		},
		options.WithDB(db),
		options.WithAllowedPermissions(allowedPermissions),
	)
	if exception != nil {
		return nil, exception
	}

	return &dtos.UpdateMyBlockPackByIdResDto{
		UpdatedAt: blockPack.UpdatedAt,
	}, nil
}

func (s *BlockPackService) UpdateMyBlockPacksByIds(
	ctx context.Context, reqDto *dtos.UpdateMyBlockPacksByIdsReqDto,
) (*dtos.UpdateMyBlockPacksByIdsResDto, *exceptions.Exception) {
	if err := validation.Validator.Struct(reqDto); err != nil {
		return nil, exceptions.BlockPack.InvalidDto().WithOrigin(err)
	}

	db := s.db.WithContext(ctx)
	allowedPermissions, exception := contexts.GetAllowedPermissions(ctx)
	if exception != nil {
		return nil, exception
	}

	input := make([]inputs.UpdateBlockPackByIdInput, len(reqDto.Body.UpdatedBlockPacks))
	for index, updatedBlockPack := range reqDto.Body.UpdatedBlockPacks {
		input[index] = inputs.UpdateBlockPackByIdInput{
			Id: updatedBlockPack.BlockPackId,
			PartialUpdateInput: inputs.PartialUpdateInput[inputs.UpdateBlockPackInput]{
				Values: inputs.UpdateBlockPackInput{
					Name:                updatedBlockPack.Values.Name,
					Icon:                updatedBlockPack.Values.Icon,
					HeaderBackgroundURL: updatedBlockPack.Values.HeaderBackgroundURL,
				},
			},
		}
	}
	exception = s.blockPackRepository.UpdateManyByIds(
		reqDto.ContextFields.UserId,
		input,
		options.WithDB(db),
		options.WithAllowedPermissions(allowedPermissions),
	)
	if exception != nil {
		return nil, exception
	}

	return &dtos.UpdateMyBlockPacksByIdsResDto{
		UpdatedAt: time.Now(),
	}, nil
}

func (s *BlockPackService) MoveMyBlockPackByParentSubShelfId(
	ctx context.Context, reqDto *dtos.MoveMyBlockPackByParentSubShelfIdReqDto,
) (*dtos.MoveMyBlockPackByParentSubShelfIdResDto, *exceptions.Exception) {
	if err := validation.Validator.Struct(reqDto); err != nil {
		return nil, exceptions.BlockPack.InvalidDto().WithOrigin(err)
	}

	db := s.db.WithContext(ctx)
	allowedPermissions, exception := contexts.GetAllowedPermissions(ctx)
	if exception != nil {
		return nil, exception
	}

	_, exception = s.blockPackRepository.UpdateOneById(
		reqDto.Body.BlockPackId,
		reqDto.ContextFields.UserId,
		inputs.PartialUpdateBlockPackInput{
			Values: inputs.UpdateBlockPackInput{
				ParentSubShelfId: &reqDto.Body.DestinationParentSubShelfId,
			},
			SetNull: nil,
		},
		options.WithDB(db),
		options.WithAllowedPermissions(allowedPermissions),
	)
	if exception != nil {
		return nil, exception
	}

	return &dtos.MoveMyBlockPackByParentSubShelfIdResDto{
		UpdatedAt: time.Now(),
	}, nil
}

func (s *BlockPackService) MoveMyBlockPacksByParentSubShelfId(
	ctx context.Context, reqDto *dtos.MoveMyBlockPacksByParentSubShelfIdReqDto,
) (*dtos.MoveMyBlockPacksByParentSubShelfIdResDto, *exceptions.Exception) {
	if err := validation.Validator.Struct(reqDto); err != nil {
		return nil, exceptions.BlockPack.InvalidDto().WithOrigin(err)
	}

	db := s.db.WithContext(ctx)
	allowedPermissions, exception := contexts.GetAllowedPermissions(ctx)
	if exception != nil {
		return nil, exception
	}

	input := make([]inputs.UpdateBlockPackByIdInput, len(reqDto.Body.BlockPackIds))
	for index, blockPackId := range reqDto.Body.BlockPackIds {
		input[index] = inputs.UpdateBlockPackByIdInput{
			Id: blockPackId,
			PartialUpdateInput: inputs.PartialUpdateInput[inputs.UpdateBlockPackInput]{
				Values: inputs.UpdateBlockPackInput{
					ParentSubShelfId: &reqDto.Body.DestinationParentSubShelfId,
				},
			},
		}
	}
	exception = s.blockPackRepository.UpdateManyByIds(
		reqDto.ContextFields.UserId,
		input,
		options.WithDB(db),
		options.WithAllowedPermissions(allowedPermissions),
	)
	if exception != nil {
		return nil, exception
	}

	return &dtos.MoveMyBlockPacksByParentSubShelfIdResDto{
		UpdatedAt: time.Now(),
	}, nil
}

func (s *BlockPackService) MoveMyBlockPacksByParentSubShelfIds(
	ctx context.Context, reqDto *dtos.MoveMyBlockPacksByParentSubShelfIdsReqDto,
) (*dtos.MoveMyBlockPacksByParentSubShelfIdsResDto, *exceptions.Exception) {
	if err := validation.Validator.Struct(reqDto); err != nil {
		return nil, exceptions.BlockPack.InvalidDto().WithOrigin(err)
	}

	db := s.db.WithContext(ctx)
	allowedPermissions, exception := contexts.GetAllowedPermissions(ctx)
	if exception != nil {
		return nil, exception
	}

	input := make([]inputs.UpdateBlockPackByIdInput, 0)
	for _, movedBlockPack := range reqDto.Body.MovedBlockPacks {
		for _, blockPackId := range movedBlockPack.BlockPackIds {
			input = append(input, inputs.UpdateBlockPackByIdInput{
				Id: blockPackId,
				PartialUpdateInput: inputs.PartialUpdateInput[inputs.UpdateBlockPackInput]{
					Values: inputs.UpdateBlockPackInput{
						ParentSubShelfId: &movedBlockPack.DestinationParentSubShelfId,
					},
				},
			})
		}
	}

	if exception = s.blockPackRepository.UpdateManyByIds(
		reqDto.ContextFields.UserId,
		input,
		options.WithDB(db),
		options.WithAllowedPermissions(allowedPermissions),
	); exception != nil {
		return nil, exception
	}

	return &dtos.MoveMyBlockPacksByParentSubShelfIdsResDto{
		UpdatedAt: time.Now(),
	}, nil
}

func (s *BlockPackService) RestoreMyBlockPackById(
	ctx context.Context, reqDto *dtos.RestoreMyBlockPackByIdReqDto,
) (*dtos.RestoreMyBlockPackByIdResDto, *exceptions.Exception) {
	if err := validation.Validator.Struct(reqDto); err != nil {
		return nil, exceptions.BlockPack.InvalidDto().WithOrigin(err)
	}

	db := s.db.WithContext(ctx)
	allowedPermissions, exception := contexts.GetAllowedPermissions(ctx)
	if exception != nil {
		return nil, exception
	}

	restoredBlockPack, exception := s.blockPackRepository.RestoreSoftDeletedOneById(
		reqDto.Body.BlockPackId,
		reqDto.ContextFields.UserId,
		options.WithDB(db),
		options.WithAllowedPermissions(allowedPermissions),
	)
	if exception != nil {
		return nil, exception
	}

	return &dtos.RestoreMyBlockPackByIdResDto{
		Id:                  restoredBlockPack.Id,
		ParentSubShelfId:    restoredBlockPack.ParentSubShelfId,
		Name:                restoredBlockPack.Name,
		Icon:                restoredBlockPack.Icon,
		HeaderBackgroundURL: restoredBlockPack.HeaderBackgroundURL,
		BlockCount:          restoredBlockPack.BlockCount,
		DeletedAt:           restoredBlockPack.DeletedAt,
		UpdatedAt:           restoredBlockPack.UpdatedAt,
		CreatedAt:           restoredBlockPack.CreatedAt,
	}, nil
}

func (s *BlockPackService) RestoreMyBlockPacksByIds(
	ctx context.Context, reqDto *dtos.RestoreMyBlockPacksByIdsReqDto,
) (*dtos.RestoreMyBlockPacksByIdsResDto, *exceptions.Exception) {
	if err := validation.Validator.Struct(reqDto); err != nil {
		return nil, exceptions.BlockPack.InvalidDto().WithOrigin(err)
	}

	db := s.db.WithContext(ctx)
	allowedPermissions, exception := contexts.GetAllowedPermissions(ctx)
	if exception != nil {
		return nil, exception
	}

	restoredBlockPacks, exception := s.blockPackRepository.RestoreSoftDeletedManyByIds(
		reqDto.Body.BlockPackIds,
		reqDto.ContextFields.UserId,
		options.WithDB(db),
		options.WithAllowedPermissions(allowedPermissions),
	)
	if exception != nil {
		return nil, exception
	}

	resDto := dtos.RestoreMyBlockPacksByIdsResDto{}
	for _, restoredBlockPack := range restoredBlockPacks {
		resDto = append(resDto, dtos.RestoreMyBlockPackByIdResDto{
			Id:                  restoredBlockPack.Id,
			ParentSubShelfId:    restoredBlockPack.ParentSubShelfId,
			Name:                restoredBlockPack.Name,
			Icon:                restoredBlockPack.Icon,
			HeaderBackgroundURL: restoredBlockPack.HeaderBackgroundURL,
			BlockCount:          restoredBlockPack.BlockCount,
			DeletedAt:           restoredBlockPack.DeletedAt,
			UpdatedAt:           restoredBlockPack.UpdatedAt,
			CreatedAt:           restoredBlockPack.CreatedAt,
		})
	}

	return &resDto, nil
}

func (s *BlockPackService) DeleteMyBlockPackById(
	ctx context.Context, reqDto *dtos.DeleteMyBlockPackByIdReqDto,
) (*dtos.DeleteMyBlockPackByIdResDto, *exceptions.Exception) {
	if err := validation.Validator.Struct(reqDto); err != nil {
		return nil, exceptions.BlockPack.InvalidDto().WithOrigin(err)
	}

	db := s.db.WithContext(ctx)
	allowedPermissions, exception := contexts.GetAllowedPermissions(ctx)
	if exception != nil {
		return nil, exception
	}

	if exception = s.blockPackRepository.SoftDeleteOneById(
		reqDto.Body.BlockPackId,
		reqDto.ContextFields.UserId,
		options.WithDB(db),
		options.WithAllowedPermissions(allowedPermissions),
	); exception != nil {
		return nil, exception
	}
	if err := s.realtimeLeaseStore.PublishBlockPackChannelRevocation(uuid.Nil, []uuid.UUID{reqDto.Body.BlockPackId}); err != nil {
		logs.NotezyLogger.Error(ctx, err, "Failed to revoke realtime BlockPack channels")
	}

	return &dtos.DeleteMyBlockPackByIdResDto{
		DeletedAt: time.Now(),
	}, nil
}

func (s *BlockPackService) DeleteMyBlockPacksByIds(
	ctx context.Context, reqDto *dtos.DeleteMyBlockPacksByIdsReqDto,
) (*dtos.DeleteMyBlockPacksByIdsResDto, *exceptions.Exception) {
	if err := validation.Validator.Struct(reqDto); err != nil {
		return nil, exceptions.BlockPack.InvalidDto().WithOrigin(err)
	}

	db := s.db.WithContext(ctx)
	allowedPermissions, exception := contexts.GetAllowedPermissions(ctx)
	if exception != nil {
		return nil, exception
	}

	if exception = s.blockPackRepository.SoftDeleteManyByIds(
		reqDto.Body.BlockPackIds,
		reqDto.ContextFields.UserId,
		options.WithDB(db),
		options.WithAllowedPermissions(allowedPermissions),
	); exception != nil {
		return nil, exception
	}
	if err := s.realtimeLeaseStore.PublishBlockPackChannelRevocation(uuid.Nil, reqDto.Body.BlockPackIds); err != nil {
		logs.NotezyLogger.Error(ctx, err, "Failed to revoke realtime BlockPack channels")
	}

	return &dtos.DeleteMyBlockPacksByIdsResDto{
		DeletedAt: time.Now(),
	}, nil
}

/* ============================== Service Methods for GraphQL BlockPack ============================== */

func (s *BlockPackService) SearchPrivateBlockPacks(
	ctx context.Context, userId uuid.UUID, gqlInput gqlmodels.SearchBlockPackInput,
) (*gqlmodels.SearchBlockPackConnection, *exceptions.Exception) {
	startTime := time.Now()
	db := s.db.WithContext(ctx)

	allowedPermissions, exception := contexts.GetAllowedPermissions(ctx)
	if exception != nil {
		return nil, exception
	}

	onlyDeleted := types.Ternary_Negative
	if gqlInput.IsDeletedAt != nil && *gqlInput.IsDeletedAt {
		onlyDeleted = types.Ternary_Positive
	}

	query := db.Model(&schemas.BlockPack{}).
		Select(`"BlockPackTable".*`).
		Joins(`INNER JOIN "SubShelfTable" ss ON ss.id = "BlockPackTable".parent_sub_shelf_id`).
		Joins(`INNER JOIN "UsersToShelvesTable" uts ON uts.root_shelf_id = ss.root_shelf_id`).
		Where("uts.user_id = ? AND uts.permission IN ?", userId, allowedPermissions).
		Scopes(s.blockPackScope.FilterOnlyDeleted(onlyDeleted))

	if gqlInput.ParentSubShelfID != nil {
		query = query.Where(`"BlockPackTable".parent_sub_shelf_id = ?`, *gqlInput.ParentSubShelfID)
	}

	if gqlInput.RootShelfID != nil {
		query = query.Where("ss.root_shelf_id = ?", *gqlInput.RootShelfID)
	}

	if len(strings.ReplaceAll(gqlInput.Query, " ", "")) > 0 {
		query = query.Where(
			`"BlockPackTable".name ILIKE ?`,
			"%"+gqlInput.Query+"%",
		)
	}
	if gqlInput.After != nil && len(strings.ReplaceAll(*gqlInput.After, " ", "")) > 0 {
		searchCursor, err := searchcursor.Decode[gqlmodels.SearchBlockPackCursorFields](*gqlInput.After)
		if err != nil {
			return nil, exceptions.Search.FailedToDecode().WithOrigin(err)
		}

		query = query.Where(`"BlockPackTable".id > ?`, searchCursor.Fields.ID)
	}

	if gqlInput.SortBy != nil && gqlInput.SortOrder != nil {
		cending := gqlmodels.SearchSortOrderAsc.String()
		if *gqlInput.SortOrder == gqlmodels.SearchSortOrderDesc {
			cending = gqlmodels.SearchSortOrderDesc.String()
		}

		switch *gqlInput.SortBy {
		case gqlmodels.SearchBlockPackSortByName:
			query = query.Order(`"BlockPackTable".name ` + cending).
				Order(`"BlockPackTable".updated_at ` + cending).
				Order(`"BlockPackTable".created_at ` + cending)
		case gqlmodels.SearchBlockPackSortByBlockCount:
			query = query.Order(`"BlockPackTable".block_count ` + cending).
				Order(`"BlockPackTable".name ` + cending).
				Order(`"BlockPackTable".updated_at ` + cending).
				Order(`"BlockPackTable".created_at ` + cending)
		case gqlmodels.SearchBlockPackSortByLastUpdate:
			query = query.Order(`"BlockPackTable".updated_at ` + cending).
				Order(`"BlockPackTable".name ` + cending).
				Order(`"BlockPackTable".created_at ` + cending)
		case gqlmodels.SearchBlockPackSortByCreatedAt:
			query = query.Order(`"BlockPackTable".created_at ` + cending).
				Order(`"BlockPackTable".name ` + cending).
				Order(`"BlockPackTable".updated_at ` + cending)
		default:
			query = query.Order(`"BlockPackTable".name ` + cending).
				Order(`"BlockPackTable".updated_at ` + cending).
				Order(`"BlockPackTable".created_at ` + cending)
		}
	}

	limit := constants.DefaultSearchLimit
	if gqlInput.First != nil && *gqlInput.First > 0 {
		limit = int(*gqlInput.First)
	}
	limit = min(limit, constants.MaxSearchLimit)
	query = query.Limit(limit + 1)

	var blockPacks []schemas.BlockPack
	if err := query.Scopes(s.blockPackScope.IncludePreloads(
		[]schemas.BlockPackRelation{
			schemas.BlockPackRelation_Blocks,
		},
	)).Find(&blockPacks).Error; err != nil {
		return nil, exceptions.BlockPack.NotFound().WithOrigin(err)
	}

	hasNextPage := len(blockPacks) > limit
	searchEdges := make([]*gqlmodels.SearchBlockPackEdge, len(blockPacks))

	for index, blockPack := range blockPacks {
		searchCursor := searchcursor.SearchCursor[gqlmodels.SearchBlockPackCursorFields]{
			Fields: gqlmodels.SearchBlockPackCursorFields{
				ID: blockPack.Id,
			},
		}
		encodedSearchCursor, err := searchCursor.Encode()
		if err != nil {
			return nil, exceptions.Search.FailedToEncode().WithOrigin(err)
		}
		if encodedSearchCursor == nil {
			return nil, exceptions.Search.FailedToUnmarshalSearchCursor()
		}

		searchEdges[index] = &gqlmodels.SearchBlockPackEdge{
			EncodedSearchCursor: *encodedSearchCursor,
			Node:                blockPack.ToPrivateBlockPack(),
		}
	}

	searchPageInfo := &gqlmodels.SearchPageInfo{
		HasNextPage:     hasNextPage,
		HasPreviousPage: gqlInput.After != nil && len(strings.ReplaceAll(*gqlInput.After, " ", "")) > 0,
	}

	if len(searchEdges) > 0 {
		searchPageInfo.StartEncodedSearchCursor = &searchEdges[0].EncodedSearchCursor
		searchPageInfo.EndEncodedSearchCursor = &searchEdges[len(searchEdges)-1].EncodedSearchCursor
	}

	searchTime := float64(time.Since(startTime).Nanoseconds()) / 1e6
	if hasNextPage {
		searchEdges = searchEdges[:limit]
	}

	return &gqlmodels.SearchBlockPackConnection{
		SearchEdges:    searchEdges,
		SearchPageInfo: searchPageInfo,
		TotalCount:     int32(len(searchEdges)),
		SearchTime:     searchTime,
	}, nil
}
