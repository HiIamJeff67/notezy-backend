package services

import (
	"context"
	"time"

	"gorm.io/gorm"

	dtos "notezy-backend/app/dtos"
	exceptions "notezy-backend/app/exceptions"
	inputs "notezy-backend/app/models/inputs"
	repositories "notezy-backend/app/models/repositories"
	schemas "notezy-backend/app/models/schemas"
	enums "notezy-backend/app/models/schemas/enums"
	blockpacksql "notezy-backend/app/models/sql/block_pack"
	validation "notezy-backend/app/validation"
	constants "notezy-backend/shared/constants"
	types "notezy-backend/shared/types"
)

/* ============================== Interface & Instance ============================== */

type BlockPackServiceInterface interface {
	GetMyBlockPackById(ctx context.Context, reqDto *dtos.GetMyBlockPackByIdReqDto) (*dtos.GetMyBlockPackByIdResDto, *exceptions.Exception)
	GetMyBlockPackAndItsParentById(ctx context.Context, reqDto *dtos.GetMyBlockPackAndItsParentByIdReqDto) (*dtos.GetMyBlockPackAndItsParentByIdResDto, *exceptions.Exception)
	GetAllMyBlockPacksByParentSubShelfId(ctx context.Context, reqDto *dtos.GetAllMyBlockPacksByParentSubShelfIdReqDto) (*dtos.GetAllMyBlockPacksByParentSubShelfIdResDto, *exceptions.Exception)
	GetAllMyBlockPacksByRootShelfId(ctx context.Context, reqDto *dtos.GetAllMyBlockPacksByRootShelfIdReqDto) (*dtos.GetAllMyBlockPacksByRootShelfIdResDto, *exceptions.Exception)
	CreateBlockPack(ctx context.Context, reqDto *dtos.CreateBlockPackReqDto) (*dtos.CreateBlockPackResDto, *exceptions.Exception)
	UpdateMyBlockPackById(ctx context.Context, reqDto *dtos.UpdateMyBlockPackByIdReqDto) (*dtos.UpdateMyBlockPackByIdResDto, *exceptions.Exception)
	MoveMyBlockPackById(ctx context.Context, reqDto *dtos.MoveMyBlockPackByIdReqDto) (*dtos.MoveMyBlockPackByIdResDto, *exceptions.Exception)
	MoveMyBlockPacksByIds(ctx context.Context, reqDto *dtos.MoveMyBlockPacksByIdsReqDto) (*dtos.MoveMyBlockPacksByIdsResDto, *exceptions.Exception)
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
		return nil, exceptions.BlockPack.InvalidDto().WithError(err)
	}

	db := s.db.WithContext(ctx)

	blockPack, exception := s.blockPackRepository.GetOneById(
		db,
		reqDto.Param.BlockPackId,
		reqDto.ContextFields.UserId,
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
		return nil, exceptions.BlockPack.InvalidDto().WithError(err)
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

	result := db.Raw(blockpacksql.GetMyBlockPackAndItsParentByIdSQL,
		reqDto.Param.BlockPackId, reqDto.ContextFields.UserId, allowedPermissions, onlyDeleted, onlyDeleted,
	).Scan(&resDto)
	if err := result.Error; err != nil {
		return nil, exceptions.BlockPack.NotFound().WithError(err)
	}

	return &resDto, nil
}

func (s *BlockPackService) GetAllMyBlockPacksByParentSubShelfId(
	ctx context.Context, reqDto *dtos.GetAllMyBlockPacksByParentSubShelfIdReqDto,
) (*dtos.GetAllMyBlockPacksByParentSubShelfIdResDto, *exceptions.Exception) {
	if err := validation.Validator.Struct(reqDto); err != nil {
		return nil, exceptions.BlockPack.InvalidDto().WithError(err)
	}

	db := s.db.WithContext(ctx)

	allowedPermissions := []enums.AccessControlPermission{
		enums.AccessControlPermission_Owner,
		enums.AccessControlPermission_Admin,
		enums.AccessControlPermission_Write,
		enums.AccessControlPermission_Read,
	}

	resDto := dtos.GetAllMyBlockPacksByParentSubShelfIdResDto{}

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
		return nil, exceptions.BlockPack.NotFound().WithError(err)
	}

	return &resDto, nil
}

func (s *BlockPackService) GetAllMyBlockPacksByRootShelfId(
	ctx context.Context, reqDto *dtos.GetAllMyBlockPacksByRootShelfIdReqDto,
) (*dtos.GetAllMyBlockPacksByRootShelfIdResDto, *exceptions.Exception) {
	if err := validation.Validator.Struct(reqDto); err != nil {
		return nil, exceptions.BlockPack.InvalidDto().WithError(err)
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
		return nil, exceptions.BlockPack.NotFound().WithError(err)
	}

	return &resDto, nil
}

func (s *BlockPackService) CreateBlockPack(
	ctx context.Context, reqDto *dtos.CreateBlockPackReqDto,
) (*dtos.CreateBlockPackResDto, *exceptions.Exception) {
	if err := validation.Validator.Struct(reqDto); err != nil {
		return nil, exceptions.BlockPack.InvalidDto().WithError(err)
	}

	db := s.db.WithContext(ctx)

	newBlockPackId, exception := s.blockPackRepository.CreateOneBySubShelfId(
		db,
		reqDto.Body.ParentSubShelfId,
		reqDto.ContextFields.UserId,
		inputs.CreateBlockPackInput{
			Name:                reqDto.Body.Name,
			Icon:                reqDto.Body.Icon,
			HeaderBackgroundURL: reqDto.Body.HeaderBackgroundURL,
		},
	)
	if exception != nil {
		return nil, exception
	}
	if newBlockPackId == nil {
		return nil, exceptions.BlockPack.FailedToCreate().WithDetails("got nil block pack id")
	}

	return &dtos.CreateBlockPackResDto{
		Id:        *newBlockPackId,
		CreatedAt: time.Now(),
	}, nil
}

func (s *BlockPackService) UpdateMyBlockPackById(
	ctx context.Context, reqDto *dtos.UpdateMyBlockPackByIdReqDto,
) (*dtos.UpdateMyBlockPackByIdResDto, *exceptions.Exception) {
	if err := validation.Validator.Struct(reqDto); err != nil {
		return nil, exceptions.BlockPack.InvalidDto().WithError(err)
	}

	db := s.db.WithContext(ctx)

	blockPack, exception := s.blockPackRepository.UpdateOneById(
		db,
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
	)
	if exception != nil {
		return nil, exception
	}

	return &dtos.UpdateMyBlockPackByIdResDto{
		UpdatedAt: blockPack.UpdatedAt,
	}, nil
}

func (s *BlockPackService) MoveMyBlockPackById(
	ctx context.Context, reqDto *dtos.MoveMyBlockPackByIdReqDto,
) (*dtos.MoveMyBlockPackByIdResDto, *exceptions.Exception) {
	if err := validation.Validator.Struct(reqDto); err != nil {
		return nil, exceptions.BlockPack.InvalidDto().WithError(err)
	}

	db := s.db.WithContext(ctx)

	allowedPermissions := []enums.AccessControlPermission{
		enums.AccessControlPermission_Owner,
		enums.AccessControlPermission_Admin,
		enums.AccessControlPermission_Write,
	}

	result := db.Exec(blockpacksql.MoveMyBlockPackByIdSQL,
		reqDto.Body.DestinationParentSubShelfId,
		reqDto.Body.BlockPackId,
		reqDto.ContextFields.UserId,
		allowedPermissions,
		reqDto.Body.DestinationParentSubShelfId,
		reqDto.ContextFields.UserId,
		allowedPermissions,
	)

	if err := result.Error; err != nil {
		return nil, exceptions.BlockPack.FailedToUpdate().WithError(err)
	}
	if result.RowsAffected == 0 {
		return nil, exceptions.BlockPack.NoChanges()
	}

	return &dtos.MoveMyBlockPackByIdResDto{
		UpdatedAt: time.Now(),
	}, nil
}

func (s *BlockPackService) MoveMyBlockPacksByIds(
	ctx context.Context, reqDto *dtos.MoveMyBlockPacksByIdsReqDto,
) (*dtos.MoveMyBlockPacksByIdsResDto, *exceptions.Exception) {
	if err := validation.Validator.Struct(reqDto); err != nil {
		return nil, exceptions.BlockPack.InvalidDto().WithError(err)
	}

	db := s.db.WithContext(ctx)

	allowedPermissions := []enums.AccessControlPermission{
		enums.AccessControlPermission_Owner,
		enums.AccessControlPermission_Admin,
		enums.AccessControlPermission_Write,
	}

	result := db.Exec(blockpacksql.MoveMyBlockPacksByIdsSQL,
		reqDto.Body.DestinationParentSubShelfId,
		reqDto.Body.BlockPackIds,
		reqDto.ContextFields.UserId,
		allowedPermissions,
		reqDto.Body.DestinationParentSubShelfId,
		reqDto.ContextFields.UserId,
		allowedPermissions,
	)

	if err := result.Error; err != nil {
		return nil, exceptions.BlockPack.FailedToUpdate().WithError(err)
	}
	if result.RowsAffected == 0 {
		return nil, exceptions.BlockPack.NoChanges()
	}

	return &dtos.MoveMyBlockPacksByIdsResDto{
		UpdatedAt: time.Now(),
	}, nil
}

func (s *BlockPackService) RestoreMyBlockPackById(
	ctx context.Context, reqDto *dtos.RestoreMyBlockPackByIdReqDto,
) (*dtos.RestoreMyBlockPackByIdResDto, *exceptions.Exception) {
	if err := validation.Validator.Struct(reqDto); err != nil {
		return nil, exceptions.BlockPack.InvalidDto().WithError(err)
	}

	db := s.db.WithContext(ctx)

	if exception := s.blockPackRepository.RestoreSoftDeletedOneById(
		db,
		reqDto.Body.BlockPackId,
		reqDto.ContextFields.UserId,
	); exception != nil {
		return nil, exception
	}

	return &dtos.RestoreMyBlockPackByIdResDto{
		UpdatedAt: time.Now(),
	}, nil
}

func (s *BlockPackService) RestoreMyBlockPacksByIds(
	ctx context.Context, reqDto *dtos.RestoreMyBlockPacksByIdsReqDto,
) (*dtos.RestoreMyBlockPacksByIdsResDto, *exceptions.Exception) {
	if err := validation.Validator.Struct(reqDto); err != nil {
		return nil, exceptions.BlockPack.InvalidDto().WithError(err)
	}

	db := s.db.WithContext(ctx)

	if exception := s.blockPackRepository.RestoreSoftDeletedManyByIds(
		db,
		reqDto.Body.BlockPackIds,
		reqDto.ContextFields.UserId,
	); exception != nil {
		return nil, exception
	}

	return &dtos.RestoreMyBlockPacksByIdsResDto{
		UpdatedAt: time.Now(),
	}, nil
}

func (s *BlockPackService) DeleteMyBlockPackById(
	ctx context.Context, reqDto *dtos.DeleteMyBlockPackByIdReqDto,
) (*dtos.DeleteMyBlockPackByIdResDto, *exceptions.Exception) {
	if err := validation.Validator.Struct(reqDto); err != nil {
		return nil, exceptions.BlockPack.InvalidDto().WithError(err)
	}

	db := s.db.WithContext(ctx)

	if exception := s.blockPackRepository.SoftDeleteOneById(
		db,
		reqDto.Body.BlockPackId,
		reqDto.ContextFields.UserId,
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
		return nil, exceptions.BlockPack.InvalidDto().WithError(err)
	}

	db := s.db.WithContext(ctx)

	if exception := s.blockPackRepository.SoftDeleteManyByIds(
		db,
		reqDto.Body.BlockPackIds,
		reqDto.ContextFields.UserId,
	); exception != nil {
		return nil, exception
	}

	return &dtos.DeleteMyBlockPacksByIdsResDto{
		DeletedAt: time.Now(),
	}, nil
}
