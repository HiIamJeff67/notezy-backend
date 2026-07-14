package services

import (
	"context"
	"time"

	pg "github.com/lib/pq"
	"gorm.io/gorm"

	dtos "github.com/HiIamJeff67/notezy-backend/app/dtos"
	exceptions "github.com/HiIamJeff67/notezy-backend/app/exceptions"
	inputs "github.com/HiIamJeff67/notezy-backend/app/models/inputs"
	repositories "github.com/HiIamJeff67/notezy-backend/app/models/repositories"
	schemas "github.com/HiIamJeff67/notezy-backend/app/models/schemas"
	enums "github.com/HiIamJeff67/notezy-backend/app/models/schemas/enums"
	scopes "github.com/HiIamJeff67/notezy-backend/app/models/scopes"
	blockpacksql "github.com/HiIamJeff67/notezy-backend/app/models/sqls/block_pack"
	options "github.com/HiIamJeff67/notezy-backend/app/options"
	validation "github.com/HiIamJeff67/notezy-backend/app/validation"
	constants "github.com/HiIamJeff67/notezy-backend/shared/constants"
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
	MoveMyBlockPackById(ctx context.Context, reqDto *dtos.MoveMyBlockPackByIdReqDto) (*dtos.MoveMyBlockPackByIdResDto, *exceptions.Exception)
	MoveMyBlockPacksByParentSubShelfId(ctx context.Context, reqDto *dtos.MoveMyBlockPacksByParentSubShelfIdReqDto) (*dtos.MoveMyBlockPacksByParentSubShelfIdResDto, *exceptions.Exception)
	MoveMyBlockPacksByParentSubShelfIds(ctx context.Context, reqDto *dtos.MoveMyBlockPacksByParentSubShelfIdsReqDto) (*dtos.MoveMyBlockPacksByParentSubShelfIdsResDto, *exceptions.Exception)
	RestoreMyBlockPackById(ctx context.Context, reqDto *dtos.RestoreMyBlockPackByIdReqDto) (*dtos.RestoreMyBlockPackByIdResDto, *exceptions.Exception)
	RestoreMyBlockPacksByIds(ctx context.Context, reqDto *dtos.RestoreMyBlockPacksByIdsReqDto) (*dtos.RestoreMyBlockPacksByIdsResDto, *exceptions.Exception)
	DeleteMyBlockPackById(ctx context.Context, reqDto *dtos.DeleteMyBlockPackByIdReqDto) (*dtos.DeleteMyBlockPackByIdResDto, *exceptions.Exception)
	DeleteMyBlockPacksByIds(ctx context.Context, reqDto *dtos.DeleteMyBlockPacksByIdsReqDto) (*dtos.DeleteMyBlockPacksByIdsResDto, *exceptions.Exception)
}

type BlockPackService struct {
	db                  *gorm.DB
	subShelfRepository  repositories.SubShelfRepositoryInterface
	blockPackRepository repositories.BlockPackRepositoryInterface
}

func NewBlockPackService(
	db *gorm.DB,
	subShelfRepository repositories.SubShelfRepositoryInterface,
	blockPackRepository repositories.BlockPackRepositoryInterface,
) BlockPackServiceInterface {
	return &BlockPackService{
		db:                  db,
		subShelfRepository:  subShelfRepository,
		blockPackRepository: blockPackRepository,
	}
}

/* ============================== Service Methods for BlockPack ============================== */

func (s *BlockPackService) GetMyBlockPackById(
	ctx context.Context, reqDto *dtos.GetMyBlockPackByIdReqDto,
) (*dtos.GetMyBlockPackByIdResDto, *exceptions.Exception) {
	if err := validation.Validator.Struct(reqDto); err != nil {
		return nil, exceptions.BlockPack.InvalidDto().WithOrigin(err)
	}

	db := s.db.WithContext(ctx)

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
		[]enums.AccessControlPermission{
			enums.AccessControlPermission_Owner,
			enums.AccessControlPermission_Admin,
			enums.AccessControlPermission_Write,
			enums.AccessControlPermission_Read,
		},
		options.WithDB(db),
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

	allowedPermissions := []enums.AccessControlPermission{
		enums.AccessControlPermission_Owner,
		enums.AccessControlPermission_Admin,
		enums.AccessControlPermission_Write,
		enums.AccessControlPermission_Read,
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

	onlyDeleted := types.Ternary_Neutral
	if reqDto.Param.AreDeleted != nil {
		if *reqDto.Param.AreDeleted {
			onlyDeleted = types.Ternary_Positive
		} else {
			onlyDeleted = types.Ternary_Negative
		}
	}

	allowedPermissions := []enums.AccessControlPermission{
		enums.AccessControlPermission_Owner,
		enums.AccessControlPermission_Admin,
		enums.AccessControlPermission_Write,
		enums.AccessControlPermission_Read,
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
	onlyDeleted := types.Ternary_Neutral
	if reqDto.Param.AreDeleted != nil {
		if *reqDto.Param.AreDeleted {
			onlyDeleted = types.Ternary_Positive
		} else {
			onlyDeleted = types.Ternary_Negative
		}
	}

	allowedPermissions := []enums.AccessControlPermission{
		enums.AccessControlPermission_Owner,
		enums.AccessControlPermission_Admin,
		enums.AccessControlPermission_Write,
		enums.AccessControlPermission_Read,
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
	exception := s.blockPackRepository.UpdateManyByIds(
		reqDto.ContextFields.UserId,
		input,
		options.WithDB(db),
	)
	if exception != nil {
		return nil, exception
	}

	return &dtos.UpdateMyBlockPacksByIdsResDto{
		UpdatedAt: time.Now(),
	}, nil
}

func (s *BlockPackService) MoveMyBlockPackById(
	ctx context.Context, reqDto *dtos.MoveMyBlockPackByIdReqDto,
) (*dtos.MoveMyBlockPackByIdResDto, *exceptions.Exception) {
	if err := validation.Validator.Struct(reqDto); err != nil {
		return nil, exceptions.BlockPack.InvalidDto().WithOrigin(err)
	}

	db := s.db.WithContext(ctx)

	_, exception := s.blockPackRepository.UpdateOneById(
		reqDto.Body.BlockPackId,
		reqDto.ContextFields.UserId,
		inputs.PartialUpdateBlockPackInput{
			Values: inputs.UpdateBlockPackInput{
				ParentSubShelfId: &reqDto.Body.DestinationParentSubShelfId,
			},
			SetNull: nil,
		},
		options.WithDB(db),
	)
	if exception != nil {
		return nil, exception
	}

	return &dtos.MoveMyBlockPackByIdResDto{
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
	exception := s.blockPackRepository.UpdateManyByIds(
		reqDto.ContextFields.UserId,
		input,
		options.WithDB(db),
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

	if exception := s.blockPackRepository.UpdateManyByIds(
		reqDto.ContextFields.UserId,
		input,
		options.WithDB(db),
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

	restoredBlockPack, exception := s.blockPackRepository.RestoreSoftDeletedOneById(
		reqDto.Body.BlockPackId,
		reqDto.ContextFields.UserId,
		options.WithDB(db),
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

	restoredBlockPacks, exception := s.blockPackRepository.RestoreSoftDeletedManyByIds(
		reqDto.Body.BlockPackIds,
		reqDto.ContextFields.UserId,
		options.WithDB(db),
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

	if exception := s.blockPackRepository.SoftDeleteOneById(
		reqDto.Body.BlockPackId,
		reqDto.ContextFields.UserId,
		options.WithDB(db),
	); exception != nil {
		return nil, exception
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

	if exception := s.blockPackRepository.SoftDeleteManyByIds(
		reqDto.Body.BlockPackIds,
		reqDto.ContextFields.UserId,
		options.WithDB(db),
	); exception != nil {
		return nil, exception
	}

	return &dtos.DeleteMyBlockPacksByIdsResDto{
		DeletedAt: time.Now(),
	}, nil
}
