package apiexceptions

import (
	"fmt"
	"net/http"

	"github.com/google/uuid"

	exceptions "github.com/HiIamJeff67/notezy-backend/contracts/types/exceptions"
)

type blockPackExceptionDomain struct {
	domainException
}

var BlockPack = blockPackExceptionDomain{
	domainException: newDomainException("BlockPack"),
}

func (blockPackExceptionDomain) NoRootBlockInBlockPack(blockPackId uuid.UUID) *exceptions.Exception {
	return exceptions.New(
		"NoRootBlockInBlockPack",
		"BlockPack",
		"Project",
		fmt.Sprintf("No root block exists in block pack %s", blockPackId),
		http.StatusInternalServerError,
		true,
	)
}
