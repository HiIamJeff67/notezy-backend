package adapters

import (
	"encoding/json"
	"testing"

	"github.com/google/uuid"

	dtos "github.com/HiIamJeff67/notezy-backend/app/dtos"
)

func TestEditableBlockAdapterFlattensForestWithSiblingPointers(t *testing.T) {
	firstRootId := uuid.New()
	secondRootId := uuid.New()
	firstChildId := uuid.New()
	secondChildId := uuid.New()

	payload := []byte(`[
		{
			"id": "` + firstRootId.String() + `",
			"type": "paragraph",
			"props": {"backgroundColor":"default","textColor":"default","textAlignment":"left"},
			"content": [],
			"children": [
				{
					"id": "` + firstChildId.String() + `",
					"type": "paragraph",
					"props": {"backgroundColor":"default","textColor":"default","textAlignment":"left"},
					"content": [],
					"children": []
				},
				{
					"id": "` + secondChildId.String() + `",
					"type": "paragraph",
					"props": {"backgroundColor":"default","textColor":"default","textAlignment":"left"},
					"content": [],
					"children": []
				}
			]
		},
		{
			"id": "` + secondRootId.String() + `",
			"type": "paragraph",
			"props": {"backgroundColor":"default","textColor":"default","textAlignment":"left"},
			"content": [],
			"children": []
		}
	]`)

	var roots []dtos.ArborizedEditableBlock
	if err := json.Unmarshal(payload, &roots); err != nil {
		t.Fatalf("failed to unmarshal arborized blocks: %v", err)
	}

	blocks, _, exception := NewEditableBlockAdapter().FlattenManyToRaw(roots)
	if exception != nil {
		t.Fatalf("failed to flatten arborized block forest: %v", exception)
	}
	if len(blocks) != 4 {
		t.Fatalf("expected 4 flattened blocks, got %d", len(blocks))
	}

	blocksById := make(map[uuid.UUID]dtos.RawFlattenedEditableBlock, len(blocks))
	for _, block := range blocks {
		blocksById[block.Id] = block
	}

	firstRoot := blocksById[firstRootId]
	if firstRoot.ParentBlockId != nil || firstRoot.PrevBlockId != nil ||
		firstRoot.NextBlockId == nil || *firstRoot.NextBlockId != secondRootId {
		t.Fatalf("unexpected first root pointers: %#v", firstRoot)
	}

	secondRoot := blocksById[secondRootId]
	if secondRoot.ParentBlockId != nil || secondRoot.NextBlockId != nil ||
		secondRoot.PrevBlockId == nil || *secondRoot.PrevBlockId != firstRootId {
		t.Fatalf("unexpected second root pointers: %#v", secondRoot)
	}

	firstChild := blocksById[firstChildId]
	if firstChild.ParentBlockId == nil || *firstChild.ParentBlockId != firstRootId ||
		firstChild.PrevBlockId != nil || firstChild.NextBlockId == nil ||
		*firstChild.NextBlockId != secondChildId {
		t.Fatalf("unexpected first child pointers: %#v", firstChild)
	}

	secondChild := blocksById[secondChildId]
	if secondChild.ParentBlockId == nil || *secondChild.ParentBlockId != firstRootId ||
		secondChild.NextBlockId != nil || secondChild.PrevBlockId == nil ||
		*secondChild.PrevBlockId != firstChildId {
		t.Fatalf("unexpected second child pointers: %#v", secondChild)
	}
}
