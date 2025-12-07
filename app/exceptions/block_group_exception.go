package exceptions

import (
	"fmt"
	"net/http"

	"github.com/google/uuid"
)

const (
	_ExceptionBaseCode_BlockGroup ExceptionCode = BlockGroupExceptionSubDomainCode * ExceptionSubDomainCodeShiftAmount

	BlockGroupExceptionSubDomainCode ExceptionCode   = 43
	ExceptionBaseCode_BlockGroup     ExceptionCode   = _ExceptionBaseCode_BlockGroup + ReservedExceptionCode
	ExceptionPrefix_BlockGroup       ExceptionPrefix = "BlockGroup"
)

type BlockGroupExceptionDomain struct {
	BaseCode ExceptionCode
	Prefix   ExceptionPrefix
	DatabaseExceptionDomain
	FileExceptionDomain
	TypeExceptionDomain
}

var BlockGroup = &BlockGroupExceptionDomain{
	BaseCode: ExceptionBaseCode_BlockGroup,
	Prefix:   ExceptionPrefix_BlockGroup,
	DatabaseExceptionDomain: DatabaseExceptionDomain{
		_BaseCode: _ExceptionBaseCode_BlockGroup,
		_Prefix:   ExceptionPrefix_BlockGroup,
	},
	FileExceptionDomain: FileExceptionDomain{
		_BaseCode: _ExceptionBaseCode_BlockGroup,
		_Prefix:   ExceptionPrefix_BlockGroup,
	},
	TypeExceptionDomain: TypeExceptionDomain{
		_BaseCode: _ExceptionBaseCode_BlockGroup,
		_Prefix:   ExceptionPrefix_BlockGroup,
	},
}

/* ============================== Block Group Type Error ============================== */

func (d *BlockGroupExceptionDomain) MoreThanOneBlockGroupDetected(expectedBlockGroupId uuid.UUID, anotherBlockGroupId uuid.UUID) *Exception {
	return &Exception{
		Code:           d.BaseCode + 1,
		Prefix:         d.Prefix,
		Reason:         "MoreThanOneBlockGroupDetected",
		IsInternal:     true,
		Message:        fmt.Sprintf("More than 1 block group detected which is invalid in this operation, got %s and %s", expectedBlockGroupId, anotherBlockGroupId),
		HTTPStatusCode: http.StatusInternalServerError,
		LastStackFrame: &GetStackTrace(2, 1)[0],
	}
}

func (d *BlockGroupExceptionDomain) NoRootBlockInBlockGroup(blockGroupId uuid.UUID) *Exception {
	return &Exception{
		Code:           d.BaseCode + 2,
		Prefix:         d.Prefix,
		Reason:         "NoRootBlockInBlockGroup",
		IsInternal:     true,
		Message:        fmt.Sprintf("No root block in the block group of %s", blockGroupId),
		HTTPStatusCode: http.StatusInternalServerError,
		LastStackFrame: &GetStackTrace(2, 1)[0],
	}
}

func (d *BlockGroupExceptionDomain) RepeatedRootBlockInBlockGroupDetected(blockGroupId uuid.UUID, blockId uuid.UUID) *Exception {
	return &Exception{
		Code:           d.BaseCode + 2,
		Prefix:         d.Prefix,
		Reason:         "RepeatedRootBlockInBlockGroupDetected",
		IsInternal:     true,
		Message:        fmt.Sprintf("Block group of %s has multiple blocks as the root block, one of the id of the repeated block is %s", blockGroupId, blockId),
		HTTPStatusCode: http.StatusInternalServerError,
		LastStackFrame: &GetStackTrace(2, 1)[0],
	}
}
