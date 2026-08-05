package corestrategies

import (
	"errors"
	"sync"
	"time"

	durablejobeventscontract "github.com/HiIamJeff67/notezy-backend/contracts/durablejob/v1/events"
	"github.com/google/uuid"
)

const (
	maximumPendingHints    = 1_000
	MaximumDispatchBatch   = 32
	MaximumDispatchWorkers = 8
	MaximumRequestAttempts = 3
)

type YjsMaintenanceStrategy struct {
	mutex    sync.Mutex
	pending  map[uuid.UUID]durablejobeventscontract.YjsMaintenanceHintData
	requests map[uuid.UUID]MaintenanceRequest
	attempts map[uuid.UUID]int
	inFlight map[uuid.UUID]struct{}
	notify   chan struct{}
}

type MaintenanceRequest struct {
	Hint    durablejobeventscontract.YjsMaintenanceHintData
	Attempt int
}

func NewYjsMaintenanceStrategy() *YjsMaintenanceStrategy {
	return &YjsMaintenanceStrategy{
		pending:  make(map[uuid.UUID]durablejobeventscontract.YjsMaintenanceHintData),
		requests: make(map[uuid.UUID]MaintenanceRequest),
		attempts: make(map[uuid.UUID]int),
		inFlight: make(map[uuid.UUID]struct{}),
		notify:   make(chan struct{}, 1),
	}
}

func (s *YjsMaintenanceStrategy) Enqueue(hint durablejobeventscontract.YjsMaintenanceHintData) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	if _, exists := s.pending[hint.BlockPackId]; !exists && len(s.pending) >= maximumPendingHints {
		return errors.New("Yjs maintenance strategy queue is full")
	}
	if current, exists := s.pending[hint.BlockPackId]; !exists || hint.LatestUpdateSequence >= current.LatestUpdateSequence {
		s.pending[hint.BlockPackId] = hint
	}
	select {
	case s.notify <- struct{}{}:
	default:
	}

	return nil
}

func (s *YjsMaintenanceStrategy) DequeueBatch(limit int) []durablejobeventscontract.YjsMaintenanceHintData {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	if len(s.pending) == 0 {
		return nil
	}
	if limit > len(s.pending) {
		limit = len(s.pending)
	}
	result := make([]durablejobeventscontract.YjsMaintenanceHintData, 0, limit)
	for len(result) < limit {
		var selectedId uuid.UUID
		var selected durablejobeventscontract.YjsMaintenanceHintData
		selectedScore := int64(-1)
		for blockPackId, hint := range s.pending {
			if _, exists := s.inFlight[blockPackId]; exists {
				continue
			}
			score := hint.UncompactedUpdateCount*4 + (hint.LatestUpdateSequence-hint.ProjectedUntilSequence)*3
			if hint.LastCompactedAt == nil {
				score += 100_000
			} else if age := time.Since(*hint.LastCompactedAt); age > 0 {
				score += int64(age / time.Minute)
			}
			score += int64((hint.SnapshotBytes + hint.StateVectorBytes) / 1024)
			if score > selectedScore {
				selectedId = blockPackId
				selected = hint
				selectedScore = score
			}
		}
		if selectedId == uuid.Nil {
			break
		}
		delete(s.pending, selectedId)
		result = append(result, selected)
	}

	return result
}

func (s *YjsMaintenanceStrategy) PendingCount() int {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	return len(s.pending)
}

func (s *YjsMaintenanceStrategy) Track(requestId uuid.UUID, hint durablejobeventscontract.YjsMaintenanceHintData) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	attempt := s.attempts[hint.BlockPackId] + 1
	s.attempts[hint.BlockPackId] = attempt
	s.requests[requestId] = MaintenanceRequest{Hint: hint, Attempt: attempt}
	s.inFlight[hint.BlockPackId] = struct{}{}
}

func (s *YjsMaintenanceStrategy) Complete(requestId uuid.UUID) (MaintenanceRequest, bool) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	request, exists := s.requests[requestId]
	if exists {
		delete(s.attempts, request.Hint.BlockPackId)
		delete(s.inFlight, request.Hint.BlockPackId)
	}
	delete(s.requests, requestId)
	s.signalPending(request.Hint.BlockPackId)

	return request, exists
}

func (s *YjsMaintenanceStrategy) Fail(requestId uuid.UUID) (MaintenanceRequest, bool) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	request, exists := s.requests[requestId]
	if exists {
		delete(s.requests, requestId)
		delete(s.inFlight, request.Hint.BlockPackId)
		if request.Attempt >= MaximumRequestAttempts {
			delete(s.attempts, request.Hint.BlockPackId)
		}
	}
	s.signalPending(request.Hint.BlockPackId)

	return request, exists
}

func (s *YjsMaintenanceStrategy) Notify() <-chan struct{} {
	return s.notify
}

func (s *YjsMaintenanceStrategy) signalPending(blockPackId uuid.UUID) {
	if _, exists := s.pending[blockPackId]; !exists {
		return
	}
	select {
	case s.notify <- struct{}{}:
	default:
	}
}
