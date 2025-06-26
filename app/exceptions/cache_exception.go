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

const (
	ExceptionReason_FailedToConnectToServer     ExceptionReason = "Failed_To_Connect_To_Server"
	ExceptionReason_FailedToDisconnectToServer  ExceptionReason = "Failed_To_Disconnect_To_Server"
	ExceptionReason_ClientInstanceDoesNotExist  ExceptionReason = "Client_Instance_Does_Not_Exist"
	ExceptionReason_CLientConfigDoesNotExist    ExceptionReason = "Client_Config_Does_Not_Exist"
	ExceptionReason_InvalidCacheDataStruct      ExceptionReason = "Invalid_Cache_Data_Struct"
	ExceptionReason_FailedToConvertStructToJson ExceptionReason = "Failed_To_Convert_Struct_To_Json"
	ExceptionReason_FailedToConvertJsonToStruct ExceptionReason = "Failed_To_Convert_Json_To_Struct"
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
		Reason:         ExceptionReason_NotFound,
		Message:        fmt.Sprintf("Cannot find the %s in the cache server", convertCamelCaseToSentenceCase(cachePurpose.String())),
		HTTPStatusCode: http.StatusNotFound,
	}
}

func (d *CacheExceptionSubDomain) FailedToCreate(cachePurpose shared.ValidCachePurpose) *Exception {
	return &Exception{
		Code:           d.BaseCode + 2,
		Prefix:         d.Prefix,
		Reason:         ExceptionReason_FailedToCreate,
		Message:        fmt.Sprintf("Failed to set the %s to the cache server", convertCamelCaseToSentenceCase(cachePurpose.String())),
		HTTPStatusCode: http.StatusInternalServerError,
	}
}

func (d *CacheExceptionSubDomain) FailedToUpdate(cachePurpose shared.ValidCachePurpose) *Exception {
	return &Exception{
		Code:           d.BaseCode + 3,
		Prefix:         d.Prefix,
		Reason:         ExceptionReason_FailedToUpdate,
		Message:        fmt.Sprintf("Failed to update the %s in the cache server", convertCamelCaseToSentenceCase(cachePurpose.String())),
		HTTPStatusCode: http.StatusInternalServerError,
	}
}

func (d *CacheExceptionSubDomain) FailedToDelete(cachePurpose shared.ValidCachePurpose) *Exception {
	return &Exception{
		Code:           d.BaseCode + 4,
		Prefix:         d.Prefix,
		Reason:         ExceptionReason_FailedToDelete,
		Message:        fmt.Sprintf("Failed to delete the %s in the cache server", convertCamelCaseToSentenceCase(cachePurpose.String())),
		HTTPStatusCode: http.StatusInternalServerError,
	}
}

/* ============================== Handling Connection of the Servers ============================== */

func (d *CacheExceptionSubDomain) FailedToConnectToServer(serverNumber int) *Exception {
	return &Exception{
		Code:           d.BaseCode + 1,
		Prefix:         d.Prefix,
		Reason:         ExceptionReason_FailedToConnectToServer,
		Message:        fmt.Sprintf("Error on connecting to the redis client server of %v", serverNumber),
		HTTPStatusCode: http.StatusBadGateway,
	}
}

func (d *CacheExceptionSubDomain) FailedToDisconnectToServer(serverNumber int) *Exception {
	return &Exception{
		Code:           d.BaseCode + 2,
		Prefix:         d.Prefix,
		Reason:         ExceptionReason_FailedToDisconnectToServer,
		Message:        fmt.Sprintf("Error on disconnecting to the redis client server of %v", serverNumber),
		HTTPStatusCode: http.StatusBadGateway,
	}
}

func (d *CacheExceptionSubDomain) ClientInstanceDoesNotExist(serverNumber int) *Exception {
	return &Exception{
		Code:           d.BaseCode + 3,
		Prefix:         d.Prefix,
		Reason:         ExceptionReason_ClientInstanceDoesNotExist,
		Message:        fmt.Sprintf("The client instance with server number of %v does not exist", serverNumber),
		HTTPStatusCode: http.StatusBadGateway,
	}
}

func (d *CacheExceptionSubDomain) ClientConfigDoesNotExist() *Exception {
	return &Exception{
		Code:           d.BaseCode + 4,
		Prefix:         d.Prefix,
		Reason:         ExceptionReason_CLientConfigDoesNotExist,
		Message:        "The config of the client instance does not exist",
		HTTPStatusCode: http.StatusBadGateway,
	}
}

/* ============================== Handling Cached Data Type ============================== */

func (d *CacheExceptionSubDomain) InvalidCacheDataStruct(cachedDataStruct any) *Exception {
	return &Exception{
		Code:           d.BaseCode + 11,
		Prefix:         d.Prefix,
		Reason:         ExceptionReason_InvalidCacheDataStruct,
		Message:        fmt.Sprintf("Invalid cached data struct detected %v", cachedDataStruct),
		HTTPStatusCode: http.StatusInternalServerError,
	}
}

func (d *CacheExceptionSubDomain) FailedToConvertStructToJson() *Exception {
	return &Exception{
		Code:           d.BaseCode + 12,
		Prefix:         d.Prefix,
		Reason:         ExceptionReason_FailedToConvertStructToJson,
		Message:        "Failed to convert struct to json",
		HTTPStatusCode: http.StatusForbidden,
	}
}

func (d *CacheExceptionSubDomain) FailedToConvertJsonToStruct() *Exception {
	return &Exception{
		Code:           d.BaseCode + 13,
		Prefix:         d.Prefix,
		Reason:         ExceptionReason_FailedToConvertJsonToStruct,
		Message:        "Failed to convert json to struct",
		HTTPStatusCode: http.StatusForbidden,
	}
}
