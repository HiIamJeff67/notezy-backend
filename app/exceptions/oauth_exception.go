package exceptions

import (
	"fmt"
	"net/http"
	traces "notezy-backend/app/monitor/traces"
)

const (
	_ExceptionBaseCode_OAuth ExceptionCode = OAuthExceptionSubDomainCode * ExceptionSubDomainCodeShiftAmount

	OAuthExceptionSubDomainCode ExceptionCode   = 45
	ExceptionBaseCode_OAuth     ExceptionCode   = _ExceptionBaseCode_OAuth + ReservedExceptionCode
	ExceptionPrefix_OAuth       ExceptionPrefix = "OAuth"
)

type OAuthExceptionDomain struct {
	BaseCode ExceptionCode
	Prefix   ExceptionPrefix
	APIExceptionDomain
	TypeExceptionDomain
}

var OAuth = &OAuthExceptionDomain{
	BaseCode: ExceptionBaseCode_OAuth,
	Prefix:   ExceptionPrefix_OAuth,
	APIExceptionDomain: APIExceptionDomain{
		_BaseCode: _ExceptionBaseCode_OAuth,
		_Prefix:   ExceptionPrefix_OAuth,
	},
	TypeExceptionDomain: TypeExceptionDomain{
		_BaseCode: _ExceptionBaseCode_OAuth,
		_Prefix:   ExceptionPrefix_OAuth,
	},
}

/* ============================== Handling OAuth Errors ============================== */

func (d *OAuthExceptionDomain) InvalidAuthenticationCode(authenticationCode string) *Exception {
	return &Exception{
		Code:           d.BaseCode + 1,
		Prefix:         d.Prefix,
		Reason:         "InvalidAuthenticationCode",
		IsInternal:     false,
		Message:        fmt.Sprintf("Invalid or non-existent authentication code of %s", authenticationCode),
		HTTPStatusCode: http.StatusBadRequest,
		LastTrace:      traces.GetTrace(1),
	}
}

func (d *OAuthExceptionDomain) FailedToExchangeToken(authenticationCode string) *Exception {
	return &Exception{
		Code:           d.BaseCode + 2,
		Prefix:         d.Prefix,
		Reason:         "FailedToExchangeToken",
		IsInternal:     true,
		Message:        fmt.Sprintf("Failed to exchange token to google with authentication code of %s", authenticationCode),
		HTTPStatusCode: http.StatusInternalServerError,
		LastTrace:      traces.GetTrace(1),
	}
}

func (d *OAuthExceptionDomain) FailedToParseResposneFromOAuthThirdParty(thirdPartyName string) *Exception {
	return &Exception{
		Code:           d.BaseCode + 3,
		Prefix:         d.Prefix,
		Reason:         "FailedToParseResposneFromOAuthThirdParty",
		IsInternal:     true,
		Message:        fmt.Sprintf("Failed to parse response from oauth third party of %s", thirdPartyName),
		HTTPStatusCode: http.StatusInternalServerError,
		LastTrace:      traces.GetTrace(1),
	}
}
