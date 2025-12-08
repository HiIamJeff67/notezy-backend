package exceptions

const (
	_ExceptionBaseCode_Block ExceptionCode = BlockExceptionSubDomainCode * ExceptionSubDomainCodeShiftAmount

	BlockExceptionSubDomainCode ExceptionCode   = 44
	ExceptionBaseCode_Block     ExceptionCode   = _ExceptionBaseCode_Block + ReservedExceptionCode
	ExceptionPrefix_Block       ExceptionPrefix = "Block"
)

type BlockExceptionDomain struct {
	BaseCode ExceptionCode
	Prefix   ExceptionPrefix
	DatabaseExceptionDomain
	FileExceptionDomain
	TypeExceptionDomain
}

var Block = &BlockExceptionDomain{
	BaseCode: ExceptionBaseCode_Block,
	Prefix:   ExceptionPrefix_Block,
	DatabaseExceptionDomain: DatabaseExceptionDomain{
		_BaseCode: _ExceptionBaseCode_Block,
		_Prefix:   ExceptionPrefix_Block,
	},
	FileExceptionDomain: FileExceptionDomain{
		_BaseCode: _ExceptionBaseCode_Block,
		_Prefix:   ExceptionPrefix_Block,
	},
	TypeExceptionDomain: TypeExceptionDomain{
		_BaseCode: _ExceptionBaseCode_Block,
		_Prefix:   ExceptionPrefix_Block,
	},
}
