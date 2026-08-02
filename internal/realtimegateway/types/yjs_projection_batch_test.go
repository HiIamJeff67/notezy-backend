package realtimetypes

import (
	"bytes"
	"testing"

	"github.com/google/uuid"
)

func TestYjsProjectionBatchInputMarshalBytes(t *testing.T) {
	blockPackId := uuid.New()
	expected := YjsProjectionBatchInput{
		BlockPackId: blockPackId,
		State: YjsDocumentState{
			Snapshot:               []byte{1, 2},
			StateVector:            []byte{3},
			LastUpdateSequence:     2,
			CompactedUntilSequence: 1,
			ProjectedUntilSequence: 0,
			Updates: []YjsDocumentUpdate{
				{UpdateSequence: 2, Payload: []byte{4, 5}},
			},
		},
	}

	payload, err := expected.MarshalBytes()
	if err != nil {
		t.Fatalf("MarshalBytes() error = %v", err)
	}
	if !bytes.Equal(payload[:16], blockPackId[:]) {
		t.Fatalf("MarshalBytes() block pack id does not match")
	}
}

func TestYjsProjectionBatchResultUnmarshalBytes(t *testing.T) {
	blockPackId := uuid.New()
	payload := append(append([]byte{}, blockPackId[:]...), 0, 0, 0, 3, '{', '}', '\n')

	var result YjsProjectionBatchResult
	if err := result.UnmarshalBytes(payload); err != nil {
		t.Fatalf("UnmarshalBytes() error = %v", err)
	}
	if result.BlockPackId != blockPackId {
		t.Fatalf("UnmarshalBytes() block pack id = %s, want %s", result.BlockPackId, blockPackId)
	}
	if string(result.Payload) != "{}\n" {
		t.Fatalf("UnmarshalBytes() payload = %q, want %q", result.Payload, "{}\n")
	}
}
