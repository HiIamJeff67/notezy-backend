package exceptions

import (
	"fmt"
	"go-gorm-api/app/logs"
	"net/http"
)

/* ============================== Exception Field Type Definition ============================== */
type ExceptionCode int
type ExceptionPrefix string
type ExceptionReason string

const (
	// the first 3 digits are the class of exceptions
	// the last 5 digits are the individual labels for each exceptions
	MaxExceptionCode = 99999999	// 999 99999
	MinExceptionCode = 0 		// 000 00000
	// reserve some codes for general use purpose
	// see the below general exceptions ex. NotFound(), FailedToCreate()
	ReservedExceptionCode = 100	// *** **100, the codes >= *** **100 will be use in the general domain
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
)

// global reason for common domain use
// if some individual domain require a custom reason, 
// just create one with ExceptionReason type privately whic means its variable name in lower case
const (
	ExceptionNotFound ExceptionReason = "Not_Found"
	ExceptionFailedToCreate ExceptionReason = "Failed_To_Create"
	ExceptionFailedToUpdate ExceptionReason = "Failed_To_Update"
	ExceptionFailedToDelete ExceptionReason = "Failed_To_Delete"
)

func IsExceptionCode(exceptionCode int) bool {
	return exceptionCode >= MinExceptionCode && exceptionCode <= MaxExceptionCode
}
/* ============================== Exception Field Type Definition ============================== */

/* ============================== General Exception Structure Definition ============================== */
type Exception struct {
	Code			ExceptionCode	// custom exception code
	Prefix  		ExceptionPrefix // custom exception prefix
	Reason 			ExceptionReason	// custom exception reason
	Message 		string			// custom exception message
	HTTPStatusCode 	int				// http status code
	Detials 		any				// additional error details (optional)
	Error			error			// original error (optional)
}

func (e *Exception) GetString() string {
	if e.Error != nil {
		return fmt.Sprintf("[%v]%s: %v", e.Code, e.Reason, e.Error)
	}
	return fmt.Sprintf("[%v]%s: %s", e.Code, e.Reason, e.Message)
}

func (e *Exception) WithDetials(detials any) *Exception {
	e.Detials = detials
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
/* ============================== General Exception Structure Definition ============================== */

/* ============================== Exception Domain Definition ============================== */
type APIExceptionDomain struct {
	BaseCode ExceptionCode
	Prefix ExceptionPrefix
}

func (d *APIExceptionDomain) NotFound() *Exception {
	return &Exception{
		Code: d.BaseCode + 1, 
		Prefix: d.Prefix, 
		Reason: ExceptionNotFound,
		Message: fmt.Sprintf("%s not found", d.Prefix), 
		HTTPStatusCode: http.StatusNotFound,
	}
}

func (d *APIExceptionDomain) FailedToCreate() *Exception {
	return &Exception{
		Code: d.BaseCode + 2, 
		Prefix: d.Prefix, 
		Reason: ExceptionFailedToCreate, 
		Message: fmt.Sprintf("Failed to create the %s", d.Prefix), 
		HTTPStatusCode: http.StatusBadRequest,
	}
}

func (d *APIExceptionDomain) FailedToUpdate() *Exception {
	return &Exception{
		Code: d.BaseCode + 3, 
		Prefix: d.Prefix, 
		Reason: ExceptionFailedToUpdate, 
		Message: fmt.Sprintf("Failed to update the %s", d.Prefix), 
		HTTPStatusCode: http.StatusBadRequest,
	}
}

func (d *APIExceptionDomain) FailedToDelete() *Exception {
	return &Exception{
		Code: d.BaseCode + 4, 
		Prefix: d.Prefix, 
		Reason: ExceptionFailedToDelete, 
		Message: fmt.Sprintf("Failed to delete the %s", d.Prefix), 
		HTTPStatusCode: http.StatusBadRequest,
	}
}
/* ============================== Exception Domain Definition ============================== */