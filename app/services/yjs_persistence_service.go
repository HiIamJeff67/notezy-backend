package services

import (
	"context"

	"github.com/google/uuid"
	"gorm.io/gorm"

	inputs "github.com/HiIamJeff67/notezy-backend/app/models/inputs"
	repositories "github.com/HiIamJeff67/notezy-backend/app/models/repositories"
	options "github.com/HiIamJeff67/notezy-backend/app/options"
	realtimetypes "github.com/HiIamJeff67/notezy-backend/app/realtime/types"
)

type YjsPersistenceServiceInterface interface {
	LoadDocument(ctx context.Context, blockPackId uuid.UUID) (*realtimetypes.YjsDocumentState, error)
	AppendUpdate(ctx context.Context, blockPackId uuid.UUID, persistenceBatchId uuid.UUID, originConnectionId *uuid.UUID, payload []byte) (int64, error)
	GetCompactableYjsDocumentWithUpdates(ctx context.Context, blockPackId uuid.UUID) (*realtimetypes.YjsCompactionInput, error)
	ApplyCompactedYjsDocument(ctx context.Context, blockPackId uuid.UUID, result realtimetypes.YjsCompactionResult) (bool, error)
}

type YjsPersistenceService struct {
	db                     *gorm.DB
	blockPackYjsRepository repositories.BlockPackYjsRepositoryInterface
}

func NewYjsPersistenceService(db *gorm.DB) YjsPersistenceServiceInterface {
	return &YjsPersistenceService{
		db:                     db,
		blockPackYjsRepository: repositories.NewBlockPackYjsRepository(),
	}
}

func (s *YjsPersistenceService) LoadDocument(
	ctx context.Context, blockPackId uuid.UUID,
) (*realtimetypes.YjsDocumentState, error) {
	db := s.db.WithContext(ctx)

	document, updates, err := s.blockPackYjsRepository.LoadDocumentAndUpdates(
		blockPackId,
		options.WithDB(db),
	)
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
	ctx context.Context,
	blockPackId uuid.UUID,
	persistenceBatchId uuid.UUID,
	originConnectionId *uuid.UUID,
	payload []byte,
) (int64, error) {
	db := s.db.WithContext(ctx)

	return s.blockPackYjsRepository.AppendUpdate(
		blockPackId,
		inputs.AppendBlockPackYjsUpdateInput{
			PersistenceBatchId: persistenceBatchId,
			OriginConnectionId: originConnectionId,
			Payload:            payload,
		},
		options.WithDB(db),
	)
}

func (s *YjsPersistenceService) GetCompactableYjsDocumentWithUpdates(
	ctx context.Context, blockPackId uuid.UUID,
) (*realtimetypes.YjsCompactionInput, error) {
	db := s.db.WithContext(ctx)

	document, updates, err := s.blockPackYjsRepository.GetCompactableYjsDocumentWithUpdates(
		blockPackId,
		options.WithDB(db),
	)
	if err != nil || document == nil {
		return nil, err
	}

	input := &realtimetypes.YjsCompactionInput{
		Snapshot:                   document.Snapshot,
		StateVector:                document.StateVector,
		BaseCompactedUntilSequence: document.CompactedUntilSequence,
		CutoffSequence:             document.LastUpdateSequence,
		Updates:                    make([]realtimetypes.YjsDocumentUpdate, len(updates)),
	}
	for index, update := range updates {
		input.Updates[index] = realtimetypes.YjsDocumentUpdate{
			UpdateSequence: update.UpdateSequence,
			Payload:        update.Payload,
		}
	}

	return input, nil
}

func (s *YjsPersistenceService) ApplyCompactedYjsDocument(
	ctx context.Context, blockPackId uuid.UUID, result realtimetypes.YjsCompactionResult,
) (bool, error) {
	db := s.db.WithContext(ctx)

	return s.blockPackYjsRepository.ApplyCompactedYjsDocument(
		blockPackId,
		inputs.ApplyCompactedBlockPackYjsDocumentInput{
			BaseCompactedUntilSequence: result.BaseCompactedUntilSequence,
			CutoffSequence:             result.CutoffSequence,
			Snapshot:                   result.Snapshot,
			StateVector:                result.StateVector,
		},
		options.WithDB(db),
	)
}
