package services

import (
	"gorm.io/gorm"

	repositories "notezy-backend/app/models/repositories"
)

/* ============================== Interface & Instance ============================== */

type BlockPackServiceInterface interface {
}

type BlockPackService struct {
	db                   *gorm.DB
	blockPackRepository  repositories.BlockPackRepositoryInterface
	blockGroupRepository repositories.BlockGroupRepositoryInterface
}

func NewBlockPackService(
	db *gorm.DB,
	blockPackRepository repositories.BlockPackRepositoryInterface,
	blockGroupRepository repositories.BlockGroupRepositoryInterface,
) BlockPackServiceInterface {
	return &BlockPackService{
		db:                   db,
		blockPackRepository:  blockPackRepository,
		blockGroupRepository: blockGroupRepository,
	}
}

/* ============================== Service Methods for Block Pack ============================== */
