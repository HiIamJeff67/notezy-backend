package exceptions

const (
	_ExceptionBaseCode_Material ExceptionCode = MaterialExceptionSubDomainCode * ExceptionSubDomainCodeShiftAmount

	MaterialExceptionSubDomainCode ExceptionCode   = 41
	ExceptionBaseCode_Material     ExceptionCode   = _ExceptionBaseCode_Material + ReservedExceptionCode
	ExceptionPrefix_Material       ExceptionPrefix = "Material"
)

type MaterialExceptionDomain struct {
	BaseCode ExceptionCode
	Prefix   ExceptionPrefix
	DatabaseExceptionDomain
	TypeExceptionDomain
}

var Material = &MaterialExceptionDomain{
	BaseCode: ExceptionBaseCode_Material,
	Prefix:   ExceptionPrefix_Material,
	DatabaseExceptionDomain: DatabaseExceptionDomain{
		_BaseCode: _ExceptionBaseCode_Material,
		_Prefix:   ExceptionPrefix_Material,
	},
	TypeExceptionDomain: TypeExceptionDomain{
		_BaseCode: _ExceptionBaseCode_Material,
		_Prefix:   ExceptionPrefix_Material,
	},
}
