package exceptions

const (
	_ExceptionBaseCode_Item ExceptionCode = ItemExceptionSubDomainCode * ExceptionSubDomainCodeShiftAmount

	ItemExceptionSubDomainCode ExceptionCode   = 47
	ExceptionBaseCode_Item     ExceptionCode   = _ExceptionBaseCode_Item + ReservedExceptionCode
	ExceptionPrefix_Item       ExceptionPrefix = "Item"
)

type ItemExceptionDomain struct {
	BaseCode ExceptionCode
	Prefix   ExceptionPrefix
	APIExceptionDomain
	DatabaseExceptionDomain
	TypeExceptionDomain
	FileExceptionDomain
}

var Item = &ItemExceptionDomain{
	BaseCode: ExceptionBaseCode_Item,
	Prefix:   ExceptionPrefix_Item,
	APIExceptionDomain: APIExceptionDomain{
		_BaseCode: _ExceptionBaseCode_Item,
		_Prefix:   ExceptionPrefix_Item,
	},
	DatabaseExceptionDomain: DatabaseExceptionDomain{
		_BaseCode: _ExceptionBaseCode_Item,
		_Prefix:   ExceptionPrefix_Item,
	},
	TypeExceptionDomain: TypeExceptionDomain{
		_BaseCode: _ExceptionBaseCode_Item,
		_Prefix:   ExceptionPrefix_Item,
	},
	FileExceptionDomain: FileExceptionDomain{
		_BaseCode: _ExceptionBaseCode_Item,
		_Prefix:   ExceptionPrefix_Item,
	},
}
