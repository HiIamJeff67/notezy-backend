package exceptions

import "net/http"

const (
	_ExceptionBaseCode_Material ExceptionCode = MaterialExceptionSubDomainCode * ExceptionSubDomainCodeShiftAmount

	MaterialExceptionSubDomainCode ExceptionCode   = 41
	ExceptionBaseCode_Material     ExceptionCode   = _ExceptionBaseCode_Material + ReservedExceptionCode
	ExceptionPrefix_Material       ExceptionPrefix = "Material"
)

type MaterialExceptionDomain struct {
	BaseCode ExceptionCode
	Prefix   ExceptionPrefix
	APIExceptionDomain
	DatabaseExceptionDomain
	TypeExceptionDomain
}

var Material = &MaterialExceptionDomain{
	BaseCode: ExceptionBaseCode_Material,
	Prefix:   ExceptionPrefix_Material,
	APIExceptionDomain: APIExceptionDomain{
		_BaseCode: _ExceptionBaseCode_Material,
		_Prefix:   ExceptionPrefix_Material,
	},
	DatabaseExceptionDomain: DatabaseExceptionDomain{
		_BaseCode: _ExceptionBaseCode_Material,
		_Prefix:   ExceptionPrefix_Material,
	},
	TypeExceptionDomain: TypeExceptionDomain{
		_BaseCode: _ExceptionBaseCode_Material,
		_Prefix:   ExceptionPrefix_Material,
	},
}

func (d *MaterialExceptionDomain) NoPermission() *Exception {
	return &Exception{
		Code:           d.BaseCode + 1,
		Prefix:         d.Prefix,
		Reason:         "NoPermission",
		IsInternal:     false,
		Message:        "You have no permission to do this",
		HTTPStatusCode: http.StatusBadRequest,
		LastStackFrame: &GetStackTrace(2, 1)[0],
	}
}
