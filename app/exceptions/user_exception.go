package exceptions

import (
	"fmt"
	"net/http"
)

const (
	// define this bcs the general domain has some general exception that has be defined

	_ExceptionBaseCode_User ExceptionCode = UserExceptionSubDomainCode * ExceptionSubDomainCodeShiftAmount // the actual base for the exceptions of user

	// if you need to code a custom exception of users,
	// use the ExceptionBaseCode_User, instead of _ExceptionBaseCode_User
	// the exception codes that we can actually customize here is shifted with ReservedExceptionCode

	UserExceptionSubDomainCode ExceptionCode   = 32
	ExceptionBaseCode_User     ExceptionCode   = _ExceptionBaseCode_User + ReservedExceptionCode
	ExceptionPrefix_User       ExceptionPrefix = "User"
)

type UserExceptionDomain struct {
	BaseCode ExceptionCode
	Prefix   ExceptionPrefix
	// as the down layer of DatabaseExceptionDomain
	// so that we don't make methods for DatabaseExceptionDomain
	// instead we make methods for UserExceptionDomain
	DatabaseExceptionDomain
	APIExceptionDomain
	GraphQLExceptionDomain
	TypeExceptionDomain
}

var User = &UserExceptionDomain{
	BaseCode: ExceptionBaseCode_User,
	Prefix:   ExceptionPrefix_User,
	DatabaseExceptionDomain: DatabaseExceptionDomain{
		_BaseCode: _ExceptionBaseCode_User,
		_Prefix:   ExceptionPrefix_User,
	},
	APIExceptionDomain: APIExceptionDomain{
		_BaseCode: _ExceptionBaseCode_User,
		_Prefix:   ExceptionPrefix_User,
	},
	GraphQLExceptionDomain: GraphQLExceptionDomain{
		_BaseCode: _ExceptionBaseCode_User,
		_Prefix:   ExceptionPrefix_User,
	},
	TypeExceptionDomain: TypeExceptionDomain{
		_BaseCode: _ExceptionBaseCode_User,
		_Prefix:   ExceptionPrefix_User,
	},
}

/* ============================== Unique Constaints ============================== */

func (d *UserExceptionDomain) DuplicateName(name string) *Exception {
	return &Exception{
		Code:           d.BaseCode + 1,
		Prefix:         d.Prefix,
		Reason:         "DuplicateName",
		IsInternal:     false,
		Message:        fmt.Sprintf("The name of %s is already be used", name),
		HTTPStatusCode: http.StatusConflict,
		LastStackFrame: &GetStackTrace(2, 1)[0],
	}
}

func (d *UserExceptionDomain) DuplicateEmail(email string) *Exception {
	return &Exception{
		Code:           d.BaseCode + 2,
		Prefix:         d.Prefix,
		Reason:         "DuplicateEmail",
		IsInternal:     false,
		Message:        fmt.Sprintf("The email of %s is already be used", email),
		HTTPStatusCode: http.StatusConflict,
		LastStackFrame: &GetStackTrace(2, 1)[0],
	}
}
