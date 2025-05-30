package exceptions

import (
	"fmt"
	"net/http"
	"notezy-backend/app/logs"
	"strings"
	"time"
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
// ExceptionPrefix_User ExceptionPrefix = "User"
// ExceptionPrefix_UserInfo ExceptionPrefix = "UserInfo"
// ExceptionPrefix_UserAccount ExceptionPrefix = "UserAccount"
// ExceptionPrefix_UserSetting ExceptionPrefix = "UserSetting"
// ExceptionPrefix_UsersToBadges ExceptionPrefix = "UsersToBadges"
// ExceptionPrefix_Badge ExceptionPrefix = "Badge"

// ExceptionPrefix_Cache ExceptionPrefix = "Cache"
// ExceptionPrefix_Util ExceptionPrefix = "Util"
)

// global reason for common domain use
// if some individual domain require a custom reason,
// just create one with ExceptionReason type privately whic means its variable name in lower case
const (
	ExceptionReason_NotFound                  ExceptionReason = "Not_Found"
	ExceptionReason_FailedToCreate            ExceptionReason = "Failed_To_Create"
	ExceptionReason_FailedToUpdate            ExceptionReason = "Failed_To_Update"
	ExceptionReason_FailedToDelete            ExceptionReason = "Failed_To_Delete"
	ExceptionReason_FailedToCommitTransaction ExceptionReason = "Failed_To_Commit_Transaction"
	ExceptionReason_InvalidInput              ExceptionReason = "Invalid_Input"
	ExceptionReason_Timeout                   ExceptionReason = "Timeout"
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
	BaseCode ExceptionCode
	Prefix   ExceptionPrefix
}

func (d *DatabaseExceptionDomain) NotFound(optionalMessage ...string) *Exception {
	message := fmt.Sprintf("%s not found", strings.ToLower(string(d.Prefix)))
	if len(optionalMessage) > 0 && len(strings.ReplaceAll(optionalMessage[0], " ", "")) > 0 {
		message = optionalMessage[0]
	}

	return &Exception{
		Code:           d.BaseCode + 1,
		Prefix:         d.Prefix,
		Reason:         ExceptionReason_NotFound,
		Message:        message,
		HTTPStatusCode: http.StatusNotFound,
	}
}

func (d *DatabaseExceptionDomain) FailedToCreate(optionalMessage ...string) *Exception {
	message := fmt.Sprintf("Failed to create the %s", strings.ToLower(string(d.Prefix)))
	if len(optionalMessage) > 0 && len(strings.ReplaceAll(optionalMessage[0], " ", "")) > 0 {
		message = optionalMessage[0]
	}

	return &Exception{
		Code:           d.BaseCode + 2,
		Prefix:         d.Prefix,
		Reason:         ExceptionReason_FailedToCreate,
		Message:        message,
		HTTPStatusCode: http.StatusInternalServerError,
	}
}

func (d *DatabaseExceptionDomain) FailedToUpdate(optionalMessage ...string) *Exception {
	message := fmt.Sprintf("Failed to update the %s", strings.ToLower(string(d.Prefix)))
	if len(optionalMessage) > 0 && len(strings.ReplaceAll(optionalMessage[0], " ", "")) > 0 {
		message = optionalMessage[0]
	}

	return &Exception{
		Code:           d.BaseCode + 3,
		Prefix:         d.Prefix,
		Reason:         ExceptionReason_FailedToUpdate,
		Message:        message,
		HTTPStatusCode: http.StatusInternalServerError,
	}
}

func (d *DatabaseExceptionDomain) FailedToDelete(optionalMessage ...string) *Exception {
	message := fmt.Sprintf("Failed to delete the %s", strings.ToLower(string(d.Prefix)))
	if len(optionalMessage) > 0 && len(strings.ReplaceAll(optionalMessage[0], " ", "")) > 0 {
		message = optionalMessage[0]
	}

	return &Exception{
		Code:           d.BaseCode + 4,
		Prefix:         d.Prefix,
		Reason:         ExceptionReason_FailedToDelete,
		Message:        message,
		HTTPStatusCode: http.StatusInternalServerError,
	}
}

func (d *DatabaseExceptionDomain) FailedToCommitTransaction(optionalMessage ...string) *Exception {
	message := fmt.Sprintf("Failed to commit the transaction in %s", strings.ToLower(string(d.Prefix)))
	if len(optionalMessage) > 0 && len(strings.ReplaceAll(optionalMessage[0], " ", "")) > 0 {
		message = optionalMessage[0]
	}

	return &Exception{
		Code:           d.BaseCode + 5,
		Prefix:         d.Prefix,
		Reason:         ExceptionReason_FailedToCommitTransaction,
		Message:        message,
		HTTPStatusCode: http.StatusInternalServerError,
	}
}

func (d *DatabaseExceptionDomain) InvalidInput(optionalMessage ...string) *Exception {
	message := fmt.Sprintf("Invalid input object detected in %s", strings.ToLower(string(d.Prefix)))
	if len(optionalMessage) > 0 && len(strings.ReplaceAll(optionalMessage[0], " ", "")) > 0 {
		message = optionalMessage[0]
	}

	return &Exception{
		Code:           d.BaseCode + 6,
		Prefix:         d.Prefix,
		Reason:         ExceptionReason_InvalidInput,
		Message:        message,
		HTTPStatusCode: http.StatusBadRequest,
	}
}

/* ============================== API Exception Domain Definition ============================== */
type APIExceptionDomain struct {
	BaseCode ExceptionCode
	Prefix   ExceptionPrefix
}

func (d *APIExceptionDomain) NotFound(optionalMessage ...string) *Exception {
	message := fmt.Sprintf("%s not found", strings.ToLower(string(d.Prefix)))
	if len(optionalMessage) > 0 && len(strings.ReplaceAll(optionalMessage[0], " ", "")) > 0 {
		message = optionalMessage[0]
	}

	return &Exception{
		Code:           d.BaseCode + 1,
		Prefix:         d.Prefix,
		Reason:         ExceptionReason_NotFound,
		Message:        message,
		HTTPStatusCode: http.StatusNotFound,
	}
}

func (d *APIExceptionDomain) FailedToCreate(optionalMessage ...string) *Exception {
	message := fmt.Sprintf("Failed to create the %s", strings.ToLower(string(d.Prefix)))
	if len(optionalMessage) > 0 && len(strings.ReplaceAll(optionalMessage[0], " ", "")) > 0 {
		message = optionalMessage[0]
	}

	return &Exception{
		Code:           d.BaseCode + 2,
		Prefix:         d.Prefix,
		Reason:         ExceptionReason_FailedToCreate,
		Message:        message,
		HTTPStatusCode: http.StatusInternalServerError,
	}
}

func (d *APIExceptionDomain) FailedToUpdate(optionalMessage ...string) *Exception {
	message := fmt.Sprintf("Failed to update the %s", strings.ToLower(string(d.Prefix)))
	if len(optionalMessage) > 0 && len(strings.ReplaceAll(optionalMessage[0], " ", "")) > 0 {
		message = optionalMessage[0]
	}

	return &Exception{
		Code:           d.BaseCode + 3,
		Prefix:         d.Prefix,
		Reason:         ExceptionReason_FailedToUpdate,
		Message:        message,
		HTTPStatusCode: http.StatusInternalServerError,
	}
}

func (d *APIExceptionDomain) FailedToDelete(optionalMessage ...string) *Exception {
	message := fmt.Sprintf("Failed to delete the %s", strings.ToLower(string(d.Prefix)))
	if len(optionalMessage) > 0 && len(strings.ReplaceAll(optionalMessage[0], " ", "")) > 0 {
		message = optionalMessage[0]
	}

	return &Exception{
		Code:           d.BaseCode + 4,
		Prefix:         d.Prefix,
		Reason:         ExceptionReason_FailedToDelete,
		Message:        message,
		HTTPStatusCode: http.StatusInternalServerError,
	}
}

func (d *APIExceptionDomain) Timeout(time time.Duration, optionalMessage ...string) *Exception {
	message := fmt.Sprintf("Timeout in %s with %v", strings.ToLower(string(d.Prefix)), time)
	if len(optionalMessage) > 0 && len(strings.ReplaceAll(optionalMessage[0], " ", "")) > 0 {
		message = optionalMessage[0]
	}

	return &Exception{
		Code:           d.BaseCode + 5,
		Prefix:         d.Prefix,
		Reason:         ExceptionReason_Timeout,
		Message:        message,
		HTTPStatusCode: http.StatusRequestTimeout,
	}
}
