package durablejobexceptions

type routineTaskExceptionDomain struct {
	domainException
}

var RoutineTask = routineTaskExceptionDomain{
	domainException: newDomainException("RoutineTask"),
}
