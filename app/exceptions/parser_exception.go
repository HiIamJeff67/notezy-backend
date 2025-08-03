package exceptions

import (
	"fmt"
	"net/http"
)

const (
	_ExceptionBaseCode_Parser ExceptionCode = ParserExceptionSubDomainCode * ExceptionSubDomainCodeShiftAmount

	ParserExceptionSubDomainCode ExceptionCode   = 39
	ExceptionBaseCode_Parser     ExceptionCode   = _ExceptionBaseCode_Parser + ReservedExceptionCode
	ExceptionPrefix_Parser       ExceptionPrefix = "Parser"
)

type ParserExceptionDomain struct {
	BaseCode ExceptionCode
	Prefix   ExceptionPrefix
	APIExceptionDomain
	TypeExceptionDomain
}

var Parser = &ParserExceptionDomain{
	BaseCode: ExceptionBaseCode_Parser,
	Prefix:   ExceptionPrefix_Parser,
	APIExceptionDomain: APIExceptionDomain{
		_BaseCode: _ExceptionBaseCode_Parser,
		_Prefix:   ExceptionPrefix_Parser,
	},
	TypeExceptionDomain: TypeExceptionDomain{
		_BaseCode: _ExceptionBaseCode_Parser,
		_Prefix:   ExceptionPrefix_Parser,
	},
}

func (d *ParserExceptionDomain) FailedToParseFromStringToTime(timeString string) *Exception {
	return &Exception{
		Code:           d.BaseCode + 1,
		Prefix:         d.Prefix,
		Reason:         "FailedToParseFromStringToTime",
		IsInternal:     false,
		Message:        fmt.Sprintf("Invalid time format: %s", timeString),
		HTTPStatusCode: http.StatusUnauthorized,
		LastStackFrame: &GetStackTrace(2, 1)[0],
	}
}
