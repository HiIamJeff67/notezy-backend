package services

import (
	"context"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"notezy-backend/app/adapters"
	dtos "notezy-backend/app/dtos"
	exceptions "notezy-backend/app/exceptions"
	inputs "notezy-backend/app/models/inputs"
	repositories "notezy-backend/app/models/repositories"
	schemas "notezy-backend/app/models/schemas"
	enums "notezy-backend/app/models/schemas/enums"
	validation "notezy-backend/app/validation"
	constants "notezy-backend/shared/constants"
	types "notezy-backend/shared/types"
)

/* ============================== Interface & Instance ============================== */

type BlockGroupServiceInterface interface {
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
		return nil, exceptions.Material.InvalidDto().WithError(err)
	}

	db := s.db.WithContext(ctx)

	blockGroup, exception := s.blockGroupRepository.GetOneById(
		db,
		reqDto.Param.BlockGroupId,
		reqDto.ContextFields.UserId,
		nil,
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
		return nil, exceptions.Material.InvalidDto().WithError(err)
	}

	db := s.db.WithContext(ctx)

	allowedPermissions := []enums.AccessControlPermission{
		enums.AccessControlPermission_Owner,
		enums.AccessControlPermission_Admin,
		enums.AccessControlPermission_Write,
		enums.AccessControlPermission_Read,
	}

	onlyDeleted := types.Ternary_Negative

	blockGroup, exception := s.blockGroupRepository.CheckPermissionAndGetOneById(
		db,
		reqDto.Param.BlockGroupId,
		reqDto.ContextFields.UserId,
		nil,
		allowedPermissions,
		onlyDeleted,
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

func (s *BlockGroupService) GetMyBlockGroupsByPrevBlockGroupId(
	ctx context.Context, reqDto *dtos.GetMyBlockGroupsByPrevBlockGroupIdReqDto,
) (*dtos.GetMyBlockGroupsByPrevBlockGroupIdResDto, *exceptions.Exception) {
	if err := validation.Validator.Struct(reqDto); err != nil {
		return nil, exceptions.Material.InvalidDto().WithError(err)
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
		return nil, exceptions.Material.InvalidDto().WithError(err)
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

func (s *BlockGroupService) GetMyBlocksGroupsAndTheirBlocksByBlockPackId() {}

func (s *BlockGroupService) CreateBlockGroupByBlockPackId(
	ctx context.Context, reqDto *dtos.CreateBlockGroupByBlockPackIdReqDto,
) (*dtos.CreateBlockGroupByBlockPackIdResDto, *exceptions.Exception) {
	if err := validation.Validator.Struct(reqDto); err != nil {
		return nil, exceptions.Material.InvalidDto().WithError(err)
	}

	db := s.db.WithContext(ctx)

	newBlockGroupId, exception := s.blockGroupRepository.CreateOneByBlockPackId(
		db,
		reqDto.Body.BlockPackId,
		reqDto.ContextFields.UserId,
		inputs.CreateBlockGroupInput{
			PrevBlockGroupId: reqDto.Body.PrevBlockGroupId,
		},
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

func (s *BlockGroupService) CreateBlockGroupAndItsBlocksByBlockPackId(
	ctx context.Context, reqDto *dtos.CreateBlockGroupAndItsBlocksByBlockPackIdReqDto,
) (*dtos.CreateBlockGroupAndItsBlocksByBlockPackIdResDto, *exceptions.Exception) {
	if err := validation.Validator.Struct(reqDto); err != nil {
		return nil, exceptions.Material.InvalidDto().WithError(err)
	}

	tx := s.db.WithContext(ctx).Begin()

	newBlockGroupId, exception := s.blockGroupRepository.CreateOneByBlockPackId(
		tx,
		reqDto.Body.BlockPackId,
		reqDto.ContextFields.UserId,
		inputs.CreateBlockGroupInput{
			PrevBlockGroupId: reqDto.Body.PrevBlockGroupId,
		},
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
		tx,
		*newBlockGroupId,
		reqDto.ContextFields.UserId,
		constants.MaxBatchCreateBlockSize,
		input,
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
		CreatedAt: time.Now(),
	}, nil
}

// TODO: use concurrency to run seperate create operation using FlattenRaw
func (s *BlockGroupService) CreateBlockGroupsAndTheirBlocksByBlockPackId() {

}

func (s *BlockGroupService) MoveMyBlockGroupById() {}

func (s *BlockGroupService) MoveMyBlockGroupsByIds() {}

func (s *BlockGroupService) RestoreMyBlockGroupById() {}

func (s *BlockGroupService) RestoreMyBlockGroupsByIds() {}

func (s *BlockGroupService) DeleteMyBlockGroupById() {}

func (s *BlockGroupService) DeleteMyBlockGroupsByIds() {}
