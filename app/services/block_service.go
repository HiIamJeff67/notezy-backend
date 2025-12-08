package services

import (
	"gorm.io/gorm"

	repositories "notezy-backend/app/models/repositories"
)

// This service includes the business logic of block group, block, etc.

/* ============================== Interface & Instance ============================== */

type BlockServiceInterface interface {
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

func (s *BlockService) CreateBlocksByBlockGroupId() {}
