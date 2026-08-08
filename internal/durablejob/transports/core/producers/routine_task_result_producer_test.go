package coreproducers

import (
	"context"
	"testing"

	durablejobroutinetask "github.com/HiIamJeff67/notezy-backend/internal/durablejob/routinetask"
)

func TestRoutineTaskResultProducerRejectsUnsupportedResultKind(t *testing.T) {
	producer := NewRoutineTaskResultProducer(nil)
	err := producer.Produce(context.Background(), durablejobroutinetask.RoutineTaskResult{
		Kind: durablejobroutinetask.RoutineTaskResultKind("unsupported"),
	})
	if err == nil {
		t.Fatal("expected unsupported result kind error")
	}
}
