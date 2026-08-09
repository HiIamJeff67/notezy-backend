package apiexceptions

import (
	"fmt"
	"net/http"

	exceptions "github.com/HiIamJeff67/notezy-backend/contracts/types/exceptions"
)

type ShelfException struct {
	CoreException
}

func NewShelfException() ShelfException {
	return ShelfException{
		CoreException: NewCoreException("Shelf"),
	}
}

func (ShelfException) DuplicateName(name string) *exceptions.Exception {
	return exceptions.New(
		"DuplicateName",
		"Shelf",
		"Create",
		fmt.Sprintf("The name of %s is already in use", name),
		http.StatusConflict,
	)
}

func (ShelfException) MaximumDepthExceeded(currentDepth int32, maxDepth int32) *exceptions.Exception {
	return exceptions.New(
		"MaximumDepthExceeded",
		"Shelf",
		"Validate",
		fmt.Sprintf("The current depth of %d exceeds the maximum depth of %d", currentDepth, maxDepth),
		http.StatusBadRequest,
	)
}

func (ShelfException) InsertParentIntoItsChildren(destination any, target any) *exceptions.Exception {
	return exceptions.New(
		"InsertParentIntoItsChildren",
		"Shelf",
		"Validate",
		fmt.Sprintf("Cannot insert parent %v into child %v", target, destination),
		http.StatusBadRequest,
	)
}
