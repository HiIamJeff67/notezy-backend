package exceptions

import (
	"fmt"
	"net/http"
)

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
	FileExceptionDomain
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
	FileExceptionDomain: FileExceptionDomain{
		_BaseCode: _ExceptionBaseCode_Material,
		_Prefix:   ExceptionPrefix_Material,
	},
}

/* ============================== Handling Material and File Integration Errors ============================== */

func (d *MaterialExceptionDomain) MaterialTypeNotMatch(
	id string,
	currentType interface{},
	expectedType interface{},
) *Exception {
	return &Exception{
		Code:           d.BaseCode + 1,
		Prefix:         d.Prefix,
		Reason:         "MaterialTypeNotMatch",
		IsInternal:     true,
		Message:        fmt.Sprintf("Try to manipulate a material with id of %s which type required to %v, but got type of %v", id, expectedType, currentType),
		HTTPStatusCode: http.StatusInternalServerError,
		LastStackFrame: &GetStackTrace(2, 1)[0],
	}
}

func (d *MaterialExceptionDomain) MaterialContentTypeNotAllowedInMaterialType(
	id string,
	materialType string,
	currentContentType string,
	expectedContentTypes []string,
) *Exception {
	return &Exception{
		Code:       d.BaseCode + 2,
		Prefix:     d.Prefix,
		Reason:     "MaterialContentTypeNotAllowedInMaterialType",
		IsInternal: false,
		Message: fmt.Sprintf("The type of material with id of %s is %s which allowed content type to be %s, but got %v",
			id, materialType, expectedContentTypes, currentContentType,
		),
		HTTPStatusCode: http.StatusUnsupportedMediaType,
		LastStackFrame: &GetStackTrace(2, 1)[0],
	}
}
