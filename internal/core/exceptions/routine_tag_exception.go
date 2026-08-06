package apiexceptions

type routineTagExceptionDomain struct {
	domainException
}

var RoutineTag = routineTagExceptionDomain{
	domainException: newDomainException("RoutineTag"),
}
