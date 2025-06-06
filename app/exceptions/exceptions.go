package exceptions

import (
	"fmt"
	"net/http"
	"notezy-backend/app/logs"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

/* ============================== Exception Field Type Definition ============================== */
type ExceptionCode int
type ExceptionPrefix string
type ExceptionReason string

const (
	// the first 3 digits are the class of exceptions
	// the last 5 digits are the individual labels for each exceptions
	ExceptionDomainCodeShiftAmount    = 10000000
	ExceptionSubDomainCodeShiftAmount = 100000
	MaxExceptionCode                  = 99999999 // 999 99999
	MinExceptionCode                  = 0        // 000 00000
	// reserve some codes for general use purpose
	// see the below general exceptions ex. NotFound(), FailedToCreate()
	ReservedExceptionCode       = 100 // *** **100, the codes >= *** **100 will be use in the general domain
	DatabaseExceptionDomainCode = 1
	APIExceptionDomainCode      = 2
)

// all the domain prefix shown here, defined in their corresponded files
const (
// ExceptionPrefix_User ExceptionPrefix = "User"                         1
// ExceptionPrefix_UserInfo ExceptionPrefix = "UserInfo"                 2
// ExceptionPrefix_UserAccount ExceptionPrefix = "UserAccount"           3
// ExceptionPrefix_UserSetting ExceptionPrefix = "UserSetting"           4
// ExceptionPrefix_UsersToBadges ExceptionPrefix = "UsersToBadges"       5
// ExceptionPrefix_Badge ExceptionPrefix = "Badge"                       6

// ExceptionPrefix_Cache ExceptionPrefix = "Cache"	   					 1
// ExceptionPrefix_Util ExceptionPrefix = "Util"       					 2
// ExceptionPrefix_Auth       ExceptionPrefix = "Auth" 					 3
)

// global reason for common domain use
// if some individual domain require a custom reason,
// just create one with ExceptionReason type privately whic means its variable name in lower case
const (
	ExceptionReason_UndefinedError            ExceptionReason = "Undefined_Error"
	ExceptionReason_NotFound                  ExceptionReason = "Not_Found"
	ExceptionReason_FailedToCreate            ExceptionReason = "Failed_To_Create"
	ExceptionReason_FailedToUpdate            ExceptionReason = "Failed_To_Update"
	ExceptionReason_FailedToDelete            ExceptionReason = "Failed_To_Delete"
	ExceptionReason_FailedToCommitTransaction ExceptionReason = "Failed_To_Commit_Transaction"
	ExceptionReason_InvalidInput              ExceptionReason = "Invalid_Input"
	ExceptionReason_Timeout                   ExceptionReason = "Timeout"
	ExceptionReason_InvalidDto                ExceptionReason = "Invalid_Dto"
	ExceptionReason_NotImplemented            ExceptionReason = "NotImplement"
	ExceptionReason_InvalidType               ExceptionReason = "Invalid_Type"
)

func IsExceptionCode(exceptionCode int) bool {
	return exceptionCode >= MinExceptionCode && exceptionCode <= MaxExceptionCode
}

/* ============================== Exception Field Type Definition ============================== */

/* ============================== General Exception Structure Definition ============================== */
type Exception struct {
	Code           ExceptionCode   // custom exception code
	Prefix         ExceptionPrefix // custom exception prefix
	Reason         ExceptionReason // custom exception reason
	Message        string          // custom exception message
	HTTPStatusCode int             // http status code
	Details        any             // additional error details (optional)
	Error          error           // original error (optional)
}

func (e *Exception) GetString() string {
	if e.Error != nil {
		return fmt.Sprintf("[%v]%s: %v", e.Code, e.Reason, e.Error)
	}
	return fmt.Sprintf("[%v]%s: %s", e.Code, e.Reason, e.Message)
}

func (e *Exception) GetGinH() *gin.H {
	return &gin.H{
		"code":    e.Code,
		"prefix":  e.Prefix,
		"reason":  e.Reason,
		"message": e.Message,
		"status":  e.HTTPStatusCode,
		"details": e.Details,
		"error":   e.Error,
	}
}

func (e *Exception) WithDetails(details any) *Exception {
	e.Details = details
	return e
}

func (e *Exception) WithError(err error) *Exception {
	e.Error = err
	return e
}

func (e *Exception) Log() *Exception {
	if e.Error != nil {
		logs.FError("[%d] %s: %v", e.Code, e.Message, e.Error)
	} else {
		logs.FError("[%d] %s", e.Code, e.Message)
	}
	return e
}

func (e *Exception) Panic() {
	if e.Error != nil {
		panic(fmt.Sprintf("[%d] %s: %s", e.Code, e.Message, e.Error.Error()))
	} else {
		panic(fmt.Sprintf("[%d] %s", e.Code, e.Message))
	}
}

func (e *Exception) PanicVerbose() {
	if e.Error != nil {
		panic(fmt.Sprintf("[%d] %s: %v", e.Code, e.Message, e.Error))
	} else {
		panic(fmt.Sprintf("[%d] %s", e.Code, e.Message))
	}
}

/* ============================== General Exception Structure Definition ============================== */

/* ============================== Database Exception Domain Definition ============================== */
type DatabaseExceptionDomain struct {
	_BaseCode ExceptionCode
	_Prefix   ExceptionPrefix
}

func (d *DatabaseExceptionDomain) UndefinedError(optionalMessage ...string) *Exception {
	message := fmt.Sprintf("Undefined error happened in %s", strings.ToLower(string(d._Prefix)))
	if len(optionalMessage) > 0 && len(strings.ReplaceAll(optionalMessage[0], " ", "")) > 0 {
		message = optionalMessage[0]
	}

	return &Exception{
		Code:           d._BaseCode + 0,
		Prefix:         d._Prefix,
		Reason:         ExceptionReason_UndefinedError,
		Message:        message,
		HTTPStatusCode: http.StatusBadRequest,
	}
}

func (d *DatabaseExceptionDomain) NotFound(optionalMessage ...string) *Exception {
	message := fmt.Sprintf("%s not found", strings.ToLower(string(d._Prefix)))
	if len(optionalMessage) > 0 && len(strings.ReplaceAll(optionalMessage[0], " ", "")) > 0 {
		message = optionalMessage[0]
	}

	return &Exception{
		Code:           d._BaseCode + 1,
		Prefix:         d._Prefix,
		Reason:         ExceptionReason_NotFound,
		Message:        message,
		HTTPStatusCode: http.StatusNotFound,
	}
}

func (d *DatabaseExceptionDomain) FailedToCreate(optionalMessage ...string) *Exception {
	message := fmt.Sprintf("Failed to create the %s", strings.ToLower(string(d._Prefix)))
	if len(optionalMessage) > 0 && len(strings.ReplaceAll(optionalMessage[0], " ", "")) > 0 {
		message = optionalMessage[0]
	}

	return &Exception{
		Code:           d._BaseCode + 2,
		Prefix:         d._Prefix,
		Reason:         ExceptionReason_FailedToCreate,
		Message:        message,
		HTTPStatusCode: http.StatusInternalServerError,
	}
}

func (d *DatabaseExceptionDomain) FailedToUpdate(optionalMessage ...string) *Exception {
	message := fmt.Sprintf("Failed to update the %s", strings.ToLower(string(d._Prefix)))
	if len(optionalMessage) > 0 && len(strings.ReplaceAll(optionalMessage[0], " ", "")) > 0 {
		message = optionalMessage[0]
	}

	return &Exception{
		Code:           d._BaseCode + 3,
		Prefix:         d._Prefix,
		Reason:         ExceptionReason_FailedToUpdate,
		Message:        message,
		HTTPStatusCode: http.StatusInternalServerError,
	}
}

func (d *DatabaseExceptionDomain) FailedToDelete(optionalMessage ...string) *Exception {
	message := fmt.Sprintf("Failed to delete the %s", strings.ToLower(string(d._Prefix)))
	if len(optionalMessage) > 0 && len(strings.ReplaceAll(optionalMessage[0], " ", "")) > 0 {
		message = optionalMessage[0]
	}

	return &Exception{
		Code:           d._BaseCode + 4,
		Prefix:         d._Prefix,
		Reason:         ExceptionReason_FailedToDelete,
		Message:        message,
		HTTPStatusCode: http.StatusInternalServerError,
	}
}

func (d *DatabaseExceptionDomain) FailedToCommitTransaction(optionalMessage ...string) *Exception {
	message := fmt.Sprintf("Failed to commit the transaction in %s", strings.ToLower(string(d._Prefix)))
	if len(optionalMessage) > 0 && len(strings.ReplaceAll(optionalMessage[0], " ", "")) > 0 {
		message = optionalMessage[0]
	}

	return &Exception{
		Code:           d._BaseCode + 5,
		Prefix:         d._Prefix,
		Reason:         ExceptionReason_FailedToCommitTransaction,
		Message:        message,
		HTTPStatusCode: http.StatusInternalServerError,
	}
}

func (d *DatabaseExceptionDomain) InvalidInput(optionalMessage ...string) *Exception {
	message := fmt.Sprintf("Invalid input object detected in %s", strings.ToLower(string(d._Prefix)))
	if len(optionalMessage) > 0 && len(strings.ReplaceAll(optionalMessage[0], " ", "")) > 0 {
		message = optionalMessage[0]
	}

	return &Exception{
		Code:           d._BaseCode + 6,
		Prefix:         d._Prefix,
		Reason:         ExceptionReason_InvalidInput,
		Message:        message,
		HTTPStatusCode: http.StatusBadRequest,
	}
}

func (d *DatabaseExceptionDomain) NotImplemented(optionalMessage ...string) *Exception {
	message := fmt.Sprintf("Not yet implemented the methods in %s", strings.ToLower(string(d._Prefix)))
	if len(optionalMessage) > 0 && len(strings.ReplaceAll(optionalMessage[0], " ", "")) > 0 {
		message = optionalMessage[0]
	}

	return &Exception{
		Code:           d._BaseCode + 7,
		Prefix:         d._Prefix,
		Reason:         ExceptionReason_NotImplemented,
		Message:        message,
		HTTPStatusCode: http.StatusNotImplemented,
	}
}

func (d *DatabaseExceptionDomain) InvalidType(value any) *Exception {
	return &Exception{
		Code:           d._BaseCode + 8,
		Prefix:         d._Prefix,
		Reason:         ExceptionReason_InvalidType,
		Message:        fmt.Sprintf("Invalid type in %s", strings.ToLower(string(d._Prefix))),
		HTTPStatusCode: http.StatusInternalServerError,
		Details: map[string]any{
			"actualType": fmt.Sprintf("%T", value),
			"value":      value,
		},
	}
}

/* ============================== API Exception Domain Definition ============================== */
type APIExceptionDomain struct {
	_BaseCode ExceptionCode
	_Prefix   ExceptionPrefix
}

func (d *APIExceptionDomain) UndefinedError(optionalMessage ...string) *Exception {
	message := fmt.Sprintf("Undefined error happened in %s", strings.ToLower(string(d._Prefix)))
	if len(optionalMessage) > 0 && len(strings.ReplaceAll(optionalMessage[0], " ", "")) > 0 {
		message = optionalMessage[0]
	}

	return &Exception{
		Code:           d._BaseCode + 0,
		Prefix:         d._Prefix,
		Reason:         ExceptionReason_UndefinedError,
		Message:        message,
		HTTPStatusCode: http.StatusBadRequest,
	}
}

func (d *APIExceptionDomain) Timeout(time time.Duration, optionalMessage ...string) *Exception {
	message := fmt.Sprintf("Timeout in %s with %v", strings.ToLower(string(d._Prefix)), time)
	if len(optionalMessage) > 0 && len(strings.ReplaceAll(optionalMessage[0], " ", "")) > 0 {
		message = optionalMessage[0]
	}

	return &Exception{
		Code:           d._BaseCode + 1,
		Prefix:         d._Prefix,
		Reason:         ExceptionReason_Timeout,
		Message:        message,
		HTTPStatusCode: http.StatusRequestTimeout,
	}
}

func (d *APIExceptionDomain) InvalidDto(optionalMessage ...string) *Exception {
	message := fmt.Sprintf("Invalid dto detected in %s", strings.ToLower(string(d._Prefix)))
	if len(optionalMessage) > 0 && len(strings.ReplaceAll(optionalMessage[0], " ", "")) > 0 {
		message = optionalMessage[0]
	}

	return &Exception{
		Code:           d._BaseCode + 2,
		Prefix:         d._Prefix,
		Reason:         ExceptionReason_InvalidDto,
		Message:        message,
		HTTPStatusCode: http.StatusRequestTimeout,
	}
}

func (d *APIExceptionDomain) NotImplemented(optionalMessage ...string) *Exception {
	message := fmt.Sprintf("Not yet implemented the methods in %s", strings.ToLower(string(d._Prefix)))
	if len(optionalMessage) > 0 && len(strings.ReplaceAll(optionalMessage[0], " ", "")) > 0 {
		message = optionalMessage[0]
	}

	return &Exception{
		Code:           d._BaseCode + 3,
		Prefix:         d._Prefix,
		Reason:         ExceptionReason_NotImplemented,
		Message:        message,
		HTTPStatusCode: http.StatusNotImplemented,
	}
}

func (d *APIExceptionDomain) InvalidType(value any) *Exception {
	return &Exception{
		Code:           d._BaseCode + 4,
		Prefix:         d._Prefix,
		Reason:         ExceptionReason_InvalidType,
		Message:        fmt.Sprintf("Invalid type in %s", strings.ToLower(string(d._Prefix))),
		HTTPStatusCode: http.StatusInternalServerError,
		Details: map[string]any{
			"actualType": fmt.Sprintf("%T", value),
			"value":      value,
		},
	}
}
