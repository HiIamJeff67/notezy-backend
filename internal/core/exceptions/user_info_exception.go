package apiexceptions

type userInfoExceptionDomain struct {
	domainException
}

var UserInfo = userInfoExceptionDomain{
	domainException: newDomainException("UserInfo"),
}
