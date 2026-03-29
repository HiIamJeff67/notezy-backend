package exceptions

import (
	"fmt"
	"net/http"
	traces "notezy-backend/app/monitor/traces"

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
		LastTrace:      traces.GetTrace(1),
	}
}

func (d *BlockGroupExceptionDomain) NoRootBlockInBlockGroup(blockGroupId uuid.UUID) *Exception {
	return &Exception{
		Code:           d.BaseCode + 2,
		Prefix:         d.Prefix,
		Reason:         "NoRootBlockInBlockGroup",
		IsInternal:     true,
		Message:        fmt.Sprintf("No root blocks in the block group of %s", blockGroupId),
		HTTPStatusCode: http.StatusInternalServerError,
		LastTrace:      traces.GetTrace(1),
	}
}

func (d *BlockGroupExceptionDomain) RepeatedRootBlockInBlockGroupDetected(blockGroupId uuid.UUID, blockId uuid.UUID) *Exception {
	return &Exception{
		Code:           d.BaseCode + 3,
		Prefix:         d.Prefix,
		Reason:         "RepeatedRootBlockInBlockGroupDetected",
		IsInternal:     true,
		Message:        fmt.Sprintf("Block group of %s has multiple blocks as the root block, one of the id of the repeated block is %s", blockGroupId, blockId),
		HTTPStatusCode: http.StatusInternalServerError,
		LastTrace:      traces.GetTrace(1),
	}
}

func (d *BlockGroupExceptionDomain) BrokenBlockGroupsLinkedListDetected(blockPackId uuid.UUID, blockGroupIds []uuid.UUID) *Exception {
	return &Exception{
		Code:           d.BaseCode + 4,
		Prefix:         d.Prefix,
		Reason:         "BrokenBlockGroupsLinkedListDetected",
		IsInternal:     true,
		Message:        fmt.Sprintf("Block groups of ids within %v in block pack of %s is broken and not in a valid linked list structure", blockGroupIds, blockPackId),
		HTTPStatusCode: http.StatusInternalServerError,
		LastTrace:      traces.GetTrace(1),
	}
}

func (d *BlockGroupExceptionDomain) DuplicateBlockGroupsWithSamePrevBlockGroupId(blockPackId uuid.UUID) *Exception {
	return &Exception{
		Code:           d.BaseCode + 5,
		Prefix:         d.Prefix,
		Reason:         "DuplicateBlockGroupsWithSamePrevBlockGroupId",
		IsInternal:     true,
		Message:        fmt.Sprintf("There're more than 2 block groups in block pack of %s have the same prev block group id", blockPackId),
		HTTPStatusCode: http.StatusInternalServerError,
		LastTrace:      traces.GetTrace(1),
	}
}
