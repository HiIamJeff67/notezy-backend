package apiexceptions

import (
	"fmt"
	"net/http"
	"time"

	exceptions "github.com/HiIamJeff67/notezy-backend/contracts/types/exceptions"
)

type RoutineException struct {
	CoreException
}

func NewRoutineException() RoutineException {
	return RoutineException{
		CoreException: NewCoreException("Routine"),
	}
}

func (RoutineException) QueriedTimeRangeTooLarge(from time.Time, to time.Time) *exceptions.Exception {
	return exceptions.New(
		"QueriedTimeRangeTooLarge",
		"Routine",
		"Search",
		fmt.Sprintf("Cannot query the time range from %s to %s because it is too large", from, to),
		http.StatusBadRequest,
	)
}

func (RoutineException) FailedToLinkRoutineTags() *exceptions.Exception {
	return exceptions.New(
		"FailedToLinkRoutineTags",
		"Routine",
		"Link",
		"Cannot link the given routine tags to the target routine",
		http.StatusBadRequest,
	)
}

func (RoutineException) FailedToLinkItems() *exceptions.Exception {
	return exceptions.New(
		"FailedToLinkItems",
		"Routine",
		"Link",
		"Cannot link the given items to the target routine",
		http.StatusBadRequest,
	)
}
