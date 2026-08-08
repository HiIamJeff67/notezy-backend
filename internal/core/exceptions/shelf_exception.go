package apiexceptions

import (
	"fmt"
	"net/http"

	exceptions "github.com/HiIamJeff67/notezy-backend/contracts/types/exceptions"
)

type shelfExceptionDomain struct {
	domainException
}

var Shelf = shelfExceptionDomain{
	domainException: newDomainException("Shelf"),
}

func (shelfExceptionDomain) DuplicateName(name string) *exceptions.Exception {
	return exceptions.New(
		"DuplicateName",
		"Shelf",
		"Create",
		fmt.Sprintf("The name of %s is already in use", name),
		http.StatusConflict,
	)
}

func (shelfExceptionDomain) MaximumDepthExceeded(currentDepth int32, maxDepth int32) *exceptions.Exception {
	return exceptions.New(
		"MaximumDepthExceeded",
		"Shelf",
		"Validate",
		fmt.Sprintf("The current depth of %d exceeds the maximum depth of %d", currentDepth, maxDepth),
		http.StatusBadRequest,
	)
}

func (shelfExceptionDomain) InsertParentIntoItsChildren(destination any, target any) *exceptions.Exception {
	return exceptions.New(
		"InsertParentIntoItsChildren",
		"Shelf",
		"Validate",
		fmt.Sprintf("Cannot insert parent %v into child %v", target, destination),
		http.StatusBadRequest,
	)
}
