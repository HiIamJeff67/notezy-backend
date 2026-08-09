package apiexceptions

import (
	"fmt"
	"net/http"

	"github.com/google/uuid"

	exceptions "github.com/HiIamJeff67/notezy-backend/contracts/types/exceptions"
)

type BlockPackException struct {
	CoreException
}

func NewBlockPackException() BlockPackException {
	return BlockPackException{
		CoreException: NewCoreException("BlockPack"),
	}
}

func (BlockPackException) NoRootBlockInBlockPack(blockPackId uuid.UUID) *exceptions.Exception {
	return exceptions.New(
		"NoRootBlockInBlockPack",
		"BlockPack",
		"Project",
		fmt.Sprintf("No root block exists in block pack %s", blockPackId),
		http.StatusInternalServerError,
		true,
	)
}
