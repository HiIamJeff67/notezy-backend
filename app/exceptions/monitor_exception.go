package exceptions

import (
	"net/http"
	traces "notezy-backend/app/traces"
)

const (
	_ExceptionBaseCode_Monitor ExceptionCode = MonitorExceptionSubDomainCode * ExceptionSubDomainCodeShiftAmount

	MonitorExceptionSubDomainCode ExceptionCode   = 12
	ExceptionBaseCode_Monitor     ExceptionCode   = _ExceptionBaseCode_Monitor + ReservedExceptionCode
	ExceptionPrefix_Monitor       ExceptionPrefix = "Monitor"
)

type MonitorExceptionDomain struct {
	BaseCode ExceptionCode
	Prefix   ExceptionPrefix
	APIExceptionDomain
}

var Monitor = &MonitorExceptionDomain{
	BaseCode: ExceptionBaseCode_Monitor,
	Prefix:   ExceptionPrefix_Monitor,
	APIExceptionDomain: APIExceptionDomain{
		_BaseCode: _ExceptionBaseCode_Monitor,
		_Prefix:   ExceptionPrefix_Monitor,
	},
}

/* ============================== Handling Initialization Error ============================== */

func (d *MonitorExceptionDomain) FailedToInitializeRequestCounter() *Exception {
	return &Exception{
		Code:           d.BaseCode + 1,
		Prefix:         d.Prefix,
		Reason:         "FailedToInitializeRequestCounter",
		IsInternal:     true,
		Message:        "Failed to initialize the request counter via metrics meter",
		HTTPStatusCode: http.StatusInternalServerError,
		LastTrace:      traces.GetTrace(1),
	}
}
