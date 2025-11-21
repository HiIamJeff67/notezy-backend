package exceptions

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
}
