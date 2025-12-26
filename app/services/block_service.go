package services

import (
	"context"

	"gorm.io/gorm"

	dtos "notezy-backend/app/dtos"
	exceptions "notezy-backend/app/exceptions"
	repositories "notezy-backend/app/models/repositories"
	schemas "notezy-backend/app/models/schemas"
	options "notezy-backend/app/options"
	validation "notezy-backend/app/validation"
	types "notezy-backend/shared/types"
)

/* ============================== Interface & Instance ============================== */

type BlockServiceInterface interface {
	GetMyBlockById(ctx context.Context, reqDto *dtos.GetMyBlockByIdReqDto) (*dtos.GetMyBlockByIdResDto, *exceptions.Exception)
	GetAllMyBlocks(ctx context.Context, reqDto *dtos.GetAllMyBlocksReqDto) (*dtos.GetAllMyBlocksResDto, *exceptions.Exception)
}

type BlockService struct {
	db                   *gorm.DB
	blockGroupRepository repositories.BlockGroupRepositoryInterface
	blockRepository      repositories.BlockRepositoryInterface
}

func NewBlockService(
	db *gorm.DB,
	blockGroupRepository repositories.BlockGroupRepositoryInterface,
	blockRepository repositories.BlockRepositoryInterface,
) BlockServiceInterface {
	return &BlockService{
		db:                   db,
		blockGroupRepository: blockGroupRepository,
		blockRepository:      blockRepository,
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

func (s *BlockService) CreateBlock() {}

func (s *BlockService) CreateBlocks() {}

func (s *BlockService) UpdateMyBlockById() {}

func (s *BlockService) UpdateMyBlocksByIds() {}

func (s *BlockService) RestoreMyBlockById() {}

func (s *BlockService) RestoreMyBlocksByIds() {}

func (s *BlockService) DeleteMyBlockById() {}

func (s *BlockService) DeleteMyBlocksByIds() {}
