package durablejobexceptions

type stationExceptionDomain struct {
	domainException
}

var Station = stationExceptionDomain{
	domainException: newDomainException("Station"),
}
