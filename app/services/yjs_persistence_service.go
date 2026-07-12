package services

import (
	"context"

	"github.com/google/uuid"
	"gorm.io/gorm"

	repositories "github.com/HiIamJeff67/notezy-backend/app/models/repositories"
	realtimetypes "github.com/HiIamJeff67/notezy-backend/app/realtime/types"
)

type YjsPersistenceServiceInterface interface {
	LoadDocument(ctx context.Context, blockPackId uuid.UUID) (*realtimetypes.YjsDocumentState, error)
	AppendUpdate(ctx context.Context, blockPackId uuid.UUID, originConnectionId uuid.UUID, payload []byte) (int64, error)
}

type YjsPersistenceService struct {
	blockPackYjsRepository repositories.BlockPackYjsRepositoryInterface
}

func NewYjsPersistenceService(db *gorm.DB) YjsPersistenceServiceInterface {
	return &YjsPersistenceService{
		blockPackYjsRepository: repositories.NewBlockPackYjsRepository(db),
	}
}

func (s *YjsPersistenceService) LoadDocument(
	ctx context.Context, blockPackId uuid.UUID,
) (*realtimetypes.YjsDocumentState, error) {
	document, updates, err := s.blockPackYjsRepository.LoadDocumentAndUpdates(ctx, blockPackId)
	if err != nil {
		return nil, err
	}

	state := &realtimetypes.YjsDocumentState{
		Snapshot:               document.Snapshot,
		StateVector:            document.StateVector,
		LastUpdateSequence:     document.LastUpdateSequence,
		CompactedUntilSequence: document.CompactedUntilSequence,
		ProjectedUntilSequence: document.ProjectedUntilSequence,
		Updates:                make([]realtimetypes.YjsDocumentUpdate, len(updates)),
	}
	for index, update := range updates {
		state.Updates[index] = realtimetypes.YjsDocumentUpdate{
			UpdateSequence: update.UpdateSequence,
			Payload:        update.Payload,
		}
	}

	return state, nil
}

func (s *YjsPersistenceService) AppendUpdate(
	ctx context.Context, blockPackId uuid.UUID, originConnectionId uuid.UUID, payload []byte,
) (int64, error) {
	return s.blockPackYjsRepository.AppendUpdate(ctx, blockPackId, originConnectionId, payload)
}
