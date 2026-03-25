package exceptions

import (
	"fmt"
	"net/http"
	traces "notezy-backend/app/traces"

	"github.com/google/uuid"
)

const (
	_ExceptionBaseCode_BlockPack ExceptionCode = BlockPackExceptionSubDomainCode * ExceptionSubDomainCodeShiftAmount

	BlockPackExceptionSubDomainCode ExceptionCode   = 42
	ExceptionBaseCode_BlockPack     ExceptionCode   = _ExceptionBaseCode_BlockPack + ReservedExceptionCode
	ExceptionPrefix_BlockPack       ExceptionPrefix = "BlockPack"
)

type BlockPackExceptionDomain struct {
	BaseCode ExceptionCode
	Prefix   ExceptionPrefix
	DatabaseExceptionDomain
	FileExceptionDomain
	TypeExceptionDomain
}

var BlockPack = &BlockPackExceptionDomain{
	BaseCode: ExceptionBaseCode_BlockPack,
	Prefix:   ExceptionPrefix_BlockPack,
	DatabaseExceptionDomain: DatabaseExceptionDomain{
		_BaseCode: _ExceptionBaseCode_BlockPack,
		_Prefix:   ExceptionPrefix_BlockPack,
	},
	FileExceptionDomain: FileExceptionDomain{
		_BaseCode: _ExceptionBaseCode_BlockPack,
		_Prefix:   ExceptionPrefix_BlockPack,
	},
	TypeExceptionDomain: TypeExceptionDomain{
		_BaseCode: _ExceptionBaseCode_BlockPack,
		_Prefix:   ExceptionPrefix_BlockPack,
	},
}

func (d *BlockPackExceptionDomain) NoRootBlockGroupInBlockPack(blockPackId uuid.UUID) *Exception {
	return &Exception{
		Code:           d.BaseCode + 1,
		Prefix:         d.Prefix,
		Reason:         "NoRootBlockGroupInBlockPack",
		IsInternal:     true,
		Message:        fmt.Sprintf("No root block groups in the block pack of %s", blockPackId),
		HTTPStatusCode: http.StatusInternalServerError,
		LastTrace:      traces.GetTrace(1),
	}
}
