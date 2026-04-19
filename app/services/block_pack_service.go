package services

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	pg "github.com/lib/pq"
	"gorm.io/gorm"

	dtos "notezy-backend/app/dtos"
	exceptions "notezy-backend/app/exceptions"
	inputs "notezy-backend/app/models/inputs"
	repositories "notezy-backend/app/models/repositories"
	schemas "notezy-backend/app/models/schemas"
	enums "notezy-backend/app/models/schemas/enums"
	blockpacksql "notezy-backend/app/models/sqls/block_pack"
	options "notezy-backend/app/options"
	validation "notezy-backend/app/validation"
	constants "notezy-backend/shared/constants"
	types "notezy-backend/shared/types"
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
	MoveMyBlockPacksByIds(ctx context.Context, reqDto *dtos.MoveMyBlockPacksByIdsReqDto) (*dtos.MoveMyBlockPacksByIdsResDto, *exceptions.Exception)
	BatchMoveMyBlockPacksByIds(ctx context.Context, reqDto *dtos.BatchMoveMyBlockPacksByIdsReqDto) (*dtos.BatchMoveMyBlockPacksByIdsResDto, *exceptions.Exception)
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

/* ============================== Service Methods for Block Pack ============================== */

func (s *BlockPackService) GetMyBlockPackById(
	ctx context.Context, reqDto *dtos.GetMyBlockPackByIdReqDto,
) (*dtos.GetMyBlockPackByIdResDto, *exceptions.Exception) {
	if err := validation.Validator.Struct(reqDto); err != nil {
		return nil, exceptions.BlockPack.InvalidDto().WithOrigin(err)
	}

	db := s.db.WithContext(ctx)

	blockPack, exception := s.blockPackRepository.GetOneById(
		reqDto.Param.BlockPackId,
		reqDto.ContextFields.UserId,
		options.WithDB(db),
		options.WithOnlyDeleted(types.Ternary_Negative),
	)
	if exception != nil {
		return nil, exception
	}

	return &dtos.GetMyBlockPackByIdResDto{
		Id:                  blockPack.Id,
		ParentSubShelfId:    blockPack.ParentSubShelfId,
		Name:                blockPack.Name,
		Icon:                blockPack.Icon,
		HeaderBackgroundURL: blockPack.HeaderBackgroundURL,
		BlockCount:          blockPack.BlockCount,
		DeletedAt:           blockPack.DeletedAt,
		UpdatedAt:           blockPack.UpdatedAt,
		CreatedAt:           blockPack.CreatedAt,
	}, nil
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
	onlyDeleted := types.Ternary_Negative
	resDto := dtos.GetMyBlockPackAndItsParentByIdResDto{}

	err := db.Raw(blockpacksql.GetMyBlockPackAndItsParentByIdSQL,
		reqDto.Param.BlockPackId, reqDto.ContextFields.UserId, pg.Array(allowedPermissions), onlyDeleted,
	).Row().
		Scan(&resDto.Id,
			&resDto.Name,
			&resDto.Icon,
			&resDto.HeaderBackgroundURL,
			&resDto.BlockCount,
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

	allowedPermissions := []enums.AccessControlPermission{
		enums.AccessControlPermission_Owner,
		enums.AccessControlPermission_Admin,
		enums.AccessControlPermission_Write,
		enums.AccessControlPermission_Read,
	}

	resDto := dtos.GetMyBlockPacksByParentSubShelfIdResDto{}

	result := db.Model(&schemas.BlockPack{}).
		Joins("LEFT JOIN \"SubShelfTable\" ss ON \"BlockPackTable\".parent_sub_shelf_id = ss.id").
		Joins("LEFT JOIN \"UsersToShelvesTable\" uts ON ss.root_shelf_id = uts.root_shelf_id").
		Where("ss.id = ? AND uts.user_id = ? AND uts.permission IN ?",
			reqDto.Param.ParentSubShelfId,
			reqDto.ContextFields.UserId,
			allowedPermissions,
		).Where("\"BlockPackTable\".deleted_at IS NULL").
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

	allowedPermissions := []enums.AccessControlPermission{
		enums.AccessControlPermission_Owner,
		enums.AccessControlPermission_Admin,
		enums.AccessControlPermission_Write,
		enums.AccessControlPermission_Read,
	}

	resDto := dtos.GetAllMyBlockPacksByRootShelfIdResDto{}

	result := db.Model(&schemas.BlockPack{}).
		Joins("LEFT JOIN \"SubShelfTable\" ss ON \"BlockPackTable\".parent_sub_shelf_id = ss.id").
		Joins("LEFT JOIN \"UsersToShelvesTable\" uts ON ss.root_shelf_id = uts.root_shelf_id").
		Where("ss.root_shelf_id = ? AND uts.user_id = ? AND uts.permission IN ?",
			reqDto.Param.RootShelfId, reqDto.ContextFields.UserId, allowedPermissions,
		).Where("\"BlockPackTable\".deleted_at IS NULL").
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

	db := s.db.WithContext(ctx)

	newBlockPackId, exception := s.blockPackRepository.CreateOneBySubShelfId(
		reqDto.Body.ParentSubShelfId,
		reqDto.ContextFields.UserId,
		inputs.CreateBlockPackInput{
			Name:                reqDto.Body.Name,
			Icon:                reqDto.Body.Icon,
			HeaderBackgroundURL: reqDto.Body.HeaderBackgroundURL,
		},
		options.WithDB(db),
	)
	if exception != nil {
		return nil, exception
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

	db := s.db.WithContext(ctx)

	input := make([]inputs.BulkCreateBlockPackInput, len(reqDto.Body.CreatedBlockPacks))
	for index, createdBlockPack := range reqDto.Body.CreatedBlockPacks {
		input[index] = inputs.BulkCreateBlockPackInput{
			ParentSubShelfId:    createdBlockPack.ParentSubShelfId,
			Name:                createdBlockPack.Name,
			Icon:                createdBlockPack.Icon,
			HeaderBackgroundURL: createdBlockPack.HeaderBackgroundURL,
		}
	}
	newBlockPackIds, exception := s.blockPackRepository.BulkCreateManyBySubShelfIds(
		reqDto.ContextFields.UserId,
		input,
		options.WithDB(db),
	)
	if exception != nil {
		return nil, exception
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

	input := make([]inputs.BulkUpdateBlockPackInput, len(reqDto.Body.UpdatedBlockPacks))
	for index, updatedBlockPack := range reqDto.Body.UpdatedBlockPacks {
		input[index] = inputs.BulkUpdateBlockPackInput{
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
	exception := s.blockPackRepository.BulkUpdateManyByIds(
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

	allowedPermissions := []enums.AccessControlPermission{
		enums.AccessControlPermission_Owner,
		enums.AccessControlPermission_Admin,
		enums.AccessControlPermission_Write,
	}

	if !s.subShelfRepository.HasPermission(
		reqDto.Body.DestinationParentSubShelfId,
		reqDto.ContextFields.UserId,
		allowedPermissions,
		options.WithDB(db),
		options.WithOnlyDeleted(types.Ternary_Negative),
	) {
		return nil, exceptions.Shelf.NoPermission("get the destination sub shelf")
	}

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

func (s *BlockPackService) MoveMyBlockPacksByIds(
	ctx context.Context, reqDto *dtos.MoveMyBlockPacksByIdsReqDto,
) (*dtos.MoveMyBlockPacksByIdsResDto, *exceptions.Exception) {
	if err := validation.Validator.Struct(reqDto); err != nil {
		return nil, exceptions.BlockPack.InvalidDto().WithOrigin(err)
	}

	db := s.db.WithContext(ctx)

	allowedPermissions := []enums.AccessControlPermission{
		enums.AccessControlPermission_Owner,
		enums.AccessControlPermission_Admin,
		enums.AccessControlPermission_Write,
	}

	if !s.subShelfRepository.HasPermission(
		reqDto.Body.DestinationParentSubShelfId,
		reqDto.ContextFields.UserId,
		allowedPermissions,
		options.WithDB(db),
		options.WithOnlyDeleted(types.Ternary_Negative),
	) {
		return nil, exceptions.Shelf.NoPermission("get the destination sub shelf")
	}

	input := make([]inputs.BulkUpdateBlockPackInput, len(reqDto.Body.BlockPackIds))
	for index, blockPackId := range reqDto.Body.BlockPackIds {
		input[index] = inputs.BulkUpdateBlockPackInput{
			Id: blockPackId,
			PartialUpdateInput: inputs.PartialUpdateInput[inputs.UpdateBlockPackInput]{
				Values: inputs.UpdateBlockPackInput{
					ParentSubShelfId: &reqDto.Body.DestinationParentSubShelfId,
				},
			},
		}
	}
	exception := s.blockPackRepository.BulkUpdateManyByIds(
		reqDto.ContextFields.UserId,
		input,
		options.WithDB(db),
	)
	if exception != nil {
		return nil, exception
	}

	return &dtos.MoveMyBlockPacksByIdsResDto{
		UpdatedAt: time.Now(),
	}, nil
}

func (s *BlockPackService) BatchMoveMyBlockPacksByIds(
	ctx context.Context, reqDto *dtos.BatchMoveMyBlockPacksByIdsReqDto,
) (*dtos.BatchMoveMyBlockPacksByIdsResDto, *exceptions.Exception) {
	if err := validation.Validator.Struct(reqDto); err != nil {
		return nil, exceptions.BlockPack.InvalidDto().WithOrigin(err)
	}

	db := s.db.WithContext(ctx)

	allowedPermissions := []enums.AccessControlPermission{
		enums.AccessControlPermission_Owner,
		enums.AccessControlPermission_Admin,
		enums.AccessControlPermission_Write,
	}

	var blockPackIds []uuid.UUID
	destinationParentSubShelfIds := make([]uuid.UUID, len(reqDto.Body.MovedBlockPacks))
	for index, movedBlockPack := range reqDto.Body.MovedBlockPacks {
		blockPackIds = append(blockPackIds, movedBlockPack.BlockPackIds...)
		destinationParentSubShelfIds[index] = movedBlockPack.DestinationParentSubShelfId
	}

	isBlockPackValid := make(map[uuid.UUID]bool)
	validBlockPacks, exception := s.blockPackRepository.CheckPermissionsAndGetManyByIds(
		blockPackIds,
		reqDto.ContextFields.UserId,
		nil,
		allowedPermissions,
		options.WithDB(db),
		options.WithOnlyDeleted(types.Ternary_Negative),
	)
	if exception != nil {
		return nil, exception
	}
	for _, validBlockPack := range validBlockPacks {
		isBlockPackValid[validBlockPack.Id] = true
	}

	isDestinationParentSubShelfValid := make(map[uuid.UUID]bool)
	validDestinationParentSubShelves, exception := s.subShelfRepository.CheckPermissionsAndGetManyByIds(
		destinationParentSubShelfIds,
		reqDto.ContextFields.UserId,
		nil,
		allowedPermissions,
		options.WithDB(db),
		options.WithOnlyDeleted(types.Ternary_Negative),
	)
	if exception != nil {
		return nil, exception
	}
	for _, validDestinationParentSubShelf := range validDestinationParentSubShelves {
		isDestinationParentSubShelfValid[validDestinationParentSubShelf.Id] = true
	}

	var valuePlaceholders []string
	var valueArgs []interface{}
	for _, movedBlockPack := range reqDto.Body.MovedBlockPacks {
		if !isDestinationParentSubShelfValid[movedBlockPack.DestinationParentSubShelfId] {
			continue
		}

		for _, blockPackId := range movedBlockPack.BlockPackIds {
			if !isBlockPackValid[blockPackId] {
				continue
			}

			valuePlaceholders = append(valuePlaceholders, "(?::uuid, ?::uuid)")
			valueArgs = append(valueArgs,
				blockPackId,
				movedBlockPack.DestinationParentSubShelfId,
			)
		}
	}

	sql := fmt.Sprintf(`
		UPDATE "BlockPackTable" AS bp
		SET
			parent_sub_shelf_id = COALESCE(bp.parent_sub_shelf_id, v.dest_parent_sub_shelf_id::uuid),
			updated_at = NOW()
		FROM (VALUES %s) AS v(id, dest_parent_sub_shelf_id)
		WHERE bp.id = v.id::uuid AND bp.deleted_at IS NULL
	`, strings.Join(valuePlaceholders, ","))
	result := s.db.Exec(sql, valueArgs...)
	if exception := exceptions.Cover(nil, []types.Pair[bool, *exceptions.Exception]{
		{First: result.Error != nil, Second: exceptions.BlockPack.FailedToUpdate().WithOrigin(result.Error)},
		{First: result.RowsAffected == 0, Second: exceptions.BlockPack.NoChanges()},
	}); exception != nil {
		return nil, exception
	}

	return &dtos.BatchMoveMyBlockPacksByIdsResDto{
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
