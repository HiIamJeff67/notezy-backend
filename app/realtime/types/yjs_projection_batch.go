package realtimetypes

import (
	"encoding/binary"
	"errors"
	"math"

	"github.com/google/uuid"
)

type YjsProjectionBatchInput struct {
	BlockPackId uuid.UUID
	State       YjsDocumentState
}

func (b YjsProjectionBatchInput) MarshalBytes() ([]byte, error) {
	if b.BlockPackId == uuid.Nil {
		return nil, errors.New("invalid yjs projection batch block pack id")
	}

	statePayload, err := b.State.MarshalBytes()
	if err != nil {
		return nil, err
	}
	if len(statePayload) > math.MaxUint32 {
		return nil, errors.New("yjs projection batch state exceeds the payload limit")
	}

	payload := make([]byte, 20+len(statePayload))
	copy(payload[:16], b.BlockPackId[:])
	binary.BigEndian.PutUint32(payload[16:20], uint32(len(statePayload)))
	copy(payload[20:], statePayload)

	return payload, nil
}

type YjsProjectionBatchResult struct {
	BlockPackId uuid.UUID
	Payload     []byte
}

func (b *YjsProjectionBatchResult) UnmarshalBytes(payload []byte) error {
	*b = YjsProjectionBatchResult{}
	if len(payload) < 20 {
		return errors.New("invalid yjs projection batch result")
	}

	copy(b.BlockPackId[:], payload[:16])
	if b.BlockPackId == uuid.Nil {
		return errors.New("invalid yjs projection batch result block pack id")
	}

	resultLength := binary.BigEndian.Uint32(payload[16:20])
	if uint64(resultLength) != uint64(len(payload)-20) {
		return errors.New("invalid yjs projection batch result size")
	}
	b.Payload = append(b.Payload, payload[20:]...)

	return nil
}
