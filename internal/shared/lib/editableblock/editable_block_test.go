package editableblock

import (
	"encoding/json"
	"testing"

	"github.com/google/uuid"
)

func TestFlattenManyToRaw(t *testing.T) {
	rootId := uuid.New()
	childId := uuid.New()
	blocks, totalSize, err := FlattenManyToRaw([]ArborizedBlock{
		{
			Id:      rootId,
			Type:    "paragraph",
			Props:   map[string]string{"textColor": "default"},
			Content: []string{"root"},
			Children: []ArborizedBlock{
				{
					Id:      childId,
					Type:    "paragraph",
					Props:   map[string]string{},
					Content: []string{"child"},
				},
			},
		},
	})
	if err != nil {
		t.Fatalf("flatten blocks: %v", err)
	}
	if len(blocks) != 2 || totalSize == 0 {
		t.Fatalf("unexpected flatten result: %#v, totalSize=%d", blocks, totalSize)
	}
	if blocks[1].ParentBlockId == nil || *blocks[1].ParentBlockId != rootId {
		t.Fatalf("expected child parent id %s, got %#v", rootId, blocks[1].ParentBlockId)
	}
	if !json.Valid(blocks[0].Props) || !json.Valid(blocks[0].Content) {
		t.Fatalf("expected valid JSON in raw block: %#v", blocks[0])
	}
}
