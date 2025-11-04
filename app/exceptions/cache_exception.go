package exceptions

import (
	"fmt"
	"net/http"
	"unicode"
)

const (
	_ExceptionBaseCode_Cache ExceptionCode = CacheExceptionSubDomainCode * ExceptionSubDomainCodeShiftAmount

	CacheExceptionSubDomainCode ExceptionCode   = 3
	ExceptionBaseCode_Cache     ExceptionCode   = _ExceptionBaseCode_Cache + ReservedExceptionCode
	ExceptionPrefix_Cache       ExceptionPrefix = "Cache"
)

type CacheExceptionSubDomain struct {
	BaseCode           ExceptionCode
	Prefix             ExceptionPrefix
	APIExceptionDomain APIExceptionDomain
}

var Cache = &CacheExceptionSubDomain{
	BaseCode: ExceptionBaseCode_Cache,
	Prefix:   ExceptionPrefix_Cache,
	APIExceptionDomain: APIExceptionDomain{
		_BaseCode: _ExceptionBaseCode_Cache,
		_Prefix:   ExceptionPrefix_Cache,
	},
}

/* ============================== Temporary Function to Convert Camel Case to Sentence Case ============================== */

func convertCamelCaseToSentenceCase(camelCaseString string) string {
	var result []rune
	for index, r := range camelCaseString {
		if unicode.IsUpper(r) && index != 0 {
			result = append(result, ' ')
		}
		result = append(result, unicode.ToLower(r))
	}
	return string(result)
}

/* ============================== Handling Cached Data in the Servers (overriding methods) ============================== */
func (d *CacheExceptionSubDomain) NotFound(cachePurpose string) *Exception {
	return &Exception{
		Code:           d.BaseCode + 1,
		Prefix:         d.Prefix,
		Reason:         "NotFound",
		IsInternal:     true,
		Message:        fmt.Sprintf("Cannot find the %s in the cache server", convertCamelCaseToSentenceCase(cachePurpose)),
		HTTPStatusCode: http.StatusNotFound,
		LastStackFrame: &GetStackTrace(2, 1)[0],
	}
}

func (d *CacheExceptionSubDomain) FailedToCreate(cachePurpose string) *Exception {
	return &Exception{
		Code:           d.BaseCode + 2,
		Prefix:         d.Prefix,
		Reason:         "FailedToCreate",
		IsInternal:     true,
		Message:        fmt.Sprintf("Failed to set the %s to the cache server", convertCamelCaseToSentenceCase(cachePurpose)),
		HTTPStatusCode: http.StatusInternalServerError,
		LastStackFrame: &GetStackTrace(2, 1)[0],
	}
}

func (d *CacheExceptionSubDomain) FailedToUpdate(cachePurpose string) *Exception {
	return &Exception{
		Code:           d.BaseCode + 3,
		Prefix:         d.Prefix,
		Reason:         "FailedToUpdate",
		IsInternal:     true,
		Message:        fmt.Sprintf("Failed to update the %s in the cache server", convertCamelCaseToSentenceCase(cachePurpose)),
		HTTPStatusCode: http.StatusInternalServerError,
		LastStackFrame: &GetStackTrace(2, 1)[0],
	}
}

func (d *CacheExceptionSubDomain) FailedToDelete(cachePurpose string) *Exception {
	return &Exception{
		Code:           d.BaseCode + 4,
		Prefix:         d.Prefix,
		Reason:         "FailedToDelete",
		IsInternal:     true,
		Message:        fmt.Sprintf("Failed to delete the %s in the cache server", convertCamelCaseToSentenceCase(cachePurpose)),
		HTTPStatusCode: http.StatusInternalServerError,
		LastStackFrame: &GetStackTrace(2, 1)[0],
	}
}

/* ============================== Handling Connection of the Servers ============================== */

func (d *CacheExceptionSubDomain) RedisServerNumberNotFound() *Exception {
	return &Exception{
		Code:           d.BaseCode + 11,
		Prefix:         d.Prefix,
		Reason:         "RedisServerNumberNotFound",
		IsInternal:     true,
		Message:        "Redis server number not found or out of range",
		HTTPStatusCode: http.StatusInternalServerError,
		LastStackFrame: &GetStackTrace(2, 1)[0],
	}
}

func (d *CacheExceptionSubDomain) BackendServerNameNotReferenced(cachePurpose string) *Exception {
	return &Exception{
		Code:           d.BaseCode + 12,
		Prefix:         d.Prefix,
		Reason:         "BackendServerNameNotReferenced",
		IsInternal:     true,
		Message:        fmt.Sprintf("The backend server name is not referenced to %s", cachePurpose),
		HTTPStatusCode: http.StatusInternalServerError,
		LastStackFrame: &GetStackTrace(2, 1)[0],
	}
}

func (d *CacheExceptionSubDomain) FailedToConnectToServer(serverNumber int) *Exception {
	return &Exception{
		Code:           d.BaseCode + 13,
		Prefix:         d.Prefix,
		Reason:         "FailedToConnectToServer",
		IsInternal:     true,
		Message:        fmt.Sprintf("Error on connecting to the redis client server of %v", serverNumber),
		HTTPStatusCode: http.StatusBadGateway,
		LastStackFrame: &GetStackTrace(2, 1)[0],
	}
}

func (d *CacheExceptionSubDomain) FailedToDisconnectToServer(serverNumber int) *Exception {
	return &Exception{
		Code:           d.BaseCode + 14,
		Reason:         "FailedToDisconnectToServer",
		Prefix:         d.Prefix,
		Message:        fmt.Sprintf("Error on disconnecting to the redis client server of %v", serverNumber),
		HTTPStatusCode: http.StatusBadGateway,
		LastStackFrame: &GetStackTrace(2, 1)[0],
	}
}

func (d *CacheExceptionSubDomain) ClientInstanceDoesNotExist() *Exception {
	return &Exception{
		Code:           d.BaseCode + 15,
		Prefix:         d.Prefix,
		Reason:         "ClientInstanceDoesNotExist",
		IsInternal:     true,
		Message:        "The client instance does not exist, maybe the authentication is expired",
		HTTPStatusCode: http.StatusBadGateway,
		LastStackFrame: &GetStackTrace(2, 1)[0],
	}
}

func (d *CacheExceptionSubDomain) ClientConfigDoesNotExist() *Exception {
	return &Exception{
		Code:           d.BaseCode + 16,
		Prefix:         d.Prefix,
		Reason:         "ClientConfigDoesNotExist",
		IsInternal:     true,
		Message:        "The config of the client instance does not exist",
		HTTPStatusCode: http.StatusBadGateway,
		LastStackFrame: &GetStackTrace(2, 1)[0],
	}
}

/* ============================== Handling Cached Data Type ============================== */

func (d *CacheExceptionSubDomain) InvalidCacheDataStruct(cachedDataStruct any) *Exception {
	return &Exception{
		Code:           d.BaseCode + 21,
		Prefix:         d.Prefix,
		Reason:         "InvalidCacheDataStruct",
		IsInternal:     true,
		Message:        fmt.Sprintf("Invalid cached data struct detected %v", cachedDataStruct),
		HTTPStatusCode: http.StatusInternalServerError,
		LastStackFrame: &GetStackTrace(2, 1)[0],
	}
}

func (d *CacheExceptionSubDomain) FailedToConvertStructToJson() *Exception {
	return &Exception{
		Code:           d.BaseCode + 22,
		Prefix:         d.Prefix,
		Reason:         "FailedToConvertStructToJson",
		IsInternal:     true,
		Message:        "Failed to convert struct to json",
		HTTPStatusCode: http.StatusForbidden,
		LastStackFrame: &GetStackTrace(2, 1)[0],
	}
}

func (d *CacheExceptionSubDomain) FailedToConvertJsonToStruct() *Exception {
	return &Exception{
		Code:           d.BaseCode + 23,
		Prefix:         d.Prefix,
		Reason:         "FailedToConvertJsonToStruct",
		IsInternal:     true,
		Message:        "Failed to convert json to struct",
		HTTPStatusCode: http.StatusForbidden,
		LastStackFrame: &GetStackTrace(2, 1)[0],
	}
}

func (d *CacheExceptionSubDomain) InvalidFormattedKey() *Exception {
	return &Exception{
		Code:           d.BaseCode + 24,
		Prefix:         d.Prefix,
		Reason:         "InvalidFormattedKey",
		IsInternal:     true,
		Message:        "Invalid formattedkey",
		HTTPStatusCode: http.StatusInternalServerError,
		LastStackFrame: &GetStackTrace(2, 1)[0],
	}
}
