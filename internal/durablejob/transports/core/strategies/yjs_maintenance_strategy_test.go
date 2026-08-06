package corestrategies

import (
	"testing"
	"time"

	"github.com/google/uuid"

	coreeventscontract "github.com/HiIamJeff67/notezy-backend/contracts/core/v1/events"
)

func TestYjsMaintenanceStrategyCoalescesHintsByBlockPack(t *testing.T) {
	strategy := NewYjsMaintenanceStrategy()
	blockPackId := uuid.New()
	documentId := uuid.New()
	older := time.Now().UTC().Add(-time.Minute)
	newer := time.Now().UTC()

	if err := strategy.Enqueue(coreeventscontract.YjsMaintenanceHintData{
		BlockPackId: blockPackId, DocumentId: documentId, LatestUpdateSequence: 4, LastCompactedAt: &older,
	}); err != nil {
		t.Fatal(err)
	}
	if err := strategy.Enqueue(coreeventscontract.YjsMaintenanceHintData{
		BlockPackId: blockPackId, DocumentId: documentId, LatestUpdateSequence: 9, LastCompactedAt: &newer,
	}); err != nil {
		t.Fatal(err)
	}

	hints := strategy.DequeueBatch(10)
	if len(hints) != 1 || hints[0].LatestUpdateSequence != 9 {
		t.Fatalf("expected one latest coalesced hint, got %#v", hints)
	}
}

func TestYjsMaintenanceStrategyPrioritizesMaintenanceLag(t *testing.T) {
	strategy := NewYjsMaintenanceStrategy()
	firstBlockPackId := uuid.New()
	secondBlockPackId := uuid.New()
	if err := strategy.Enqueue(coreeventscontract.YjsMaintenanceHintData{
		BlockPackId: firstBlockPackId, DocumentId: uuid.New(), LatestUpdateSequence: 20,
		ProjectedUntilSequence: 1, UncompactedUpdateCount: 1,
	}); err != nil {
		t.Fatal(err)
	}
	if err := strategy.Enqueue(coreeventscontract.YjsMaintenanceHintData{
		BlockPackId: secondBlockPackId, DocumentId: uuid.New(), LatestUpdateSequence: 8,
		ProjectedUntilSequence: 0, UncompactedUpdateCount: 10,
	}); err != nil {
		t.Fatal(err)
	}

	hints := strategy.DequeueBatch(1)
	if len(hints) != 1 || hints[0].BlockPackId != secondBlockPackId {
		t.Fatalf("expected high uncompacted update count to win, got %#v", hints)
	}
}

func TestYjsMaintenanceStrategyPreservesBlockPackOrdering(t *testing.T) {
	strategy := NewYjsMaintenanceStrategy()
	blockPackId := uuid.New()
	firstHint := coreeventscontract.YjsMaintenanceHintData{
		BlockPackId: blockPackId, DocumentId: uuid.New(), LatestUpdateSequence: 3,
	}
	secondHint := firstHint
	secondHint.LatestUpdateSequence = 4
	requestId := uuid.New()

	strategy.Track(requestId, firstHint)
	if err := strategy.Enqueue(secondHint); err != nil {
		t.Fatal(err)
	}
	if hints := strategy.DequeueBatch(1); len(hints) != 0 {
		t.Fatalf("expected active BlockPack to remain queued, got %d hints", len(hints))
	}

	strategy.Complete(requestId)
	hints := strategy.DequeueBatch(1)
	if len(hints) != 1 || hints[0].LatestUpdateSequence != 4 {
		t.Fatalf("expected latest hint after completion, got %#v", hints)
	}
}
