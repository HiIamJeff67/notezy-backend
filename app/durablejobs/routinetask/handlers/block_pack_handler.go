package handlers

import (
	"github.com/HiIamJeff67/notezy-backend/app/models/repositories"
	"gorm.io/gorm"
)

type BlockPackHandlerInterface interface {
}

type BlockPackHandler struct {
	db                  *gorm.DB
	blockPackRepository repositories.BlockPackRepositoryInterface
}

func NewBlockPackHandler(
	db *gorm.DB,
	blockPackRepository repositories.BlockPackRepositoryInterface,
) BlockPackHandlerInterface {
	return BlockPackHandler{
		db:                  db,
		blockPackRepository: blockPackRepository,
	}
}

func (h *BlockPackHandler) HandleCreateBlockPacks() {

}

func (h *BlockPackHandler) HandleResetBlockPacks() {

}
