package apiexceptions

type materialExceptionDomain struct {
	domainException
}

var Material = materialExceptionDomain{
	domainException: newDomainException("Material"),
}
