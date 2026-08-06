package apiexceptions

import (
	"fmt"
	"net/http"
	"time"

	exceptions "github.com/HiIamJeff67/notezy-backend/shared/exceptions"
)

type routineExceptionDomain struct {
	domainException
}

var Routine = routineExceptionDomain{
	domainException: newDomainException("Routine"),
}

func (routineExceptionDomain) QueriedTimeRangeTooLarge(from time.Time, to time.Time) *exceptions.Exception {
	return exceptions.New(
		"QueriedTimeRangeTooLarge",
		"Routine",
		"Search",
		fmt.Sprintf("Cannot query the time range from %s to %s because it is too large", from, to),
		http.StatusBadRequest,
	)
}

func (routineExceptionDomain) FailedToLinkRoutineTags() *exceptions.Exception {
	return exceptions.New(
		"FailedToLinkRoutineTags",
		"Routine",
		"Link",
		"Cannot link the given routine tags to the target routine",
		http.StatusBadRequest,
	)
}

func (routineExceptionDomain) FailedToLinkItems() *exceptions.Exception {
	return exceptions.New(
		"FailedToLinkItems",
		"Routine",
		"Link",
		"Cannot link the given items to the target routine",
		http.StatusBadRequest,
	)
}
