package exceptions

const (
	_ExceptionBaseCode_Station ExceptionCode = StationExceptionSubDomainCode * ExceptionSubDomainCodeShiftAmount

	StationExceptionSubDomainCode ExceptionCode   = 48
	ExceptionBaseCode_Station     ExceptionCode   = _ExceptionBaseCode_Station + ReservedExceptionCode
	ExceptionPrefix_Station       ExceptionPrefix = "Station"
)

type StationExceptionDomain struct {
	BaseCode ExceptionCode
	Prefix   ExceptionPrefix
	APIExceptionDomain
	DatabaseExceptionDomain
	TypeExceptionDomain
	FileExceptionDomain
}

var Station = &StationExceptionDomain{
	BaseCode: ExceptionBaseCode_Station,
	Prefix:   ExceptionPrefix_Station,
	APIExceptionDomain: APIExceptionDomain{
		_BaseCode: _ExceptionBaseCode_Station,
		_Prefix:   ExceptionPrefix_Station,
	},
	DatabaseExceptionDomain: DatabaseExceptionDomain{
		_BaseCode: _ExceptionBaseCode_Station,
		_Prefix:   ExceptionPrefix_Station,
	},
	TypeExceptionDomain: TypeExceptionDomain{
		_BaseCode: _ExceptionBaseCode_Station,
		_Prefix:   ExceptionPrefix_Station,
	},
	FileExceptionDomain: FileExceptionDomain{
		_BaseCode: _ExceptionBaseCode_Station,
		_Prefix:   ExceptionPrefix_Station,
	},
}
