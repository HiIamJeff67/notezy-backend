package exceptions

import (
	"fmt"
	"net/http"
	"unicode"

	shared "notezy-backend/shared"
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
func (d *CacheExceptionSubDomain) NotFound(cachePurpose shared.ValidCachePurpose) *Exception {
	return &Exception{
		Code:           d.BaseCode + 1,
		Prefix:         d.Prefix,
		Message:        fmt.Sprintf("Cannot find the %s in the cache server", convertCamelCaseToSentenceCase(cachePurpose.String())),
		HTTPStatusCode: http.StatusNotFound,
	}
}

func (d *CacheExceptionSubDomain) FailedToCreate(cachePurpose shared.ValidCachePurpose) *Exception {
	return &Exception{
		Code:           d.BaseCode + 2,
		Prefix:         d.Prefix,
		Message:        fmt.Sprintf("Failed to set the %s to the cache server", convertCamelCaseToSentenceCase(cachePurpose.String())),
		HTTPStatusCode: http.StatusInternalServerError,
	}
}

func (d *CacheExceptionSubDomain) FailedToUpdate(cachePurpose shared.ValidCachePurpose) *Exception {
	return &Exception{
		Code:           d.BaseCode + 3,
		Prefix:         d.Prefix,
		Message:        fmt.Sprintf("Failed to update the %s in the cache server", convertCamelCaseToSentenceCase(cachePurpose.String())),
		HTTPStatusCode: http.StatusInternalServerError,
	}
}

func (d *CacheExceptionSubDomain) FailedToDelete(cachePurpose shared.ValidCachePurpose) *Exception {
	return &Exception{
		Code:           d.BaseCode + 4,
		Prefix:         d.Prefix,
		Message:        fmt.Sprintf("Failed to delete the %s in the cache server", convertCamelCaseToSentenceCase(cachePurpose.String())),
		HTTPStatusCode: http.StatusInternalServerError,
	}
}

/* ============================== Handling Connection of the Servers ============================== */

func (d *CacheExceptionSubDomain) FailedToConnectToServer(serverNumber int) *Exception {
	return &Exception{
		Code:           d.BaseCode + 1,
		Prefix:         d.Prefix,
		Message:        fmt.Sprintf("Error on connecting to the redis client server of %v", serverNumber),
		HTTPStatusCode: http.StatusBadGateway,
	}
}

func (d *CacheExceptionSubDomain) FailedToDisconnectToServer(serverNumber int) *Exception {
	return &Exception{
		Code:           d.BaseCode + 2,
		Prefix:         d.Prefix,
		Message:        fmt.Sprintf("Error on disconnecting to the redis client server of %v", serverNumber),
		HTTPStatusCode: http.StatusBadGateway,
	}
}

func (d *CacheExceptionSubDomain) ClientInstanceDoesNotExist(serverNumber int) *Exception {
	return &Exception{
		Code:           d.BaseCode + 3,
		Prefix:         d.Prefix,
		Message:        fmt.Sprintf("The client instance with server number of %v does not exist", serverNumber),
		HTTPStatusCode: http.StatusBadGateway,
	}
}

func (d *CacheExceptionSubDomain) ClientConfigDoesNotExist() *Exception {
	return &Exception{
		Code:           d.BaseCode + 4,
		Prefix:         d.Prefix,
		Message:        "The config of the client instance does not exist",
		HTTPStatusCode: http.StatusBadGateway,
	}
}

/* ============================== Handling Cached Data Type ============================== */

func (d *CacheExceptionSubDomain) InvalidCacheDataStruct(cachedDataStruct any) *Exception {
	return &Exception{
		Code:           d.BaseCode + 11,
		Prefix:         d.Prefix,
		Message:        fmt.Sprintf("Invalid cached data struct detected %v", cachedDataStruct),
		HTTPStatusCode: http.StatusInternalServerError,
	}
}

func (d *CacheExceptionSubDomain) FailedToConvertStructToJson() *Exception {
	return &Exception{
		Code:           d.BaseCode + 12,
		Prefix:         d.Prefix,
		Message:        "Failed to convert struct to json",
		HTTPStatusCode: http.StatusForbidden,
	}
}

func (d *CacheExceptionSubDomain) FailedToConvertJsonToStruct() *Exception {
	return &Exception{
		Code:           d.BaseCode + 13,
		Prefix:         d.Prefix,
		Message:        "Failed to convert json to struct",
		HTTPStatusCode: http.StatusForbidden,
	}
}
