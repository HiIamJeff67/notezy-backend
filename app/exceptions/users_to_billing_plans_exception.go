package exceptions

const (
	_ExceptionBaseCode_UsersToBillingPlans ExceptionCode = UsersToBillingPlansExceptionSubDomainCode * ExceptionSubDomainCodeShiftAmount

	UsersToBillingPlansExceptionSubDomainCode ExceptionCode   = 46
	ExceptionBaseCode_UsersToBillingPlans     ExceptionCode   = _ExceptionBaseCode_UsersToBillingPlans + ReservedExceptionCode
	ExceptionPrefix_UsersToBillingPlans       ExceptionPrefix = "UsersToBillingPlans"
)

type UsersToBillingPlansExceptionDomain struct {
	BaseCode ExceptionCode
	Prefix   ExceptionPrefix
	DatabaseExceptionDomain
}

var UsersToBillingPlans = &UsersToBillingPlansExceptionDomain{
	BaseCode: ExceptionBaseCode_UsersToBillingPlans,
	Prefix:   ExceptionPrefix_UsersToBillingPlans,
	DatabaseExceptionDomain: DatabaseExceptionDomain{
		_BaseCode: _ExceptionBaseCode_UsersToBillingPlans,
		_Prefix:   ExceptionPrefix_UsersToBillingPlans,
	},
}
