package realtimetypes

import (
	"testing"

	"github.com/google/uuid"
)

func TestYjsCompactionBatchInputMarshalAndUnmarshalBytes(t *testing.T) {
	blockPackId := uuid.MustParse("6c6a5f1f-5f9f-4b05-b3c0-3ab7a3d7a4e0")
	input := YjsCompactionBatchInput{
		BlockPackId: blockPackId,
		Input: YjsCompactionInput{
			Snapshot:                   []byte{1},
			StateVector:                []byte{2},
			BaseCompactedUntilSequence: 0,
			CutoffSequence:             1,
			Updates: []YjsDocumentUpdate{{
				UpdateSequence: 1,
				Payload:        []byte{3},
			}},
		},
	}

	payload, err := input.MarshalBytes()
	if err != nil {
		t.Fatalf("MarshalBytes() error = %v", err)
	}

	var actual YjsCompactionBatchInput
	if err := actual.UnmarshalBytes(payload); err != nil {
		t.Fatalf("UnmarshalBytes() error = %v", err)
	}
	if actual.BlockPackId != blockPackId {
		t.Fatalf("block pack id = %s, want %s", actual.BlockPackId, blockPackId)
	}
	if actual.Input.CutoffSequence != input.Input.CutoffSequence {
		t.Fatalf("cutoff sequence = %d, want %d", actual.Input.CutoffSequence, input.Input.CutoffSequence)
	}
}

func TestYjsCompactionBatchInputMarshalBytesRejectsNilBlockPackId(t *testing.T) {
	input := YjsCompactionBatchInput{
		Input: YjsCompactionInput{
			BaseCompactedUntilSequence: 0,
			CutoffSequence:             1,
			Updates:                    []YjsDocumentUpdate{{UpdateSequence: 1}},
		},
	}

	if _, err := input.MarshalBytes(); err == nil {
		t.Fatal("MarshalBytes() error = nil, want nil block pack id error")
	}
}
