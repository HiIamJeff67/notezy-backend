package exceptions

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
