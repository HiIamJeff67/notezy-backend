package realtimetypes

import (
	"encoding/binary"
	"errors"

	"github.com/google/uuid"
)

type YjsCompactionBatchInput struct {
	BlockPackId uuid.UUID
	Input       YjsCompactionInput
}

func (b YjsCompactionBatchInput) MarshalBytes() ([]byte, error) {
	if b.BlockPackId == uuid.Nil {
		return nil, errors.New("invalid yjs compaction batch block pack id")
	}

	inputPayload, err := b.Input.MarshalBytes()
	if err != nil {
		return nil, err
	}

	payload := make([]byte, 20+len(inputPayload))
	copy(payload[0:16], b.BlockPackId[:])
	binary.BigEndian.PutUint32(payload[16:20], uint32(len(inputPayload)))
	copy(payload[20:], inputPayload)

	return payload, nil
}

func (b *YjsCompactionBatchInput) UnmarshalBytes(payload []byte) error {
	*b = YjsCompactionBatchInput{}
	if len(payload) < 20 {
		return errors.New("invalid yjs compaction batch input")
	}

	copy(b.BlockPackId[:], payload[0:16])
	if b.BlockPackId == uuid.Nil {
		return errors.New("invalid yjs compaction batch block pack id")
	}

	inputLength := binary.BigEndian.Uint32(payload[16:20])
	if uint64(inputLength) != uint64(len(payload)-20) {
		return errors.New("invalid yjs compaction batch input size")
	}

	return b.Input.UnmarshalBytes(payload[20:])
}

type YjsCompactionBatchResult struct {
	BlockPackId uuid.UUID
	Result      YjsCompactionResult
}

func (b YjsCompactionBatchResult) MarshalBytes() ([]byte, error) {
	if b.BlockPackId == uuid.Nil {
		return nil, errors.New("invalid yjs compaction batch block pack id")
	}

	resultPayload, err := b.Result.MarshalBytes()
	if err != nil {
		return nil, err
	}

	payload := make([]byte, 20+len(resultPayload))
	copy(payload[0:16], b.BlockPackId[:])
	binary.BigEndian.PutUint32(payload[16:20], uint32(len(resultPayload)))
	copy(payload[20:], resultPayload)

	return payload, nil
}

func (b *YjsCompactionBatchResult) UnmarshalBytes(payload []byte) error {
	*b = YjsCompactionBatchResult{}
	if len(payload) < 20 {
		return errors.New("invalid yjs compaction batch result")
	}

	copy(b.BlockPackId[:], payload[0:16])
	if b.BlockPackId == uuid.Nil {
		return errors.New("invalid yjs compaction batch block pack id")
	}

	resultLength := binary.BigEndian.Uint32(payload[16:20])
	if uint64(resultLength) != uint64(len(payload)-20) {
		return errors.New("invalid yjs compaction batch result size")
	}

	return b.Result.UnmarshalBytes(payload[20:])
}
