package exceptions

import (
	"fmt"
	"net/http"
)

const (
	_ExceptionBaseCode_Storage ExceptionCode = StorageExceptionSubDomainCode * ExceptionSubDomainCodeShiftAmount

	StorageExceptionSubDomainCode ExceptionCode   = 8
	ExceptionBaseCode_Storage     ExceptionCode   = _ExceptionBaseCode_Storage + ReservedExceptionCode
	ExceptionPrefix_Storage       ExceptionPrefix = "Storage"
)

type StorageExceptionDomain struct {
	BaseCode ExceptionCode
	Prefix   ExceptionPrefix
	APIExceptionDomain
	TypeExceptionDomain
}

var Storage = &StorageExceptionDomain{
	BaseCode: ExceptionBaseCode_Storage,
	Prefix:   ExceptionPrefix_Storage,
	APIExceptionDomain: APIExceptionDomain{
		_BaseCode: _ExceptionBaseCode_Storage,
		_Prefix:   ExceptionPrefix_Storage,
	},
	TypeExceptionDomain: TypeExceptionDomain{
		_BaseCode: _ExceptionBaseCode_Storage,
		_Prefix:   ExceptionPrefix_Storage,
	},
}

/* ============================== Handling IO Reader or Writter Error ============================== */

func (d *StorageExceptionDomain) FailedToReadObjectBytes() *Exception {
	return &Exception{
		Code:           d.BaseCode + 1,
		Prefix:         d.Prefix,
		Reason:         "FailedToReadObjectBytes",
		IsInternal:     true,
		Message:        "Failed to read object into bytes",
		HTTPStatusCode: http.StatusInternalServerError,
		LastStackFrame: &GetStackTrace(2, 1)[0],
	}
}

func (d *StorageExceptionDomain) FailedToWriteObjectBytes() *Exception {
	return &Exception{
		Code:           d.BaseCode + 2,
		Prefix:         d.Prefix,
		Reason:         "FailedToWriteObjectBytes",
		IsInternal:     true,
		Message:        "Failed to write object bytes",
		HTTPStatusCode: http.StatusInternalServerError,
		LastStackFrame: &GetStackTrace(2, 1)[0],
	}
}

func (d *StorageExceptionDomain) ObjectTooLarge(size int64, maxSize int64) *Exception {
	return &Exception{
		Code:           d.BaseCode + 3,
		Prefix:         d.Prefix,
		Reason:         "ObjectTooLarge",
		IsInternal:     true,
		Message:        fmt.Sprintf("Object with size of %d is larger than the max size of %d", size, maxSize),
		HTTPStatusCode: http.StatusInternalServerError,
		LastStackFrame: &GetStackTrace(2, 1)[0],
	}
}

/* ============================== Handling CRUD Operations Error for Storage ============================== */

func (d *StorageExceptionDomain) FailedToPutObject(object any) *Exception {
	return &Exception{
		Code:           d.BaseCode + 11,
		Prefix:         d.Prefix,
		Reason:         "FailedToPutObject",
		IsInternal:     true,
		Message:        fmt.Sprintf("Failed to put object of %v", object),
		HTTPStatusCode: http.StatusInternalServerError,
		LastStackFrame: &GetStackTrace(2, 1)[0],
	}
}

func (d *StorageExceptionDomain) FailedToGetObject(key string) *Exception {
	return &Exception{
		Code:           d.BaseCode + 12,
		Prefix:         d.Prefix,
		Reason:         "FailedToGetObject",
		IsInternal:     true,
		Message:        fmt.Sprintf("Failed to get object with key of %s", key),
		HTTPStatusCode: http.StatusNotFound,
		LastStackFrame: &GetStackTrace(2, 1)[0],
	}
}

func (d *StorageExceptionDomain) FailedToDeleteObject(key string) *Exception {
	return &Exception{
		Code:           d.BaseCode + 13,
		Prefix:         d.Prefix,
		Reason:         "FailedToDeleteObject",
		IsInternal:     true,
		Message:        fmt.Sprintf("Failed to delete object with key of %s", key),
		HTTPStatusCode: http.StatusInternalServerError,
		LastStackFrame: &GetStackTrace(2, 1)[0],
	}
}

func (d *StorageExceptionDomain) FailedToPresignPutObject(object any) *Exception {
	return &Exception{
		Code:           d.BaseCode + 14,
		Prefix:         d.Prefix,
		Reason:         "FailedToPresignPutObject",
		IsInternal:     true,
		Message:        fmt.Sprintf("Failed to presigned put object of %v", object),
		HTTPStatusCode: http.StatusInternalServerError,
		LastStackFrame: &GetStackTrace(2, 1)[0],
	}
}

func (d *StorageExceptionDomain) FailedToPresignedGetObject(object any) *Exception {
	return &Exception{
		Code:           d.BaseCode + 15,
		Prefix:         d.Prefix,
		Reason:         "FailedToPresignedGetObject",
		IsInternal:     true,
		Message:        fmt.Sprintf("Failed to presigned get object of %v", object),
		HTTPStatusCode: http.StatusInternalServerError,
		LastStackFrame: &GetStackTrace(2, 1)[0],
	}
}
