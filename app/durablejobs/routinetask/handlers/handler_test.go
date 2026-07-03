package handlers

import (
	"testing"

	adapters "github.com/HiIamJeff67/notezy-backend/app/adapters"
	dtos "github.com/HiIamJeff67/notezy-backend/app/dtos"

	"github.com/google/uuid"
)

func TestFlattenArborizedBlockRejectsZeroBlock(t *testing.T) {
	_, _, _, exception := flattenArborizedBlock(
		adapters.NewEditableBlockAdapter(),
		uuid.New(),
		&dtos.ArborizedEditableBlock{},
	)
	if exception == nil {
		t.Fatal("expected invalid zero arborized block to be rejected")
	}
}
